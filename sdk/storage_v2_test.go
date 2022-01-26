// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
)

const (
	testBcDriveInfoJson = `{
    "drive": {
      "multisig": "415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309",
      "multisigAddress": "9048760066A50F0F65820D3008A79CF73E1034A564BF44AB3E",
      "owner": "CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE",
      "rootHash": "0100000000000000000000000000000000000000000000000000000000000000",
      "size": [
        1000,
        0
      ],
      "usedSize": [
        0,
        0
      ],
      "metaFilesSize": [
        20,
        0
      ],
      "replicatorCount": 5,
	  "ownerCumulativeUploadSize": [
		  10,
		  0
	  ],
      "activeDataModifications": [
        {
          "id": "0100000000000000000000000000000000000000000000000000000000000000",
          "owner": "CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE",
          "downloadDataCdi": "0100000000000000000000000000000000000000000000000000000000000000",
          "expectedUploadSize": [
            100,
            0
          ],
		  "actualUploadSize": [
			50,
			0
		  ],
		  "folderName": "C://MyStorage",
		  "readyForApproval": false
        }
      ],
      "completedDataModifications": [
        {
          "id": "0100000000000000000000000000000000000000000000000000000000000000",
          "owner": "CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE",
          "downloadDataCdi": "0100000000000000000000000000000000000000000000000000000000000000",
		  "expectedUploadSize": [
		    100,
		    0
		  ],
		  "actualUploadSize": [
		    50,
		    0
		  ],
		  "folderName": "C://MyStorage",
		  "readyForApproval": false,
          "state": 0
        }
      ],
	  "confirmedUsedSizes": [
		{
		  "replicator": "E01D208E8539FEF6FD2E23F9CCF1300FF61199C3FE24F9FBCE30941090BD4A64",
		  "size": [
			1000,
			0
		  ]
		}
	  ],
	  "replicators": [
		"36E7F50C8B8BC9A4FC6325B2359E0E5DB50C75A914B5292AD726FD5AE3992691",
		"E01D208E8539FEF6FD2E23F9CCF1300FF61199C3FE24F9FBCE30941090BD4A64"
	  ],
	  "verifications": [
		{
		  "verificationTrigger": "0100000000000000000000000000000000000000000000000000000000000000",
		  "expiration": [
				0,
				0
			],
		  "expired": true,
		  "shards": [],
	    }
	  ]
    }
}`

	testBcDriveInfoJsonArr = "[" + testBcDriveInfoJson + ", " + testBcDriveInfoJson + "]"
)

const (
	testReplicatorInfoJson = `{
        "replicator": {
            "key": "36E7F50C8B8BC9A4FC6325B2359E0E5DB50C75A914B5292AD726FD5AE3992691",
            "version": 1,
            "capacity": [
                1000,
                0
            ],
            "drives": [
                {
                    "drive": "415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309",
                    "lastApprovedDataModificationId": "0100000000000000000000000000000000000000000000000000000000000000",
                    "dataModificationIdIsValid": false,
                    "initialDownloadWork": [
					  0,
					  0
					]
                }
            ]
        }
    }`

	testReplicatorInfoJsonArr = "[" + testReplicatorInfoJson + ", " + testReplicatorInfoJson + "]"
)

const (
	testDownloadChannelInfoJson = `{
		"downloadChannelInfo": {
			"id": "0200000000000000000000000000000000000000000000000000000000000000",
			"consumer": "5830A8E6AC1AD2775F38EA43E86BE7B686E833F27B5D22B9AD3542B3BBDF33AB",
			"drive": "415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309",
			"downloadSize": [
				500,
				0
			],
			"downloadApprovalCount": 0,
			"listOfPublicKeys": [
				"36E7F50C8B8BC9A4FC6325B2359E0E5DB50C75A914B5292AD726FD5AE3992691",
				"E01D208E8539FEF6FD2E23F9CCF1300FF61199C3FE24F9FBCE30941090BD4A64"
			],
			"cumulativePayments": [
				{
					"replicator": "36E7F50C8B8BC9A4FC6325B2359E0E5DB50C75A914B5292AD726FD5AE3992691",
					"payment": [
						300,
						0
					]
				},
				{
					"replicator": "E01D208E8539FEF6FD2E23F9CCF1300FF61199C3FE24F9FBCE30941090BD4A64",
					"payment": [
						300,
						0
					]
				}
			]
		}
	}`

	testDownloadChannelInfoJsonArr = "[" + testDownloadChannelInfoJson + ", " + testDownloadChannelInfoJson + "]"
)

var testBcDriveAccount, _ = NewAccountFromPublicKey("415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309", PublicTest)
var testBcDriveOwnerAccount, _ = NewAccountFromPublicKey("CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE", PublicTest)
var testReplicatorV2Account1, _ = NewAccountFromPublicKey("36E7F50C8B8BC9A4FC6325B2359E0E5DB50C75A914B5292AD726FD5AE3992691", PublicTest)
var testReplicatorV2Account2, _ = NewAccountFromPublicKey("E01D208E8539FEF6FD2E23F9CCF1300FF61199C3FE24F9FBCE30941090BD4A64", PublicTest)
var testConsumerAccount, _ = NewAccountFromPublicKey("5830A8E6AC1AD2775F38EA43E86BE7B686E833F27B5D22B9AD3542B3BBDF33AB", PublicTest)

