-- Insert sample users
INSERT INTO users (id, first_name, last_name, phone_number, address, pin, balance, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'John', 'Doe', '081234567890', 'Jl. Sudirman No. 1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 1000000, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440001', 'Jane', 'Smith', '081234567891', 'Jl. Thamrin No. 2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 500000, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440002', 'Bob', 'Johnson', '081234567892', 'Jl. Gatot Subroto No. 3', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 750000, NOW(), NOW());

-- Insert sample top ups
INSERT INTO top_ups (id, user_id, amount, balance_before, balance_after, created_at) VALUES
('660e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 500000, 500000, 1000000, NOW()),
('660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 500000, 0, 500000, NOW()),
('660e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 750000, 0, 750000, NOW());

-- Insert sample transactions
INSERT INTO transactions (id, user_id, transaction_type, amount, remarks, balance_before, balance_after, status, created_at) VALUES
('770e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 'CREDIT', 500000, 'Top Up', 500000, 1000000, 'SUCCESS', NOW()),
('770e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'CREDIT', 500000, 'Top Up', 0, 500000, 'SUCCESS', NOW()),
('770e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 'CREDIT', 750000, 'Top Up', 0, 750000, 'SUCCESS', NOW());

-- Insert sample payments
INSERT INTO payments (id, user_id, amount, remarks, balance_before, balance_after, created_at) VALUES
('880e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 50000, 'Pulsa Telkomsel', 1000000, 950000, NOW());

-- Insert sample payment transaction
INSERT INTO transactions (id, user_id, transaction_type, amount, remarks, balance_before, balance_after, status, created_at) VALUES
('770e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440000', 'DEBIT', 50000, 'Pulsa Telkomsel', 1000000, 950000, 'SUCCESS', NOW());
