package service

import (
	"context"
	crypNotify "webapi_erc20/common/cryp-notify"
	response "webapi_erc20/common/rsp"
	"webapi_erc20/define"
	"webapi_erc20/entity"

	"errors"
	"fmt"
	"math/big"
	"strings"
	"webapi_erc20/common/config"
	"webapi_erc20/common/ethereum"
	"webapi_erc20/common/logs"
	"webapi_erc20/common/utils"
	"webapi_erc20/dao"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	ethereumI "github.com/ethereum/go-ethereum"
	ethereumTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/ethclient"

	"gorm.io/gorm"

	"github.com/robfig/cron/v3"

	"github.com/gofrs/uuid"
)

// type TransactionUseCaseCond struct {
// 	dig.In

// 	BlockHeightRepo repository.BlockHeightRepository
// 	TransRepo       repository.TransactionRepository

// 	AddressUseCase  usecase.AddressUseCase
// 	TokensUseCase   usecase.TokensUseCase
// 	WithdrawUseCase usecase.WithdrawUseCase

// 	DB *gorm.DB `name:"dbM"`
// }

// type transactionUseCase struct {
// 	TransactionUseCaseCond
// }

// func NewTransactionUseCase(cond TransactionUseCaseCond) usecase.TransactionUseCase {
// 	uc := &transactionUseCase{
// 		TransactionUseCaseCond: cond,
// 	}

// 	return uc
// }

func GetByTxHash(ctx context.Context, req entity.TransGetTxHashReq) (entity.TransGetTxHashResp, response.Status, error) {
	// 撈取
	tx, err := dao.GetTransactionInstance().GetByTxHash(ctx, req.TxHash)
	if err != nil && err != gorm.ErrRecordNotFound {
		return entity.TransGetTxHashResp{}, response.CodeInternalError, fmt.Errorf("TransRepo.GetByTxHash error: %s", err)
	}

	if err == gorm.ErrRecordNotFound {
		tx, err := dao.GetWithdrawInstance().GetByTxHash(ctx, req.TxHash)
		if err != nil && err != gorm.ErrRecordNotFound {
			return entity.TransGetTxHashResp{}, response.CodeInternalError, fmt.Errorf("TransRepo.GetByTxHash error: %s", err)
		}

		if err == gorm.ErrRecordNotFound {
			return entity.TransGetTxHashResp{}, response.CodeTxNotFound, errors.New(response.CodeTxNotFound.Messages)
		}

		return entity.TransGetTxHashResp{
			BlockHeight: 0,
			TxHash:      tx.TxHash,
			CryptoType:  tx.CryptoType,
			ChainType:   tx.ChainType,
			FromAddress: tx.FromAddress,
			ToAddress:   tx.ToAddress,
			Amount:      tx.Amount,
			Fee:         decimal.Zero,
			FeeCrypto:   "",
			Status:      define.TxStatusWaitConfirm,
			Memo:        tx.Memo,
		}, response.Status{}, nil
	}

	if tx.CryptoType != req.CryptoType {
		respStatus := response.CodeParamInvalid.WithMsg(", req.crypto_type is not eq tx.crypto_type")
		return entity.TransGetTxHashResp{}, respStatus, errors.New(respStatus.Messages)
	}

	result := entity.TransGetTxHashResp{
		BlockHeight: tx.BlockHeight,
		TxHash:      tx.TxHash,
		CryptoType:  tx.CryptoType,
		ChainType:   tx.ChainType,
		FromAddress: tx.FromAddress,
		ToAddress:   tx.ToAddress,
		Amount:      tx.Amount,
		Fee:         tx.Fee,
		FeeCrypto:   tx.FeeCrypto,
		Status:      tx.Status,
		Memo:        tx.Memo,
	}

	return result, response.Status{}, nil
}

func RunListenBlock() {
	c := cron.New()
	c.AddFunc("*/1 * * * *", runListenBlock)
	c.Start()
}

