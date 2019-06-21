# sdk
--
    import "github.com/go-xpx-catapult-sdk/sdk"

Package sdk provides a client library for the Catapult REST API.

## Usage

```go
const (
	AddressSize                              int = 25
	AmountSize                               int = 8
	KeySize                                  int = 32
	Hash256                                  int = 32
	MosaicSize                               int = 8
	NamespaceSize                            int = 8
	SizeSize                                 int = 4
	SignerSize                               int = KeySize
	SignatureSize                            int = 64
	VersionSize                              int = 2
	TypeSize                                 int = 2
	MaxFeeSize                               int = 8
	DeadLineSize                             int = 8
	DurationSize                             int = 8
	TransactionHeaderSize                    int = SizeSize + SignerSize + SignatureSize + VersionSize + TypeSize + MaxFeeSize + DeadLineSize
	PropertyTypeSize                         int = 2
	PropertyModificationTypeSize             int = 1
	AccountPropertiesAddressModificationSize int = PropertyModificationTypeSize + AddressSize
	AccountPropertiesMosaicModificationSize  int = PropertyModificationTypeSize + MosaicSize
	AccountPropertiesEntityModificationSize  int = PropertyModificationTypeSize + TypeSize
	AccountPropertyAddressHeader             int = TransactionHeaderSize + PropertyTypeSize
	AccountPropertyMosaicHeader              int = TransactionHeaderSize + PropertyTypeSize
	AccountPropertyEntityTypeHeader          int = TransactionHeaderSize + PropertyTypeSize
	LinkActionSize                           int = 1
	AccountLinkTransactionSize               int = TransactionHeaderSize + KeySize + LinkActionSize
	AliasActionSize                          int = 1
	AliasTransactionHeader                   int = TransactionHeaderSize + NamespaceSize + AliasActionSize
	AggregateBondedHeader                    int = TransactionHeaderSize + SizeSize
	HashTypeSize                             int = 1
	LockSize                                 int = TransactionHeaderSize + MosaicSize + AmountSize + DurationSize + Hash256
	MetadataTypeSize                         int = 1
	MetadataHeaderSize                       int = TransactionHeaderSize + MetadataTypeSize
	ModificationsSizeSize                    int = 1
	ModifyContractHeaderSize                 int = TransactionHeaderSize + DurationSize + Hash256 + 3*ModificationsSizeSize
	MinApprovalSize                          int = 1
	MinRemovalSize                           int = 1
	ModifyMultisigHeaderSize                 int = TransactionHeaderSize + MinApprovalSize + MinRemovalSize + ModificationsSizeSize
	MosaicNonceSize                          int = 4
	MosaicPropertySize                       int = 4
	MosaicDefinitionTransactionSize          int = TransactionHeaderSize + MosaicNonceSize + MosaicSize + DurationSize + MosaicPropertySize
	MosaicSupplyDirectionSize                int = 1
	MosaicSupplyChangeTransactionSize        int = TransactionHeaderSize + MosaicSize + AmountSize + MosaicSupplyDirectionSize
	NamespaceTypeSize                        int = 1
	NamespaceNameSizeSize                    int = 1
	RegisterNamespaceHeaderSize              int = TransactionHeaderSize + NamespaceTypeSize + DurationSize + NamespaceSize + NamespaceNameSizeSize
	SecretLockSize                           int = TransactionHeaderSize + MosaicSize + AmountSize + DurationSize + HashTypeSize + Hash256 + AddressSize
	ProofSizeSize                            int = 2
	SecretProofHeaderSize                    int = TransactionHeaderSize + HashTypeSize + Hash256 + ProofSizeSize
	MosaicsSizeSize                          int = 1
	MessageSizeSize                          int = 2
	TransferHeaderSize                       int = TransactionHeaderSize + AddressSize + MosaicsSizeSize + MessageSizeSize
)
```

```go
const EmptyPublicKey = "0000000000000000000000000000000000000000000000000000000000000000"
```

```go
const LevyMutable = 0x04
```

```go
const NUM_CHECKSUM_BYTES = 4
```

```go
const NamespaceBit uint64 = 1 << 63
```

```go
const Supply_Mutable = 0x01
```

```go
const TimestampNemesisBlockMilliseconds int64 = 1459468800 * 1000
```

```go
const Transferable = 0x02
```

```go
const WebsocketReconnectionDefaultTimeout = time.Second * 5
```

```go
var (
	ErrResourceNotFound              = newRespError("resource is not found")
	ErrArgumentNotValid              = newRespError("argument is not valid")
	ErrInvalidRequest                = newRespError("request is not valid")
	ErrInternalError                 = newRespError("response is nil")
	ErrNotAcceptedResponseStatusCode = newRespError("not accepted response status code")
)
```
Catapult REST API errors

```go
var (
	ErrMetadataEmptyAddresses    = errors.New("list adresses ids must not by empty")
	ErrMetadataNilAdress         = errors.New("adress must not be blank")
	ErrMetadataEmptyMosaicIds    = errors.New("list mosaics ids must not by empty")
	ErrMetadataNilMosaicId       = errors.New("mosaicId must not be nil")
	ErrMetadataEmptyNamespaceIds = errors.New("list namespaces ids must not by empty")
	ErrMetadataNilNamespaceId    = errors.New("namespaceId must not be nil")
)
```
Metadata errors

```go
var (
	ErrNilAssetId            = errors.New("assetId must not be nil")
	ErrEmptyAssetIds         = errors.New("list blockchain ids must not by empty")
	ErrUnknownBlockchainType = errors.New("Not supported Blockchain Type")
)
```
Common errors

```go
var (
	ErrEmptyMosaicIds        = errors.New("list mosaics ids must not by empty")
	ErrNilMosaicId           = errors.New("mosaicId must not be nil")
	ErrWrongBitMosaicId      = errors.New("mosaicId has 64th bit")
	ErrInvalidOwnerPublicKey = errors.New("public owner key is invalid")
	ErrNilMosaicProperties   = errors.New("mosaic properties must not be nil")
)
```
Mosaic errors

```go
var (
	ErrNamespaceTooManyPart = errors.New("too many parts")
	ErrNilNamespaceId       = errors.New("namespaceId is nil or zero")
	ErrWrongBitNamespaceId  = errors.New("namespaceId doesn't have 64th bit")
	ErrEmptyNamespaceIds    = errors.New("list namespace ids must not by empty")
	ErrInvalidNamespaceName = errors.New("namespace name is invalid")
)
```
Namespace errors

```go
var (
	ErrNilOrZeroHeight = errors.New("block height should not be nil or zero")
	ErrNilOrZeroLimit  = errors.New("limit should not be nil or zero")
)
```
Blockchain errors

```go
var (
	ErrEmptyAddressesIds = errors.New("list of addresses should not be empty")
	ErrNilAddress        = errors.New("address is nil")
	ErrBlankAddress      = errors.New("address is blank")
	ErrNilAccount        = errors.New("account should not be nil")
	ErrInvalidAddress    = errors.New("wrong address")
)
```
plain errors

```go
var (
	ErrInvalidReputationConfig = errors.New("default reputation should be greater than 0 and less than 1")
)
```
reputations error

```go
var XemMosaicId, _ = NewMosaicId(0x0DC67FBE1CAD29E3)
```
mosaic id for XEM mosaic

```go
var XpxMosaicId, _ = NewMosaicId(0x0DC67FBE1CAD29E3)
```
mosaic id for XPX mosaic

#### func  ExtractVersion

```go
func ExtractVersion(version uint64) uint8
```

#### func  GenerateChecksum

```go
func GenerateChecksum(b []byte) ([]byte, error)
```

#### func  NewReputationConfig

```go
func NewReputationConfig(minInter uint64, defaultRep float64) (*reputationConfig, error)
```

#### type AbstractTransaction

```go
type AbstractTransaction struct {
	*TransactionInfo
	NetworkType NetworkType
	Deadline    *Deadline
	Type        TransactionType
	Version     TransactionVersion
	MaxFee      Amount
	Signature   string
	Signer      *PublicAccount
}
```


#### func (*AbstractTransaction) HasMissingSignatures

```go
func (tx *AbstractTransaction) HasMissingSignatures() bool
```

#### func (*AbstractTransaction) IsConfirmed

```go
func (tx *AbstractTransaction) IsConfirmed() bool
```

#### func (*AbstractTransaction) IsUnannounced

```go
func (tx *AbstractTransaction) IsUnannounced() bool
```

#### func (*AbstractTransaction) IsUnconfirmed

```go
func (tx *AbstractTransaction) IsUnconfirmed() bool
```

#### func (*AbstractTransaction) String

```go
func (tx *AbstractTransaction) String() string
```

#### func (*AbstractTransaction) ToAggregate

```go
func (tx *AbstractTransaction) ToAggregate(signer *PublicAccount)
```

#### type Account

```go
type Account struct {
	*PublicAccount
	*crypto.KeyPair
}
```


#### func  NewAccount

```go
func NewAccount(networkType NetworkType) (*Account, error)
```
returns new Account generated for passed NetworkType

#### func  NewAccountFromPrivateKey

```go
func NewAccountFromPrivateKey(pKey string, networkType NetworkType) (*Account, error)
```
returns new Account from private key for passed NetworkType

#### func (*Account) DecryptMessage

