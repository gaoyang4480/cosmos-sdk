package genutil

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ sdk.AppModule      = AppModule{}
	_ sdk.AppModuleBasic = AppModuleBasic{}
)

// module name
const ModuleName = "genutil"

// app module basics object
type AppModuleBasic struct {
	gvi GenesisValidatorsInfo
}

// NewAppModuleBasic creates a new AppModuleBasic object
func NewAppModuleBasic(gvi GenesisValidatorsInfo) AppModuleBasic {
	return AppModuleBasic{
		gvi: gvi,
	}
}

// module name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {}

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return moduleCdc.MustMarshalJSON(GenesisState{})
}

// module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := moduleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data, gvi)
}

//___________________________
// app module
type AppModule struct {
	AppModuleBasic
	accountKeeper AccountKeeper
	stakingKeeper StakingKeeper
	deliverTx     deliverTxfn
}

// NewAppModule creates a new AppModule object
func NewAppModule(accountKeeper AccountKeeper,
	stakingKeeper StakingKeeper, deliverTx deliverTxfn) AppModule {

	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		accountKeeper:  accountKeeper,
		stakingKeeper:  stakingKeeper,
		deliverTx:      deliverTx,
	}
}

// register invariants
func (AppModule) RegisterInvariants(_ sdk.InvariantRouter) {}

// module message route name
func (AppModule) Route() string { return "" }

// module handler
func (AppModule) NewHandler() sdk.Handler { return nil }

// module querier route name
func (AppModule) QuerierRoute() string { return "" }

// module querier
func (am AppModule) NewQuerierHandler() sdk.Querier { return nil }

// module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	moduleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, moduleCdc, am.stakingKeeper, am.deliverTx, genesisState)
}

// module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return nil
}

// module begin-block
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) sdk.Tags {
	return sdk.EmptyTags()
}

// module end-block
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) ([]abci.ValidatorUpdate, sdk.Tags) {
	return []abci.ValidatorUpdate{}, sdk.EmptyTags()
}
