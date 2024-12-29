package config

import (
	"github.com/goletan/config-library/pkg"
	logger "github.com/goletan/logger-library/pkg"
	"github.com/goletan/resilience-library/internal/types"
	"go.uber.org/zap"
)

func LoadResilienceConfig(log *logger.ZapLogger) (*types.ResilienceConfig, error) {
	var cfg types.ResilienceConfig

	if err := config.LoadConfig("Resilience", &cfg, log); err != nil {
		log.Error("Failed to load resilience-library configuration", zap.Error(err))
		return nil, err
	}

	return &cfg, nil
}
