-- migrations/001_create_tables.sql

-- Таблица ресторанов
CREATE TABLE IF NOT EXISTS restaurants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица блюд
CREATE TABLE IF NOT EXISTS meals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    calories INT NOT NULL CHECK (calories > 0),
    description TEXT,
    price DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для поиска
CREATE INDEX IF NOT EXISTS idx_meals_restaurant_id ON meals(restaurant_id);
CREATE INDEX IF NOT EXISTS idx_meals_calories ON meals(calories);

-- Таблица пользователей (тестовые)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id VARCHAR(255) NOT NULL UNIQUE,
    api_key VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_api_key ON users(api_key);
CREATE INDEX IF NOT EXISTS idx_users_device_id ON users(device_id);

-- Таблица сохраненных наборов блюд
CREATE TABLE IF NOT EXISTS meal_collections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    total_calories INT NOT NULL CHECK (total_calories > 0),
    meal_ids UUID[] NOT NULL, -- массив ID блюд
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_meal_collections_user_id ON meal_collections(user_id);
CREATE INDEX IF NOT EXISTS idx_meal_collections_created_at ON meal_collections(created_at DESC);
