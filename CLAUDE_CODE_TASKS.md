# 📋 Примеры задач для Claude Code

Вот готовые задачи, которые ты можешь дать Claude Code для автоматизации разработки.

## 🟢 Easy (15-30 мин каждая)

### Task 1.1: Добавить новый API endpoint для получения информации о ресторане

```
Добавь новый endpoint:

GET /api/restaurants/:id

Что должен делать:
1. Получить информацию о ресторане по ID
2. Вернуть ресторан + все его блюда
3. Требует аутентификацию (X-API-Key)

Response JSON:
{
  "success": true,
  "data": {
    "restaurant": {
      "id": "...",
      "name": "...",
      "created_at": "..."
    },
    "meals": [
      {
        "id": "...",
        "name": "...",
        "calories": 400,
        ...
      }
    ]
  }
}

Файлы для изменения:
- internal/api/handlers.go (добавь метод GetRestaurant)
- internal/db/queries.go (если нужна новая query)
```

### Task 1.2: Улучшить валидацию API ключа

```
Улучши middleware аутентификации в internal/api/middleware.go:

Требования:
1. Проверь что X-API-Key не пустой
2. Проверь длину ключа (должен быть >= 10 символов)
3. Добавь логирование неудачных попыток (только в dev режиме)
4. Верни понятный error message

Error response должен быть:
{
  "error": "unauthorized",
  "message": "Invalid or missing API key",
  "code": "INVALID_API_KEY"
}
```

### Task 1.3: Добавить новые тестовые рестораны и блюда

```
Добавь в migrations/002_seed_test_data.sql:

1. Новый ресторан: "Паста Фабрика"
2. 5-6 блюд к нему (с разной калорийностью 200-600 ккал)
3. Реалистичные названия блюд (на русском)
4. Разные цены (200-500 рублей)

Примеры блюд:
- Лазанья (450 ккал)
- Паста Примавера (380 ккал)
- и т.д.

После этого пересоздай БД:
docker-compose down
docker-compose up -d
```

### Task 1.4: Улучшить дизайн фронтенда

```
Изменения в frontend/style.css:

1. Сделай карточки результатов более красивыми:
   - Добавь border-radius и shadows
   - Цветные индикаторы калорийности (зеленый, оранжевый, красный)

2. Добавь анимации:
   - Fade-in при загрузке результатов
   - Hover эффекты на кнопках

3. Улучши мобильный дизайн:
   - Сделай боковую панель stack'ируемой на мобилке
```

### Task 1.5: Добавить сортировку результатов

```
Добавь в frontend/app.js функцию сортировки результатов:

1. Добавь select с опциями:
   - "По близости к максимуму" (по умолчанию)
   - "По количеству блюд (меньше всего)"
   - "По количеству блюд (больше всего)"

2. Реализуй функцию которая переупорядочивает results после сортировки
3. Сохраняй выбранную сортировку в localStorage
```

## 🟡 Medium (1-2 часа каждая)

### Task 2.1: Добавить фильтр по цене

```
Добавь поддержку фильтра по цене в API:

1. Обнови endpoint /api/suggest чтобы принимал max_price (опционально):
   POST /api/suggest
   {
     "restaurant_id": "...",
     "max_calories": 1500,
     "max_price": 1000.00  // новый параметр, опционально
   }

2. В алгоритме (internal/algorithm/knapsack.go) учитывай цену
   - Если max_price задан, фильтруй результаты по цене
   - Возвращай сумму цен каждого набора в response

3. На фронтенде (frontend/app.js):
   - Добавь input для max_price
   - Показывай общую цену каждого набора

Response:
{
  "solutions": [
    {
      "meals": [...],
      "total_calories": 1250,
      "total_price": 850.00  // новое поле
    }
  ]
}
```

### Task 2.2: Добавить поиск по названию блюда

