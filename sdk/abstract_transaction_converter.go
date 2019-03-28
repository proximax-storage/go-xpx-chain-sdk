package sdk

import (
	"math/big"
	"time"
)

func newAbstractTransactionConverterImpl(accountFactory AccountFactory) abstractTransactionConverter {
	return &abstractTransactionConverterImpl{accountFactory: accountFactory}
}

type abstractTransactionConverter interface {
	Convert(abstractTransactionDTO, *TransactionInfo) (*AbstractTransaction, error)
}

type abstractTransactionConverterImpl struct {
	accountFactory AccountFactory
}

func (c *abstractTransactionConverterImpl) Convert(dto abstractTransactionDTO, tInfo *TransactionInfo) (*AbstractTransaction, error) {
	t, err := TransactionTypeFromRaw(dto.Type)
	if err != nil {
		return nil, err
	}

	nt := ExtractNetworkType(dto.Version)

	tv := TransactionVersion(ExtractVersion(dto.Version))

	pa, err := c.accountFactory.NewAccountFromPublicKey(dto.Signer, nt)
	if err != nil {
		return nil, err
	}

	var d *Deadline
	if dto.Deadline != nil {
		d = &Deadline{time.Unix(0, dto.Deadline.toBigInt().Int64()*int64(time.Millisecond))}
	}

	var f *big.Int
	if dto.Fee != nil {
		f = dto.Fee.toBigInt()
	}

	return &AbstractTransaction{
		tInfo,
		nt,
		d,
		t,
		tv,
		f,
		dto.Signature,
		pa,
	}, nil
}
