package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/headersync/internal/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/codec"
	mctype "github.com/ontio/multi-chain/core/types"
	mcc "github.com/ontio/multi-chain/common"
)

var _ Keeper = (*BaseKeeper)(nil)

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	BaseViewKeeper
	SyncGenesisHeader(ctx sdk.Context, genesisHeader []byte) sdk.Error
	SyncBlockHeaders(ctx sdk.Context, headers [][]byte) sdk.Error
	ProcessHeader(ctx sdk.Context, header *mctype.Header) sdk.Error
}


// BaseKeeper manages transfers between accounts. It implements the Keeper interface.
type BaseKeeper struct {
	cdc 	*codec.Codec
	storeKey      sdk.StoreKey
	paramSpace    params.Subspace
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper( cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,) BaseKeeper {

	ps := paramSpace.WithKeyTable(types.ParamKeyTable())
	return BaseKeeper{
		cdc: cdc,
		storeKey:             key,
		paramSpace:     ps,
	}
}

// DelegateCoins performs delegation by deducting amt coins from an account with
// address addr. For vesting accounts, delegations amounts are tracked for both
// vesting and vested coins.
// The coins are then transferred from the delegator address to a ModuleAccount address.
// If any of the delegation amounts are negative, an error is returned.
func (keeper BaseKeeper) SyncGenesisHeader(ctx sdk.Context, genesisHeaderBytes []byte) sdk.Error {
	genesisHeader := &mctype.Header{}
	//source := mcc.NewZeroCopySource(genesisHeaderBytes)
	//err := genesisHeader.Deserialization(source)
	//if err != nil {
	//	return sdk.ErrInternal(fmt.Sprintf("sync multichain genesis header, deserialize header error"))
	//}
	source := mcc.NewZeroCopySource(genesisHeaderBytes)
	if err := genesisHeader.Deserialization(source); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("GenesisHeader deserialization err:%v", err))
	}

	if err := keeper.SetBlockHeader(ctx, genesisHeader); err != nil {
		return err
	}
	if err := keeper.UpdateConsensusPeer(ctx, genesisHeader); err != nil {
		return err
	}
	return nil
}


func (keeper BaseKeeper) SyncBlockHeaders(ctx sdk.Context, headers [][]byte) sdk.Error {
	for _, headerBytes := range headers {
		header := &mctype.Header{}
		source := mcc.NewZeroCopySource(headerBytes)
		if err := header.Deserialization(source); err != nil {
			return sdk.ErrInternal(fmt.Sprintf("BlockHeader deserialization err:%v", err))
		}
		h, err := keeper.GetHeaderByHeight(ctx, header.ChainID, header.Height)
		if err == nil {
			continue
		}

		if h == nil {
			if err := keeper.ProcessHeader(ctx, header); err != nil {
				return sdk.ErrInternal(fmt.Sprintf("SyncBlockHeader error:%s", err))
			}
		}


		//if err := keeper.VerifyHeader(ctx, header); err != nil {
		//
		//}
	}
	return nil
}


func (keeper BaseKeeper) ProcessHeader(ctx sdk.Context, header *mctype.Header) sdk.Error {
	if err := keeper.VerifyHeader(ctx, header); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("processHeader, %v", err))
	}
	if err := keeper.SetBlockHeader(ctx, header); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("processHeader, %v", err))
	}
	if err := keeper.UpdateConsensusPeer(ctx, header); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("processHeader, %v", err))
	}
	return nil
}



// BaseViewKeeper implements a read only keeper implementation of ViewKeeper.
type BaseViewKeeper interface {

	GetHeaderByHeight(ctx sdk.Context, chainId uint64, height uint32) (*mctype.Header, sdk.Error)
	GetHeaderByHash(ctx sdk.Context, chainId uint64, hash mcc.Uint256) (*mctype.Header, sdk.Error)

	GetConsensusPeers(ctx sdk.Context, chainId uint64, height uint32) (*types.ConsensusPeers, sdk.Error)
	GetKeyHeights(ctx sdk.Context, chainId uint64) (*types.KeyHeights, sdk.Error)
}

