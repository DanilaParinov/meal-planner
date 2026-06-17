# 🍽️ Meal Planner

Сервис для подбора наборов блюд на основе калорийного лимита и доступных ресторанов.

## Стек технологий

- **Backend**: Go + Gin
- **Database**: PostgreSQL 15
- **Frontend**: HTML + CSS + Vanilla JS
- **Deployment**: Yandex Cloud (позже)

## Требования (для локальной разработки)

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15 (или через Docker)

## Быстрый старт

### 1. Клонируй репозиторий и перейди в папку

```bash
cd meal-planner
```

### 2. Подготовь переменные окружения

```bash
cp .env.example .env
```

Если нужно изменить параметры, отредактируй `.env`.

### 3. Запусти PostgreSQL в Docker

```bash
docker-compose up -d
```

Это запустит:
- PostgreSQL на `localhost:5432`
- pgAdmin на `http://localhost:5050` (опционально, для просмотра БД)

Убедись, что контейнер здоров:
```bash
docker-compose ps
```

### 4. Установи зависимости Go

```bash
go mod download
go mod tidy
```

### 5. Запусти сервер

```bash
cd cmd/server
go run main.go
```

Вывод должен быть похож на:
```
2024/01/15 10:30:45 ✓ Connected to database
2024/01/15 10:30:45 ✓ Executed migration: 001_create_tables.sql
2024/01/15 10:30:45 ✓ Executed migration: 002_seed_test_data.sql
2024/01/15 10:30:45 🚀 Starting Meal Planner API (env: development, port: 8080)
2024/01/15 10:30:45 ✓ Server listening on http://localhost:8080
2024/01/15 10:30:45 ✓ UI available at http://localhost:8080/ui
```

### 6. Открой приложение в браузере

```
http://localhost:8080/ui
```

или просто

```
http://localhost:8080
```

### 7. Залогинься тестовым ключом

Доступные тестовые ключи (из `.env`):
- `test-user-abc123xyz`
- `test-user-def456uvw`
- `test-user-ghi789rst`

## Структура проекта

```
meal-planner/
├── cmd/server/          # Точка входа (main.go)
├── internal/
│   ├── api/            # HTTP handlers, middleware
│   ├── db/             # Database layer
│   ├── models/         # Struct definitions
│   ├── algorithm/      # Knapsack algorithm
│   └── config/         # Configuration
├── migrations/         # SQL миграции
├── frontend/           # HTML, CSS, JS
├── docker-compose.yml  # PostgreSQL setup
├── go.mod
└── README.md
```

## API Endpoints

### Публичные

- `GET /health` — проверка статуса сервера

### Защищенные (требуют `X-API-Key` header)

#### Рестораны
```
GET /api/restaurants
```
Возвращает список всех доступных ресторанов.

#### Подбор блюд
```
POST /api/suggest
Content-Type: application/json
X-API-Key: <your-api-key>

{
  "restaurant_id": "uuid",
  "max_calories": 1500
}
```

Возвращает массив комбинаций блюд, отсортированных по близости к максимуму калорий.

#### История наборов
```
GET /api/collections
X-API-Key: <your-api-key>
```

Возвращает все сохраненные пользователем наборы (последние 50).

```
POST /api/collections
Content-Type: application/json
X-API-Key: <your-api-key>

{
  "restaurant_id": "uuid",
  "meal_ids": ["meal-id-1", "meal-id-2"],
  "total_calories": 1250
}
```

Сохраняет новый набор блюд.

## Примеры использования

### cURL

```bash
# Получить список ресторанов
curl -H "X-API-Key: test-user-abc123xyz" \
  http://localhost:8080/api/restaurants

# Найти комбинации
curl -X POST http://localhost:8080/api/suggest \
  -H "X-API-Key: test-user-abc123xyz" \
  -H "Content-Type: application/json" \
  -d '{
    "restaurant_id": "restaurant-uuid",
    "max_calories": 1500
  }'

# Сохранить набор
curl -X POST http://localhost:8080/api/collections \
  -H "X-API-Key: test-user-abc123xyz" \
  -H "Content-Type: application/json" \
  -d '{
    "restaurant_id": "restaurant-uuid",
    "meal_ids": ["meal-1", "meal-2"],
    "total_calories": 1250
  }'
```

## Тестовые данные

БД автоматически заполняется при первом запуске:

**Рестораны (4):**
- Пицца Хаус
- Суши Мастер
- Бургерная
- Салат Бар

**Блюда:** ~15-20 блюд на каждый ресторан

**Пользователи (3):** с разными API ключами

## Разработка

### Структура кода

- `internal/config/` — загрузка переменных окружения
- `internal/db/` — подключение и запросы к БД
- `internal/models/` — structs для JSON (request/response)
- `internal/algorithm/` — алгоритм подбора блюд (knapsack)
- `internal/api/` — HTTP обработчики и middleware

### Добавление нового endpoint'а

1. Добавь метод в `internal/api/handlers.go`
2. Зарегистрируй маршрут в `RegisterRoutes()`
3. Протестируй через cURL или UI

### Добавление нового блюда в БД

Создай SQL миграцию или используй прямой INSERT:

```sql
INSERT INTO meals (restaurant_id, name, calories, description, price)
SELECT id, 'Новое блюдо', 300, 'Описание', 200.00
FROM restaurants WHERE name = 'Пицца Хаус';
```

## Развертывание (Yandex Cloud)

_В разработке — см. Phase 3 дорожной карты_

## Проблемы и отладка

### Не могу подключиться к БД

Убедись, что контейнер работает:
```bash
docker-compose ps
```

Если не запущен:
```bash
docker-compose up -d
```

### Порт 5432 уже занят

Измени порт в `docker-compose.yml`:
```yaml
ports:
  - "5433:5432"  # Используем 5433 вместо 5432
```

И обнови `.env`:
```
DB_PORT=5433
```

### Миграции не выполняются

Убедись, что папка `migrations/` находится в текущей директории при запуске:
```bash
pwd  # Должен быть в корне meal-planner/
go run cmd/server/main.go
```

## Дорожная карта

- [x] Phase 1: Основы (БД, API, тестовые данные)
- [x] Phase 2: Алгоритм подбора
- [x] Phase 3: Фронтенд (HTML/CSS/JS)
- [ ] Phase 4: Пайплайн фото → БД
- [ ] Deployment на Yandex Cloud
- [ ] Добавление цены в подбор
- [ ] Пользовательские предпочтения

## Лицензия

MIT

## Вопросы?

Смотри комментарии в коде или открой issue.
