package ethereum

import (
	"context"
	"fmt"
	"math/big"
	entity "webapi_erc20/dao"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethereumTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

const (
	TransactionStatusSuccess = 1
	TransactionStatusFail    = 0
)

func GetTransactionFree(gasPrice decimal.Decimal, gasUsed int64) decimal.Decimal {
	return gasPrice.Mul(decimal.New(gasUsed, 0))
}

func GetTransactionReceipt(ctx context.Context, client *ethclient.Client, txHashStr string) (*ethereumTypes.Receipt, error) {
	txHash := common.HexToHash(txHashStr)

	result, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("client.TransactionReceipt error: %s", err)
	}

	return result, nil
}

func GetTransactionByHash(ctx context.Context, client *ethclient.Client, txHashStr string) (*ethereumTypes.Transaction, error) {
	txHash := common.HexToHash(txHashStr)

	result, _, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("client.TransactionByHash error: %s", err)
	}

	return result, nil
}

func GetTransactionByTokens(ctx context.Context, client *ethclient.Client, blockHeight int64, contractAddress []common.Address) ([]ethereumTypes.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(blockHeight),
		ToBlock:   big.NewInt(blockHeight),
		Addresses: contractAddress,
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return []ethereumTypes.Log{}, fmt.Errorf("client.FilterLogs error: %s", err)
	}

	return logs, nil
}

type EstimateTxFee struct {
	GasLimit    uint64
	GasPriceWei *big.Int

	TxFeeWei decimal.Decimal
}

func GetEstimateTxFee(ctx context.Context, client *ethclient.Client, token entity.Tokens) (EstimateTxFee, error) {
	// 取得平均瓦斯價格
	suggestGasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return EstimateTxFee{}, fmt.Errorf("SuggestGasPrice error: %s", err)
	}

	gasPrice := decimal.NewFromBigInt(suggestGasPrice, 0)

	if gasPrice.LessThan(token.GasPrice) {
		gasPrice = token.GasPrice
	}

	return EstimateTxFee{
		GasLimit:    uint64(token.GasLimit),
		GasPriceWei: gasPrice.BigInt(),
		TxFeeWei:    GetTransactionFree(gasPrice, token.GasLimit),
	}, nil
}

func GetNetworkID(ctx context.Context, client *ethclient.Client) (*big.Int, error) {
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("client.NetworkID error: %s", err)
	}

	return chainID, nil
}

func MakeSignedTx(ctx context.Context, client *ethclient.Client,
	fromAddress, toAddress, privateKey string, estimateTxFee EstimateTxFee, txAmount *big.Int) (*ethereumTypes.Transaction, error) {
	fromAddr := common.HexToAddress(fromAddress)
	toAddr := common.HexToAddress(toAddress)

	gasMsg := ethereum.CallMsg{
		From:     fromAddr,
		To:       &toAddr,
		Gas:      estimateTxFee.GasLimit,
		GasPrice: estimateTxFee.GasPriceWei,
		Value:    txAmount,
	}

	gas, err := client.EstimateGas(ctx, gasMsg)
	if err != nil {
		return nil, fmt.Errorf("EstimateGas, gas: %d, gasPrice: %d, error: %s",
			estimateTxFee.GasLimit, estimateTxFee.GasPriceWei, err)
	}

	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return nil, fmt.Errorf("client.PendingNonceAt error: %s", err)
	}

	legacyTx := &ethereumTypes.LegacyTx{
		To:       &toAddr,
		Nonce:    nonce,
		Value:    txAmount,
		Gas:      gas,
		GasPrice: estimateTxFee.GasPriceWei,
	}

	prv, err := HexToECDSA(privateKey)
	if err != nil {
		return nil, fmt.Errorf("HexToECDSA error: %s", err)
	}

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("client.NetworkID error: %s", err)
	}
	signer := ethereumTypes.NewEIP155Signer(chainID)

	signedTx, err := ethereumTypes.SignNewTx(prv, signer, legacyTx)
	if err != nil {
		return nil, fmt.Errorf("ethereumTypes.SignNewTx error: %s", err)
	}

	return signedTx, nil
}

func MakeSignNewTxTokens(ctx context.Context, client *ethclient.Client,
	fromAddress, toAddress, privateKey, contractAddress string,
	estimateTxFee EstimateTxFee, txAmount *big.Int) (*ethereumTypes.Transaction, error) {

	fromAddr := common.HexToAddress(fromAddress)
	toAddr := common.HexToAddress(toAddress)
	contractAddr := common.HexToAddress(contractAddress)

	methodID := getTransferContract()
	paddedToAddress := common.LeftPadBytes(toAddr.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(txAmount.Bytes(), 32)

	// 打包合約資料
	var txDate []byte
	txDate = append(txDate, methodID...)
	txDate = append(txDate, paddedToAddress...)
	txDate = append(txDate, paddedAmount...)

	txValue := big.NewInt(0)
	gasMsg := ethereum.CallMsg{
		From:     fromAddr,
		To:       &contractAddr,
		Gas:      estimateTxFee.GasLimit,
		GasPrice: estimateTxFee.GasPriceWei,
		Value:    txValue,
		Data:     txDate,
	}

	gas, err := client.EstimateGas(ctx, gasMsg)
	if err != nil {
		return nil, fmt.Errorf("EstimateGas, gas: %d, gasPrice: %d, error: %s",
			estimateTxFee.GasLimit, estimateTxFee.GasPriceWei, err)
	}

	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return nil, fmt.Errorf("client.PendingNonceAt error: %s", err)
	}

	legacyTx := &ethereumTypes.LegacyTx{
		To:       &contractAddr,
		Nonce:    nonce,
		Value:    txValue,
		Gas:      gas,
		GasPrice: estimateTxFee.GasPriceWei,
		Data:     txDate,
	}

	prv, err := HexToECDSA(privateKey)
	if err != nil {
		return nil, fmt.Errorf("HexToECDSA error: %s", err)
	}

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("client.NetworkID error: %s", err)
	}
	signer := ethereumTypes.NewEIP155Signer(chainID)

	signedTx, err := ethereumTypes.SignNewTx(prv, signer, legacyTx)
	if err != nil {
		return nil, fmt.Errorf("ethereumTypes.SignNewTx error: %s", err)
	}

	return signedTx, nil
}

func getTransferContract() []byte {
	transferFnSignature := []byte("transfer(address,uint256)")

	hash2 := crypto.Keccak256Hash(transferFnSignature)
	methodID := hash2[:4]

	return methodID
}

func SendTransaction(ctx context.Context, client *ethclient.Client, signedTx *ethereumTypes.Transaction) error {
	err := client.SendTransaction(ctx, signedTx)
	if err != nil {
		return fmt.Errorf("client.SendTransaction error: %s", err)
	}

	return nil
}