func runListenBlock() {

	uid, err := uuid.NewV4()
	if err != nil {
		logs.Debugf("uuid.NewV4 error: %s", err)
		return
	}

	ctx := context.WithValue(context.Background(), utils.LogUUID, uid.String())
	ctx, cancelCtx := context.WithTimeout(ctx, 3*utils.Time30S)
	defer cancelCtx()

	client, err := ethereum.GetClientNoCtx()
	if err != nil {
		logs.Debugf("ethereum.GetClientNoCtx error: %s", err)
		return
	}

	blockHeight, isHaveNewBlock, err := checkNewBlock(ctx, client)
	if err != nil {
		logs.Debugf("checkNewBlock error: %s", err)
		return
	}

	logs.Debugf("blockHeight:%v, isHaveNewBlock:%v", blockHeight, isHaveNewBlock)

	// 注意: 處理遞歸終止條件
	if !isHaveNewBlock {
		return
	}

	blockIsFail := false
	var trans []dao.Transaction

	block, err := ethereum.GetBlockByNumber(ctx, client, blockHeight)
	if err != nil {

		logs.Debugf("ethereum.GetBlockByNumber error: %s", err)

		if strings.Index(err.Error(), ethereum.ErrBlockFailed) > -1 {
			// sometimes there are bad blocks in the chain that require special judgment
			// bad blocks ex: https://mumbai.polygonscan.com/block/33364784
			// github: https://github.com/ethereum/go-ethereum/blob/master/ethclient/ethclient.go#L133 ~ L142
			blockIsFail = true
		} else {
			return
		}
	}

	if !blockIsFail {
		// makeTransaction (監聽區並建立 trans struct)
		trans, err = makeTransaction(ctx, client, block)
		if err != nil {
			logs.Debugf("makeTransaction error: %s", err)
			return
		}
	}

	err = createTransAndUpdateBlockHeight(ctx, blockHeight, trans)
	if err != nil {
		logs.Debugf("createTransAndUpdateBlockHeight blockHeight:%v, error: %s", blockHeight, err)
		return
	}

	// 注意: 用遞歸確保，同步到最新區塊資料
	runListenBlock()
}

func checkNewBlock(ctx context.Context, client *ethclient.Client) (int64, bool, error) {
	logs.Debugf("checkNewBlock")

	latestHeight, err := ethereum.GetBlockNumberLatest(ctx, client)
	if err != nil {
		return 0, false, fmt.Errorf("ethereum.GetBlockNumberLatest error: %s", err)
	}

	dbBlockHeight, err := dao.GetBlockHeightInstance().Get(ctx)
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, false, fmt.Errorf("BlockHeightRepo.Get error: %s", err)
	}
	logs.Debugf("latestHeight:%v, dbBlockHeight:%v", latestHeight, dbBlockHeight)

	if err == gorm.ErrRecordNotFound {
		dbBlockHeight = dao.BlockHeight{
			BlockHeight: latestHeight,
		}

		_, err := dao.GetBlockHeightInstance().Create(ctx, dbBlockHeight)
		if err != nil {
			return 0, false, fmt.Errorf("BlockHeightRepo.Create error: %s", err)
		}
	}

	blockHeight := dbBlockHeight.BlockHeight
	isHaveNewBlock := false

	// 注意: 用 <= 判斷原因是，更新邏輯是 blockHeight + 1
	if blockHeight <= latestHeight {
		isHaveNewBlock = true
	}

	logs.Debugf("latestHeight:%v, isHaveNewBlock:%v", latestHeight, isHaveNewBlock)

	return blockHeight, isHaveNewBlock, nil
}

// 建立轉帳單
func makeTransaction(ctx context.Context, client *ethclient.Client, block *ethereumTypes.Block) ([]dao.Transaction, error) {
	networkID, err := ethereum.GetNetworkID(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("ethereum.GetNetworkID error: %s", err)
	}

	// 建立eth轉帳單
	transETH, err := makeTransactionByETH(ctx, block, networkID)
	if err != nil {
		return nil, fmt.Errorf("makeTransactionByETH error: %s", err)
	}

	// token轉帳
	transTokens, err := makeTransactionByTokens(ctx, client, block)
	if err != nil {
		return nil, fmt.Errorf("makeTransactionByTokens error: %s", err)
	}

	size := len(transETH) + len(transTokens)
	result := make([]dao.Transaction, 0, size)

	result = append(result, transETH...)
	result = append(result, transTokens...)

	return result, nil
}