```go
func (a *Account) DecryptMessage(encryptedMessage *SecureMessage, senderPublicAccount *PublicAccount) (*PlainMessage, error)
```

#### func (*Account) EncryptMessage

```go
func (a *Account) EncryptMessage(message string, recipientPublicAccount *PublicAccount) (*SecureMessage, error)
```

#### func (*Account) Sign

```go
func (a *Account) Sign(tx Transaction) (*SignedTransaction, error)
```

#### func (*Account) SignCosignatureTransaction

```go
func (a *Account) SignCosignatureTransaction(tx *CosignatureTransaction) (*CosignatureSignedTransaction, error)
```

#### func (*Account) SignWithCosignatures

```go
func (a *Account) SignWithCosignatures(tx *AggregateTransaction, cosignatories []*Account) (*SignedTransaction, error)
```
sign AggregateTransaction with current Account and with every passed cosignatory
Account's returns announced Aggregate SignedTransaction

#### type AccountInfo

```go
type AccountInfo struct {
	Address         *Address
	AddressHeight   Height
	PublicKey       string
	PublicKeyHeight Height
	AccountType     AccountType
	LinkedAccount   *PublicAccount
	Mosaics         []*Mosaic
	Reputation      float64
}
```


#### func (*AccountInfo) String

```go
func (a *AccountInfo) String() string
```

#### type AccountLinkAction

```go
type AccountLinkAction uint8
```


```go
const (
	AccountLink AccountLinkAction = iota
	AccountUnlink
)
```
AccountLinkAction enums

#### type AccountLinkTransaction

```go
type AccountLinkTransaction struct {
	AbstractTransaction
	RemoteAccount *PublicAccount
	LinkAction    AccountLinkAction
}
```


#### func  NewAccountLinkTransaction

```go
func NewAccountLinkTransaction(deadline *Deadline, remoteAccount *PublicAccount, linkAction AccountLinkAction, networkType NetworkType) (*AccountLinkTransaction, error)
```
returns AccountLinkTransaction from passed PublicAccount and AccountLinkAction

#### func (*AccountLinkTransaction) GetAbstractTransaction

```go
func (tx *AccountLinkTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*AccountLinkTransaction) Size

```go
func (tx *AccountLinkTransaction) Size() int
```

#### func (*AccountLinkTransaction) String

```go
func (tx *AccountLinkTransaction) String() string
```

#### type AccountProperties

```go
type AccountProperties struct {
	Address            *Address
	AllowedAddresses   []*Address
	AllowedMosaicId    []*MosaicId
	AllowedEntityTypes []TransactionType
	BlockedAddresses   []*Address
	BlockedMosaicId    []*MosaicId
	BlockedEntityTypes []TransactionType
}
```


#### func (*AccountProperties) String

```go
func (a *AccountProperties) String() string
```

#### type AccountPropertiesAddressModification

```go
type AccountPropertiesAddressModification struct {
	ModificationType PropertyModificationType
	Address          *Address
}
```


#### func (*AccountPropertiesAddressModification) String

```go
func (mod *AccountPropertiesAddressModification) String() string
```

#### type AccountPropertiesAddressTransaction

```go
type AccountPropertiesAddressTransaction struct {
	AbstractTransaction
	PropertyType  PropertyType
	Modifications []*AccountPropertiesAddressModification
}
```


#### func  NewAccountPropertiesAddressTransaction

```go
func NewAccountPropertiesAddressTransaction(deadline *Deadline, propertyType PropertyType,
	modifications []*AccountPropertiesAddressModification, networkType NetworkType) (*AccountPropertiesAddressTransaction, error)
```
returns AccountPropertiesAddressTransaction from passed PropertyType and
AccountPropertiesAddressModification's

#### func (*AccountPropertiesAddressTransaction) GetAbstractTransaction

```go
func (tx *AccountPropertiesAddressTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*AccountPropertiesAddressTransaction) Size

```go
func (tx *AccountPropertiesAddressTransaction) Size() int
```

#### func (*AccountPropertiesAddressTransaction) String

```go
func (tx *AccountPropertiesAddressTransaction) String() string
```

#### type AccountPropertiesEntityTypeModification

```go
type AccountPropertiesEntityTypeModification struct {
	ModificationType PropertyModificationType
	EntityType       TransactionType
}
```


#### func (*AccountPropertiesEntityTypeModification) String

```go
func (mod *AccountPropertiesEntityTypeModification) String() string
```

#### type AccountPropertiesEntityTypeTransaction

```go
type AccountPropertiesEntityTypeTransaction struct {
	AbstractTransaction
	PropertyType  PropertyType
	Modifications []*AccountPropertiesEntityTypeModification
}
```


#### func  NewAccountPropertiesEntityTypeTransaction

```go
func NewAccountPropertiesEntityTypeTransaction(deadline *Deadline, propertyType PropertyType,
	modifications []*AccountPropertiesEntityTypeModification, networkType NetworkType) (*AccountPropertiesEntityTypeTransaction, error)
```
returns AccountPropertiesEntityTypeTransaction from passed PropertyType and
AccountPropertiesEntityTypeModification's

#### func (*AccountPropertiesEntityTypeTransaction) GetAbstractTransaction

```go
func (tx *AccountPropertiesEntityTypeTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*AccountPropertiesEntityTypeTransaction) Size

```go
func (tx *AccountPropertiesEntityTypeTransaction) Size() int
```

#### func (*AccountPropertiesEntityTypeTransaction) String

```go
func (tx *AccountPropertiesEntityTypeTransaction) String() string
```

#### type AccountPropertiesMosaicModification

```go
type AccountPropertiesMosaicModification struct {
	ModificationType PropertyModificationType
	AssetId          AssetId
}
```


#### func (*AccountPropertiesMosaicModification) String

```go
func (mod *AccountPropertiesMosaicModification) String() string
```

#### type AccountPropertiesMosaicTransaction

```go
type AccountPropertiesMosaicTransaction struct {
	AbstractTransaction
	PropertyType  PropertyType
	Modifications []*AccountPropertiesMosaicModification
}
```


#### func  NewAccountPropertiesMosaicTransaction

```go
func NewAccountPropertiesMosaicTransaction(deadline *Deadline, propertyType PropertyType,
	modifications []*AccountPropertiesMosaicModification, networkType NetworkType) (*AccountPropertiesMosaicTransaction, error)
```
returns AccountPropertiesMosaicTransaction from passed PropertyType and
AccountPropertiesMosaicModification's

#### func (*AccountPropertiesMosaicTransaction) GetAbstractTransaction

```go
func (tx *AccountPropertiesMosaicTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*AccountPropertiesMosaicTransaction) Size

```go
func (tx *AccountPropertiesMosaicTransaction) Size() int
```

#### func (*AccountPropertiesMosaicTransaction) String

```go
func (tx *AccountPropertiesMosaicTransaction) String() string
```

#### type AccountService

```go
type AccountService service
```


#### func (*AccountService) AggregateBondedTransactions

```go
func (a *AccountService) AggregateBondedTransactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]*AggregateTransaction, error)
```
returns an array of AggregateTransaction's where passed account is signer or
cosigner

#### func (*AccountService) GetAccountInfo

```go
func (a *AccountService) GetAccountInfo(ctx context.Context, address *Address) (*AccountInfo, error)
```

#### func (*AccountService) GetAccountProperties

```go
func (a *AccountService) GetAccountProperties(ctx context.Context, address *Address) (*AccountProperties, error)
```

#### func (*AccountService) GetAccountsInfo

```go
func (a *AccountService) GetAccountsInfo(ctx context.Context, addresses ...*Address) ([]*AccountInfo, error)
```

#### func (*AccountService) GetAccountsProperties

```go
func (a *AccountService) GetAccountsProperties(ctx context.Context, addresses ...*Address) ([]*AccountProperties, error)
```

#### func (*AccountService) GetMultisigAccountGraphInfo

```go
func (a *AccountService) GetMultisigAccountGraphInfo(ctx context.Context, address *Address) (*MultisigAccountGraphInfo, error)
```

#### func (*AccountService) GetMultisigAccountInfo

```go
func (a *AccountService) GetMultisigAccountInfo(ctx context.Context, address *Address) (*MultisigAccountInfo, error)
```

#### func (*AccountService) IncomingTransactions

```go
func (a *AccountService) IncomingTransactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]Transaction, error)
```
returns an array of Transaction's for which passed account is receiver

#### func (*AccountService) OutgoingTransactions

```go
func (a *AccountService) OutgoingTransactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]Transaction, error)
```
returns an array of Transaction's for which passed account is sender

#### func (*AccountService) Transactions

```go
func (a *AccountService) Transactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]Transaction, error)
```
returns an array of confirmed Transaction's for which passed account is sender
or receiver.

#### func (*AccountService) UnconfirmedTransactions

```go
func (a *AccountService) UnconfirmedTransactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]Transaction, error)
```
returns an array of confirmed Transaction's for which passed account is sender
or receiver. unconfirmed transactions are those transactions that have not yet
been included in a block. they are not guaranteed to be included in any block.

#### type AccountTransactionsOption

```go
type AccountTransactionsOption struct {
	PageSize int              `url:"pageSize,omitempty"`
	Id       string           `url:"id,omitempty"`
	Ordering TransactionOrder `url:"ordering,omitempty"`
}
```


