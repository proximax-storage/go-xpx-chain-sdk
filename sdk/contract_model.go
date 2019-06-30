package sdk

import "github.com/proximax-storage/go-xpx-utils/str"

type ContractInfo struct {
	Multisig        string
	MultisigAddress *Address
	Start           Height
	Duration        Duration
	Content         Hash
	Customers       []string
	Executors       []string
	Verifiers       []string
}

func (ref *ContractInfo) String() string {
	return str.StructToString(
		"ContractInfo",
		str.NewField("Multisig", str.StringPattern, ref.Multisig),
		str.NewField("MultisigAddress", str.StringPattern, ref.MultisigAddress),
		str.NewField("Start", str.IntPattern, ref.Start),
		str.NewField("Duration", str.IntPattern, ref.Duration),
		str.NewField("Content", str.StringPattern, ref.Content),
		str.NewField("Customers", str.StringPattern, ref.Customers),
		str.NewField("Executors", str.StringPattern, ref.Executors),
		str.NewField("Verifiers", str.StringPattern, ref.Verifiers),
	)
}
