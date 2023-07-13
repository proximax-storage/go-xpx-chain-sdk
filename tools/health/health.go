package main

import (
	"errors"
	"net"
	"sync"
	"time"

	crypto "github.com/proximax-storage/go-xpx-crypto"

	"go.uber.org/multierr"
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

	buf, err := readFromConn(sc.conn, ChainInfoSize)
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

func (ncp *NodeConnectorPool) WaitHeightAll(height uint64) error {
	var multiErr error
	wg := &sync.WaitGroup{}
	for _, connector := range ncp.nodeConnectors {
		wg.Add(1)
		go func(nodeConnector *NodeConnector) {
			defer wg.Done()
			err := nodeConnector.WaitHeight(height)
			if err != nil {
				multiErr = multierr.Append(multiErr, err)
			}
		}(connector)
	}
	wg.Wait()

	return multiErr
}
