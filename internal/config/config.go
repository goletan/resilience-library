// /resilience/internal/config/config.go
package config

import (
	"github.com/goletan/config"
	"github.com/goletan/resilience/internal/types"
	"go.uber.org/zap"
)

var cfg types.ResilienceConfig

func LoadResilienceConfig(logger *zap.Logger) (*types.ResilienceConfig, error) {
	if err := config.LoadConfig("Resilience", &cfg, logger); err != nil {
		logger.Error("Failed to load resilience configuration", zap.Error(err))
		return nil, err
	}

	return &cfg, nil
}
