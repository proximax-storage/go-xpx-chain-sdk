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

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"
	"github.com/proximax-storage/go-xpx-chain-sdk/tools"
	sync "github.com/proximax-storage/go-xpx-chain-sync"
)

var (
	ErrNoUrl              = errors.New("url is not provided")
	ErrNoSignerPrivateKey = errors.New("signer private key is not provided")
	ErrNoReplicatorKeys   = errors.New("replicator public keys are not provided")
)

func main() {
	url := flag.String("url", "http://127.0.0.1:3000", "ProximaX Chain REST Url")
	feeStrategy := flag.String("feeStrategy", "middle", "Fee calculation strategy (low, middle, high)")
	signerPrivateKey := flag.String("signerPrivateKey", "", "Transaction signer private key")
	replicatorKeys := flag.String("replicatorKeys", "", "List of replicator public keys separated by whitespaces")
	flag.Parse()

	if err := rebuildTree(*url, *signerPrivateKey, *replicatorKeys, tools.ParseFeeStrategy(feeStrategy)); err != nil {
		fmt.Printf("Replicator tree rebuild failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Replicator tree rebuilt successfully!")
}

func rebuildTree(url, signerPrivateKey string, replicatorKeys string, feeStrategy sdk.FeeCalculationStrategy) error {
	if url == "" {
		return ErrNoUrl
	}

	if signerPrivateKey == "" {
		return ErrNoSignerPrivateKey
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

	signerAccount, err := client.NewAccountFromPrivateKey(signerPrivateKey)
	if err != nil {
		return err
	}

	replicatorAccounts, err := parseReplicatorKeys(replicatorKeys, client)
	if err != nil {
		return err
	}

	replicatorTreeRebuildTx, err := client.NewReplicatorTreeRebuildTransaction(
		sdk.NewDeadline(time.Hour),
		replicatorAccounts,
	)
	if err != nil {
		return err
	}

	res, err := sync.Announce(ctx, cfg, ws, signerAccount, replicatorTreeRebuildTx)
	if err != nil {
		return err
	}

	return res.Err()
}

func parseReplicatorKeys(keysStr string, client *sdk.Client) ([]*sdk.PublicAccount, error) {
	keysStr = strings.TrimSpace(keysStr)
	if keysStr == "" {
		return nil, ErrNoReplicatorKeys
	}

	keysStrArr := strings.Split(keysStr, " ")
	keys := make([]*sdk.PublicAccount, 0, len(keysStrArr))

	for _, keyStr := range keysStrArr {
		key, err := client.NewAccountFromPublicKey(keyStr)
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	return keys, nil
}
