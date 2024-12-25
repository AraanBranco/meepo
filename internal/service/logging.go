package service

import (
	"fmt"
	"net/http"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func ConfigureLogging(configPreset string) error {
	var cfg zap.Config
	switch configPreset {
	case "development":
		cfg = zap.NewDevelopmentConfig()
	case "production":
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	default:
		return fmt.Errorf("unexpected log_config: %v", configPreset)
	}

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zap.L().Info("Request: ", zap.String("Path:", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
