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
)

const (
	create = "create"
	change = "change"
)

var (
	ErrNoUrl          = errors.New("url is not provided")
	ErrUnknownCommand = errors.New("unknown command")

	ErrZeroBeta                  = errors.New("zero beta")
	ErrZeroAlpha                 = errors.New("zero alpha")
	ErrZeroWindowSize            = errors.New("zero window size")
	ErrUnknownMosaicName         = errors.New("unknown mosaic name")
	ErrZeroSlashingPeriod        = errors.New("zero slashing period")
	ErrZeroCurrencyDeposit       = errors.New("zero currency deposit")
	ErrEmptySlashingAccount      = errors.New("empty slashing account")
	ErrZeroInitialMosaicsMinting = errors.New("zero initial mosaics minting")
)

var sender *sdk.Account
var ownerAccount *sdk.PublicAccount

func main() {
	// common
	url := flag.String("url", "http://127.0.0.1:3000", "ProximaX Chain REST Url")
	feeStrategy := flag.String("feeStrategy", tools.MiddleFeeStrategy, "fee calculation strategy (low, middle, high)")
	senderPrivateKey := flag.String("sender", "", "transaction sender private key")
	ownerPublicKey := flag.String("owner", "", "liquidity provider owner public key (empty if owner is not multisig)")

	providerMosaicName := flag.String("mosaic", "", "Name of a mosaic (storage, streaming or sc units)")

	// create lp
	slashingAccount := flag.String("slashingAcc", "", "slashing account public key")
	currencyDeposit := flag.Uint64("deposit", 100000, "Amount of currency deposit")
	initialMosaicsMinting := flag.Uint64("initial", 100000, "Amount of initial mosaics minting")
	slashingPeriod := flag.Uint("slashingPeriod", 500, "slashing period")
	windowSize := flag.Uint("ws", 5, "window size")
	alpha := flag.Uint("a", 500, "alpha")
	beta := flag.Uint("b", 500, "beta")

	// change lp
	currencyBalanceIncrease := flag.Bool("currencyBalanceIncrease", false, "currency balance increase")
	currencyBalanceChange := flag.Uint64("currencyBalanceChange", 0, "currency balance change")
	mosaicBalanceIncrease := flag.Bool("mosaicBalanceIncrease", false, "mosaic balance increase")
	mosaicBalanceChange := flag.Uint64("mosaicBalanceChange", 0, "mosaic balance change")

	flag.Parse()

	if *url == "" {
		fmt.Println(ErrNoUrl)
		os.Exit(1)
	}

	ctx := context.Background()
	cfg, err := sdk.NewConfig(ctx, []string{*url})
	if err != nil {
		fmt.Printf("Cannot create client: %s\n", err)
		os.Exit(1)
	}
	cfg.FeeCalculationStrategy = tools.ParseFeeStrategy(feeStrategy)

	client := sdk.NewClient(http.DefaultClient, cfg)
	ws, err := websocket.NewClient(cfg)
	if err != nil {
		fmt.Printf("Cannot create websocket client: %s\n", err)
		os.Exit(1)
	}

	if senderPrivateKey == nil || *senderPrivateKey == "" {
		fmt.Println("Transaction sender not specified")
		os.Exit(1)
	}

	sender, err = client.NewAccountFromPrivateKey(*senderPrivateKey)
	if err != nil {
		fmt.Printf("Cannot create sender account from private key: %s\n", err)
		os.Exit(1)
	}

	if ownerPublicKey != nil && *ownerPublicKey != "" {
		ownerAccount, err = client.NewAccountFromPublicKey(*ownerPublicKey)
		if err != nil {
			fmt.Printf("Cannot create owner account from public key: %s\n", err)
			os.Exit(1)
		}
	}

	var mosacInfo *sdk.MosaicInfo
	switch strings.ToLower(*providerMosaicName) {
	case "storage":
		mosacInfo, err = client.Resolve.GetMosaicInfoByAssetId(ctx, sdk.StorageNamespaceId)
	case "streaming":
		mosacInfo, err = client.Resolve.GetMosaicInfoByAssetId(ctx, sdk.StreamingNamespaceId)
	case "supercontract":
		fallthrough
	case "sc":
		mosacInfo, err = client.Resolve.GetMosaicInfoByAssetId(ctx, sdk.SuperContractNamespaceId)
	default:
		err = ErrUnknownMosaicName
	}
	if err != nil {
		fmt.Printf("%s: %s\n", err, *providerMosaicName)
		os.Exit(1)
	}

	arg := os.Args[len(os.Args)-1]
	switch arg {
	case create:
		err := newLiquidityProvider(
			client,
			ws,
			cfg,
			mosacInfo.MosaicId,
			*currencyDeposit,
			*initialMosaicsMinting,
			uint32(*slashingPeriod),
			uint16(*windowSize),
			*slashingAccount,
			uint32(*alpha),
			uint32(*beta),
		)
		if err != nil {
			fmt.Printf("Cannot create liquidity provider: %s\n", err)
			os.Exit(1)
		}

		return
	case change:
		err := manualRateChange(
			client,
			ws,
			cfg,
			mosacInfo.MosaicId,
			*currencyBalanceIncrease,
			*currencyBalanceChange,
			*mosaicBalanceIncrease,
			*mosaicBalanceChange,
		)
		if err != nil {
			fmt.Printf("Cannot change liquidity provider: %s\n", err)
			os.Exit(1)
		}

		return
	default:
		fmt.Printf("%s: %s\n", ErrUnknownCommand, arg)
		os.Exit(1)
	}
}

