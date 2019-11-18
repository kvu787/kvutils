package util

import (
	"crypto/rand"
	"math"
	"math/big"
)

func GetRandUint64() (uint64, error) {
	maxUint64 := &big.Int{}
	maxUint64.SetUint64(math.MaxUint64)
	bigInt, err := rand.Int(rand.Reader, maxUint64)
	if err != nil {
		return 0, err
	}
	return bigInt.Uint64(), nil
}
