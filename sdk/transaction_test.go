// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"bytes"
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/proximax-storage/go-xpx-utils/tests"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

const transactionId = "5B55E02EACCB7B00015DB6E1"
const transactionHash = "7D354E056A10E7ADAC66741D1021B0E79A57998EAD7E17198821141CE87CF63F"

var transaction = &TransferTransaction{
	AbstractTransaction: AbstractTransaction{
		Type:        Transfer,
		Version:     TransferVersion,
		NetworkType: MijinTest,
		Signature:   "ADF80CBC864B65A8D94205E9EC6640FA4AE0E3011B27F8A93D93761E454A9853BF0AB1ECB3DF62E1D2D267D3F1913FAB0E2225CE5EA3937790B78FFA1288870C",
		Signer:      &PublicAccount{&Address{MijinTest, "SBJ5D7TFIJWPY56JBEX32MUWI5RU6KVKZYITQ2HA"}, "27F6BEF9A7F75E33AE2EB2EBA10EF1D6BEA4D30EBD5E39AF8EE06E96E11AE2A9"},
		Fee:         uint64DTO{1, 0}.toBigInt(),
		Deadline:    &Deadline{time.Unix(0, uint64DTO{1094650402, 17}.toBigInt().Int64()*int64(time.Millisecond))},
		TransactionInfo: &TransactionInfo{
			Height:              uint64DTO{42, 0}.toBigInt(),
			Hash:                "45AC1259DABD7163B2816232773E66FC00342BB8DD5C965D4B784CD575FDFAF1",
			MerkleComponentHash: "45AC1259DABD7163B2816232773E66FC00342BB8DD5C965D4B784CD575FDFAF1",
			Index:               0,
			Id:                  "5B686E97F0C0EA00017B9437",
		},
	},
	Mosaics: []*Mosaic{
		{bigIntToMosaicId(uint64DTO{3646934825, 3576016193}.toBigInt()), uint64DTO{10000000, 0}.toBigInt()},
	},
	Recipient: &Address{MijinTest, "SBJUINHAC3FKCMVLL2WHBQFPPXYEHOMQY6E2SPVR"},
	Message:   &Message{Type: 0, Payload: ""},
}

var fakeDeadline = &Deadline{time.Unix(1459468800, 1000000)}

const transactionJson = `
{
   "meta":{
      "height":[42, 0],
      "hash":"45AC1259DABD7163B2816232773E66FC00342BB8DD5C965D4B784CD575FDFAF1",
      "merkleComponentHash":"45AC1259DABD7163B2816232773E66FC00342BB8DD5C965D4B784CD575FDFAF1",
      "index":0,
      "id":"5B686E97F0C0EA00017B9437"
   },
   "transaction":{
      "signature":"ADF80CBC864B65A8D94205E9EC6640FA4AE0E3011B27F8A93D93761E454A9853BF0AB1ECB3DF62E1D2D267D3F1913FAB0E2225CE5EA3937790B78FFA1288870C",
      "signer":"27F6BEF9A7F75E33AE2EB2EBA10EF1D6BEA4D30EBD5E39AF8EE06E96E11AE2A9",
      "version":36867,
      "type":16724,
      "fee":[
         1,
         0
      ],
      "deadline":[
         1094650402,
         17
      ],
      "recipient":"90534434E016CAA132AB5EAC70C0AF7DF043B990C789A93EB1",
      "message":{
         "type":0,
         "payload":""
      },
      "mosaics":[
         {
            "id":[
               3646934825,
               3576016193
            ],
            "amount":[
               10000000,
               0
            ]
         }
      ]
   }
}
`

var status = &TransactionStatus{
	&Deadline{time.Unix(uint64DTO{1, 0}.toBigInt().Int64(), int64(time.Millisecond))},
	"confirmed",
	"Success",
	"7D354E056A10E7ADAC66741D1021B0E79A57998EAD7E17198821141CE87CF63F",
	uint64DTO{1, 0}.toBigInt(),
}

const statusJson = `{
	"group": "confirmed",
	"status": "Success",
	"hash": "7D354E056A10E7ADAC66741D1021B0E79A57998EAD7E17198821141CE87CF63F",
	"deadline": [1,0],
	"height": [1, 0]
}`

