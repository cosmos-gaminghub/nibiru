module github.com/cosmos-gaminghub/nibiru

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.44.2
	github.com/cosmos/ibc-go v1.2.0
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.2.1
	github.com/syndtr/goleveldb v1.0.1-0.20210305035536-64b5b1c73954 // indirect
	github.com/tendermint/tendermint v0.34.13
	github.com/tendermint/tm-db v0.6.4
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tendermint/tendermint => github.com/tendermint/tendermint v0.34.13
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)