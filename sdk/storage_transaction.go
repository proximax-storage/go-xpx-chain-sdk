package sdk

func NewModifyDriveTransaction(
	deadline *Deadline, priceDelta Amount, durationDelta Duration,
	sizeDelta SizeDelta, minReplicatorsDelta int8, minApproversDelta int8, replicasDelta int8,
	networkType NetworkType) (*ModifyDriveTransaction, error) {

	mctx := ModifyDriveTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     ModifyDriveVersion,
			Deadline:    deadline,
			Type:        ModifyDrive,
			NetworkType: networkType,
		},
		PriceDelta:          priceDelta,
		DurationDelta:       durationDelta,
		SizeDelta:           sizeDelta,
		ReplicasDelta:       replicasDelta,
		MinReplicatorsDelta: minReplicatorsDelta,
		MinApproversDelta:   minApproversDelta,
	}

	return &mctx, nil
}

func NewJoinToDriveTransaction(deadline *Deadline, driveKey *PublicAccount, networkType NetworkType) (*JoinToDriveTransaction, error) {

	tx := JoinToDriveTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     JoinToDriveVersion,
			Deadline:    deadline,
			Type:        JoinToDrive,
			NetworkType: networkType,
		},
		DriveKey: driveKey,
	}

	return &tx, nil
}

func NewDriveFileSystemTransaction(
	deadline *Deadline, rootHash *Hash, xorRootHash *Hash, addActions []*AddAction, removeActions []*RemoveAction, networkType NetworkType) (*DriveFileSystemTransaction, error) {

	tx := DriveFileSystemTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     DriveFileSystemVersion,
			Deadline:    deadline,
			Type:        DriveFileSystem,
			NetworkType: networkType,
		},
		RootHash:      rootHash,
		XorRootHash:   xorRootHash,
		AddActions:    addActions,
		RemoveActions: removeActions,
	}

	return &tx, nil
}

func NewFilesDepositTransaction(
	deadline *Deadline, driveKey *PublicAccount, files []*File, networkType NetworkType) (*FilesDepositTransaction, error) {

	tx := FilesDepositTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     FilesDepositVersion,
			Deadline:    deadline,
			Type:        FilesDeposit,
			NetworkType: networkType,
		},
		DriveKey: driveKey,
		Files:    files,
	}

	return &tx, nil
}

func NewEndDriveTransaction(
	deadline *Deadline, networkType NetworkType) (*EndDriveTransaction, error) {

	tx := EndDriveTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     EndDriveVersion,
			Deadline:    deadline,
			Type:        EndDrive,
			NetworkType: networkType,
		},
	}

	return &tx, nil
}
