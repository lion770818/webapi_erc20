package service

import (
	"context"
	"errors"
	"fmt"
	"webapi_erc20/common/ethereum"
	response "webapi_erc20/common/rsp"
	"webapi_erc20/dao"
	"webapi_erc20/define"
	"webapi_erc20/entity"

	"webapi_erc20/common/logs"

	"gorm.io/gorm"
)

func Create(ctx context.Context, req entity.AddressCreateReq) (entity.AddressCreateResp, response.Status, error) {
	newAddr, err := ethereum.GenerateAddress()
	if err != nil {
		return entity.AddressCreateResp{}, response.CodeInternalError, fmt.Errorf("ethereum.GenerateAddress error: %s", err)
	}

	addr := dao.Address{
		MerchantType: define.MerchantID2Type[req.MerchantID],
		Address:      newAddr.Address,
		ChainType:    req.ChainType,
	}

	_, err = dao.GetAddressInstance().Create(ctx, addr)
	if err != nil {
		return entity.AddressCreateResp{}, response.CodeInternalError, fmt.Errorf("AddressRepo.Create error: %s", err)
	}

	result := entity.AddressCreateResp{
		Address:   newAddr.Address,
		SecretKey: newAddr.PrivateKey,
		PublicKey: newAddr.PublicKey,
	}

	return result, response.Status{}, nil
}

func GetByAddress(ctx context.Context, addr string) (dao.Address, error) {
	return dao.GetAddressInstance().GetByAddress(ctx, addr)
}

// 取得錢包餘額
func GetBalance(ctx context.Context, req entity.AddressGetBalanceReq) (entity.AddressGetBalanceResp, response.Status, error) {

	// 檢查地址合法性
	if ethereum.IsValidateAddressFail(req.Address) {
		return entity.AddressGetBalanceResp{}, response.CodeAddressInvalidLength, errors.New(response.CodeAddressInvalidLength.Messages)
	}

	if req.CryptoType == define.CryptoType {
		// 如果是 ETH, 則呼叫原生 getBalance
		balance, err := ethereum.GetBalanceETH(req.Address)
		if err != nil {
			return entity.AddressGetBalanceResp{}, response.CodeInternalError, fmt.Errorf("ethereum.GetBalanceETH error: %s", err)
		}
		logs.Debugf("原生eth 錢包資訊 address:%v, balance:%v", req.Address, balance)

		result := entity.AddressGetBalanceResp{
			Balance: balance,
		}

		return result, response.Status{}, nil
	}

	// 如果是 其他token, 則先去DB撈取 token的合約之訊
	tokens, err := dao.GetTokenInstance().GetByCryptoType(ctx, req.CryptoType)
	if err != nil && err != gorm.ErrRecordNotFound {
		return entity.AddressGetBalanceResp{}, response.CodeInternalError, fmt.Errorf("TokensUseCase.GetByCryptoType error: %s", err)
	}

	if err == gorm.ErrRecordNotFound {
		return entity.AddressGetBalanceResp{}, response.CodeCryptoNotFound, errors.New(response.CodeCryptoNotFound.Messages)
	}

	// 透過合約去獲取 token錢包餘額
	tokenBalance, err := ethereum.GetBalanceToken(tokens.ContractAddr, tokens.ContractAbi, req.Address)
	if err != nil {
		return entity.AddressGetBalanceResp{}, response.CodeInternalError, fmt.Errorf("ethereum.GetBalanceToken error: %s", err)
	}

	logs.Debugf("tokens:%+v, tokenBalance:%v", tokens, tokenBalance)

	// 小數點位數轉換
	tokenBalance = ethereum.ConvertBalanceToAmount(tokenBalance, tokens.Decimals)

	result := entity.AddressGetBalanceResp{
		Balance: tokenBalance,
	}

	return result, response.Status{}, nil
}
