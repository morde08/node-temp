package observer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// InitGenesis initializes the observer module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	genesisObservers := genState.Observers
	observerCount := uint64(0)
	for _, mapper := range genesisObservers {
		if mapper != nil {
			k.SetObserverMapper(ctx, mapper)
			observerCount += uint64(len(mapper.ObserverList))
		}
	}

	// If core params are defined set them, otherwise set default
	if len(genState.CoreParamsList.CoreParams) > 0 {
		k.SetCoreParams(ctx, genState.CoreParamsList)
	} else {
		k.SetCoreParams(ctx, types.GetCoreParams())
	}

	// Set all the nodeAccount
	for _, elem := range genState.NodeAccountList {
		if elem != nil {
			k.SetNodeAccount(ctx, *elem)
		}
	}

	params := types.DefaultParams()
	if genState.Params != nil {
		params = *genState.Params
	}
	k.SetParams(ctx, params)

	// Set if defined

	crosschainFlags := types.DefaultCrosschainFlags()
	if genState.CrosschainFlags != nil {
		crosschainFlags.IsOutboundEnabled = genState.CrosschainFlags.IsOutboundEnabled
		crosschainFlags.IsInboundEnabled = genState.CrosschainFlags.IsInboundEnabled
		if genState.CrosschainFlags.BlockHeaderVerificationFlags != nil {
			crosschainFlags.BlockHeaderVerificationFlags = genState.CrosschainFlags.BlockHeaderVerificationFlags
		}
		if genState.CrosschainFlags.GasPriceIncreaseFlags != nil {
			crosschainFlags.GasPriceIncreaseFlags = genState.CrosschainFlags.GasPriceIncreaseFlags
		}
		k.SetCrosschainFlags(ctx, *crosschainFlags)
	} else {
		k.SetCrosschainFlags(ctx, *types.DefaultCrosschainFlags())
	}

	// Set if defined
	if genState.Keygen != nil {
		k.SetKeygen(ctx, *genState.Keygen)
	}

	ballotListForHeight := make(map[int64][]string)
	if len(genState.Ballots) > 0 {
		for _, ballot := range genState.Ballots {
			if ballot != nil {
				k.SetBallot(ctx, ballot)
				ballotListForHeight[ballot.BallotCreationHeight] = append(ballotListForHeight[ballot.BallotCreationHeight], ballot.BallotIdentifier)
			}
		}
	}

	for height, ballotList := range ballotListForHeight {
		k.SetBallotList(ctx, &types.BallotListForHeight{
			Height:           height,
			BallotsIndexList: ballotList,
		})
	}

	if genState.LastObserverCount != nil {
		k.SetLastObserverCount(ctx, genState.LastObserverCount)
	} else {
		k.SetLastObserverCount(ctx, &types.LastObserverCount{LastChangeHeight: 0, Count: observerCount})
	}

	if genState.Tss != nil {
		tss := *genState.Tss
		k.SetTSS(ctx, tss)
		for _, chain := range common.DefaultChainsList() {
			k.SetPendingNonces(ctx, types.PendingNonces{
				NonceLow:  0,
				NonceHigh: 0,
				ChainId:   chain.ChainId,
				Tss:       tss.TssPubkey,
			})
		}
	}

	// Get all chain nonces

	for _, elem := range genState.TssHistory {
		k.SetTSSHistory(ctx, elem)
	}

	for _, elem := range genState.TssFundMigrators {
		k.SetFundMigrator(ctx, elem)
	}

	for _, elem := range genState.BlameList {
		k.SetBlame(ctx, elem)
	}

	// Set all the pending nonces
	for _, pendingNonce := range genState.PendingNonces {
		k.SetPendingNonces(ctx, pendingNonce)
	}
	for _, chainNonce := range genState.ChainNonces {
		k.SetChainNonces(ctx, chainNonce)
	}
	for _, elem := range genState.NonceToCctx {
		k.SetNonceToCctx(ctx, elem)
	}

}

// ExportGenesis returns the observer module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)

	coreParams, found := k.GetAllCoreParams(ctx)
	if !found {
		coreParams = types.CoreParamsList{}
	}

	// Get all node accounts
	nodeAccountList := k.GetAllNodeAccount(ctx)
	nodeAccounts := make([]*types.NodeAccount, len(nodeAccountList))
	for i, elem := range nodeAccountList {
		elem := elem
		nodeAccounts[i] = &elem
	}

	// Get all crosschain flags
	cf := types.DefaultCrosschainFlags()
	crosschainFlags, found := k.GetCrosschainFlags(ctx)
	if found {
		cf = &crosschainFlags
	}

	kn := &types.Keygen{}
	keygen, found := k.GetKeygen(ctx)
	if found {
		kn = &keygen
	}

	oc := &types.LastObserverCount{}
	observerCount, found := k.GetLastObserverCount(ctx)
	if found {
		oc = &observerCount
	}

	// Get tss
	tss := &types.TSS{}
	t, found := k.GetTSS(ctx)
	if found {
		tss = &t
	}

	var pendingNonces []types.PendingNonces
	p, err := k.GetAllPendingNonces(ctx)
	if err == nil {
		pendingNonces = p
	}
	return &types.GenesisState{
		Ballots:           k.GetAllBallots(ctx),
		Observers:         k.GetAllObserverMappers(ctx),
		CoreParamsList:    coreParams,
		Params:            &params,
		NodeAccountList:   nodeAccounts,
		CrosschainFlags:   cf,
		Keygen:            kn,
		LastObserverCount: oc,
		Tss:               tss,
		PendingNonces:     pendingNonces,
		TssHistory:        k.GetAllTSS(ctx),
		TssFundMigrators:  k.GetAllTssFundMigrators(ctx),
		BlameList:         k.GetAllBlame(ctx),
		ChainNonces:       k.GetAllChainNonces(ctx),
		NonceToCctx:       k.GetAllNonceToCctx(ctx),
	}
}
