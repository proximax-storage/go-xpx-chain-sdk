package main

import (
	"errors"
	"log"
	"math"
	"net"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	crypto "github.com/proximax-storage/go-xpx-crypto"
)

const AvgSecondsPerBlock = 15 * time.Second

type (
	ChainScore struct {
		High uint64
		Low  uint64
	}

	ChainInfo struct {
		Height     uint64
		ChainScore ChainScore
	}

	NodeInfo struct {
		IdentityKey string
		Endpoint    string
	}

	NodeConnector struct {
		NodeInfo *NodeInfo

		auth *AuthPacketHandler
		conn net.Conn
	}

	NodeConnectorPool struct {
		nodeConnectors []*NodeConnector
	}
)

func NewNodeInfo(pKey, addr string) *NodeInfo {
	return &NodeInfo{
		IdentityKey: pKey,
		Endpoint:    addr,
	}
}

func NewNodeConnector(client *crypto.KeyPair, info *NodeInfo, mode ConnectionSecurityMode) (*NodeConnector, error) {
	serverKey, err := crypto.NewPublicKeyfromHex(info.IdentityKey)
	if err != nil {
		return nil, err
	}

	connection, err := net.Dial("tcp", info.Endpoint)
	if err != nil {
		return nil, err
	}

	auth := NewAuthPacketHandler(client, serverKey, mode, connection)
	err = auth.Start()
	if err != nil {
		return nil, err
	}

	return &NodeConnector{
		NodeInfo: info,
		auth:     auth,
		conn:     connection,
	}, nil
}

func NewServerConnectorPool(client *crypto.KeyPair, nodeInfos []*NodeInfo, mode ConnectionSecurityMode) (*NodeConnectorPool, error) {
	ncp := &NodeConnectorPool{nodeConnectors: make([]*NodeConnector, 0, len(nodeInfos))}
	for _, info := range nodeInfos {
		nc, err := NewNodeConnector(client, info, mode)
		if err != nil {
			return nil, err
		}
		ncp.nodeConnectors = append(ncp.nodeConnectors, nc)
	}

	return ncp, nil
}

func (sc *NodeConnector) Close() error {
	return sc.conn.Close()
}

func (sc *NodeConnector) ChainInfo() (*ChainInfo, error) {
	chainInfo := NewPacketHeader(ChainInfoPacketType)
	_, err := sc.conn.Write(chainInfo.Bytes())
	if err != nil {
		return nil, err
	}

	buf, err := readFromConn(sc.conn, ChainInfoResponseSize)
	if err != nil {
		return nil, err
	}

	chainInfoResponse := &ChainInfoResponse{}
	err = chainInfoResponse.Parse(buf)
	if err != nil {
		return nil, err
	}

	return &ChainInfo{
		Height: chainInfoResponse.Height,
		ChainScore: ChainScore{
			chainInfoResponse.ScoreHigh,
			chainInfoResponse.ScoreLow,
		},
	}, nil
}

func (sc *NodeConnector) WaitHeight(height uint64) error {
	ci, err := sc.ChainInfo()
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
			ci, err = sc.ChainInfo()
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

func (sc *NodeConnector) GetLastHash(height uint64) (*sdk.Hash, error) {
	chainInfo := NewBlockHashesRequest(height, 1)
	_, err := sc.conn.Write(chainInfo.Bytes())
	if err != nil {
		return nil, err
	}

	buf, err := readFromConn(sc.conn, BlockHashesResponseSize+HashSize)
	if err != nil {
		return nil, err
	}

	blockHashesResponse := &BlockHashesResponse{}
	err = blockHashesResponse.Parse(buf)
	if err != nil {
		return nil, err
	}

	return blockHashesResponse.Hashes[0], nil
}

func (ncp *NodeConnectorPool) WaitHeightAll(height uint64) error {
	countOfRich := 0
	ticker := time.NewTicker(AvgSecondsPerBlock)
	for {
		select {
		case <-ticker.C:
			minHeight := uint64(math.MaxUint64)
			countOfRich = 0
			for _, connector := range ncp.nodeConnectors {
				ci, err := connector.ChainInfo()
				if err != nil {
					log.Printf("cannot get the height from %s:%s\n", connector.NodeInfo.Endpoint, err)
					continue
				}

				if ci.Height < minHeight {
					minHeight = ci.Height
				}

				if ci.Height == height {
					countOfRich++
				}
			}

			if countOfRich == len(ncp.nodeConnectors) {
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
			for _, connector := range ncp.nodeConnectors {
				h, err := connector.GetLastHash(height)
				if err != nil {
					log.Printf("cannot get the last hash from %s:%s\n", connector.NodeInfo.Endpoint, err)
					continue
				}

				hashes[*h] += 1
				if len(hashes) > 1 {
					break
				}
			}

			if len(hashes) == 1 {
				for _, u := range hashes {
					if u == len(ncp.nodeConnectors) {
						return nil
					}
				}
			}
		}
	}
}
