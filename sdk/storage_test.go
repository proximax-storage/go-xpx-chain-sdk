package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
)

const (
	testDriveInfoJson = `{
  "drive": {
    "multisig": "415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309",
    "multisigAddress": "9048760066A50F0F65820D3008A79CF73E1034A564BF44AB3E",
    "start": [
      2073,
      0
    ],
    "end": [
      0,
      0
    ],
    "state": 3,
    "owner": "CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE",
    "rootHash": "0100000000000000000000000000000000000000000000000000000000000000",
    "duration": [
      3,
      0
    ],
    "billingPeriod": [
      1,
      0
    ],
    "billingPrice": [
      50,
      0
    ],
    "size": [
      10000,
      0
    ],
    "replicas": 1,
    "minReplicators": 1,
    "percentApprovers": 100,
    "billingHistory": [
      {
        "start": [
          2084,
          0
        ],
        "end": [
          2085,
          0
        ],
        "payments": [
          {
            "receiver": "36E7F50C8B8BC9A4FC6325B2359E0E5DB50C75A914B5292AD726FD5AE3992691",
            "amount": [
              10,
              0
            ],
            "height": [
              2085,
              0
            ]
          }
        ]
      }
    ],
    "files": [
      {
        "fileHash": "AA2D2427E105A9B60DF634553849135DF629F1408A018D02B07A70CAFFB43093",
        "size": [
          50,
          0
        ]
      }
    ],
    "replicators": [
      {
        "replicator": "36E7F50C8B8BC9A4FC6325B2359E0E5DB50C75A914B5292AD726FD5AE3992691",
        "start": [
          2077,
          0
        ],
        "end": [
          0,
          0
        ],
        "activeFilesWithoutDeposit": [
			"AA2D2427E105A9B60DF634553849135DF629F1408A018D02B07A70CAFFB43093"
		],
        "inactiveFilesWithoutDeposit": []
      }
    ],
    "removedReplicators": [],
    "uploadPayments": [
      {
        "receiver": "CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE",
        "amount": [
          9999925,
          0
        ],
        "height": [
          2098,
          0
        ]
      }
    ]
  }
}`

	testDriveInfoJsonArr = "[" + testDriveInfoJson + ", " + testDriveInfoJson + "]"
)

var testDriveAccount, _ = NewAccountFromPublicKey("415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309", PublicTest)
var testDriveOwnerAccount, _ = NewAccountFromPublicKey("CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE", PublicTest)
var testReplicatorAccount, _ = NewAccountFromPublicKey("36E7F50C8B8BC9A4FC6325B2359E0E5DB50C75A914B5292AD726FD5AE3992691", PublicTest)
var testFileHash = stringToHashPanic("AA2D2427E105A9B60DF634553849135DF629F1408A018D02B07A70CAFFB43093")

var (
	testDriveInfo = &Drive{
		Duration:         Duration(3),
		RootHash:         &Hash{1},
		State:            Finished,
		DriveSize:        StorageSize(10000),
		BillingPeriod:    Duration(1),
		BillingPrice:     Amount(50),
		MinReplicators:   1,
		PercentApprovers: 100,
		Start:            Height(2073),
		Replicas:         1,
		OwnerAccount:     testDriveOwnerAccount,
		DriveAccount:     testDriveAccount,
		Files: map[Hash]StorageSize{
			*testFileHash: StorageSize(50),
		},
		Replicators: map[string]*ReplicatorInfo{
			testReplicatorAccount.PublicKey: &ReplicatorInfo{
				Start:   Height(2077),
				End:     Height(0),
				Account: testReplicatorAccount,
				Index:   0,
				ActiveFilesWithoutDeposit: map[Hash]bool{
					*testFileHash: true,
				},
			},
		},
		UploadPayments: []*PaymentInformation{
			&PaymentInformation{
				Amount:   Amount(9999925),
				Receiver: testDriveOwnerAccount,
				Height:   Height(2098),
			},
		},
		BillingHistory: []*BillingDescription{
			&BillingDescription{
				Start: Height(2084),
				End:   Height(2085),
				Payments: []*PaymentInformation{
					&PaymentInformation{
						Amount:   Amount(10),
						Height:   Height(2085),
						Receiver: testReplicatorAccount,
					},
				},
			},
		},
	}
)

