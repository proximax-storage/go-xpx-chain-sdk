package tools

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

const (
	LowFeeStrategy    = "low"
	MiddleFeeStrategy = "middle"
	HighFeeStrategy   = "high"
)

var (
	ErrEmptyMosaic        = errors.New("empty mosaic id")
	ErrUnknownFeeStrategy = errors.New("unknown fee calculation strategy")
)

func MosaicIdFromString(mosaicId string) (*sdk.MosaicId, error) {
	if mosaicId == "" {
		return nil, ErrEmptyMosaic
	}

	mId, err := strconv.ParseUint(mosaicId, 16, 64)
	if err != nil {
		return nil, err
	}

	return sdk.NewMosaicId(mId)
}

func ParseFeeStrategy(feeStrategy *string) sdk.FeeCalculationStrategy {
	fee := sdk.FeeCalculationStrategy(0)
	switch *feeStrategy {
	case LowFeeStrategy:
		fee = sdk.LowCalculationStrategy
	case MiddleFeeStrategy:
		fee = sdk.MiddleCalculationStrategy
	case HighFeeStrategy:
		fee = sdk.HighCalculationStrategy
	default:
		fmt.Printf("%s: %s\n", ErrUnknownFeeStrategy, *feeStrategy)
		os.Exit(1)
	}

	return fee
}
