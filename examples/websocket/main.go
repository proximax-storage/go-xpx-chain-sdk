// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket"
	"time"
)

const (
	privateKey = "A97B139EB641BCC841A610231870925EB301BA680D07BBCF9AEE83FAA5E9FB43"
)

var (
	//baseUrls = []string{"http://192.168.88.15:3000"}
	baseUrls = []string{"http://127.0.0.1:3000", "http://127.0.0.1:3001", "http://127.0.0.1:3002"}
)

// WebSockets make possible receiving notifications when a transaction or event occurs in the blockchain.
// The notification is received in real time without having to poll the API waiting for a reply.
func main() {
	cfg, err := sdk.NewDefaultConfig(baseUrls)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	client := sdk.NewClient(nil, cfg)
	err = client.SetupConfigFromRest(ctx)
	if err != nil {
		panic(err)
	}

	wsc, err := websocket.NewClient(ctx, cfg)
	if err != nil {
		panic(err)
	}

	//Starting listening messages from websocket
	go wsc.Listen()

	destAccount, _ := client.NewAccountFromPrivateKey(privateKey)
	address := destAccount.PublicAccount.Address

	fmt.Println(fmt.Sprintf("destination address: %s", address.Address))

	// Register handlers functions for needed topics

	if err := wsc.AddBlockHandlers(BlocksHandler1, BlocksHandler2); err != nil {
		panic(err)
	}

	if err := wsc.AddConfirmedAddedHandlers(address, ConfirmedAddedHandler1, ConfirmedAddedHandler2); err != nil {
		panic(err)
	}

	if err := wsc.AddUnconfirmedAddedHandlers(address, UnconfirmedAddedHandler1, UnconfirmedAddedHandler2); err != nil {
		panic(err)
	}

	if err := wsc.AddUnconfirmedRemovedHandlers(address, UnconfirmedRemovedHandler1, UnconfirmedRemovedHandler2); err != nil {
		panic(err)
	}

	if err := wsc.AddPartialAddedHandlers(address, PartialAddedHandler1, PartialAddedHandler2); err != nil {
		panic(err)
	}

	if err := wsc.AddPartialRemovedHandlers(address, PartialRemovedHandler1, PartialRemovedHandler2); err != nil {
		panic(err)
	}

	if err := wsc.AddStatusHandlers(address, StatusHandler1, StatusHandler2); err != nil {
		panic(err)
	}

	if err := wsc.AddCosignatureHandlers(address, CosignatureHandler1, CosignatureHandler2); err != nil {
		panic(err)
	}

	//Running the goroutine which will close websocket connection and listening after 2 minutes.
	go func() {
		timer := time.NewTimer(time.Minute * 2)

		for range timer.C {
			if err := wsc.Close(); err != nil {
				panic(err)
			}
			return
		}
	}()

	//Publish test transactions
	doTransferTransaction(address, client)
	time.Sleep(time.Second * 30)
	doBondedAggregateTransaction(address, client)

	<-time.NewTimer(time.Minute * 5).C
}

// publish test transfer transaction
func doTransferTransaction(address *sdk.Address, client *sdk.Client) {

	fmt.Println("start publishing transfer transaction")
	acc, err := client.NewAccountFromPrivateKey(privateKey)

	ttx, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour*1),
		address,
		[]*sdk.Mosaic{sdk.Xem(10000000)},
		sdk.NewPlainMessage("my test transaction"),
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

// publish test aggregated transaction
func doBondedAggregateTransaction(address *sdk.Address, client *sdk.Client) {

	fmt.Println("start publishing bonded aggregated transaction")
	acc, err := client.NewAccountFromPrivateKey(privateKey)

	ttx1, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour*1),
		address,
		[]*sdk.Mosaic{sdk.Xem(50)},
		sdk.NewPlainMessage("first transaction"),
	)

	if err != nil {
		panic(err)
	}

	ttx1.ToAggregate(acc.PublicAccount)

	ttx2, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour*1),
		address,
		[]*sdk.Mosaic{sdk.Xem(90)},
		sdk.NewPlainMessage("second transaction"),
	)

	if err != nil {
		panic(err)
	}

	ttx2.ToAggregate(acc.PublicAccount)

	bondedTx, err := client.NewBondedAggregateTransaction(
		sdk.NewDeadline(time.Hour*3),
		[]sdk.Transaction{ttx1, ttx2},
	)

	if err != nil {
		panic(err)
	}

	signedBondedTx, err := acc.Sign(bondedTx)
	if err != nil {
		panic(err)
	}

	lockFound, err := client.NewLockFundsTransaction(sdk.NewDeadline(time.Hour*3), sdk.XpxRelative(10), sdk.Duration(240), signedBondedTx)
	if err != nil {
		panic(err)
	}

	signedLockFound, err := acc.Sign(lockFound)
	if err != nil {
		panic(err)
	}

	_, err = client.Transaction.Announce(context.Background(), signedLockFound)

	resp, err := client.Transaction.AnnounceAggregateBonded(context.Background(), signedBondedTx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", resp)
	fmt.Printf("Content: \t\t%v\n", signedBondedTx.Hash)
	fmt.Printf("Signer: \t%X\n\n", acc.KeyPair.PublicKey.Raw)

	fmt.Println("bonded aggregated transaction successfully published")
}

// Examples of handler functions for different websocket topics.
//
// If handler function will return true, this handler will be removed from handler storage for topic and
// won't be called nex time for topic message.
// If handler will return true, it will be called for next topic message.
// If all handlers for the topic will be removed, client will unsubscribe from topic on the websocket server

func BlocksHandler1(blockInfo *sdk.BlockInfo) bool {
	fmt.Println("called BlockHandler1")
	//fmt.Println(blockInfo.String())
	return true
}

func BlocksHandler2(blockInfo *sdk.BlockInfo) bool {
	fmt.Println("called BlockHandler2")
	//fmt.Println(blockInfo.String())
	return false
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