var (
	aggregateTransactionSerializationCorr = []byte{0xd1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x90, 0x41, 0x41, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x55, 0x0, 0x0, 0x0, 0x55, 0x0, 0x0, 0x0, 0x84, 0x6b, 0x44, 0x39, 0x15, 0x45, 0x79, 0xa5, 0x90, 0x3b, 0x14, 0x59, 0xc9, 0xcf, 0x69, 0xcb, 0x81, 0x53, 0xf6, 0xd0, 0x11, 0xa, 0x7a, 0xe, 0xd6, 0x1d, 0xe2, 0x9a, 0xe4, 0x81, 0xb, 0xf2, 0x3, 0x90, 0x54, 0x41, 0x90, 0x50, 0xb9, 0x83, 0x7e, 0xfa, 0xb4, 0xbb, 0xe8, 0xa4, 0xb9, 0xbb, 0x32, 0xd8, 0x12, 0xf9, 0x88, 0x5c, 0x0, 0xd8, 0xfc, 0x16, 0x50, 0xe1, 0x42, 0x1, 0x0, 0x1, 0x0, 0xe3, 0x29, 0xad, 0x1c, 0xbe, 0x7f, 0xc6, 0xd, 0x80, 0x96, 0x98, 0x0, 0x0, 0x0, 0x0, 0x0}

	cosignatureTransactionSigningCorr = "bf3bc39f2292c028cb0ffa438a9f567a7c4d793d2f8522c8deac74befbcb61af6414adf27b2176d6a24fef612aa6db2f562176a11c46ba6d5e05430042cb5705"

	mosaicDefinitionTransactionSerializationCorr = []byte{0x90, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x90, 0x4d, 0x41, 0xa, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xe3, 0x29, 0xad, 0x1c, 0xbe, 0x7f, 0xc6, 0xd, 0x1, 0x7, 0x4, 0x2, 0x10, 0x27, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	mosaicSupplyChangeTransactionSerializationCorr = []byte{137, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		2, 144, 77, 66, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 136, 105, 116, 110, 155, 26, 112, 87, 1, 10, 0, 0, 0, 0, 0, 0, 0}

	transferTransactionSerializationCorr = []byte{165, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		3, 144, 84, 65, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 144, 232, 254, 189, 103, 29, 212, 27, 238, 148, 236, 59, 165, 131, 28, 182, 8, 163, 18, 194, 242, 3, 186, 132, 172,
		1, 0, 1, 0, 103, 43, 0, 0, 206, 86, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0}

	transferTransactionToAggregateCorr = []byte{85, 0, 0, 0, 154, 73, 54, 100, 6, 172, 169, 82, 184, 139, 173, 245, 241, 233, 190, 108, 228, 150, 129, 65, 3, 90, 96, 190, 80, 50, 115, 234,
		101, 69, 107, 36, 3, 144, 84, 65, 144, 232, 254, 189, 103, 29, 212, 27, 238, 148, 236, 59, 165, 131, 28, 182, 8, 163, 18, 194, 242, 3, 186, 132, 172, 1, 0, 1, 0, 103, 43, 0, 0, 206, 86, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0}

	transferTransactionSigningCorr = "A5000000773891AD01DD4CDF6E3A55C186C673E256D7DF9D471846F1943CC3529E4E02B38B9AF3F8D13784645FF5FAAFA94A321B" +
		"94933C673D12DE60E4BC05ABA56F750E1026D70E1954775749C6811084D6450A3184D977383F0E4282CD47118AF377550390544100000" +
		"00000000000010000000000000090E8FEBD671DD41BEE94EC3BA5831CB608A312C2F203BA84AC01000100672B0000CE56000064000000" +
		"00000000"

	modifyMultisigAccountTransactionSerializationCorr = []byte{189, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		3, 144, 85, 65, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2,
		0, 104, 179, 251, 177, 135, 41, 193, 253, 226, 37, 197, 127, 140, 224, 128, 250, 130, 143, 0, 103, 228, 81, 163, 253, 129, 250, 98, 136, 66, 176, 183, 99, 0, 207, 137, 63, 252, 196, 124, 51, 231, 246, 138, 177, 219, 86, 54, 92, 21, 107, 7, 54, 130, 74, 12, 30, 39, 63, 158, 0, 184, 223, 143, 1, 235}

	modifyContractTransactionSerializationCorr = []byte{105, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		3, 144, 87, 65, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0,
		2, 0, 0, 0, 0, 0, 0, 0,
		207, 137, 63, 252, 196, 124, 51, 231, 246, 138, 177, 219, 86, 54, 92, 21, 107, 7, 54, 130, 74, 12, 30, 39, 63, 158, 0, 184, 223, 143, 1, 235,
		2, 2, 2,
		0, 104, 179, 251, 177, 135, 41, 193, 253, 226, 37, 197, 127, 140, 224, 128, 250, 130, 143, 0, 103, 228, 81, 163, 253, 129, 250, 98, 136, 66, 176, 183, 99, 0, 207, 137, 63, 252, 196, 124, 51, 231, 246, 138, 177, 219, 86, 54, 92, 21, 107, 7, 54, 130, 74, 12, 30, 39, 63, 158, 0, 184, 223, 143, 1, 235,
		0, 104, 179, 251, 177, 135, 41, 193, 253, 226, 37, 197, 127, 140, 224, 128, 250, 130, 143, 0, 103, 228, 81, 163, 253, 129, 250, 98, 136, 66, 176, 183, 99, 0, 207, 137, 63, 252, 196, 124, 51, 231, 246, 138, 177, 219, 86, 54, 92, 21, 107, 7, 54, 130, 74, 12, 30, 39, 63, 158, 0, 184, 223, 143, 1, 235,
		0, 104, 179, 251, 177, 135, 41, 193, 253, 226, 37, 197, 127, 140, 224, 128, 250, 130, 143, 0, 103, 228, 81, 163, 253, 129, 250, 98, 136, 66, 176, 183, 99, 0, 207, 137, 63, 252, 196, 124, 51, 231, 246, 138, 177, 219, 86, 54, 92, 21, 107, 7, 54, 130, 74, 12, 30, 39, 63, 158, 0, 184, 223, 143, 1, 235}

	modifyBodyMetadatatTransactionSerializationCorr = []byte{0x11, 0x0, 0x0, 0x0, 0x0, 0x4, 0x5, 0x0, 0x6b, 0x65, 0x79, 0x31, 0x76, 0x61, 0x6c, 0x75, 0x65, 0xc, 0x0, 0x0, 0x0, 0x1, 0x4, 0x0, 0x0, 0x6b, 0x65, 0x79, 0x32}

	modifyAddressHeaderTransactionSerializationCorr = []byte{0xaf, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x90, 0x3d, 0x41, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x90, 0x49, 0xe1, 0x4b, 0xeb, 0xca, 0x93, 0x75, 0x8e, 0xb3, 0x68, 0x5, 0xba, 0xe7, 0x60, 0xa5, 0x72, 0x39, 0x97, 0x6f, 0x0, 0x9a, 0x54, 0x5c, 0xad}

	modifyAddressTransactionSerializationCorr = append(modifyAddressHeaderTransactionSerializationCorr, modifyBodyMetadatatTransactionSerializationCorr...)

	modifyMosaicHeaderTransactionSerializationCorr = []byte{0x9e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x90, 0x3d, 0x42, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x88, 0x69, 0x74, 0x6e, 0x9b, 0x1a, 0x70, 0x57}

	modifyMosaicTransactionSerializationCorr = append(modifyMosaicHeaderTransactionSerializationCorr, modifyBodyMetadatatTransactionSerializationCorr...)

	modifyNamespaceHeaderTransactionSerializationCorr = []byte{0x9e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x90, 0x3d, 0x43, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x88, 0x69, 0x74, 0x6e, 0x9b, 0x1a, 0x70, 0x57}

	modifyNamespaceTransactionSerializationCorr = append(modifyNamespaceHeaderTransactionSerializationCorr, modifyBodyMetadatatTransactionSerializationCorr...)

	registerRootNamespaceTransactionSerializationCorr = []byte{150, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		2, 144, 78, 65, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 16, 39, 0, 0, 0, 0, 0, 0, 126, 233, 179, 184, 175, 223, 83, 192, 12, 110, 101, 119, 110, 97, 109, 101, 115, 112, 97, 99, 101}

	registerSubNamespaceTransactionSerializationCorr = []byte{0x96, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x90, 0x4e, 0x41, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x7e, 0xe9, 0xb3, 0xb8, 0xaf, 0xdf, 0x53, 0x40, 0x3, 0x12, 0x98, 0x1b, 0x78, 0x79, 0xa3, 0xf1, 0xc, 0x73, 0x75, 0x62, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65}

	lockFundsTransactionSerializationCorr = []byte{0xb0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x90, 0x48, 0x41, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xe3, 0x29, 0xad, 0x1c, 0xbe, 0x7f, 0xc6, 0xd, 0x80, 0x96, 0x98, 0x0, 0x0, 0x0, 0x0, 0x0, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x84, 0x98, 0xb3, 0x8d, 0x89, 0xc1, 0xdc, 0x8a, 0x44, 0x8e, 0xa5, 0x82, 0x49, 0x38, 0xff, 0x82, 0x89, 0x26, 0xcd, 0x9f, 0x77, 0x47, 0xb1, 0x84, 0x4b, 0x59, 0xb4, 0xb6, 0x80, 0x7e, 0x87, 0x8b}

	secretProofTransactionSigningCorr = "9F0000007164460939684C5BCA586C4F8FCD457677743E17E3AE5CCC6CFAD7865DB47DDDF5C228B24B99F5CCE18BE0788270A5A5B4750A1DBEB213B2DEB4F220DB75E5071026D70E1954775749C6811084D6450A3184D977383F0E4282CD47118AF37755019052420000000000000000010000000000000000B778A39A3663719DFC5E48C9D78431B1E45C2AF9DF538782BF199C189DABEAC704009A493664"

	secretProofTransactionToAggregateCorr = []byte{0x4f, 0x0, 0x0, 0x0, 0x9a, 0x49, 0x36, 0x64, 0x6, 0xac, 0xa9, 0x52, 0xb8, 0x8b, 0xad, 0xf5, 0xf1, 0xe9, 0xbe, 0x6c, 0xe4, 0x96, 0x81, 0x41, 0x3, 0x5a, 0x60, 0xbe, 0x50, 0x32, 0x73, 0xea, 0x65, 0x45, 0x6b, 0x24, 0x1, 0x90, 0x52, 0x42, 0x0, 0xb7, 0x78, 0xa3, 0x9a, 0x36, 0x63, 0x71, 0x9d, 0xfc, 0x5e, 0x48, 0xc9, 0xd7, 0x84, 0x31, 0xb1, 0xe4, 0x5c, 0x2a, 0xf9, 0xdf, 0x53, 0x87, 0x82, 0xbf, 0x19, 0x9c, 0x18, 0x9d, 0xab, 0xea, 0xc7, 0x4, 0x0, 0x9a, 0x49, 0x36, 0x64}

	secretProofTransactionSerializationCorr = []byte{0x9f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x90, 0x52, 0x42, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xb7, 0x78, 0xa3, 0x9a, 0x36, 0x63, 0x71, 0x9d, 0xfc, 0x5e, 0x48, 0xc9, 0xd7, 0x84, 0x31, 0xb1, 0xe4, 0x5c, 0x2a, 0xf9, 0xdf, 0x53, 0x87, 0x82, 0xbf, 0x19, 0x9c, 0x18, 0x9d, 0xab, 0xea, 0xc7, 0x4, 0x0, 0x9a, 0x49, 0x36, 0x64}

	secretLockTransactionToAggregateCorr = []byte{0x7a, 0x0, 0x0, 0x0, 0x9a, 0x49, 0x36, 0x64, 0x6, 0xac, 0xa9, 0x52, 0xb8, 0x8b, 0xad, 0xf5, 0xf1, 0xe9, 0xbe, 0x6c, 0xe4, 0x96, 0x81, 0x41, 0x3, 0x5a, 0x60, 0xbe, 0x50, 0x32, 0x73, 0xea, 0x65, 0x45, 0x6b, 0x24, 0x1, 0x90, 0x52, 0x41, 0xe3, 0x29, 0xad, 0x1c, 0xbe, 0x7f, 0xc6, 0xd, 0x80, 0x96, 0x98, 0x0, 0x0, 0x0, 0x0, 0x0, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xb7, 0x78, 0xa3, 0x9a, 0x36, 0x63, 0x71, 0x9d, 0xfc, 0x5e, 0x48, 0xc9, 0xd7, 0x84, 0x31, 0xb1, 0xe4, 0x5c, 0x2a, 0xf9, 0xdf, 0x53, 0x87, 0x82, 0xbf, 0x19, 0x9c, 0x18, 0x9d, 0xab, 0xea, 0xc7, 0x90, 0xe8, 0xfe, 0xbd, 0x67, 0x1d, 0xd4, 0x1b, 0xee, 0x94, 0xec, 0x3b, 0xa5, 0x83, 0x1c, 0xb6, 0x8, 0xa3, 0x12, 0xc2, 0xf2, 0x3, 0xba, 0x84, 0xac}

	secretLockTransactionSerializationCorr = []byte{0xca, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x90, 0x52, 0x41, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xe3, 0x29, 0xad, 0x1c, 0xbe, 0x7f, 0xc6, 0xd, 0x80, 0x96, 0x98, 0x0, 0x0, 0x0, 0x0, 0x0, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xb7, 0x78, 0xa3, 0x9a, 0x36, 0x63, 0x71, 0x9d, 0xfc, 0x5e, 0x48, 0xc9, 0xd7, 0x84, 0x31, 0xb1, 0xe4, 0x5c, 0x2a, 0xf9, 0xdf, 0x53, 0x87, 0x82, 0xbf, 0x19, 0x9c, 0x18, 0x9d, 0xab, 0xea, 0xc7, 0x90, 0xe8, 0xfe, 0xbd, 0x67, 0x1d, 0xd4, 0x1b, 0xee, 0x94, 0xec, 0x3b, 0xa5, 0x83, 0x1c, 0xb6, 0x8, 0xa3, 0x12, 0xc2, 0xf2, 0x3, 0xba, 0x84, 0xac}

	lockFundsTransactionToAggregateCorr = []byte{0x60, 0x0, 0x0, 0x0, 0x9a, 0x49, 0x36, 0x64, 0x6, 0xac, 0xa9, 0x52, 0xb8, 0x8b, 0xad, 0xf5, 0xf1, 0xe9, 0xbe, 0x6c, 0xe4, 0x96, 0x81, 0x41, 0x3, 0x5a, 0x60, 0xbe, 0x50, 0x32, 0x73, 0xea, 0x65, 0x45, 0x6b, 0x24, 0x1, 0x90, 0x48, 0x41, 0xe3, 0x29, 0xad, 0x1c, 0xbe, 0x7f, 0xc6, 0xd, 0x80, 0x96, 0x98, 0x0, 0x0, 0x0, 0x0, 0x0, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x84, 0x98, 0xb3, 0x8d, 0x89, 0xc1, 0xdc, 0x8a, 0x44, 0x8e, 0xa5, 0x82, 0x49, 0x38, 0xff, 0x82, 0x89, 0x26, 0xcd, 0x9f, 0x77, 0x47, 0xb1, 0x84, 0x4b, 0x59, 0xb4, 0xb6, 0x80, 0x7e, 0x87, 0x8b}

	secretLockTransactionSigningCorr = "CA000000ACCEA3B564B8D2375EE8EAB968E2BCB9101D0D338D7354CE02844F096C35A78A8753393C4605C9657C580E080E6281885B99BB33F41571AAC416B06F792EF0031026D70E1954775749C6811084D6450A3184D977383F0E4282CD47118AF377550190524100000000000000000100000000000000E329AD1CBE7FC60D8096980000000000640000000000000000B778A39A3663719DFC5E48C9D78431B1E45C2AF9DF538782BF199C189DABEAC790E8FEBD671DD41BEE94EC3BA5831CB608A312C2F203BA84AC"

	lockFundsTransactionSigningCorr = "B0000000D01A4B5AC16A61A129D68B9AF507CDCAA376444E9DDDB9FDEAEB25E26B2800CE45265BE388D6ACB6D61625CB19AC7DA22D59C7164A1F82D561F37C19342A85051026D70E1954775749C6811084D6450A3184D977383F0E4282CD47118AF377550190484100000000000000000100000000000000E329AD1CBE7FC60D809698000000000064000000000000008498B38D89C1DC8A448EA5824938FF828926CD9F7747B1844B59B4B6807E878B"
)

func TestTransactionService_GetTransaction_TransferTransaction(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     fmt.Sprintf("/transaction/%s", transactionId),
		RespBody: transactionJson,
	})

	cl := mockServer.getPublicTestClientUnsafe()

	tx, err := cl.Transaction.GetTransaction(context.Background(), transactionId)

	assert.Nilf(t, err, "TransactionService.GetTransaction returned error: %v", err)

	tests.ValidateStringers(t, transaction, tx)
}

