package types

import (
	"fmt"
	"strconv"
)

const (
	// denom id must be longer thant 3 length
	MIN_DENOM_ID = uint64(100)
)

type DenomID uint64

func (id DenomID) String() string {
	return fmt.Sprintf("%d", id)
}

func (id DenomID) Uint64() uint64 {
	return uint64(id)
}

func ToDenomID(id string) (DenomID, error) {
	denomID, err := strconv.ParseUint(id, 10, 64)
	return DenomID(denomID), err
}
