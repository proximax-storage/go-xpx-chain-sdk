package health

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/tools/health/packets"
	crypto "github.com/proximax-storage/go-xpx-crypto"
)

const DefaultAvgSecondsPerBlock = 15 * time.Second

var (
	ErrReturnedZeroHashes   = errors.New("returned zero hashes")
	ErrTimedOut             = errors.New("timed out")
	ErrTooLargeWaitTime     = errors.New("too large wait time")
	ErrNodeGotStuck         = errors.New("node got stuck")
	ErrNodeNotReachedHeight = errors.New("node does not reached the height")
)

type (
	NodeInfo struct {
		IdentityKey  *crypto.PublicKey
		Endpoint     string
		FriendlyName string
	}

	NodeHealthChecker struct {
		nodeInfo *NodeInfo

		handler *Handler
	}

	ChainScore struct {
		High uint64
		Low  uint64
	}

	ChainInfo struct {
		Height     uint64
		ChainScore ChainScore
	}
)

func NewNodeInfo(pKey, addr, friendlyName string) (*NodeInfo, error) {
	k, err := crypto.NewPublicKeyfromHex(pKey)
	if err != nil {
		return nil, err
	}

	return &NodeInfo{
		IdentityKey:  k,
		Endpoint:     addr,
		FriendlyName: friendlyName,
	}, nil
}

func (ni *NodeInfo) String() string {
	return fmt.Sprintf("%s=%s", ni.Endpoint, ni.IdentityKey)
}

func NewNodeHealthChecker(client *crypto.KeyPair, info *NodeInfo, mode packets.ConnectionSecurityMode) (*NodeHealthChecker, error) {
	nodeIo, err := NewNodeTcpIo(info)
	if err != nil {
		return nil, err
	}

	nhc := &NodeHealthChecker{
		nodeInfo: info,
		handler:  NewHandler(nodeIo),
	}

	err = nhc.handler.AuthHandle(client, info.IdentityKey, mode)
	if err != nil {
		nhc.Close()
		return nil, err
	}

	return nhc, nil
}

func (nhc *NodeHealthChecker) Close() error {
	return nhc.handler.Close()
}

func (nhc *NodeHealthChecker) ChainInfo() (*ChainInfo, error) {
	chainInfoReq := packets.NewPacketHeader(packets.ChainInfoPacketType)
	chainInfoResp := &packets.ChainInfoResponse{}
	err := nhc.handler.CommonHandle(&chainInfoReq, chainInfoResp)
	if err != nil {
		return nil, err
	}

	return &ChainInfo{
		Height: chainInfoResp.Height,
		ChainScore: ChainScore{
			High: chainInfoResp.ScoreHigh,
			Low:  chainInfoResp.ScoreLow,
		},
	}, nil
}

func (nhc *NodeHealthChecker) LastBlockHash() (sdk.Hash, error) {
	ci, err := nhc.ChainInfo()
	if err != nil {
		return sdk.Hash{}, err
	}

	return nhc.BlockHash(ci.Height)
}

func (nhc *NodeHealthChecker) BlockHash(height uint64) (sdk.Hash, error) {
	ci, err := nhc.ChainInfo()
	if err != nil {
		return sdk.Hash{}, err
	}

	if ci.Height < height {
		return sdk.Hash{}, ErrNodeNotReachedHeight
	}

	blockHashesReq := packets.NewBlockHashesRequest(height, 1)
	blockHashesResp := &packets.BlockHashesResponse{}
	err = nhc.handler.CommonHandle(blockHashesReq, blockHashesResp)
	if err != nil {
		return sdk.Hash{}, err
	}

	if len(blockHashesResp.Hashes) == 0 {
		return sdk.Hash{}, ErrReturnedZeroHashes
	}

	return blockHashesResp.Hashes[0], nil
}

func (nhc *NodeHealthChecker) NodeList() ([]*NodeInfo, error) {
	ndReq := packets.NewPacketHeader(packets.NodeDiscoveryPullPeersPacketType)
	ndResp := &packets.NodeDiscoveryPullPeersResponse{}
	err := nhc.handler.CommonHandle(&ndReq, ndResp)
	if err != nil {
		return nil, err
	}

	ni := make([]*NodeInfo, 0, len(ndResp.NetworkNodes))
	for _, node := range ndResp.NetworkNodes {
		if node.Host == "" || node.Port == 0 {
			continue
		}

		ni = append(ni, &NodeInfo{
			IdentityKey:  node.IdentityKey,
			Endpoint:     node.Host + ":" + strconv.Itoa(int(node.Port)),
			FriendlyName: node.FriendlyName,
		})
	}

	return ni, nil
}

// WaitHeight waits when a node will reach the expectedHeight
// In error case returns the last reached height
func (nhc *NodeHealthChecker) WaitHeight(expectedHeight uint64) (uint64, error) {
	ci, err := nhc.ChainInfo()
	if err != nil {
		return 0, err
	}

	if ci.Height >= expectedHeight {
		return ci.Height, err
	}
	prevHeight := ci.Height

	multiplier := expectedHeight - ci.Height
	if multiplier > 5 {
		multiplier = 5
	}

	avgSecondsPerBlock := DefaultAvgSecondsPerBlock
	tickerDuration := DefaultAvgSecondsPerBlock * time.Duration(multiplier)
	periodicTicker := time.NewTicker(tickerDuration)
	defer periodicTicker.Stop()

	log.Printf("Start waiting for node %s=%s to reach height: %d, current: %d\n", nhc.nodeInfo.Endpoint, nhc.nodeInfo.IdentityKey, expectedHeight, ci.Height)
	maxRetryCount := uint8(3)
	updateTicker := false
	for {
		select {
		case <-periodicTicker.C:
			for retryCount := uint8(1); retryCount <= maxRetryCount; retryCount++ {
				ci, err = nhc.ChainInfo()
				if err == nil {
					break
				}

				log.Printf("Cannot get height of %s=%s: %s\n", nhc.nodeInfo.Endpoint, nhc.nodeInfo.IdentityKey, err)
				if retryCount <= maxRetryCount {
					log.Printf("Retrying to get chain height from %s (attempt %d/%d)\n", nhc.nodeInfo.Endpoint, retryCount, maxRetryCount)
					time.Sleep(time.Second)
					continue
				}

				return prevHeight, err
			}

			if ci.Height == prevHeight {
				return ci.Height, ErrNodeGotStuck
			}

			if ci.Height >= expectedHeight {
				log.Printf("Node %s=%v has reached the required height\n", nhc.nodeInfo.Endpoint, nhc.nodeInfo.IdentityKey)
				return ci.Height, nil
			}

			log.Printf("Still waiting for node %s=%s to reach height: %d, current: %d\n", nhc.nodeInfo.Endpoint, nhc.nodeInfo.IdentityKey, expectedHeight, ci.Height)

			avgSecondsPerBlock = tickerDuration / time.Duration(ci.Height-prevHeight)
			if expectedHeight-ci.Height < multiplier {
				updateTicker = true
				tickerDuration = time.Duration(expectedHeight-ci.Height) * avgSecondsPerBlock
			} else if avgSecondsPerBlock != DefaultAvgSecondsPerBlock {
				updateTicker = true
				tickerDuration = time.Duration(multiplier) * avgSecondsPerBlock
			}

			if updateTicker {
				if tickerDuration < DefaultAvgSecondsPerBlock {
					tickerDuration = DefaultAvgSecondsPerBlock
				}
				updateTicker = false
				periodicTicker.Reset(tickerDuration)
			}

			prevHeight = ci.Height
		}
	}
}
