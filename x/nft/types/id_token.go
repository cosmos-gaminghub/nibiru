package types

import (
	"fmt"
	"strconv"
)

const (
	// token id must be longer thant 3 length
	MIN_TOKEN_ID = uint64(100)
)

type TokenID uint64

func (id TokenID) String() string {
	return fmt.Sprintf("%d", id)
}

func (id TokenID) Uint64() uint64 {
	return uint64(id)
}

func ToTokenID(id string) (TokenID, error) {
	tokenID, err := strconv.ParseUint(id, 10, 64)
	return TokenID(tokenID), err
}
