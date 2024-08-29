package health

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
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
	// Health checkers by pubKeys
	validCheckersMu sync.Mutex
	validCheckers   map[string]*NodeHealthChecker

	// pubKey and last known height
	knownStuckNodesMu sync.Mutex
	knownStuckNodes   map[string]uint64

	client        *crypto.KeyPair
	mode          packets.ConnectionSecurityMode
	maxConnection int
}

func NewNodeHealthCheckerPool(client *crypto.KeyPair, mode packets.ConnectionSecurityMode, maxConnection int) *NodeHealthCheckerPool {
	return &NodeHealthCheckerPool{
		validCheckers:     make(map[string]*NodeHealthChecker),
		knownStuckNodesMu: sync.Mutex{},
		knownStuckNodes:   make(map[string]uint64),
		client:            client,
		mode:              mode,
		maxConnection:     maxConnection,
	}
}

func (ncp *NodeHealthCheckerPool) ResetPeers() {
	ncp.validCheckers = map[string]*NodeHealthChecker{}
}

func (ncp *NodeHealthCheckerPool) ReviewConnections() {
	log.Printf("Reviewing exist valid checkers...")
	defer log.Printf("Finished reviewing exist valid checkers")

	ncp.validCheckersMu.Lock()
	defer ncp.validCheckersMu.Unlock()

	toDeleteMu := sync.Mutex{}
	toDelete := map[string]*NodeHealthChecker{}

	wg := sync.WaitGroup{}
	for s, checker := range ncp.validCheckers {
		wg.Add(1)
		go func(s string, c *NodeHealthChecker) {
			defer wg.Done()

			_, err := c.ChainInfo()
			if err != nil {
				toDeleteMu.Lock()
				toDelete[s] = c
				toDeleteMu.Unlock()
			}
		}(s, checker)
	}

	wg.Wait()

	for s, checker := range toDelete {
		log.Printf("Delete %s from valid checkers", s)

		checker.Close()
		delete(ncp.validCheckers, s)
	}
}

// ConnectToNodes connects to nodes. Returns map with pubKey as key and nodes info as value and error
func (ncp *NodeHealthCheckerPool) ConnectToNodes(nodeInfos []*NodeInfo, discover bool) (failedConnectionsNodes map[string]*NodeInfo, err error) {
	log.Printf("Start conecting to nodes...")
	defer log.Printf("Finished conecting to nodes...")

	ncp.ReviewConnections()

	ncp.validCheckersMu.Lock()
	if len(ncp.validCheckers) >= ncp.maxConnection {
		ncp.validCheckersMu.Unlock()
		return
	}
	ncp.validCheckersMu.Unlock()

	chInfo := make(chan *NodeInfo, len(nodeInfos)*5)
	for _, info := range nodeInfos {
		chInfo <- info
	}

	failedConnectionsNodesMutex := sync.Mutex{}
	failedConnectionsNodes = make(map[string]*NodeInfo)

	connectedNodesMutex := sync.Mutex{}
	connectedNodes := make(map[string]*NodeHealthChecker)

	waiting := int32(0)
	handled := make(map[string]struct{})
	for {
		select {
		case info, ok := <-chInfo:
			if !ok || info == nil {
				if len(connectedNodes) == 0 {
					return nil, ErrCannotConnect
				}

				ncp.validCheckers = connectedNodes
				return failedConnectionsNodes, nil
			}

			if _, ok := handled[info.IdentityKey.String()]; ok {
				v := atomic.LoadInt32(&waiting)
				if v == 0 && len(chInfo) == 0 {
					close(chInfo)
				}

				continue
			}

			handled[info.IdentityKey.String()] = struct{}{}
			atomic.AddInt32(&waiting, 1)
			go func(info *NodeInfo) {
				defer atomic.AddInt32(&waiting, -1)

				checker, err := ncp.MaybeConnectToNode(info)
				if err != nil {
					failedConnectionsNodesMutex.Lock()
					failedConnectionsNodes[info.IdentityKey.String()] = info
					failedConnectionsNodesMutex.Unlock()

					chInfo <- info
					return
				}

				connectedNodesMutex.Lock()
				connectedNodes[info.IdentityKey.String()] = checker
				connectedNodesMutex.Unlock()

				if !discover {
					chInfo <- info
					return
				}

				nodeList, err := checker.NodeList()
				if err != nil {
					chInfo <- info
					log.Printf("Error getting list of nodes from %s=%v: %s\n", checker.nodeInfo.Endpoint, checker.nodeInfo.IdentityKey, err)
					return
				}

				log.Printf("Node %s=%v returned %d nodes\n", checker.nodeInfo.Endpoint, checker.nodeInfo.IdentityKey, len(nodeList))
				if len(nodeList) == 0 {
					chInfo <- info
				}

				for _, nodeInfo := range nodeList {
					chInfo <- nodeInfo
				}
			}(info)
		}
	}
}

