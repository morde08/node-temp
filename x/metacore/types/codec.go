package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgGasBalanceVoter{}, "metacore/GasBalanceVoter", nil)

	cdc.RegisterConcrete(&MsgGasPriceVoter{}, "metacore/GasPriceVoter", nil)

	cdc.RegisterConcrete(&MsgNonceVoter{}, "metacore/NonceVoter", nil)

	cdc.RegisterConcrete(&MsgReceiveConfirmation{}, "metacore/ReceiveConfirmation", nil)

	cdc.RegisterConcrete(&MsgSendVoter{}, "metacore/SendVoter", nil)

	cdc.RegisterConcrete(&MsgSetNodeKeys{}, "metacore/SetNodeKeys", nil)

}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgGasBalanceVoter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgGasPriceVoter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNonceVoter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgReceiveConfirmation{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSendVoter{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetNodeKeys{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	//amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