#### type AccountType

```go
type AccountType uint8
```


```go
const (
	UnlinkedAccount AccountType = iota
	MainAccount
	RemoteAccount
	RemoteUnlinkedAccount
)
```
AccountType enums

#### type Address

```go
type Address struct {
	Type    NetworkType
	Address string
}
```


#### func  EncodedStringToAddresses

```go
func EncodedStringToAddresses(addresses ...string) ([]*Address, error)
```

#### func  NewAddress

```go
func NewAddress(address string, networkType NetworkType) *Address
```
returns Address from passed address string for passed NetworkType

#### func  NewAddressFromBase32

```go
func NewAddressFromBase32(encoded string) (*Address, error)
```

#### func  NewAddressFromNamespace

```go
func NewAddressFromNamespace(namespaceId *NamespaceId) (*Address, error)
```
returns new Address from namespace identifier

#### func  NewAddressFromPublicKey

```go
func NewAddressFromPublicKey(pKey string, networkType NetworkType) (*Address, error)
```
returns an Address from public key for passed NetworkType

#### func  NewAddressFromRaw

```go
func NewAddressFromRaw(address string) (*Address, error)
```
returns Address from passed address string

#### func (*Address) Pretty

```go
func (ad *Address) Pretty() string
```

#### func (*Address) String

```go
func (ad *Address) String() string
```

#### type AddressAliasTransaction

```go
type AddressAliasTransaction struct {
	AliasTransaction
	Address *Address
}
```


#### func  NewAddressAliasTransaction

```go
func NewAddressAliasTransaction(deadline *Deadline, address *Address, namespaceId *NamespaceId, actionType AliasActionType, networkType NetworkType) (*AddressAliasTransaction, error)
```
returns AddressAliasTransaction from passed Address, NamespaceId and
AliasActionType

#### func (*AddressAliasTransaction) Size

```go
func (tx *AddressAliasTransaction) Size() int
```

#### func (*AddressAliasTransaction) String

```go
func (tx *AddressAliasTransaction) String() string
```

#### type AddressMetadataInfo

```go
type AddressMetadataInfo struct {
	MetadataInfo
	Address *Address
}
```


#### type AggregateTransaction

```go
type AggregateTransaction struct {
	AbstractTransaction
	InnerTransactions []Transaction
	Cosignatures      []*AggregateTransactionCosignature
}
```


#### func  NewBondedAggregateTransaction

```go
func NewBondedAggregateTransaction(deadline *Deadline, innerTxs []Transaction, networkType NetworkType) (*AggregateTransaction, error)
```
returns bounded AggregateTransaction from passed array of transactions to be
included in

#### func  NewCompleteAggregateTransaction

```go
func NewCompleteAggregateTransaction(deadline *Deadline, innerTxs []Transaction, networkType NetworkType) (*AggregateTransaction, error)
```
returns complete AggregateTransaction from passed array of own Transaction's to
be included in

#### func (*AggregateTransaction) GetAbstractTransaction

```go
func (tx *AggregateTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*AggregateTransaction) Size

```go
func (tx *AggregateTransaction) Size() int
```

#### func (*AggregateTransaction) String

```go
func (tx *AggregateTransaction) String() string
```

#### type AggregateTransactionCosignature

```go
type AggregateTransactionCosignature struct {
	Signature string
	Signer    *PublicAccount
}
```


#### func (*AggregateTransactionCosignature) String

```go
func (agt *AggregateTransactionCosignature) String() string
```

#### type AliasActionType

```go
type AliasActionType uint8
```


```go
const (
	AliasLink AliasActionType = iota
	AliasUnlink
)
```
AliasActionType enums

#### type AliasTransaction

```go
type AliasTransaction struct {
	AbstractTransaction
	ActionType  AliasActionType
	NamespaceId *NamespaceId
}
```


#### func (*AliasTransaction) GetAbstractTransaction

```go
func (tx *AliasTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*AliasTransaction) Size

```go
func (tx *AliasTransaction) Size() int
```

#### func (*AliasTransaction) String

```go
func (tx *AliasTransaction) String() string
```

#### type AliasType

```go
type AliasType uint8
```


```go
const (
	NoneAliasType AliasType = iota
	MosaicAliasType
	AddressAliasType
)
```
AliasType enums

#### type Amount

```go
type Amount = baseInt64
```


#### type AssetId

```go
type AssetId interface {
	fmt.Stringer
	Type() AssetIdType
	Id() uint64
	Equals(AssetId) bool
	// contains filtered or unexported methods
}
```


#### type AssetIdType

```go
type AssetIdType uint8
```


```go
const (
	NamespaceAssetIdType AssetIdType = iota
	MosaicAssetIdType
)
```
AssetIdType enums

#### type BlockInfo

```go
type BlockInfo struct {
	NetworkType
	Hash                  string
	GenerationHash        string
	TotalFee              Amount
	NumTransactions       uint64
	Signature             string
	Signer                *PublicAccount
	Version               uint8
	Type                  uint64
	Height                Height
	Timestamp             *Timestamp
	Difficulty            Difficulty
	FeeMultiplier         uint32
	PreviousBlockHash     string
	BlockTransactionsHash string
	BlockReceiptsHash     string
	StateHash             string
	Beneficiary           *PublicAccount
}
```


#### func  MapBlock

```go
func MapBlock(m []byte) (*BlockInfo, error)
```

#### func (*BlockInfo) String

```go
func (b *BlockInfo) String() string
```

#### type BlockMapper

```go
type BlockMapper interface {
	MapBlock(m []byte) (*BlockInfo, error)
}
```


#### type BlockMapperFn

```go
type BlockMapperFn func(m []byte) (*BlockInfo, error)
```


#### func (BlockMapperFn) MapBlock

```go
func (p BlockMapperFn) MapBlock(m []byte) (*BlockInfo, error)
```

#### type BlockchainService

```go
type BlockchainService service
```


#### func (*BlockchainService) GetBlockByHeight

```go
func (b *BlockchainService) GetBlockByHeight(ctx context.Context, height Height) (*BlockInfo, error)
```
returns BlockInfo for passed block's height

#### func (*BlockchainService) GetBlockTransactions

```go
func (b *BlockchainService) GetBlockTransactions(ctx context.Context, height Height) ([]Transaction, error)
```
returns Transaction's inside of block at passed height

#### func (*BlockchainService) GetBlockchainHeight

```go
func (b *BlockchainService) GetBlockchainHeight(ctx context.Context) (Height, error)
```

#### func (*BlockchainService) GetBlockchainScore

```go
func (b *BlockchainService) GetBlockchainScore(ctx context.Context) (*ChainScore, error)
```

#### func (*BlockchainService) GetBlockchainStorage

```go
func (b *BlockchainService) GetBlockchainStorage(ctx context.Context) (*BlockchainStorageInfo, error)
```

#### func (*BlockchainService) GetBlocksByHeightWithLimit

```go
func (b *BlockchainService) GetBlocksByHeightWithLimit(ctx context.Context, height Height, limit Amount) ([]*BlockInfo, error)
```
returns BlockInfo's for range block height - (block height + limit) Example:
GetBlocksByHeightWithLimit(ctx, 1, 25) => [BlockInfo25, BlockInfo24, ...,
BlockInfo1]

#### type BlockchainStorageInfo

```go
type BlockchainStorageInfo struct {
	NumBlocks       int `json:"numBlocks"`
	NumTransactions int `json:"numTransactions"`
	NumAccounts     int `json:"numAccounts"`
}
```


#### func (*BlockchainStorageInfo) String

```go
func (b *BlockchainStorageInfo) String() string
```

#### type BlockchainTimestamp

```go
type BlockchainTimestamp struct {
}
```


#### func  NewBlockchainTimestamp

```go
func NewBlockchainTimestamp(milliseconds int64) *BlockchainTimestamp
```
returns new BlockchainTimestamp from passed milliseconds value

#### func (BlockchainTimestamp) String

```go
func (m BlockchainTimestamp) String() string
```

#### func (*BlockchainTimestamp) ToTimestamp

```go
func (t *BlockchainTimestamp) ToTimestamp() *Timestamp
```

#### type ChainScore

```go
type ChainScore [2]uint64
```


#### func  NewChainScore

```go
func NewChainScore(scoreLow uint64, scoreHigh uint64) *ChainScore
```
returns new ChainScore from passed low and high score

#### func (*ChainScore) String

```go
func (m *ChainScore) String() string
```

#### type Client

```go
type Client struct {

	// Services for communicating to the Catapult REST APIs
	Blockchain  *BlockchainService
	Mosaic      *MosaicService
	Namespace   *NamespaceService
	Network     *NetworkService
	Transaction *TransactionService
	Resolve     *ResolverService
	Account     *AccountService
	Contract    *ContractService
	Metadata    *MetadataService
}
```

Catapult API Client configuration

#### func  NewClient

```go
func NewClient(httpClient *http.Client, conf *Config) *Client
```
returns catapult http.Client from passed existing client and configuration if
passed client is nil, http.DefaultClient will be used

#### type Config

```go
type Config struct {
	BaseURLs              []*url.URL
	UsedBaseUrl           *url.URL
	WsReconnectionTimeout time.Duration
	NetworkType
}
```

Provides service configuration

#### func  NewConfig

