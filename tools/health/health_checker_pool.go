package health

import (
	"log"
	"math"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/tools/health/packets"
	crypto "github.com/proximax-storage/go-xpx-crypto"
)

type NodeHealthCheckerPool struct {
	// nodeHealthChecker is endpoint as key and health checker as value
	nodeHealthCheckers map[string]*NodeHealthChecker

	client *crypto.KeyPair
	mode   packets.ConnectionSecurityMode
}

func NewNodeHealthCheckerPool(client *crypto.KeyPair, nodeInfos []*NodeInfo, mode packets.ConnectionSecurityMode, findConnected bool) (*NodeHealthCheckerPool, error) {
	ncp := &NodeHealthCheckerPool{
		nodeHealthCheckers: make(map[string]*NodeHealthChecker),
		client:             client,
		mode:               mode,
	}

	for _, info := range nodeInfos {
		if _, ok := ncp.nodeHealthCheckers[info.Endpoint]; ok {
			continue
		}

		nc, err := NewNodeHealthChecker(client, info, mode)
		if err != nil {
			return nil, err
		}

		ncp.nodeHealthCheckers[info.Endpoint] = nc
	}

	if findConnected {
		err := ncp.CollectConnectedNodes()
		if err != nil {
			return nil, err
		}
	}

	return ncp, nil
}

func (ncp *NodeHealthCheckerPool) CollectConnectedNodes() error {
	toCheck := make([]*NodeHealthChecker, 0, len(ncp.nodeHealthCheckers)*5)
	for _, checker := range ncp.nodeHealthCheckers {
		toCheck = append(toCheck, checker)
	}

	for len(toCheck) > 0 {
		nodeList, err := toCheck[0].NodeList()
		if err != nil {
			return err
		}

		for _, info := range nodeList {
			if _, ok := ncp.nodeHealthCheckers[info.Endpoint]; ok {
				continue
			}

			nc, err := NewNodeHealthChecker(ncp.client, info, ncp.mode)
			if err != nil {
				return err
			}

			toCheck = append(toCheck, nc)
			ncp.nodeHealthCheckers[info.Endpoint] = nc
		}

		toCheck = toCheck[1:]
	}

	return nil
}

func (ncp *NodeHealthCheckerPool) WaitHeightAll(expectedHeight uint64) error {
	log.Printf("Start waititing when all nodes will have the same height - %d\n", expectedHeight)

	reached := map[string]struct{}{}
	ticker := time.NewTicker(AvgSecondsPerBlock)
	for {
		select {
		case <-ticker.C:
			minHeight := uint64(math.MaxUint64)
			for _, connector := range ncp.nodeHealthCheckers {
				ci, err := connector.ChainInfo()
				if err != nil {
					log.Printf("cannot get the expectedHeight from %s:%s\n", connector.nodeInfo.Endpoint, err)
					continue
				}

				if ci.Height < minHeight {
					minHeight = ci.Height
				}

				if ci.Height == expectedHeight {
					reached[connector.nodeInfo.Endpoint] = struct{}{}
					log.Printf("Node %s reached the required height\n", connector.nodeInfo.Endpoint)
				} else {
					log.Printf("Node %s has not reach the required height. Expected : %d, currtent: %d\n", connector.nodeInfo.Endpoint, expectedHeight, ci.Height)
				}
			}

			if len(reached) == len(ncp.nodeHealthCheckers) {
				log.Println("All nodes reached the required height")
				return nil
			}

			log.Printf("%d nodes did not reached the required height. Continue waiting\n", len(ncp.nodeHealthCheckers)-len(reached))
			ticker = time.NewTicker(time.Duration(expectedHeight-minHeight) * AvgSecondsPerBlock)
		}
	}
}

func (ncp *NodeHealthCheckerPool) WaitAllHashesEqual(height uint64) error {
	log.Printf("Start waititing when all nodes will have the same block hash at %d height\n", height)

	for {
		select {
		case <-time.After(AvgSecondsPerBlock):
			hashes := make(map[sdk.Hash]int)
			for _, connector := range ncp.nodeHealthCheckers {
				h, err := connector.BlockHash(height)
				if err != nil {
					log.Printf("Cannot get the last hash from %s:%s\n", connector.nodeInfo.Endpoint, err)
					break
				}

				hashes[h] += 1
				if len(hashes) > 1 {
					break
				}
			}

			if len(hashes) != 1 {
				log.Printf("Hashes are not the same. Collected hashes(hash:count of returned): %v\n", hashes)
				log.Println("Continue waiting")
			}

			for _, u := range hashes {
				if u == len(ncp.nodeHealthCheckers) {
					log.Printf("All nodes have the same hash of the block at %d height\n", height)
					return nil
				} else {
					log.Printf("%d nodes did not provide the hash of the block at %d height\n", len(ncp.nodeHealthCheckers)-u, height)
					log.Println("Continue waiting")
				}
			}
		}
	}
}
