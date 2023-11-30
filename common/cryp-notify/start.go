package cryp_notify

var (
	MerchantType2URL = make(map[int]string)
)

func Start(notifyURL string) {
	MerchantType2URL = map[int]string{
		0: notifyURL,
	}

	initTransaction()
}
