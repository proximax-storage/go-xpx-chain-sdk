package sdk

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
)

type ContractService service

func (ref *ContractService) GetContractInfo(ctx context.Context, contractPubKeys ...string) ([]*ContractInfo, error) {
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
