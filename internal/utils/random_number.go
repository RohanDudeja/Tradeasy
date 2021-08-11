package utils

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

func GetRandNum() (string, error) {
	nBig, e := rand.Int(rand.Reader, big.NewInt(9999))
	if e != nil {
		return "", e
	}
	return strconv.FormatInt(nBig.Int64()+1000, 10), nil
}
