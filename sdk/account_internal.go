package sdk

import (
	"encoding/base32"
	"encoding/hex"
	"errors"

	"github.com/proximax-storage/go-xpx-crypto"
)

var addressNet = map[uint8]NetworkType{
	96:  Mijin,
	144: MijinTest,
	145: AliasAddress,
	184: Public,
	168: PublicTest,
	200: Private,
	176: PrivateTest,
}

type propertiesDTO struct {
	PropertyType PropertyType `json:"propertyType"`
	MosaicIds    mosaicIdDTOs
	Addresses    []string
	EntityTypes  []EntityType
}

func (d *propertiesDTO) UnmarshalJSON(data []byte) error {
	temp := struct {
		PropertyType PropertyType `json:"propertyType"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	d.PropertyType = temp.PropertyType

	if temp.PropertyType&AllowAddress != 0 {
		addresses := struct {
			Addresses []string `json:"values"`
		}{}
		if err := json.Unmarshal(data, &addresses); err != nil {
			return err
		}
		d.Addresses = addresses.Addresses
	} else if temp.PropertyType&AllowMosaic != 0 {
		mosaicIds := struct {
			MosaicIds mosaicIdDTOs `json:"values"`
		}{}
		if err := json.Unmarshal(data, &mosaicIds); err != nil {
			return err
		}
		d.MosaicIds = mosaicIds.MosaicIds
	} else if temp.PropertyType&AllowTransaction != 0 {
		entityTypes := struct {
			EntityTypes []EntityType `json:"values"`
		}{}
		if err := json.Unmarshal(data, &entityTypes); err != nil {
			return err
		}
		d.EntityTypes = entityTypes.EntityTypes
	} else {
		return errors.New("not supported PropertyType")
	}

	return nil
}

type accountPropertiesDTO struct {
	AccountProperties struct {
		Address    string           `json:"address"`
		Properties []*propertiesDTO `json:"properties"`
	} `json:"accountProperties"`
}

func (ref *accountPropertiesDTO) toStruct() (*AccountProperties, error) {
	var err error = nil
	properties := AccountProperties{
		AllowedAddresses:   make([]*Address, 0),
		AllowedMosaicId:    make([]*MosaicId, 0),
		AllowedEntityTypes: make([]EntityType, 0),
		BlockedAddresses:   make([]*Address, 0),
		BlockedMosaicId:    make([]*MosaicId, 0),
		BlockedEntityTypes: make([]EntityType, 0),
	}

	properties.Address, err = NewAddressFromBase32(ref.AccountProperties.Address)
	if err != nil {
		return nil, err
	}

	for _, p := range ref.AccountProperties.Properties {
		switch p.PropertyType {
		case AllowAddress:
			properties.AllowedAddresses, err = EncodedStringToAddresses(p.Addresses...)
		case AllowMosaic:
			properties.AllowedMosaicId, err = p.MosaicIds.toStruct()
		case AllowTransaction:
			properties.AllowedEntityTypes = p.EntityTypes
		case BlockAddress:
			properties.BlockedAddresses, err = EncodedStringToAddresses(p.Addresses...)
		case BlockMosaic:
			properties.BlockedMosaicId, err = p.MosaicIds.toStruct()
		case BlockTransaction:
			properties.BlockedEntityTypes = p.EntityTypes
		}

		if err != nil {
			return nil, err
		}
	}

	return &properties, nil
}

type accountPropertiesDTOs []*accountPropertiesDTO

func (a accountPropertiesDTOs) toStruct() ([]*AccountProperties, error) {
	var (
		accountProperties = make([]*AccountProperties, len(a))
		err               error
	)

	for idx, dto := range a {
		accountProperties[idx], err = dto.toStruct()
		if err != nil {
			return nil, err
		}
	}

	return accountProperties, nil
}

type reputationDTO struct {
	PositiveInteractions uint64DTO `json:"positiveInteractions"`
	NegativeInteractions uint64DTO `json:"negativeInteractions"`
}

func (ref *reputationDTO) toFloat(repConfig *reputationConfig) float64 {
	posInter := ref.PositiveInteractions.toUint64()
	negInter := ref.NegativeInteractions.toUint64()

	if posInter < repConfig.minInteractions {
		return repConfig.defaultReputation
	}

	rep := (posInter - negInter) / posInter

	return float64(rep)
}

type accountInfoDTO struct {
	Account struct {
		Address                string         `json:"address"`
		AddressHeight          uint64DTO      `json:"addressHeight"`
		PublicKey              string         `json:"publicKey"`
		PublicKeyHeight        uint64DTO      `json:"publicKeyHeight"`
		AccountType            AccountType    `json:"accountType"`
		LinkedAccountPublicKey string         `json:"linkedAccountKey"`
		Mosaics                []*mosaicDTO   `json:"mosaics"`
		Reputation             *reputationDTO `json:"reputation"`
	} `json:"account"`
}

func (dto *accountInfoDTO) toStruct(repConfig *reputationConfig) (*AccountInfo, error) {
	var (
		ms  = make([]*Mosaic, len(dto.Account.Mosaics))
		err error
	)

	for idx, m := range dto.Account.Mosaics {
		ms[idx], err = m.toStruct()
		if err != nil {
			return nil, err
		}
	}

	add, err := NewAddressFromBase32(dto.Account.Address)
	if err != nil {
		return nil, err
	}

	var linkedAccount *PublicAccount = nil

	if dto.Account.AccountType != UnlinkedAccount && dto.Account.AccountType != RemoteUnlinkedAccount {
		linkedAccount, err = NewAccountFromPublicKey(dto.Account.LinkedAccountPublicKey, add.Type)
		if err != nil {
			return nil, err
		}
	}

	acc := &AccountInfo{
		Address:         add,
		AddressHeight:   dto.Account.AddressHeight.toStruct(),
		PublicKey:       dto.Account.PublicKey,
		PublicKeyHeight: dto.Account.PublicKeyHeight.toStruct(),
		AccountType:     dto.Account.AccountType,
		LinkedAccount:   linkedAccount,
		Mosaics:         ms,
		Reputation:      repConfig.defaultReputation,
	}

	if dto.Account.Reputation != nil {
		acc.Reputation = dto.Account.Reputation.toFloat(repConfig)
	}

	return acc, nil
}

type accountInfoDTOs []*accountInfoDTO

func (a accountInfoDTOs) toStruct(repConfig *reputationConfig) ([]*AccountInfo, error) {
	var (
		accountInfos = make([]*AccountInfo, len(a))
		err          error
	)

	for idx, dto := range a {
		accountInfos[idx], err = dto.toStruct(repConfig)
		if err != nil {
			return nil, err
		}
	}

	return accountInfos, nil
}

type multisigAccountInfoDTO struct {
	Multisig struct {
		Account          string   `json:"account"`
		MinApproval      int32    `json:"minApproval"`
		MinRemoval       int32    `json:"minRemoval"`
		Cosignatories    []string `json:"cosignatories"`
		MultisigAccounts []string `json:"multisigAccounts"`
	} `json:"multisig"`
}

func (dto *multisigAccountInfoDTO) toStruct(networkType NetworkType) (*MultisigAccountInfo, error) {
	cs := make([]*PublicAccount, len(dto.Multisig.Cosignatories))
	ms := make([]*PublicAccount, len(dto.Multisig.MultisigAccounts))

	acc, err := NewAccountFromPublicKey(dto.Multisig.Account, networkType)
	if err != nil {
		return nil, err
	}

	for i, c := range dto.Multisig.Cosignatories {
		cs[i], err = NewAccountFromPublicKey(c, networkType)
		if err != nil {
			return nil, err
		}
	}

	for i, m := range dto.Multisig.MultisigAccounts {
		ms[i], err = NewAccountFromPublicKey(m, networkType)
		if err != nil {
			return nil, err
		}
	}

	return &MultisigAccountInfo{
		Account:          *acc,
		MinApproval:      dto.Multisig.MinApproval,
		MinRemoval:       dto.Multisig.MinRemoval,
		Cosignatories:    cs,
		MultisigAccounts: ms,
	}, nil
}

type multisigAccountGraphInfoDTOEntry struct {
	Level     int32                    `json:"level"`
	Multisigs []multisigAccountInfoDTO `json:"multisigEntries"`
}

type multisigAccountGraphInfoDTOS []multisigAccountGraphInfoDTOEntry

func (dto multisigAccountGraphInfoDTOS) toStruct(networkType NetworkType) (*MultisigAccountGraphInfo, error) {
	var (
		ms  = make(map[int32][]*MultisigAccountInfo)
		err error
	)

	for _, m := range dto {
		mAccInfos := make([]*MultisigAccountInfo, len(m.Multisigs))

		for idx, c := range m.Multisigs {
			mAccInfos[idx], err = c.toStruct(networkType)
			if err != nil {
				return nil, err
			}
		}

		ms[m.Level] = mAccInfos
	}

	return &MultisigAccountGraphInfo{ms}, nil
}

type addresses struct {
	Addresses []*Address
}

func (ref *addresses) MarshalJSON() (buf []byte, err error) {
	buf = []byte(`{"addresses":[`)
	for i, address := range ref.Addresses {
		b := []byte(`"` + address.Address + `"`)
		if i > 0 {
			buf = append(buf, ',')
		}
		buf = append(buf, b...)
	}

	buf = append(buf, ']', '}')
	return
}

func (ref *addresses) UnmarshalJSON([]byte) error {
	return nil
}

// generateEncodedAddress convert publicKey to address
func generateEncodedAddress(pKey string, version NetworkType) (string, error) {
	// step 1: sha3 hash of the public key
	pKeyD, err := hex.DecodeString(pKey)
	if err != nil {
		return "", err
	}
	sha3PublicKeyHash, err := crypto.HashesSha3_256(pKeyD)
	if err != nil {
		return "", err
	}
	// step 2: ripemd160 hash of (1)
	ripemd160StepOneHash, err := crypto.HashesRipemd160(sha3PublicKeyHash)
	if err != nil {
		return "", err
	}

	// step 3: add version byte in front of (2)
	versionPrefixedRipemd160Hash := append([]byte{uint8(version)}, ripemd160StepOneHash...)

	// step 4: get the checksum of (3)
	stepThreeChecksum, err := GenerateChecksum(versionPrefixedRipemd160Hash)
	if err != nil {
		return "", err
	}

	// step 5: concatenate (3) and (4)
	concatStepThreeAndStepSix := append(versionPrefixedRipemd160Hash, stepThreeChecksum...)

	// step 6: base32 encode (5)
	return base32.StdEncoding.EncodeToString(concatStepThreeAndStepSix), nil
}

type accountNamesDTO struct {
	Names   []string `json:"names"`
	Address string   `json:"address"`
}

type accountNamesDTOs []*accountNamesDTO

func (m *accountNamesDTO) toStruct() (*AccountName, error) {

	address, err := NewAddressFromBase32(m.Address)
	if err != nil {
		return nil, err
	}
	return &AccountName{
		Address: address,
		Names:   m.Names,
	}, nil
}

func (m *accountNamesDTOs) toStruct() ([]*AccountName, error) {
	dtos := *m
	accNames := make([]*AccountName, 0, len(dtos))

	for _, dto := range dtos {
		accName, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		accNames = append(accNames, accName)
	}

	return accNames, nil
}

type harvesterDTO struct {
	Harvester struct {
		Key                    string    `json:"key"`
		Owner                  string    `json:"owner"`
		Address                string    `json:"address"`
		DisabledHeight         uint64DTO `json:"disabledHeight"`
		LastSigningBlockHeight uint64DTO `json:"lastSigningBlockHeight"`
		EffectiveBalance       uint64DTO `json:"effectiveBalance"`
		CanHarvest             bool      `json:"canHarvest"`
		Activity               float64   `json:"activity"`
		Greed                  float64   `json:"greed"`
	} `json:"harvester"`
}

func (ref *harvesterDTO) toStruct() (*Harvester, error) {
	address, err := NewAddressFromBase32(ref.Harvester.Address)
	if err != nil {
		return nil, err
	}

	return &Harvester{
		Key:                    ref.Harvester.Key,
		Owner:                  ref.Harvester.Owner,
		Address:                address,
		DisabledHeight:         ref.Harvester.DisabledHeight.toStruct(),
		LastSigningBlockHeight: ref.Harvester.LastSigningBlockHeight.toStruct(),
		EffectiveBalance:       ref.Harvester.EffectiveBalance.toStruct(),
		CanHarvest:             ref.Harvester.CanHarvest,
		Activity:               ref.Harvester.Activity,
		Greed:                  ref.Harvester.Greed,
	}, nil
}

type harvesterDTOs []*harvesterDTO

type harvestersPageDTO struct {
	Harvesters []harvesterDTO `json:"data"`

	Pagination struct {
		TotalEntries uint64 `json:"totalEntries"`
		PageNumber   uint64 `json:"pageNumber"`
		PageSize     uint64 `json:"pageSize"`
		TotalPages   uint64 `json:"totalPages"`
	} `json:"pagination"`
}

func (ref *harvestersPageDTO) toStruct() (*HarvestersPage, error) {
	page := &HarvestersPage{
		Harvesters: make([]*Harvester, len(ref.Harvesters)),
		Pagination: &Pagination{
			TotalEntries: ref.Pagination.TotalEntries,
			PageNumber:   ref.Pagination.PageNumber,
			PageSize:     ref.Pagination.PageSize,
			TotalPages:   ref.Pagination.TotalPages,
		},
	}

	var err error
	for i, h := range ref.Harvesters {
		page.Harvesters[i], err = h.toStruct()
		if err != nil {
			return nil, err
		}
	}

	return page, nil
}
