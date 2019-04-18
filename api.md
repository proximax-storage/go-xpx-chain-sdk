## Usage

```go
const NUM_CHECKSUM_BYTES = 4
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
	ErrEmptyMosaicIds        = errors.New("list mosaics ids must not by empty")
	ErrNilMosaicId           = errors.New("mosaicId must not be nil")
	ErrNilMosaicAmount       = errors.New("amount must be not nil")
	ErrInvalidOwnerPublicKey = errors.New("public owner key is invalid")
	ErrNilMosaicProperties   = errors.New("mosaic properties must not be nil")
)
```
Mosaic errors

```go
var (
	ErrNamespaceTooManyPart = errors.New("too many parts")
	ErrNilNamespaceId       = errors.New("namespaceId is nil or zero")
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
var TimestampNemesisBlock = time.Unix(1459468800, 0)
```

```go
var XemMosaicId, _ = NewMosaicId(big.NewInt(0x0DC67FBE1CAD29E3))
```
mosaic id for XEM mosaic

```go
var XpxMosaicId, _ = NewMosaicId(big.NewInt(0x0DC67FBE1CAD29E3))
```
mosaic id for XPX mosaic

#### func  BigIntegerToHex

```go
func BigIntegerToHex(id *big.Int) string
```
TODO analog JAVA Uint64.bigIntegerToHex

#### func  ExtractVersion

```go
func ExtractVersion(version uint64) uint8
```

#### func  FromBigInt

```go
func FromBigInt(int *big.Int) []uint32
```
TODO

#### func  GenerateChecksum

```go
func GenerateChecksum(b []byte) ([]byte, error)
```

#### func  GenerateNamespacePath

```go
func GenerateNamespacePath(name string) ([]*big.Int, error)
```
GenerateNamespacePath create list NamespaceId from string returns an array of
big ints representation if namespace ids from passed namespace path to create
root namespace pass namespace name in format like `rootname` to create child
namespace pass namespace name in format like `rootna.childname` to create grand
child namespace pass namespace name in format like
`rootna.childname.grandchildname`

#### func  IntToHex

```go
func IntToHex(u uint32) string
```
TODO

#### func  NewReputationConfig

```go
func NewReputationConfig(minInter uint64, defaultRep float64) (*reputationConfig, error)
```
TODO

#### type AbstractTransaction

