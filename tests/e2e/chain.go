package e2e

import (
	"fmt"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	nibiru "github.com/cosmos-gaminghub/nibiru/app"
	"github.com/cosmos-gaminghub/nibiru/app/params"
)

const (
	keyringPassphrase = "testpassphrase"
	keyringAppName    = "testnet"
)

var (
	encodingConfig params.EncodingConfig
	cdc            codec.Codec
)

func init() {
	encodingConfig = nibiru.MakeEncodingConfig()

	encodingConfig.InterfaceRegistry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&stakingtypes.MsgCreateValidator{},
	)
	encodingConfig.InterfaceRegistry.RegisterImplementations(
		(*cryptotypes.PubKey)(nil),
		&secp256k1.PubKey{},
		&ed25519.PubKey{},
	)

	cdc = encodingConfig.Marshaler

	config := sdk.GetConfig()
	nibiru.SetBech32AddressPrefixes(config)
	// bech32MainPrefix := "cosmos"
	// config.SetBech32PrefixForAccount(bech32MainPrefix, bech32MainPrefix+sdk.PrefixPublic)
	// config.SetBech32PrefixForValidator(bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator, bech32MainPrefix+sdk.Bech32MainPrefix+sdk.PrefixOperator+sdk.PrefixPublic)
	// config.SetBech32PrefixForConsensusNode(bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus, bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
	config.Seal()
}

type chain struct {
	dataDir    string
	id         string
	validators []*validator
}

func newChain() (*chain, error) {
	tmpDir, err := ioutil.TempDir("", "nibiru-e2e-testnet-")
	if err != nil {
		return nil, err
	}

	return &chain{
		id:      "chain-" + tmrand.NewRand().Str(6),
		dataDir: tmpDir,
	}, nil
}

func (c *chain) configDir() string {
	return fmt.Sprintf("%s/%s", c.dataDir, c.id)
}

func (c *chain) createAndInitValidators(count int) error {
	for i := 0; i < count; i++ {
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(); err != nil {
			return err
		}

		c.validators = append(c.validators, node)

		// create keys
		if err := node.createKey("val"); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}

	return nil
}

// func (c *chain) createAndInitValidatorsWithMnemonics(count int, mnemonics []string) error {
// 	for i := 0; i < count; i++ {
// 		// create node
// 		node := c.createValidator(i)

// 		// generate genesis files
// 		if err := node.init(); err != nil {
// 			return err
// 		}

// 		c.validators = append(c.validators, node)

// 		// create keys
// 		if err := node.createKeyFromMnemonic("val", mnemonics[i]); err != nil {
// 			return err
// 		}
// 		if err := node.createNodeKey(); err != nil {
// 			return err
// 		}
// 		if err := node.createConsensusKey(); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (c *chain) createValidator(index int) *validator {
	return &validator{
		chain:   c,
		index:   index,
		moniker: fmt.Sprintf("%s-nibiru-%d", c.id, index),
	}
}
