package config

import (
	"github.com/goletan/config/pkg"
	logger "github.com/goletan/logger/pkg"
	"github.com/goletan/resilience/internal/types"
	"go.uber.org/zap"
)

var cfg types.ResilienceConfig

func LoadResilienceConfig(log *logger.ZapLogger) (*types.ResilienceConfig, error) {
	if err := config.LoadConfig("Resilience", &cfg, log); err != nil {
		log.Error("Failed to load resilience configuration", zap.Error(err))
		return nil, err
	}

	return &cfg, nil
}