```go
func NewConfig(baseUrls []string, networkType NetworkType, wsReconnectionTimeout time.Duration) (*Config, error)
```
returns config for HTTP Client from passed node url and network type

#### func  NewConfigWithReputation

```go
func NewConfigWithReputation(baseUrls []string, networkType NetworkType, repConf *reputationConfig, wsReconnectionTimeout time.Duration) (*Config, error)
```

#### type ConfirmedAddedMapper

```go
type ConfirmedAddedMapper interface {
	MapConfirmedAdded(m []byte) (Transaction, error)
}
```


#### func  NewConfirmedAddedMapper

```go
func NewConfirmedAddedMapper(mapTransactionFunc mapTransactionFunc) ConfirmedAddedMapper
```

#### type ContractInfo

```go
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
```


#### type ContractService

```go
type ContractService service
```


#### func (*ContractService) GetContractsByAddress

```go
func (ref *ContractService) GetContractsByAddress(ctx context.Context, address string) ([]*ContractInfo, error)
```

#### func (*ContractService) GetContractsInfo

```go
func (ref *ContractService) GetContractsInfo(ctx context.Context, contractPubKeys ...string) ([]*ContractInfo, error)
```

#### type CosignatureMapper

```go
type CosignatureMapper interface {
	MapCosignature(m []byte) (*SignerInfo, error)
}
```


#### type CosignatureMapperFn

```go
type CosignatureMapperFn func(m []byte) (*SignerInfo, error)
```


#### func (CosignatureMapperFn) MapCosignature

```go
func (p CosignatureMapperFn) MapCosignature(m []byte) (*SignerInfo, error)
```

#### type CosignatureSignedTransaction

```go
type CosignatureSignedTransaction struct {
	ParentHash Hash   `json:"parentHash"`
	Signature  string `json:"signature"`
	Signer     string `json:"signer"`
}
```


#### type CosignatureTransaction

```go
type CosignatureTransaction struct {
	TransactionToCosign *AggregateTransaction
}
```


#### func  NewCosignatureTransaction

```go
func NewCosignatureTransaction(txToCosign *AggregateTransaction) (*CosignatureTransaction, error)
```
returns a CosignatureTransaction from passed AggregateTransaction

#### func  NewCosignatureTransactionFromHash

```go
func NewCosignatureTransactionFromHash(hash Hash) *CosignatureTransaction
```
returns a CosignatureTransaction from passed hash of bounded
AggregateTransaction

#### func (*CosignatureTransaction) String

```go
func (tx *CosignatureTransaction) String() string
```

#### type Deadline

```go
type Deadline struct {
	Timestamp
}
```


#### func  NewDeadline

```go
func NewDeadline(delta time.Duration) *Deadline
```
returns new Deadline from passed duration

#### func  NewDeadlineFromBlockchainTimestamp

```go
func NewDeadlineFromBlockchainTimestamp(timestamp *BlockchainTimestamp) *Deadline
```
returns new Deadline from passed BlockchainTimestamp

#### type Difficulty

```go
type Difficulty = baseInt64
```


#### type Duration

```go
type Duration = baseInt64
```


#### type Hash

```go
type Hash string
```


#### func (Hash) String

```go
func (h Hash) String() string
```

#### type HashType

```go
type HashType uint8
```


```go
const (
	/// Input is hashed using Sha-3-256.
	SHA3_256 HashType = iota
	/// Input is hashed using Keccak-256.
	KECCAK_256
	/// Input is hashed twice: first with SHA-256 and then with RIPEMD-160.
	HASH_160
	/// Input is hashed twice with SHA-256.
	SHA_256
)
```

#### func (HashType) String

```go
func (ht HashType) String() string
```

#### type Height

```go
type Height = baseInt64
```


#### type HttpError

```go
type HttpError struct {
	StatusCode int
}
```


#### type LockFundsTransaction

```go
type LockFundsTransaction struct {
	AbstractTransaction
	*Mosaic
	Duration Duration
	*SignedTransaction
}
```


#### func  NewLockFundsTransaction

```go
func NewLockFundsTransaction(deadline *Deadline, mosaic *Mosaic, duration Duration, signedTx *SignedTransaction, networkType NetworkType) (*LockFundsTransaction, error)
```
returns a LockFundsTransaction from passed Mosaic, duration in blocks and
SignedTransaction

#### func (*LockFundsTransaction) GetAbstractTransaction

```go
func (tx *LockFundsTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*LockFundsTransaction) Size

```go
func (tx *LockFundsTransaction) Size() int
```

#### func (*LockFundsTransaction) String

```go
func (tx *LockFundsTransaction) String() string
```

#### type Message

```go
type Message interface {
	Type() MessageType
	Payload() []byte
	String() string
}
```


#### type MessageType

```go
type MessageType uint8
```


```go
const (
	PlainMessageType MessageType = iota
	SecureMessageType
)
```

#### type MetadataInfo

```go
type MetadataInfo struct {
	MetadataType MetadataType
	Fields       map[string]string
}
```


#### type MetadataModification

```go
type MetadataModification struct {
	Type  MetadataModificationType
	Key   string
	Value string
}
```


#### func (*MetadataModification) Size

```go
func (m *MetadataModification) Size() int
```

#### func (*MetadataModification) String

```go
func (m *MetadataModification) String() string
```

#### type MetadataModificationType

```go
type MetadataModificationType uint8
```


```go
const (
	AddMetadata MetadataModificationType = iota
	RemoveMetadata
)
```

#### func (MetadataModificationType) String

```go
func (t MetadataModificationType) String() string
```

#### type MetadataService

```go
type MetadataService service
```


#### func (*MetadataService) GetAddressMetadatasInfo

```go
func (ref *MetadataService) GetAddressMetadatasInfo(ctx context.Context, addresses ...string) ([]*AddressMetadataInfo, error)
```

#### func (*MetadataService) GetMetadataByAddress

```go
func (ref *MetadataService) GetMetadataByAddress(ctx context.Context, address string) (*AddressMetadataInfo, error)
```

#### func (*MetadataService) GetMetadataByMosaicId

```go
func (ref *MetadataService) GetMetadataByMosaicId(ctx context.Context, mosaicId *MosaicId) (*MosaicMetadataInfo, error)
```

#### func (*MetadataService) GetMetadataByNamespaceId

```go
func (ref *MetadataService) GetMetadataByNamespaceId(ctx context.Context, namespaceId *NamespaceId) (*NamespaceMetadataInfo, error)
```

#### func (*MetadataService) GetMosaicMetadatasInfo

```go
func (ref *MetadataService) GetMosaicMetadatasInfo(ctx context.Context, mosaicIds ...*MosaicId) ([]*MosaicMetadataInfo, error)
```

#### func (*MetadataService) GetNamespaceMetadatasInfo

```go
func (ref *MetadataService) GetNamespaceMetadatasInfo(ctx context.Context, namespaceIds ...*NamespaceId) ([]*NamespaceMetadataInfo, error)
```

#### type MetadataType

```go
type MetadataType uint8
```


```go
const (
	MetadataNone MetadataType = iota
	MetadataAddressType
	MetadataMosaicType
	MetadataNamespaceType
)
```

#### func (MetadataType) String

```go
func (t MetadataType) String() string
```

#### type ModifyContractTransaction

```go
type ModifyContractTransaction struct {
	AbstractTransaction
	DurationDelta Duration
	Hash          string
	Customers     []*MultisigCosignatoryModification
	Executors     []*MultisigCosignatoryModification
	Verifiers     []*MultisigCosignatoryModification
}
```


#### func  NewModifyContractTransaction

```go
func NewModifyContractTransaction(
	deadline *Deadline, durationDelta Duration, hash string,
	customers []*MultisigCosignatoryModification,
	executors []*MultisigCosignatoryModification,
	verifiers []*MultisigCosignatoryModification,
	networkType NetworkType) (*ModifyContractTransaction, error)
```
returns ModifyContractTransaction from passed duration delta in blocks, file
hash, arrays of customers, replicators and verificators

#### func (*ModifyContractTransaction) GetAbstractTransaction

```go
func (tx *ModifyContractTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*ModifyContractTransaction) Size

```go
func (tx *ModifyContractTransaction) Size() int
```

#### func (*ModifyContractTransaction) String

```go
func (tx *ModifyContractTransaction) String() string
```

#### type ModifyMetadataAddressTransaction

```go
type ModifyMetadataAddressTransaction struct {
	ModifyMetadataTransaction
	Address *Address
}
```


#### func  NewModifyMetadataAddressTransaction

```go
func NewModifyMetadataAddressTransaction(deadline *Deadline, address *Address, modifications []*MetadataModification, networkType NetworkType) (*ModifyMetadataAddressTransaction, error)
```
returns ModifyMetadataAddressTransaction from passed Address to be modified, and
an array of MetadataModification's

#### func (*ModifyMetadataAddressTransaction) Size

```go
func (tx *ModifyMetadataAddressTransaction) Size() int
```

#### func (*ModifyMetadataAddressTransaction) String

```go
func (tx *ModifyMetadataAddressTransaction) String() string
```

#### type ModifyMetadataMosaicTransaction

```go
type ModifyMetadataMosaicTransaction struct {
	ModifyMetadataTransaction
	MosaicId *MosaicId
}
```


#### func  NewModifyMetadataMosaicTransaction