func TestTransactionService_GetTransactions(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     "/transaction",
		RespBody: "[" + transactionJson + "]",
	})

	cl := mockServer.getPublicTestClientUnsafe()

	transactions, err := cl.Transaction.GetTransactions(context.Background(), []string{
		transactionId,
	})

	assert.Nilf(t, err, "TransactionService.GetTransactions returned error: %v", err)

	for _, tx := range transactions {
		tests.ValidateStringers(t, transaction, tx)
	}
}

func TestTransactionService_GetTransactionStatus(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     "/transaction/7D354E056A10E7ADAC66741D1021B0E79A57998EAD7E17198821141CE87CF63F/status",
		RespBody: statusJson,
	})

	cl := mockServer.getPublicTestClientUnsafe()

	txStatus, err := cl.Transaction.GetTransactionStatus(context.Background(), transactionHash)

	assert.Nilf(t, err, "TransactionService.GetTransactionStatus returned error: %v", err)

	tests.ValidateStringers(t, status, txStatus)
}

func TestTransactionService_GetTransactionStatuses(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     "/transaction/statuses",
		RespBody: "[" + statusJson + "]",
	})

	cl := mockServer.getPublicTestClientUnsafe()

	txStatuses, err := cl.Transaction.GetTransactionStatuses(context.Background(), []string{transactionHash})

	assert.Nilf(t, err, "TransactionService.GetTransactionStatuses returned error: %v", err)

	for _, txStatus := range txStatuses {
		tests.ValidateStringers(t, status, txStatus)
	}
}