// 建立eth轉帳單
func makeTransactionByETH(ctx context.Context, block *ethereumTypes.Block, networkID *big.Int) ([]dao.Transaction, error) {
	transaction := make([]dao.Transaction, 0, 0)

	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			// 注意: == nil 代表為 Contract 地址，不處理略過
			continue
		}

		checkAmount := decimal.NewFromBigInt(tx.Value(), 0)
		if checkAmount.LessThanOrEqual(decimal.Zero) {
			// 注意: amount <= 0，不處理略過
			continue
		}

		fromAddr, err := getETHFromAddr(ctx, tx, networkID)
		if err != nil {
			if err == errTxChainIDNotEqualNetworkID {
				continue
			}

			return nil, fmt.Errorf("getETHFromAddr error: %s", err)
		}

		toAddr := tx.To().Hex()

		addressMap, err := getAddressMap(ctx, fromAddr, toAddr)
		if err != nil {
			// 注意: 代表此地址，不是服務產生過的，略過不處理
			continue
		}

		isNewTransaction, err := isNewTransaction(ctx, tx.Hash().String())
		if err != nil {
			return nil, fmt.Errorf("isNewTransaction error: %s", err)
		}

		if !isNewTransaction {
			continue
		}

		txType, err := getTxType(addressMap, fromAddr, toAddr)
		if err != nil {
			return nil, fmt.Errorf("geTxType error: %s", err)
		}

		trans := dao.Transaction{
			TxType:           txType,
			BlockHeight:      int64(block.NumberU64()),
			TransactionIndex: 0,
			TxHash:           tx.Hash().String(),
			CryptoType:       define.CryptoType,
			ChainType:        define.ChainType,
			ContractAddr:     "",
			FromAddress:      fromAddr,
			ToAddress:        toAddr,
			Amount:           ethereum.WeiToETH(tx.Value()),
			Gas:              int64(tx.Gas()),
			GasUsed:          0,
			GasPrice:         getGasPrice(tx, block.BaseFee()),
			Fee:              decimal.Zero,
			FeeCrypto:        "",
			Confirm:          0,
			Status:           define.TxStatusWaitConfirm,
			Memo:             "",
			NotifyStatus:     define.TxNotifyStatusNotYetProcessed,
		}

		transaction = append(transaction, trans)
	}

	return transaction, nil
}

// token轉帳
func makeTransactionByTokens(ctx context.Context, client *ethclient.Client, block *ethereumTypes.Block) ([]dao.Transaction, error) {
	contractAddress, err := GetContractAddress(ctx)
	if err != nil {
		return nil, fmt.Errorf("TokensUseCase.GetContractAddress error: %s", err)
	}

	blockHeight := int64(block.NumberU64())
	tokensLogs, err := ethereum.GetTransactionByTokens(ctx, client, blockHeight, contractAddress)
	if err != nil {
		return nil, fmt.Errorf("ethereum.GetTransactionByTokens error: %s", err)
	}

	contractAddr2Tokens, err := GetContractAddr2Tokens(ctx)
	if err != nil {
		return nil, fmt.Errorf("TokensUseCase.GetContractAddr2Token error: %s", err)
	}

	logTransferSigHash := getLogTransferSigHash()
	transaction := make([]dao.Transaction, 0, 0)

	for _, vLog := range tokensLogs {
		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			contractAddr := vLog.Address.Hex()
			tokens, ok := contractAddr2Tokens[strings.ToLower(contractAddr)]
			if !ok {
				continue
			}

			amount, err := getAmount(vLog.Data, tokens)
			if err != nil {
				return nil, fmt.Errorf("getAmount error: %s", err)
			}

			if amount.LessThanOrEqual(decimal.Zero) {
				// 注意: amount <= 0，不處理略過
				continue
			}

			txHash := vLog.TxHash.Hex()
			tx, err := ethereum.GetTransactionByHash(ctx, client, txHash)
			if err != nil {
				return nil, fmt.Errorf("ethereum.GetTransactionByHash error: %s", err)
			}

			fromAddr := common.HexToAddress(vLog.Topics[1].Hex()).Hex()
			toAddr := common.HexToAddress(vLog.Topics[2].Hex()).Hex()

			addressMap, err := getAddressMap(ctx, fromAddr, toAddr)
			if err != nil {
				// 注意: 代表此地址，不是服務產生過的，略過不處理
				continue
			}

			isNewTransaction, err := isNewTransaction(ctx, txHash)
			if err != nil {
				return nil, fmt.Errorf("isNewTransaction error: %s", err)
			}

			if !isNewTransaction {
				continue
			}

			txType, err := getTxType(addressMap, fromAddr, toAddr)
			if err != nil {
				return nil, fmt.Errorf("geTxType error: %s", err)
			}

			trans := dao.Transaction{
				TxType:           txType,
				BlockHeight:      int64(vLog.BlockNumber),
				TransactionIndex: 0,
				TxHash:           txHash,
				CryptoType:       tokens.CryptoType,
				ChainType:        tokens.ChainType,
				ContractAddr:     contractAddr,
				FromAddress:      fromAddr,
				ToAddress:        toAddr,
				Amount:           amount,
				Gas:              int64(tx.Gas()),
				GasUsed:          0,
				GasPrice:         getGasPrice(tx, block.BaseFee()),
				Fee:              decimal.Zero,
				FeeCrypto:        "",
				Confirm:          0,
				Status:           define.TxStatusWaitConfirm,
				Memo:             "",
				NotifyStatus:     define.TxNotifyStatusNotYetProcessed,
			}

			transaction = append(transaction, trans)
		}
	}

	return transaction, nil
}

