// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type filesWithoutDepositDTO struct {
	FileHash	hashDto    	`json:"fileHash"`
	Count		uint16 		`json:"count"`
}

type filesWithoutDepositDTOs []*filesWithoutDepositDTO

func (ref *filesWithoutDepositDTOs) toStruct() (map[Hash]uint16, error) {
	var (
		dtos  = *ref
		deposits = make(map[Hash]uint16)
	)

	for _, dto := range dtos {
		fileHash, err := dto.FileHash.Hash()
		if err != nil {
			return nil, err
		}

		deposits[*fileHash] = dto.Count
	}

	return deposits, nil
}

type paymentDTO struct {
	Receiver	string        	`json:"receiver"`
	Amount		uint64DTO 		`json:"amount"`
	Height		uint64DTO 		`json:"height"`
}

func (ref *paymentDTO) toStruct(networkType NetworkType) (*PaymentInformation, error) {
	receiver, err := NewAccountFromPublicKey(ref.Receiver, networkType)
	if err != nil {
		return nil, err
	}

	return &PaymentInformation{
		Receiver: 		receiver,
		Amount: 		ref.Amount.toStruct(),
		Height: 		ref.Height.toStruct(),
	}, nil
}

type paymentsDTOs []*paymentDTO

func (ref *paymentsDTOs) toStruct(networkType NetworkType) ([]*PaymentInformation, error) {
	var (
		dtos  = *ref
		payments = make([]*PaymentInformation, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		payments = append(payments, info)
	}

	return payments, nil
}

type actionDTO struct {
	Type			FileActionType        		`json:"type"`
	Height			uint64DTO        			`json:"height"`
}

type actionsDTOs []*actionDTO

func (ref *actionsDTOs) toStruct() ([]*FileAction, error) {
	var (
		dtos  = *ref
		actions = make([]*FileAction, 0, len(dtos))
	)

	for _, dto := range dtos {
		info := FileAction{
			Type:	 	dto.Type,
			Height: 	dto.Height.toStruct(),
		}

		actions = append(actions, &info)
	}

	return actions, nil
}

type replicatorDTO struct {
	Replicator			string        			`json:"replicator"`
	Start				uint64DTO 				`json:"start"`
	End					uint64DTO 				`json:"end"`
	Deposit				uint64DTO				`json:"deposit"`
	FilesWithoutDeposit	filesWithoutDepositDTOs	`json:"filesWithoutDeposit"`
}

type replicatorsDTOs []*replicatorDTO

func (ref *replicatorsDTOs) toStruct(networkType NetworkType) (map[PublicAccount]*ReplicatorInfo, error) {
	var (
		dtos  = *ref
		replicators = make(map[PublicAccount]*ReplicatorInfo)
	)

	for _, dto := range dtos {
		replicator, err := NewAccountFromPublicKey(dto.Replicator, networkType)
		if err != nil {
			return nil, err
		}

		filesWithoutDeposit, err := dto.FilesWithoutDeposit.toStruct()
		if err != nil {
			return nil, err
		}

		info := ReplicatorInfo{
			Start:					dto.Start.toStruct(),
			End:					dto.End.toStruct(),
			Deposit:				dto.Deposit.toStruct(),
			FilesWithoutDeposit:	filesWithoutDeposit,
		}

		replicators[*replicator] = &info
	}

	return replicators, nil
}

type fileDTO struct {
	FileHash	hashDto        	`json:"fileHash"`
	Deposit		uint64DTO 		`json:"deposit"`
	FileSize	uint64DTO 		`json:"size"`
	Payments	paymentsDTOs	`json:"payments"`
	Actions		actionsDTOs		`json:"actions"`
}

type filesDTOs []*fileDTO

func (ref *filesDTOs) toStruct(networkType NetworkType) (map[Hash]*FileInfo, error) {
	var (
		dtos  = *ref
		files = make(map[Hash]*FileInfo)
	)

	for _, dto := range dtos {
		fileHash, err := dto.FileHash.Hash()
		if err != nil {
			return nil, err
		}

		payments, err := dto.Payments.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		actions, err := dto.Actions.toStruct()
		if err != nil {
			return nil, err
		}

		info := FileInfo{
			Deposit:	dto.Deposit.toStruct(),
			FileSize:	dto.FileSize.toStruct(),
			Payments:	payments,
			Actions:	actions,
		}

		files[*fileHash] = &info
	}

	return files, nil
}


type billingDescriptionDTO struct {
	Start 		uint64DTO		`json:"start"`
	End       	uint64DTO 		`json:"end"`
	Payments  	paymentsDTOs 	`json:"payments"`
}

func (ref *billingDescriptionDTO) toStruct(networkType NetworkType) (*BillingDescription, error) {
	payments, err := ref.Payments.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	return &BillingDescription{
		Start: 		ref.Start.toStruct(),
		End: 		ref.End.toStruct(),
		Payments: 	payments,
	}, nil
}

type billingHistoryDTOs []*billingDescriptionDTO

func (ref *billingHistoryDTOs) toStruct(networkType NetworkType) ([]*BillingDescription, error) {
	var (
		dtos  = *ref
		history = make([]*BillingDescription, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		history = append(history, info)
	}

	return history, nil
}

type driveDTO struct {
	Drive struct {
		DriveKey 			string   	     	`json:"multisig"`
		State       		DriveState 			`json:"state"`
		Owner       		string 				`json:"owner"`
		RootHash      		hashDto 			`json:"rootHash"`
		Duration       		uint64DTO 			`json:"duration"`
		BillingPeriod		uint64DTO 			`json:"billingPeriod"`
		BillingPrice		uint64DTO 			`json:"billingPrice"`
		DriveSize       	uint64DTO 			`json:"size"`
		Replicas       		uint16 				`json:"replicas"`
		MinReplicators		uint16 				`json:"minReplicators"`
		PercentApprovers	uint8 				`json:"percentApprovers"`
		BillingHistory		billingHistoryDTOs	`json:"billingHistory"`
		Files       		filesDTOs			`json:"files"`
		Replicators       	replicatorsDTOs		`json:"replicators"`
	}
}

func (ref *driveDTO) toStruct(networkType NetworkType) (*Drive, error) {
	drive := Drive{}

	driveKey, err := NewAccountFromPublicKey(ref.Drive.DriveKey, networkType)
	if err != nil {
		return nil, err
	}

	owner, err := NewAccountFromPublicKey(ref.Drive.Owner, networkType)
	if err != nil {
		return nil, err
	}

	rootHash, err := ref.Drive.RootHash.Hash()
	if err != nil {
		return nil, err
	}

	drive.DriveKey = driveKey
	drive.State = ref.Drive.State
	drive.Owner = owner
	drive.RootHash = rootHash
	drive.Duration = ref.Drive.Duration.toStruct()
	drive.BillingPeriod = ref.Drive.BillingPeriod.toStruct()
	drive.BillingPrice = ref.Drive.BillingPrice.toStruct()
	drive.DriveSize = ref.Drive.DriveSize.toStruct()
	drive.Replicas = ref.Drive.Replicas
	drive.MinReplicators = ref.Drive.MinReplicators
	drive.PercentApprovers = ref.Drive.PercentApprovers

	billingHistory, err := ref.Drive.BillingHistory.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	drive.BillingHistory = billingHistory

	files, err := ref.Drive.Files.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	drive.Files = files

	replicators, err := ref.Drive.Replicators.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	drive.Replicators = replicators

	return &drive, nil
}

type driveDTOs []*driveDTO

func (ref *driveDTOs) toStruct(networkType NetworkType) ([]*Drive, error) {
	var (
		dtos  = *ref
		drives = make([]*Drive, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		drives = append(drives, info)
	}

	return drives, nil
}