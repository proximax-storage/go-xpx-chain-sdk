package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"
	"github.com/proximax-storage/go-xpx-chain-sdk/tools"

	sync "github.com/proximax-storage/go-xpx-chain-sync"
)

var (
	ErrNoUrl        = errors.New("url is not provided")
	ErrEmptyMosaic  = errors.New("empty mosaic id")
	ErrZeroCapacity = errors.New("capacity is zero")
	ErrEmptyKey     = errors.New("sender or receiver key is not provided")
)

func main() {
	url := flag.String("url", "http://127.0.0.1:3000", "ProximaX Chain REST Url")
	feeStrategy := flag.String("feeStrategy", tools.MiddleFeeStrategy, "fee calculation strategy (low, middle, high)")
	sender := flag.String("sender", "", "Sender private key")
	receiver := flag.String("receiver", "", "Receiver public key")
	mosaicId := flag.String("mosaic", "", "HEX provider mosaic id, e.g. 6C5D687508AC9D75")
	amount := flag.Uint64("amount", 0, "Amount of transfer mosaic")
	flag.Parse()

	if err := transfer(*url, tools.ParseFeeStrategy(feeStrategy), *sender, *receiver, *mosaicId, *amount); err != nil {
		fmt.Printf("Transfer failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Transfered successfully!")
}

func transfer(url string, feeStrategy sdk.FeeCalculationStrategy, sender, receiver, mosaicId string, amount uint64) error {
	if url == "" {
		return ErrNoUrl
	}

	if sender == "" || receiver == "" {
		return ErrEmptyKey
	}

	if mosaicId == "" {
		return ErrEmptyMosaic
	}

	if amount == 0 {
		return ErrZeroCapacity
	}

	ctx := context.Background()
	cfg, err := sdk.NewConfig(ctx, []string{url})
	if err != nil {
		return err
	}

	cfg.FeeCalculationStrategy = feeStrategy
	client := sdk.NewClient(http.DefaultClient, cfg)

	ws, err := websocket.NewClient(ctx, cfg)
	if err != nil {
		return err
	}

	senderAccount, err := client.NewAccountFromPrivateKey(sender)
	if err != nil {
		return err
	}

	receiverAccount, err := client.NewAccountFromPublicKey(receiver)
	if err != nil {
		return err
	}

	mId, err := tools.MosaicIdFromString(mosaicId)
	if err != nil {
		return err
	}

	transferTx, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		receiverAccount.Address,
		[]*sdk.Mosaic{
			{
				AssetId: mId,
				Amount:  sdk.Amount(amount),
			},
		},
		sdk.NewPlainMessage(""),
	)
	if err != nil {
		return err
	}

	res, err := sync.Announce(ctx, cfg, ws, senderAccount, transferTx)
	if err != nil {
		return err
	}

	return res.Err()
}
