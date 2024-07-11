package health

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/tools/health/packets"
	crypto "github.com/proximax-storage/go-xpx-crypto"
)

var (
	ErrorSomeNodeGotStuck  = errors.New("some nodes got stuck")
	ErrNoConnectedPeers    = errors.New("there are no connected peers")
	ErrCannotConnect       = errors.New("could not connect to any peer")
	ErrHashesAreNotTheSame = errors.New("hashes are not the same")
)

type NodeHealthCheckerPool struct {
	// Health checkers by endpoint
	validCheckers map[string]*NodeHealthChecker

	client        *crypto.KeyPair
	mode          packets.ConnectionSecurityMode
	maxConnection int
}

func NewNodeHealthCheckerPool(client *crypto.KeyPair, mode packets.ConnectionSecurityMode, maxConnection int) *NodeHealthCheckerPool {
	return &NodeHealthCheckerPool{
		validCheckers: make(map[string]*NodeHealthChecker),
		client:        client,
		mode:          mode,
		maxConnection: maxConnection,
	}
}

func (ncp *NodeHealthCheckerPool) ConnectToNodes(nodeInfos []*NodeInfo, findConnected bool) (failedConnectionsNodes []*NodeInfo, err error) {
	for _, info := range nodeInfos {
		if _, err := ncp.MaybeConnectToNode(info); err != nil {
			failedConnectionsNodes = append(failedConnectionsNodes, info)
		}

		if len(ncp.validCheckers) >= ncp.maxConnection {
			break
		}
	}

	if findConnected {
		ncp.CollectConnectedNodes()
	}

	if len(ncp.validCheckers) == 0 {
		return nil, ErrCannotConnect
	}

	return failedConnectionsNodes, nil
}

func (ncp *NodeHealthCheckerPool) CollectConnectedNodes() {
	if len(ncp.validCheckers) >= ncp.maxConnection {
		return
	}

	toCheck := make([]*NodeHealthChecker, 0, len(ncp.validCheckers)*5)
	for _, checker := range ncp.validCheckers {
		toCheck = append(toCheck, checker)
	}

	for len(toCheck) > 0 {
		checker := toCheck[0]
		nodeList, err := checker.NodeList()
		if err != nil {
			log.Printf("Error getting list of validCheckers from %s=%v: %s\n", checker.nodeInfo.Endpoint, checker.nodeInfo.IdentityKey, err)
			continue
		}

		log.Printf("Node %s=%v returned %d validCheckers\n", checker.nodeInfo.Endpoint, checker.nodeInfo.IdentityKey, len(nodeList))
		for _, info := range nodeList {
			if len(ncp.validCheckers) >= ncp.maxConnection {
				return
			}

			nc, _ := ncp.MaybeConnectToNode(info)
			if nc != nil {
				toCheck = append(toCheck, nc)
			}
		}

		toCheck = toCheck[1:]
	}
}

func (ncp *NodeHealthCheckerPool) MaybeConnectToNode(info *NodeInfo) (*NodeHealthChecker, error) {
	if _, ok := ncp.validCheckers[info.Endpoint]; ok {
		return nil, nil
	}

	nc, err := NewNodeHealthChecker(ncp.client, info, ncp.mode)
	if err != nil {
		log.Printf("Error connecting to %s: %s", info.Endpoint, err)
		return nil, err
	}

	ncp.validCheckers[info.Endpoint] = nc
	return nc, nil
}

func (ncp *NodeHealthCheckerPool) CheckHeight(expectedHeight uint64) (uint64, map[string]*NodeHealthChecker) {
	return checkHeight(expectedHeight, ncp.validCheckers)
}

func (ncp *NodeHealthCheckerPool) WaitHeightAll(expectedHeight uint64) error {
	if len(ncp.validCheckers) == 0 {
		return ErrNoConnectedPeers
	}

	log.Printf("Waiting for the network (%d validCheckers) to reach the height %d\n", len(ncp.validCheckers), expectedHeight)

	minHeight, notReached := checkHeight(expectedHeight, ncp.validCheckers)
	if minHeight >= expectedHeight {
		return nil
	}

	prevMinHeight := minHeight
	nextStuckCheck := time.Now().Add(AvgSecondsPerBlock * 20)

	ticker := time.NewTicker(AvgSecondsPerBlock)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			minHeight, notReached = checkHeight(expectedHeight, notReached)
			if minHeight >= expectedHeight {
				return nil
			}

			if time.Now().After(nextStuckCheck) && prevMinHeight == minHeight {
				for s, _ := range notReached {
					delete(ncp.validCheckers, s)
				}

				return ErrorSomeNodeGotStuck
			}

			prevMinHeight = minHeight
		}
	}
}