```go
type AbstractTransaction struct {
	*TransactionInfo
	NetworkType NetworkType
	Deadline    *Deadline
	Type        TransactionType
	Version     TransactionVersion
	Fee         *big.Int
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
returns new account generated for passed network type

#### func  NewAccountFromPrivateKey

```go
func NewAccountFromPrivateKey(pKey string, networkType NetworkType) (*Account, error)
```
returns new account from private key for passed network type

#### func (*Account) Sign

```go
func (a *Account) Sign(tx Transaction) (*SignedTransaction, error)
```
signs given transaction returns a signed transaction

#### func (*Account) SignCosignatureTransaction

```go
func (a *Account) SignCosignatureTransaction(tx *CosignatureTransaction) (*CosignatureSignedTransaction, error)
```
signs aggregate signature transaction returns signed cosignature transaction

#### func (*Account) SignWithCosignatures

```go
func (a *Account) SignWithCosignatures(tx *AggregateTransaction, cosignatories []*Account) (*SignedTransaction, error)
```
TODO

#### type AccountInfo

```go
type AccountInfo struct {
	Address          *Address
	AddressHeight    *big.Int
	PublicKey        string
	PublicKeyHeight  *big.Int
	Importance       *big.Int
	ImportanceHeight *big.Int
	Mosaics          []*Mosaic
	Reputation       float64
}
```


#### func (*AccountInfo) String

```go
func (a *AccountInfo) String() string
```

#### type AccountService

```go
type AccountService service
```


#### func (*AccountService) AggregateBondedTransactions

```go
func (a *AccountService) AggregateBondedTransactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]*AggregateTransaction, error)
```
returns an array of aggregate bounded transactions where an account is signer or
cosigner

#### func (*AccountService) GetAccountInfo

```go
func (a *AccountService) GetAccountInfo(ctx context.Context, address *Address) (*AccountInfo, error)
```
returns account info for given address

#### func (*AccountService) GetAccountsInfo

```go
func (a *AccountService) GetAccountsInfo(ctx context.Context, addresses []*Address) ([]*AccountInfo, error)
```
returns an array of account info for given addresses

#### func (*AccountService) GetMultisigAccountGraphInfo

```go
func (a *AccountService) GetMultisigAccountGraphInfo(ctx context.Context, address *Address) (*MultisigAccountGraphInfo, error)
```
returns multisig account info for given address

#### func (*AccountService) GetMultisigAccountInfo

```go
func (a *AccountService) GetMultisigAccountInfo(ctx context.Context, address *Address) (*MultisigAccountInfo, error)
```
returns multisig account info for given address

#### func (*AccountService) IncomingTransactions

```go
func (a *AccountService) IncomingTransactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]Transaction, error)
```
returns an array of transactions for which an account is receiver

#### func (*AccountService) OutgoingTransactions

```go
func (a *AccountService) OutgoingTransactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]Transaction, error)
```
returns an array of transaction for which an account is sender

#### func (*AccountService) Transactions

```go
func (a *AccountService) Transactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]Transaction, error)
```
returns an array of confirmed transactions for which an account is sender or
receiver.

#### func (*AccountService) UnconfirmedTransactions

```go
func (a *AccountService) UnconfirmedTransactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]Transaction, error)
```
returns an array of confirmed transactions for which an account is sender or
receiver. unconfirmed transactions are those transactions that have not yet been
included in a block. unconfirmed transactions are not guaranteed to be included
in any block.

#### type AccountTransactionsOption

```go
type AccountTransactionsOption struct {
	PageSize int    `url:"pageSize,omitempty"`
	Id       string `url:"id,omitempty"`
}
```


#### type Address

```go
type Address struct {
	Type    NetworkType
	Address string
}
```


#### func  NewAddress

```go
func NewAddress(address string, networkType NetworkType) *Address
```
returns address struc from passed address string for passed network type

#### func  NewAddressFromEncoded

```go
func NewAddressFromEncoded(encoded string) (*Address, error)
```
TODO

#### func  NewAddressFromPublicKey

```go
func NewAddressFromPublicKey(pKey string, networkType NetworkType) (*Address, error)
```
returns an address from public key for passed network type

#### func  NewAddressFromRaw

```go
func NewAddressFromRaw(address string) (*Address, error)
```
TODO

#### func (*Address) Pretty

```go
func (ad *Address) Pretty() string
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
returns bounded aggregate transaction or error from passed deadline and an array
of transactions to be included in

#### func  NewCompleteAggregateTransaction

```go
func NewCompleteAggregateTransaction(deadline *Deadline, innerTxs []Transaction, networkType NetworkType) (*AggregateTransaction, error)
```
returns complete aggregate transaction or error from passed deadline and an
array of transactions to be included in

#### func (*AggregateTransaction) GetAbstractTransaction

