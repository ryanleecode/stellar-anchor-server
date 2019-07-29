package stellar

import (
	"math"
	"math/big"
	"strconv"
)

const int64Len = 19

func FormatToAssetPrecision(number big.Int) int64 {
	str := number.String()
	length := len(str)

	sliceLen := int(math.Min(int64Len, float64(length)))
	slice := str[:sliceLen]

	amount, _ := strconv.Atoi(slice)

	return int64(amount)
}