func (ncp *NodeHealthCheckerPool) WaitHeight(expectedHeight uint64) (notReached map[string]uint64, reached map[string]uint64, err error) {
	if len(ncp.validCheckers) == 0 {
		return nil, nil, ErrNoConnectedPeers
	}

	log.Printf("Waiting for the network (%d validCheckers) to reach the height %d\n", len(ncp.validCheckers), expectedHeight)

	var notReachedMu sync.Mutex
	notReached = make(map[string]uint64)
	var reachedMu sync.Mutex
	reached = make(map[string]uint64)

	var wg sync.WaitGroup
	for _, checker := range ncp.validCheckers {
		wg.Add(1)
		go func(checker *NodeHealthChecker) {
			defer wg.Done()
			height, err := checker.WaitHeight(expectedHeight)
			if err != nil {
				notReachedMu.Lock()
				notReached[checker.nodeInfo.Endpoint] = height
				delete(ncp.validCheckers, checker.nodeInfo.Endpoint)
				notReachedMu.Unlock()

				return
			}

			reachedMu.Lock()
			reached[checker.nodeInfo.Endpoint] = height
			reachedMu.Unlock()
		}(checker)
	}

	wg.Wait()
	return notReached, reached, nil
}

func (ncp *NodeHealthCheckerPool) GetHashes(height uint64) (map[string]sdk.Hash, error) {
	if len(ncp.validCheckers) == 0 {
		return nil, ErrNoConnectedPeers
	}

	hashesMu := sync.Mutex{}
	hashes := make(map[string]sdk.Hash)

	wg := sync.WaitGroup{}
	for _, checker := range ncp.validCheckers {
		wg.Add(1)
		go func(checker *NodeHealthChecker) {
			defer wg.Done()

			for attemptsCount := 0; attemptsCount < 3; attemptsCount++ {
				hash, err := checker.BlockHash(height)
				if err != nil && attemptsCount < 3 {
					log.Printf("Error getting block hash from %s:%s\n", checker.nodeInfo.Endpoint, err)
					log.Printf("Retrying to get block hash from %s (attempt %d)\n", checker.nodeInfo.Endpoint, attemptsCount)

					time.Sleep(AvgSecondsPerBlock)
					continue
				}

				if err != nil {
					delete(ncp.validCheckers, checker.nodeInfo.Endpoint)
				}

				if err == nil {
					log.Printf("Node %s=%v at height %d has %s hash", checker.nodeInfo.Endpoint, checker.nodeInfo.IdentityKey, height, hash)
				}

				hashesMu.Lock()
				hashes[checker.nodeInfo.Endpoint] = hash
				hashesMu.Unlock()
			}
		}(checker)
	}

	wg.Wait()
	return hashes, nil
}

func (ncp *NodeHealthCheckerPool) CompareHashes(height uint64) (map[string]sdk.Hash, error) {
	if len(ncp.validCheckers) == 0 {
		return nil, ErrNoConnectedPeers
	}

	hashes, err := ncp.GetHashes(height)
	if err != nil {
		return nil, err
	}

	uniqueHashes := map[sdk.Hash]struct{}{}
	for _, hash := range hashes {
		uniqueHashes[hash] = struct{}{}
		if len(uniqueHashes) > 1 {
			return hashes, ErrHashesAreNotTheSame
		}
	}

	return hashes, nil
}

func (ncp *NodeHealthCheckerPool) WaitAllHashesEqual(height uint64) error {
	log.Printf("Waiting for the same block hash at %d height\n", height)

	_, err := ncp.CompareHashes(height)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(AvgSecondsPerBlock)
	for {
		select {
		case <-ticker.C:
			if _, err = ncp.CompareHashes(height); err != nil {
				ticker.Stop()
				return err
			}
		}
	}
}

func checkHeight(expectedHeight uint64, nodeHealthCheckers map[string]*NodeHealthChecker) (uint64, map[string]*NodeHealthChecker) {
	notReachedLock := sync.Mutex{}
	notReached := make(map[string]*NodeHealthChecker)

	minHeightLock := sync.Mutex{}
	minHeight := expectedHeight

	wg := sync.WaitGroup{}
	for _, checker := range nodeHealthCheckers {
		wg.Add(1)
		go func(checker *NodeHealthChecker) {
			defer wg.Done()

			ci, err := checker.ChainInfo()
			if err == nil && ci.Height >= expectedHeight {
				return
			}

			if err != nil {
				log.Printf("Error getting chain info from %s=%v: %s\n", checker.nodeInfo.Endpoint, checker.nodeInfo.IdentityKey, err)
			} else if ci.Height < expectedHeight {
				log.Printf("Node %s=%v has not reached the required height. Expected : %d, current: %d\n", checker.nodeInfo.Endpoint, checker.nodeInfo.IdentityKey, expectedHeight, ci.Height)
				minHeightLock.Lock()
				if ci.Height < minHeight {
					minHeight = ci.Height
				}
				minHeightLock.Unlock()
			}

			notReachedLock.Lock()
			notReached[checker.nodeInfo.Endpoint] = checker
			notReachedLock.Unlock()
		}(checker)
	}

	wg.Wait()

	if len(notReached) == 0 {
		log.Println("All nodes reached the required height")
	} else {
		log.Printf("%d nodes has not reach the required height\n", len(notReached))
	}

	return minHeight, notReached
}
