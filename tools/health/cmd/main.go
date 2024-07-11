package main

import (
	"errors"
	"flag"
	"log"
	"math"
	"strings"

	"github.com/proximax-storage/go-xpx-chain-sdk/tools/health"
	"github.com/proximax-storage/go-xpx-chain-sdk/tools/health/packets"
	crypto "github.com/proximax-storage/go-xpx-crypto"
)

var ErrBadPair = errors.New("bad endpoint-key pair")

func main() {
	height := flag.Uint64("height", 0, "Expected height")
	nodesArg := flag.String("nodes", "", "List of values <ip:port>=<nodePubKey>")
	discover := flag.Bool("discover", true, "Discover connected nodes (Default is true)")
	flag.Parse()

	nodeInfos, err := parseNodes(*nodesArg)
	if err != nil {
		log.Fatal(err)
	}

	client, err := crypto.NewRandomKeyPair()
	if err != nil {
		log.Fatal(err)
	}

	pool := health.NewNodeHealthCheckerPool(client, packets.NoneConnectionSecurity, math.MaxInt)

	_, err = pool.ConnectToNodes(nodeInfos, *discover)
	if err != nil {
		log.Fatal(err)
	}

	err = pool.WaitHeightAll(*height)
	if err != nil {
		log.Fatal(err)
	}

	_, err = pool.CompareHashes(*height)
	if err != nil {
		log.Fatal(err)
	}
}

func parseNodes(nodesStr string) ([]*health.NodeInfo, error) {
	nodesStr = strings.TrimSpace(nodesStr)
	endpointKeyPairs := strings.Split(nodesStr, " ")
	nodeInfos := make([]*health.NodeInfo, 0, len(endpointKeyPairs))

	var pair []string
	for _, endpointKeyPair := range endpointKeyPairs {
		pair = strings.Split(endpointKeyPair, "=")
		if len(pair) > 2 {
			return nil, errors.New(ErrBadPair.Error() + ": " + endpointKeyPair)
		}

		ni, err := health.NewNodeInfo(pair[1], pair[0])
		if err != nil {
			return nil, err
		}

		nodeInfos = append(nodeInfos, ni)
	}

	return nodeInfos, nil
}
