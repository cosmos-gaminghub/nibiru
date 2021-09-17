module github.com/cosmos-gaminghub/nibiru

go 1.16

require (
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cosmos/cosmos-sdk v0.44.0
	github.com/ethereum/go-ethereum v1.10.3
	github.com/pkg/errors v0.9.1
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/tendermint/tendermint v0.34.13
	github.com/tendermint/tm-db v0.6.4
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	github.com/cosmos/ibc-go v1.2.0
	github.com/gravity-devs/liquidity v1.4.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
	github.com/tendermint/tendermint => github.com/tendermint/tendermint v0.34.13
)
