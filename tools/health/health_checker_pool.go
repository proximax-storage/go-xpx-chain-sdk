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

type NodeHealthCheckerPool struct {
	// Health checkers by endpoint
	nodeHealthCheckers map[string]*NodeHealthChecker

	// Endpoints to which a dial attempt has been performed, whether successful or not
	dialed map[string]bool

	client        *crypto.KeyPair
	mode          packets.ConnectionSecurityMode
	maxConnection int
}

func NewNodeHealthCheckerPool(client *crypto.KeyPair, nodeInfos []*NodeInfo, mode packets.ConnectionSecurityMode, findConnected bool, maxConnection int) (*NodeHealthCheckerPool, error) {
	ncp := &NodeHealthCheckerPool{
		nodeHealthCheckers: make(map[string]*NodeHealthChecker),
		dialed:             make(map[string]bool),
		client:             client,
		mode:               mode,
		maxConnection:      maxConnection,
	}

	for _, info := range nodeInfos {
		ncp.MaybeConnectToNode(info)
		if len(ncp.nodeHealthCheckers) >= ncp.maxConnection {
			break
		}
	}

	if findConnected {
		ncp.CollectConnectedNodes()
	}

	if len(ncp.nodeHealthCheckers) == 0 {
		return nil, errors.New("could not connect to any peer")
	}

	return ncp, nil
}

func (ncp *NodeHealthCheckerPool) CollectConnectedNodes() {
	if len(ncp.nodeHealthCheckers) >= ncp.maxConnection {
		return
	}

	toCheck := make([]*NodeHealthChecker, 0, len(ncp.nodeHealthCheckers)*5)
	for _, checker := range ncp.nodeHealthCheckers {
		toCheck = append(toCheck, checker)
	}

	for len(toCheck) > 0 {
		checker := toCheck[0]
		nodeList, err := checker.NodeList()
		if err == nil {
			log.Printf("Node %s=%v returned %d nodes\n", checker.nodeInfo.Endpoint, checker.nodeInfo.IdentityKey, len(nodeList))
			for _, info := range nodeList {
				if len(ncp.nodeHealthCheckers) >= ncp.maxConnection {
					return
				}

				nc := ncp.MaybeConnectToNode(info)
				if nc != nil {
					toCheck = append(toCheck, nc)
				}
			}
		} else {
			log.Printf("Error getting list of nodes from %s=%v: %s\n", checker.nodeInfo.Endpoint, checker.nodeInfo.IdentityKey, err)
		}

		toCheck = toCheck[1:]
	}
}

func (ncp *NodeHealthCheckerPool) MaybeConnectToNode(info *NodeInfo) *NodeHealthChecker {
	if _, ok := ncp.nodeHealthCheckers[info.Endpoint]; ok {
		return nil
	}

	if _, ok := ncp.dialed[info.Endpoint]; ok {
		return nil
	}

	nc, err := NewNodeHealthChecker(ncp.client, info, ncp.mode)
	ncp.dialed[info.Endpoint] = true
	if err != nil {
		log.Printf("Error connecting to %s: %s", info.Endpoint, err)
		return nil
	}

	ncp.nodeHealthCheckers[info.Endpoint] = nc

	return nc
}

func (ncp *NodeHealthCheckerPool) CheckHeight(expectedHeight uint64, nodeHealthCheckers map[string]*NodeHealthChecker) (uint64, map[string]*NodeHealthChecker) {
	log.Println("Start checking height")
	type CheckHeightResult struct {
		Endpoint    string
		IdentityKey *crypto.PublicKey
		Height      uint64
	}

	nodeCheckCh := make(chan CheckHeightResult)
	for _, checker := range nodeHealthCheckers {
		go func(checker *NodeHealthChecker) {
			endpoint := checker.nodeInfo.Endpoint
			identityKey := checker.nodeInfo.IdentityKey
			ci, err := checker.ChainInfo()
			if err != nil {
				log.Printf("error getting chain info from %s=%v: %s\n", endpoint, identityKey, err)
				nodeCheckCh <- CheckHeightResult{endpoint, identityKey, 0}
				return
			}

			nodeCheckCh <- CheckHeightResult{endpoint, identityKey, ci.Height}
		}(checker)
	}

	nodeCount := len(nodeHealthCheckers)
	reachedCount := 0
	notReached := make(map[string]*NodeHealthChecker)
	minHeight := expectedHeight
	for {
		select {
		case res := <-nodeCheckCh:
			if res.Height >= expectedHeight {
				log.Printf("Node %s=%v has reached the required height\n", res.Endpoint, res.IdentityKey)
				reachedCount++
			} else {
				log.Printf("Node %s=%v has not reached the required height. Expected : %d, current: %d\n", res.Endpoint, res.IdentityKey, expectedHeight, res.Height)
				notReached[res.Endpoint] = ncp.nodeHealthCheckers[res.Endpoint]
			}

			if minHeight > res.Height && res.Height > 0 {
				minHeight = res.Height
			}

			notReachedCount := len(notReached)
			if reachedCount+notReachedCount == nodeCount {
				if notReachedCount > 0 {
					log.Printf("%d nodes has not reach the required height. Continue waiting\n", notReachedCount)
				} else {
					log.Println("All nodes reached the required height")
				}
				return minHeight, notReached
			}
		}
	}
}

