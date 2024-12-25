package commom

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AraanBranco/meepo/internal/config"
	"github.com/AraanBranco/meepo/internal/config/viper"
	"github.com/AraanBranco/meepo/internal/service"
	"go.uber.org/zap"
)

func ServiceSetup(ctx context.Context, cancelFn context.CancelFunc, logConfig, configPath string) (error, config.Config) {
	err := service.ConfigureLogging(logConfig)
	if err != nil {
		return fmt.Errorf("unable to configure logging: %w", err), nil
	}

	viperConfig, err := viper.NewViperConfig(configPath)
	if err != nil {
		return fmt.Errorf("unable to load config: %w", err), nil
	}

	launchTerminatingListenerGoroutine(cancelFn)

	return nil, viperConfig
}

func launchTerminatingListenerGoroutine(cancelFunc context.CancelFunc) {
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		<-sigs
		zap.L().Info("received termination")

		cancelFunc()
	}()
}
