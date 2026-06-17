# 🎉 Phase 1: Завершено!

## Что было создано

Я подготовил **полный рабочий MVP** для локальной разработки Meal Planner. Вот что входит:

### 📁 Файлы структуры проекта (17 файлов)

#### Go Backend:
1. **go.mod** — зависимости (Gin, PostgreSQL driver, dotenv)
2. **cmd/server/main.go** — точка входа, инициализация сервера
3. **internal/config/config.go** — загрузка конфигурации из .env
4. **internal/db/db.go** — подключение к БД, выполнение миграций
5. **internal/db/queries.go** — repository pattern (SQL запросы)
6. **internal/models/models.go** — struct'ы для API и БД
7. **internal/api/middleware.go** — аутентификация по API ключу
8. **internal/api/handlers.go** — HTTP endpoints (restaurants, suggest, collections)
9. **internal/algorithm/knapsack.go** — алгоритм подбора наборов (DP + backtracking)

#### БД:
10. **migrations/001_create_tables.sql** — создание таблиц (restaurants, meals, users, meal_collections)
11. **migrations/002_seed_test_data.sql** — тестовые данные (4 ресторана, ~20 блюд, 3 пользователя)
12. **.env.example** — пример переменных окружения

#### Frontend:
13. **frontend/index.html** — главная страница с двумя режимами (auth + app)
14. **frontend/style.css** — современный дизайн (responsive, темная тема)
15. **frontend/app.js** — клиентская логика (API запросы, UI управление)

#### DevOps:
16. **docker-compose.yml** — запуск PostgreSQL + pgAdmin локально
17. **README.md** — полная документация

---

## Архитектура в 30 секунд

```
Frontend (HTML/CSS/JS)
    ↓ (HTTP Requests)
    ↓
Go Backend (Gin Framework)
    ├─ HTTP Handlers (API endpoints)
    ├─ Middleware (Authentication)
    └─ Business Logic
        ├─ Knapsack Algorithm (подбор блюд)
        └─ Repository Pattern (DB queries)
    ↓ (SQL)
    ↓
PostgreSQL Database
    ├─ restaurants
    ├─ meals
    ├─ users
    └─ meal_collections
```

---

## Как запустить (4 шага)

### Шаг 1️⃣: Запусти PostgreSQL

```bash
cd meal-planner
docker-compose up -d
```

Это запустит БД на `localhost:5432`.

### Шаг 2️⃣: Скачай зависимости Go

```bash
go mod download
go mod tidy
```

### Шаг 3️⃣: Запусти сервер

```bash
cd cmd/server
go run main.go
```

Ожидаемый вывод:
```
✓ Connected to database
✓ Executed migration: 001_create_tables.sql
✓ Executed migration: 002_seed_test_data.sql
🚀 Starting Meal Planner API (env: development, port: 8080)
✓ Server listening on http://localhost:8080
✓ UI available at http://localhost:8080/ui
```

### Шаг 4️⃣: Открой UI

```
http://localhost:8080
```

Залогинься любым из тестовых ключей:
- `test-user-abc123xyz`
- `test-user-def456uvw`
- `test-user-ghi789rst`

---

## Что ты можешь делать в MVP

✅ **Просмотреть рестораны** — список всех ресторанов с блюдами

✅ **Подобрать наборы блюд** — указать лимит калорий → система найдет все подходящие комбинации из одного ресторана

✅ **Сохранить набор** — сохранить понравившийся вариант в истории

✅ **Просмотреть историю** — все сохраненные наборы с датой и калорийностью

---

## API Endpoints (готовы к тестированию)

| Метод | Endpoint | Описание |
|-------|----------|---------|
| GET | `/health` | Проверка статуса |
| GET | `/api/restaurants` | Список ресторанов |
| POST | `/api/suggest` | Подбор блюд по калориям |
| GET | `/api/collections` | История пользователя |
| POST | `/api/collections` | Сохранить набор |

**Все endpoints (кроме /health) требуют header:** `X-API-Key: <your-key>`

---

## Примеры тестирования (cURL)

### 1. Получить все рестораны
```bash
curl -H "X-API-Key: test-user-abc123xyz" \
  http://localhost:8080/api/restaurants
```

### 2. Найти комбинации на 1500 ккал
```bash
curl -X POST http://localhost:8080/api/suggest \
  -H "X-API-Key: test-user-abc123xyz" \
  -H "Content-Type: application/json" \
  -d '{
    "restaurant_id": "<restaurant-uuid-здесь>",
    "max_calories": 1500
  }'
```

(Замени `<restaurant-uuid-здесь>` на реальный ID из первого запроса)