```go
func NewModifyMetadataMosaicTransaction(deadline *Deadline, mosaicId *MosaicId, modifications []*MetadataModification, networkType NetworkType) (*ModifyMetadataMosaicTransaction, error)
```
returns ModifyMetadataMosaicTransaction from passed MosaicId to be modified, and
an array of MetadataModification's

#### func (*ModifyMetadataMosaicTransaction) Size

```go
func (tx *ModifyMetadataMosaicTransaction) Size() int
```

#### func (*ModifyMetadataMosaicTransaction) String

```go
func (tx *ModifyMetadataMosaicTransaction) String() string
```

#### type ModifyMetadataNamespaceTransaction

```go
type ModifyMetadataNamespaceTransaction struct {
	ModifyMetadataTransaction
	NamespaceId *NamespaceId
}
```


#### func  NewModifyMetadataNamespaceTransaction

```go
func NewModifyMetadataNamespaceTransaction(deadline *Deadline, namespaceId *NamespaceId, modifications []*MetadataModification, networkType NetworkType) (*ModifyMetadataNamespaceTransaction, error)
```
returns ModifyMetadataNamespaceTransaction from passed NamespaceId to be
modified, and an array of MetadataModification's

#### func (*ModifyMetadataNamespaceTransaction) Size

```go
func (tx *ModifyMetadataNamespaceTransaction) Size() int
```

#### func (*ModifyMetadataNamespaceTransaction) String

```go
func (tx *ModifyMetadataNamespaceTransaction) String() string
```

#### type ModifyMetadataTransaction

```go
type ModifyMetadataTransaction struct {
	AbstractTransaction
	MetadataType  MetadataType
	Modifications []*MetadataModification
}
```


#### func (*ModifyMetadataTransaction) GetAbstractTransaction

```go
func (tx *ModifyMetadataTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*ModifyMetadataTransaction) Size

```go
func (tx *ModifyMetadataTransaction) Size() int
```

#### func (*ModifyMetadataTransaction) String

```go
func (tx *ModifyMetadataTransaction) String() string
```

#### type ModifyMultisigAccountTransaction

```go
type ModifyMultisigAccountTransaction struct {
	AbstractTransaction
	MinApprovalDelta int8
	MinRemovalDelta  int8
	Modifications    []*MultisigCosignatoryModification
}
```


#### func  NewModifyMultisigAccountTransaction

```go
func NewModifyMultisigAccountTransaction(deadline *Deadline, minApprovalDelta int8, minRemovalDelta int8, modifications []*MultisigCosignatoryModification, networkType NetworkType) (*ModifyMultisigAccountTransaction, error)
```
returns a ModifyMultisigAccountTransaction from passed min approval and removal
deltas and array of MultisigCosignatoryModification's

#### func (*ModifyMultisigAccountTransaction) GetAbstractTransaction

```go
func (tx *ModifyMultisigAccountTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*ModifyMultisigAccountTransaction) Size

```go
func (tx *ModifyMultisigAccountTransaction) Size() int
```

#### func (*ModifyMultisigAccountTransaction) String

```go
func (tx *ModifyMultisigAccountTransaction) String() string
```

#### type Mosaic

```go
type Mosaic struct {
	AssetId AssetId
	Amount  Amount
}
```


#### func  NewMosaic

```go
func NewMosaic(assetId AssetId, amount Amount) (*Mosaic, error)
```
returns a Mosaic for passed AssetId and amount

#### func  NewMosaicNoCheck

```go
func NewMosaicNoCheck(assetId AssetId, amount Amount) *Mosaic
```
returns a Mosaic for passed AssetId and amount without validation of parameters

#### func  Xem

```go
func Xem(amount uint64) *Mosaic
```
returns XEM mosaic with passed amount

#### func  XemRelative

```go
func XemRelative(amount uint64) *Mosaic
```
returns XEM with actual passed amount

#### func  Xpx

```go
func Xpx(amount uint64) *Mosaic
```
returns XPX mosaic with passed amount

#### func  XpxRelative

```go
func XpxRelative(amount uint64) *Mosaic
```
returns XPX with actual passed amount

#### func (*Mosaic) String

```go
func (m *Mosaic) String() string
```

#### type MosaicAliasTransaction

```go
type MosaicAliasTransaction struct {
	AliasTransaction
	MosaicId *MosaicId
}
```


#### func  NewMosaicAliasTransaction

```go
func NewMosaicAliasTransaction(deadline *Deadline, mosaicId *MosaicId, namespaceId *NamespaceId, actionType AliasActionType, networkType NetworkType) (*MosaicAliasTransaction, error)
```
returns MosaicAliasTransaction from passed MosaicId, NamespaceId and
AliasActionType

#### func (*MosaicAliasTransaction) Size

```go
func (tx *MosaicAliasTransaction) Size() int
```

#### func (*MosaicAliasTransaction) String

```go
func (tx *MosaicAliasTransaction) String() string
```

#### type MosaicDefinitionTransaction

```go
type MosaicDefinitionTransaction struct {
	AbstractTransaction
	*MosaicProperties
	MosaicNonce uint32
	*MosaicId
}
```


#### func  NewMosaicDefinitionTransaction

```go
func NewMosaicDefinitionTransaction(deadline *Deadline, nonce uint32, ownerPublicKey string, mosaicProps *MosaicProperties, networkType NetworkType) (*MosaicDefinitionTransaction, error)
```
returns MosaicDefinitionTransaction from passed nonce, public key of announcer
and MosaicProperties

#### func (*MosaicDefinitionTransaction) GetAbstractTransaction

```go
func (tx *MosaicDefinitionTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*MosaicDefinitionTransaction) Size

```go
func (tx *MosaicDefinitionTransaction) Size() int
```

#### func (*MosaicDefinitionTransaction) String

```go
func (tx *MosaicDefinitionTransaction) String() string
```

#### type MosaicId

```go
type MosaicId struct {
}
```


#### func  NewMosaicId

```go
func NewMosaicId(id uint64) (*MosaicId, error)
```
returns MosaicId for passed mosaic identifier

#### func  NewMosaicIdFromNonceAndOwner

```go
func NewMosaicIdFromNonceAndOwner(nonce uint32, ownerPublicKey string) (*MosaicId, error)
```
returns MosaicId for passed nonce and public key of mosaic owner

#### func  NewMosaicIdNoCheck

```go
func NewMosaicIdNoCheck(id uint64) *MosaicId
```
TODO

#### func (*MosaicId) Equals

```go
func (m *MosaicId) Equals(id AssetId) bool
```

#### func (*MosaicId) Id

```go
func (m *MosaicId) Id() uint64
```

#### func (*MosaicId) String

```go
func (m *MosaicId) String() string
```

#### func (*MosaicId) Type

```go
func (m *MosaicId) Type() AssetIdType
```

#### type MosaicInfo

```go
type MosaicInfo struct {
	MosaicId   *MosaicId
	Supply     Amount
	Height     Height
	Owner      *PublicAccount
	Revision   uint32
	Properties *MosaicProperties
}
```


#### func (*MosaicInfo) String

```go
func (m *MosaicInfo) String() string
```

#### type MosaicMetadataInfo

```go
type MosaicMetadataInfo struct {
	MetadataInfo
	MosaicId *MosaicId
}
```


#### type MosaicName

```go
type MosaicName struct {
	MosaicId *MosaicId
	Names    []string
}
```


#### func (*MosaicName) String

```go
func (m *MosaicName) String() string
```

#### type MosaicProperties

```go
type MosaicProperties struct {
	SupplyMutable bool
	Transferable  bool
	LevyMutable   bool
	Divisibility  uint8
	Duration      Duration
}
```

structure which includes several properties for defining mosaic `SupplyMutable`
- is supply of defined mosaic can be changed in future `Transferable` - if this
property is set to "false", only transfer transactions having the creator as
sender or as recipient can transfer mosaics of that type. If set to "true" the
mosaics can be transferred to and from arbitrary accounts `LevyMutable` - if
this property is set to "true", whenever other users transact with your mosaic,
owner gets a levy fee from them `Divisibility` - divisibility determines up to
what decimal place the mosaic can be divided into `Duration` - duration in
blocks mosaic will be available. After the renew mosaic is inactive and can be
renewed

#### func  NewMosaicProperties

```go
func NewMosaicProperties(supplyMutable bool, transferable bool, levyMutable bool, divisibility uint8, duration Duration) *MosaicProperties
```
returns MosaicProperties from actual values

#### func (*MosaicProperties) String

```go
func (mp *MosaicProperties) String() string
```

#### type MosaicService

```go
type MosaicService service
```


#### func (*MosaicService) GetMosaicInfo

```go
func (ref *MosaicService) GetMosaicInfo(ctx context.Context, mosaicId *MosaicId) (*MosaicInfo, error)
```

#### func (*MosaicService) GetMosaicInfos

```go
func (ref *MosaicService) GetMosaicInfos(ctx context.Context, mscIds []*MosaicId) ([]*MosaicInfo, error)
```

#### func (*MosaicService) GetMosaicsNames

```go
func (ref *MosaicService) GetMosaicsNames(ctx context.Context, mscIds ...*MosaicId) ([]*MosaicName, error)
```
GetMosaicsNames Get readable names for a set of mosaics post @/mosaic/names

#### type MosaicSupplyChangeTransaction

```go
type MosaicSupplyChangeTransaction struct {
	AbstractTransaction
	MosaicSupplyType
	AssetId
	Delta Amount
}
```


#### func  NewMosaicSupplyChangeTransaction

```go
func NewMosaicSupplyChangeTransaction(deadline *Deadline, assetId AssetId, supplyType MosaicSupplyType, delta Duration, networkType NetworkType) (*MosaicSupplyChangeTransaction, error)
```
returns MosaicSupplyChangeTransaction from passed AssetId, MosaicSupplyTypeand
supply delta

#### func (*MosaicSupplyChangeTransaction) GetAbstractTransaction

```go
func (tx *MosaicSupplyChangeTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*MosaicSupplyChangeTransaction) Size

