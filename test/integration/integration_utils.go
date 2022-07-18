package integration

import (
	"context"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"testing"
)

func GetEntityIsSupportedAtVersion(client *sdk.Client, entityType sdk.EntityType, version sdk.EntityVersion) bool {
	config, err := client.Network.GetNetworkConfig(context.Background())
	if err != nil {
		return false
	}
	val, ok := config.SupportedEntityVersions.Entities[entityType]
	if ok {
		for _, s := range val.SupportedVersions {
			if s == version {
				return true
			}
		}
	}
	return false

}

func SkipIfEntityNotSupportedAtVersion(client *sdk.Client, t *testing.T, entityType sdk.EntityType, version sdk.EntityVersion) {
	if !GetEntityIsSupportedAtVersion(client, entityType, version) {
		t.SkipNow()
	}
}
