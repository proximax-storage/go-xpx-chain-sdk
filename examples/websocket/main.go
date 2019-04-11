// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket"
	"math/big"
	"net/http"
	"sync"
	"time"
)

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

	wsc, err := websocket.NewCatapultWebSocketClient(wsBaseUrl)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	//Starting listening messages from websocket
	wg.Add(1)
	go wsc.Listen(&wg)

	//if err := wsc.AddBlockHandlers(BlocksHandler1, BlocksHandler2); err != nil {
	//	panic(err)
	//}

	//if err := wsc.AddConfirmedAddedHandlers(address, ConfirmedAddedHandler1, ConfirmedAddedHandler2); err != nil {
	//	panic(err)
	//}

	//if err := wsc.AddUnconfirmedAddedHandlers(address, UnconfirmedAddedHandler1, UnconfirmedAddedHandler2); err != nil {
	//	panic(err)
	//}

	//if err := wsc.AddUnconfirmedRemovedHandlers(address, UnconfirmedRemovedHandler1, UnconfirmedRemovedHandler2); err != nil {
	//	panic(err)
	//}

	//if err := wsc.AddPartialAddedHandlers(address, PartialAddedHandler1, PartialAddedHandler2); err != nil {
	//	panic(err)
	//}

	//if err := wsc.AddPartialRemovedHandlers(address, PartialRemovedHandler1, PartialRemovedHandler2); err != nil {
	//	panic(err)
	//}

	//if err := wsc.AddStatusHandlers(address, StatusHandler1, StatusHandler2); err != nil {
	//	panic(err)
	//}

	//if err := wsc.AddCosignatureHandlers(address, CosignatureHandler1, CosignatureHandler2); err != nil {
	//	panic(err)
	//}

	time.Sleep(time.Second * 5)

	//doTransferTransaction(address)

	doBondedAggregateTransaction(address)

	wg.Wait()
}

// publish transfer transaction
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

// publish aggregated transaction
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

// examples of handler functions for different websocket topics

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
	return false
}

func UnconfirmedAddedHandler2(tr sdk.Transaction) bool {
	fmt.Println("called UnconfirmedAddedHandler2")
	//fmt.Println(tr.String())
	return false
}

func UnconfirmedRemovedHandler1(removed *sdk.UnconfirmedRemoved) bool {
	fmt.Println("called UnconfirmedRemovedHandler1")
	//fmt.Println(tr.String())
	return false
}

func UnconfirmedRemovedHandler2(removed *sdk.UnconfirmedRemoved) bool {
	fmt.Println("called UnconfirmedRemovedHandler2")
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
	return false
}

func PartialAddedHandler1(tr *sdk.AggregateTransaction) bool {
	fmt.Println("called PartialAddedHandler1")
	//fmt.Println(tr.String())
	return false
}

func PartialAddedHandler2(tr *sdk.AggregateTransaction) bool {
	fmt.Println("called PartialAddedHandler2")
	//fmt.Println(tr.String())
	return false
}

func PartialRemovedHandler1(i *sdk.PartialRemovedInfo) bool {
	fmt.Println("called PartialRemovedHandler1")
	//fmt.Println(tr.String())
	return false
}

func PartialRemovedHandler2(i *sdk.PartialRemovedInfo) bool {
	fmt.Println("called PartialRemovedHandler1")
	//fmt.Println(tr.String())
	return false
}

func StatusHandler1(s *sdk.StatusInfo) bool {
	fmt.Println("called StatusHandler1")
	//fmt.Println(tr.String())
	return false
}

func StatusHandler2(s *sdk.StatusInfo) bool {
	fmt.Println("called StatusHandler2")
	//fmt.Println(tr.String())
	return false
}

func CosignatureHandler1(tr *sdk.SignerInfo) bool {
	fmt.Println("called CosignatureHandler1")
	//fmt.Println(tr.String())
	return false
}

func CosignatureHandler2(tr *sdk.SignerInfo) bool {
	fmt.Println("called CosignatureHandler2")
	//fmt.Println(tr.String())
	return false
}
