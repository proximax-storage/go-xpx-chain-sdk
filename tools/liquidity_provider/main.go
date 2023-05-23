package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"
	sync "github.com/proximax-storage/go-xpx-chain-sync"
)

const (
	create = "create"
	change = "change"

	low    = "low"
	middle = "middle"
	high   = "high"
)

var (
	ErrNoUrl              = errors.New("url is not provided")
	ErrUnknownCommand     = errors.New("unknown command")
	ErrUnknownFeeStrategy = errors.New("unknown fee calculation strategy")

	ErrZeroBeta                  = errors.New("zero beta")
	ErrZeroAlpha                 = errors.New("zero alpha")
	ErrEmptyMosaic               = errors.New("empty mosaic id")
	ErrZeroWindowSize            = errors.New("zero window size")
	ErrZeroSlashingPeriod        = errors.New("zero slashing period")
	ErrZeroCurrencyDeposit       = errors.New("zero currency deposit")
	ErrEmptySlashingAccount      = errors.New("empty slashing account")
	ErrZeroInitialMosaicsMinting = errors.New("zero initial mosaics minting")
)

var sender *sdk.Account

func main() {
	// common
	url := flag.String("url", "http://127.0.0.1:3000", "ProximaX Chain REST Url")
	feeStrategy := flag.String("feeStrategy", middle, "fee calculation strategy (low, middle, high)")
	txSender := flag.String("sender", "", "transaction sender")

	providerMosaicId := flag.String("mosaic", "", "HEX provider mosaic id, e.g. 0x6C5D687508AC9D75")

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

	switch *feeStrategy {
	case low:
		cfg.FeeCalculationStrategy = sdk.LowCalculationStrategy
	case middle:
		cfg.FeeCalculationStrategy = sdk.MiddleCalculationStrategy
	case high:
		cfg.FeeCalculationStrategy = sdk.HighCalculationStrategy
	default:
		fmt.Printf("%s: %s\n", ErrUnknownFeeStrategy, *feeStrategy)
		os.Exit(1)
	}

	client := sdk.NewClient(http.DefaultClient, cfg)
	ws, err := websocket.NewClient(ctx, cfg)
	if err != nil {
		fmt.Printf("Cannot create websocket client: %s\n", err)
		os.Exit(1)
	}

	if txSender == nil || *txSender == "" {
		fmt.Println("Missed transaction sender account")
		os.Exit(1)
	}

	sender, err = client.NewAccountFromPrivateKey(*txSender)
	if err != nil {
		fmt.Printf("Cannot create txSender from private key: %s\n", err)
		os.Exit(1)
	}

	arg := os.Args[len(os.Args)-1]
	switch arg {
	case create:
		err := newLiquidityProvider(
			client,
			ws,
			cfg,
			*providerMosaicId,
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
			*providerMosaicId,
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
	mosaicId string,
	currencyDeposit uint64,
	initialMosaicsMinting uint64,
	slashingPeriod uint32,
	windowSize uint16,
	slashingAccount string,
	alpha uint32,
	beta uint32) error {

	mId, err := mosaicIdFromString(mosaicId)
	if err != nil {
		return err
	}

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
		mId,
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

	return announce(context.Background(), cfg, ws, tx)
}

func manualRateChange(
	client *sdk.Client,
	ws websocket.CatapultClient,
	cfg *sdk.Config,
	mosaicId string,
	currencyBalanceIncrease bool,
	currencyBalanceChange uint64,
	mosaicBalanceIncrease bool,
	mosaicBalanceChange uint64) error {

	mId, err := mosaicIdFromString(mosaicId)
	if err != nil {
		return err
	}

	tx, err := client.NewManualRateChangeTransaction(
		sdk.NewDeadline(time.Hour),
		mId,
		currencyBalanceIncrease,
		sdk.Amount(currencyBalanceChange),
		mosaicBalanceIncrease,
		sdk.Amount(mosaicBalanceChange),
	)
	if err != nil {
		return err
	}

	return announce(context.Background(), cfg, ws, tx)
}

func announce(ctx context.Context, cfg *sdk.Config, ws websocket.CatapultClient, tx sdk.Transaction) error {
	res, err := sync.Announce(ctx, cfg, ws, sender, tx)
	if err != nil {
		return err
	}

	return res.Err()
}

func mosaicIdFromString(mosaicId string) (*sdk.MosaicId, error) {
	if mosaicId == "" {
		return nil, ErrEmptyMosaic
	}

	mId, err := strconv.ParseUint(mosaicId, 16, 64)
	if err != nil {
		return nil, err
	}

	return sdk.NewMosaicId(mId)
}
