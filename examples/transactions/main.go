// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/proximax-storage/proximax-nem2-sdk-go/sdk"
	"time"
)

const (
	baseUrl            = "http://localhost:3000"
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

	customerAcc, err := sdk.NewAccountFromPrivateKey(customerPrivateKey, networkType)

	// timeout in milliseconds
	// 15000 ms = 15 seconds
	ws, err := sdk.NewConnectWs(baseUrl, 15000)
	if err != nil {
		panic(err)
	}

	fmt.Println("websocket negotiated uid:", ws.Uid)

	// The UnconfirmedAdded channel notifies when a transaction related to an
	// address is in unconfirmed state and waiting to be included in a block.
	// The message contains the transaction.
	chUnconfirmedAdded, _ := ws.Subscribe.UnconfirmedAdded(customerAcc.Address)
	go func() {
		for {
			data := <-chUnconfirmedAdded.Ch
			fmt.Printf("UnconfirmedAdded Tx Content: %v \n", data.GetAbstractTransaction().Hash)
			chUnconfirmedAdded.Unsubscribe()
		}
	}()
	//
	//// The confirmedAdded channel notifies when a transaction related to an
	//// address is included in a block. The message contains the transaction.
	chConfirmedAdded, _ := ws.Subscribe.ConfirmedAdded(customerAcc.Address)
	go func() {
		for {
			data := <-chConfirmedAdded.Ch
			fmt.Printf("ConfirmedAdded Tx Content: %v \n", data.GetAbstractTransaction().Hash)
			chConfirmedAdded.Unsubscribe()
			fmt.Println("Successful transfer!")
		}
	}()

	//The status channel notifies when a transaction related to an address rises an error.
	//The message contains the error message and the transaction hash.
	chStatus, _ := ws.Subscribe.Status(customerAcc.Address)

	go func() {
		for {
			data := <-chStatus.Ch
			chStatus.Unsubscribe()
			fmt.Printf("Content: %v \n", data.Hash)
			panic(fmt.Sprint("Status: ", data.Status))
		}
	}()

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
	chBlock, _ := ws.Subscribe.Block()

	for {
		data := <-chBlock.Ch
		fmt.Printf("Block received with height: %v \n", data.Height)
	}
}