```go
func (tx *AggregateTransaction) GetAbstractTransaction() *AbstractTransaction
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

#### type BlockInfo

```go
type BlockInfo struct {
	NetworkType
	Hash                  string
	GenerationHash        string
	TotalFee              *big.Int
	NumTransactions       uint64
	Signature             string
	Signer                *PublicAccount
	Version               uint8
	Type                  uint64
	Height                *big.Int
	Timestamp             *big.Int
	Difficulty            *big.Int
	PreviousBlockHash     string
	BlockTransactionsHash string
}
```

Models Block

#### func (*BlockInfo) String

```go
func (b *BlockInfo) String() string
```

#### type BlockchainService

```go
type BlockchainService service
```


#### func (*BlockchainService) GetBlockByHeight

```go
func (b *BlockchainService) GetBlockByHeight(ctx context.Context, height *big.Int) (*BlockInfo, error)
```
returns info for block with passed height

#### func (*BlockchainService) GetBlockTransactions

```go
func (b *BlockchainService) GetBlockTransactions(ctx context.Context, height *big.Int) ([]Transaction, error)
```
get transactions inside of block with passed height

#### func (*BlockchainService) GetBlockchainHeight

```go
func (b *BlockchainService) GetBlockchainHeight(ctx context.Context) (*big.Int, error)
```
returns blockchain height

#### func (*BlockchainService) GetBlockchainScore

```go
func (b *BlockchainService) GetBlockchainScore(ctx context.Context) (*big.Int, error)
```
returns blockchain score

#### func (*BlockchainService) GetBlockchainStorage

```go
func (b *BlockchainService) GetBlockchainStorage(ctx context.Context) (*BlockchainStorageInfo, error)
```
returns storage information

#### func (*BlockchainService) GetBlocksByHeightWithLimit

```go
func (b *BlockchainService) GetBlocksByHeightWithLimit(ctx context.Context, height, limit *big.Int) ([]*BlockInfo, error)
```
TODO

#### type BlockchainStorageInfo

```go
type BlockchainStorageInfo struct {
	NumBlocks       int `json:"numBlocks"`
	NumTransactions int `json:"numTransactions"`
	NumAccounts     int `json:"numAccounts"`
}
```

Blockchain Storage

#### func (*BlockchainStorageInfo) String

```go
func (b *BlockchainStorageInfo) String() string
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
	Account     *AccountService
	Contract    *ContractService
}
```

Catapult API Client configuration

#### func  NewClient

```go
func NewClient(httpClient *http.Client, conf *Config) *Client
```
returns catapult http client from passed existing client and configuration if
passed client is nil, http.DefaultClient will be used

#### func (*Client) Do

```go
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error)
```
TODO Do sends an API Request and returns a parsed response

#### func (*Client) DoNewRequest

```go
func (s *Client) DoNewRequest(ctx context.Context, method string, path string, body interface{}, v interface{}) (*http.Response, error)
```
TODO DoNewRequest creates new request, Do it & return result in V

#### func (*Client) NewRequest

```go
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error)
```
TODO

#### type ClientWebsocket

```go
type ClientWebsocket struct {
	Uid string

	Subscribe *SubscribeService
}
```

Catapult Websocket Client configuration

#### func  NewConnectWs

```go
func NewConnectWs(host string, timeout time.Duration) (*ClientWebsocket, error)
```
returns entity which you can use to reach different subscribe services from
passed host url and waiting timeout

#### type Config

```go
type Config struct {
	BaseURL *url.URL
	NetworkType
}
```

Provides service configuration

#### func  NewConfig

```go
func NewConfig(baseUrl string, networkType NetworkType) (*Config, error)
```
returns config for HTTP Client from passed base url and network type

#### func  NewConfigWithReputation

```go
func NewConfigWithReputation(baseUrl string, networkType NetworkType, repConf *reputationConfig) (*Config, error)
```
TODO

#### type ContractInfo

```go
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
```


#### type ContractService

```go
type ContractService service
```


#### func (*ContractService) GetContractsByAddress

```go
func (ref *ContractService) GetContractsByAddress(ctx context.Context, address string) ([]*ContractInfo, error)
```
returns an array of contract infos for passed customer address

#### func (*ContractService) GetContractsInfo

```go
func (ref *ContractService) GetContractsInfo(ctx context.Context, contractPubKeys ...string) ([]*ContractInfo, error)
```
returns an array of contract infos for passed public keys

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
returns a cosignature transaction from passed aggreagate bounded

#### func  NewCosignatureTransactionFromHash

```go
func NewCosignatureTransactionFromHash(hash Hash) *CosignatureTransaction
```
returns a cosignature transaction from passed hash of aggreagate bounded

#### func (*CosignatureTransaction) String

```go
func (tx *CosignatureTransaction) String() string
```

#### type Deadline

```go
type Deadline struct {
	time.Time
}
```


#### func  NewDeadline

```go
func NewDeadline(d time.Duration) *Deadline
```
returns

#### func (*Deadline) GetInstant

```go
func (d *Deadline) GetInstant() int64
```

#### type ErrorInfo

```go
type ErrorInfo struct {
	Error error
}
```


#### type Hash

```go
type Hash string
```


#### func (Hash) String

```go
func (h Hash) String() string
```

#### type HashInfo

```go
type HashInfo struct {
	Hash Hash `json:"hash"`
}
```

structure for Subscribe status

#### type HashType

```go
type HashType uint8
```


```go
const SHA3_256 HashType = 0
```

#### func (HashType) String

```go
func (ht HashType) String() string
```

#### type LockFundsTransaction

```go
type LockFundsTransaction struct {
	AbstractTransaction
	*Mosaic
	Duration *big.Int
	*SignedTransaction
}
```


#### func  NewLockFundsTransaction

```go
func NewLockFundsTransaction(deadline *Deadline, mosaic *Mosaic, duration *big.Int, signedTx *SignedTransaction, networkType NetworkType) (*LockFundsTransaction, error)
```
returns a lock funds transaction from passed deadline, mosaic, duration, and
signed transaction

#### func (*LockFundsTransaction) GetAbstractTransaction

```go
func (tx *LockFundsTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*LockFundsTransaction) String

