package sdk

import (
	"bytes"
	"encoding/hex"
	jsonLib "encoding/json"
	"sync"
)

type hashDto string

func (dto *hashDto) Hash() (*Hash, error) {
	s := string(*dto)

	if len(s) == 0 {
		return nil, nil
	}

	return StringToHash(s)
}

type hashDtos []hashDto

func (h *hashDtos) toStruct() ([]*Hash, error) {
	dtos := *h
	hashes := make([]*Hash, 0, len(dtos))

	for _, dto := range dtos {
		status, err := dto.Hash()
		if err != nil {
			return nil, err
		}

		hashes = append(hashes, status)
	}

	return hashes, nil
}

type signatureDto string

func (dto *signatureDto) Signature() (*Signature, error) {
	s := string(*dto)

	if len(s) == 0 {
		return nil, nil
	}

	return StringToSignature(s)
}

type transactionStatusDTOs []*transactionStatusDTO

func (t *transactionStatusDTOs) toStruct() ([]*TransactionStatus, error) {
	dtos := *t
	statuses := make([]*TransactionStatus, 0, len(dtos))

	for _, dto := range dtos {
		status, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

type abstractTransactionDTO struct {
	Type      EntityType              `json:"type"`
	Version   int64                   `json:"version"`
	MaxFee    *uint64DTO              `json:"maxFee"`
	Deadline  *blockchainTimestampDTO `json:"deadline"`
	Signature string                  `json:"signature"`
	Signer    string                  `json:"signer"`
}

func (dto *abstractTransactionDTO) toStruct(tInfo *TransactionInfo) (*AbstractTransaction, error) {
	nt := ExtractNetworkType(dto.Version)

	tv := EntityVersion(ExtractVersion(dto.Version))

	pa, err := NewAccountFromPublicKey(dto.Signer, nt)
	if err != nil {
		return nil, err
	}

	var d *Deadline
	if dto.Deadline != nil {
		d = NewDeadlineFromBlockchainTimestamp(dto.Deadline.toStruct())
	}

	var f Amount
	if dto.MaxFee != nil {
		f = dto.MaxFee.toStruct()
	}

	return &AbstractTransaction{
		*tInfo,
		nt,
		d,
		dto.Type,
		tv,
		f,
		dto.Signature,
		pa,
	}, nil
}

type transactionsPageDTO struct {
	Transactions []jsonLib.RawMessage `json:"data"`

	Pagination struct {
		TotalEntries uint64 `json:"totalEntries"`
		PageNumber   uint64 `json:"pageNumber"`
		PageSize     uint64 `json:"pageSize"`
		TotalPages   uint64 `json:"totalPages"`
	} `json:"pagination"`
}

func (t *transactionsPageDTO) toStruct(generationHash *Hash) (*TransactionsPage, error) {
	var wg sync.WaitGroup
	page := &TransactionsPage{
		Transactions: make([]Transaction, len(t.Transactions)),
		Pagination: Pagination{
			TotalEntries: t.Pagination.TotalEntries,
			PageNumber:   t.Pagination.PageNumber,
			PageSize:     t.Pagination.PageSize,
			TotalPages:   t.Pagination.TotalPages,
		},
	}

	errs := make([]error, len(t.Transactions))
	for i, t := range t.Transactions {
		wg.Add(1)
		go func(i int, t jsonLib.RawMessage) {
			defer wg.Done()
			page.Transactions[i], errs[i] = MapTransaction(bytes.NewBuffer([]byte(t)), generationHash)
		}(i, t)
	}

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return page, err
		}
	}

	return page, nil
}

type transactionInfoDTO struct {
	Height              uint64DTO `json:"height"`
	Index               uint32    `json:"index"`
	Id                  string    `json:"id"`
	TransactionHash     hashDto   `json:"hash"`
	MerkleComponentHash hashDto   `json:"merkleComponentHash"`
	AggregateHash       hashDto   `json:"aggregateHash,omitempty"`
	UniqueAggregateHash hashDto   `json:"uniqueAggregateHash,omitempty"`
	AggregateId         string    `json:"aggregateId,omitempty"`
}

