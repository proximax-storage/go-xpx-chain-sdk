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

const AvgSecondsPerBlock = 15 * time.Second

var (
	ErrReturnedZeroHashes = errors.New("returned zero hashes")
	ErrTimedOut           = errors.New("timed out")
	ErrTooLargeWaitTime   = errors.New("too large wait time")
	ErrNodeGotStuck       = errors.New("node got stuck")
)

type (
	NodeInfo struct {
		IdentityKey *crypto.PublicKey
		Endpoint    string
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

func NewNodeInfo(pKey, addr string) (*NodeInfo, error) {
	k, err := crypto.NewPublicKeyfromHex(pKey)
	if err != nil {
		return nil, err
	}

	return &NodeInfo{
		IdentityKey: k,
		Endpoint:    addr,
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
	blockHashesReq := packets.NewBlockHashesRequest(height, 1)
	blockHashesResp := &packets.BlockHashesResponse{}
	err := nhc.handler.CommonHandle(blockHashesReq, blockHashesResp)
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
			IdentityKey: node.IdentityKey,
			Endpoint:    node.Host + ":" + strconv.Itoa(int(node.Port)),
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
	lastHeight := ci.Height

	multiplier := expectedHeight - ci.Height
	if multiplier > 4 {
		multiplier = 4
	}

	tickerDuration := AvgSecondsPerBlock * time.Duration(multiplier)
	periodicTicker := time.NewTicker(tickerDuration)
	defer periodicTicker.Stop()

	var retryCount uint8
	maxRetryCount := uint8(3)
	for {
		select {
		case <-time.After(tickerDuration):
			ci, err := nhc.ChainInfo()
			if err != nil {
				if retryCount < maxRetryCount {
					retryCount++
					log.Printf("Retrying to get chain height from %s (attempt %d/%d)\n", nhc.nodeInfo.Endpoint, retryCount, maxRetryCount)
					continue
				}

				return lastHeight, err
			}
			retryCount = 0

			if ci.Height == lastHeight {
				return lastHeight, ErrNodeGotStuck
			}

			if ci.Height >= expectedHeight {
				log.Printf("Node %s=%v has reached the required height\n", nhc.nodeInfo.Endpoint, nhc.nodeInfo.IdentityKey)
				return ci.Height, nil
			}

			lastHeight = ci.Height
			log.Printf("Waiting for node %s=%s to reach height: %d, current: %d\n", nhc.nodeInfo.Endpoint, nhc.nodeInfo.IdentityKey, expectedHeight, ci.Height)

			duration := time.Duration(expectedHeight-ci.Height) * AvgSecondsPerBlock
			if duration < tickerDuration {
				// add AvgSecondsPerBlock just as an extra time
				periodicTicker.Reset(duration + AvgSecondsPerBlock)
			}
		}
	}
}
