// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	prepareDriveTransactionSerializationCorr = []byte{0xbf, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x90, 0x5a, 0x41, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xba, 0xfd, 0x56, 0x0, 0x0, 0x0, 0x0, 0xce, 0x2, 0x70, 0x4f, 0xab, 0x5, 0xd5, 0xbb, 0x89, 0x81, 0x39, 0x7d, 0x56, 0x3f, 0x7c, 0xc9, 0x31, 0x28, 0x64, 0xf0, 0x65, 0x19, 0x5e, 0x3d, 0xb8, 0xe0, 0x84, 0x7d, 0x9f, 0x20, 0x79, 0x21, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x64}

	prepareDriveTransactionToAggregateCorr = []byte{0x6f, 0x0, 0x0, 0x0, 0xce, 0x2, 0x70, 0x4f, 0xab, 0x5, 0xd5, 0xbb, 0x89, 0x81, 0x39, 0x7d, 0x56, 0x3f, 0x7c, 0xc9, 0x31, 0x28, 0x64, 0xf0, 0x65, 0x19, 0x5e, 0x3d, 0xb8, 0xe0, 0x84, 0x7d, 0x9f, 0x20, 0x79, 0x21, 0x3, 0x0, 0x0, 0x90, 0x5a, 0x41, 0xce, 0x2, 0x70, 0x4f, 0xab, 0x5, 0xd5, 0xbb, 0x89, 0x81, 0x39, 0x7d, 0x56, 0x3f, 0x7c, 0xc9, 0x31, 0x28, 0x64, 0xf0, 0x65, 0x19, 0x5e, 0x3d, 0xb8, 0xe0, 0x84, 0x7d, 0x9f, 0x20, 0x79, 0x21, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x64}

	prepareDriveTransactionSigningCorr = "BF000000D9A4C6ED8691BEEF8545FF14D1033020DA1E3D5CFE461D5F468695EA1DCE67AC2B8D1564AD1FD2A88DB34BC52F2F143844512E7416D58857786492850309390BCE02704FAB05D5BB8981397D563F7CC9312864F065195E3DB8E0847D9F207921030000905A41000000000000000000BAFD5600000000CE02704FAB05D5BB8981397D563F7CC9312864F065195E3DB8E0847D9F20792101000000000000000100000000000000010000000000000001000000000000000100010064"

	joinToDriveTransactionSerializationCorr = []byte{0x9a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x42, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xba, 0xfd, 0x56, 0x0, 0x0, 0x0, 0x0, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27}

	joinToDriveTransactionToAggregateCorr = []byte{0x4a, 0x0, 0x0, 0x0, 0xce, 0x2, 0x70, 0x4f, 0xab, 0x5, 0xd5, 0xbb, 0x89, 0x81, 0x39, 0x7d, 0x56, 0x3f, 0x7c, 0xc9, 0x31, 0x28, 0x64, 0xf0, 0x65, 0x19, 0x5e, 0x3d, 0xb8, 0xe0, 0x84, 0x7d, 0x9f, 0x20, 0x79, 0x21, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x42, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27}

	joinToDriveTransactionSigningCorr = "9A000000CF1075D3043255982F89529C20D101CF51BAFD7BAB167BA6D468F62EF282170A64E8BBA97BDC48F553B936A382417F7094E96E1D4E79D366E83ED65956EA4E0CCE02704FAB05D5BB8981397D563F7CC9312864F065195E3DB8E0847D9F207921010000905A42000000000000000000BAFD5600000000FC5CDB2478117F48BA0C1687178EF69C016BBA89D342C1ACD155E1F5AE0F4727"

	driveFileSystemTransactionSerializationCorr = []byte{0x2e, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x43, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xba, 0xfd, 0x56, 0x0, 0x0, 0x0, 0x0, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	driveFileSystemTransactionToAggregateCorr = []byte{0xde, 0x0, 0x0, 0x0, 0xce, 0x2, 0x70, 0x4f, 0xab, 0x5, 0xd5, 0xbb, 0x89, 0x81, 0x39, 0x7d, 0x56, 0x3f, 0x7c, 0xc9, 0x31, 0x28, 0x64, 0xf0, 0x65, 0x19, 0x5e, 0x3d, 0xb8, 0xe0, 0x84, 0x7d, 0x9f, 0x20, 0x79, 0x21, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x43, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	driveFileSystemTransactionSigningCorr = "2E010000A3D39A93273EAC2FF0A7EB1017F53FDE702D9D958E4A8C9FA1212A65CF3FD9340B629F58BA9B6F06C8EA4F54A84CBA5BD8AA268B02248B984138E469D0292907CE02704FAB05D5BB8981397D563F7CC9312864F065195E3DB8E0847D9F207921010000905A43000000000000000000BAFD5600000000FC5CDB2478117F48BA0C1687178EF69C016BBA89D342C1ACD155E1F5AE0F472700000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010001000300000000000000000000000000000000000000000000000000000000000000080000000000000004000000000000000000000000000000000000000000000000000000000000000900000000000000"

	filesDepositTransactionSerializationCorr = []byte{0xbc, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x44, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xba, 0xfd, 0x56, 0x0, 0x0, 0x0, 0x0, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27, 0x1, 0x0, 0xaa, 0x2d, 0x24, 0x27, 0xe1, 0x5, 0xa9, 0xb6, 0xd, 0xf6, 0x34, 0x55, 0x38, 0x49, 0x13, 0x5d, 0xf6, 0x29, 0xf1, 0x40, 0x8a, 0x1, 0x8d, 0x2, 0xb0, 0x7a, 0x70, 0xca, 0xff, 0xb4, 0x30, 0x93}

	filesDepositTransactionToAggregateCorr = []byte{0x6c, 0x0, 0x0, 0x0, 0xce, 0x2, 0x70, 0x4f, 0xab, 0x5, 0xd5, 0xbb, 0x89, 0x81, 0x39, 0x7d, 0x56, 0x3f, 0x7c, 0xc9, 0x31, 0x28, 0x64, 0xf0, 0x65, 0x19, 0x5e, 0x3d, 0xb8, 0xe0, 0x84, 0x7d, 0x9f, 0x20, 0x79, 0x21, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x44, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27, 0x1, 0x0, 0xaa, 0x2d, 0x24, 0x27, 0xe1, 0x5, 0xa9, 0xb6, 0xd, 0xf6, 0x34, 0x55, 0x38, 0x49, 0x13, 0x5d, 0xf6, 0x29, 0xf1, 0x40, 0x8a, 0x1, 0x8d, 0x2, 0xb0, 0x7a, 0x70, 0xca, 0xff, 0xb4, 0x30, 0x93}

	filesDepositTransactionSigningCorr = "BC000000AED8F6293E6DCFC2A8E6EEC652B30ACBBDB72249C4AF05C4AE748787167E01AA992589C4251D57438641D9B8F32B7E9D78C8D74D02E506A3D94744C6A8B7E50CCE02704FAB05D5BB8981397D563F7CC9312864F065195E3DB8E0847D9F207921010000905A44000000000000000000BAFD5600000000FC5CDB2478117F48BA0C1687178EF69C016BBA89D342C1ACD155E1F5AE0F47270100AA2D2427E105A9B60DF634553849135DF629F1408A018D02B07A70CAFFB43093"

	endDriveTransactionSerializationCorr = []byte{0x9a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x45, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xba, 0xfd, 0x56, 0x0, 0x0, 0x0, 0x0, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27}

	endDriveTransactionToAggregateCorr = []byte{0x4a, 0x0, 0x0, 0x0, 0xce, 0x2, 0x70, 0x4f, 0xab, 0x5, 0xd5, 0xbb, 0x89, 0x81, 0x39, 0x7d, 0x56, 0x3f, 0x7c, 0xc9, 0x31, 0x28, 0x64, 0xf0, 0x65, 0x19, 0x5e, 0x3d, 0xb8, 0xe0, 0x84, 0x7d, 0x9f, 0x20, 0x79, 0x21, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x45, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27}

	endDriveTransactionSigningCorr = "9A0000003AF2763FAED4E5EF6BC8ADB08ECCE1C7882A5AFA702CB0A747FD8968A058AD10C79D4FD53AF51FD8D6EA87B2742E1C5019FCF645DCDAFF0D9EB5A7984F145501CE02704FAB05D5BB8981397D563F7CC9312864F065195E3DB8E0847D9F207921010000905A45000000000000000000BAFD5600000000FC5CDB2478117F48BA0C1687178EF69C016BBA89D342C1ACD155E1F5AE0F4727"

	driveFilesRewardTransactionSerializationCorr = []byte{0xa4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xba, 0xfd, 0x56, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27, 0x63, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	driveFilesRewardTransactionToAggregateCorr = []byte{0x54, 0x0, 0x0, 0x0, 0xce, 0x2, 0x70, 0x4f, 0xab, 0x5, 0xd5, 0xbb, 0x89, 0x81, 0x39, 0x7d, 0x56, 0x3f, 0x7c, 0xc9, 0x31, 0x28, 0x64, 0xf0, 0x65, 0x19, 0x5e, 0x3d, 0xb8, 0xe0, 0x84, 0x7d, 0x9f, 0x20, 0x79, 0x21, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x46, 0x1, 0x0, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27, 0x63, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	driveFilesRewardTransactionSigningCorr = "A4000000C9F10FA29AFE3064A917993DE84B4C81EB34F8AB6D2907EFC11A24FBE82C64D139383032B234B79C752D75F027BA1B4965B061F90C5CECB9DA7F4D0734EE1603CE02704FAB05D5BB8981397D563F7CC9312864F065195E3DB8E0847D9F207921010000905A46000000000000000000BAFD56000000000100FC5CDB2478117F48BA0C1687178EF69C016BBA89D342C1ACD155E1F5AE0F47276300000000000000"

	startDriveVerificationTransactionSerializationCorr = []byte{0x9a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x47, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xba, 0xfd, 0x56, 0x0, 0x0, 0x0, 0x0, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27}

	startDriveVerificationTransactionToAggregateCorr = []byte{0x4a, 0x0, 0x0, 0x0, 0xce, 0x2, 0x70, 0x4f, 0xab, 0x5, 0xd5, 0xbb, 0x89, 0x81, 0x39, 0x7d, 0x56, 0x3f, 0x7c, 0xc9, 0x31, 0x28, 0x64, 0xf0, 0x65, 0x19, 0x5e, 0x3d, 0xb8, 0xe0, 0x84, 0x7d, 0x9f, 0x20, 0x79, 0x21, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x47, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27}

	startDriveVerificationTransactionSigningCorr = "9A0000007DA8D1C763F3D5538B10E5CC026041281D8FFDED44A987B576AC6272756DD1F49364F2682B4C805EC77BBF6789C2579EA14A2F9D04C856647ECA2D9CCD5B830BCE02704FAB05D5BB8981397D563F7CC9312864F065195E3DB8E0847D9F207921010000905A47000000000000000000BAFD5600000000FC5CDB2478117F48BA0C1687178EF69C016BBA89D342C1ACD155E1F5AE0F4727"

	endDriveVerificationTransactionSerializationCorr = []byte{0xbe, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x48, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xba, 0xfd, 0x56, 0x0, 0x0, 0x0, 0x0, 0x44, 0x0, 0x0, 0x0, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	endDriveVerificationTransactionToAggregateCorr = []byte{0x6e, 0x0, 0x0, 0x0, 0xce, 0x2, 0x70, 0x4f, 0xab, 0x5, 0xd5, 0xbb, 0x89, 0x81, 0x39, 0x7d, 0x56, 0x3f, 0x7c, 0xc9, 0x31, 0x28, 0x64, 0xf0, 0x65, 0x19, 0x5e, 0x3d, 0xb8, 0xe0, 0x84, 0x7d, 0x9f, 0x20, 0x79, 0x21, 0x1, 0x0, 0x0, 0x90, 0x5a, 0x48, 0x44, 0x0, 0x0, 0x0, 0xfc, 0x5c, 0xdb, 0x24, 0x78, 0x11, 0x7f, 0x48, 0xba, 0xc, 0x16, 0x87, 0x17, 0x8e, 0xf6, 0x9c, 0x1, 0x6b, 0xba, 0x89, 0xd3, 0x42, 0xc1, 0xac, 0xd1, 0x55, 0xe1, 0xf5, 0xae, 0xf, 0x47, 0x27, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	endDriveVerificationTransactionSigningCorr = "BE00000080286C1006AC01CF5C508DE6EFEC26AE810A835EA2929FDDA8F68A4D4B8F3069FCCE660F2282F44028760FB0BC4964964211ABAF4C1659E16CC46AF34B6A0101CE02704FAB05D5BB8981397D563F7CC9312864F065195E3DB8E0847D9F207921010000905A48000000000000000000BAFD560000000044000000FC5CDB2478117F48BA0C1687178EF69C016BBA89D342C1ACD155E1F5AE0F47277B00000000000000000000000000000000000000000000000000000000000000"
)

var testDriveOwner, _ = NewAccountFromPrivateKey("9A49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24", MijinTest, &Hash{})
var testDrive, _ = NewAccountFromPrivateKey("AA49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24", MijinTest, &Hash{})

func TestPrepareDriveTransactionSerialization(t *testing.T) {
	tx, err := NewPrepareDriveTransaction(
		fakeDeadline,
		testDriveOwner.PublicAccount,
		Duration(1),
		Duration(1),
		Amount(1),
		StorageSize(1),
		uint16(1),
		uint16(1),
		uint8(100),
		MijinTest,
	)
	assert.Nilf(t, err, "NewPrepareDriveTransaction returned error: %s", err)

	b, err := tx.Bytes()

	assert.Nilf(t, err, "PrepareDriveTransaction.Bytes returned error: %s", err)
	assert.Equal(t, prepareDriveTransactionSerializationCorr, b)
}

func TestPrepareDriveTransactionToAggregate(t *testing.T) {
	tx, err := NewPrepareDriveTransaction(
		fakeDeadline,
		testDriveOwner.PublicAccount,
		Duration(1),
		Duration(1),
		Amount(1),
		StorageSize(1),
		uint16(1),
		uint16(1),
		uint8(100),
		MijinTest,
	)
	assert.Nilf(t, err, "NewPrepareDriveTransaction returned error: %s", err)
	tx.Signer = testDriveOwner.PublicAccount

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, prepareDriveTransactionToAggregateCorr, b)
}

func TestPrepareDriveTransactionSigning(t *testing.T) {
	tx, err := NewPrepareDriveTransaction(
		fakeDeadline,
		testDriveOwner.PublicAccount,
		Duration(1),
		Duration(1),
		Amount(1),
		StorageSize(1),
		uint16(1),
		uint16(1),
		uint8(100),
		MijinTest,
	)
	assert.Nilf(t, err, "NewPrepareDriveTransaction returned error: %s", err)

	b, err := testDriveOwner.Sign(tx)

	assert.Nilf(t, err, "Sign returned error: %s", err)
	assert.Equal(t, prepareDriveTransactionSigningCorr, b.Payload)
	assert.Equal(t, stringToHashPanic("c9d1bff05ea5500512d838ce8b551a3c2b77a4d4d367b430490bef9ff983d3e1"), b.Hash)
}

func TestJoinToDriveTransactionSerialization(t *testing.T) {
	tx, err := NewJoinToDriveTransaction(
		fakeDeadline,
		testDrive.PublicAccount,
		MijinTest,
	)
	assert.Nilf(t, err, "NewJoinToDriveTransaction returned error: %s", err)

	b, err := tx.Bytes()

	assert.Nilf(t, err, "JoinToDriveTransaction.Bytes returned error: %s", err)
	assert.Equal(t, joinToDriveTransactionSerializationCorr, b)
}

func TestJoinToDriveTransactionToAggregate(t *testing.T) {
	tx, err := NewJoinToDriveTransaction(
		fakeDeadline,
		testDrive.PublicAccount,
		MijinTest,
	)
	assert.Nilf(t, err, "NewJoinToDriveTransaction returned error: %s", err)
	tx.Signer = testDriveOwner.PublicAccount

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, joinToDriveTransactionToAggregateCorr, b)
}

func TestJoinToDriveTransactionSigning(t *testing.T) {
	tx, err := NewJoinToDriveTransaction(
		fakeDeadline,
		testDrive.PublicAccount,
		MijinTest,
	)
	assert.Nilf(t, err, "NewJoinToDriveTransaction returned error: %s", err)
	assert.Nilf(t, err, "NewJoinToDriveTransaction returned error: %s", err)

	b, err := testDriveOwner.Sign(tx)

	assert.Nilf(t, err, "Sign returned error: %s", err)
	assert.Equal(t, joinToDriveTransactionSigningCorr, b.Payload)
	assert.Equal(t, stringToHashPanic("23b0b1570a69208379687ac48ed4d93520393b81645724f250fc172a8be14bf7"), b.Hash)
}

func TestDriveFileSystemTransactionSerialization(t *testing.T) {
	tx, err := NewDriveFileSystemTransaction(
		fakeDeadline,
		testDrive.PublicAccount.PublicKey,
		&Hash{0},
		&Hash{1},
		[]*Action{
			{
				FileSize: 8,
				FileHash: &Hash{3},
			},
		},
		[]*Action{
			{
				FileSize: 9,
				FileHash: &Hash{4},
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewDriveFileSystemTransaction returned error: %s", err)

	b, err := tx.Bytes()

	assert.Nilf(t, err, "DriveFileSystemTransaction.Bytes returned error: %s", err)
	assert.Equal(t, driveFileSystemTransactionSerializationCorr, b)
}

func TestDriveFileSystemTransactionToAggregate(t *testing.T) {
	tx, err := NewDriveFileSystemTransaction(
		fakeDeadline,
		testDrive.PublicAccount.PublicKey,
		&Hash{0},
		&Hash{1},
		[]*Action{
			{
				FileSize: 8,
				FileHash: &Hash{3},
			},
		},
		[]*Action{
			{
				FileSize: 9,
				FileHash: &Hash{4},
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewDriveFileSystemTransaction returned error: %s", err)
	tx.Signer = testDriveOwner.PublicAccount

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, driveFileSystemTransactionToAggregateCorr, b)
}

func TestDriveFileSystemTransactionSigning(t *testing.T) {
	tx, err := NewDriveFileSystemTransaction(
		fakeDeadline,
		testDrive.PublicAccount.PublicKey,
		&Hash{0},
		&Hash{1},
		[]*Action{
			{
				FileSize: 8,
				FileHash: &Hash{3},
			},
		},
		[]*Action{
			{
				FileSize: 9,
				FileHash: &Hash{4},
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewDriveFileSystemTransaction returned error: %s", err)

	b, err := testDriveOwner.Sign(tx)

	assert.Nilf(t, err, "Sign returned error: %s", err)
	assert.Equal(t, driveFileSystemTransactionSigningCorr, b.Payload)
	assert.Equal(t, stringToHashPanic("c44b25527afeffb67d35c19f37858d922f2926620b7b903409c4f0b9bf8d0952"), b.Hash)
}

func TestFilesDepositTransactionSerialization(t *testing.T) {
	tx, err := NewFilesDepositTransaction(
		fakeDeadline,
		testDrive.PublicAccount,
		[]*File{
			{
				FileHash: testFileHash,
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewFilesDepositTransaction returned error: %s", err)

	b, err := tx.Bytes()

	assert.Nilf(t, err, "FilesDepositTransaction.Bytes returned error: %s", err)
	assert.Equal(t, filesDepositTransactionSerializationCorr, b)
}

func TestFilesDepositTransactionToAggregate(t *testing.T) {
	tx, err := NewFilesDepositTransaction(
		fakeDeadline,
		testDrive.PublicAccount,
		[]*File{
			{
				FileHash: testFileHash,
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewFilesDepositTransaction returned error: %s", err)
	tx.Signer = testDriveOwner.PublicAccount

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, filesDepositTransactionToAggregateCorr, b)
}

func TestFilesDepositTransactionSigning(t *testing.T) {
	tx, err := NewFilesDepositTransaction(
		fakeDeadline,
		testDrive.PublicAccount,
		[]*File{
			{
				FileHash: testFileHash,
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewFilesDepositTransaction returned error: %s", err)

	b, err := testDriveOwner.Sign(tx)

	assert.Nilf(t, err, "Sign returned error: %s", err)
	assert.Equal(t, filesDepositTransactionSigningCorr, b.Payload)
	assert.Equal(t, stringToHashPanic("78a194caa9ce4985bf28f0c0adcfb813ba199c7ef5918a344ab4899d904183a0"), b.Hash)
}

func TestEndDriveTransactionSerialization(t *testing.T) {
	tx, err := NewEndDriveTransaction(
		fakeDeadline,
		testDrive.PublicAccount,
		MijinTest,
	)
	assert.Nilf(t, err, "NewEndDriveTransaction returned error: %s", err)

	b, err := tx.Bytes()

	assert.Nilf(t, err, "EndDriveTransaction.Bytes returned error: %s", err)
	assert.Equal(t, endDriveTransactionSerializationCorr, b)
}

func TestEndDriveTransactionToAggregate(t *testing.T) {
	tx, err := NewEndDriveTransaction(
		fakeDeadline,
		testDrive.PublicAccount,
		MijinTest,
	)
	assert.Nilf(t, err, "NewEndDriveTransaction returned error: %s", err)
	tx.Signer = testDriveOwner.PublicAccount

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, endDriveTransactionToAggregateCorr, b)
}

func TestEndDriveTransactionSigning(t *testing.T) {
	tx, err := NewEndDriveTransaction(
		fakeDeadline,
		testDrive.PublicAccount,
		MijinTest,
	)
	assert.Nilf(t, err, "NewEndDriveTransaction returned error: %s", err)

	b, err := testDriveOwner.Sign(tx)

	assert.Nilf(t, err, "Sign returned error: %s", err)
	assert.Equal(t, endDriveTransactionSigningCorr, b.Payload)
	assert.Equal(t, stringToHashPanic("413563292e36107832d7d6a0f26a885818d344c1e5f806b4c9903ecbca0e9448"), b.Hash)
}

func TestDriveFilesRewardTransactionSerialization(t *testing.T) {
	tx, err := NewDriveFilesRewardTransaction(
		fakeDeadline,
		[]*UploadInfo{
			{
				Participant:  testDrive.PublicAccount,
				UploadedSize: 99,
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewDriveFilesRewardTransaction returned error: %s", err)

	b, err := tx.Bytes()

	assert.Nilf(t, err, "DriveFilesRewardTransaction.Bytes returned error: %s", err)
	assert.Equal(t, driveFilesRewardTransactionSerializationCorr, b)
}

func TestDriveFilesRewardTransactionToAggregate(t *testing.T) {
	tx, err := NewDriveFilesRewardTransaction(
		fakeDeadline,
		[]*UploadInfo{
			{
				Participant:  testDrive.PublicAccount,
				UploadedSize: 99,
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewDriveFilesRewardTransaction returned error: %s", err)
	tx.Signer = testDriveOwner.PublicAccount

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, driveFilesRewardTransactionToAggregateCorr, b)
}

func TestDriveFilesRewardTransactionSigning(t *testing.T) {
	tx, err := NewDriveFilesRewardTransaction(
		fakeDeadline,
		[]*UploadInfo{
			{
				Participant:  testDrive.PublicAccount,
				UploadedSize: 99,
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewDriveFilesRewardTransaction returned error: %s", err)

	b, err := testDriveOwner.Sign(tx)

	assert.Nilf(t, err, "Sign returned error: %s", err)
	assert.Equal(t, driveFilesRewardTransactionSigningCorr, b.Payload)
	assert.Equal(t, stringToHashPanic("c41abc296ded3421ff42611bd48a2abc31305b6bb10d9f70c62e76dcc960a261"), b.Hash)
}

func TestStartDriveVerificationTransactionSerialization(t *testing.T) {
	tx, err := NewStartDriveVerificationTransaction(
		fakeDeadline,
		testDrive.PublicAccount,
		MijinTest,
	)
	assert.Nilf(t, err, "NewStartDriveVerificationTransaction returned error: %s", err)

	b, err := tx.Bytes()

	assert.Nilf(t, err, "StartDriveVerificationTransaction.Bytes returned error: %s", err)
	assert.Equal(t, startDriveVerificationTransactionSerializationCorr, b)
}

func TestStartDriveVerificationTransactionToAggregate(t *testing.T) {
	tx, err := NewStartDriveVerificationTransaction(
		fakeDeadline,
		testDrive.PublicAccount,
		MijinTest,
	)
	assert.Nilf(t, err, "NewStartDriveVerificationTransaction returned error: %s", err)
	tx.Signer = testDriveOwner.PublicAccount

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, startDriveVerificationTransactionToAggregateCorr, b)
}

func TestStartDriveVerificationTransactionSigning(t *testing.T) {
	tx, err := NewStartDriveVerificationTransaction(
		fakeDeadline,
		testDrive.PublicAccount,
		MijinTest,
	)
	assert.Nilf(t, err, "NewStartDriveVerificationTransaction returned error: %s", err)

	b, err := testDriveOwner.Sign(tx)

	assert.Nilf(t, err, "Sign returned error: %s", err)
	assert.Equal(t, startDriveVerificationTransactionSigningCorr, b.Payload)
	assert.Equal(t, stringToHashPanic("9130e4ce8c998951adf2b4459a96a9cd9fb902780a5ee651d3c623ab55dada8b"), b.Hash)
}

func TestEndDriveVerificationTransactionSerialization(t *testing.T) {
	tx, err := NewEndDriveVerificationTransaction(
		fakeDeadline,
		[]*FailureVerification{
			{
				Replicator:  testDrive.PublicAccount,
				BlochHashes: []*Hash{{123}},
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewEndDriveVerificationTransaction returned error: %s", err)

	b, err := tx.Bytes()

	assert.Nilf(t, err, "EndDriveVerificationTransaction.Bytes returned error: %s", err)
	assert.Equal(t, endDriveVerificationTransactionSerializationCorr, b)
}

func TestEndDriveVerificationTransactionToAggregate(t *testing.T) {
	tx, err := NewEndDriveVerificationTransaction(
		fakeDeadline,
		[]*FailureVerification{
			{
				Replicator:  testDrive.PublicAccount,
				BlochHashes: []*Hash{{123}},
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewEndDriveVerificationTransaction returned error: %s", err)
	tx.Signer = testDriveOwner.PublicAccount

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, endDriveVerificationTransactionToAggregateCorr, b)
}

func TestEndDriveVerificationTransactionSigning(t *testing.T) {
	tx, err := NewEndDriveVerificationTransaction(
		fakeDeadline,
		[]*FailureVerification{
			{
				Replicator:  testDrive.PublicAccount,
				BlochHashes: []*Hash{{123}},
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewEndDriveVerificationTransaction returned error: %s", err)

	b, err := testDriveOwner.Sign(tx)

	assert.Nilf(t, err, "Sign returned error: %s", err)
	assert.Equal(t, endDriveVerificationTransactionSigningCorr, b.Payload)
	assert.Equal(t, stringToHashPanic("cec4d0b335a2c7e0e5ab170b8f2073d3739a7ea1ed57b409fc10fa649699a421"), b.Hash)
}
