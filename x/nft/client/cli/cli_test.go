package cli_test

import (
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos-gaminghub/nibiru/testutil/network"
	"github.com/cosmos-gaminghub/nibiru/x/nft/client/cli"
	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	nftcli "github.com/irisnet/irismod/modules/nft/client/cli"
	irismodtypes "github.com/irisnet/irismod/modules/nft/types"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

type CliTestSuite struct {
	suite.Suite

	network *network.Network
}

func (s *CliTestSuite) SetupSuite() {
	s.T().Log("setting up cli test suite")
	s.network = network.New(s.T())
	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *CliTestSuite) TearDownSuite() {
	s.T().Log("tearing down cli test suite")
	s.network.Cleanup()
}

func TestCliTestSuite(t *testing.T) {
	suite.Run(t, new(CliTestSuite))
}

func (s *CliTestSuite) TestIssueDenomMintEditTransferBurnNFTQueryNFT() {
	var (
		val       = s.network.Validators[0]
		val2      = s.network.Validators[1]
		from      = val.Address
		ctx       = val.ClientCtx
		recepient = val2.Address
		denomID   = "nbr"
		tokenID   = "1"
		denomName = "nibiru token"
		schema    = "nibiru token schema"
		tokenName = "token name"
		tokenURI  = "https://cosmosgaminghub.com/"
		tokenData = "data"

		args         []string
		out          sdktestutil.BufferWriter
		err          error
		resp         sdk.TxResponse
		expectedCode = uint32(0)
	)

	//------test GetCmdIssueDenom()-------------
	args = []string{
		denomID,
		fmt.Sprintf("--%s=%s", nftcli.FlagDenomName, denomName),
		fmt.Sprintf("--%s=%s", nftcli.FlagSchema, schema),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.Config.BondDenom, sdk.NewInt(10))).String()),
	}

	out, err = clitestutil.ExecTestCLICmd(ctx, cli.GetCmdIssueDenom(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &resp))
	s.Require().Equal(expectedCode, resp.Code)

	//------test GetCmdMintNFT()-------------
	args = []string{
		denomID,
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, tokenName),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenURI, tokenURI),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenData, tokenData),
		fmt.Sprintf("--%s=%s", nftcli.FlagRecipient, from.String()),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.Config.BondDenom, sdk.NewInt(10))).String()),
	}

	out, err = clitestutil.ExecTestCLICmd(ctx, cli.GetCmdMintNFT(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &resp))
	s.Require().Equal(expectedCode, resp.Code)

	//------test GetCmdQueryNFT()-------------
	respType := proto.Message(&irismodtypes.BaseNFT{})
	args = []string{
		denomID,
		tokenID,
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	out, err = clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryNFT(), args)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), respType))
	nftItem := respType.(*irismodtypes.BaseNFT)
	s.Require().Equal(&irismodtypes.BaseNFT{
		Id:    types.TOKEN_ID_PREFIX + tokenID,
		Name:  tokenName,
		URI:   tokenURI,
		Data:  tokenData,
		Owner: from.String(),
	}, nftItem)

	//------test GetCmdEditNFT()-------------
	args = []string{
		denomID,
		tokenID,
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, tokenName),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenData, tokenData),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.Config.BondDenom, sdk.NewInt(10))).String()),
	}

	out, err = clitestutil.ExecTestCLICmd(ctx, cli.GetCmdEditNFT(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &resp))
	s.Require().Equal(expectedCode, resp.Code)

	//------test GetCmdTransferNFT()-------------
	args = []string{
		recepient.String(),
		denomID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.Config.BondDenom, sdk.NewInt(10))).String()),
	}

	out, err = clitestutil.ExecTestCLICmd(ctx, cli.GetCmdTransferNFT(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &resp))
	s.Require().Equal(expectedCode, resp.Code)

	//------test GetCmdBurnNFT()-------------
	client := ctx.Client
	ctx = val2.ClientCtx
	ctx.Client = client
	args = []string{
		denomID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, recepient),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.Config.BondDenom, sdk.NewInt(10))).String()),
	}

	out, err = clitestutil.ExecTestCLICmd(ctx, cli.GetCmdBurnNFT(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &resp))
	s.Require().Equal(expectedCode, resp.Code)
}
