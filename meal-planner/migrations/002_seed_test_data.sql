-- migrations/002_seed_test_data.sql

-- Добавляем ресторанов
INSERT INTO restaurants (name) VALUES
    ('Пицца Хаус'),
    ('Суши Мастер'),
    ('Бургерная'),
    ('Салат Бар')
ON CONFLICT (name) DO NOTHING;

-- Достаем ID ресторанов (для использования в INSERT'ах)
-- Примечание: в реальных миграциях лучше использовать sub-queries или явные ID

-- Пицца Хаус (примерно 250-550 ккал на блюдо)
INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Маргарита', 400, 'Классическая пицца с моцареллой', 450.00
FROM restaurants WHERE name = 'Пицца Хаус'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Пепперони', 520, 'Пицца с острой колбасой', 520.00
FROM restaurants WHERE name = 'Пицца Хаус'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Спагетти Болоньезе', 550, 'Классические спагетти', 380.00
FROM restaurants WHERE name = 'Пицца Хаус'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Паста Карбонара', 480, 'Паста сливочная с беконом', 420.00
FROM restaurants WHERE name = 'Пицца Хаус'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Десерт Тирамису', 300, 'Классический итальянский десерт', 250.00
FROM restaurants WHERE name = 'Пицца Хаус'
ON CONFLICT DO NOTHING;

-- Суши Мастер (примерно 200-400 ккал)
INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Филадельфия', 350, 'Рол с лососем и сливочным сыром', 380.00
FROM restaurants WHERE name = 'Суши Мастер'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Калифорния', 380, 'Рол с крабом и авокадо', 420.00
FROM restaurants WHERE name = 'Суши Мастер'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Мисо суп', 150, 'Традиционный японский суп', 180.00
FROM restaurants WHERE name = 'Суши Мастер'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Темпура', 320, 'Жареные овощи и морепродукты', 350.00
FROM restaurants WHERE name = 'Суши Мастер'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Рис с овощами', 200, 'Рис басмати с сезонными овощами', 220.00
FROM restaurants WHERE name = 'Суши Мастер'
ON CONFLICT DO NOTHING;

-- Бургерная (300-650 ккал)
INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Классик бургер', 650, 'Классический бургер с говядиной', 320.00
FROM restaurants WHERE name = 'Бургерная'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Чикен бургер', 520, 'Бургер с куриной грудкой', 280.00
FROM restaurants WHERE name = 'Бургерная'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Картофель фри', 320, 'Хрустящий картофель', 120.00
FROM restaurants WHERE name = 'Бургерная'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Молочный шейк', 280, 'Ванильный молочный коктейль', 150.00
FROM restaurants WHERE name = 'Бургерная'
ON CONFLICT DO NOTHING;

-- Салат Бар (100-350 ккал)
INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Цезарь с курицей', 300, 'Классический салат с курицей', 280.00
FROM restaurants WHERE name = 'Салат Бар'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Греческий салат', 220, 'Салат с фетой и помидорами', 250.00
FROM restaurants WHERE name = 'Салат Бар'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Салат Байсь', 150, 'Легкий салат с зеленью', 200.00
FROM restaurants WHERE name = 'Салат Бар'
ON CONFLICT DO NOTHING;

INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Смузи ягодный', 180, 'Смузи из свежих ягод', 160.00
FROM restaurants WHERE name = 'Салат Бар'
ON CONFLICT DO NOTHING;

-- Добавляем тестовых пользователей
INSERT INTO users (device_id, api_key) VALUES
    ('device-laptop-001', 'test-user-abc123xyz'),
    ('device-mobile-002', 'test-user-def456uvw'),
    ('device-tablet-003', 'test-user-ghi789rst')
ON CONFLICT (api_key) DO NOTHING;
