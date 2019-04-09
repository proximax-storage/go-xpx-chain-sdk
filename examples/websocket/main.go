// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"math/big"
	"net/http"
	"sync"
	"time"
)

const (
//wsBaseUrl     = "ws://bcstage1.xpxsirius.io:3000/ws"
//baseUrl     = "http://bcstage1.xpxsirius.io:3000"
//networkType = sdk.PublicTest
//privateKey  = "809CD6699B7F38063E28F606BD3A8AECA6E13B1E688FE8E733D13DB843BC14B7"
)

//const (
//	wsBaseUrl     = "ws://192.168.88.41:3000/ws"
//	baseUrl     = "http://192.168.88.41:3000"
//	networkType = sdk.MijinTest
//	privateKey  = "A97B139EB641BCC841A610231870925EB301BA680D07BBCF9AEE83FAA5E9FB43"
//)

const (
	wsBaseUrl   = "ws://127.0.0.1:3000/ws"
	baseUrl     = "http://127.0.0.1:3000"
	networkType = sdk.MijinTest
	privateKey  = "A97B139EB641BCC841A610231870925EB301BA680D07BBCF9AEE83FAA5E9FB43"
)

// WebSockets make possible receiving notifications when a transaction or event occurs in the blockchain.
// The notification is received in real time without having to poll the API waiting for a reply.
func main() {

	destAccount, _ := sdk.NewAccountFromPrivateKey(privateKey, networkType)
	address := destAccount.PublicAccount.Address

	fmt.Println(fmt.Sprintf("destination address: %s", address.Address))

	wsc, err := sdk.NewCatapultWebSocketClient(wsBaseUrl)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	//Starting listening messages from websocket
	go wsc.Listen(&wg)

	if err := wsc.AddBlockHandlers(BlocksHandler1, BlocksHandler2); err != nil {
		panic(err)
	}

	//if err := wsc.AddConfirmedAddedHandlers(address, ConfirmedAddedHandler1, ConfirmedAddedHandler2); err != nil {
	//	panic(err)
	//}

	if err = wsc.AddUnconfirmedAddedHandlers(address, UnconfirmedAddedHandler1, UnconfirmedAddedHandler2); err != nil {
		panic(err)
	}

	//time.Sleep(time.Second * 5)

	//doTransferTransaction(address)

	//doBondedAggregateTransaction(address)

	wg.Wait()
}

// test publish transfer transaction
func doTransferTransaction(address *sdk.Address) {

	fmt.Println("start publishing transfer transaction")

	conf, err := sdk.NewConfig(baseUrl, networkType)
	if err != nil {
		panic(err)
	}

	acc, err := sdk.NewAccountFromPrivateKey(privateKey, networkType)

	// Use the default http client
	client := sdk.NewClient(http.DefaultClient, conf)

	ttx, err := sdk.NewTransferTransaction(
		sdk.NewDeadline(time.Hour*1),
		address,
		[]*sdk.Mosaic{sdk.Xem(10000000)},
		sdk.NewPlainMessage("my test transaction"),
		networkType,
	)

	stx, err := acc.Sign(ttx)
	if err != nil {
		panic(fmt.Errorf("TransaferTransaction signing returned error: %s", err))
	}

	restTx, err := client.Transaction.Announce(context.Background(), stx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", restTx)
	fmt.Printf("Content: \t\t%v\n", stx.Hash)
	fmt.Printf("Signer: \t%X\n\n", acc.KeyPair.PublicKey.Raw)

	fmt.Println("transfer transaction successfully published")
}

// test publish aggregated transaction
func doBondedAggregateTransaction(address *sdk.Address) {

	fmt.Println("start publishing bonded aggregated transaction")

	conf, err := sdk.NewConfig(baseUrl, networkType)
	if err != nil {
		panic(err)
	}

	acc, err := sdk.NewAccountFromPrivateKey(privateKey, networkType)

	// Use the default http client
	client := sdk.NewClient(http.DefaultClient, conf)

	ttx1, err := sdk.NewTransferTransaction(
		sdk.NewDeadline(time.Hour*1),
		address,
		[]*sdk.Mosaic{sdk.Xem(50)},
		sdk.NewPlainMessage("first transaction"),
		networkType,
	)

	if err != nil {
		panic(err)
	}

	ttx1.ToAggregate(acc.PublicAccount)

	ttx2, err := sdk.NewTransferTransaction(
		sdk.NewDeadline(time.Hour*1),
		address,
		[]*sdk.Mosaic{sdk.Xem(90)},
		sdk.NewPlainMessage("second transaction"),
		networkType,
	)

	if err != nil {
		panic(err)
	}

	ttx2.ToAggregate(acc.PublicAccount)

	bondedTx, err := sdk.NewBondedAggregateTransaction(
		sdk.NewDeadline(time.Hour*3),
		[]sdk.Transaction{ttx1, ttx2},
		networkType,
	)

	if err != nil {
		panic(err)
	}

	signedBondedTx, err := acc.Sign(bondedTx)
	if err != nil {
		panic(err)
	}

	lockFound, err := sdk.NewLockFundsTransaction(sdk.NewDeadline(time.Hour*3), sdk.XpxRelative(10), big.NewInt(240), signedBondedTx, networkType)
	if err != nil {
		panic(err)
	}

	signedLockFound, err := acc.Sign(lockFound)
	if err != nil {
		panic(err)
	}

	_, err = client.Transaction.Announce(context.Background(), signedLockFound)

	time.Sleep(time.Second * 30)

	resp, err := client.Transaction.AnnounceAggregateBonded(context.Background(), signedBondedTx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", resp)
	fmt.Printf("Content: \t\t%v\n", signedBondedTx.Hash)
	fmt.Printf("Signer: \t%X\n\n", acc.KeyPair.PublicKey.Raw)

	fmt.Println("bonded aggregated transaction successfully published")
}

// Examples of handler functions for different channels

func BlocksHandler1(blockInfo *sdk.BlockInfo) bool {
	fmt.Println("called BlockHandler1")
	//fmt.Println(blockInfo.String())
	return true
}

func BlocksHandler2(blockInfo *sdk.BlockInfo) bool {
	fmt.Println("called BlockHandler2")
	//fmt.Println(blockInfo.String())
	return true
}

func UnconfirmedAddedHandler1(tr sdk.Transaction) bool {
	fmt.Println("called UnconfirmedAddedHandler1")
	//fmt.Println(tr.String())
	return true
}

func UnconfirmedAddedHandler2(tr sdk.Transaction) bool {
	fmt.Println("called UnconfirmedAddedHandler2")
	//fmt.Println(tr.String())
	return false
}

func ConfirmedAddedHandler1(tr sdk.Transaction) bool {
	fmt.Println("called ConfirmedAddedHandler1")
	//fmt.Println(tr.String())
	return false
}

func ConfirmedAddedHandler2(tr sdk.Transaction) bool {
	fmt.Println("called ConfirmedAddedHandler2")
	//fmt.Println(tr.String())
	return true
}
