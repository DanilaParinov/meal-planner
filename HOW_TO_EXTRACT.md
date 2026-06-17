# 📦 Как распаковать архив meal-planner.tar.gz

## Что в архиве?

Архив `meal-planner.tar.gz` содержит **ВСЕ файлы** проекта:
- ✅ 8 файлов Go кода
- ✅ 2 SQL миграции
- ✅ 3 файла фронтенда (HTML/CSS/JS)
- ✅ Конфигурационные файлы
- ✅ 8 файлов документации (.md)
- ✅ docker-compose.yml и go.mod

**Размер архива:** 41 KB

---

## 🖥️ Windows (PowerShell или 7-Zip)

### Способ 1: PowerShell (встроенный)

```powershell
# Перейди в папку где находится архив
cd Downloads  # или где у тебя архив

# Распакуй архив
tar -xzf meal-planner.tar.gz

# Готово! Папка meal-planner создана
ls meal-planner
```

### Способ 2: 7-Zip (если установлен)

1. Правый клик на `meal-planner.tar.gz`
2. 7-Zip → Extract Here
3. Готово!

### Способ 3: Windows Explorer + WSL2

Если у тебя WSL2 установлен:

```bash
# В WSL терминале
cd ~/Downloads
tar -xzf meal-planner.tar.gz
```

---

## 🍎 macOS (Terminal)

```bash
# Перейди в папку где архив
cd ~/Downloads  # или где у тебя архив

# Распакуй
tar -xzf meal-planner.tar.gz

# Проверь
ls -la meal-planner
```

---

## 🐧 Linux (Terminal)

```bash
# Перейди в папку с архивом
cd ~/Downloads

# Распакуй
tar -xzf meal-planner.tar.gz

# Проверь содержимое
ls meal-planner
```

---

## ✅ Проверка что всё распаковалось

После распаковки должна быть такая структура:

```
meal-planner/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers.go
│   │   └── middleware.go
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   ├── db.go
│   │   └── queries.go
│   ├── models/
│   │   └── models.go
│   └── algorithm/
│       └── knapsack.go
├── migrations/
│   ├── 001_create_tables.sql
│   └── 002_seed_test_data.sql
├── frontend/
│   ├── index.html
│   ├── style.css
│   └── app.js
├── docker-compose.yml
├── go.mod
├── .env.example
├── .claudeignore
├── START_HERE.md
├── CLAUDE.md
├── CLAUDE_CODE_SETUP.md
├── CLAUDE_CODE_TASKS.md
├── README.md
└── [другие .md файлы]
```

Проверь что все эти файлы есть!

---

## 🚀 Что делать после распаковки

1. **Открой папку в VSCode:**
   ```bash
   cd meal-planner
   code .
   ```

2. **Установи Claude Code расширение** (если еще не установил)
   - Ctrl+Shift+X → поиск "Claude Code" → Install

3. **Запусти PostgreSQL:**
   ```bash
   docker-compose up -d
   ```

4. **Запусти сервер:**
   ```bash
   go mod download
   cd cmd/server
   go run main.go
   ```

5. **Открой Claude Code** и начни разработку!

---

## 🆘 Если распаковка не работает

### Ошибка: "команда tar не найдена" (Windows)

→ Используй 7-Zip или обновись на Windows 10+ (там tar встроен)

### Ошибка: "Permission denied" (Mac/Linux)

→ Проверь права: `ls -l meal-planner.tar.gz`

### Архив поврежден

→ Скачай снова из outputs

---

## 📊 Размеры файлов

```
meal-planner.tar.gz        41 KB  (сжатый архив)
meal-planner/ (распакованный)  ~200 KB
```

Очень компактно! 🎉

---

## ✨ Готово!

После распаковки у тебя есть **полный рабочий проект** со всеми файлами.

**Следующий шаг:** Читай START_HERE.md в папке проекта

**Удачи! 🚀**
