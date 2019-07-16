package sdk

import (
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	testAddressMetadataInfoJson = `{
	"metadata": {
		"metadataId": "90936FF3536858CBEA8EE0EAAB99FE9EC4EF5EF1F66366569A",
		"metadataType": 1,
    	"fields": [
      		{
        		"key": "jora229",
        		"value": "I Love you"
      		}
    	]
  	}
}`
	testMosaicMetadataInfoJson = `{
	"metadata": {
		"metadataId": [
			1942264427,
			1236639028
    	],
		"metadataType": 2,
    	"fields": [
      		{
        		"key": "hello",
        		"value": "hell"
      		}
    	]
  	}
}`
	testNamespaceMetadataInfoJson = `{
	"metadata": {
		"metadataId": [
      		94107370,
      		2911018100
    	],
		"metadataType": 3,
    	"fields": [
      		{
        		"key": "hello",
        		"value": "world"
      		}
    	]
  	}
}`
	testNotFoundInfoJson = `{
	"code": "ResourceNotFound",
	"message": "no resource exists with id '4829DDBFBD2121A007E13AF98558092D193312C10B26CA95CD147FB40CEE3E68'"
}`

	testAddressMetadataInfoJsonArr   = "[" + testAddressMetadataInfoJson + "]"
	testMosaicMetadataInfoJsonArr    = "[" + testMosaicMetadataInfoJson + "]"
	testNamespaceMetadataInfoJsonArr = "[" + testNamespaceMetadataInfoJson + "]"
)

var (
	testMetadataInfoPubKey         = "SCJW742TNBMMX2UO4DVKXGP6T3CO6XXR6ZRWMVU2"
	testMetadataAddress, _         = NewAddressFromRaw(testMetadataInfoPubKey)
	testMetadataInfoMosaicId, _    = NewMosaicId(0x49B59D3473C49A6B)
	testMetadataInfoNamespaceId, _ = NewNamespaceId(0xAD829C74059BF6EA)

	testAddressMetadataInfo = &AddressMetadataInfo{
		MetadataInfo: MetadataInfo{
			MetadataType: MetadataAddressType,
			Fields: map[string]string{
				"jora229": "I Love you",
			},
		},
		Address: testMetadataAddress,
	}

	testMosaicMetadataInfo = &MosaicMetadataInfo{
		MetadataInfo: MetadataInfo{
			MetadataType: MetadataMosaicType,
			Fields: map[string]string{
				"hello": "hell",
			},
		},
		MosaicId: testMetadataInfoMosaicId,
	}

	testNamespaceMetadataInfo = &NamespaceMetadataInfo{
		MetadataInfo: MetadataInfo{
			MetadataType: MetadataNamespaceType,
			Fields: map[string]string{
				"hello": "world",
			},
		},
		NamespaceId: testMetadataInfoNamespaceId,
	}

	testAddressMetadataInfo_NotFound = &AddressMetadataInfo{
		MetadataInfo: MetadataInfo{
			MetadataType: MetadataAddressType,
			Fields:       make(map[string]string),
		},
		Address: testMetadataAddress,
	}

	testMosaicMetadataInfo_NotFound = &MosaicMetadataInfo{
		MetadataInfo: MetadataInfo{
			MetadataType: MetadataMosaicType,
			Fields:       make(map[string]string),
		},
		MosaicId: testMetadataInfoMosaicId,
	}

	testNamespaceMetadataInfo_NotFound = &NamespaceMetadataInfo{
		MetadataInfo: MetadataInfo{
			MetadataType: MetadataNamespaceType,
			Fields:       make(map[string]string),
		},
		NamespaceId: testMetadataInfoNamespaceId,
	}
)

func TestMetadataService_GetAddressMetadatasInfo(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                metadatasInfoRoute,
		AcceptedHttpMethods: []string{http.MethodPost},
		RespHttpCode:        200,
		RespBody:            testAddressMetadataInfoJsonArr,
		ReqJsonBodyStruct: struct {
			Addresses []string `json:"metadataIds"`
		}{},
	})
	metadataClient := mock.getPublicTestClientUnsafe().Metadata

	defer mock.Close()

	infos, err := metadataClient.GetAddressMetadatasInfo(ctx, testMetadataInfoPubKey)
	assert.Nil(t, err)
	assert.NotNil(t, infos)
	assert.Equal(t, 1, len(infos))
	assert.Equal(t, testAddressMetadataInfo, infos[0])
}

