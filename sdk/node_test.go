package sdk

import (
	"net/http"
	"testing"

	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
)

const (
	testNodeInfoJson = `{
  "publicKey": "460458B98E2BAA36A8E95DE9B320379E89898885B71CF0174E02F1324FAFFAC1",
  "port": 7900,
  "networkIdentifier": 168,
  "version": 0,
  "roles": 2,
  "host": "catapult-api-node",
  "friendlyName": "api-node-0"
}`
	testNodeTimeJson = `{
  "communicationTimestamps": {
    "sendTimestamp": [
      3905031395,
      33
    ],
    "receiveTimestamp": [
      3905031395,
      33
    ]
  }
}`

	testNodeInfoJsonArr = "[" + testNodeInfoJson + ", " + testNodeInfoJson + "]"
)

var testPublicKey, _ = NewAccountFromPublicKey("460458B98E2BAA36A8E95DE9B320379E89898885B71CF0174E02F1324FAFFAC1", PublicTest)
var testNodeTime = NewBlockchainTimestamp(int64(uint64DTO{3905031395, 33}.toUint64()))

var (
	testNodeInfo = &NodeInfo{
		Account:      testPublicKey,
		Port:         7900,
		NetworkType:  PublicTest,
		Roles:        2,
		Host:         "catapult-api-node",
		FriendlyName: "api-node-0",
	}
)

func TestNodeService_GetNodeInfo(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                nodeInfoRoute,
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testNodeInfoJson,
	})
	nodeClient := mock.getPublicTestClientUnsafe().Node

	defer mock.Close()

	info, err := nodeClient.GetNodeInfo(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, testNodeInfo, info)
}

func TestNodeService_GetNodeTime(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                nodeTimeRoute,
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testNodeTimeJson,
	})
	nodeClient := mock.getPublicTestClientUnsafe().Node

	defer mock.Close()

	nodeTime, err := nodeClient.GetNodeTime(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, nodeTime)
	assert.Equal(t, testNodeTime, nodeTime)
}

func TestNodeService_GetNodePeers(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                nodePeersRoute,
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testNodeInfoJsonArr,
	})
	nodeClient := mock.getPublicTestClientUnsafe().Node

	defer mock.Close()

	infos, err := nodeClient.GetNodePeers(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, infos)
	assert.Equal(t, len(infos), 2)
	assert.Equal(t, []*NodeInfo{testNodeInfo, testNodeInfo}, infos)
}
