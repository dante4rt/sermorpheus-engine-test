package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port               string
	Environment        string
	DatabaseURL        string
	BSCRPCUrl          string
	BSCWebSocketURL    string
	USDTContract       string
	USDTDecimals       int
	PlatformFeePercent float64
}

func Load() *Config {
	platformFee, _ := strconv.ParseFloat(getEnv("PLATFORM_FEE_PERCENT", "1.2"), 64)
	usdtDecimals, _ := strconv.Atoi(getEnv("USDT_DECIMALS", "6"))

	return &Config{
		Port:               getEnv("PORT", "8080"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/sermorpheus?sslmode=disable"),
		BSCRPCUrl:          getEnv("BSC_RPC_URL", "https://data-seed-prebsc-1-s1.binance.org:8545"),
		BSCWebSocketURL:    getEnv("BSC_WSS_URL", "wss://bsc-testnet.drpc.org"),
		USDTContract:       getEnv("USDT_CONTRACT", "0xCD60747D9Bbb1da2AfB2F834391f0FF6ccb15f1a"),
		USDTDecimals:       usdtDecimals,
		PlatformFeePercent: platformFee,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