```go
func (tx *LockFundsTransaction) String() string
```

#### type Message

```go
type Message struct {
	Type    uint8
	Payload string
}
```

Message

#### func  NewPlainMessage

```go
func NewPlainMessage(payload string) *Message
```
The transaction message of 1024 characters.

#### func (*Message) String

```go
func (m *Message) String() string
```

#### type ModifyContractTransaction

```go
type ModifyContractTransaction struct {
	AbstractTransaction
	DurationDelta int64
	Hash          string
	Customers     []*MultisigCosignatoryModification
	Executors     []*MultisigCosignatoryModification
	Verifiers     []*MultisigCosignatoryModification
}
```

ModifyContractTransaction

#### func  NewModifyContractTransaction

```go
func NewModifyContractTransaction(
	deadline *Deadline, durationDelta int64, hash string,
	customers []*MultisigCosignatoryModification,
	executors []*MultisigCosignatoryModification,
	verifiers []*MultisigCosignatoryModification,
	networkType NetworkType) (*ModifyContractTransaction, error)
```

#### func (*ModifyContractTransaction) GetAbstractTransaction

```go
func (tx *ModifyContractTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*ModifyContractTransaction) String

```go
func (tx *ModifyContractTransaction) String() string
```

#### type ModifyMultisigAccountTransaction

```go
type ModifyMultisigAccountTransaction struct {
	AbstractTransaction
	MinApprovalDelta uint8
	MinRemovalDelta  uint8
	Modifications    []*MultisigCosignatoryModification
}
```


#### func  NewModifyMultisigAccountTransaction

```go
func NewModifyMultisigAccountTransaction(deadline *Deadline, minApprovalDelta uint8, minRemovalDelta uint8, modifications []*MultisigCosignatoryModification, networkType NetworkType) (*ModifyMultisigAccountTransaction, error)
```
returns a modify multisig transaction or error from passed deadline and multisig
modification properties

#### func (*ModifyMultisigAccountTransaction) GetAbstractTransaction

```go
func (tx *ModifyMultisigAccountTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*ModifyMultisigAccountTransaction) String

```go
func (tx *ModifyMultisigAccountTransaction) String() string
```

#### type Mosaic

```go
type Mosaic struct {
	MosaicId *MosaicId
	Amount   *big.Int
}
```


#### func  NewMosaic

```go
func NewMosaic(mosaicId *MosaicId, amount *big.Int) (*Mosaic, error)
```
returns a mosaic for passed mosaic id and amount

#### func  Xem

```go
func Xem(amount int64) *Mosaic
```
returns XEM mosaic with passed amount

#### func  XemRelative

```go
func XemRelative(amount int64) *Mosaic
```
returns XEM with actual passed amount

#### func  Xpx

```go
func Xpx(amount int64) *Mosaic
```
returns XPX mosaic with passed amount

#### func  XpxRelative

```go
func XpxRelative(amount int64) *Mosaic
```
returns XPX with actual passed amount

#### func (*Mosaic) String

```go
func (m *Mosaic) String() string
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
returns mosaic definidion transaction or error from passed deadline, nonce,
public key of announcer and mosaic properties

#### func (*MosaicDefinitionTransaction) GetAbstractTransaction

```go
func (tx *MosaicDefinitionTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*MosaicDefinitionTransaction) String

```go
func (tx *MosaicDefinitionTransaction) String() string
```

#### type MosaicId

```go
type MosaicId big.Int
```

MosaicId

#### func  NewMosaicId

```go
func NewMosaicId(id *big.Int) (*MosaicId, error)
```
returns mosaic id corresponding passed big int id

#### func  NewMosaicIdFromNonceAndOwner

```go
func NewMosaicIdFromNonceAndOwner(nonce uint32, ownerPublicKey string) (*MosaicId, error)
```
returns mosaic id for passed nonce and public key of owner

#### func (*MosaicId) Equals

```go
func (m *MosaicId) Equals(id *MosaicId) bool
```

