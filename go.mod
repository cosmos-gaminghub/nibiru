module github.com/cosmos-gaminghub/nibiru

go 1.16

require (
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/cosmos/cosmos-sdk v0.42.5
	github.com/ethereum/go-ethereum v1.10.3
	github.com/go-git/go-git/v5 v5.1.0
	github.com/goccy/go-yaml v1.8.9
	github.com/gookit/color v1.4.2
	github.com/gorilla/mux v1.8.0
	github.com/imdario/mergo v0.3.12
	github.com/mattn/go-zglob v0.0.3
	github.com/otiai10/copy v1.6.0
	github.com/pelletier/go-toml v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/radovskyb/watcher v1.0.7
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.34.10
	github.com/tendermint/tm-db v0.6.4
	golang.org/x/mod v0.4.2
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
