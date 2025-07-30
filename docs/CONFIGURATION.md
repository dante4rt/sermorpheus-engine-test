# Environment Configuration Guide

## Configuration Overview

The **sermorpheus-engine-test** application uses environment variables for configuration, making it flexible across different deployment environments (development, staging, production).

## Configuration Files

### `.env.example`

Template file showing all available configuration options with sensible defaults.

### `.env`

Your actual environment configuration (not committed to git for security).

## Available Configuration Options

### Server Configuration

```bash
PORT=8080                    # Server port
ENVIRONMENT=development      # Environment mode (development/staging/production)
```

### Database Configuration

```bash
DATABASE_URL=postgres://user:password@localhost:5432/sermorpheus?sslmode=disable
```

### Blockchain Configuration

```bash
# BSC Network Settings
BSC_RPC_URL=https://data-seed-prebsc-1-s1.binance.org:8545
BSC_WSS_URL=wss://bsc-testnet.drpc.org

# USDT Contract Settings  
USDT_CONTRACT=0xCD60747D9Bbb1da2AfB2F834391f0FF6ccb15f1a
USDT_DECIMALS=6

# Payment Processing
PLATFORM_FEE_PERCENT=1.2    # Platform fee percentage
```

## Network Configurations

### BSC Testnet (Default)

```bash
BSC_RPC_URL=https://data-seed-prebsc-1-s1.binance.org:8545
USDT_CONTRACT=0xCD60747D9Bbb1da2AfB2F834391f0FF6ccb15f1a
USDT_DECIMALS=6
```

### Alternative BSC Testnet Endpoints

```bash
# Option 1: Binance official
BSC_RPC_URL=https://data-seed-prebsc-1-s1.binance.org:8545

# Option 2: Public API
BSC_RPC_URL=https://bsc-testnet.public.blastapi.io

# Option 3: QuickNode (requires API key)
BSC_RPC_URL=https://your-endpoint.bsc-testnet.quiknode.pro/your-api-key/
```

## Setup Instructions

### 1. Copy Environment Template

```bash
cp .env.example .env
```

### 2. Configure Your Environment

Edit `.env` with your specific settings:

```bash
# Update database URL
DATABASE_URL=postgres://youruser:yourpass@localhost:5432/yourdatabase?sslmode=disable

# Choose your BSC RPC endpoint
BSC_RPC_URL=https://your-preferred-bsc-endpoint.com

# Adjust platform fee if needed
PLATFORM_FEE_PERCENT=1.5
```

### 3. Verify Configuration

```bash
# Test compilation
go build ./cmd/server

# Test startup (will show config values in logs)
go run cmd/server/main.go
```

## Configuration Validation

The application validates configuration at startup:

- **BSC RPC URL**: Tests connection to blockchain network
- **Database URL**: Validates database connectivity  
- **USDT Contract**: Verifies contract address format
- **Platform Fee**: Ensures reasonable percentage (0-10%)

## Environment-Specific Configurations

### Development

```bash
ENVIRONMENT=development
BSC_RPC_URL=https://data-seed-prebsc-1-s1.binance.org:8545  # Testnet
USDT_CONTRACT=0xCD60747D9Bbb1da2AfB2F834391f0FF6ccb15f1a      # Testnet USDT
```

## Security Best Practices

1. **Never commit `.env`** - Add to `.gitignore`
2. **Use different contracts** for testnet vs mainnet
3. **Validate RPC endpoints** before production use
4. **Monitor platform fee** settings
5. **Use secure database credentials**

## Configuration Architecture

The configuration follows a clean architecture pattern:

```
├── internal/config/config.go    # Configuration structure & loading
├── .env.example                 # Template with defaults
├── .env                        # Your actual config (gitignored)
└── cmd/server/main.go          # Config injection into services
```

This approach ensures:

- ✅ **No hardcoded values** in business logic
- ✅ **Environment-specific** configurations  
- ✅ **Secure credential** management
- ✅ **Easy deployment** across environments
- ✅ **Testable configuration** with defaults
