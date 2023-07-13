package main

import (
	"errors"
	"flag"
	"log"
	"strings"

	crypto "github.com/proximax-storage/go-xpx-crypto"
)

var ErrBadPair = errors.New("bad endpoint-key pair")

func main() {
	height := flag.Uint64("height", 0, "Expected height")
	nodesArg := flag.String("nodes", "", "List of values <ip:port>=<nodePubKey>")
	flag.Parse()

	nodeInfos, err := parseNodes(*nodesArg)
	if err != nil {
		log.Fatal(err)
	}

	client, err := crypto.NewRandomKeyPair()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := NewServerConnectorPool(client, nodeInfos, NoneConnectionSecurity)
	if err != nil {
		log.Fatal(err)
	}

	err = pool.WaitHeightAll(*height)
	if err != nil {
		log.Fatal(err)
	}
}

func parseNodes(nodesStr string) ([]*NodeInfo, error) {
	endpointKeyPairs := strings.Split(nodesStr, " ")
	nodeInfos := make([]*NodeInfo, 0, len(endpointKeyPairs))

	var pair []string
	for _, endpointKeyPair := range endpointKeyPairs {
		pair = strings.Split(endpointKeyPair, "=")
		if len(pair) > 2 {
			return nil, errors.New(ErrBadPair.Error() + ": " + endpointKeyPair)
		}

		nodeInfos = append(nodeInfos, NewNodeInfo(pair[1], pair[0]))
	}

	return nodeInfos, nil
}
