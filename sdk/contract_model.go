package sdk

type ContractInfo struct {
	Multisig        string
	MultisigAddress *Address
	Start           Height
	Duration        Duration
	Content         string
	Customers       []string
	Executors       []string
	Verifiers       []string
}