func TestMetadataService_GetMosaicMetadatasInfo(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                metadatasInfoRoute,
		AcceptedHttpMethods: []string{http.MethodPost},
		RespHttpCode:        200,
		RespBody:            testMosaicMetadataInfoJsonArr,
		ReqJsonBodyStruct: struct {
			MosaicIds []string `json:"metadataIds"`
		}{},
	})
	metadataClient := mock.getPublicTestClientUnsafe().Metadata

	defer mock.Close()

	infos, err := metadataClient.GetMosaicMetadatasInfo(ctx, testMetadataInfoMosaicId)
	assert.Nil(t, err)
	assert.NotNil(t, infos)
	assert.Equal(t, 1, len(infos))
	assert.Equal(t, testMosaicMetadataInfo, infos[0])
}

func TestMetadataService_GetNamespaceMetadatasInfo(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                metadatasInfoRoute,
		AcceptedHttpMethods: []string{http.MethodPost},
		RespHttpCode:        200,
		RespBody:            testNamespaceMetadataInfoJsonArr,
		ReqJsonBodyStruct: struct {
			NamespaceIds []string `json:"metadataIds"`
		}{},
	})
	metadataClient := mock.getPublicTestClientUnsafe().Metadata

	defer mock.Close()

	infos, err := metadataClient.GetNamespaceMetadatasInfo(ctx, testMetadataInfoNamespaceId)
	assert.Nil(t, err)
	assert.NotNil(t, infos)
	assert.Equal(t, 1, len(infos))
	assert.Equal(t, testNamespaceMetadataInfo, infos[0])
}

func TestMetadataService_GetMetadatasByAddress(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(metadataByAccountRoute, testMetadataInfoPubKey),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testAddressMetadataInfoJson,
	})
	metadataClient := mock.getPublicTestClientUnsafe().Metadata

	defer mock.Close()

	info, err := metadataClient.GetMetadataByAddress(ctx, testMetadataInfoPubKey)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, testAddressMetadataInfo, info)
}

func TestMetadataService_GetMetadatasByMosaicId(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(metadataByMosaicRoute, testMetadataInfoMosaicId.toHexString()),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testMosaicMetadataInfoJson,
	})
	metadataClient := mock.getPublicTestClientUnsafe().Metadata

	defer mock.Close()

	info, err := metadataClient.GetMetadataByMosaicId(ctx, testMetadataInfoMosaicId)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, testMosaicMetadataInfo, info)
}

func TestMetadataService_GetMetadatasByNamespaceId(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(metadataByNamespaceRoute, testMetadataInfoNamespaceId.toHexString()),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testNamespaceMetadataInfoJson,
	})
	metadataClient := mock.getPublicTestClientUnsafe().Metadata

	defer mock.Close()

	info, err := metadataClient.GetMetadataByNamespaceId(ctx, testMetadataInfoNamespaceId)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, testNamespaceMetadataInfo, info)
}

func TestMetadataService_GetMetadatasByAddress_NotFound(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(metadataByAccountRoute, testMetadataInfoPubKey),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        404,
		RespBody:            testNotFoundInfoJson,
	})
	metadataClient := mock.getPublicTestClientUnsafe().Metadata

	defer mock.Close()

	info, err := metadataClient.GetMetadataByAddress(ctx, testMetadataInfoPubKey)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, testAddressMetadataInfo_NotFound, info)
}

func TestMetadataService_GetMetadatasByMosaicId_NotFound(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(metadataByMosaicRoute, testMetadataInfoMosaicId.toHexString()),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        404,
		RespBody:            testNotFoundInfoJson,
	})
	metadataClient := mock.getPublicTestClientUnsafe().Metadata

	defer mock.Close()

	info, err := metadataClient.GetMetadataByMosaicId(ctx, testMetadataInfoMosaicId)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, testMosaicMetadataInfo_NotFound, info)
}

func TestMetadataService_GetMetadatasByNamespaceId_NotFound(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(metadataByNamespaceRoute, testMetadataInfoNamespaceId.toHexString()),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        404,
		RespBody:            testNotFoundInfoJson,
	})
	metadataClient := mock.getPublicTestClientUnsafe().Metadata

	defer mock.Close()

	info, err := metadataClient.GetMetadataByNamespaceId(ctx, testMetadataInfoNamespaceId)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, testNamespaceMetadataInfo_NotFound, info)
}
