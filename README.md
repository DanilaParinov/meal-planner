# 🍽️ Meal Planner

Сервис для подбора наборов блюд на основе калорийного лимита.

## Стек

- **Backend**: Go + Gin
- **Database**: PostgreSQL 15 (Docker)
- **Frontend**: HTML + CSS + Vanilla JS

---

## Первый запуск

### Проверь зависимости

```powershell
go version        # нужен Go 1.21+
docker --version  # нужен Docker Desktop (запущенный)
```

### Создай .env

```powershell
cp .env.example .env
```

### Собери бинарник

```powershell
go build -o main.exe ./cmd/server/
```

### Запусти

```powershell
docker-compose up -d
.\main.exe
```

---

## Регулярный запуск

```powershell
# 1. Запустить БД (если не запущена)
docker-compose up -d

# 2. Запустить сервер
.\main.exe
```

**Приложение:** `http://localhost:8080`  
**Админ-панель:** `http://localhost:8080/admin.html`  
**Тестовые API ключи (для локальной разработки):** `test-user-abc123xyz`, `test-user-def456uvw`

Для реального пользователя или администратора создайте отдельный аккаунт:
```bash
go run cmd/createuser/main.go -name "Имя" [-admin]
```
Ключ будет выведен один раз в терминал — сохраните его.

### Остановить

```powershell
# Сервер: Ctrl+C в терминале

# БД
docker-compose down
```

### Если порт 8080 уже занят

Это значит, что предыдущий сервер не был остановлен. Найди и убей процесс:

```powershell
Get-Process -Id (Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue).OwningProcess | Stop-Process -Force
```

---

## Структура проекта

```
meal-planner/
├── cmd/server/main.go       # точка входа
├── internal/
│   ├── api/                 # HTTP handlers, middleware
│   ├── db/                  # подключение и SQL запросы
│   ├── models/              # structs
│   ├── algorithm/           # алгоритм подбора блюд
│   └── config/              # конфигурация
├── migrations/              # SQL миграции (применяются автоматически)
├── frontend/                # HTML, CSS, JS
├── menu-photos/             # фото меню (не в git)
├── docker-compose.yml
└── .env.example              # шаблон, скопируй в .env
```

## Данные в БД

Рестораны и блюда загружаются автоматически при первом запуске из `migrations/003_import_restaurants.sql`:

| Ресторан | Блюд |
|---|---|
| Бостон | 58 |
| Магнум | 64 |
| Mozza | 78 |
| Торро Гриль | 79 |

## API

| Метод | Путь | Описание |
|---|---|---|
| GET | `/health` | статус сервера |
| GET | `/api/restaurants` | список ресторанов |
| POST | `/api/suggest` | подбор блюд по калориям |
| GET | `/api/collections` | история пользователя |
| POST | `/api/collections` | сохранить набор |
| POST | `/api/admin/restaurants` | создать ресторан |
| POST | `/api/admin/meals` | добавить блюдо |

Все `/api/*` эндпоинты (кроме `/health`) требуют заголовок `X-API-Key`.