```go
func (tx *MosaicSupplyChangeTransaction) Size() int
```

#### func (*MosaicSupplyChangeTransaction) String

```go
func (tx *MosaicSupplyChangeTransaction) String() string
```

#### type MosaicSupplyType

```go
type MosaicSupplyType uint8
```


```go
const (
	Decrease MosaicSupplyType = iota
	Increase
)
```

#### func (MosaicSupplyType) String

```go
func (tx MosaicSupplyType) String() string
```

#### type MultisigAccountGraphInfo

```go
type MultisigAccountGraphInfo struct {
	MultisigAccounts map[int32][]*MultisigAccountInfo
}
```


#### type MultisigAccountInfo

```go
type MultisigAccountInfo struct {
	Account          PublicAccount
	MinApproval      int32
	MinRemoval       int32
	Cosignatories    []*PublicAccount
	MultisigAccounts []*PublicAccount
}
```


#### func (*MultisigAccountInfo) String

```go
func (ref *MultisigAccountInfo) String() string
```

#### type MultisigCosignatoryModification

```go
type MultisigCosignatoryModification struct {
	Type MultisigCosignatoryModificationType
	*PublicAccount
}
```


#### func (*MultisigCosignatoryModification) String

```go
func (m *MultisigCosignatoryModification) String() string
```

#### type MultisigCosignatoryModificationType

```go
type MultisigCosignatoryModificationType uint8
```


```go
const (
	Add MultisigCosignatoryModificationType = iota
	Remove
)
```

#### func (MultisigCosignatoryModificationType) String

```go
func (t MultisigCosignatoryModificationType) String() string
```

#### type NamespaceAlias

```go
type NamespaceAlias struct {
	Type AliasType
}
```

NamespaceAlias contains aliased mosaicId or address and type of alias

#### func (*NamespaceAlias) Address

```go
func (ref *NamespaceAlias) Address() *Address
```

#### func (*NamespaceAlias) MosaicId

```go
func (ref *NamespaceAlias) MosaicId() *MosaicId
```

#### func (*NamespaceAlias) String

```go
func (ref *NamespaceAlias) String() string
```

#### type NamespaceId

```go
type NamespaceId struct {
}
```


#### func  GenerateNamespacePath

```go
func GenerateNamespacePath(name string) ([]*NamespaceId, error)
```
returns an array of big ints representation if namespace ids from passed
namespace path to create root namespace pass namespace name in format like
'rootname' to create child namespace pass namespace name in format like
'rootname.childname' to create grand child namespace pass namespace name in
format like 'rootname.childname.grandchildname'

#### func  NewNamespaceId

```go
func NewNamespaceId(id uint64) (*NamespaceId, error)
```
returns new NamespaceId from passed namespace identifier

#### func  NewNamespaceIdFromName

```go
func NewNamespaceIdFromName(namespaceName string) (*NamespaceId, error)
```
returns namespace id from passed namespace name should be used for creating
root, child and grandchild namespace ids to create root namespace pass namespace
name in format like 'rootname' to create child namespace pass namespace name in
format like 'rootname.childname' to create grand child namespace pass namespace
name in format like 'rootname.childname.grandchildname'

#### func  NewNamespaceIdNoCheck

```go
func NewNamespaceIdNoCheck(id uint64) *NamespaceId
```
returns new NamespaceId from passed namespace identifier TODO

#### func (*NamespaceId) Equals

```go
func (m *NamespaceId) Equals(id AssetId) bool
```

#### func (*NamespaceId) Id

```go
func (m *NamespaceId) Id() uint64
```

#### func (*NamespaceId) String

```go
func (m *NamespaceId) String() string
```

#### func (*NamespaceId) Type

```go
func (m *NamespaceId) Type() AssetIdType
```

#### type NamespaceInfo

```go
type NamespaceInfo struct {
	NamespaceId *NamespaceId
	Active      bool
	TypeSpace   NamespaceType
	Depth       int
	Levels      []*NamespaceId
	Alias       *NamespaceAlias
	Parent      *NamespaceInfo
	Owner       *PublicAccount
	StartHeight Height
	EndHeight   Height
}
```


#### func (*NamespaceInfo) String

```go
func (ref *NamespaceInfo) String() string
```

#### type NamespaceMetadataInfo

```go
type NamespaceMetadataInfo struct {
	MetadataInfo
	NamespaceId *NamespaceId
}
```


#### type NamespaceName

```go
type NamespaceName struct {
	NamespaceId *NamespaceId
	Name        string
	ParentId    *NamespaceId /* Optional NamespaceId my be nil */
}
```


#### func (*NamespaceName) String

```go
func (n *NamespaceName) String() string
```

#### type NamespaceService

```go
type NamespaceService service
```

NamespaceService provides a set of methods for obtaining information about the
namespace

#### func (*NamespaceService) GetLinkedAddress

```go
func (ref *NamespaceService) GetLinkedAddress(ctx context.Context, namespaceId *NamespaceId) (*Address, error)
```
GetLinkedAddress @/namespace/%s

#### func (*NamespaceService) GetLinkedMosaicId

```go
func (ref *NamespaceService) GetLinkedMosaicId(ctx context.Context, namespaceId *NamespaceId) (*MosaicId, error)
```
GetLinkedMosaicId @/namespace/%s

#### func (*NamespaceService) GetNamespaceInfo

```go
func (ref *NamespaceService) GetNamespaceInfo(ctx context.Context, nsId *NamespaceId) (*NamespaceInfo, error)
```

#### func (*NamespaceService) GetNamespaceInfosFromAccount

```go
func (ref *NamespaceService) GetNamespaceInfosFromAccount(ctx context.Context, address *Address, nsId *NamespaceId,
	pageSize int) ([]*NamespaceInfo, error)
```
returns NamespaceInfo's corresponding to passed Address and NamespaceId with
maximum limit TODO: fix pagination

#### func (*NamespaceService) GetNamespaceInfosFromAccounts

```go
func (ref *NamespaceService) GetNamespaceInfosFromAccounts(ctx context.Context, addrs []*Address, nsId *NamespaceId,
	pageSize int) ([]*NamespaceInfo, error)
```
returns NamespaceInfo's corresponding to passed Address's and NamespaceId with
maximum limit TODO: fix pagination

#### func (*NamespaceService) GetNamespaceNames

```go
func (ref *NamespaceService) GetNamespaceNames(ctx context.Context, nsIds []*NamespaceId) ([]*NamespaceName, error)
```

#### type NamespaceType

```go
type NamespaceType uint8
```


```go
const (
	Root NamespaceType = iota
	Sub
)
```

#### type NetworkService

```go
type NetworkService service
```


#### func (*NetworkService) GetNetworkType

```go
func (ref *NetworkService) GetNetworkType(ctx context.Context) (NetworkType, error)
```

#### type NetworkType

```go
type NetworkType uint8
```


```go
const (
	Mijin           NetworkType = 96
	MijinTest       NetworkType = 144
	Public          NetworkType = 184
	PublicTest      NetworkType = 168
	Private         NetworkType = 200
	PrivateTest     NetworkType = 176
	NotSupportedNet NetworkType = 0
	AliasAddress    NetworkType = 145
)
```

#### func  ExtractNetworkType

```go
func ExtractNetworkType(version uint64) NetworkType
```

#### func  NetworkTypeFromString

```go
func NetworkTypeFromString(networkType string) NetworkType
```

#### func (NetworkType) String

```go
func (nt NetworkType) String() string
```

#### type PartialAddedMapper

```go
type PartialAddedMapper interface {
	MapPartialAdded(m []byte) (*AggregateTransaction, error)
}
```


#### func  NewPartialAddedMapper

```go
func NewPartialAddedMapper(mapTransactionFunc mapTransactionFunc) PartialAddedMapper
```

#### type PartialRemovedInfo

```go
type PartialRemovedInfo struct {
	Meta *TransactionInfo
}
```


#### func  MapPartialRemoved

```go
func MapPartialRemoved(m []byte) (*PartialRemovedInfo, error)
```

#### type PartialRemovedMapper

```go
type PartialRemovedMapper interface {
	MapPartialRemoved(m []byte) (*PartialRemovedInfo, error)
}
```


#### type PartialRemovedMapperFn

```go
type PartialRemovedMapperFn func(m []byte) (*PartialRemovedInfo, error)
```


#### func (PartialRemovedMapperFn) MapPartialRemoved

```go
func (p PartialRemovedMapperFn) MapPartialRemoved(m []byte) (*PartialRemovedInfo, error)
```

#### type PlainMessage

```go
type PlainMessage struct {
}
```


