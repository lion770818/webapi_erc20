package define

const (
	CryptoType = "ETH"

	ChainType = "ETH"
)

const (
	MerchantTypeOPDev = iota
	MerchantTypeOPPre
	MerchantTypeOP
	MerchantTypeQA

	MerchantTypeOPDevName = "OP_DEV"
	MerchantTypeOPPreName = "OP_PRE"
	MerchantTypeOPName    = "OP"
	MerchantTypeQAName    = "QA"
)

var (
	MerchantType2Name = map[int]string{
		MerchantTypeOPDev: MerchantTypeOPDevName,
		MerchantTypeOPPre: MerchantTypeOPPreName,
		MerchantTypeOP:    MerchantTypeOPName,
		MerchantTypeQA:    MerchantTypeQAName,
	}

	MerchantID2Type = map[string]int{
		MerchantTypeOPDevName: MerchantTypeOPDev,
		MerchantTypeOPPreName: MerchantTypeOPPre,
		MerchantTypeOPName:    MerchantTypeOP,
		MerchantTypeQAName:    MerchantTypeQA,
	}
)

const (
	TxTypeDeposit  = iota + 1 // 儲值
	TxTypeWithdraw            // 提款
)

const (
	TxStatusWaitConfirm = iota
	TxStatusSuccess
	TxStatusFail
)

const (
	TxNotifyStatusNotYetProcessed = iota
	TxNotifyStatusWaitNotify
	TxNotifyStatusSuccess
	TxNotifyStatusFail
)