func (dto *transactionInfoDTO) toStruct() (*TransactionInfo, error) {
	transactionHash, err := dto.TransactionHash.Hash()
	if err != nil {
		return nil, err
	}
	merkleComponentHash, err := dto.MerkleComponentHash.Hash()
	if err != nil {
		return nil, err
	}
	aggregateHash, err := dto.AggregateHash.Hash()
	if err != nil {
		return nil, err
	}
	uniqueAggregateHash, err := dto.UniqueAggregateHash.Hash()
	if err != nil {
		return nil, err
	}

	ref := TransactionInfo{
		dto.Height.toStruct(),
		dto.Index,
		dto.Id,
		transactionHash,
		merkleComponentHash,
		aggregateHash,
		uniqueAggregateHash,
		dto.AggregateId,
	}

	return &ref, nil
}

type accountPropertiesAddressModificationDTO struct {
	ModificationType PropertyModificationType `json:"type"`
	Address          string                   `json:"value"`
}

func (dto *accountPropertiesAddressModificationDTO) toStruct() (*AccountPropertiesAddressModification, error) {
	a, err := NewAddressFromBase32(dto.Address)
	if err != nil {
		return nil, err
	}

	return &AccountPropertiesAddressModification{
		dto.ModificationType,
		a,
	}, nil
}

type accountPropertiesAddressTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		PropertyType  PropertyType                               `json:"propertyType"`
		Modifications []*accountPropertiesAddressModificationDTO `json:"modifications"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *accountPropertiesAddressTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	ms := make([]*AccountPropertiesAddressModification, len(dto.Tx.Modifications))
	for i, m := range dto.Tx.Modifications {
		ms[i], err = m.toStruct()

		if err != nil {
			return nil, err
		}
	}

	return &AccountPropertiesAddressTransaction{
		*atx,
		dto.Tx.PropertyType,
		ms,
	}, nil
}

type accountPropertiesMosaicModificationDTO struct {
	ModificationType PropertyModificationType `json:"type"`
	AssetId          assetIdDTO               `json:"value"`
}

func (dto *accountPropertiesMosaicModificationDTO) toStruct() (*AccountPropertiesMosaicModification, error) {
	assetId, err := dto.AssetId.toStruct()
	if err != nil {
		return nil, err
	}

	return &AccountPropertiesMosaicModification{
		dto.ModificationType,
		assetId,
	}, nil
}

type accountPropertiesMosaicTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		PropertyType  PropertyType                              `json:"propertyType"`
		Modifications []*accountPropertiesMosaicModificationDTO `json:"modifications"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *accountPropertiesMosaicTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	ms := make([]*AccountPropertiesMosaicModification, len(dto.Tx.Modifications))
	for i, m := range dto.Tx.Modifications {
		ms[i], err = m.toStruct()

		if err != nil {
			return nil, err
		}
	}

	return &AccountPropertiesMosaicTransaction{
		*atx,
		dto.Tx.PropertyType,
		ms,
	}, nil
}

type accountPropertiesEntityTypeModificationDTO struct {
	ModificationType PropertyModificationType `json:"type"`
	EntityType       EntityType               `json:"value"`
}

func (dto *accountPropertiesEntityTypeModificationDTO) toStruct() (*AccountPropertiesEntityTypeModification, error) {
	return &AccountPropertiesEntityTypeModification{
		dto.ModificationType,
		dto.EntityType,
	}, nil
}

type accountPropertiesEntityTypeTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		PropertyType  PropertyType                                  `json:"propertyType"`
		Modifications []*accountPropertiesEntityTypeModificationDTO `json:"modifications"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *accountPropertiesEntityTypeTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	ms := make([]*AccountPropertiesEntityTypeModification, len(dto.Tx.Modifications))
	for i, m := range dto.Tx.Modifications {
		ms[i], err = m.toStruct()

		if err != nil {
			return nil, err
		}
	}

	return &AccountPropertiesEntityTypeTransaction{
		*atx,
		dto.Tx.PropertyType,
		ms,
	}, nil
}

type aliasTransactionDTO struct {
	abstractTransactionDTO
	NamespaceId namespaceIdDTO  `json:"namespaceId"`
	ActionType  AliasActionType `json:"aliasAction"`
}

func (dto *aliasTransactionDTO) toStruct(tInfo *TransactionInfo) (*AliasTransaction, error) {
	atx, err := dto.abstractTransactionDTO.toStruct(tInfo)
	if err != nil {
		return nil, err
	}

	namespaceId, err := dto.NamespaceId.toStruct()
	if err != nil {
		return nil, err
	}

	return &AliasTransaction{
		*atx,
		dto.ActionType,
		namespaceId,
	}, nil
}

type addressAliasTransactionDTO struct {
	Tx struct {
		aliasTransactionDTO
		Address string `json:"address"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *addressAliasTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.aliasTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	a, err := NewAddressFromBase32(dto.Tx.Address)
	if err != nil {
		return nil, err
	}

	return &AddressAliasTransaction{
		*atx,
		a,
	}, nil
}

type mosaicAliasTransactionDTO struct {
	Tx struct {
		aliasTransactionDTO
		MosaicId *mosaicIdDTO `json:"mosaicId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *mosaicAliasTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.aliasTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	mosaicId, err := dto.Tx.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	return &MosaicAliasTransaction{
		*atx,
		mosaicId,
	}, nil
}

type accountLinkTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		RemoteAccountKey string            `json:"remoteAccountKey"`
		Action           AccountLinkAction `json:"action"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *accountLinkTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	acc, err := NewAccountFromPublicKey(dto.Tx.RemoteAccountKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &AccountLinkTransaction{
		*atx,
		acc,
		dto.Tx.Action,
	}, nil
}

type networkConfigTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		ApplyHeightDelta        uint64DTO `json:"applyHeightDelta"`
		NetworkConfig           string    `json:"networkConfig"`
		SupportedEntityVersions string    `json:"supportedEntityVersions"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *networkConfigTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	applyHeightDelta := dto.Tx.ApplyHeightDelta.toUint64()

	s := NewSupportedEntities()

	err = s.UnmarshalBinary([]byte(dto.Tx.SupportedEntityVersions))
	if err != nil {
		return nil, err
	}

	c := NewNetworkConfig()

	err = c.UnmarshalBinary([]byte(dto.Tx.NetworkConfig))
	if err != nil {
		return nil, err
	}

	return &NetworkConfigTransaction{
		*atx,
		Duration(applyHeightDelta),
		c,
		s,
	}, nil
}

type blockchainUpgradeTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		UpgradePeriod        uint64DTO `json:"upgradePeriod"`
		NewBlockChainVersion uint64DTO `json:"newBlockChainVersion"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *blockchainUpgradeTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	upgradePeriod := dto.Tx.UpgradePeriod.toUint64()
	newBlockChainVersion := dto.Tx.NewBlockChainVersion.toUint64()

	return &BlockchainUpgradeTransaction{
		*atx,
		Duration(upgradePeriod),
		BlockChainVersion(newBlockChainVersion),
	}, nil
}

type aggregateTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Cosignatures      []*aggregateTransactionCosignatureDTO `json:"cosignatures"`
		InnerTransactions []map[string]interface{}              `json:"transactions"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *aggregateTransactionDTO) toStruct(generationHash *Hash) (Transaction, error) {
	txsr, err := json.Marshal(dto.Tx.InnerTransactions)
	if err != nil {
		return nil, err
	}

	txs, err := MapTransactions(bytes.NewBuffer(txsr), generationHash)
	if err != nil {
		return nil, err
	}

	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	as := make([]*AggregateTransactionCosignature, len(dto.Tx.Cosignatures))
	for i, a := range dto.Tx.Cosignatures {
		as[i], err = a.toStruct(atx.NetworkType)
	}
	if err != nil {
		return nil, err
	}

	for _, tx := range txs {
		iatx := tx.GetAbstractTransaction()
		iatx.Deadline = atx.Deadline
		iatx.Signature = atx.Signature
		iatx.MaxFee = atx.MaxFee
		iatx.TransactionInfo = atx.TransactionInfo
	}

	agtx := AggregateTransaction{
		*atx,
		txs,
		as,
	}

	return &agtx, agtx.UpdateUniqueAggregateHash(generationHash)
}

type accountMetadataTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		TargetKey         string    `json:"targetKey"`
		ScopedMetadataKey uint64DTO `json:"scopedMetadataKey"`
		ValueSizeDelta    int16     `json:"valueSizeDelta"`
		Value             string    `json:"value"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *accountMetadataTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}
	acc, err := NewAccountFromPublicKey(dto.Tx.TargetKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	value, err := hex.DecodeString(dto.Tx.Value)
	if err != nil {
		return nil, err
	}

	tx := BasicMetadataTransaction{
		*atx,
		acc,
		dto.Tx.ScopedMetadataKey.toStruct(),
		value,
		dto.Tx.ValueSizeDelta,
	}

	return &AccountMetadataTransaction{
		tx,
	}, nil
}

type mosaicMetadataTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		TargetKey         string       `json:"targetKey"`
		ScopedMetadataKey uint64DTO    `json:"scopedMetadataKey"`
		MosaicId          *mosaicIdDTO `json:"targetMosaicId"`
		ValueSizeDelta    int16        `json:"valueSizeDelta"`
		Value             string       `json:"value"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *mosaicMetadataTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}
	acc, err := NewAccountFromPublicKey(dto.Tx.TargetKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	value, err := hex.DecodeString(dto.Tx.Value)
	if err != nil {
		return nil, err
	}

	tx := BasicMetadataTransaction{
		*atx,
		acc,
		dto.Tx.ScopedMetadataKey.toStruct(),
		value,
		dto.Tx.ValueSizeDelta,
	}
	mosaicId, err := dto.Tx.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	return &MosaicMetadataTransaction{
		tx,
		mosaicId,
	}, nil
}

type namespaceMetadataTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		TargetKey         string          `json:"targetKey"`
		ScopedMetadataKey uint64DTO       `json:"scopedMetadataKey"`
		NamespaceId       *namespaceIdDTO `json:"targetNamespaceId"`
		ValueSizeDelta    int16           `json:"valueSizeDelta"`
		Value             string          `json:"value"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *namespaceMetadataTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}
	acc, err := NewAccountFromPublicKey(dto.Tx.TargetKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	value, err := hex.DecodeString(dto.Tx.Value)
	if err != nil {
		return nil, err
	}

	tx := BasicMetadataTransaction{
		*atx,
		acc,
		dto.Tx.ScopedMetadataKey.toStruct(),
		value,
		dto.Tx.ValueSizeDelta,
	}
	namespaceId, err := dto.Tx.NamespaceId.toStruct()
	if err != nil {
		return nil, err
	}

	return &NamespaceMetadataTransaction{
		tx,
		namespaceId,
	}, nil
}

type modifyMetadataTransactionDTO struct {
	abstractTransactionDTO
	MetadataType  MetadataType               `json:"metadataType"`
	Modifications []*metadataModificationDTO `json:"modifications"`
}