```
Добавь новый endpoint для поиска:

GET /api/meals/search?q=паста

1. Параметр q (query) — часть названия блюда
2. Опционально: restaurant_id для поиска в конкретном ресторане
3. Возвращает все блюда которые содержат q в названии

Response:
{
  "success": true,
  "data": [
    {
      "id": "...",
      "name": "Паста Болоньезе",
      "calories": 550,
      "restaurant": {...}
    }
  ]
}

Требования:
- Case-insensitive поиск
- Максимум 50 результатов
- Отсортировано по релевантности (название в начале лучше)
```

### Task 2.3: Добавить статистику пользователя

```
Создай новый endpoint:

GET /api/user/stats

Возвращает:
{
  "success": true,
  "data": {
    "total_collections": 15,
    "total_unique_meals": 28,
    "average_calories": 1380,
    "favorite_restaurant": "Пицца Хаус",
    "last_collection": "2024-01-15T10:30:00Z"
  }
}

Для этого:
1. Добавь method в internal/db/queries.go для подсчета статистики
2. Добавь handler в internal/api/handlers.go
3. Зарегистрируй маршрут в RegisterRoutes()
```

### Task 2.4: Улучшить алгоритм подбора (оптимизация)

```
Оптимизируй internal/algorithm/knapsack.go:

1. Добавь кеширование результатов (in-memory):
   - Если один пользователь ищет одно и то же — верни закешированный результат
   - Кеш действует 5 минут

2. Добавь лимит на максимальное количество блюд в наборе:
   - По умолчанию максимум 5 блюд в одном наборе
   - Параметр max_items_per_set (опционально в API)

3. Добавь параметр для фильтра по калорийности одного блюда:
   - Исключай блюда если они > max_calories сами по себе
   - (это логично — нельзя выбрать блюдо которое уже превышает лимит)
```

### Task 2.5: Добавить логирование и мониторинг

```
Добавь логирование в Go приложение:

1. Используй стандартный log пакет
2. Логируй:
   - Все успешные API запросы (с методом, путем, статусом)
   - Все ошибки (с методом, пути, сообщением)
   - Время выполнения алгоритма подбора
   - Попытки авторизации (неудачные)

3. Сохраняй логи в файл logs/app.log (если в prod режиме)

Пример лога:
"2024-01-15 10:30:45 [INFO] POST /api/suggest - user: abc123 - 125ms"
"2024-01-15 10:30:50 [ERROR] /api/collections - Invalid meal_id: xyz"
```

## 🔴 Hard (3+ часа каждая)

### Task 3.1: Поддержка смешанных ресторанов

```
Добавь в API опцию для подбора блюд из разных ресторанов:

Новый параметр:
POST /api/suggest
{
  "restaurant_id": null,  // null = любой ресторан
  "max_calories": 1500,
  "max_restaurants": 3    // максимум из 3 ресторанов
}

Что нужно сделать:
1. Обнови algorithm/knapsack.go чтобы работал с несколькими ресторанами
2. Ограничь комбинации (максимум max_restaurants)
3. Отсортируй результаты по разнообразию ресторанов
4. Верни информацию о каждом ресторане в результате

Response:
{
  "solutions": [
    {
      "meals": [...],
      "total_calories": 1250,
      "restaurants": [
        {"id": "...", "name": "Пицца Хаус", "meals_count": 2},
        {"id": "...", "name": "Суши Мастер", "meals_count": 1}
      ]
    }
  ]
}

Внимание: это добавит сложность алгоритму!
```

### Task 3.2: Интеграция с Google Vision API (подготовка)

```
Создай структуру для будущей интеграции с Vision API:

1. Создай новый package: internal/vision/
2. Структуры для API:
   - type MenuPhoto struct { ... }
   - type ExtractedMeal struct { ... }
   - type MealExtractionResult struct { ... }

3. Создай interface для abstracting Vision provider:
   type VisionProvider interface {
       ExtractMealsFromImage(ctx context.Context, imagePath string) (*MealExtractionResult, error)
   }

4. Создай mock implementation для тестирования

5. Подготовь конфигурацию для Google API ключа:
   - GOOGLE_VISION_API_KEY в .env
   - Загрузка в config/config.go

Это подготовит terrain для Phase 4 (анализ фото)
```