#### func  NewPlainMessage

```go
func NewPlainMessage(payload string) *PlainMessage
```

#### func  NewPlainMessageFromEncodedData

```go
func NewPlainMessageFromEncodedData(encodedData []byte, recipient *xpxcrypto.PrivateKey, sender *xpxcrypto.PublicKey) (*PlainMessage, error)
```

#### func (*PlainMessage) Message

```go
func (m *PlainMessage) Message() string
```

#### func (*PlainMessage) Payload

```go
func (m *PlainMessage) Payload() []byte
```

#### func (*PlainMessage) String

```go
func (m *PlainMessage) String() string
```

#### func (*PlainMessage) Type

```go
func (m *PlainMessage) Type() MessageType
```

#### type Proof

```go
type Proof struct {
	Data []byte
}
```


#### func  NewProofFromBytes

```go
func NewProofFromBytes(proof []byte) *Proof
```

#### func  NewProofFromHexString

```go
func NewProofFromHexString(hexProof string) (*Proof, error)
```

#### func  NewProofFromString

```go
func NewProofFromString(proof string) *Proof
```

#### func  NewProofFromUint16

```go
func NewProofFromUint16(number uint16) *Proof
```

#### func  NewProofFromUint32

```go
func NewProofFromUint32(number uint32) *Proof
```

#### func  NewProofFromUint64

```go
func NewProofFromUint64(number uint64) *Proof
```

#### func  NewProofFromUint8

```go
func NewProofFromUint8(number uint8) *Proof
```

#### func (*Proof) ProofString

```go
func (p *Proof) ProofString() string
```
bytes representation of Proof

#### func (*Proof) Secret

```go
func (p *Proof) Secret(hashType HashType) (*Secret, error)
```
returns Secret generated from Proof with passed HashType

#### func (*Proof) Size

```go
func (p *Proof) Size() int
```
bytes length of Proof

#### func (*Proof) String

```go
func (p *Proof) String() string
```

#### type PropertyModificationType

```go
type PropertyModificationType uint8
```


```go
const (
	AddProperty PropertyModificationType = iota
	RemoveProperty
)
```
PropertyModificationType enums

#### type PropertyType

```go
type PropertyType uint8
```


```go
const (
	AllowAddress     PropertyType = 0x01
	AllowMosaic      PropertyType = 0x02
	AllowTransaction PropertyType = 0x04
	Sentinel         PropertyType = 0x05
	BlockAddress     PropertyType = 0x80 + 0x01
	BlockMosaic      PropertyType = 0x80 + 0x02
	BlockTransaction PropertyType = 0x80 + 0x04
)
```
Account property type 0x01 The property type is an address. 0x02 The property
type is mosaic id. 0x04 The property type is a transaction type. 0x05 Property
type sentinel. 0x80 + type The property is interpreted as a blocking operation.

#### type PublicAccount

```go
type PublicAccount struct {
	Address   *Address
	PublicKey string
}
```


#### func  NewAccountFromPublicKey

```go
func NewAccountFromPublicKey(pKey string, networkType NetworkType) (*PublicAccount, error)
```
returns a PublicAccount from public key for passed NetworkType

#### func (*PublicAccount) String

```go
func (ref *PublicAccount) String() string
```

#### type RegisterNamespaceTransaction

```go
type RegisterNamespaceTransaction struct {
	AbstractTransaction
	*NamespaceId
	NamespaceType
	NamspaceName string
	Duration     Duration
	ParentId     *NamespaceId
}
```


#### func  NewRegisterRootNamespaceTransaction

```go
func NewRegisterRootNamespaceTransaction(deadline *Deadline, namespaceName string, duration Duration, networkType NetworkType) (*RegisterNamespaceTransaction, error)
```
returns a RegisterNamespaceTransaction from passed namespace name and duration
in blocks

#### func  NewRegisterSubNamespaceTransaction

```go
func NewRegisterSubNamespaceTransaction(deadline *Deadline, namespaceName string, parentId *NamespaceId, networkType NetworkType) (*RegisterNamespaceTransaction, error)
```
returns a RegisterNamespaceTransaction from passed namespace name and parent
NamespaceId

#### func (*RegisterNamespaceTransaction) GetAbstractTransaction

```go
func (tx *RegisterNamespaceTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*RegisterNamespaceTransaction) Size

```go
func (tx *RegisterNamespaceTransaction) Size() int
```

#### func (*RegisterNamespaceTransaction) String

```go
func (tx *RegisterNamespaceTransaction) String() string
```

#### type ResolverService

```go
type ResolverService struct {
	NamespaceService *NamespaceService
	MosaicService    *MosaicService
}
```

TODO: Implement resolving namespace to account

#### func (*ResolverService) GetMosaicInfoByAssetId

```go
func (ref *ResolverService) GetMosaicInfoByAssetId(ctx context.Context, assetId AssetId) (*MosaicInfo, error)
```

#### func (*ResolverService) GetMosaicInfosByAssetIds

```go
func (ref *ResolverService) GetMosaicInfosByAssetIds(ctx context.Context, assetIds ...AssetId) ([]*MosaicInfo, error)
```

#### type RespErr

```go
type RespErr struct {
}
```


#### func (*RespErr) Error

```go
func (r *RespErr) Error() string
```

#### type Secret

```go
type Secret struct {
	Hash []byte
	Type HashType
}
```


#### func  NewSecret

```go
func NewSecret(hash []byte, hashType HashType) (*Secret, error)
```
returns Secret from passed hash and HashType

#### func  NewSecretFromHexString

```go
func NewSecretFromHexString(hash string, hashType HashType) (*Secret, error)
```
returns Secret from passed hex string hash and HashType

#### func (*Secret) HashString

```go
func (s *Secret) HashString() string
```

#### func (*Secret) String

```go
func (s *Secret) String() string
```

#### type SecretLockTransaction

```go
type SecretLockTransaction struct {
	AbstractTransaction
	*Mosaic
	Duration  Duration
	Secret    *Secret
	Recipient *Address
}
```


#### func  NewSecretLockTransaction

```go
func NewSecretLockTransaction(deadline *Deadline, mosaic *Mosaic, duration Duration, secret *Secret, recipient *Address, networkType NetworkType) (*SecretLockTransaction, error)
```
returns a SecretLockTransaction from passed Mosaic, duration in blocks, Secret
and mosaic recipient Address

#### func (*SecretLockTransaction) GetAbstractTransaction

```go
func (tx *SecretLockTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*SecretLockTransaction) Size

```go
func (tx *SecretLockTransaction) Size() int
```

#### func (*SecretLockTransaction) String

```go
func (tx *SecretLockTransaction) String() string
```

#### type SecretProofTransaction

```go
type SecretProofTransaction struct {
	AbstractTransaction
	HashType
	Proof *Proof
}
```


#### func  NewSecretProofTransaction

```go
func NewSecretProofTransaction(deadline *Deadline, hashType HashType, proof *Proof, networkType NetworkType) (*SecretProofTransaction, error)
```
returns a SecretProofTransaction from passed HashType and Proof

#### func (*SecretProofTransaction) GetAbstractTransaction

