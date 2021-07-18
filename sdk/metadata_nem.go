package sdk

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/net"
	"golang.org/x/crypto/sha3"
	"net/http"
)

type MetadataNemService service

func (ref *MetadataNemService) GetMetadataNemInfo(ctx context.Context, computedHash *Hash) (*MetadataNemTupleInfo, error) {
	if computedHash == nil {
		return nil, ErrNilHash
	}

	url := net.NewUrl(fmt.Sprintf(metadataEntryHashRoute, computedHash.String()))

	dto := &metadataNemInfoDTO{}

	resp, err := ref.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	mscInfo, err := dto.toStruct(ref.client.config.NetworkType)
	if err != nil {
		return nil, err
	}

	return mscInfo, nil
}

func (ref *MetadataNemService) GetMetadataNemInfos(ctx context.Context, hashes []*Hash) ([]*MetadataNemTupleInfo, error) {
	if len(hashes) == 0 {
		return nil, ErrNilHashes
	}

	dtos := metadataNemInfoDTOs(make([]*metadataNemInfoDTO, 0))

	resp, err := ref.client.doNewRequest(ctx, http.MethodPost, metadataEntriesRoute, &computedHashes{hashes}, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{400: ErrInvalidRequest, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	mscInfos, err := dtos.toStruct(ref.client.config.NetworkType)
	if err != nil {
		return nil, err
	}

	return mscInfos, nil
}

func CalculateUniqueAccountMetadataId(sourceAddress *Address, targetAccount *PublicAccount, key ScopedMetadataKey) (*Hash, error) {
	return calculate(sourceAddress, targetAccount, key, 0, 0)
}

func CalculateUniqueMosaicMetadataId(sourceAddress *Address, targetAccount *PublicAccount, key ScopedMetadataKey, mosaic *MosaicId) (*Hash, error) {
	return calculate(sourceAddress, targetAccount, key, mosaic.baseInt64, 1)
}

func CalculateUniqueNamespaceMetadataId(sourceAddress *Address, targetAccount *PublicAccount, key ScopedMetadataKey, namespace *NamespaceId) (*Hash, error) {
	return calculate(sourceAddress, targetAccount, key, namespace.baseInt64, 2)
}

func calculate(sourceAddress *Address, targetAccount *PublicAccount, key ScopedMetadataKey, targetId baseInt64, metadataType uint8) (*Hash, error) {
	result := sha3.New256()
	source, err := sourceAddress.Decode()
	if err != nil {
		return nil, err
	}

	targetKey, err := hex.DecodeString(targetAccount.PublicKey)
	if err != nil {
		return nil, err
	}

	if _, err := result.Write(source[:]); err != nil {
		return nil, err
	}

	if _, err := result.Write(targetKey[:]); err != nil {
		return nil, err
	}

	if _, err := result.Write(key.toLittleEndian()); err != nil {
		return nil, err
	}

	if _, err := result.Write(targetId.toLittleEndian()); err != nil {
		return nil, err
	}

	if _, err := result.Write([]byte{metadataType}); err != nil {
		return nil, err
	}

	hash, err := bytesToHash(result.Sum(nil))
	if err != nil {
		return nil, err
	}

	return hash, nil
}
