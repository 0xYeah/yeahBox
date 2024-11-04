package unit_tools

import (
	"fmt"
	"math/big"
)

// UnitFormatWith1024 单位计算公式 以byte 为单位传入，根据数值不同换算成进制单位MB、GB等
// fSed 小数点保留位数
func UnitFormatWith1024(baseWithByte *big.Float, fSed int) string {
	k := big.NewFloat(1024)
	m := new(big.Float).Mul(k, k)
	g := new(big.Float).Mul(m, k)
	t := new(big.Float).Mul(g, k)
	p := new(big.Float).Mul(t, k)
	e := new(big.Float).Mul(p, k)
	z := new(big.Float).Mul(e, k)
	y := new(big.Float).Mul(z, k)

	format := fmt.Sprintf("%%.%df", fSed)

	newStr := ""
	switch {
	case baseWithByte.Cmp(y) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, y)), "YB")
	case baseWithByte.Cmp(z) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, z)), "ZB")
	case baseWithByte.Cmp(e) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, e)), "EB")
	case baseWithByte.Cmp(p) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, p)), "PB")
	case baseWithByte.Cmp(t) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, t)), "TB")
	case baseWithByte.Cmp(g) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, g)), "GB")
	case baseWithByte.Cmp(m) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, m)), "MB")
	case baseWithByte.Cmp(k) >= 0:
		newStr = fmt.Sprintf("%s%s", fmt.Sprintf(format, new(big.Float).Quo(baseWithByte, k)), "KB")
	default:
		newStr = fmt.Sprintf("%s", baseWithByte)
	}

	return newStr
}
