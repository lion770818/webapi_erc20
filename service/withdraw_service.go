package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"webapi_erc20/common/ethereum"
	"webapi_erc20/common/logs"
	response "webapi_erc20/common/rsp"
	"webapi_erc20/dao"
	"webapi_erc20/define"
	"webapi_erc20/entity"

	"github.com/ethereum/go-ethereum/core"
	"gorm.io/gorm"
)

// 轉帳
func Withdraw(ctx context.Context, req entity.WithdrawCreateReq) (entity.WithdrawCreateResp, response.Status, error) {
	if ethereum.IsValidateAddressFail(req.FromAddress) || ethereum.IsValidateAddressFail(req.ToAddress) {
		return entity.WithdrawCreateResp{}, response.CodeAddressInvalidLength, errors.New(response.CodeAddressInvalidLength.Messages)
	}

	// 取得db內token資訊
	tokens, err := dao.GetTokenInstance().GetByCryptoType(ctx, req.CryptoType)
	if err != nil && err != gorm.ErrRecordNotFound {
		return entity.WithdrawCreateResp{}, response.CodeInternalError, fmt.Errorf("TokensUseCase.GetByCryptoType error: %s", err)
	}

	if err == gorm.ErrRecordNotFound {
		return entity.WithdrawCreateResp{}, response.CodeCryptoNotFound, errors.New(response.CodeCryptoNotFound.Messages)
	}

	txHash, err := sendTransaction(ctx, req, tokens)
	if err != nil {
		if strings.Index(err.Error(), core.ErrInsufficientFundsForTransfer.Error()) > -1 {
			return entity.WithdrawCreateResp{}, response.CodeBalanceInsufficient, errors.New(response.CodeBalanceInsufficient.Messages)
		}

		if strings.Index(err.Error(), ethereum.ErrTransferAmountExceedsBalance) > -1 {
			return entity.WithdrawCreateResp{}, response.CodeBalanceInsufficient, errors.New(response.CodeBalanceInsufficient.Messages)
		}

		return entity.WithdrawCreateResp{}, response.CodeInternalError, fmt.Errorf("sendTransaction error: %s", err)
	}

	respStatus, err := create(ctx, req, txHash)
	if err != nil {
		return entity.WithdrawCreateResp{}, respStatus, fmt.Errorf("create error: %s", err)
	}

	result := entity.WithdrawCreateResp{
		TxHash: txHash,
		Memo:   req.Memo,
	}

	return result, response.Status{}, nil
}

// 建立交易單
func sendTransaction(ctx context.Context, req entity.WithdrawCreateReq, tokens dao.Tokens) (string, error) {
	if req.CryptoType == define.CryptoType {
		//送出eth交易單到鏈上
		return sendTransactionByETH(ctx, req, tokens)
	}

	return sendTransactionByToken(ctx, req, tokens)
}

// 送出eth交易單到鏈上
func sendTransactionByETH(ctx context.Context, req entity.WithdrawCreateReq, tokens dao.Tokens) (string, error) {
	client, err := ethereum.GetClientNoCtx()
	if err != nil {
		return "", fmt.Errorf("ethereum.GetClientNoCtx error: %s", err)
	}

	// 獲取平均燃氣價格
	estimateTxFee, err := ethereum.GetEstimateTxFee(ctx, client, tokens)
	if err != nil {
		return "", fmt.Errorf("ethereum.GetEstimateTxFee error: %s", err)
	}

	// 轉帳金額處理
	amount := ethereum.ConvWei(req.Amount, tokens.Decimals)

	// 交易簽章
	signedTx, err := ethereum.MakeSignedTx(ctx, client, req.FromAddress, req.ToAddress, req.SecretKey,
		estimateTxFee, amount.BigInt())
	if err != nil {
		return "", fmt.Errorf("ethereum.MakeSignedTx error: %s", err)
	}

	// 送交意到鏈上
	err = ethereum.SendTransaction(ctx, client, signedTx)
	if err != nil {
		return "", fmt.Errorf("ethereum.SendTransaction error: %s", err)
	}

	logs.Debugf("req:%+v, amount:%v, txHash:%v, estimateTxFee:%v, decimals:%v, gasPrice:%v, gasLimit:%v",
		req, amount, signedTx.Hash().String(), estimateTxFee, tokens.Decimals, tokens.GasPrice, tokens.GasLimit)

	return signedTx.Hash().String(), nil
}

// 送出erc20 token 交易單到鏈上
func sendTransactionByToken(ctx context.Context, req entity.WithdrawCreateReq, tokens dao.Tokens) (string, error) {
	client, err := ethereum.GetClientNoCtx()
	if err != nil {
		return "", fmt.Errorf("ethereum.GetClientNoCtx error: %s", err)
	}

	// 獲取平均燃氣價格
	estimateTxFee, err := ethereum.GetEstimateTxFee(ctx, client, tokens)
	if err != nil {
		return "", fmt.Errorf("ethereum.GetEstimateTxFee error: %s", err)
	}

	amount := ethereum.ConvWei(req.Amount, tokens.Decimals)

	// 使用 token 內的 合約地址 來處理 交易單
	signedTx, err := ethereum.MakeSignNewTxTokens(ctx, client, req.FromAddress, req.ToAddress,
		req.SecretKey, tokens.ContractAddr, estimateTxFee, amount.BigInt())
	if err != nil {
		return "", fmt.Errorf("ethereum.MakeSignNewTxTokens error: %s", err)
	}

	// 將交易單送到鏈上
	err = ethereum.SendTransaction(ctx, client, signedTx)
	if err != nil {
		return "", fmt.Errorf("ethereum.SendTransaction error: %s", err)
	}

	logs.Debugf("req:%+v, amount:%v, txHash:%v, estimateTxFee:%v, decimals:%v, gasPrice:%v, gasLimit:%v, contractAddr:%v",
		req, amount, signedTx.Hash().String(), estimateTxFee, tokens.Decimals, tokens.GasPrice, tokens.GasLimit, tokens.ContractAddr)

	return signedTx.Hash().String(), nil
}

func create(ctx context.Context, req entity.WithdrawCreateReq, txHash string) (response.Status, error) {
	withdraw := dao.Withdraw{
		MerchantType: define.MerchantID2Type[req.MerchantID],
		TxHash:       txHash,
		CryptoType:   req.CryptoType,
		ChainType:    req.ChainType,
		FromAddress:  req.FromAddress,
		ToAddress:    req.ToAddress,
		Amount:       req.Amount,
		Memo:         req.Memo,
	}
	_, err := dao.GetWithdrawInstance().Create(ctx, withdraw)
	if err != nil {
		return response.CodeInternalError, fmt.Errorf("WithdrawRepo.Create error: %s", err)
	}

	return response.Status{}, nil
}

// func GetByTxHash(ctx context.Context, txHash string) (dao.Withdraw, error) {
// 	return dao.GetWithdrawInstance().GetByTxHash(ctx, txHash)
// }
