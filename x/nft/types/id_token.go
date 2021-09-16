package types

import (
	"fmt"
	"strconv"
)

const (
	// token id must be longer thant 3 length
	MIN_TOKEN_ID = uint64(1)
	// iris nft module restrict only number denom id, therefore appending alphanumeric prefix
	TOKEN_ID_PREFIX = "tokenid"
)

type TokenID uint64

func (id TokenID) ToIris() string {
	return TOKEN_ID_PREFIX + fmt.Sprintf("%d", id)
}

func (id TokenID) Uint64() uint64 {
	return uint64(id)
}

func FromIrisTokenID(id string) (TokenID, error) {
	tokenID, err := strconv.ParseUint(id[7:], 10, 64)
	return TokenID(tokenID), err
}

func ToTokenID(id string) (TokenID, error) {
	tokenID, err := strconv.ParseUint(id, 10, 64)
	return TokenID(tokenID), err
}