var (
	testBcDriveInfo = &BcDrive{
		MultisigAccount:           testBcDriveAccount,
		Owner:                     testBcDriveOwnerAccount,
		RootHash:                  &Hash{1},
		Size:                      StorageSize(1000),
		UsedSize:                  StorageSize(0),
		MetaFilesSize:             StorageSize(20),
		ReplicatorCount:           5,
		OwnerCumulativeUploadSize: 10,
		ActiveDataModifications: []*ActiveDataModification{
			{
				Id:                 &Hash{1},
				Owner:              testBcDriveOwnerAccount,
				DownloadDataCdi:    &Hash{1},
				ExpectedUploadSize: StorageSize(100),
				ActualUploadSize:   StorageSize(50),
				FolderName:         "C://MyStorage",
				ReadyForApproval:   false,
			},
		},
		CompletedDataModifications: []*CompletedDataModification{
			{
				ActiveDataModification: &ActiveDataModification{
					Id:                 &Hash{1},
					Owner:              testBcDriveOwnerAccount,
					DownloadDataCdi:    &Hash{1},
					ExpectedUploadSize: StorageSize(100),
					ActualUploadSize:   StorageSize(50),
					FolderName:         "C://MyStorage",
					ReadyForApproval:   false,
				},
				State: DataModificationState(Succeeded),
			},
		},
		ConfirmedUsedSizes: []*ConfirmedUsedSize{
			{
				Replicator: testReplicatorV2Account2,
				Size:       StorageSize(1000),
			},
		},
		Replicators: []*PublicAccount{
			testReplicatorV2Account1,
			testReplicatorV2Account2,
		},
		Verifications: []*Verification{
			{
				VerificationTrigger: &Hash{1},
				Expiration:          blockchainTimestampDTO{0, 0}.toStruct().ToTimestamp(),
				Expired:             true,
				Shards:              []*Shard{},
			},
		},
	}

	testReplicatorInfo = &Replicator{
		ReplicatorAccount: testReplicatorV2Account1,
		Version:           1,
		Capacity:          StorageSize(1000),
		Drives: []*DriveInfo{
			{
				DriveKey:                       testBcDriveAccount,
				LastApprovedDataModificationId: &Hash{1},
				DataModificationIdIsValid:      false,
				InitialDownloadWork:            0,
			},
		},
	}

	testDownloadChannelInfo = &DownloadChannel{
		Id:                    &Hash{2},
		Consumer:              testConsumerAccount,
		Drive:                 testBcDriveAccount,
		DownloadSize:          StorageSize(500),
		DownloadApprovalCount: 0,
		ListOfPublicKeys: []*PublicAccount{
			testReplicatorV2Account1,
			testReplicatorV2Account2,
		},
		CumulativePayments: []*Payment{
			{
				Replicator: testReplicatorV2Account1,
				Payment:    Amount(300),
			},
			{
				Replicator: testReplicatorV2Account2,
				Payment:    Amount(300),
			},
		},
	}
)

var (
	testBcDrivesPage = &BcDrivesPage{
		BcDrives: []*BcDrive{testBcDriveInfo, testBcDriveInfo},
	}
	testReplicatorsPage = &ReplicatorsPage{
		Replicators: []*Replicator{testReplicatorInfo, testReplicatorInfo},
	}
	testDownloadChannelsPage = &DownloadChannelsPage{
		DownloadChannels: []*DownloadChannel{testDownloadChannelInfo, testDownloadChannelInfo},
	}
)

func TestStorageV2Service_GetDrive(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(driveRouteV2, testBcDriveAccount.PublicKey),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testBcDriveInfoJson,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().StorageV2

	defer mock.Close()

	bcdrive, err := exchangeClient.GetDrive(ctx, testBcDriveAccount)
	assert.Nil(t, err)
	assert.NotNil(t, bcdrive)
	assert.Equal(t, testBcDriveInfo, bcdrive)
}

func TestStorageV2Service_GetDrives(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                drivesRouteV2,
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            `{ "data":` + testBcDriveInfoJsonArr + `}`,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().StorageV2

	defer mock.Close()

	bcdrives, err := exchangeClient.GetDrives(ctx, nil)
	assert.Nil(t, err)
	assert.NotNil(t, bcdrives)
	assert.Equal(t, testBcDrivesPage, bcdrives)
}

func TestStorageV2Service_GetReplicator(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(replicatorRouteV2, testReplicatorV2Account1.PublicKey),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testReplicatorInfoJson,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().StorageV2

	defer mock.Close()

	replicator, err := exchangeClient.GetReplicator(ctx, testReplicatorV2Account1)
	assert.Nil(t, err)
	assert.NotNil(t, replicator)
	assert.Equal(t, testReplicatorInfo, replicator)
}

func TestStorageV2Service_GetReplicators(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                replicatorsRouteV2,
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            `{ "data":` + testReplicatorInfoJsonArr + `}`,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().StorageV2

	defer mock.Close()

	replicators, err := exchangeClient.GetReplicators(ctx, nil)
	assert.Nil(t, err)
	assert.NotNil(t, replicators)
	assert.Equal(t, testReplicatorsPage, replicators)
}

func TestStorageV2Service_GetDownloadChannelInfo(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(downloadChannelRouteV2, testDownloadChannelInfo.Id),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testDownloadChannelInfoJson,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().StorageV2

	defer mock.Close()

	downloadChannelInfo, err := exchangeClient.GetDownloadChannelInfo(ctx, testDownloadChannelInfo.Id)
	assert.Nil(t, err)
	assert.NotNil(t, downloadChannelInfo)
	assert.Equal(t, testDownloadChannelInfo, downloadChannelInfo)
}

func TestStorageV2Service_GetDownloadChannels(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                downloadChannelsRouteV2,
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            `{ "data":` + testDownloadChannelInfoJsonArr + `}`,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().StorageV2

	defer mock.Close()

	downloadChannels, err := exchangeClient.GetDownloadChannels(ctx, nil)
	assert.Nil(t, err)
	assert.NotNil(t, downloadChannels)
	assert.Equal(t, testDownloadChannelsPage, downloadChannels)
}
