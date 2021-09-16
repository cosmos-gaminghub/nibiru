package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgIssueDenom{}, "cosmminden/MsgIssueDenom", nil)
	cdc.RegisterConcrete(&MsgMintNFT{}, "cosmminden/MsgMintNFT", nil)
	cdc.RegisterConcrete(&MsgEditNFT{}, "cosmminden/MsgEditNFT", nil)
	cdc.RegisterConcrete(&MsgTransferNFT{}, "cosmminden/MsgTransferNFT", nil)
	cdc.RegisterConcrete(&MsgBurnNFT{}, "cosmminden/MsgBurnNFT", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgIssueDenom{},
		&MsgMintNFT{},
		&MsgEditNFT{},
		&MsgTransferNFT{},
		&MsgBurnNFT{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
