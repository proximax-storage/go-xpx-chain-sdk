package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
)

const (
	testHashLockInfoJson = `{
  "meta": {
    "id": "5df25d2284631392c297388c"
  },
  "lock": {
    "account": "CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE",
    "accountAddress": "9052838D173BA9DC1822C12F8CC2EA1A3B88939772B7F26D84",
    "mosaicId": [
      519256100,
      642862634
    ],
    "amount": [
      10000000,
      0
    ],
    "height": [
      256,
      0
    ],
    "status": 1,
    "hash": "67829ABA183FDA679273373C9973F23F0D8611371ED31C23C6D80FCAD0AE5C87"
  }
}`

	testSecretLockInfoJson_Used = `{
  "meta": {
    "id": "5df25d1c84631392c2973877"
  },
  "lock": {
    "account": "CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE",
    "accountAddress": "9052838D173BA9DC1822C12F8CC2EA1A3B88939772B7F26D84",
    "mosaicId": [
      519256100,
      642862634
    ],
    "amount": [
      10,
      0
    ],
    "height": [
      254,
      0
    ],
    "status": 1,
    "hashAlgorithm": 4,
    "secret": "0000000000000000000000000000000000000000000000000000000000000000",
    "recipient": "90CFA4D204CC396ED38A1BA693CB2482B58152E175BFE8B5BB",
    "compositeHash": "B8C1A5FBAA5AB8AB62444212CABB59E2E357DD9099001A29E81C166606810AA6"
  }
}`

	testSecretLockInfoJson_Unused = `{
  "meta": {
    "id": "5df25d1c84631392c2973877"
  },
  "lock": {
    "account": "CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE",
    "accountAddress": "9052838D173BA9DC1822C12F8CC2EA1A3B88939772B7F26D84",
    "mosaicId": [
      519256100,
      642862634
    ],
    "amount": [
      10,
      0
    ],
    "height": [
      254,
      0
    ],
    "status": 0,
    "hashAlgorithm": 4,
    "secret": "0000000000000000000000000000000000000000000000000000000000000000",
    "recipient": "90CFA4D204CC396ED38A1BA693CB2482B58152E175BFE8B5BB",
    "compositeHash": "B8C1A5FBAA5AB8AB62444212CABB59E2E357DD9099001A29E81C166606810AA6"
  }
}`

	testHashLockInfoJsonArr   = "[" + testHashLockInfoJson + ", " + testHashLockInfoJson + "]"
	testSecretLockInfoJsonArr = "[" + testSecretLockInfoJson_Used + ", " + testSecretLockInfoJson_Unused + "]"
)

var testLockAccount, _ = NewAccountFromPublicKey("CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE", PublicTest)
var testRecipientAddress, _ = NewAddressFromBase32("90CFA4D204CC396ED38A1BA693CB2482B58152E175BFE8B5BB")
var testLockMosaicId = newMosaicIdPanic(uint64DTO{519256100, 642862634}.toUint64())
var testHashLockHash = stringToHashPanic("67829ABA183FDA679273373C9973F23F0D8611371ED31C23C6D80FCAD0AE5C87")
var testCompositeHashLockHash = stringToHashPanic("B8C1A5FBAA5AB8AB62444212CABB59E2E357DD9099001A29E81C166606810AA6")

var (
	testLockHashInfo = &HashLockInfo{
		CommonLockInfo: CommonLockInfo{
			Account:  testLockAccount,
			MosaicId: testLockMosaicId,
			Height:   Height(256),
			Amount:   Height(10000000),
			Status:   Used,
		},
		Hash: testHashLockHash,
	}

	testLockSecretInfo_Used = &SecretLockInfo{
		CommonLockInfo: CommonLockInfo{
			Account:  testLockAccount,
			MosaicId: testLockMosaicId,
			Height:   Height(254),
			Amount:   Height(10),
			Status:   Used,
		},
		CompositeHash: testCompositeHashLockHash,
		HashAlgorithm: Internal_Hash_Type,
		Recipient:     testRecipientAddress,
		Secret:        &Hash{},
	}

	testLockSecretInfo_Unused = &SecretLockInfo{
		CommonLockInfo: CommonLockInfo{
			Account:  testLockAccount,
			MosaicId: testLockMosaicId,
			Height:   Height(254),
			Amount:   Height(10),
			Status:   Unused,
		},
		CompositeHash: testCompositeHashLockHash,
		HashAlgorithm: Internal_Hash_Type,
		Recipient:     testRecipientAddress,
		Secret:        &Hash{},
	}
)

func TestLockService_GetHashLockInfosByAccount(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(hashLocksRoute, testLockAccount.PublicKey),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testHashLockInfoJsonArr,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().Lock

	defer mock.Close()

	hashLocks, err := exchangeClient.GetHashLockInfosByAccount(ctx, testLockAccount)
	assert.Nil(t, err)
	assert.NotNil(t, hashLocks)
	assert.Equal(t, len(hashLocks), 2)
	assert.Equal(t, []*HashLockInfo{testLockHashInfo, testLockHashInfo}, hashLocks)
}

func TestLockService_GetHashLockInfo(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(hashLockRoute, testHashLockHash.String()),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testHashLockInfoJson,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().Lock

	defer mock.Close()

	hashLock, err := exchangeClient.GetHashLockInfo(ctx, testHashLockHash)
	assert.Nil(t, err)
	assert.NotNil(t, hashLock)
	assert.Equal(t, testLockHashInfo, hashLock)
}

func TestLockService_GetSecretLockInfosByAccount(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(secretLocksByAccountRoute, testLockAccount.PublicKey),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testSecretLockInfoJsonArr,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().Lock

	defer mock.Close()

	secretLocks, err := exchangeClient.GetSecretLockInfosByAccount(ctx, testLockAccount)
	assert.Nil(t, err)
	assert.NotNil(t, secretLocks)
	assert.Equal(t, len(secretLocks), 2)
	assert.Equal(t, []*SecretLockInfo{testLockSecretInfo_Used, testLockSecretInfo_Unused}, secretLocks)
}

func TestLockService_GetSecretLockInfo(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(secretLockRoute, testCompositeHashLockHash.String()),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testSecretLockInfoJson_Used,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().Lock

	defer mock.Close()

	secretLock, err := exchangeClient.GetSecretLockInfo(ctx, testCompositeHashLockHash)
	assert.Nil(t, err)
	assert.NotNil(t, secretLock)
	assert.Equal(t, testLockSecretInfo_Used, secretLock)
}

func TestLockService_GetSecretLockInfosBySecret(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(secretLocksBySecretRoute, Hash{}.String()),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testSecretLockInfoJsonArr,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().Lock

	defer mock.Close()

	secretLocks, err := exchangeClient.GetSecretLockInfosBySecret(ctx, &Hash{})
	assert.Nil(t, err)
	assert.NotNil(t, secretLocks)
	assert.Equal(t, len(secretLocks), 2)
	assert.Equal(t, []*SecretLockInfo{testLockSecretInfo_Used, testLockSecretInfo_Unused}, secretLocks)
}
