package config

import (
	"github.com/goletan/config/pkg"
	observability "github.com/goletan/observability/pkg"
	"github.com/goletan/resilience/internal/types"
	"go.uber.org/zap"
)

var cfg types.ResilienceConfig

func LoadResilienceConfig(obs *observability.Observability) (*types.ResilienceConfig, error) {
	if err := config.LoadConfig("Resilience", &cfg, obs); err != nil {
		obs.Logger.Error("Failed to load resilience configuration", zap.Error(err))
		return nil, err
	}

	return &cfg, nil
}