func newLiquidityProvider(
	client *sdk.Client,
	ws websocket.CatapultClient,
	cfg *sdk.Config,
	mosaicId *sdk.MosaicId,
	currencyDeposit uint64,
	initialMosaicsMinting uint64,
	slashingPeriod uint32,
	windowSize uint16,
	slashingAccount string,
	alpha uint32,
	beta uint32) error {

	if currencyDeposit == 0 {
		return ErrZeroCurrencyDeposit
	}

	if slashingPeriod == 0 {
		return ErrZeroSlashingPeriod
	}

	if initialMosaicsMinting == 0 {
		return ErrZeroInitialMosaicsMinting
	}

	if windowSize == 0 {
		return ErrZeroWindowSize
	}

	if slashingAccount == "" {
		return ErrEmptySlashingAccount
	}

	if alpha == 0 {
		return ErrZeroAlpha
	}

	if beta == 0 {
		return ErrZeroBeta
	}

	sa, err := client.NewAccountFromPublicKey(slashingAccount)
	if err != nil {
		return err
	}

	tx, err := client.NewCreateLiquidityProviderTransaction(
		sdk.NewDeadline(time.Hour),
		mosaicId,
		sdk.Amount(currencyDeposit),
		sdk.Amount(initialMosaicsMinting),
		slashingPeriod,
		windowSize,
		sa,
		alpha,
		beta,
	)
	if err != nil {
		return err
	}

	return announce(client, context.Background(), cfg, ws, tx)
}

func manualRateChange(
	client *sdk.Client,
	ws websocket.CatapultClient,
	cfg *sdk.Config,
	mosaicId *sdk.MosaicId,
	currencyBalanceIncrease bool,
	currencyBalanceChange uint64,
	mosaicBalanceIncrease bool,
	mosaicBalanceChange uint64) error {

	tx, err := client.NewManualRateChangeTransaction(
		sdk.NewDeadline(time.Hour),
		mosaicId,
		currencyBalanceIncrease,
		sdk.Amount(currencyBalanceChange),
		mosaicBalanceIncrease,
		sdk.Amount(mosaicBalanceChange),
	)
	if err != nil {
		return err
	}

	return announce(client, context.Background(), cfg, ws, tx)
}

func announce(client *sdk.Client, ctx context.Context, cfg *sdk.Config, ws websocket.CatapultClient, tx sdk.Transaction) error {
	var err error
	if ownerAccount != nil {
		tx.GetAbstractTransaction().ToAggregate(ownerAccount)
		var aggregateTx *sdk.AggregateTransaction
		aggregateTx, err = client.NewBondedAggregateTransaction(
			sdk.NewDeadline(time.Hour*48),
			[]sdk.Transaction{tx},
		)
		if err != nil {
			return err
		}

		signedABT, err := sender.Sign(aggregateTx)
		if err != nil {
			return err
		}

		lockFundsTx, err := client.NewLockFundsTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.XpxRelative(10),
			sdk.Duration(11519),
			signedABT,
		)
		if err != nil {
			return err
		}

		signedLockFundsTx, err := sender.Sign(lockFundsTx)
		if err != nil {
			return err
		}

		lockFundsTxHash, err := client.Transaction.Announce(ctx, signedLockFundsTx)
		if err != nil {
			return err
		}

		fmt.Println("LockFundsTx Hash:", lockFundsTxHash)

		time.Sleep(90 * time.Second)

		hash, err := client.Transaction.AnnounceAggregateBonded(ctx, signedABT)
		if err != nil {
			return err
		}

		fmt.Println("ABT hash:", hash)
	} else {
		signedTx, err := sender.Sign(tx)
		if err != nil {
			return err
		}

		_, err = client.Transaction.Announce(ctx, signedTx)
		if err != nil {
			return err
		}
	}
	return nil
}
