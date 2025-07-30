-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create some sample payment addresses for testing
INSERT INTO payment_addresses (id, address, private_key, is_used, created_at, updated_at) VALUES 
(uuid_generate_v4(), '0x742d35C6634C0532925a3b8D39391d5B5d33CD1D', 'sample_private_key_1', false, NOW(), NOW()),
(uuid_generate_v4(), '0x8b5f9A9c9C8E8A8B9A9c9C8E8A8B9A9c9C8E8A8B', 'sample_private_key_2', false, NOW(), NOW()),
(uuid_generate_v4(), '0x9c8E8A8B9A9c9C8E8A8B9A9c9C8E8A8B9A9c9C8E', 'sample_private_key_3', false, NOW(), NOW()),
(uuid_generate_v4(), '0xA8B9A9c9C8E8A8B9A9c9C8E8A8B9A9c9C8E8A8B9', 'sample_private_key_4', false, NOW(), NOW()),
(uuid_generate_v4(), '0xB9A9c9C8E8A8B9A9c9C8E8A8B9A9c9C8E8A8B9A9c', 'sample_private_key_5', false, NOW(), NOW());

-- Insert sample USDT rate
INSERT INTO usdt_rates (id, idr_to_usdt_rate, created_at) VALUES 
(uuid_generate_v4(), 15420.50, NOW());

-- Insert sample events
INSERT INTO events (id, name, description, location, schedule, price_idr, quota, available_quota, created_at, updated_at) VALUES 
(uuid_generate_v4(), 'Tech Conference 2025', 'Annual technology conference with leading speakers', 'Jakarta Convention Center', '2025-08-15 09:00:00', 500000, 1000, 1000, NOW(), NOW()),
(uuid_generate_v4(), 'Music Festival', 'Three-day music festival featuring local and international artists', 'Gelora Bung Karno Stadium', '2025-09-20 16:00:00', 750000, 5000, 5000, NOW(), NOW()),
(uuid_generate_v4(), 'Startup Pitch Day', 'Startup pitch competition for emerging entrepreneurs', 'Bandung Digital Valley', '2025-08-30 10:00:00', 250000, 200, 200, NOW(), NOW());