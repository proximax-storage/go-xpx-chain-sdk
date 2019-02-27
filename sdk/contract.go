package sdk

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-utils/net"
	"golang.org/x/net/context"
	"net/http"
)

type ContractService service

func (ref *ContractService) GetContractsInfo(ctx context.Context, contractPubKeys ...string) ([]*ContractInfo, error) {
	if contractPubKeys == nil {
		return nil, errors.New("contract public key should not be nil")
	}

	pubKeys := struct {
		PublicKeys []string `json:"publicKeys"`
	}{
		PublicKeys: contractPubKeys,
	}

	dtos := contractInfoDTOs(make([]*contractInfoDTO, 0))

	resp, err := ref.client.DoNewRequest(ctx, http.MethodPost, contractsInfoRoute, pubKeys, &dtos)
	if err != nil {
		return nil, errors.Wrapf(err, "within POST request %s", contractsInfoRoute)
	}

	if err = handleResponseStatusCode(resp, map[int]error{409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	infos, err := dtos.toStruct(ref.client.config.NetworkType)
	if err != nil {
		return nil, errors.Wrap(err, "within converting dto to []*ContractInfo")
	}

	return infos, nil
}

func (ref *ContractService) GetContractsByAddress(ctx context.Context, address string) ([]*ContractInfo, error) {
	if len(address) == 0 {
		return nil, errors.New("address should not be blank")
	}

	url := net.NewUrl(fmt.Sprintf(contractsByAccountRoute, address))

	dtos := contractInfoDTOs(make([]*contractInfoDTO, 0))

	resp, err := ref.client.DoNewRequest(ctx, http.MethodGet, url.Encode(), nil, &dtos)
	if err != nil {
		return nil, errors.Wrapf(err, "within GET request %s", url.Encode())
	}

	if err = handleResponseStatusCode(resp, map[int]error{409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	infos, err := dtos.toStruct(ref.client.config.NetworkType)
	if err != nil {
		return nil, errors.Wrap(err, "within converting dto to []*ContractInfo")
	}

	return infos, nil
}
