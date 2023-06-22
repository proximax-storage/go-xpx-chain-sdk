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

const (
	low    = "low"
	middle = "middle"
	high   = "high"
)

var (
	ErrNoUrl                  = errors.New("url is not provided")
	ErrZeroCapacity           = errors.New("capacity is zero")
	ErrUnknownFeeStrategy     = errors.New("unknown fee calculation strategy")
	ErrNoReplicatorPrivateKey = errors.New("replicator private key is not provided")
)

func main() {
	url := flag.String("url", "http://127.0.0.1:3000", "ProximaX Chain REST Url")
	capacity := flag.Uint64("capacity", 0, "capacity of replicator")
	feeStrategy := flag.String("feeStrategy", middle, "fee calculation strategy (low, middle, high)")
	replicatorPrivateKey := flag.String("privateKey", "", "Replicator private key")
	flag.Parse()

	if err := onboard(*url, *replicatorPrivateKey, tools.ParseFeeStrategy(feeStrategy), *capacity); err != nil {
		fmt.Printf("Replicator onboarding failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Replicator onboarded successfully!!!")
}

func onboard(url, replicatorPrivateKey string, feeStrategy sdk.FeeCalculationStrategy, capacity uint64) error {
	if url == "" {
		return ErrNoUrl
	}

	if replicatorPrivateKey == "" {
		return ErrNoReplicatorPrivateKey
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

	replicatorOnboardingTx, err := client.NewReplicatorOnboardingTransaction(
		sdk.NewDeadline(time.Hour),
		sdk.Amount(capacity),
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
