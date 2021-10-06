package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)


var testActiveDataModification = &ActiveDataModification{&Hash{1}, testDriveAccount, &Hash{2}, 12}
var testCompletedDataModification = CompletedDataModification{testActiveDataModification, Succeeded}
var testBcDrive = BcDrive{testDriveAccount, testDriveOwnerAccount, &Hash{3}, 12, 13, 14, 3, []*ActiveDataModification{testActiveDataModification}, []*CompletedDataModification{&testCompletedDataModification}}
var testDriveInfov2 = &DriveInfo{&Hash{4}, false, 1, 1}
var driveInfoMap = map[string]*DriveInfo{"CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE" : testDriveInfov2}
var testReplicator = Replicator{testReplicatorAccount, 2, Amount(10), "blskeys", driveInfoMap}

func TestActiveDataModificationString(t *testing.T) {
	expectedResult := fmt.Sprintf(
		`
			"Id": 0100000000000000000000000000000000000000000000000000000000000000,
			"Owner": Address:  [Type=168, Address=VBEHMADGUUHQ6ZMCBUYARJ44647BANFFMRKLLUPF], PublicKey: "415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309",
			"DownloadDataCdi": 0200000000000000000000000000000000000000000000000000000000000000,
			"UploadSize": 12,
		`)
	assert.Equal(t, expectedResult, testActiveDataModification.String())
}

func TestCompletedDataModificationString(t *testing.T){
	expectedResult := fmt.Sprintf(
		`
			"ActiveDataModification": 
			"Id": 0100000000000000000000000000000000000000000000000000000000000000,
			"Owner": Address:  [Type=168, Address=VBEHMADGUUHQ6ZMCBUYARJ44647BANFFMRKLLUPF], PublicKey: "415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309",
			"DownloadDataCdi": 0200000000000000000000000000000000000000000000000000000000000000,
			"UploadSize": 12,
		,
			"State:" 0,
		`)
	assert.Equal(t, expectedResult, testCompletedDataModification.String())
}

func TestBcDriveString(t *testing.T){
	expectedResult := fmt.Sprintf(
		`
		"BcDriveAccount": Address:  [Type=168, Address=VBEHMADGUUHQ6ZMCBUYARJ44647BANFFMRKLLUPF], PublicKey: "415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309",
		"OwnerAccount": Address:  [Type=168, Address=VBJIHDIXHOU5YGBCYEXYZQXKDI5YRE4XOI5EALMN], PublicKey: "CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE",
		"RootHash": 0300000000000000000000000000000000000000000000000000000000000000,
		"DriveSize": 12,
		"UsedSize": 13,
		"MetaFilesSize": 14,
		"ReplicatorCount": 3,
		"ActiveDataModifications": [
			"Id": 0100000000000000000000000000000000000000000000000000000000000000,
			"Owner": Address:  [Type=168, Address=VBEHMADGUUHQ6ZMCBUYARJ44647BANFFMRKLLUPF], PublicKey: "415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309",
			"DownloadDataCdi": 0200000000000000000000000000000000000000000000000000000000000000,
			"UploadSize": 12,
		],
		"CompletedDataModifications": [
			"ActiveDataModification": 
			"Id": 0100000000000000000000000000000000000000000000000000000000000000,
			"Owner": Address:  [Type=168, Address=VBEHMADGUUHQ6ZMCBUYARJ44647BANFFMRKLLUPF], PublicKey: "415C7C61822B063F62A4876A6F6BA2DAAE114AB298D7AC7FC56FDBA95872C309",
			"DownloadDataCdi": 0200000000000000000000000000000000000000000000000000000000000000,
			"UploadSize": 12,
		,
			"State:" 0,
		],
		`)
	assert.Equal(t, expectedResult, testBcDrive.String())
}

func TestDriveInfoString(t *testing.T){
	expectedResult := fmt.Sprintf(
		`
		    "LastApprovedDataModificationId": 0400000000000000000000000000000000000000000000000000000000000000,
			"DataModificationIdIsValid": false,
			"InitialDownloadWork": 1,
			"Index": 1
		`)
	assert.Equal(t, expectedResult, testDriveInfov2.String())
}

func TestReplicatorString(t *testing.T){
	expectedResult := fmt.Sprintf(
		`
		ReplicatorAccount: Address:  [Type=168, Address=VDQPWCXBJRL5JW4VAWVWGWXDCLREERSI24PBUSUL], PublicKey: "36E7F50C8B8BC9A4FC6325B2359E0E5DB50C75A914B5292AD726FD5AE3992691", 
		Version: 2,
		Capacity: 10,
		BLSKey: blskeys,
		Drives: map[CFC31B3080B36BC3D59DF4AB936AC72F4DC15CE3C3E1B1EC5EA41415A4C33FEE:
		    "LastApprovedDataModificationId": 0400000000000000000000000000000000000000000000000000000000000000,
			"DataModificationIdIsValid": false,
			"InitialDownloadWork": 1,
			"Index": 1
		],
		`)
	assert.Equal(t, expectedResult, testReplicator.String())
}