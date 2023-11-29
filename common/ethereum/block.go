package ethereum

import (
	"context"
	"fmt"
	"math/big"

	ethereumTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrBlockFailed = "server returned"
)

func GetBlockNumberLatest(ctx context.Context, client *ethclient.Client) (int64, error) {
	// 注意: nil = 取最新區塊
	result, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("client.BlockByNumber by latest block error: %s", err)
	}

	return int64(result.NumberU64()), nil
}

func GetBlockByNumber(ctx context.Context, client *ethclient.Client, blockNumber int64) (*ethereumTypes.Block, error) {
	result, err := client.BlockByNumber(ctx, big.NewInt(blockNumber))
	if err != nil {
		return nil, fmt.Errorf("client.BlockByNumber error: %s", err)
	}

	return result, nil
}
