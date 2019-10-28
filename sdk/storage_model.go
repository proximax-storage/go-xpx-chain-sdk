package sdk

import "fmt"

type DriveState uint8

const (
	NotStarted DriveState = iota
	Pending
	InProgress
	Finished
)

type PaymentInformation struct {
	Receiver *PublicAccount
	Amount   Amount
	Height   Height
}

type BillingDescription struct {
	Start    Height
	End      Height
	Payments []*PaymentInformation
}

type ReplicatorInfo struct {
	Start               Height
	End                 Height
	Deposit             Amount
	FilesWithoutDeposit map[Hash]uint16
}

type FileInfo struct {
	FileSize StorageSize
	Deposit  Amount
	Payments []*PaymentInformation
}

type Drive struct {
	State            DriveState
	Owner            *PublicAccount
	RootHash         *Hash
	Duration         Duration
	BillingPeriod    Duration
	BillingPrice     Amount
	DriveSize        StorageSize
	Replicas         uint16
	MinReplicators   uint16
	PercentApprovers uint8
	BillingHistory   []BillingDescription
	Files            map[Hash]*FileInfo
	Replicators      map[PublicAccount]*ReplicatorInfo
}

// Prepare Drive Transaction
type PrepareDriveTransaction struct {
	AbstractTransaction
	Owner            *PublicAccount
	Duration         Duration
	BillingPeriod    Duration
	BillingPrice     Amount
	DriveSize        StorageSize
	Replicas         uint16
	MinReplicators   uint16
	PercentApprovers uint8
}

// Join Drive Transaction

type JoinToDriveTransaction struct {
	AbstractTransaction
	DriveKey *PublicAccount
}

// Drive File System Transaction
type File struct {
	FileHash *Hash
}

func (file *File) String() string {
	return fmt.Sprintf(
		`
			"FileHash": %s,
		`,
		file.FileHash,
	)
}

type AddAction struct {
	File
	FileSize StorageSize
}

func (action *AddAction) String() string {
	return fmt.Sprintf(
		`
			"FileHash": %s,
			"FileSize": %s,
		`,
		action.FileHash,
		action.FileSize,
	)
}

type RemoveAction struct {
	File
}

type DriveFileSystemTransaction struct {
	AbstractTransaction
	DriveKey      *PublicAccount
	NewRootHash   *Hash
	OldRootHash   *Hash
	AddActions    []*AddAction
	RemoveActions []*RemoveAction
}

// Files Deposit Transaction
type FilesDepositTransaction struct {
	AbstractTransaction
	DriveKey *PublicAccount
	Files    []*File
}

// End Drive Transaction

type EndDriveTransaction struct {
	AbstractTransaction
	DriveKey *PublicAccount
}
