package ratelimiter

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/goletan/config"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var (
	rateLimiterInstance *rate.Limiter
	once                sync.Once
	logger              *zap.Logger
)

// Init initializes the rate limiter module by utilizing the configuration library.
// This function loads the relevant configuration parameters and sets up the rate limiter accordingly.
func Init(configFile string, log *zap.Logger) {
	logger = log
	var cfg ResilienceConfig

	configPathsEnv := os.Getenv("RESILIENCE_CONFIG_PATHS")
	var configPaths []string
	if configPathsEnv != "" {
		configPaths = strings.Split(configPathsEnv, ",")
	} else {
		configPaths = []string{"."}
	}

	// Load the configuration
	err := config.LoadConfig("resilience", configPaths, &cfg, logger)
	if err != nil {
		logger.Fatal("Failed to load rate limiter configuration", zap.Error(err))
	}

	InitRateLimiter(cfg.RateLimiter.RPS, cfg.RateLimiter.Burst)
}

// ResilienceConfig encapsulates the configuration parameters for the rate limiter.
type ResilienceConfig struct {
	RateLimiter struct {
		RPS   int `mapstructure:"rps"`
		Burst int `mapstructure:"burst"`
	} `mapstructure:"rate_limiter"`
}

// InitRateLimiter initializes the rate limiter to regulate the rate of incoming requests.
// This function employs a "once" mechanism to ensure that the rate limiter is only initialized once.
func InitRateLimiter(rps int, burst int) {
	once.Do(func() {
		rateLimiterInstance = rate.NewLimiter(rate.Limit(rps), burst)
		logger.Info("Rate limiter initialized", zap.Int("rps", rps), zap.Int("burst", burst))
	})
}

// GetRateLimiterInstance guarantees that the rate limiter is properly initialized and returns its instance.
func GetRateLimiterInstance() *rate.Limiter {
	return rateLimiterInstance
}

// ExecuteWithRateLimiting enforces rate limiting when executing the provided operation.
// If the rate limit is exceeded, it logs a warning and returns an error; otherwise, it proceeds to execute the operation.
func ExecuteWithRateLimiting(ctx context.Context, operation func() error) error {
	limiter := GetRateLimiterInstance()
	if err := limiter.Wait(ctx); err != nil {
		logger.Warn("Rate limit exceeded", zap.Error(err))
		return err
	}
	return operation()
}
