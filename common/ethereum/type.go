package ethereum

import (
	"math"
	"math/big"

	"github.com/shopspring/decimal"
)

var (
	ErrTransferAmountExceedsBalance = "transfer amount exceeds balance"

	ethToWeiUnit = int32(18)
	WeiToETHUnit = int32(-18)
)

func ConvWei(v decimal.Decimal, decimals int) decimal.Decimal {
	return v.Mul(decimal.New(1, int32(decimals)))
}

func ETHToWei(v decimal.Decimal) decimal.Decimal {
	return v.Mul(decimal.New(1, ethToWeiUnit))
}

func WeiToETH(v *big.Int) decimal.Decimal {
	return decimal.NewFromBigInt(v, WeiToETHUnit)
}

func ConvertBalanceToAmount(amount decimal.Decimal, decimals int) decimal.Decimal {
	return amount.Div(decimal.NewFromFloat(math.Pow10(decimals)))
}

func ConvertBalanceToAmountByBigInt(amount *big.Int, decimals int) decimal.Decimal {
	return ConvertBalanceToAmount(decimal.NewFromBigInt(amount, 0), decimals)
}