func isNewTransaction(ctx context.Context, txHash string) (bool, error) {
	_, err := dao.GetTransactionInstance().GetByTxHash(ctx, txHash)
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, fmt.Errorf("TransRepo.GetByTxHash, txHash: %s, error: %s", txHash, err)
	}

	if err == gorm.ErrRecordNotFound {
		return true, nil
	}

	return false, nil
}

var errTxChainIDNotEqualNetworkID = errors.New("error tx chainID not equal networkID")

func getETHFromAddr(ctx context.Context, tx *ethereumTypes.Transaction, networkID *big.Int) (string, error) {
	var fromAddr common.Address
	var err error

	// if tx.ChainId().String() != networkID.String() {
	// 	logs.Errorf("ChainId != networkID message:%v, hash:%v, tx.ChainId:%v, networkID:%v", errTxChainIDNotEqualNetworkID, tx.Hash(), tx.ChainId().String(), networkID.String())
	// 	return "", errTxChainIDNotEqualNetworkID
	// }

	switch tx.Type() {
	case ethereumTypes.LegacyTxType, ethereumTypes.AccessListTxType:
		// 注意: 舊版 NewEIP155Signer 不支援處理多個 Tokens Transferred，應使用 NewEIP2930Signer 處理
		fromAddr, err = ethereumTypes.NewEIP2930Signer(tx.ChainId()).Sender(tx)
		if err != nil {
			return "", fmt.Errorf("tx.AsMessage(NewEIP2930Signer) ChainId: %s, error: %s",
				tx.ChainId().String(), err)
		}
	case ethereumTypes.DynamicFeeTxType:
		fromAddr, err = ethereumTypes.NewLondonSigner(tx.ChainId()).Sender(tx)
		if err != nil {
			return "", fmt.Errorf("tx.AsMessage(NewLondonSigner) ChainId: %s, error: %s",
				tx.ChainId().String(), err)
		}
	}

	return fromAddr.Hex(), nil
}

func getGasPrice(tx *ethereumTypes.Transaction, baseFee *big.Int) decimal.Decimal {
	switch tx.Type() {
	case ethereumTypes.LegacyTxType, ethereumTypes.AccessListTxType:
		return ethereum.WeiToETH(tx.GasPrice())
	case ethereumTypes.DynamicFeeTxType:
		return ethereum.WeiToETH(baseFee).
			Add(ethereum.WeiToETH(tx.GasTipCap()))
	}

	return decimal.Zero
}

func getAddressMap(ctx context.Context, fromAddr, toAddr string) (map[string]struct{}, error) {
	isExistFromAddress, isExistToAddress := false, false
	address2Struct := make(map[string]struct{})

	fromAddress, err := GetByAddress(ctx, fromAddr)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("AddressUseCase.GetByAddress by fromAddress error: %s", err)
	}

	if fromAddress.ID > 0 {
		isExistFromAddress = true
		address2Struct[fromAddr] = struct{}{}
	}

	toAddress, err := GetByAddress(ctx, toAddr)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("AddressUseCase.GetByAddress by toAddress error: %s", err)
	}

	if toAddress.ID > 0 {
		isExistToAddress = true
		address2Struct[toAddr] = struct{}{}
	}

	if !isExistFromAddress && !isExistToAddress {
		return nil, errors.New("address not found")
	}

	return address2Struct, nil
}