func (dto *modifyMetadataTransactionDTO) toStruct(tInfo *TransactionInfo) (*ModifyMetadataTransaction, error) {
	atx, err := dto.abstractTransactionDTO.toStruct(tInfo)
	if err != nil {
		return nil, err
	}

	ms, err := metadataDTOArrayToStruct(dto.Modifications, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &ModifyMetadataTransaction{
		*atx,
		dto.MetadataType,
		ms,
	}, nil
}

type modifyMetadataAddressTransactionDTO struct {
	Tx struct {
		modifyMetadataTransactionDTO
		Address string `json:"metadataId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyMetadataAddressTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.modifyMetadataTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	a, err := NewAddressFromBase32(dto.Tx.Address)
	if err != nil {
		return nil, err
	}

	return &ModifyMetadataAddressTransaction{
		*atx,
		a,
	}, nil
}

type modifyMetadataMosaicTransactionDTO struct {
	Tx struct {
		modifyMetadataTransactionDTO
		MosaicId *mosaicIdDTO `json:"metadataId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyMetadataMosaicTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.modifyMetadataTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	mosaicId, err := dto.Tx.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	return &ModifyMetadataMosaicTransaction{
		*atx,
		mosaicId,
	}, nil
}

type modifyMetadataNamespaceTransactionDTO struct {
	Tx struct {
		modifyMetadataTransactionDTO
		NamespaceId *namespaceIdDTO `json:"metadataId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyMetadataNamespaceTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.modifyMetadataTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	namespaceId, err := dto.Tx.NamespaceId.toStruct()
	if err != nil {
		return nil, err
	}

	return &ModifyMetadataNamespaceTransaction{
		*atx,
		namespaceId,
	}, nil
}

type mosaicDefinitionTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Properties  mosaicPropertiesDTO `json:"properties"`
		MosaicNonce int64               `json:"mosaicNonce"`
		MosaicId    *mosaicIdDTO        `json:"mosaicId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *mosaicDefinitionTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	mosaicId, err := dto.Tx.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	properties, err := dto.Tx.Properties.toStruct()
	if err != nil {
		return nil, err
	}

	return &MosaicDefinitionTransaction{
		*atx,
		properties,
		uint32(dto.Tx.MosaicNonce),
		mosaicId,
	}, nil
}

type mosaicSupplyChangeTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		MosaicSupplyType `json:"direction"`
		AssetId          *assetIdDTO `json:"mosaicId"`
		Delta            uint64DTO   `json:"delta"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *mosaicSupplyChangeTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	assetId, err := dto.Tx.AssetId.toStruct()
	if err != nil {
		return nil, err
	}

	return &MosaicSupplyChangeTransaction{
		*atx,
		dto.Tx.MosaicSupplyType,
		assetId,
		dto.Tx.Delta.toStruct(),
	}, nil
}

type transferTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Message messageDTO   `json:"message"`
		Mosaics []*mosaicDTO `json:"mosaics"`
		Address string       `json:"recipient"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *transferTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	mosaics := make([]*Mosaic, len(dto.Tx.Mosaics))

	for i, mosaic := range dto.Tx.Mosaics {
		msc, err := mosaic.toStruct()
		if err != nil {
			return nil, err
		}

		mosaics[i] = msc
	}

	a, err := NewAddressFromBase32(dto.Tx.Address)
	if err != nil {
		return nil, err
	}

	m, err := dto.Tx.Message.toStruct()
	if err != nil {
		return nil, err
	}

	return &TransferTransaction{
		*atx,
		m,
		mosaics,
		a,
	}, nil
}

type harvesterTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		HarvesterKey string `json:"harvesterKey"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *harvesterTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	acc, err := NewAccountFromPublicKey(dto.Tx.HarvesterKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &HarvesterTransaction{
		*atx,
		acc,
	}, nil
}

type modifyMultisigAccountTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		MinApprovalDelta int8                                  `json:"minApprovalDelta"`
		MinRemovalDelta  int8                                  `json:"minRemovalDelta"`
		Modifications    []*multisigCosignatoryModificationDTO `json:"modifications"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyMultisigAccountTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	ms, err := multisigCosignatoryDTOArrayToStruct(dto.Tx.Modifications, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &ModifyMultisigAccountTransaction{
		*atx,
		dto.Tx.MinApprovalDelta,
		dto.Tx.MinRemovalDelta,
		ms,
	}, nil
}

type modifyContractTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DurationDelta uint64DTO                             `json:"durationDelta"`
		Hash          hashDto                               `json:"hash"`
		Customers     []*multisigCosignatoryModificationDTO `json:"customers"`
		Executors     []*multisigCosignatoryModificationDTO `json:"executors"`
		Verifiers     []*multisigCosignatoryModificationDTO `json:"verifiers"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyContractTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	customers, err := multisigCosignatoryDTOArrayToStruct(dto.Tx.Customers, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	executors, err := multisigCosignatoryDTOArrayToStruct(dto.Tx.Executors, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	verifiers, err := multisigCosignatoryDTOArrayToStruct(dto.Tx.Verifiers, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	hash, err := dto.Tx.Hash.Hash()
	if err != nil {
		return nil, err
	}

	return &ModifyContractTransaction{
		*atx,
		dto.Tx.DurationDelta.toStruct(),
		hash,
		customers,
		executors,
		verifiers,
	}, nil
}

type registerNamespaceTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Id            namespaceIdDTO `json:"namespaceId"`
		NamespaceType `json:"namespaceType"`
		NamspaceName  string    `json:"name"`
		Duration      uint64DTO `json:"duration"`
		ParentId      namespaceIdDTO
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *registerNamespaceTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	d := Duration(0)
	var n *NamespaceId = nil

	if dto.Tx.NamespaceType == Root {
		d = dto.Tx.Duration.toStruct()
	} else {
		n, err = dto.Tx.ParentId.toStruct()
		if err != nil {
			return nil, err
		}
	}

	nsId, err := dto.Tx.Id.toStruct()
	if err != nil {
		return nil, err
	}

	return &RegisterNamespaceTransaction{
		*atx,
		nsId,
		dto.Tx.NamespaceType,
		dto.Tx.NamspaceName,
		d,
		n,
	}, nil
}

type lockFundsTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		AssetId  assetIdDTO `json:"mosaicId"`
		Amount   uint64DTO  `json:"amount"`
		Duration uint64DTO  `json:"duration"`
		Hash     hashDto    `json:"hash"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *lockFundsTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	assetId, err := dto.Tx.AssetId.toStruct()
	if err != nil {
		return nil, err
	}

	mosaic, err := NewMosaic(assetId, dto.Tx.Amount.toStruct())
	if err != nil {
		return nil, err
	}

	hash, err := dto.Tx.Hash.Hash()
	if err != nil {
		return nil, err
	}

	return &LockFundsTransaction{
		*atx,
		mosaic,
		dto.Tx.Duration.toStruct(),
		&SignedTransaction{AggregateBonded, "", hash},
	}, nil
}

type secretLockTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		AssetId   *assetIdDTO `json:"mosaicId"`
		Amount    *uint64DTO  `json:"amount"`
		HashType  HashType    `json:"hashAlgorithm"`
		Duration  uint64DTO   `json:"duration"`
		Secret    string      `json:"secret"`
		Recipient string      `json:"recipient"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *secretLockTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	a, err := NewAddressFromBase32(dto.Tx.Recipient)
	if err != nil {
		return nil, err
	}

	assetId, err := dto.Tx.AssetId.toStruct()
	if err != nil {
		return nil, err
	}

	mosaic, err := NewMosaic(assetId, dto.Tx.Amount.toStruct())
	if err != nil {
		return nil, err
	}

	secret, err := NewSecretFromHexString(dto.Tx.Secret, dto.Tx.HashType)
	if err != nil {
		return nil, err
	}

	return &SecretLockTransaction{
		*atx,
		mosaic,
		dto.Tx.Duration.toStruct(),
		secret,
		a,
	}, nil
}

type secretProofTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		HashType  `json:"hashAlgorithm"`
		Proof     string `json:"proof"`
		Recipient string `json:"recipient"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *secretProofTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	proof, err := NewProofFromHexString(dto.Tx.Proof)
	if err != nil {
		return nil, err
	}

	a, err := NewAddressFromBase32(dto.Tx.Recipient)
	if err != nil {
		return nil, err
	}

	return &SecretProofTransaction{
		*atx,
		dto.Tx.HashType,
		proof,
		a,
	}, nil
}