var (
	testDrivesPage = &DrivesPage{
		Drives: []Drive{*testDriveInfo, *testDriveInfo},
	}
)

func TestStorageService_GetDrive(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(driveRoute, testDriveAccount.PublicKey),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testDriveInfoJson,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().Storage

	defer mock.Close()

	drive, err := exchangeClient.GetDrive(ctx, testDriveAccount)
	assert.Nil(t, err)
	assert.NotNil(t, drive)
	assert.Equal(t, testDriveInfo, drive)
}

func TestStorageService_GetDrives(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                drivesRoute,
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            `{ "data":` + testDriveInfoJsonArr + `}`,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().Storage

	defer mock.Close()

	drives, err := exchangeClient.GetDrives(ctx, nil, nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, drives)
	assert.Equal(t, testDrivesPage, drives)
}

func TestStorageService_GetAccountDrives(test *testing.T) {
	for _, filter := range []DriveParticipantFilter{ReplicatorDrive, OwnerDrive, AllDriveRoles} {
		test.Run(fmt.Sprintf("Test for filter %s", filter), func(t *testing.T) {
			mock := newSdkMockWithRouter(&mock.Router{
				Path:                fmt.Sprintf(drivesOfAccountRoute, testDriveOwnerAccount.PublicKey, filter),
				AcceptedHttpMethods: []string{http.MethodGet},
				RespHttpCode:        200,
				RespBody:            testDriveInfoJsonArr,
			})
			exchangeClient := mock.getPublicTestClientUnsafe().Storage

			defer mock.Close()

			drives, err := exchangeClient.GetAccountDrives(ctx, testDriveOwnerAccount, filter)
			assert.Nil(t, err)
			assert.NotNil(t, drives)
			assert.Equal(t, len(drives), 2)
			assert.Equal(t, []*Drive{testDriveInfo, testDriveInfo}, drives)
		})
	}
}

func TestStorageService_GetVerificationStatus_Available_NotActive(t *testing.T) {
	compositeHash, err := CalculateCompositeHash(&Hash{}, testDriveAccount.Address)
	assert.Nil(t, err)

	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(secretLockRoute, compositeHash.String()),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        404,
		RespBody:            testNotFoundInfoJson,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().Storage

	defer mock.Close()

	state, err := exchangeClient.GetVerificationStatus(ctx, testDriveAccount)
	assert.Nil(t, err)
	assert.NotNil(t, state)
	assert.False(t, state.Active)
	assert.True(t, state.Available)
}

func TestStorageService_GetVerificationStatus_NotAvaialble_NotActive(t *testing.T) {
	compositeHash, err := CalculateCompositeHash(&Hash{}, testDriveAccount.Address)
	assert.Nil(t, err)

	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(secretLockRoute, compositeHash.String()),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testSecretLockInfoJson_Used,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().Storage

	defer mock.Close()

	state, err := exchangeClient.GetVerificationStatus(ctx, testDriveAccount)
	assert.Nil(t, err)
	assert.NotNil(t, state)
	assert.False(t, state.Active)
	assert.False(t, state.Available)
}

func TestStorageService_GetVerificationStatus_NotAvailable_Active(t *testing.T) {
	compositeHash, err := CalculateCompositeHash(&Hash{}, testDriveAccount.Address)
	assert.Nil(t, err)

	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(secretLockRoute, compositeHash.String()),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testSecretLockInfoJson_Unused,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().Storage

	defer mock.Close()

	state, err := exchangeClient.GetVerificationStatus(ctx, testDriveAccount)
	assert.Nil(t, err)
	assert.NotNil(t, state)
	assert.True(t, state.Active)
	assert.False(t, state.Available)
}