#### func (*MosaicId) String

```go
func (m *MosaicId) String() string
```

#### type MosaicInfo

```go
type MosaicInfo struct {
	MosaicId   *MosaicId
	Supply     *big.Int
	Height     *big.Int
	Owner      *PublicAccount
	Revision   uint32
	Properties *MosaicProperties
}
```

MosaicInfo info structure contains its properties, the owner and the namespace
to which it belongs to.

#### func (*MosaicInfo) String

```go
func (m *MosaicInfo) String() string
```

#### type MosaicProperties

```go
type MosaicProperties struct {
	SupplyMutable bool
	Transferable  bool
	LevyMutable   bool
	Divisibility  uint8
	Duration      *big.Int
}
```


#### func  NewMosaicProperties

```go
func NewMosaicProperties(supplyMutable bool, transferable bool, levyMutable bool, divisibility uint8, duration *big.Int) *MosaicProperties
```
TODO

#### func (*MosaicProperties) String

```go
func (mp *MosaicProperties) String() string
```

#### type MosaicService

```go
type MosaicService service
```


#### func (*MosaicService) GetMosaic

```go
func (ref *MosaicService) GetMosaic(ctx context.Context, mosaicId *MosaicId) (*MosaicInfo, error)
```
returns a mosaic info for passed mosaic id

#### func (*MosaicService) GetMosaics

```go
func (ref *MosaicService) GetMosaics(ctx context.Context, mscIds []*MosaicId) ([]*MosaicInfo, error)
```
returns an array of mosaic infos for passed mosaic ids

#### type MosaicSupplyChangeTransaction

```go
type MosaicSupplyChangeTransaction struct {
	AbstractTransaction
	MosaicSupplyType
	*MosaicId
	Delta *big.Int
}
```


#### func  NewMosaicSupplyChangeTransaction

```go
func NewMosaicSupplyChangeTransaction(deadline *Deadline, mosaicId *MosaicId, supplyType MosaicSupplyType, delta *big.Int, networkType NetworkType) (*MosaicSupplyChangeTransaction, error)
```
returns mosaic supply change transaction or error from passed deadline, mosaic
id, supply type and supply delta

#### func (*MosaicSupplyChangeTransaction) GetAbstractTransaction

```go
func (tx *MosaicSupplyChangeTransaction) GetAbstractTransaction() *AbstractTransaction
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

#### type NamespaceId

```go
type NamespaceId big.Int
```


#### func  NewNamespaceId

```go
func NewNamespaceId(id *big.Int) (*NamespaceId, error)
```
returns namespace id from passed big int representation

#### func  NewNamespaceIdFromName

```go
func NewNamespaceIdFromName(namespaceName string) (*NamespaceId, error)
```
returns namespace id from passed namespace name should be used for creating
root, child and grandchild namespace ids to create root namespace pass namespace
name in format like `rootname` to create child namespace pass namespace name in
format like `rootna.childname` to create grand child namespace pass namespace
name in format like `rootna.childname.grandchildname`

#### func (*NamespaceId) String

```go
func (m *NamespaceId) String() string
```

#### type NamespaceIds

```go
type NamespaceIds struct {
	List []*NamespaceId
}
```


#### func (*NamespaceIds) Decode

```go
func (ref *NamespaceIds) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator)
```
TODO

#### func (*NamespaceIds) Encode

```go
func (ref *NamespaceIds) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream)
```
TODO

#### func (*NamespaceIds) IsEmpty

```go
func (ref *NamespaceIds) IsEmpty(ptr unsafe.Pointer) bool
```
TODO

#### func (*NamespaceIds) MarshalJSON

```go
func (ref *NamespaceIds) MarshalJSON() (buf []byte, err error)
```
TODO

#### type NamespaceInfo

```go
type NamespaceInfo struct {
	NamespaceId *NamespaceId
	FullName    string
	Active      bool
	Index       int
	MetaId      string
	TypeSpace   NamespaceType
	Depth       int
	Levels      []*NamespaceId
	Parent      *NamespaceInfo
	Owner       *PublicAccount
	StartHeight *big.Int
	EndHeight   *big.Int
}
```


#### func (*NamespaceInfo) String

```go
func (ref *NamespaceInfo) String() string
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

#### func (*NamespaceService) GetNamespace