func (ncp *NodeHealthCheckerPool) MaybeConnectToNode(info *NodeInfo) (*NodeHealthChecker, error) {
	log.Printf("Start maybe connect to %s=%v", info.Endpoint, info.IdentityKey)
	defer log.Printf("Finished maybe connect to %s=%v", info.Endpoint, info.IdentityKey)

	ncp.validCheckersMu.Lock()
	vd, ok := ncp.validCheckers[info.IdentityKey.String()]
	ncp.validCheckersMu.Unlock()
	if ok {
		return vd, nil
	}

	log.Printf("Dialing %s=%v", info.Endpoint, info.IdentityKey)
	nc, err := NewNodeHealthChecker(ncp.client, info, ncp.mode)
	if err != nil {
		log.Printf("Error connecting to %s: %s", info.Endpoint, err)
		return nil, err
	}
	log.Printf("Connected to %s=%v", info.Endpoint, info.IdentityKey)

	ncp.validCheckersMu.Lock()
	defer ncp.validCheckersMu.Unlock()

	existNc, ok := ncp.validCheckers[info.IdentityKey.String()]
	if ok {
		return existNc, nil
	}

	ncp.validCheckers[info.IdentityKey.String()] = nc
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
	nextStuckCheck := time.Now().Add(DefaultAvgSecondsPerBlock * 20)

	ticker := time.NewTicker(DefaultAvgSecondsPerBlock)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			minHeight, notReached = checkHeight(expectedHeight, notReached)
			if minHeight >= expectedHeight {
				return nil
			}

			if time.Now().After(nextStuckCheck) && prevMinHeight == minHeight {
				//for s, _ := range notReached {
				//	delete(ncp.validCheckers, s)
				//}

				return ErrorSomeNodeGotStuck
			}

			prevMinHeight = minHeight
		}
	}
}

func (ncp *NodeHealthCheckerPool) WaitHeight(expectedHeight uint64) (notReached map[NodeInfo]uint64, reached map[NodeInfo]uint64, err error) {
	if len(ncp.validCheckers) == 0 {
		return nil, nil, ErrNoConnectedPeers
	}

	log.Printf("Waiting for the network (%d validCheckers) to reach the height %d\n", len(ncp.validCheckers), expectedHeight)

	var notReachedMu sync.Mutex
	notReached = make(map[NodeInfo]uint64)
	var reachedMu sync.Mutex
	reached = make(map[NodeInfo]uint64)

	var wg sync.WaitGroup
	for _, checker := range ncp.validCheckers {
		wg.Add(1)
		go func(checker *NodeHealthChecker) {
			defer wg.Done()

			if h, ok := ncp.knownStuckNodes[checker.nodeInfo.IdentityKey.String()]; ok {
				ci, err := checker.ChainInfo()
				if err != nil || h == ci.Height {
					notReachedMu.Lock()
					notReached[*checker.nodeInfo] = h
					notReachedMu.Unlock()

					return
				}

				ncp.knownStuckNodesMu.Lock()
				delete(ncp.knownStuckNodes, checker.nodeInfo.IdentityKey.String())
				ncp.knownStuckNodesMu.Unlock()
			}

			height, err := checker.WaitHeight(expectedHeight)
			if err != nil {
				notReachedMu.Lock()
				notReached[*checker.nodeInfo] = height
				notReachedMu.Unlock()

				if errors.Is(err, ErrNodeGotStuck) {
					ncp.knownStuckNodesMu.Lock()
					ncp.knownStuckNodes[checker.nodeInfo.IdentityKey.String()] = height
					ncp.knownStuckNodesMu.Unlock()
				}

				return
			}

			reachedMu.Lock()
			reached[*checker.nodeInfo] = height
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

			maxAttempts := 3
			for attemptsCount := 1; attemptsCount <= maxAttempts; attemptsCount++ {
				hash, err := checker.BlockHash(height)
				if err != nil && attemptsCount <= maxAttempts {
					if errors.Is(err, ErrNodeNotReachedHeight) {
						log.Printf("Skip getting block hash from %s: %s\n", checker.nodeInfo.Endpoint, err)
						return
					}

					log.Printf("Error getting block hash from %s:%s\n", checker.nodeInfo.Endpoint, err)
					log.Printf("Retrying to get block hash from %s (attempt %d)\n", checker.nodeInfo.Endpoint, attemptsCount)

					time.Sleep(time.Second * 3)
					continue
				}

				if err != nil {
					log.Printf("Cannot get block hash from %s:%s\n", checker.nodeInfo.Endpoint, err)
				}

				hashesMu.Lock()
				hashes[checker.nodeInfo.Endpoint] = hash
				hashesMu.Unlock()

				if err == nil {
					log.Printf("Node %s=%v at height %d has %s hash", checker.nodeInfo.Endpoint, checker.nodeInfo.IdentityKey, height, hash)
					return
				}
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

	ticker := time.NewTicker(DefaultAvgSecondsPerBlock)
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

// Returns a min height and map with pubKey as key and nodes info as value for not reached nodes
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
			notReached[checker.nodeInfo.IdentityKey.String()] = checker
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
