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
	baseUrls             = []string{"http://127.0.0.1:3000"}
	HarvesterAccountKey  = "7AA907C3D80B3815BE4B4E1470DEEE8BB83BFEB330B9A82197603D09BA947230"
	HarvesterNodeKey     = "D6430327F90FAAD41F4BC69E51EB6C9D4C78B618D0A4B616478BD05E7A480950"
	NetworkType          = ""
	TestAccountKey       = "819F72066B17FFD71B8B4142C5AEAE4B997B0882ABDF2C263B02869382BD93A0" //"819F72066B17FFD71B8B4142C5AEAE4B997B0882ABDF2C263B02869382BD93A0"
	RemoteTestAccountKey = "bf1132005751e82f8cDc54a4961649df3e73dd00054e179f4cf5633e1e4bcb8d"
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

	if err != nil {
		panic(fmt.Errorf("Customer account #0 returned error: %s", err))
	}

	//customerAccRemote, err := sdk.NewAccountFromPrivateKey(RemoteTestAccountKey, actualNetworkType, client.GenerationHash())

	if err != nil {
		panic(fmt.Errorf("Customer account #0 returned error: %s", err))
	}
	wg.Add(12)
	err = ws.AddUnconfirmedAddedHandlers(customerAcc.Address, func(transaction sdk.Transaction) bool {
		defer wg.Done()
		fmt.Printf("UnconfirmedAdded Tx Content: %s \n", transaction.GetAbstractTransaction().TransactionHash)
		return true
	})

	if err != nil {
		panic(err)
	}

	err = ws.AddConfirmedAddedHandlers(customerAcc.Address, func(transaction sdk.Transaction) bool {
		defer wg.Done()
		fmt.Printf("ConfirmedAdded Tx Content: %s \n", transaction.GetAbstractTransaction().TransactionHash)
		fmt.Println("Successful transfer!")
		return true
	})

	if err != nil {
		panic(err)
	}

	err = ws.AddStatusHandlers(customerAcc.Address, func(info *sdk.StatusInfo) bool {
		defer wg.Done()
		fmt.Printf("Main Account: Content: %v \n", info.Hash)
		fmt.Printf("Status: %s", info.Status)
		return true
	})

	if err != nil {
		panic(err)
	}

	AnnounceAccountLink(client, customerAcc, customerAccRemote)
	AnnounceNodeLink(client, customerAcc, actualNetworkType)
	AnnounceTransferMessage(client, customerAcc, customerAccRemote, actualNetworkType)
	GetNodeUnlockedAccounts(client)

	wg.Wait()
}
func AnnounceAccountLink(client *sdk.Client, customerAcc *sdk.Account, customerAccRemote *sdk.Account) {
	transaction, err := client.NewAccountLinkTransaction(sdk.NewDeadline(time.Hour*1),
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
	fmt.Printf("%s\n", restTx)
}

func AnnounceNodeLink(client *sdk.Client, customerAcc *sdk.Account, actualNetworkType sdk.NetworkType) {
	harvestingAccount, err := sdk.NewAccountFromPrivateKey(HarvesterNodeKey, actualNetworkType, client.GenerationHash())
	nodeLinkTransaction, err := client.NewNodeKeyLinkTransaction(sdk.NewDeadline(time.Hour*1),
		harvestingAccount.PublicAccount.PublicKey,
		sdk.AccountLink)

	signedNodeLinkTransaction, err := customerAcc.Sign(nodeLinkTransaction)

	if err != nil {
		panic(fmt.Errorf("Node link transaction signing returned error: %s", err))
	}

	restTx, err := client.Transaction.Announce(context.Background(), signedNodeLinkTransaction)
	if err != nil {
		panic(fmt.Errorf("Cannot announce node link: %s", err))
	}
	fmt.Printf("%s\n", restTx)
}

func AnnounceTransferMessage(client *sdk.Client, customerAcc *sdk.Account, customerAccRemote *sdk.Account, actualNetworkType sdk.NetworkType) {
	harvestingAccount, err := sdk.NewAccountFromPrivateKey(HarvesterNodeKey, actualNetworkType, client.GenerationHash())
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
}

func GetNodeUnlockedAccounts(client *sdk.Client) {
	key, err := client.Node.GetNodeUnlockedAccounts(context.Background())
	if err != nil {
		panic(fmt.Errorf("Cannot retrieve unlocked accounts: %s", err))
	}
	fmt.Printf("%v", key)
}
