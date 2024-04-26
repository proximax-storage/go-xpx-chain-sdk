package main

import (
	"context"
	"crypto/rand"
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
	ErrNoUrl                  = errors.New("url is not provided")
	ErrZeroCapacity           = errors.New("capacity is zero")
	ErrNoReplicatorPrivateKey = errors.New("replicator private key is not provided")
	ErrNoNodeBootPrivateKey   = errors.New("node boot private key is not provided")
)

func main() {
	url := flag.String("url", "http://127.0.0.1:3000", "ProximaX Chain REST Url")
	capacity := flag.Uint64("capacity", 0, "capacity of replicator")
	feeStrategy := flag.String("feeStrategy", "middle", "fee calculation strategy (low, middle, high)")
	replicatorPrivateKey := flag.String("replicatorPrivateKey", "", "Replicator private key")
	nodeBootPrivateKey := flag.String("nodeBootPrivateKey", "", "Node boot private key")
	flag.Parse()

	if err := onboard(*url, *replicatorPrivateKey, *nodeBootPrivateKey, tools.ParseFeeStrategy(feeStrategy), *capacity); err != nil {
		fmt.Printf("Replicator onboarding failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Replicator onboarded successfully!!!")
}

func onboard(url, replicatorPrivateKey string, nodeBootPrivateKey string, feeStrategy sdk.FeeCalculationStrategy, capacity uint64) error {
	if url == "" {
		return ErrNoUrl
	}

	if replicatorPrivateKey == "" {
		return ErrNoReplicatorPrivateKey
	}

	if nodeBootPrivateKey == "" {
		return ErrNoNodeBootPrivateKey
	}

	if capacity == 0 {
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

	replicatorAccount, err := client.NewAccountFromPrivateKey(replicatorPrivateKey)
	if err != nil {
		return err
	}

	nodeAccount, err := client.NewAccountFromPrivateKey(nodeBootPrivateKey)
	if err != nil {
		return err
	}

	var message sdk.Hash
	_, err = rand.Read(message[:])
	if err != nil {
		return err
	}

	messageSignature, err := nodeAccount.SignData(message[:])
	if err != nil {
		return err
	}

	replicatorOnboardingTx, err := client.NewReplicatorOnboardingTransaction(
		sdk.NewDeadline(time.Hour),
		sdk.Amount(capacity),
		nodeAccount.PublicAccount,
		&message,
		messageSignature,
	)
	if err != nil {
		return err
	}

	res, err := sync.Announce(ctx, cfg, ws, replicatorAccount, replicatorOnboardingTx)
	if err != nil {
		return err
	}

	return res.Err()
}
