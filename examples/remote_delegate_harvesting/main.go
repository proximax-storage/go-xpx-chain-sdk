// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"
)

var (
	baseUrls            = []string{"http://127.0.0.1:3000"}
	HarvesterAccountKey = "7AA907C3D80B3815BE4B4E1470DEEE8BB83BFEB330B9A82197603D09BA947230"
	HarvesterNodeKey    = "D6430327F90FAAD41F4BC69E51EB6C9D4C78B618D0A4B616478BD05E7A480950"
	NetworkType         = ""
	TestAccountKey      = "819F72066B17FFD71B8B4142C5AEAE4B997B0882ABDF2C263B02869382BD93A0"
)

func main() {
	ctx := context.Background()

	conf, err := sdk.NewConfig(ctx, baseUrls)
	if err != nil {
		panic(err)
	}
	ws, err := websocket.NewClient(ctx, conf)
	if err != nil {
		panic(err)
	}

	client := sdk.NewClient(nil, conf)

	actualNetworkType, err := client.Network.GetNetworkType(context.Background())
	if err != nil {
		fmt.Printf("Network.GetNetworkType returned error: %s", err)
		return
	}

	wg := new(sync.WaitGroup)
	go ws.Listen()

	customerAcc, err := sdk.NewAccountFromPrivateKey(TestAccountKey, actualNetworkType, client.GenerationHash())
	fmt.Printf("Main Acc address: %s \n", customerAcc.Address)
	fmt.Printf("Main harvester address: %s \n", customerAcc.Address)
	if err != nil {
		panic(fmt.Errorf("Customer account #0 returned error: %s", err))
	}

	customerAccRemote, err := client.NewAccount()

	if err != nil {
		panic(fmt.Errorf("Customer account #0 returned error: %s", err))
	}
	fmt.Printf("Customer Acc address: %s \n", customerAccRemote.Address)
	wg.Add(1)
	err = ws.AddUnconfirmedAddedHandlers(customerAcc.Address, func(transaction sdk.Transaction) bool {
		defer wg.Done()
		fmt.Printf("Main Acc: UnconfirmedAdded Tx Content: %s \n", transaction.GetAbstractTransaction().TransactionHash)
		return true
	})

	if err != nil {
		panic(err)
	}

	wg.Add(1)
	err = ws.AddUnconfirmedAddedHandlers(customerAccRemote.Address, func(transaction sdk.Transaction) bool {
		defer wg.Done()
		fmt.Printf("Remote acc: UnconfirmedAdded Tx Content: %s \n", transaction.GetAbstractTransaction().TransactionHash)
		return true
	})

	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(fmt.Errorf("Remote customer account returned error: %s", err))
	}

	//
	//// The confirmedAdded channel notifies when a transaction related to an
	//// address is included in a block. The message contains the transaction.

	wg.Add(1)
	err = ws.AddConfirmedAddedHandlers(customerAcc.Address, func(transaction sdk.Transaction) bool {
		defer wg.Done()
		fmt.Printf("Main Acc: ConfirmedAdded Tx Content: %s \n", transaction.GetAbstractTransaction().TransactionHash)
		fmt.Println("Successful transfer!")
		return true
	})

	if err != nil {
		panic(err)
	}

	wg.Add(2)
	err = ws.AddConfirmedAddedHandlers(customerAccRemote.Address, func(transaction sdk.Transaction) bool {
		defer wg.Done()
		fmt.Printf("Remote Acc: ConfirmedAdded Tx Content: %s \n", transaction.GetAbstractTransaction().TransactionHash)
		fmt.Println("Successful transfer!")
		return true
	})

	if err != nil {
		panic(err)
	}

	//The status channel notifies when a transaction related to an address rises an error.
	//The message contains the error message and the transaction hash.

	wg.Add(2)
	err = ws.AddStatusHandlers(customerAcc.Address, func(info *sdk.StatusInfo) bool {
		defer wg.Done()
		fmt.Printf("Main Account: Content: %v \n", info.Hash)
		fmt.Printf("Status: %s", info.Status)
		return true
	})

	if err != nil {
		panic(err)
	}

	wg.Add(1)
	err = ws.AddStatusHandlers(customerAccRemote.Address, func(info *sdk.StatusInfo) bool {
		defer wg.Done()
		fmt.Printf("Remote Account: Content: %v \n", info.Hash)
		fmt.Printf("Status: %s", info.Status)
		return true
	})

	if err != nil {
		panic(err)
	}
	/*transaction, err := client.NewAccountLinkTransaction(sdk.NewDeadline(time.Hour*1),
		customerAccRemote.PublicAccount,
		sdk.AccountLink)

	signedAccountLinKTransaction, err := customerAcc.Sign(transaction)

	if err != nil {
		panic(fmt.Errorf("Account link transaction signing returned error: %s", err))
	}
	restTx, err := client.Transaction.Announce(context.Background(), signedAccountLinKTransaction)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s\n", restTx)

	nodeLinkTransaction, err := client.NewNodeKeyLinkTransaction(sdk.NewDeadline(time.Hour*1),
		HarvesterNodeKey,
		sdk.AccountLink)

	signedNodeLinkTransaction, err := customerAcc.Sign(nodeLinkTransaction)

	if err != nil {
		panic(fmt.Errorf("Node link transaction signing returned error: %s", err))
	}

	restTx, err = client.Transaction.Announce(context.Background(), signedNodeLinkTransaction)
	if err != nil {
		panic(fmt.Errorf("Cannot announce node link: %s", err))
	}
	fmt.Printf("%s\n", restTx)
	*/
	/*
		key, err := client.Node.GetNodeUnlockedAccounts(context.Background())
		if err != nil {
			panic(fmt.Errorf("Cannot retrieve unlocked accounts: %s", err))
		}
		fmt.Printf("%v", key)
	*/

	harvestingAccount, err := sdk.NewAccountFromPrivateKey(HarvesterNodeKey, actualNetworkType, client.GenerationHash())
	fmt.Printf("harvest Acc address: %s \n", harvestingAccount.Address)

	message, err := sdk.NewPersistentHarvestingDelegationMessageFromPlainText(customerAccRemote.PrivateKey, harvestingAccount.KeyPair.PublicKey)
	persistentDelegationLinkTransaction, err := client.NewTransferTransaction(sdk.NewDeadline(time.Hour*1),
		harvestingAccount.Address,
		[]*sdk.Mosaic{},
		message)

	signedPersistentDelegationLinkTransaction, err := customerAcc.Sign(persistentDelegationLinkTransaction)

	if err != nil {
		panic(fmt.Errorf("Transfer transaction signing returned error: %s", err))
	}

	restTx, err := client.Transaction.Announce(context.Background(), signedPersistentDelegationLinkTransaction)
	if err != nil {
		panic(fmt.Errorf("Transfer transaction announcing returned error: %s", err))
	}
	fmt.Printf("%s\n", restTx)

	wg.Wait()
}