```go
func (ref *NamespaceService) GetNamespace(ctx context.Context, nsId *NamespaceId) (*NamespaceInfo, error)
```
returns a namespace info for passed namespace id

#### func (*NamespaceService) GetNamespaceNames

```go
func (ref *NamespaceService) GetNamespaceNames(ctx context.Context, nsIds []*NamespaceId) ([]*NamespaceName, error)
```
returns an array of namespace names for passed namespace ids

#### func (*NamespaceService) GetNamespacesFromAccount

```go
func (ref *NamespaceService) GetNamespacesFromAccount(ctx context.Context, address *Address, nsId *NamespaceId,
	pageSize int) ([]*NamespaceInfo, error)
```
TODO what is nsId returns an array of namespace infos for passed owner address
also it is possible to use pagination

#### func (*NamespaceService) GetNamespacesFromAccounts

```go
func (ref *NamespaceService) GetNamespacesFromAccounts(ctx context.Context, addrs []*Address, nsId *NamespaceId,
	pageSize int) ([]*NamespaceInfo, error)
```
TODO what is nsId returns an array of namespace infos for passed owner addresses
also it is possible to use pagination

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
returns blockchain network type

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
)
```

#### func  ExtractNetworkType

```go
func ExtractNetworkType(version uint64) NetworkType
```
TODO

#### func  NetworkTypeFromString

```go
func NetworkTypeFromString(networkType string) NetworkType
```

#### func (NetworkType) String

```go
func (nt NetworkType) String() string
```

#### type PartialRemovedInfo

```go
type PartialRemovedInfo struct {
	Meta SubscribeHash `json:"meta"`
}
```

structure for Subscribe PartialRemoved

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
returns a public account from public key for passed network type

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
	Duration     *big.Int
	ParentId     *NamespaceId
}
```


#### func  NewRegisterRootNamespaceTransaction

```go
func NewRegisterRootNamespaceTransaction(deadline *Deadline, namespaceName string, duration *big.Int, networkType NetworkType) (*RegisterNamespaceTransaction, error)
```
returns a register root namespace transaction or error from passed deadline,
namespace name and duration

#### func  NewRegisterSubNamespaceTransaction

```go
func NewRegisterSubNamespaceTransaction(deadline *Deadline, namespaceName string, parentId *NamespaceId, networkType NetworkType) (*RegisterNamespaceTransaction, error)
```
returns a register sub namespace transaction or error from passed deadline,
namespace name and parent namespace id

#### func (*RegisterNamespaceTransaction) GetAbstractTransaction

```go
func (tx *RegisterNamespaceTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*RegisterNamespaceTransaction) String

```go
func (tx *RegisterNamespaceTransaction) String() string
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

#### type SecretLockTransaction

```go
type SecretLockTransaction struct {
	AbstractTransaction
	*Mosaic
	HashType
	Duration  *big.Int
	Secret    string
	Recipient *Address
}
```


#### func  NewSecretLockTransaction

```go
func NewSecretLockTransaction(deadline *Deadline, mosaic *Mosaic, duration *big.Int, hashType HashType, secret string, recipient *Address, networkType NetworkType) (*SecretLockTransaction, error)
```
returns a secret lock transaction from passed deadline, mosaic, duration, type
of hashing, secret hashed string and mosaic recipient

#### func (*SecretLockTransaction) GetAbstractTransaction

```go
func (tx *SecretLockTransaction) GetAbstractTransaction() *AbstractTransaction
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
	Secret string
	Proof  string
}
```


#### func  NewSecretProofTransaction

```go
func NewSecretProofTransaction(deadline *Deadline, hashType HashType, secret string, proof string, networkType NetworkType) (*SecretProofTransaction, error)
```
returns a secret proof transaction from passed deadline, type of hashing, secret
hashed string and secret proof string

#### func (*SecretProofTransaction) GetAbstractTransaction

```go
func (tx *SecretProofTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*SecretProofTransaction) String

```go
func (tx *SecretProofTransaction) String() string
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


#### type StatusInfo

```go
type StatusInfo struct {
	Status string `json:"status"`
	Hash   Hash   `json:"hash"`
}
```


#### type SubscribeBlock

```go
type SubscribeBlock struct {
	Ch chan *BlockInfo
}
```


```go
var Block *SubscribeBlock
```

#### func (*SubscribeBlock) Unsubscribe

