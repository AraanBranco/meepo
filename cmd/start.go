package cmd

import (
	"context"
	"fmt"
	"net/http"

	commom "github.com/AraanBranco/meepo/cmd/common"
	"github.com/AraanBranco/meepo/internal/api/handlers"
	"github.com/AraanBranco/meepo/internal/config"
	"github.com/AraanBranco/meepo/internal/service"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	logConfig  string
	configPath string
)

const serviceName string = "api"

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start the application",
	Example: "meepo start -c config.yaml -l production",
	Run: func(cmd *cobra.Command, args []string) {
		runApi()
	},
}

func init() {
	startCmd.Flags().StringVarP(&logConfig, "log-config", "l", "production", "preset of configurations used by the logs. possible values are \"development\" or \"production\".")
	startCmd.Flags().StringVarP(&configPath, "config-path", "c", "config/config.yaml", "path of the configuration YAML file")
}

func runApi() {
	ctx, cancelFn := context.WithCancel(context.Background())

	err, config := commom.ServiceSetup(ctx, cancelFn, logConfig, configPath)
	if err != nil {
		zap.L().With(zap.Error(err)).Fatal("unable to setup service")
	}

	router := mux.NewRouter()

	shutdownManagementServerFn := runServer(config, router)

	<-ctx.Done()

	err = shutdownManagementServerFn()
	if err != nil {
		zap.L().With(zap.Error(err)).Fatal("failed to shutdown management server")
	}
}

func runServer(configs config.Config, r *mux.Router) func() error {

	// Apis
	r.HandleFunc("/", handlers.Default).Methods(http.MethodGet)

	r.HandleFunc("/new-lobby", func(w http.ResponseWriter, r *http.Request) {
		handlers.NewLobby(w, r, configs)
	}).Methods(http.MethodPost)

	r.HandleFunc("/status-lobby/{referenceId}", func(w http.ResponseWriter, r *http.Request) {
		referenceID := mux.Vars(r)["referenceId"]
		handlers.StatusLobby(w, r, configs, referenceID)
	}).Methods(http.MethodGet)

	// Mdlws
	r.Use(service.LoggingMiddleware)
	r.Use(service.JsonMiddleware)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", configs.GetString("api.port")),
		Handler: r,
	}

	go func() {
		zap.L().Info(fmt.Sprintf("started HTTP management server at :%s", configs.GetString("api.port")))
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			zap.L().With(zap.Error(err)).Fatal("failed to start HTTP management server")
		}
	}()

	return func() error {
		shutdownCtx, cancelShutdownFn := context.WithTimeout(context.Background(), configs.GetDuration("api.gracefulShutdownTimeout"))
		defer cancelShutdownFn()

		zap.L().Info("stopping HTTP management server")
		return httpServer.Shutdown(shutdownCtx)
	}
}
