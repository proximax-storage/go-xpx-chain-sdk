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
	baseUrls = []string{"http://127.0.0.1:3000"}
	Hash, _  = sdk.StringToHash("86258172F90639811F2ABD055747D1E11B55A64B68AED2CEA9A34FBD6C0BE790")
)

// WebSockets make possible receiving notifications when a transaction or event occurs in the blockchain.
// The notification is received in real time without having to poll the API waiting for a reply.
func main() {
	ctx := context.Background()

	conf, err := sdk.NewConfig(ctx, baseUrls)
	if err != nil {
		panic(err)
	}

	ws, err := websocket.NewClient(conf)
	if err != nil {
		panic(err)
	}

	client := sdk.NewClient(nil, conf)

	customerAcc, err := client.NewAccount()
	wg := new(sync.WaitGroup)
	go ws.Listen(ctx)

	// The UnconfirmedAdded channel notifies when a transaction related to an
	// address is in unconfirmed state and waiting to be included in a block.
	// The message contains the transaction.

	unconfSub, _, err := ws.NewUnConfirmedAddedSubscription(customerAcc.Address)
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		transaction := <-unconfSub
		fmt.Printf("UnconfirmedAdded Tx Content: %s \n", transaction.GetAbstractTransaction().TransactionHash)
	}()

	//
	//// The confirmedAdded channel notifies when a transaction related to an
	//// address is included in a block. The message contains the transaction.

	confSub, _, err := ws.NewConfirmedAddedSubscription(customerAcc.Address)
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		transaction := <-confSub
		fmt.Printf("ConfirmedAdded Tx Content: %s \n", transaction.GetAbstractTransaction().TransactionHash)
		fmt.Println("Successful transfer!")
	}()

	//The status channel notifies when a transaction related to an address rises an error.
	//The message contains the error message and the transaction hash.

	statusSub, _, err := ws.NewStatusSubscription(customerAcc.Address)
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		info := <-statusSub
		fmt.Printf("Content: %v \n", info.Hash)
		panic(fmt.Sprint("Status: ", info.Status))
	}()

	time.Sleep(time.Second * 5)

	executorAcc, err := client.NewAccount()
	verifierAcc, err := client.NewAccount()
	println("Customer PublickKey:", customerAcc.PublicAccount.PublicKey)
	println("Executor PublickKey:", executorAcc.PublicAccount.PublicKey)
	println("Verifier PublickKey:", verifierAcc.PublicAccount.PublicKey)

	mctx, err := client.NewModifyContractTransaction(
		sdk.NewDeadline(time.Hour*1),
		sdk.Duration(2),
		Hash,
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

	blockSub, _, err := ws.NewBlockSubscription()
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		info := <-blockSub
		fmt.Printf("Block received with height: %v \n", info.Height)
	}()

	wg.Wait()

	if err := ws.Close(); err != nil {
		panic(err)
	}
}
