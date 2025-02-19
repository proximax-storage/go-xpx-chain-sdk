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
	ErrNoUrl                  = errors.New("url is not provided")
	ErrNoReplicatorPrivateKey = errors.New("replicator private key is not provided")
	ErrNoDrivePublicKey       = errors.New("drive public key is not provided")
)

func main() {
	url := flag.String("url", "http://127.0.0.1:3000", "ProximaX Chain REST Url")
	feeStrategy := flag.String("feeStrategy", "middle", "fee calculation strategy (low, middle, high)")
	replicatorPrivateKey := flag.String("replicatorPrivateKey", "", "Replicator private key")
	drivePublicKey := flag.String("drivePublicKey", "", "Drive public key")
	flag.Parse()

	if err := offboard(*url, *replicatorPrivateKey, *drivePublicKey, tools.ParseFeeStrategy(feeStrategy)); err != nil {
		fmt.Printf("Replicator offboarding failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Replicator offboarded successfully!!!")
}

func offboard(url, replicatorPrivateKey string, drivePublicKey string, feeStrategy sdk.FeeCalculationStrategy) error {
	if url == "" {
		return ErrNoUrl
	}

	if replicatorPrivateKey == "" {
		return ErrNoReplicatorPrivateKey
	}

	if drivePublicKey == "" {
		return ErrNoDrivePublicKey
	}

	ctx := context.Background()
	cfg, err := sdk.NewConfig(ctx, []string{url})
	if err != nil {
		return err
	}

	cfg.FeeCalculationStrategy = feeStrategy
	client := sdk.NewClient(http.DefaultClient, cfg)

	ws, err := websocket.NewClient(cfg)
	if err != nil {
		return err
	}

	replicatorAccount, err := client.NewAccountFromPrivateKey(replicatorPrivateKey)
	if err != nil {
		return err
	}

	drivePublicAccount, err := client.NewAccountFromPublicKey(drivePublicKey)
	if err != nil {
		return err
	}

	replicatorOffboardingTx, err := client.NewReplicatorOffboardingTransaction(
		sdk.NewDeadline(time.Hour),
		drivePublicAccount,
	)
	if err != nil {
		return err
	}

	res, err := sync.Announce(ctx, cfg, ws, replicatorAccount, replicatorOffboardingTx)
	if err != nil {
		return err
	}

	return res.Err()
}