### 3. Сохранить набор
```bash
curl -X POST http://localhost:8080/api/collections \
  -H "X-API-Key: test-user-abc123xyz" \
  -H "Content-Type: application/json" \
  -d '{
    "restaurant_id": "<restaurant-uuid>",
    "meal_ids": ["<meal-id-1>", "<meal-id-2>"],
    "total_calories": 1250
  }'
```

---

## Алгоритм подбора (Knapsack)

**Как это работает:**

Пользователь указывает максимум калорий (например, 1500) → система находит **все возможные подмножества блюд** из выбранного ресторана, сумма калорий которых ≤ 1500.

**Сложность:** O(2^n) в худшем случае, но оптимизировано:
- Рекурсия с backtracking
- Ранний выход если превышаем лимит
- Максимум 20 результатов возвращаются (top N по близости к максимуму)

**Пример:**
```
Лимит: 1500 ккал
Ресторан "Пицца Хаус" меню:
  - Маргарита: 400 ккал
  - Спагетти: 550 ккал
  - Десерт: 300 ккал

Найденные комбинации (отсортированы по близости к 1500):
1. Маргарита + Спагетти + Десерт = 1250 ккал ← ближе к максимуму
2. Маргарита + Спагетти = 950 ккал
3. Спагетти + Десерт = 850 ккал
4. Маргарита + Десерт = 700 ккал
... и т.д.
```

---

## Тестовые данные в БД

**Рестораны (4):**
- Пицца Хаус (Паста, Пицца)
- Суши Мастер (Роллы, Супы)
- Бургерная (Бургеры, Напитки)
- Салат Бар (Салаты, Смузи)

**Блюда:** ~5-6 блюд на ресторан, с разными калорийностью (150-650 ккал)

**Пользователи:** 3 тестовых с разными API ключами

---

## Что дальше? (Phase 2+)

### Phase 2: Улучшение алгоритма ✅ (готов, но не обязателен)
- [ ] Кеширование результатов (Redis)
- [ ] Параллельная обработка (goroutines)
- [ ] Поддержка смешанных ресторанов (опция)

### Phase 3: Фото → БД (не спешим)
- [ ] CLI утилита для анализа фото меню
- [ ] Интеграция Google Vision API
- [ ] Админ-панель для валидации

### Phase 4: Deployment
- [ ] Docker контейнеризация
- [ ] Развертывание на Yandex Cloud
- [ ] CI/CD pipeline (GitHub Actions)

### Phase 5: Фичи
- [ ] Фильтр по цене
- [ ] Пользовательские предпочтения (аллергены, вегетарианское и т.д.)
- [ ] Рекомендации (ML?)
- [ ] Интеграция с доставкой

---

## Структура файлов на диске

```
meal-planner/
├── cmd/
│   └── server/
│       └── main.go                    ← точка входа
├── internal/
│   ├── api/
│   │   ├── handlers.go               ← HTTP endpoints
│   │   └── middleware.go             ← auth
│   ├── db/
│   │   ├── db.go                     ← подключение
│   │   └── queries.go                ← SQL запросы
│   ├── models/
│   │   └── models.go                 ← structs
│   ├── algorithm/
│   │   └── knapsack.go               ← алгоритм
│   └── config/
│       └── config.go                 ← env vars
├── migrations/
│   ├── 001_create_tables.sql         ← schema
│   └── 002_seed_test_data.sql        ← data
├── frontend/
│   ├── index.html
│   ├── style.css
│   └── app.js
├── docker-compose.yml                ← PostgreSQL
├── go.mod
├── .env.example
└── README.md
```

---

## Быстрая отладка

| Проблема | Решение |
|----------|---------|
| **Не могу подключиться к БД** | `docker-compose ps` — убедись что контейнер работает |
| **Порт 5432 занят** | Измени в docker-compose.yml: `5433:5432` |
| **Миграции не выполняются** | Убедись что работаешь из корня `meal-planner/` |
| **Фронтенд не загружается** | Проверь что идешь на `http://localhost:8080/ui` |
| **API ошибка 401** | Проверь что передаешь правильный `X-API-Key` |

---

## Кратко для быстрого старта

```bash
# 1. В отдельном терминале: запусти БД
docker-compose up -d

# 2. В основном терминале: запусти сервер
cd cmd/server && go run main.go

# 3. Открой браузер
http://localhost:8080

# 4. Залогинься и тестируй!
```

---

**Всё готово к разработке! 🚀**

Следующий шаг: протестируй MVP локально, убедись что всё работает, потом можно переходить на Phase 2 (опционально улучшить алгоритм) или сразу на Phase 4 (пайплайн фото).

Есть вопросы? Давай разбираться! 👍