```go
func (tx *SecretProofTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*SecretProofTransaction) Size

```go
func (tx *SecretProofTransaction) Size() int
```

#### func (*SecretProofTransaction) String

```go
func (tx *SecretProofTransaction) String() string
```

#### type SecureMessage

```go
type SecureMessage struct {
}
```


#### func  NewSecureMessage

```go
func NewSecureMessage(encodedData []byte) *SecureMessage
```

#### func  NewSecureMessageFromPlaintText

```go
func NewSecureMessageFromPlaintText(plaintText string, sender *xpxcrypto.PrivateKey, recipient *xpxcrypto.PublicKey) (*SecureMessage, error)
```

#### func (*SecureMessage) Payload

```go
func (m *SecureMessage) Payload() []byte
```

#### func (*SecureMessage) String

```go
func (m *SecureMessage) String() string
```

#### func (*SecureMessage) Type

```go
func (m *SecureMessage) Type() MessageType
```

#### type SignedTransaction

```go
type SignedTransaction struct {
	TransactionType `json:"transactionType"`
	Payload         string `json:"payload"`
	Hash            Hash   `json:"hash"`
}
```


#### type SignerInfo

```go
type SignerInfo struct {
	Signer     string `json:"signer"`
	Signature  string `json:"signature"`
	ParentHash Hash   `json:"parentHash"`
}
```


#### func  MapCosignature

```go
func MapCosignature(m []byte) (*SignerInfo, error)
```

#### type StatusInfo

```go
type StatusInfo struct {
	Status string `json:"status"`
	Hash   Hash   `json:"hash"`
}
```


#### func  MapStatus

```go
func MapStatus(m []byte) (*StatusInfo, error)
```

#### type StatusMapper

```go
type StatusMapper interface {
	MapStatus(m []byte) (*StatusInfo, error)
}
```


#### type StatusMapperFn

```go
type StatusMapperFn func(m []byte) (*StatusInfo, error)
```


#### func (StatusMapperFn) MapStatus

```go
func (p StatusMapperFn) MapStatus(m []byte) (*StatusInfo, error)
```

#### type Timestamp

```go
type Timestamp struct {
	time.Time
}
```


#### func  NewTimestamp

```go
func NewTimestamp(milliseconds int64) *Timestamp
```
returns new Timestamp from passed milliseconds value

#### func (*Timestamp) ToBlockchainTimestamp

```go
func (t *Timestamp) ToBlockchainTimestamp() *BlockchainTimestamp
```

#### type Transaction

```go
type Transaction interface {
	GetAbstractTransaction() *AbstractTransaction
	String() string
	// number of bytes of serialized transaction
	Size() int
	// contains filtered or unexported methods
}
```


#### func  MapTransaction

```go
func MapTransaction(b *bytes.Buffer) (Transaction, error)
```

#### func  MapTransactions

```go
func MapTransactions(b *bytes.Buffer) ([]Transaction, error)
```

#### type TransactionHashesDTO

```go
type TransactionHashesDTO struct {
	Hashes []string `json:"hashes"`
}
```


#### type TransactionIdsDTO

```go
type TransactionIdsDTO struct {
	Ids []string `json:"transactionIds"`
}
```


#### type TransactionInfo

```go
type TransactionInfo struct {
	Height              Height
	Index               uint32
	Id                  string
	Hash                Hash
	MerkleComponentHash Hash
	AggregateHash       Hash
	AggregateId         string
}
```


#### func (*TransactionInfo) String

```go
func (ti *TransactionInfo) String() string
```

#### type TransactionOrder

```go
type TransactionOrder string
```


```go
const (
	TRANSACTION_ORDER_ASC  TransactionOrder = "id"
	TRANSACTION_ORDER_DESC TransactionOrder = "-id"
)
```

#### type TransactionService

```go
type TransactionService struct {
	BlockchainService *BlockchainService
}
```


#### func (*TransactionService) Announce

```go
func (txs *TransactionService) Announce(ctx context.Context, tx *SignedTransaction) (string, error)
```
returns transaction hash after announcing passed SignedTransaction

#### func (*TransactionService) AnnounceAggregateBonded

```go
func (txs *TransactionService) AnnounceAggregateBonded(ctx context.Context, tx *SignedTransaction) (string, error)
```
returns transaction hash after announcing passed aggregate bounded
SignedTransaction

#### func (*TransactionService) AnnounceAggregateBondedCosignature

```go
func (txs *TransactionService) AnnounceAggregateBondedCosignature(ctx context.Context, c *CosignatureSignedTransaction) (string, error)
```
returns transaction hash after announcing passed CosignatureSignedTransaction

#### func (*TransactionService) GetTransaction

```go
func (txs *TransactionService) GetTransaction(ctx context.Context, id string) (Transaction, error)
```
returns Transaction for passed transaction id or hash

#### func (*TransactionService) GetTransactionEffectiveFee

```go
func (txs *TransactionService) GetTransactionEffectiveFee(ctx context.Context, transactionId string) (int, error)
```
Gets a transaction's effective paid fee

#### func (*TransactionService) GetTransactionStatus

```go
func (txs *TransactionService) GetTransactionStatus(ctx context.Context, id string) (*TransactionStatus, error)
```
returns TransactionStatus for passed transaction id or hash

#### func (*TransactionService) GetTransactionStatuses

```go
func (txs *TransactionService) GetTransactionStatuses(ctx context.Context, hashes []string) ([]*TransactionStatus, error)
```
returns an array of TransactionStatus's for passed transaction ids or hashes

#### func (*TransactionService) GetTransactions

```go
func (txs *TransactionService) GetTransactions(ctx context.Context, ids []string) ([]Transaction, error)
```
returns an array of Transaction's for passed array of transaction ids or hashes

#### type TransactionStatus

```go
type TransactionStatus struct {
	Deadline *Deadline
	Group    string
	Status   string
	Hash     Hash
	Height   Height
}
```


#### func (*TransactionStatus) String

```go
func (ts *TransactionStatus) String() string
```

#### type TransactionType

```go
type TransactionType uint16
```


```go
const (
	AccountPropertyAddress    TransactionType = 0x4150
	AccountPropertyMosaic     TransactionType = 0x4250
	AccountPropertyEntityType TransactionType = 0x4350
	AddressAlias              TransactionType = 0x424e
	AggregateBonded           TransactionType = 0x4241
	AggregateCompleted        TransactionType = 0x4141
	LinkAccount               TransactionType = 0x414c
	Lock                      TransactionType = 0x4148
	MetadataAddress           TransactionType = 0x413d
	MetadataMosaic            TransactionType = 0x423d
	MetadataNamespace         TransactionType = 0x433d
	ModifyContract            TransactionType = 0x4157
	ModifyMultisig            TransactionType = 0x4155
	MosaicAlias               TransactionType = 0x434e
	MosaicDefinition          TransactionType = 0x414d
	MosaicSupplyChange        TransactionType = 0x424d
	RegisterNamespace         TransactionType = 0x414e
	SecretLock                TransactionType = 0x4152
	SecretProof               TransactionType = 0x4252
	Transfer                  TransactionType = 0x4154
)
```

#### func (TransactionType) String

```go
func (t TransactionType) String() string
```

#### type TransactionVersion

```go
type TransactionVersion uint8
```


```go
const (
	AccountPropertyAddressVersion    TransactionVersion = 1
	AccountPropertyMosaicVersion     TransactionVersion = 1
	AccountPropertyEntityTypeVersion TransactionVersion = 1
	AddressAliasVersion              TransactionVersion = 1
	AggregateBondedVersion           TransactionVersion = 2
	AggregateCompletedVersion        TransactionVersion = 2
	LinkAccountVersion               TransactionVersion = 2
	LockVersion                      TransactionVersion = 1
	MetadataAddressVersion           TransactionVersion = 1
	MetadataMosaicVersion            TransactionVersion = 1
	MetadataNamespaceVersion         TransactionVersion = 1
	ModifyContractVersion            TransactionVersion = 3
	ModifyMultisigVersion            TransactionVersion = 3
	MosaicAliasVersion               TransactionVersion = 1
	MosaicDefinitionVersion          TransactionVersion = 3
	MosaicSupplyChangeVersion        TransactionVersion = 2
	RegisterNamespaceVersion         TransactionVersion = 2
	SecretLockVersion                TransactionVersion = 1
	SecretProofVersion               TransactionVersion = 1
	TransferVersion                  TransactionVersion = 3
)
```

#### type TransferTransaction

```go
type TransferTransaction struct {
	AbstractTransaction
	Message   Message
	Mosaics   []*Mosaic
	Recipient *Address
}
```


#### func  NewTransferTransaction

```go
func NewTransferTransaction(deadline *Deadline, recipient *Address, mosaics []*Mosaic, message Message, networkType NetworkType) (*TransferTransaction, error)
```
returns a TransferTransaction from passed transfer recipient Adderess, array of
Mosaic's to transfer and transfer Message

#### func  NewTransferTransactionWithNamespace

```go
func NewTransferTransactionWithNamespace(deadline *Deadline, recipient *NamespaceId, mosaics []*Mosaic, message Message, networkType NetworkType) (*TransferTransaction, error)
```
returns TransferTransaction from passed recipient NamespaceId, Mosaic's and
transfer Message

#### func (*TransferTransaction) GetAbstractTransaction

```go
func (tx *TransferTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*TransferTransaction) MessageSize

```go
func (tx *TransferTransaction) MessageSize() int
```

#### func (*TransferTransaction) Size

```go
func (tx *TransferTransaction) Size() int
```

#### func (*TransferTransaction) String

```go
func (tx *TransferTransaction) String() string
```

#### type UnconfirmedAddedMapper

```go
type UnconfirmedAddedMapper interface {
	MapUnconfirmedAdded(m []byte) (Transaction, error)
}
```


#### func  NewUnconfirmedAddedMapper

```go
func NewUnconfirmedAddedMapper(mapTransactionFunc mapTransactionFunc) UnconfirmedAddedMapper
```

#### type UnconfirmedRemoved

```go
type UnconfirmedRemoved struct {
	Meta *TransactionInfo
}
```


#### func  MapUnconfirmedRemoved

```go
func MapUnconfirmedRemoved(m []byte) (*UnconfirmedRemoved, error)
```

#### type UnconfirmedRemovedMapper

```go
type UnconfirmedRemovedMapper interface {
	MapUnconfirmedRemoved(m []byte) (*UnconfirmedRemoved, error)
}
```


#### type UnconfirmedRemovedMapperFn

```go
type UnconfirmedRemovedMapperFn func(m []byte) (*UnconfirmedRemoved, error)
```


#### func (UnconfirmedRemovedMapperFn) MapUnconfirmedRemoved

```go
func (p UnconfirmedRemovedMapperFn) MapUnconfirmedRemoved(m []byte) (*UnconfirmedRemoved, error)
```

#### type VarSize

```go
type VarSize uint32
```


```go
const (
	ByteSize  VarSize = 1
	ShortSize VarSize = 2
	IntSize   VarSize = 4
)
```

#### type WsMessageInfo

```go
type WsMessageInfo struct {
	Address     *Address
	ChannelName string
}
```


#### type WsMessageInfoDTO

```go
type WsMessageInfoDTO struct {
	Meta wsMessageInfoMetaDTO `json:"meta"`
}
```


#### func (*WsMessageInfoDTO) ToStruct

```go
func (dto *WsMessageInfoDTO) ToStruct() (*WsMessageInfo, error)
```
