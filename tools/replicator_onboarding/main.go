package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"

	sync "github.com/proximax-storage/go-xpx-chain-sync"
)

var (
	ErrNoUrl                  = errors.New("url is not provided")
	ErrZeroCapacity           = errors.New("capacity is zero")
	ErrNoReplicatorPrivateKey = errors.New("replicator private key is not provided")
)

func main() {
	fmt.Println(`	
		░░░░░▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░
		░░░▓▓▓▓▓▓▒▒▒▒▒▒▓▓░░░░░░░
		░░▓▓▓▓▒░░▒▒▓▓▒▒▓▓▓▓░░░░░
		░▓▓▓▓▒░░▓▓▓▒▄▓░▒▄▄▄▓░░░░
		▓▓▓▓▓▒░░▒▀▀▀▀▒░▄░▄▒▓▓░░░
		▓▓▓▓▓▒░░▒▒▒▒▒▓▒▀▒▀▒▓▒▓░░
		▓▓▓▓▓▒▒░░░▒▒▒░░▄▀▀▀▄▓▒▓░
		▓▓▓▓▓▓▒▒░░░▒▒▓▀▄▄▄▄▓▒▒▒▓
		░▓█▀▄▒▓▒▒░░░▒▒░░▀▀▀▒▒▒▒░
		░░▓█▒▒▄▒▒▒▒▒▒▒░░▒▒▒▒▒▒▓░
		░░░▓▓▓▓▒▒▒▒▒▒▒▒░░░▒▒▒▓▓░
		░░░░░▓▓▒░░▒▒▒▒▒▒▒▒▒▒▒▓▓░
		░░░░░░▓▒▒░░░░▒▒▒▒▒▒▒▓▓░░
	`)

	url := flag.String("url", "http://127.0.0.1:3000", "ProximaX Chain REST Url")
	capacity := flag.Uint64("capacity", 0, "capacity of replicator")
	replicatorPrivateKey := flag.String("privateKey", "", "Replicator private key")
	cosignersPrivateKey := flag.String("cosigners", "", "Cosigners private keys separated by ','")
	flag.Parse()

	if err := onboard(*url, *replicatorPrivateKey, *cosignersPrivateKey, *capacity); err != nil {
		fmt.Printf("Replicator onboarding failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Replicator onboarded successfully!!!")
}

func onboard(url, replicatorPrivateKey, cosignersStr string, capacity uint64) error {
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

	cfg.FeeCalculationStrategy = sdk.MiddleCalculationStrategy
	client := sdk.NewClient(http.DefaultClient, cfg)

	ws, err := websocket.NewClient(ctx, cfg)
	if err != nil {
		return err
	}

	cosigners := strings.Split(cosignersStr, ",")
	cosignersAccs := make([]*sdk.Account, len(cosigners))
	for i, cosigner := range cosigners {
		cosignersAccs[i], err = client.NewAccountFromPrivateKey(cosigner)
		if err != nil {
			return err
		}
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

	switch len(cosigners) {
	case 0:
		res, err := sync.Announce(ctx, cfg, ws, replicatorAccount, replicatorOnboardingTx)
		if err != nil {
			return err
		}

		return res.Err()
	case 1:
		res, err := sync.Announce(ctx, cfg, ws, cosignersAccs[0], replicatorOnboardingTx)
		if err != nil {
			return err
		}

		return res.Err()
	default:
		gp, ctx := errgroup.WithContext(ctx)

		res, err := sync.Announce(ctx, cfg, ws, replicatorAccount, replicatorOnboardingTx)
		if err != nil {
			return err
		}

		if res.Err() != nil {
			return res.Err()
		}

		for _, acc := range cosignersAccs[1:] {
			gp.Go(func() error {
				syncer, err := sync.NewTransactionSyncer(ctx, cfg, acc, sync.WithWsClient(ws), sync.WithClient(client))
				if err != nil {
					return err
				}

				return syncer.CoSign(ctx, res.Hash(), false)
			})
		}

		return nil
	}
}
