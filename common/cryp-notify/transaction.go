package cryp_notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"webapi_erc20/common/logs"

	"github.com/moul/http2curl"

	"github.com/go-resty/resty/v2"
)

var (
	Transaction *transaction
)

func initTransaction() {
	Transaction = &transaction{}
}

type transaction struct {
}

func (trans *transaction) CreateTransactionNotify(ctx context.Context, host string, req CreateTransactionNotifyReq) (string, int, error) {
	type CommonResult struct {
		Status struct {
			Code int    `json:"code"`
			Msg  string `json:"messages"`
		} `json:"status"`

		Data interface{} `json:"data"`
	}

	url := fmt.Sprintf("%s/transaction", host)

	body, err := json.Marshal(req)
	if err != nil {
		return "", 0, err
	}

	result := CommonResult{}

	client := resty.New()
	res, err := client.R().SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&result).
		Post(url)
	if err != nil {
		return "", 0, err
	}

	curl, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(string(body)))
	if err != nil {
		return "", 0, err
	}

	curl.Header.Set("Content-Type", "application/json")

	command, err := http2curl.GetCurlCommand(curl)
	if err != nil {
		return "", 0, err
	}

	if res.StatusCode() != http.StatusOK {
		return command.String(), 0, fmt.Errorf("call post transaction api fail, res.Status: %s", res.Status())
	}

	if result.Status.Code != 0 {
		logs.Errorf("createTransactionNotify fail, result.Code != 0 Code:%v", result.Status.Code)
		return command.String(), 3, nil // 3=fail
	}

	return command.String(), 0, nil // 0=success
}
