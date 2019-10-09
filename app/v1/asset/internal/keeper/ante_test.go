package keeper

import (
	"testing"

	"github.com/irisnet/irishub/app/v1/asset/internal/types"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

// TestAssetAnteHandler tests the ante handler of asset
func TestAssetAnteHandler(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, paramsKeeper.Subspace(types.DefaultParamSpace))

	// init params
	keeper.Init(ctx)

	// set test accounts
	addr1 := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr2"))
	acc1 := ak.NewAccountWithAddress(ctx, addr1)
	acc2 := ak.NewAccountWithAddress(ctx, addr2)
	baseAcc1, _ := acc1.(*sdk.BaseAccount)
	baseAcc2, _ := acc2.(*sdk.BaseAccount)

	// get asset fees
	gatewayCreateFee := GetGatewayCreateFee(ctx, keeper, "mon")
	nativeTokenIssueFee := GetTokenIssueFee(ctx, keeper, "sym")
	gatewayTokenIssueFee := GetGatewayTokenIssueFee(ctx, keeper, "sym")
	nativeTokenMintFee := GetTokenMintFee(ctx, keeper, "sym")

	// construct msgs
	msgCreateGateway := types.NewMsgCreateGateway(addr1, "mon", "i", "d", "w")
	msgIssueNativeToken := types.MsgIssueToken{Owner: addr1, Source: types.AssetSource(0x00), Symbol: "sym"}
	msgIssueGatewayToken := types.MsgIssueToken{Owner: addr2, Source: types.AssetSource(0x02), Symbol: "sym"}
	msgMintNativeToken := types.MsgMintToken{Owner: addr2, TokenId: "i.sym"}
	msgNonAsset1 := sdk.NewTestMsg(addr1)
	msgNonAsset2 := sdk.NewTestMsg(addr2)

	// construct test txs
	tx1 := auth.StdTx{Msgs: []sdk.Msg{msgCreateGateway, msgIssueNativeToken, msgIssueGatewayToken, msgMintNativeToken}}
	tx2 := auth.StdTx{Msgs: []sdk.Msg{msgCreateGateway, msgIssueNativeToken, msgNonAsset1, msgIssueGatewayToken, msgMintNativeToken}}
	tx3 := auth.StdTx{Msgs: []sdk.Msg{msgNonAsset2, msgCreateGateway, msgIssueNativeToken, msgIssueGatewayToken, msgMintNativeToken}}

	// set signers and construct an ante handler
	ctx = ctx.WithKeySignerAccs([]*sdk.BaseAccount{baseAcc1, baseAcc2})
	anteHandler := NewAnteHandler(keeper)

	// assert that the ante handler will return with `abort` set to true
	_ = acc1.SetCoins(sdk.Coins{gatewayCreateFee.Add(nativeTokenIssueFee)})
	_, res, abort := anteHandler(ctx, tx1, false)
	require.Equal(t, true, abort)
	require.Equal(t, false, res.IsOK())

	// assert that the ante handler will return with `abort` set to true
	_ = acc1.SetCoins(acc1.GetCoins().Add(sdk.Coins{gatewayTokenIssueFee}))
	_, res, abort = anteHandler(ctx, tx1, false)
	require.Equal(t, true, abort)
	require.Equal(t, false, res.IsOK())

	// assert that the ante handler will return with `abort` set to false
	_ = acc1.SetCoins(acc1.GetCoins().Add(sdk.Coins{nativeTokenMintFee}))
	_, res, abort = anteHandler(ctx, tx1, false)
	require.Equal(t, false, abort)
	require.Equal(t, true, res.IsOK())

	// assert that the ante handler will return with `abort` set to false
	_ = acc1.SetCoins(sdk.Coins{gatewayCreateFee.Add(nativeTokenIssueFee)})
	_, res, abort = anteHandler(ctx, tx2, false)
	require.Equal(t, false, abort)
	require.Equal(t, true, res.IsOK())

	// assert that the ante handler will return with `abort` set to false
	_ = acc1.SetCoins(sdk.Coins{})
	_, res, abort = anteHandler(ctx, tx3, false)
	require.Equal(t, false, abort)
	require.Equal(t, true, res.IsOK())

	// assert that the ante handler will return with `abort` set to true
	ctx = ctx.WithKeySignerAccs([]*sdk.BaseAccount{})
	_, res, abort = anteHandler(ctx, tx3, false)
	require.Equal(t, true, abort)
	require.Equal(t, false, res.IsOK())
}