func TestAggregateTransactionSerialization(t *testing.T) {
	p, err := NewAccountFromPublicKey("846B4439154579A5903B1459C9CF69CB8153F6D0110A7A0ED61DE29AE4810BF2", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	ttx, err := NewTransferTransaction(
		fakeDeadline,
		NewAddress("SBILTA367K2LX2FEXG5TFWAS7GEFYAGY7QLFBYKC", MijinTest),
		[]*Mosaic{Xem(10000000)},
		NewPlainMessage(""),
		MijinTest,
	)

	assert.Nilf(t, err, "NewTransferTransaction returned error: %s", err)

	ttx.Signer = p

	atx, err := NewCompleteAggregateTransaction(fakeDeadline, []Transaction{ttx}, MijinTest)

	assert.Nilf(t, err, "NewCompleteAggregateTransaction returned error: %s", err)

	b, err := atx.generateBytes()

	assert.Nilf(t, err, "AggregateTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, aggregateTransactionSerializationCorr, b)
}

func TestAggregateTransactionSigningWithMultipleCosignatures(t *testing.T) {
	p, err := NewAccountFromPublicKey("B694186EE4AB0558CA4AFCFDD43B42114AE71094F5A1FC4A913FE9971CACD21D", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	ttx, err := NewTransferTransaction(
		fakeDeadline,
		NewAddress("SBILTA367K2LX2FEXG5TFWAS7GEFYAGY7QLFBYKC", MijinTest),
		[]*Mosaic{},
		NewPlainMessage("test-message"),
		MijinTest,
	)

	ttx.Signer = p

	atx, err := NewCompleteAggregateTransaction(fakeDeadline, []Transaction{ttx}, MijinTest)

	assert.Nilf(t, err, "NewCompleteAggregateTransaction returned error: %s", err)

	acc1, err := NewAccountFromPrivateKey("2a2b1f5d366a5dd5dc56c3c757cf4fe6c66e2787087692cf329d7a49a594658b", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPrivateKey returned error: %s", err)

	acc2, err := NewAccountFromPrivateKey("b8afae6f4ad13a1b8aad047b488e0738a437c7389d4ff30c359ac068910c1d59", MijinTest) // TODO from original repo: "bug with private key"

	assert.Nilf(t, err, "NewAccountFromPrivateKey returned error: %s", err)

	stx, err := acc1.SignWithCosignatures(atx, []*Account{acc2})

	assert.Nilf(t, err, "Account.SignWithCosignatures returned error: %s", err)
	assert.Equal(t, "2d010000", stx.Payload[0:8])
	assert.Equal(t, "5100000051000000", stx.Payload[240:256])

	//if !reflect.DeepEqual(stx.Payload[320:474], "039054419050B9837EFAB4BBE8A4B9BB32D812F9885C00D8FC1650E1420D000000746573742D6D65737361676568B3FBB18729C1FDE225C57F8CE080FA828F0067E451A3FD81FA628842B0B763") {
	//	t.Errorf("AggregateTransaction signing returned wrong payload: \n %s", stx.Payload[320:474])
	//} this test is not working in original repo and commented out too
}

func TestCosignatureTransactionSigning(t *testing.T) {
	rtx := "{\"meta\":{\"hash\":\"671653C94E2254F2A23EFEDB15D67C38332AED1FBD24B063C0A8E675582B6A96\",\"height\":[18160,0],\"id\":\"5A0069D83F17CF0001777E55\",\"index\":0,\"merkleComponentHash\":\"81E5E7AE49998802DABC816EC10158D3A7879702FF29084C2C992CD1289877A7\"},\"transaction\":{\"cosignatures\":[{\"signature\":\"5780C8DF9D46BA2BCF029DCC5D3BF55FE1CB5BE7ABCF30387C4637DDEDFC2152703CA0AD95F21BB9B942F3CC52FCFC2064C7B84CF60D1A9E69195F1943156C07\",\"signer\":\"A5F82EC8EBB341427B6785C8111906CD0DF18838FB11B51CE0E18B5E79DFF630\"}],\"deadline\":[3266625578,11],\"fee\":[1,0],\"signature\":\"939673209A13FF82397578D22CC96EB8516A6760C894D9B7535E3A1E068007B9255CFA9A914C97142A7AE18533E381C846B69D2AE0D60D1DC8A55AD120E2B606\",\"signer\":\"7681ED5023141D9CDCF184E5A7B60B7D466739918ED5DA30F7E71EA7B86EFF2D\",\"transactions\":[{\"meta\":{\"aggregateHash\":\"3D28C804EDD07D5A728E5C5FFEC01AB07AFA5766AE6997B38526D36015A4D006\",\"aggregateId\":\"5A0069D83F17CF0001777E55\",\"height\":[18160,0],\"id\":\"5A0069D83F17CF0001777E56\",\"index\":0},\"transaction\":{\"message\":{\"payload\":\"746573742D6D657373616765\",\"type\":0},\"mosaics\":[{\"amount\":[3863990592,95248],\"id\":[3646934825,3576016193]}],\"recipient\":\"9050B9837EFAB4BBE8A4B9BB32D812F9885C00D8FC1650E142\",\"signer\":\"B4F12E7C9F6946091E2CB8B6D3A12B50D17CCBBF646386EA27CE2946A7423DCF\",\"type\":16724,\"version\":36867}}],\"type\":16705,\"version\":36867}}"
	b := bytes.NewBufferString(rtx)
	tx, err := MapTransaction(b)

	assert.Nilf(t, err, "MapTransaction returned error: %s", err)

	atx := tx.(*AggregateTransaction)

	acc, err := NewAccountFromPrivateKey("26b64cb10f005e5988a36744ca19e20d835ccc7c105aaa5f3b212da593180930", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPrivateKey returned error: %s", err)

	ctx, err := NewCosignatureTransaction(atx)

	assert.Nilf(t, err, "NewCosignatureTransaction returned error: %s", err)

	cstx, err := acc.SignCosignatureTransaction(ctx)

	assert.Nilf(t, err, "Account.SignCosignatureTransaction signing returned error: %s", err)
	assert.Equal(t, cosignatureTransactionSigningCorr, cstx.Signature)
}

func TestModifyAddressMetadataTransactionSerialization(t *testing.T) {
	acc, err := NewAccountFromPublicKey("68b3fbb18729c1fde225c57f8ce080fa828f0067e451a3fd81fa628842b0b763", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	tx, err := NewModifyMetadataAddressTransaction(
		fakeDeadline,
		acc.Address,
		[]*MetadataModification{
			{
				AddMetadata,
				"key1",
				"value",
			},
			{
				RemoveMetadata,
				"key2",
				"",
			},
		},
		MijinTest,
	)

	assert.Nilf(t, err, "NewModifyMetadataAddressTransaction returned error: %s", err)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "NewModifyMetadataAddressTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, modifyAddressTransactionSerializationCorr, b)
}

func TestModifyMosaicMetadataTransactionSerialization(t *testing.T) {
	id := bigIntToMosaicId(big.NewInt(6300565133566699912))
	tx, err := NewModifyMetadataMosaicTransaction(
		fakeDeadline,
		id,
		[]*MetadataModification{
			{
				AddMetadata,
				"key1",
				"value",
			},
			{
				RemoveMetadata,
				"key2",
				"",
			},
		},
		MijinTest,
	)

	assert.Nilf(t, err, "NewModifyMetadataMosaicTransaction returned error: %s", err)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "NewModifyMetadataMosaicTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, modifyMosaicTransactionSerializationCorr, b)
}

func TestModifyNamespaceMetadataTransactionSerialization(t *testing.T) {
	id := bigIntToNamespaceId(big.NewInt(6300565133566699912))
	tx, err := NewModifyMetadataNamespaceTransaction(
		fakeDeadline,
		id,
		[]*MetadataModification{
			{
				AddMetadata,
				"key1",
				"value",
			},
			{
				RemoveMetadata,
				"key2",
				"",
			},
		},
		MijinTest,
	)

	assert.Nilf(t, err, "NewModifyMetadataNamespaceTransaction returned error: %s", err)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "NewModifyMetadataNamespaceTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, modifyNamespaceTransactionSerializationCorr, b)
}

func TestMosaicDefinitionTransactionSerialization(t *testing.T) {
	account, err := NewAccountFromPrivateKey("C06B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPrivateKey returned error: %s", err)

	tx, err := NewMosaicDefinitionTransaction(
		fakeDeadline,
		0,
		account.PublicAccount.PublicKey,
		NewMosaicProperties(true, true, true, 4, big.NewInt(10000)),
		MijinTest)

	tx.Fee = big.NewInt(10)

	assert.Nilf(t, err, "NewMosaicDefinitionTransaction returned error: %s", err)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "MosaicDefinitionTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, mosaicDefinitionTransactionSerializationCorr, b)
}

func TestMosaicSupplyChangeTransactionSerialization(t *testing.T) {
	id := bigIntToMosaicId(big.NewInt(6300565133566699912))
	tx, err := NewMosaicSupplyChangeTransaction(fakeDeadline, id, Increase, big.NewInt(10), MijinTest)

	assert.Nilf(t, err, "NewMosaicSupplyChangeTransaction returned error: %s", err)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "MosaicSupplyChangeTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, mosaicSupplyChangeTransactionSerializationCorr, b)
}

func TestTransferTransactionSerialization(t *testing.T) {
	tx, err := NewTransferTransaction(
		fakeDeadline,
		NewAddress("SDUP5PLHDXKBX3UU5Q52LAY4WYEKGEWC6IB3VBFM", MijinTest),
		[]*Mosaic{
			{
				MosaicId: bigIntToMosaicId(big.NewInt(95442763262823)),
				Amount:   big.NewInt(100),
			},
		},
		NewPlainMessage(""),
		MijinTest,
	)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "TransferTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, transferTransactionSerializationCorr, b)
}

func TestTransferTransactionToAggregate(t *testing.T) {
	p, err := NewAccountFromPublicKey("9A49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	tx, err := NewTransferTransaction(
		fakeDeadline,
		NewAddress("SDUP5PLHDXKBX3UU5Q52LAY4WYEKGEWC6IB3VBFM", MijinTest),
		[]*Mosaic{{bigIntToMosaicId(big.NewInt(95442763262823)), big.NewInt(100)}},
		NewPlainMessage(""),
		MijinTest,
	)

	assert.Nilf(t, err, "NewTransferTransaction returned error: %s", err)

	tx.Signer = p

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, transferTransactionToAggregateCorr, b)
}

func TestTransferTransactionSigning(t *testing.T) {
	a, err := NewAccountFromPrivateKey("787225aaff3d2c71f4ffa32d4f19ec4922f3cd869747f267378f81f8e3fcb12d", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPrivateKey returned error: %s", err)

	tx, err := NewTransferTransaction(
		fakeDeadline,
		NewAddress("SDUP5PLHDXKBX3UU5Q52LAY4WYEKGEWC6IB3VBFM", MijinTest),
		[]*Mosaic{{bigIntToMosaicId(big.NewInt(95442763262823)), big.NewInt(100)}},
		NewPlainMessage(""),
		MijinTest,
	)

	assert.Nilf(t, err, "NewTransferTransaction returned error: %s", err)

	stx, err := a.Sign(tx)

	assert.Nilf(t, err, "Account.Sign returned error: %s", err)
	assert.Equal(t, transferTransactionSigningCorr, stx.Payload)
	assert.Equal(t, "350AE56BC97DB805E2098AB2C596FA4C6B37EF974BF24DFD61CD9F77C7687424", stx.Hash.String())
}

func TestModifyMultisigAccountTransactionSerialization(t *testing.T) {
	acc1, err := NewAccountFromPublicKey("68b3fbb18729c1fde225c57f8ce080fa828f0067e451a3fd81fa628842b0b763", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	acc2, err := NewAccountFromPublicKey("cf893ffcc47c33e7f68ab1db56365c156b0736824a0c1e273f9e00b8df8f01eb", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	tx, err := NewModifyMultisigAccountTransaction(
		fakeDeadline,
		2,
		1,
		[]*MultisigCosignatoryModification{
			{
				Add,
				acc1,
			},
			{
				Add,
				acc2,
			},
		},
		MijinTest,
	)

	assert.Nilf(t, err, "NewModifyMultisigAccountTransaction returned error: %s", err)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "ModifyMultisigAccountTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, modifyMultisigAccountTransactionSerializationCorr, b)
}

func TestModifyContractTransactionSerialization(t *testing.T) {
	acc1, err := NewAccountFromPublicKey("68b3fbb18729c1fde225c57f8ce080fa828f0067e451a3fd81fa628842b0b763", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	acc2, err := NewAccountFromPublicKey("cf893ffcc47c33e7f68ab1db56365c156b0736824a0c1e273f9e00b8df8f01eb", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	tx, err := NewModifyContractTransaction(
		fakeDeadline,
		2,
		"cf893ffcc47c33e7f68ab1db56365c156b0736824a0c1e273f9e00b8df8f01eb",
		[]*MultisigCosignatoryModification{
			{
				Add,
				acc1,
			},
			{
				Add,
				acc2,
			},
		},
		[]*MultisigCosignatoryModification{
			{
				Add,
				acc1,
			},
			{
				Add,
				acc2,
			},
		},
		[]*MultisigCosignatoryModification{
			{
				Add,
				acc1,
			},
			{
				Add,
				acc2,
			},
		},
		MijinTest,
	)

	assert.Nilf(t, err, "NewModifyContractTransaction returned error: %s", err)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "ModifyContractTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, modifyContractTransactionSerializationCorr, b)
}

func TestRegisterRootNamespaceTransactionSerialization(t *testing.T) {
	tx, err := NewRegisterRootNamespaceTransaction(
		fakeDeadline,
		"newnamespace",
		big.NewInt(10000),
		MijinTest,
	)

	assert.Nilf(t, err, "NewRegisterRootNamespaceTransaction returned error: %s", err)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "RegisterNamespaceTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, registerRootNamespaceTransactionSerializationCorr, b)
}

func TestRegisterSubNamespaceTransactionSerialization(t *testing.T) {
	tx, err := NewRegisterSubNamespaceTransaction(
		fakeDeadline,
		"subnamespace",
		bigIntToNamespaceId(big.NewInt(4635294387305441662)),
		MijinTest,
	)

	assert.Nilf(t, err, "NewRegisterSubNamespaceTransaction returned error: %s", err)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "RegisterNamespaceTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, registerSubNamespaceTransactionSerializationCorr, b)
}

func TestLockFundsTransactionSerialization(t *testing.T) {
	stx := &SignedTransaction{AggregateBonded, "payload", "8498B38D89C1DC8A448EA5824938FF828926CD9F7747B1844B59B4B6807E878B"}

	tx, err := NewLockFundsTransaction(fakeDeadline, XemRelative(10), big.NewInt(100), stx, MijinTest)

	assert.Nilf(t, err, "NewLockFundsTransaction returned error: %s", err)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "LockFundsTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, lockFundsTransactionSerializationCorr, b)
}

func TestLockFundsTransactionToAggregate(t *testing.T) {
	p, err := NewAccountFromPublicKey("9A49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	stx := &SignedTransaction{AggregateBonded, "payload", "8498B38D89C1DC8A448EA5824938FF828926CD9F7747B1844B59B4B6807E878B"}

	tx, err := NewLockFundsTransaction(fakeDeadline, XemRelative(10), big.NewInt(100), stx, MijinTest)

	assert.Nilf(t, err, "NewLockFundsTransaction returned error: %s", err)

	tx.Signer = p

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, lockFundsTransactionToAggregateCorr, b)
}

func TestLockFundsTransactionSigning(t *testing.T) {
	acc, err := NewAccountFromPrivateKey("787225aaff3d2c71f4ffa32d4f19ec4922f3cd869747f267378f81f8e3fcb12d", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPrivateKey returned error: %s", err)

	stx := &SignedTransaction{AggregateBonded, "payload", "8498B38D89C1DC8A448EA5824938FF828926CD9F7747B1844B59B4B6807E878B"}

	tx, err := NewLockFundsTransaction(fakeDeadline, XemRelative(10), big.NewInt(100), stx, MijinTest)

	assert.Nilf(t, err, "NewLockFundsTransaction returned error: %s", err)

	b, err := signTransactionWith(tx, acc)

	assert.Nilf(t, err, "signTransactionWith returned error: %s", err)
	assert.Equal(t, lockFundsTransactionSigningCorr, b.Payload)
	assert.Equal(t, "A01F0533159B3AE4C684BC72E641EF6456BE8A07DEB01838A9696E975FAD4017", b.Hash.String())
}

func TestSecretLockTransactionSerialization(t *testing.T) {
	s := "b778a39a3663719dfc5e48c9d78431b1e45c2af9df538782bf199c189dabeac7"

	ad, err := NewAddressFromRaw("SDUP5PLHDXKBX3UU5Q52LAY4WYEKGEWC6IB3VBFM")

	assert.Nilf(t, err, "NewAddressFromRaw returned error: %s", err)

	tx, err := NewSecretLockTransaction(fakeDeadline, XemRelative(10), big.NewInt(100), SHA3_256, s, ad, MijinTest)

	assert.Nilf(t, err, "NewSecretLockTransaction returned error: %s", err)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "SecretLockTransaction.generateBytes returned error: %s", err)
	assert.Equal(t, secretLockTransactionSerializationCorr, b)
}

func TestSecretLockTransactionToAggregate(t *testing.T) {
	p, err := NewAccountFromPublicKey("9A49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	ad, err := NewAddressFromRaw("SDUP5PLHDXKBX3UU5Q52LAY4WYEKGEWC6IB3VBFM")

	assert.Nilf(t, err, "NewAddressFromRaw returned error: %s", err)

	s := "b778a39a3663719dfc5e48c9d78431b1e45c2af9df538782bf199c189dabeac7"

	tx, err := NewSecretLockTransaction(fakeDeadline, XemRelative(10), big.NewInt(100), SHA3_256, s, ad, MijinTest)

	assert.Nilf(t, err, "NewSecretLockTransaction returned error: %s", err)

	tx.Signer = p

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, secretLockTransactionToAggregateCorr, b)
}

func TestSecretLockTransactionSigning(t *testing.T) {
	s := "b778a39a3663719dfc5e48c9d78431b1e45c2af9df538782bf199c189dabeac7"

	acc, err := NewAccountFromPrivateKey("787225aaff3d2c71f4ffa32d4f19ec4922f3cd869747f267378f81f8e3fcb12d", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPrivateKey returned error: %s", err)

	ad, err := NewAddressFromRaw("SDUP5PLHDXKBX3UU5Q52LAY4WYEKGEWC6IB3VBFM")

	assert.Nilf(t, err, "NewAddressFromRaw returned error: %s", err)

	tx, err := NewSecretLockTransaction(fakeDeadline, XemRelative(10), big.NewInt(100), SHA3_256, s, ad, MijinTest)

	assert.Nilf(t, err, "NewSecretLockTransaction returned error: %s", err)

	b, err := acc.Sign(tx)

	assert.Nilf(t, err, "Sign returned error: %s", err)
	assert.Equal(t, secretLockTransactionSigningCorr, b.Payload)
	assert.Equal(t, "D61DD2A3C4D736EF211323370F41F4C76B8FB7EF1E3C1737C0D0BD14B24312D2", b.Hash.String())
}

func TestSecretProofTransactionSerialization(t *testing.T) {
	s := "b778a39a3663719dfc5e48c9d78431b1e45c2af9df538782bf199c189dabeac7"
	ss := "9a493664"

	tx, err := NewSecretProofTransaction(fakeDeadline, SHA3_256, s, ss, MijinTest)

	assert.Nilf(t, err, "NewSecretProofTransaction returned error: %s", err)

	b, err := tx.generateBytes()

	assert.Nilf(t, err, "generateBytes returned error: %s", err)
	assert.Equal(t, secretProofTransactionSerializationCorr, b)
}

func TestSecretProofTransactionToAggregate(t *testing.T) {
	p, err := NewAccountFromPublicKey("9A49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	s := "b778a39a3663719dfc5e48c9d78431b1e45c2af9df538782bf199c189dabeac7"
	ss := "9a493664"

	tx, err := NewSecretProofTransaction(fakeDeadline, SHA3_256, s, ss, MijinTest)

	assert.Nilf(t, err, "NewSecretProofTransaction returned error: %s", err)

	tx.Signer = p

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, secretProofTransactionToAggregateCorr, b)
}

func TestSecretProofTransactionSigning(t *testing.T) {
	acc, err := NewAccountFromPrivateKey("787225aaff3d2c71f4ffa32d4f19ec4922f3cd869747f267378f81f8e3fcb12d", MijinTest)

	assert.Nilf(t, err, "NewAccountFromPrivateKey returned error: %s", err)

	s := "b778a39a3663719dfc5e48c9d78431b1e45c2af9df538782bf199c189dabeac7"
	ss := "9a493664"

	tx, err := NewSecretProofTransaction(fakeDeadline, SHA3_256, s, ss, MijinTest)

	assert.Nilf(t, err, "NewSecretProofTransaction returned error: %s", err)

	b, err := signTransactionWith(tx, acc)

	assert.Nilf(t, err, "signTransactionWith returned error: %s", err)
	assert.Equal(t, secretProofTransactionSigningCorr, b.Payload)
}

func TestDeadline(t *testing.T) {
	if !time.Now().Before(NewDeadline(time.Hour * 2).Time) {
		t.Error("now is before deadline localtime")
	}

	if !time.Now().Add(time.Hour * 2).Add(-time.Second).Before(NewDeadline(time.Hour * 2).Time) {
		t.Error("now plus 2 hours is before deadline localtime")
	}

	if !time.Now().Add(time.Hour * 2).Add(time.Second * 2).After(NewDeadline(time.Hour * 2).Time) {
		t.Error("now plus 2 hours and 2 seconds is after deadline localtime")
	}
}

func TestMapTransaction_MosaicDefinitionTransaction(t *testing.T) {
	txr := `{"transaction":{"signer":"9C2086FE49B7A00578009B705AD719DB7E02A27870C67966AAA40540C136E248","version":43010,"type":16717,"parentId":[2031553063,2912841636],"mosaicId":[1477049789,645988887],"properties":[{"key":0,"value":[2,0]},{"key":1,"value":[0,0]},{"key":2,"value":[5,0]}],"name":"storage_0"}}`

	_, err := MapTransaction(bytes.NewBuffer([]byte(txr)))

	assert.Nilf(t, err, "MapTransaction returned error: %s", err)
}

func TestMapTransactions(t *testing.T) {
	txr := `[{"transaction":{"signer":"9C2086FE49B7A00578009B705AD719DB7E02A27870C67966AAA40540C136E248","version":43010,"type":16717,"parentId":[2031553063,2912841636],"mosaicId":[1477049789,645988887],"properties":[{"key":0,"value":[2,0]},{"key":1,"value":[0,0]},{"key":2,"value":[5,0]}],"name":"storage_0"}}, {"transaction":{"signer":"9C2086FE49B7A00578009B705AD719DB7E02A27870C67966AAA40540C136E248","version":43010,"type":16717,"parentId":[2031553063,2912841636],"mosaicId":[1477049789,645988887],"properties":[{"key":0,"value":[2,0]},{"key":1,"value":[0,0]},{"key":2,"value":[5,0]}],"name":"storage_0"}}]`

	txs, err := MapTransactions(bytes.NewBuffer([]byte(txr)))

	assert.Nilf(t, err, "MapTransaction returned error: %s", err)
	assert.True(t, len(txs) == 2)
}
