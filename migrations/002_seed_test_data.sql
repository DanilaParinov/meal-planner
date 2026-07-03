-- migrations/002_seed_test_data.sql

-- Добавляем тестовых пользователей
INSERT INTO users (device_id, api_key) VALUES
    ('device-laptop-001', 'test-user-abc123xyz'),
    ('device-mobile-002', 'test-user-def456uvw'),
    ('device-tablet-003', 'test-user-ghi789rst')
ON CONFLICT (api_key) DO NOTHING;
