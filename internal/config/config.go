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
	USDTContract       string
	PlatformFeePercent float64
}

func Load() *Config {
	platformFee, _ := strconv.ParseFloat(getEnv("PLATFORM_FEE_PERCENT", "1.2"), 64)

	return &Config{
		Port:               getEnv("PORT", "8080"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/sermorpheus?sslmode=disable"),
		BSCRPCUrl:          getEnv("BSC_RPC_URL", "https://bsc-testnet.public.blastapi.io"),
		USDTContract:       getEnv("USDT_CONTRACT", "0xCD60747D9Bbb1da2AfB2F834391f0FF6ccb15f1a"),
		PlatformFeePercent: platformFee,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
