// Copyright 2023 ProximaX Limited. All rights reserved.
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
	testSuperContractV2InfoJson = `{
  "supercontractv2": {
    "superContractKey": "415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309",
    "driveKey": "CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE",
    "executionPaymentKey": "36E7F50C8B8BC9A4FC6325B2359E0E5DB50C75A914B5292AD726FD5AE3992691",
    "assignee": "E01D208E8539FEF6FD2E23F9CCF1300FF61199C3FE24F9FBCE30941090BD4A64",
    "creator": "5830A8E6AC1AD2775F38EA43E86BE7B686E833F27B5D22B9AD3542B3BBDF33AB",
    "deploymentBaseModificationsInfo": "AA2D2427E105A9B60DF634553849135DF629F1408A018D02B07A70CAFFB43093",
    "automaticExecutionsInfo": {
      "automaticExecutionFileName": "abc",
      "automaticExecutionsFunctionName": "def",
      "automaticExecutionsNextBlockToCheck": [
        1,
        0
      ],
      "automaticExecutionCallPayment": [
        1,
        0
      ],
      "automaticDownloadCallPayment": [
        1,
        0
      ],
      "automatedExecutionsNumber": 1,
      "automaticExecutionsPrepaidSince": [
        1,
        0
      ]
    },
    "requestCalls" : [
      {
        "callId": "0100000000000000000000000000000000000000000000000000000000000000",
        "caller": "7130A8E6AC1AD2775F38EA43E86BE7B686E833F27B5D22B9AD3542B3BBDF33AB",
        "fileName": "xyz",
        "functionName": "wst",
        "actualArguments": "uvw",
        "executionCallPayment": [
          1,
          0
        ],
        "downloadCallPayment": [
          1,
          0
        ],
        "servicePayments": [
          {
            "mosaicId": [
              1,
              0
            ],
            "amount": [
              1,
              0
            ]
          }
        ],
        "blockHeight": [
          1,
          0
        ]
      }
    ],
    "executorsInfo": [
      {
        "executorKey": "2130A8E6AC1AD2775F38EA43E86BE7B686E833F27B5D22B9AD3542B3BBDF33AB",
        "nextBatchToApprove": [
			1,
			0
		],
        "proofOfExecution": {
          "startBatchId": [
			1,
			0
		  ],
          "T": "test",
          "R": "test"
        }
      }
    ],
    "batches": [
      {
        "batchId": [
			1,
			0
		],
        "success": true,
        "poExVerificationInformation": "test",
        "completedCalls": [
          {
            "callId": "0100000000000000000000000000000000000000000000000000000000000000",
            "caller": "7130A8E6AC1AD2775F38EA43E86BE7B686E833F27B5D22B9AD3542B3BBDF33AB",
            "status": 1,
            "executionWork": [
              1,
              0
            ],
            "downloadWork": [
              1,
              0
            ]
          }
        ]
      }
    ],
    "releasedTransactions": [
			"0100000000000000000000000000000000000000000000000000000000000000"
    ]
  }
}`

	testSuperContractV2InfoJsonArr = "[" + testSuperContractV2InfoJson + ", " + testSuperContractV2InfoJson + "]"
)

var testSCKey, _ = NewAccountFromPublicKey("415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309", PublicTest)
var testDriveKey, _ = NewAccountFromPublicKey("CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE", PublicTest)
var testExecutionPaymentKey, _ = NewAccountFromPublicKey("36E7F50C8B8BC9A4FC6325B2359E0E5DB50C75A914B5292AD726FD5AE3992691", PublicTest)
var testAssignee, _ = NewAccountFromPublicKey("E01D208E8539FEF6FD2E23F9CCF1300FF61199C3FE24F9FBCE30941090BD4A64", PublicTest)
var testCreator, _ = NewAccountFromPublicKey("5830A8E6AC1AD2775F38EA43E86BE7B686E833F27B5D22B9AD3542B3BBDF33AB", PublicTest)
var testDeploymentBaseModificationsInfo = stringToHashPanic("AA2D2427E105A9B60DF634553849135DF629F1408A018D02B07A70CAFFB43093")
var testCaller, _ = NewAccountFromPublicKey("7130A8E6AC1AD2775F38EA43E86BE7B686E833F27B5D22B9AD3542B3BBDF33AB", PublicTest)
var testExecutorKey, _ = NewAccountFromPublicKey("2130A8E6AC1AD2775F38EA43E86BE7B686E833F27B5D22B9AD3542B3BBDF33AB", PublicTest)

