package bot

import (
	"context"

	commom "github.com/AraanBranco/meepo/cmd/common"
	"github.com/AraanBranco/meepo/internal/config"
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
	Short:   "Starts meepo bot service",
	Example: "meepo start bot -c config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		startBot()
	},
}

func init() {
	BotCmd.Flags().StringVarP(&configPath, "config-path", "c", "config/config.yaml", "path of the configuration YAML file")
}

func startBot() {
	ctx, cancelFn := context.WithCancel(context.Background())

	err, config := commom.ServiceSetup(ctx, cancelFn, logConfig, configPath)
	if err != nil {
		zap.L().With(zap.Error(err)).Fatal("unable to setup service")
	}

	shutdownManagementServerFn := runBot(config)

	<-ctx.Done()

	err = shutdownManagementServerFn()
	if err != nil {
		zap.L().With(zap.Error(err)).Fatal("failed to shutdown management server")
	}
}

func runBot(configs config.Config) func() error {
	// TODO: Implement bot service

	return func() error {
		return nil
	}
}