```go
func (s *SubscribeBlock) Unsubscribe() error
```

#### type SubscribeBonded

```go
type SubscribeBonded struct {
	Ch chan *AggregateTransaction
}
```


#### func (*SubscribeBonded) Unsubscribe

```go
func (s *SubscribeBonded) Unsubscribe() error
```

#### type SubscribeError

```go
type SubscribeError struct {
	Ch chan *ErrorInfo
}
```


#### func (*SubscribeError) Unsubscribe

```go
func (s *SubscribeError) Unsubscribe() error
```

#### type SubscribeHash

```go
type SubscribeHash struct {
	Ch chan *HashInfo
}
```


#### func (*SubscribeHash) Unsubscribe

```go
func (s *SubscribeHash) Unsubscribe() error
```

#### type SubscribePartialRemoved

```go
type SubscribePartialRemoved struct {
	Ch chan *PartialRemovedInfo
}
```


#### func (*SubscribePartialRemoved) Unsubscribe

```go
func (s *SubscribePartialRemoved) Unsubscribe() error
```

#### type SubscribeService

```go
type SubscribeService serviceWs
```


#### func (*SubscribeService) Block

```go
func (c *SubscribeService) Block() (*SubscribeBlock, error)
```
returns entity from which you access channel with block infos block info gets
into channel when new block is harvested

#### func (*SubscribeService) ConfirmedAdded

```go
func (c *SubscribeService) ConfirmedAdded(add *Address) (*SubscribeTransaction, error)
```
returns an entity from which you can access channel with Transaction infos for
passed address Transaction info gets into channel when it is included in a block

#### func (*SubscribeService) Cosignature

```go
func (c *SubscribeService) Cosignature(add *Address) (*SubscribeSigner, error)
```
returns an entity from which you can access channel with cosignature signed
transaction related to passed address is added to an aggregate bounded
transaction with partial state

#### func (*SubscribeService) Error

```go
func (c *SubscribeService) Error(add *Address) *SubscribeError
```
returns an entity from which you can access channel with errors related to
passed address

#### func (*SubscribeService) PartialAdded

```go
func (c *SubscribeService) PartialAdded(add *Address) (*SubscribeBonded, error)
```
returns an entity from which you can access channel with Aggregate Bonded
Transaction info Aggregate Bonded Transaction info gets into channel when it is
in partial state and waiting for actors to send all required cosignature
transactions

#### func (*SubscribeService) PartialRemoved

```go
func (c *SubscribeService) PartialRemoved(add *Address) (*SubscribePartialRemoved, error)
```
returns an entity from which you can access channel with Aggregate Bonded
Transaction hash for passed address Aggregate Bonded Transaction hash gets into
channel when it was in partial state but not anymore

#### func (*SubscribeService) Status

```go
func (c *SubscribeService) Status(add *Address) (*SubscribeStatus, error)
```
returns an entity from which you can access channel with Transaction status
infos for passed address Transaction info gets into channel when it rises an
error

#### func (*SubscribeService) UnconfirmedAdded

```go
func (c *SubscribeService) UnconfirmedAdded(add *Address) (*SubscribeTransaction, error)
```
returns an entity from which you can access channel with Transaction infos for
passed address Transaction info gets into channel when it is in unconfirmed
state and waiting to be included into a block

#### func (*SubscribeService) UnconfirmedRemoved

```go
func (c *SubscribeService) UnconfirmedRemoved(add *Address) (*SubscribeHash, error)
```
returns an entity from which you can access channel with Transaction infos for
passed address Transaction info gets into channel when it was in unconfirmed
state but not anymore

#### type SubscribeSigner

```go
type SubscribeSigner struct {
	Ch chan *SignerInfo
}
```


#### func (*SubscribeSigner) Unsubscribe

```go
func (s *SubscribeSigner) Unsubscribe() error
```

#### type SubscribeStatus

```go
type SubscribeStatus struct {
	Ch chan *StatusInfo
}
```


#### func (*SubscribeStatus) Unsubscribe

```go
func (s *SubscribeStatus) Unsubscribe() error
```

#### type SubscribeTransaction

```go
type SubscribeTransaction struct {
	Ch chan Transaction
}
```


#### func (*SubscribeTransaction) Unsubscribe

```go
func (s *SubscribeTransaction) Unsubscribe() error
```

