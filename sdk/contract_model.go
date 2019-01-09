package sdk

import "math/big"

type ContractInfo struct {
	Multisig        string
	MultisigAddress *Address
	Start           *big.Int
	Duration        *big.Int
	Content         string
	Customers       []string
	Executors       []string
	Verifiers       []string
}