func getTxType(address2Struct map[string]struct{}, fromAddr, toAddr string) (int, error) {
	_, ok := address2Struct[fromAddr]
	if ok {
		return define.TxTypeWithdraw, nil
	}

	_, ok = address2Struct[toAddr]
	if ok {
		return define.TxTypeDeposit, nil
	}

	return 0, errors.New("txType not found")
}

func getLogTransferSigHash() common.Hash {
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	return logTransferSigHash
}

func getAmount(vLogData []byte, tokens dao.Tokens) (decimal.Decimal, error) {
	contractAbi, err := abi.JSON(strings.NewReader(tokens.ContractAbi))
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("abi json error: %s", err)
	}

	data, err := contractAbi.Unpack("Transfer", vLogData)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("abi json error: %s", err)
	}

	val, ok := data[0].(*big.Int)
	if !ok {
		return decimal.Decimal{}, fmt.Errorf("data[0].(*big.Int) error: %s", err)
	}

	num := decimal.NewFromBigInt(val, 0)
	amount := ethereum.ConvertBalanceToAmount(num, tokens.Decimals)

	return amount, nil
}

func createTransAndUpdateBlockHeight(ctx context.Context, blockHeight int64, trans []dao.Transaction) error {

	tx := dao.SqlSession.Begin()
	defer tx.Rollback()

	if len(trans) > 0 {

		transRepo := dao.GetTransactionInstance()
		for i := range trans {
			_, err := transRepo.Create(ctx, trans[i])
			if err != nil {

				logs.Debugf("transRepo.Create i:%v, blockHeight:%v, trans:%v, error: %s",
					i, blockHeight, trans[i], err)

				return fmt.Errorf("TransRepo.Create error: %s", err)
			}
		}
	}

	blockHeightRepo := dao.GetBlockHeightInstance()
	bh, err := blockHeightRepo.Get(ctx)
	if err != nil {
		return fmt.Errorf("BlockHeightRepo.Get error: %s", err)
	}

	bh.BlockHeight = blockHeight + 1
	err = blockHeightRepo.Update(ctx, bh)
	if err != nil {
		return fmt.Errorf("BlockHeightRepo.Update error: %s", err)
	}

	tx.Commit()

	return nil
}

func RunTransactionConfirm() cron.FuncJob {
	return func() {
		uid, err := uuid.NewV4()
		if err != nil {
			logs.Debugf("uuid.NewV4 error: %s", err)
			return
		}

		ctx := context.WithValue(context.Background(), utils.LogUUID, uid.String())
		ctx, cancelCtx := context.WithTimeout(ctx, 3*utils.Time30S)
		defer cancelCtx()

		client, err := ethereum.GetClientNoCtx()
		if err != nil {
			logs.Debugf("ethereum.GetClientNoCtx error: %s", err)
			return
		}

		trans, err := getTransactionConfirm(ctx, client)
		if err != nil {
			logs.Debugf("getTransactionConfirm error: %s", err)
			return
		}

		err = runTransactionConfirm(ctx, client, trans)
		if err != nil {
			logs.Debugf("runTransactionConfirm error: %s", err)
			return
		}
	}
}

func getTransactionConfirm(ctx context.Context, client *ethclient.Client) ([]dao.Transaction, error) {
	latestHeight, err := ethereum.GetBlockNumberLatest(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("ethereum.GetBlockNumberLatest error: %v", err)
	}

	blockHeight := latestHeight - config.GetConfig().Node.Confirm
	trans, err := dao.GetTransactionInstance().GetListByStatusAndBlockHeight(ctx, define.TxStatusWaitConfirm, 20, blockHeight)
	if err != nil {
		return nil, fmt.Errorf("TransRepo.GetListByStatusAndBlockHeight error: %s", err)
	}

	return trans, nil
}