var (
	testSuperContractInfo = &SuperContractV2{
		Account:                         testSCKey,
		DriveAccount:                    testDriveKey,
		ExecutionPaymentAccount:         testExecutionPaymentKey,
		Assignee:                        testAssignee,
		Creator:                         testCreator,
		DeploymentBaseModificationsInfo: testDeploymentBaseModificationsInfo,
		AutomaticExecutionsInfo: &AutomaticExecutionsInfo{
			AutomaticExecutionFileName:          "abc",
			AutomaticExecutionsFunctionName:     "def",
			AutomaticExecutionsNextBlockToCheck: Height(1),
			AutomaticExecutionCallPayment:       Amount(1),
			AutomaticDownloadCallPayment:        Amount(1),
			AutomatedExecutionsNumber:           1,
			AutomaticExecutionsPrepaidSince:     Height(1),
		},
		RequestedCalls: []*ContractCall{
			{
				CallId:               &Hash{1},
				Caller:               testCaller,
				FileName:             "xyz",
				FunctionName:         "wst",
				ActualArguments:      []byte{117, 118, 119},
				ExecutionCallPayment: Amount(1),
				DownloadCallPayment:  Amount(1),
				ServicePayments: []*Mosaic{
					{
						AssetId: &MosaicId{1},
						Amount:  Amount(1),
					},
				},
				BlockHeight: Height(1),
			},
		},
		ExecutorsInfo: []*ExecutorInfo{
			{
				ExecutorKey:        testExecutorKey,
				NextBatchToApprove: 1,
				PoEx: ProofOfExecution{
					StartBatchId: 1,
					T:            []byte{116, 101, 115, 116},
					R:            []byte{116, 101, 115, 116},
				},
			},
		},
		Batches: []*Batch{
			{
				BatchId:                     1,
				Success:                     true,
				PoExVerificationInformation: []byte{116, 101, 115, 116},
				CompletedCalls: []*CompletedCall{
					{
						CallId:        &Hash{1},
						Caller:        testCaller,
						Status:        1,
						ExecutionWork: Amount(1),
						DownloadWork:  Amount(1),
					},
				},
			},
		},
		ReleasedTransactions: []*Hash{
			&Hash{1},
		},
	}
)

var (
	testSuperContractsPage = &SuperContractsV2Page{
		SuperContractsV2: []*SuperContractV2{testSuperContractInfo, testSuperContractInfo},
	}
)

func TestSuperContractV2Service_GetSuperContractV2(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(superContractRouteV2, testSCKey.PublicKey),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testSuperContractV2InfoJson,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().SuperContractV2
	if exchangeClient == nil {
		fmt.Println("nil")
	}

	defer mock.Close()

	superContract, err := exchangeClient.GetSuperContractV2(ctx, testSCKey)
	assert.Nil(t, err)
	assert.NotNil(t, superContract)
	assert.Equal(t, testSuperContractInfo, superContract)
}

func TestSuperContractV2Service_GetSuperContractsV2(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                superContractsRouteV2,
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            `{ "data":` + testSuperContractV2InfoJsonArr + `}`,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().SuperContractV2

	if exchangeClient == nil {
		fmt.Println("nil")
	}
	defer mock.Close()

	supercontracts, err := exchangeClient.GetSuperContractsV2(ctx, nil)
	assert.Nil(t, err)
	assert.NotNil(t, supercontracts)
	assert.Equal(t, testSuperContractsPage, supercontracts)
}
