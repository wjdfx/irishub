package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/irisnet/irishub/app/protocol"
	v0 "github.com/irisnet/irishub/app/v0"
	v1 "github.com/irisnet/irishub/app/v1"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	appName                = "IrisApp"
	appPrometheusNamespace = "iris"
	FlagReplay             = "replay-last-block"
	DefaultSyncableHeight  = store.NumStoreEvery // Multistore saves a snapshot every 10000 blocks
	DefaultCacheSize       = store.NumRecent     // Multistore saves last 100 blocks
)

// default home directories for expected binaries
var (
	DefaultLCDHome  = os.ExpandEnv("$HOME/.irislcd")
	DefaultCLIHome  = os.ExpandEnv("$HOME/.iriscli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.iris")
)

// Extended ABCI application
type IrisApp struct {
	*BaseApp
}

func NewIrisApp(logger log.Logger, db dbm.DB, config *cfg.InstrumentationConfig, traceStore io.Writer, baseAppOptions ...func(*BaseApp)) *IrisApp {
	bApp := NewBaseApp(appName, logger, db, baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)

	// create your application object
	var app = &IrisApp{BaseApp: bApp}

	protocolKeeper := sdk.NewProtocolKeeper(protocol.KeyMain)
	engine := protocol.NewProtocolEngine(protocolKeeper)
	app.SetProtocolEngine(&engine)
	app.MountStoresIAVL(engine.GetKVStoreKeys())
	app.MountStoresTransient(engine.GetTransientStoreKeys())

	var err error
	if viper.GetBool(FlagReplay) {
		lastHeight := Replay(app.Logger)
		err = app.LoadVersion(lastHeight, protocol.KeyMain, true)
	} else {
		err = app.LoadLatestVersion(protocol.KeyMain)
	} // app is now sealed
	if err != nil {
		cmn.Exit(err.Error())
	}
	//Duplicate prometheus config
	appPrometheusConfig := *config
	//Change namespace to appName
	appPrometheusConfig.Namespace = appPrometheusNamespace
	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, app.checkInvariant, app.trackCoinFlow, &appPrometheusConfig))
	engine.Add(v1.NewProtocolV1(1, logger, protocolKeeper, app.checkInvariant, app.trackCoinFlow, &appPrometheusConfig))
	// engine.Add(v1.NewProtocolV1(1, ...))
	// engine.Add(v2.NewProtocolV1(2, ...))

	loaded, current := engine.LoadCurrentProtocol(app.GetKVStore(protocol.KeyMain))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!", current))
	}
	app.BaseApp.txDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())
	engine.GetCurrentProtocol().InitMetrics(app.cms)
	return app
}

// latest version of codec
func MakeLatestCodec() *codec.Codec {
	var cdc = v1.MakeCodec() // replace with latest protocol version
	return cdc
}

func (app *IrisApp) LastBlockHeight() int64 {
	return app.BaseApp.LastBlockHeight()
}

func (app *IrisApp) ResetOrReplay(replayHeight int64) (replay bool, height int64) {
	lastBlockHeight := app.BaseApp.LastBlockHeight()
	if replayHeight > lastBlockHeight {
		replayHeight = lastBlockHeight
	}

	app.Logger.Info("This Reset operation will change the application store, backup your node home directory before proceeding!!!")
	app.Logger.Info(fmt.Sprintf("The last block height is %v, will reset height to %v.", lastBlockHeight, replayHeight))

	app.Logger.Info("Are you sure to proceed? (y/n)")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		cmn.Exit(err.Error())
	}
	confirm := strings.ToLower(strings.TrimSpace(input))
	if confirm != "y" && confirm != "yes" {
		cmn.Exit("Reset operation aborted.")
	}

	if lastBlockHeight-replayHeight <= DefaultCacheSize {
		err := app.LoadVersion(replayHeight, protocol.KeyMain, true)

		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("wanted to load target %v but only found up to", replayHeight)) {
				app.Logger.Info(fmt.Sprintf("Can not find the target version %d, trying to load an earlier version and replay blocks", replayHeight))
			} else {
				cmn.Exit(err.Error())
			}
		} else {
			app.Logger.Info(fmt.Sprintf("The last block height is %d, loaded store at %d", lastBlockHeight, replayHeight))
			return false, replayHeight
		}
	}

	loadHeight := app.replayToHeight(replayHeight, app.Logger)
	err = app.LoadVersion(loadHeight, protocol.KeyMain, true)
	if err != nil {
		cmn.Exit(err.Error())
	}

	// If reset to another protocol version, should reload Protocol and reset txDecoder
	loaded, current := app.Engine.LoadCurrentProtocol(app.GetKVStore(protocol.KeyMain))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!", current))
	}
	app.BaseApp.txDecoder = auth.DefaultTxDecoder(app.Engine.GetCurrentProtocol().GetCodec())

	app.Logger.Info(fmt.Sprintf("The last block height is %d, want to load store at %d", lastBlockHeight, replayHeight))

	// Version 1 does not need replay
	if replayHeight == 1 {
		app.Logger.Info(fmt.Sprintf("Loaded store at %d", loadHeight))
		return false, replayHeight
	}

	app.Logger.Info(fmt.Sprintf("Loaded store at %d, start to replay to %d", loadHeight, replayHeight))
	return true, replayHeight

}

// export the state of iris for a genesis file
func (app *IrisApp) ExportAppStateAndValidators(forZeroHeight bool) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	return app.Engine.GetCurrentProtocol().ExportAppStateAndValidators(ctx, forZeroHeight)
}

func (app *IrisApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, protocol.KeyMain, false)
}

func (app *IrisApp) replayToHeight(replayHeight int64, logger log.Logger) int64 {
	loadHeight := int64(0)
	logger.Info("Please make sure the replay height is smaller than the latest block height.")
	if replayHeight >= DefaultSyncableHeight {
		loadHeight = replayHeight - replayHeight%DefaultSyncableHeight
	} else {
		// version 1 will always be kept
		loadHeight = 1
	}
	return loadHeight
}