func runTransactionConfirm(ctx context.Context, client *ethclient.Client, trans []dao.Transaction) error {
	for i := range trans {
		v := trans[i]

		receipt, err := ethereum.GetTransactionReceipt(ctx, client, v.TxHash)
		if err != nil && err != ethereumI.NotFound {
			return fmt.Errorf("ethereum.GetTransactionReceipt error: %v", err)
		}

		if err == ethereumI.NotFound {

			logs.Debugf("not found, transactions to be chained txHash:%v, error: %v", v.TxHash, err)

			continue
		}

		v.Status = define.TxStatusSuccess
		if receipt.Status != ethereum.TransactionStatusSuccess {
			v.Status = define.TxStatusFail
		}

		v.NotifyStatus = define.TxNotifyStatusWaitNotify

		gasUsed := int64(receipt.GasUsed)
		txFree := ethereum.GetTransactionFree(v.GasPrice, gasUsed)

		v.Confirm = config.GetConfig().Node.Confirm
		v.TransactionIndex = int(receipt.TransactionIndex)
		v.GasUsed = gasUsed
		v.Fee = txFree
		v.FeeCrypto = define.CryptoType

		err = dao.GetTransactionInstance().Update(ctx, v)
		if err != nil {
			return fmt.Errorf("TransRepo.Update, txHash %s, error: %v", v.TxHash, err)
		}
	}

	return nil
}

func RunTransactionNotify() cron.FuncJob {
	return func() {
		uid, err := uuid.NewV4()
		if err != nil {
			logs.Debugf("uuid.NewV4 error: %v", err)
			return
		}

		ctx := context.WithValue(context.Background(), utils.LogUUID, uid.String())
		ctx, cancelCtx := context.WithTimeout(ctx, 3*utils.Time30S)
		defer cancelCtx()

		err = runTransactionNotify(ctx)
		if err != nil {
			logs.Debugf("runTransactionNotify error: %v", err)
			return
		}
	}
}

func runTransactionNotify(ctx context.Context) error {
	trans, err := dao.GetTransactionInstance().GetListByNotifyStatus(ctx, define.TxNotifyStatusWaitNotify)
	if err != nil {
		return fmt.Errorf("TransRepo.GetListByNotifyStatus error: %s", err)
	}

	for i := range trans {
		v := trans[i]

		host, err := getTransactionNotifyHost(ctx, v)
		if err != nil {
			return fmt.Errorf("getTransactionNotifyHost, txHash: %s, error: %s", v.TxHash, err)
		}

		req := crypNotify.CreateTransactionNotifyReq{
			TxType:      v.TxType,
			BlockHeight: v.BlockHeight,
			TxHash:      v.TxHash,
			CryptoType:  v.CryptoType,
			ChainType:   v.ChainType,
			FromAddress: v.FromAddress,
			ToAddress:   v.ToAddress,
			Amount:      v.Amount,
			Fee:         v.Fee,
			FeeCrypto:   v.FeeCrypto,
			Status:      v.Status,
			Memo:        v.Memo,
		}
		curl, notifyStatus, err := crypNotify.Transaction.CreateTransactionNotify(ctx, host, req)
		if err != nil {
			logs.Debugf("crypNotify.Transaction.CreateTransactionNotify curl:%v, txHash:%v error: %v",
				curl, v.TxHash, err)
			continue
		}

		logs.Infof("crypNotify.Transaction.CreateTransactionNotify curl:%v", curl)

		v.NotifyStatus = notifyStatus

		err = dao.GetTransactionInstance().Update(ctx, v)
		if err != nil {
			return fmt.Errorf("TransRepo.Update, txHash %s, error: %s", v.TxHash, err)
		}
	}

	return nil
}

func getTransactionNotifyHost(ctx context.Context, trans dao.Transaction) (string, error) {
	merchantType, err := getMerchantType(ctx, trans)
	if err != nil {
		return "", fmt.Errorf("getMerchantType error: %s", err)
	}

	notifyURL, ok := crypNotify.MerchantType2URL[merchantType]
	if !ok {
		return "", fmt.Errorf("not found notifyURL merchantType: %d", merchantType)
	}

	return notifyURL, nil
}

func getMerchantType(ctx context.Context, trans dao.Transaction) (int, error) {
	if trans.TxType == define.TxTypeWithdraw {
		addr, err := GetByAddress(ctx, trans.FromAddress)
		if err != nil {
			return 0, fmt.Errorf("AddressUseCase.GetByAddress by fromAddress error: %s", err)
		}

		return addr.MerchantType, nil
	}

	addr, err := GetByAddress(ctx, trans.ToAddress)
	if err != nil {
		return 0, fmt.Errorf("AddressUseCase.GetByAddress by toAddress error: %s", err)
	}

	return addr.MerchantType, nil
}
