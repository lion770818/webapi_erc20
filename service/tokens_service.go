package service

import (
	"context"
	"fmt"
	"strings"
	"webapi_erc20/dao"

	"github.com/ethereum/go-ethereum/common"
)

func GetByCryptoType(ctx context.Context, cryptoType string) (dao.Tokens, error) {
	return dao.GetTokenInstance().GetByCryptoType(ctx, cryptoType)
}

func GetContractAddress(ctx context.Context) ([]common.Address, error) {
	tokens, err := dao.GetTokenInstance().GetList(ctx)
	if err != nil {
		return nil, fmt.Errorf("TokensRepo.GetContractAddress error: %s", err)
	}

	result := make([]common.Address, 0, len(tokens))

	for _, v := range tokens {
		if len(v.ContractAddr) > 0 {
			result = append(result, common.HexToAddress(v.ContractAddr))
		}
	}

	return result, nil
}

func GetContractAddr2Tokens(ctx context.Context) (map[string]dao.Tokens, error) {
	tokens, err := dao.GetTokenInstance().GetList(ctx)
	if err != nil {
		return nil, fmt.Errorf("TokensRepo.GetContractAddress error: %s", err)
	}

	contractAddr2Token := make(map[string]dao.Tokens)

	for _, v := range tokens {
		if len(v.ContractAddr) > 0 {
			contractAddr2Token[strings.ToLower(v.ContractAddr)] = v
		}
	}

	return contractAddr2Token, nil
}
