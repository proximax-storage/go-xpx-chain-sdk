package health

import (
	"errors"
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

func (nhc *NodeHealthChecker) WaitHeight(expectedHeight uint64) (uint64, error) {
	globalTicker := &time.Ticker{}
	ticker := time.NewTicker(time.Second)
	var height uint64
	for {
		select {
		case <-ticker.C:
			ci, err := nhc.ChainInfo()
			if err != nil {
				return 0, err
			}

			height = ci.Height
			if height >= expectedHeight {
				return height, nil
			}

			duration := time.Duration(expectedHeight-ci.Height) * AvgSecondsPerBlock
			if duration > time.Hour {
				return height, ErrTooLargeWaitTime
			}

			ticker = time.NewTicker(duration)
			if globalTicker == nil {
				globalTicker = time.NewTicker(duration + duration/2)
			}

			log.Printf("Waiting for node %s=%s to reach height: %d, current: %d", nhc.nodeInfo.Endpoint, nhc.nodeInfo.IdentityKey, expectedHeight, height)
		case <-globalTicker.C:
			return height, ErrTimedOut
		}
	}
}
