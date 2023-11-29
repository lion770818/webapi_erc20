package ethereum

import (
	"context"
	"fmt"
	"webapi_erc20/common/config"
	"webapi_erc20/common/utils"

	"github.com/ethereum/go-ethereum/ethclient"
)

// type EthClient struct {
// 	NodeUrl []string
// }

// var instance *EthClient

// func GetInstance() *EthClient {
// 	if instance == nil {
// 		instance = &EthClient{
// 			NodeUrl: config.GetConfig().Config.Node.Url,
// 		}
// 	}
// 	return instance
// }

func getClient(ctx context.Context) (*ethclient.Client, error) {
	var client *ethclient.Client
	var err error
	nodeUrl := config.GetConfig().Config.Node.Url

	for _, url := range nodeUrl {
		client, err = ethclient.DialContext(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("ethclient.DialContext error: %w", err)
		}

		if ping(client) {
			break
		}
	}

	return client, nil
}

func GetClientNoCtx() (*ethclient.Client, error) {
	var client *ethclient.Client
	var err error
	nodeUrl := config.GetConfig().Config.Node.Url

	for _, url := range nodeUrl {
		client, err = ethclient.Dial(url)
		if err != nil {
			return nil, fmt.Errorf("ethclient.Dial error: %w", err)
		}

		if ping(client) {
			break
		}
	}

	return client, nil
}

func ping(client *ethclient.Client) bool {
	ctx, cancel := context.WithTimeout(context.Background(), utils.Time30S)
	defer cancel()

	_, err := GetBlockNumberLatest(ctx, client)
	if err != nil {
		return false
	}

	return true
}