type aggregateTransactionCosignatureDTO struct {
	Signature string `json:"signature"`
	Signer    string
}

func (dto *aggregateTransactionCosignatureDTO) toStruct(networkType NetworkType) (*AggregateTransactionCosignature, error) {
	acc, err := NewAccountFromPublicKey(dto.Signer, networkType)
	if err != nil {
		return nil, err
	}
	return &AggregateTransactionCosignature{
		dto.Signature,
		acc,
	}, nil
}

type multisigCosignatoryModificationDTO struct {
	Type          MultisigCosignatoryModificationType `json:"type"`
	PublicAccount string                              `json:"cosignatoryPublicKey"`
}

func (dto *multisigCosignatoryModificationDTO) toStruct(networkType NetworkType) (*MultisigCosignatoryModification, error) {
	acc, err := NewAccountFromPublicKey(dto.PublicAccount, networkType)
	if err != nil {
		return nil, err
	}

	return &MultisigCosignatoryModification{
		dto.Type,
		acc,
	}, nil
}

type metadataModificationDTO struct {
	Type  MetadataModificationType `json:"modificationType"`
	Key   string                   `json:"key"`
	Value string                   `json:"value"`
}

func (dto *metadataModificationDTO) toStruct(networkType NetworkType) (*MetadataModification, error) {
	return &MetadataModification{
		dto.Type,
		dto.Key,
		dto.Value,
	}, nil
}

type transactionStatusDTO struct {
	Group    TransactionGroup       `json:"group"`
	Status   string                 `json:"status"`
	Hash     hashDto                `json:"hash"`
	Deadline blockchainTimestampDTO `json:"deadline"`
	Height   uint64DTO              `json:"height"`
}

func (dto *transactionStatusDTO) toStruct() (*TransactionStatus, error) {
	hash, err := dto.Hash.Hash()
	if err != nil {
		return nil, err

	}
	return &TransactionStatus{
		NewDeadlineFromBlockchainTimestamp(dto.Deadline.toStruct()),
		dto.Group,
		dto.Status,
		hash,
		dto.Height.toStruct(),
	}, nil
}

type mosaicModifyLevyTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		MosaicId   *mosaicIdDTO   `json:"mosaicId"`
		MosaicLevy *mosaicLevyDTO `json:"levy"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *mosaicModifyLevyTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	mosaicId, err := dto.Tx.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	levy, err := dto.Tx.MosaicLevy.toStruct()
	if err != nil {
		return nil, err
	}

	return &MosaicModifyLevyTransaction{
		*atx,
		mosaicId,
		levy,
	}, nil
}

type mosaicRemoveLevyTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		MosaicId *mosaicIdDTO `json:"mosaicId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *mosaicRemoveLevyTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	mosaicId, err := dto.Tx.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	return &MosaicRemoveLevyTransaction{
		*atx,
		mosaicId,
	}, nil
}

type TransactionIdsDTO struct {
	Ids []string `json:"transactionIds"`
}

type TransactionHashesDTO struct {
	Hashes []string `json:"hashes"`
}

type addDbrbProcessTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *addDbrbProcessTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &AddDbrbProcessTransaction{
		*atx,
	}, nil
}

type removeDbrbProcessTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *removeDbrbProcessTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &RemoveDbrbProcessTransaction{
		*atx,
	}, nil
}

type removeDbrbProcessByNetworkTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *removeDbrbProcessByNetworkTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &RemoveDbrbProcessByNetworkTransaction{
		*atx,
	}, nil
}
