package health

import (
	"errors"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/tools/health/packets"
	crypto "github.com/proximax-storage/go-xpx-crypto"
)

const AvgSecondsPerBlock = 15 * time.Second

type (
	Byter interface {
		Bytes() []byte
	}

	Parser interface {
		Parse([]byte) error
	}

	PacketHeader interface {
		Byter
		Parser
	}

	Packet interface {
		Byter
		Parser

		Header() PacketHeader
	}

	NodeIo interface {
		Write(Byter) (int, error)
		Read(Parser, int) error
		Close() error
	}

	NodeInfo struct {
		IdentityKey *crypto.PublicKey
		Endpoint    string
	}

	NodeHealthChecker struct {
		nodeInfo *NodeInfo

		handler *Handler
	}

	NodeConnectorPool struct {
		nodeHealthChecker []*NodeHealthChecker
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

func (nhc *NodeHealthChecker) LastBlockHash(height uint64) (*sdk.Hash, error) {
	blockHashesReq := packets.NewBlockHashesRequest(height, 1)
	blockHashesResp := &packets.BlockHashesResponse{}
	err := nhc.handler.CommonHandle(blockHashesReq, blockHashesResp)
	if err != nil {
		return nil, err
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
		ni = append(ni, &NodeInfo{
			IdentityKey: node.IdentityKey,
			Endpoint:    node.Host + ":" + strconv.Itoa(int(node.Port)),
		})
	}

	return ni, nil
}

func (nhc *NodeHealthChecker) WaitHeight(height uint64) error {
	ci, err := nhc.ChainInfo()
	if err != nil {
		return err
	}

	if ci.Height >= height {
		return nil
	}

	// global avg time + extra time
	delta := height - ci.Height
	globalTimer := time.After(time.Duration(delta)*AvgSecondsPerBlock + time.Minute)
	for {
		println(delta)
		select {
		case <-globalTimer:
			return errors.New("timed out")
		case <-time.After(time.Duration(delta) * AvgSecondsPerBlock):
			ci, err = nhc.ChainInfo()
			if err != nil {
				return err
			}

			if ci.Height >= height {
				return nil
			}
			delta = height - ci.Height
		}
	}
}

func NewServerConnectorPool(client *crypto.KeyPair, nodeInfos []*NodeInfo, mode packets.ConnectionSecurityMode) (*NodeConnectorPool, error) {
	ncp := &NodeConnectorPool{nodeHealthChecker: make([]*NodeHealthChecker, 0, len(nodeInfos))}
	for _, info := range nodeInfos {
		nc, err := NewNodeHealthChecker(client, info, mode)
		if err != nil {
			return nil, err
		}
		ncp.nodeHealthChecker = append(ncp.nodeHealthChecker, nc)
	}

	return ncp, nil
}

func (ncp *NodeConnectorPool) WaitHeightAll(height uint64) error {
	countOfRich := 0
	ticker := time.NewTicker(AvgSecondsPerBlock)
	for {
		select {
		case <-ticker.C:
			minHeight := uint64(math.MaxUint64)
			countOfRich = 0
			for _, connector := range ncp.nodeHealthChecker {
				ci, err := connector.ChainInfo()
				if err != nil {
					log.Printf("cannot get the height from %s:%s\n", connector.nodeInfo.Endpoint, err)
					continue
				}

				if ci.Height < minHeight {
					minHeight = ci.Height
				}

				if ci.Height == height {
					countOfRich++
				}
			}

			if countOfRich == len(ncp.nodeHealthChecker) {
				return nil
			}

			ticker = time.NewTicker(time.Duration(height-minHeight) * AvgSecondsPerBlock)
		}
	}
}

func (ncp *NodeConnectorPool) WaitAllHashesEqual(height uint64) error {
	for {
		select {
		case <-time.After(AvgSecondsPerBlock):
			hashes := make(map[sdk.Hash]int)
			for _, connector := range ncp.nodeHealthChecker {
				h, err := connector.LastBlockHash(height)
				if err != nil {
					log.Printf("cannot get the last hash from %s:%s\n", connector.nodeInfo.Endpoint, err)
					continue
				}

				hashes[*h] += 1
				if len(hashes) > 1 {
					break
				}
			}

			if len(hashes) == 1 {
				for _, u := range hashes {
					if u == len(ncp.nodeHealthChecker) {
						return nil
					}
				}
			}
		}
	}
}