#### type Transaction

```go
type Transaction interface {
	GetAbstractTransaction() *AbstractTransaction
	String() string
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
	Height              *big.Int
	Index               uint32
	Id                  string
	Hash                Hash
	MerkleComponentHash Hash
	AggregateHash       Hash
	AggregateId         string
}
```

Transaction Info

#### func (*TransactionInfo) String

```go
func (ti *TransactionInfo) String() string
```

#### type TransactionService

```go
type TransactionService service
```


#### func (*TransactionService) Announce

```go
func (txs *TransactionService) Announce(ctx context.Context, tx *SignedTransaction) (string, error)
```
returns transaction hash or error after announcing passed signed transaction

#### func (*TransactionService) AnnounceAggregateBonded

```go
func (txs *TransactionService) AnnounceAggregateBonded(ctx context.Context, tx *SignedTransaction) (string, error)
```
returns transaction hash or error after announcing passed signed aggregate
bounded transaction

#### func (*TransactionService) AnnounceAggregateBondedCosignature

```go
func (txs *TransactionService) AnnounceAggregateBondedCosignature(ctx context.Context, c *CosignatureSignedTransaction) (string, error)
```
returns transaction hash or error after announcing passed signed cosignature
transaction

#### func (*TransactionService) GetTransaction

```go
func (txs *TransactionService) GetTransaction(ctx context.Context, id string) (Transaction, error)
```
returns transaction information for passed transaction id or hash

#### func (*TransactionService) GetTransactionStatus

```go
func (txs *TransactionService) GetTransactionStatus(ctx context.Context, id string) (*TransactionStatus, error)
```
returns transaction status or error for passed transaction id or hash

#### func (*TransactionService) GetTransactionStatuses

```go
func (txs *TransactionService) GetTransactionStatuses(ctx context.Context, hashes []string) ([]*TransactionStatus, error)
```
returns an array of transaction statuses or error for passed transaction ids or
hashes

#### func (*TransactionService) GetTransactions

```go
func (txs *TransactionService) GetTransactions(ctx context.Context, ids []string) ([]Transaction, error)
```
returns an array of transaction informations for passed array of transaction ids
or hashes

#### type TransactionStatus

```go
type TransactionStatus struct {
	Deadline *Deadline
	Group    string
	Status   string
	Hash     Hash
	Height   *big.Int
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
	AggregateCompleted TransactionType = iota
	AggregateBonded
	MosaicDefinition
	MosaicSupplyChange
	ModifyMultisig
	ModifyContract
	RegisterNamespace
	Transfer
	Lock
	SecretLock
	SecretProof
)
```
TransactionType enums

#### func  TransactionTypeFromRaw

```go
func TransactionTypeFromRaw(value uint32) (TransactionType, error)
```

#### func (TransactionType) Hex

```go
func (t TransactionType) Hex() uint16
```

#### func (TransactionType) Raw

```go
func (t TransactionType) Raw() uint32
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
	AggregateCompletedVersion TransactionVersion = 2
	AggregateBondedVersion    TransactionVersion = 2
	MosaicDefinitionVersion   TransactionVersion = 3
	MosaicSupplyChangeVersion TransactionVersion = 2
	ModifyMultisigVersion     TransactionVersion = 3
	ModifyContractVersion     TransactionVersion = 3
	RegisterNamespaceVersion  TransactionVersion = 2
	TransferVersion           TransactionVersion = 3
	LockVersion               TransactionVersion = 1
	SecretLockVersion         TransactionVersion = 1
	SecretProofVersion        TransactionVersion = 1
)
```
TransactionVersion enums

#### type TransferTransaction

```go
type TransferTransaction struct {
	AbstractTransaction
	*Message
	Mosaics   []*Mosaic
	Recipient *Address
}
```


#### func  NewTransferTransaction

```go
func NewTransferTransaction(deadline *Deadline, recipient *Address, mosaics []*Mosaic, message *Message, networkType NetworkType) (*TransferTransaction, error)
```
returns a transfer transaction or error from passed deadline, transfer
recipient, array of mosaics to transfer and transfer message

#### func (*TransferTransaction) GetAbstractTransaction

```go
func (tx *TransferTransaction) GetAbstractTransaction() *AbstractTransaction
```

#### func (*TransferTransaction) String

```go
func (tx *TransferTransaction) String() string
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