func (ncp *NodeHealthCheckerPool) WaitHeightAll(expectedHeight uint64, timeout time.Duration) error {
	log.Printf("Waiting for the network (%d nodes) to reach the height %d\n", len(ncp.nodeHealthCheckers), expectedHeight)

	minHeight, notReached := ncp.CheckHeight(expectedHeight, ncp.nodeHealthCheckers)
	if minHeight >= expectedHeight {
		return nil
	}

	ticker := time.NewTicker(AvgSecondsPerBlock)
	var timeoutCh <-chan time.Time
	if timeout > 0 {
		timeoutCh = time.After(timeout)
	}

	for {
		select {
		case <-ticker.C:
			if minHeight, notReached = ncp.CheckHeight(expectedHeight, notReached); minHeight >= expectedHeight {
				ticker.Stop()
				return nil
			}
		case <-timeoutCh:
			log.Printf("Timeout reached while waiting for network to reach the height %d", expectedHeight)
			for _, node := range notReached {
				log.Printf("Removing node from checklist: %v=%v", node.nodeInfo.Endpoint, node.nodeInfo.IdentityKey)
				delete(ncp.nodeHealthCheckers, node.nodeInfo.Endpoint)
			}
			return nil
		}
	}
}

func (ncp *NodeHealthCheckerPool) CheckHashes(height uint64) bool {
	hashesCh := make(chan sdk.Hash)
	for _, checker := range ncp.nodeHealthCheckers {
		go func(checker *NodeHealthChecker) {
			hash, err := checker.BlockHash(height)
			log.Printf("Node %s=%v has %s", checker.nodeInfo.Endpoint, checker.nodeInfo.IdentityKey, hash)
			if err != nil {
				log.Printf("Error getting block hash from %s:%s\n", checker.nodeInfo.Endpoint, err)
				hashesCh <- sdk.Hash{}
				return
			}

			hashesCh <- hash
		}(checker)
	}

	nodeCount := len(ncp.nodeHealthCheckers)
	hashes := make(map[sdk.Hash]int)
	for {
		select {
		case hash := <-hashesCh:
			hashes[hash]++

			hashCount := 0
			for _, count := range hashes {
				hashCount += count
			}
			if hashCount < nodeCount {
				continue
			}

			if len(hashes) > 1 {
				log.Printf("Block hashes differ (hash:count of returned): %v\n", hashes)
				log.Println("Continue waiting")
				return false
			} else {
				for hash := range hashes {
					if hash.Empty() {
						log.Printf("Could not collect block hashes at %d height\n", height)
						return false
					} else {
						log.Printf("All nodes got the same block hash at %d height\n", height)
						return true
					}
				}
			}
		}
	}
}

func (ncp *NodeHealthCheckerPool) WaitAllHashesEqual(height uint64) error {
	log.Printf("Waiting for the same block hash at %d height\n", height)

	success := ncp.CheckHashes(height)
	if success {
		return nil
	}

	ticker := time.NewTicker(AvgSecondsPerBlock)
	for {
		select {
		case <-ticker.C:
			if success = ncp.CheckHashes(height); success {
				ticker.Stop()
				return nil
			}
		}
	}
}

func (ncp *NodeHealthCheckerPool) FindInconsistentHashesAtHeight(height uint64) map[string]sdk.Hash {
	type CheckNodeResult struct {
		Endpoint string
		Hash     sdk.Hash
	}

	nodeCount := len(ncp.nodeHealthCheckers)
	nodeCheckCh := make(chan CheckNodeResult, nodeCount/2)
	var mu sync.Mutex
	for _, checker := range ncp.nodeHealthCheckers {
		go func(checker *NodeHealthChecker) {
			var hash sdk.Hash
			var err error

			for {
				hash, err = checker.BlockHash(height)
				if err == nil {
					break
				}

				if err == ErrReturnedZeroHashes {
					mu.Lock()
					nodeCount--
					mu.Unlock()
					return
				}

				// Handle connection error
				log.Printf("Error getting block hash from %s: %s. Retrying connection...\n", checker.nodeInfo.Endpoint, err)
				nc, err := NewNodeHealthChecker(ncp.client, checker.nodeInfo, ncp.mode)
				if err != nil {
					log.Printf("Error connecting to %s: %s\n", checker.nodeInfo.Endpoint, err)
					mu.Lock()
					nodeCount--
					mu.Unlock()
					return
				}

				mu.Lock()
				ncp.nodeHealthCheckers[checker.nodeInfo.Endpoint] = nc
				mu.Unlock()
			}

			log.Printf("Node %s=%v has %s", checker.nodeInfo.Endpoint, checker.nodeInfo.IdentityKey, hash)
			nodeCheckCh <- CheckNodeResult{checker.nodeInfo.Endpoint, hash}
		}(checker)
	}

	hashes := make(map[sdk.Hash]int)
	nodeHashResults := make(map[string]sdk.Hash)
	for {
		select {
		case res := <-nodeCheckCh:
			nodeHashResults[res.Endpoint] = res.Hash
			hashes[res.Hash]++

			hashCount := 0
			for _, count := range hashes {
				hashCount += count
			}

			if hashCount < nodeCount {
				continue
			}

			if len(hashes) > 1 {
				log.Printf("Block hashes differ (hash:count of returned): %v\n", hashes)
				return nodeHashResults
			}

			log.Printf("All nodes got the same block hash at %d height\n", height)
			return nil
		}
	}
}
