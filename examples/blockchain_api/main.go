// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
)

var (
	baseUrls = []string{"http://localhost:3000"}
)

// Simple Blockchain API request
func main() {
	conf, err := sdk.NewDefaultConfig(baseUrls)
	if err != nil {
		panic(err)
	}

	// Use the default http client
	client := sdk.NewClient(nil, conf)
	client.SetupConfigFromRest(context.Background())

	// Get the chain height
	chainHeight, err := client.Blockchain.GetBlockchainHeight(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n\n", chainHeight)

	// Get the chain score
	chainScore, err := client.Blockchain.GetBlockchainScore(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n\n", chainScore)

	// Get the Block by height
	blockHeight, err := client.Blockchain.GetBlockByHeight(context.Background(), sdk.Height(9999))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n\n", blockHeight)

	// Get the Block Transactions
	transactions, err := client.Blockchain.GetBlockTransactions(context.Background(), sdk.Height(1))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n\n", transactions)

	// Get the Blockchain Storage Info
	blockchainStorageInfo, err := client.Blockchain.GetBlockchainStorage(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n\n", blockchainStorageInfo)
}
