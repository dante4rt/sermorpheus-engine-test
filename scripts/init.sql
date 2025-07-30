-- Sermorpheus Engine Database Initialization
-- PostgreSQL Schema for Online Ticket Reservation System

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create database schema (if running as separate script)
-- Note: When using GORM auto-migration, these tables will be created automatically
-- This script is for manual setup or reference

-- Events table
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255) NOT NULL,
    schedule TIMESTAMP WITH TIME ZONE NOT NULL,
    price_idr DECIMAL(15,2) NOT NULL,
    quota INTEGER NOT NULL,
    available_quota INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Customers table
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL REFERENCES customers(id),
    event_id UUID NOT NULL REFERENCES events(id),
    quantity INTEGER NOT NULL,
    total_idr DECIMAL(15,2) NOT NULL,
    usdt_rate DECIMAL(15,6) NOT NULL,
    usdt_amount DECIMAL(15,6) NOT NULL,
    payment_address VARCHAR(42),
    status VARCHAR(20) DEFAULT 'pending',
    payment_locked_at TIMESTAMP WITH TIME ZONE,
    payment_confirmed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tickets table
CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    event_id UUID NOT NULL REFERENCES events(id),
    customer_id UUID NOT NULL REFERENCES customers(id),
    ticket_code VARCHAR(100) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Payment addresses table
CREATE TABLE IF NOT EXISTS payment_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    address VARCHAR(42) UNIQUE NOT NULL,
    private_key TEXT NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- USDT rates table
CREATE TABLE IF NOT EXISTS usdt_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    idr_to_usdt_rate DECIMAL(15,6) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Blockchain transactions table
CREATE TABLE IF NOT EXISTS blockchain_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    tx_hash VARCHAR(66) UNIQUE,
    from_address VARCHAR(42),
    to_address VARCHAR(42),
    amount DECIMAL(15,6),
    confirmations INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_events_deleted_at ON events(deleted_at);
CREATE INDEX IF NOT EXISTS idx_events_schedule ON events(schedule);
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_payment_address ON transactions(payment_address);
CREATE INDEX IF NOT EXISTS idx_tickets_ticket_code ON tickets(ticket_code);
CREATE INDEX IF NOT EXISTS idx_payment_addresses_is_used ON payment_addresses(is_used);
CREATE INDEX IF NOT EXISTS idx_usdt_rates_created_at ON usdt_rates(created_at);
CREATE INDEX IF NOT EXISTS idx_blockchain_transactions_tx_hash ON blockchain_transactions(tx_hash);
CREATE INDEX IF NOT EXISTS idx_blockchain_transactions_status ON blockchain_transactions(status);

-- Insert sample data for testing

-- Sample USDT rate
INSERT INTO usdt_rates (id, idr_to_usdt_rate, created_at) VALUES 
(uuid_generate_v4(), 16394.58, NOW())
ON CONFLICT DO NOTHING;

-- Sample events
INSERT INTO events (id, name, description, location, schedule, price_idr, quota, available_quota, created_at, updated_at) VALUES 
(uuid_generate_v4(), 'Web3 Developer Conference 2025', 'The biggest blockchain and Web3 developer conference in Southeast Asia', 'Jakarta Convention Center', '2025-08-15 09:00:00+07', 500000, 1000, 1000, NOW(), NOW()),
(uuid_generate_v4(), 'Crypto Music Festival', 'Three-day music festival with cryptocurrency payment integration', 'Gelora Bung Karno Stadium', '2025-09-20 16:00:00+07', 750000, 5000, 5000, NOW(), NOW()),
(uuid_generate_v4(), 'DeFi Startup Pitch Day', 'Decentralized finance startup pitch competition', 'Bandung Digital Valley', '2025-08-30 10:00:00+07', 250000, 200, 200, NOW(), NOW()),
(uuid_generate_v4(), 'NFT Art Workshop', 'Learn to create and mint NFT artwork', 'Bali Creative Hub', '2025-09-05 14:00:00+07', 150000, 50, 50, NOW(), NOW()),
(uuid_generate_v4(), 'Blockchain Gaming Summit', 'Explore the future of blockchain gaming', 'Surabaya Tech Center', '2025-10-12 11:00:00+07', 300000, 300, 300, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Sample payment addresses (for testing only - in production these should be generated dynamically)
INSERT INTO payment_addresses (id, address, private_key, is_used, created_at, updated_at) VALUES 
(uuid_generate_v4(), '0x742d35C6634C0532925a3b8D39391d5B5d33CD1D', 'sample_private_key_1', false, NOW(), NOW()),
(uuid_generate_v4(), '0x8b5f9A9c9C8E8A8B9A9c9C8E8A8B9A9c9C8E8A8B', 'sample_private_key_2', false, NOW(), NOW()),
(uuid_generate_v4(), '0x9c8E8A8B9A9c9C8E8A8B9A9c9C8E8A8B9A9c9C8E', 'sample_private_key_3', false, NOW(), NOW()),
(uuid_generate_v4(), '0xA8B9A9c9C8E8A8B9A9c9C8E8A8B9A9c9C8E8A8B9', 'sample_private_key_4', false, NOW(), NOW()),
(uuid_generate_v4(), '0xB9A9c9C8E8A8B9A9c9C8E8A8B9A9c9C8E8A8B9A9c', 'sample_private_key_5', false, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Print completion message
DO $$
BEGIN
    RAISE NOTICE 'Sermorpheus Engine database initialization completed successfully!';
    RAISE NOTICE 'Schema created with sample data for testing.';
    RAISE NOTICE 'Use the seeder command to generate actual payment addresses: go run cmd/seeder/main.go';
END
$$;