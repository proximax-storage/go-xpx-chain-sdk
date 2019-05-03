// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket"
	"sync"
	"time"
)

const (
	baseUrl            = "http://127.0.0.1:3000"
	networkType        = sdk.MijinTest
	customerPrivateKey = "0F3CC33190A49ABB32E7172E348EA927F975F8829107AAA3D6349BB10797D4F6"
	executorPrivateKey = "68B3FBB18729C1FDE225C57F8CE080FA828F0067E451A3FD81FA628842B0B763"
	verifierPrivateKey = "CF893FFCC47C33E7F68AB1DB56365C156B0736824A0C1E273F9E00B8DF8F01EB"
	hash               = sdk.Hash("037AFE3810F034C237E15B098D00A3D703B9558142BBCC561A197F72412903A6")
	multisig           = "3FE21823C74BAFEAA99100767D0AE573AD90FA362F21A5C2F7A5BDC0840E9660"
)

// WebSockets make possible receiving notifications when a transaction or event occurs in the blockchain.
// The notification is received in real time without having to poll the API waiting for a reply.
func main() {

	conf, err := sdk.NewConfig(baseUrl, networkType)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	customerAcc, err := sdk.NewAccountFromPrivateKey(customerPrivateKey, networkType)

	ws, err := websocket.NewClient(ctx, conf)
	if err != nil {
		panic(err)
	}

	wg := new(sync.WaitGroup)
	go ws.Listen()

	// The UnconfirmedAdded channel notifies when a transaction related to an
	// address is in unconfirmed state and waiting to be included in a block.
	// The message contains the transaction.

	wg.Add(1)
	err = ws.AddUnconfirmedAddedHandlers(customerAcc.Address, func(transaction sdk.Transaction) bool {
		defer wg.Done()
		fmt.Printf("UnconfirmedAdded Tx Content: %v \n", transaction.GetAbstractTransaction().Hash)
		return true
	})

	if err != nil {
		panic(err)
	}

	//
	//// The confirmedAdded channel notifies when a transaction related to an
	//// address is included in a block. The message contains the transaction.

	wg.Add(1)
	err = ws.AddConfirmedAddedHandlers(customerAcc.Address, func(transaction sdk.Transaction) bool {
		defer wg.Done()
		fmt.Printf("ConfirmedAdded Tx Content: %v \n", transaction.GetAbstractTransaction().Hash)
		fmt.Println("Successful transfer!")
		return true
	})

	if err != nil {
		panic(err)
	}

	//The status channel notifies when a transaction related to an address rises an error.
	//The message contains the error message and the transaction hash.

	wg.Add(1)
	err = ws.AddStatusHandlers(customerAcc.Address, func(info *sdk.StatusInfo) bool {
		defer wg.Done()
		fmt.Printf("Content: %v \n", info.Hash)
		panic(fmt.Sprint("Status: ", info.Status))
		return true
	})

	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 5)

	// Use the default http client
	client := sdk.NewClient(nil, conf)

	executorAcc, err := sdk.NewAccountFromPrivateKey(executorPrivateKey, networkType)
	verifierAcc, err := sdk.NewAccountFromPrivateKey(verifierPrivateKey, networkType)
	println("Customer PublickKey:", customerAcc.PublicAccount.PublicKey)
	println("Executor PublickKey:", executorAcc.PublicAccount.PublicKey)
	println("Verifier PublickKey:", verifierAcc.PublicAccount.PublicKey)

	mctx, err := sdk.NewModifyContractTransaction(
		sdk.NewDeadline(time.Hour*1),
		2,
		hash.String(),
		[]*sdk.MultisigCosignatoryModification{
			{
				sdk.Add,
				customerAcc.PublicAccount,
			},
		},
		[]*sdk.MultisigCosignatoryModification{
			{
				sdk.Add,
				executorAcc.PublicAccount,
			},
		},
		[]*sdk.MultisigCosignatoryModification{
			{
				sdk.Add,
				verifierAcc.PublicAccount,
			},
		},
		networkType,
	)

	stx, err := customerAcc.Sign(mctx)
	if err != nil {
		panic(fmt.Errorf("TransaferTransaction signing returned error: %s", err))
	}

	// Get the chain height
	restTx, err := client.Transaction.Announce(context.Background(), stx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", restTx)
	fmt.Printf("Content: \t\t%v\n", stx.Hash)
	fmt.Printf("Signer: \t%X\n\n", customerAcc.KeyPair.PublicKey.Raw)

	// The block channel notifies for every new block.
	// The message contains the block information.

	wg.Add(1)
	err = ws.AddBlockHandlers(func(info *sdk.BlockInfo) bool {
		defer wg.Done()
		fmt.Printf("Block received with height: %v \n", info.Height)
		return true
	})

	if err != nil {
		panic(err)
	}

	wg.Wait()

	if err := ws.Close(); err != nil {
		panic(err)
	}
}
