package bot

import (
	"context"
	"errors"
	"fmt"
	"time"

	commom "github.com/AraanBranco/meepow/cmd/common"
	"github.com/AraanBranco/meepow/internal/config"
	"github.com/AraanBranco/meepow/internal/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	logConfig  string
	configPath string
)

const serviceName string = "bot"

var BotCmd = &cobra.Command{
	Use:     "bot",
	Short:   "Starts meepow bot service",
	Example: "meepow start bot -c config.yaml -l production",
	Run: func(cmd *cobra.Command, args []string) {
		startBot()
	},
}

func init() {
	BotCmd.Flags().StringVarP(&logConfig, "log-config", "l", "production", "preset of configurations used by the logs. possible values are \"development\" or \"production\".")
	BotCmd.Flags().StringVarP(&configPath, "config-path", "c", "config/config.yaml", "path of the configuration YAML file")
}

func startBot() {
	ctx, cancelFn := context.WithCancel(context.Background())

	err, config := commom.ServiceSetup(ctx, cancelFn, logConfig, configPath)
	if err != nil {
		zap.L().With(zap.Error(err)).Fatal("unable to setup service")
	}

	fmt.Println("Starting bot service")
	shutdownBotServerFn := runBot(config)

	waitForShutdown(ctx, shutdownBotServerFn)
}
func runBot(configs config.Config) func() error {
	fmt.Println("Bot service started")
	botManager := service.NewBotManager(configs)

	referenceLobbyID := botManager.Config.GetString("reference.id")
	fmt.Println("Lobby ID: ", referenceLobbyID)

	lobbyData, err := botManager.GetLobbyData(referenceLobbyID)
	if err != nil {
		return logAndReturnError(botManager.Logger, err, "failed to get lobby data")
	}

	if lobbyData.ReferenceID == "" {
		return logAndReturnError(botManager.Logger, errors.New("not_found"), "Lobby not found")
	}

	botManager.Logger.Info("Lobby data", zap.Any("lobbyData", lobbyData))

	botManager.Logger.Info("Starting bot")
	go botManager.StartupBot(lobbyData)

	return func() error {
		shutdownCtx, cancelShutdownFn := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancelShutdownFn()

		zap.L().Info("stopping bot service")
		botManager.SteamClient.Disconnect()
		<-shutdownCtx.Done()
		return shutdownCtx.Err()
	}
}

func logAndReturnError(logger *zap.Logger, err error, message string) func() error {
	logger.With(zap.Error(err)).Fatal(message)
	return func() error {
		return err
	}
}

func waitForShutdown(ctx context.Context, shutdownFn func() error) {
	<-ctx.Done()

	err := shutdownFn()
	if err != nil {
		zap.L().With(zap.Error(err)).Fatal("failed to shutdown bot server")
	}
}