### Task 3.3: Админ-панель для управления меню

```
Создай простую админ-панель:

1. Новый путь: /admin (доступен только через специальный админ ключ)

2. Функционал:
   - Просмотр всех ресторанов и блюд в таблице
   - Добавить новое блюдо (форма)
   - Редактировать блюдо (изменить калории, цену)
   - Удалить блюдо
   - Просмотр статистики

3. Frontend:
   - Создай frontend/admin.html
   - Таблица с блюдами
   - Формы для CRUD операций

4. Backend:
   - POST /api/admin/meals (create)
   - PUT /api/admin/meals/:id (update)
   - DELETE /api/admin/meals/:id (delete)
   - GET /api/admin/stats (статистика)

5. Безопасность:
   - Используй отдельный ADMIN_API_KEY из .env
   - Логируй все действия админа
```

### Task 3.4: Кеширование с Redis (опционально)

```
Добавь кеширование результатов с Redis:

1. Установи redis-go библиотеку
2. Создай internal/cache/cache.go с interface:
   type Cache interface {
       Get(key string) (interface{}, error)
       Set(key string, value interface{}, ttl time.Duration) error
       Delete(key string) error
   }

3. Реализуй Redis версию

4. Обнови algorithm/knapsack.go:
   - При запросе подбора, сначала проверь кеш
   - Если есть в кеше (на 5 мин) — верни из кеша
   - Иначе вычисли и положи в кеш

5. Добавь управление кешем:
   - DELETE /api/admin/cache (очистить весь кеш)
   - GET /api/admin/cache/stats (статистика кеша)

Это значительно улучшит производительность!
```

### Task 3.5: Unit тесты для алгоритма

```
Напиши unit тесты для internal/algorithm/knapsack.go:

1. Создай internal/algorithm/knapsack_test.go

2. Тесты:
   - TestSolveBasic: простой случай (3 блюда, легко комбинировать)
   - TestSolveOptimized: оптимизированная версия
   - TestNoSolutions: случай когда нет решений
   - TestExactMatch: точное совпадение с максимумом
   - TestLargeMealSet: большой набор блюд (20+)
   - TestSingleMeal: одно блюдо в наборе
   - TestEmptyMeals: пустой список блюд

3. Используй table-driven tests:
   type testCase struct {
       name string
       meals []Meal
       maxCals int
       expectedCount int
   }

4. Убедись что все тесты pass:
   go test ./internal/algorithm/...
```

## 🎯 Как выбрать задачу

1. **Первый раз с Claude Code?** Начни с Task 1.1 или 1.2 (Easy)
2. **Уверен в Go и архитектуре?** Прыгай на Task 2.x (Medium)
3. **Готов к челленджу?** Попробуй Task 3.x (Hard)

## 💡 Советы

- ✅ После каждой задачи тестируй локально
- ✅ Коммитить только когда всё работает
- ✅ Проси Claude Code объяснить код если не понимаешь
- ❌ Не пытайся делать несколько задач одновременно
- ❌ Не скипай тестирование

## 📊 Прогресс

- [ ] Task 1.1
- [ ] Task 1.2
- [ ] Task 1.3
- [ ] Task 1.4
- [ ] Task 1.5
- [ ] Task 2.1
- [ ] Task 2.2
- [ ] Task 2.3
- [ ] Task 2.4
- [ ] Task 2.5
- [ ] Task 3.1
- [ ] Task 3.2
- [ ] Task 3.3
- [ ] Task 3.4
- [ ] Task 3.5

## 🚀 После завершения всех задач

1. Развертывание на Yandex Cloud
2. Настройка CI/CD (GitHub Actions)
3. Production-ready подготовка

Удачи! 💪
