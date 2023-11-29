package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"webapi_erc20/common/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/common"
)

func GetBalanceETH(addrStr string) (decimal.Decimal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), utils.Time30S)
	defer cancel()

	client, err := getClient(ctx)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("client error: %s", err)
	}

	result, err := getBalance(ctx, client, addrStr)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("get balance wei error: %s", err)
	}

	balance := WeiToETH(result)

	return balance, nil
}

func GetBalanceWei(ctx context.Context, client *ethclient.Client, addrStr string) (decimal.Decimal, error) {
	result, err := getBalance(ctx, client, addrStr)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("get balance wei error: %s", err)
	}

	balance := decimal.NewFromBigInt(result, 0)

	return balance, nil
}

func getBalance(ctx context.Context, client *ethclient.Client, addrStr string) (*big.Int, error) {
	addr := common.HexToAddress(addrStr)

	// 注意: nil = 取最新餘額
	result, err := client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("BalanceAt error: %s", err)
	}

	return result, nil
}

func GetBalanceToken(contractAddr, contractAbi string, address string) (decimal.Decimal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), utils.Time30S)
	defer cancel()

	client, err := getClient(ctx)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("client error: %s", err)
	}

	tokenAddress := common.HexToAddress(contractAddr)
	token, err := NewTokenCaller(tokenAddress, client, contractAbi)
	if err != nil {
		return decimal.Decimal{}, err
	}

	account := common.HexToAddress(address)
	balance, err := token.BalanceOf(&bind.CallOpts{}, account)
	if err != nil {
		// 注意: 代表尚未在合約交易過，所以在合約查無此地址
		if err == bind.ErrNoCode {
			return decimal.Decimal{}, nil
		}

		return decimal.Decimal{}, err
	}

	return decimal.NewFromBigInt(balance, 0), nil
}
