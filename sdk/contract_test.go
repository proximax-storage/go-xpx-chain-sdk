package sdk

import (
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	contractClient = mockServer.getPublicTestClientUnsafe().Contract
)

const (
	testContractInfoJson = `{
	"contract": {
		"multisig": "EB8923957301F796C884977234D20B0388A3AD6F865F1ACC7D3A94AFF597D59D",
		"multisigAddress": "905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA34",
		"start": [
			36,
			0
		],
		"duration": [
			11,
			0
		],
		"hash": "D8E06B597BEE34263E9C970A50B5341783EFF67EF00637644C114447BE1905DA",
		"customers": [
			"A93FEB7F051F4258C73FE0BD009F50F1E71DBA8A88B6E248ECDF560D9A9AB7C3"
		],
		"executors": [
			"8599BA6DB5B81BB69F96B88DD80A3B9EB7BBF8849CBD979100E89D69C30356E0"
		],
		"verifiers": [
			"3DCB6E5EFF4D63A38902EF948E895B01D6EA497EBF84B1460C14CA5BEDCAD9F3"
		]
	}
}`

	testContractInfoJsonArr = "[" + testContractInfoJson + "]"

	testContractInfoPubKey = "EB8923957301F796C884977234D20B0388A3AD6F865F1ACC7D3A94AFF597D59D"
)

var (
	testContractInfo = &ContractInfo{
		Multisig:        "EB8923957301F796C884977234D20B0388A3AD6F865F1ACC7D3A94AFF597D59D",
		MultisigAddress: NewAddress("905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA34", PublicTest),
		Start:           uint64DTO{36, 0}.toBigInt(),
		Duration:        uint64DTO{11, 0}.toBigInt(),
		Content:         "D8E06B597BEE34263E9C970A50B5341783EFF67EF00637644C114447BE1905DA",
		Customers:       []string{"A93FEB7F051F4258C73FE0BD009F50F1E71DBA8A88B6E248ECDF560D9A9AB7C3"},
		Executors: []string{
			"8599BA6DB5B81BB69F96B88DD80A3B9EB7BBF8849CBD979100E89D69C30356E0",
		},
		Verifiers: []string{
			"3DCB6E5EFF4D63A38902EF948E895B01D6EA497EBF84B1460C14CA5BEDCAD9F3",
		},
	}
)

func TestContractService_GetContractsInfo(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:                contractsInfoRoute,
		AcceptedHttpMethods: []string{http.MethodPost},
		RespHttpCode:        200,
		RespBody:            testContractInfoJsonArr,
		ReqJsonBodyStruct: struct {
			PublicKeys []string `json:"publicKeys"`
		}{},
	})

	infos, err := contractClient.GetContractsInfo(ctx, testContractInfoPubKey)
	assert.Nil(t, err)
	assert.NotNil(t, infos)
	assert.Equal(t, 1, len(infos))
	assert.Equal(t, testContractInfo, infos[0])
}

func TestContractService_GetContractsByAddress(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:                fmt.Sprintf(contractsByAccountRoute, "8599BA6DB5B81BB69F96B88DD80A3B9EB7BBF8849CBD979100E89D69C30356E0"),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testContractInfoJsonArr,
	})

	infos, err := contractClient.GetContractsByAddress(ctx, "8599BA6DB5B81BB69F96B88DD80A3B9EB7BBF8849CBD979100E89D69C30356E0")
	assert.Nil(t, err)
	assert.NotNil(t, infos)
	assert.Equal(t, 1, len(infos))
	assert.Equal(t, testContractInfo, infos[0])
}
