# WarehouseControl — CRUD склада с историей и ролями

## Описание проекта

WarehouseControl — это мини-система для учета товаров на складе, реализованная как backend-сервис на языке Go. Система поддерживает базовые CRUD-операции для инвентаря, логирование всех изменений с историей (кто, когда, что изменил) и ролевую модель доступа. 

Ключевые особенности:
- **CRUD-операции для товаров**: Добавление, просмотр, обновление и удаление товаров.
- **История изменений**: Все действия логируются в базе данных с использованием триггеров PostgreSQL (антипаттерн для демонстрации — в реальных проектах избегайте триггеров для бизнес-логики).
- **Роли пользователей**: 
  - `admin`: Полный доступ (CRUD + просмотр истории).
  - `manager`: Чтение и редактирование (CRUD без удаления? — зависит от реализации).
  - `viewer`: Только просмотр списка товаров.
- **Авторизация**: JWT-токены, роль передается в токене и проверяется на каждом запросе через middleware.
- **Простой интерфейс**: Статические HTML-страницы (`login.html` и `main.html`) для входа, просмотра/редактирования товаров и истории изменений.

Проект предназначен для демонстрации типичного корпоративного backend-сервиса в логистике или торговле. Реализация использует триггеры в БД для истории, чтобы показать их недостатки (сложность отладки, производительность), и не рекомендуется для продакшена.

## Технологии

- **Язык программирования**: Go (Golang) — для backend-логики, обработчиков и сервисов.
- **Веб-фреймворк**: Gin (предполагается по структуре handler'ов) — для маршрутизации и HTTP-запросов.
- **База данных**: PostgreSQL — для хранения товаров, пользователей, истории.
- **Авторизация**: JWT (JSON Web Tokens) для аутентификации и передачи ролей.
- **Конфигурация**: YAML-файлы для настроек (config/config.yaml).
- **Валидация**: Кастомный валидатор (internal/validator).
- **Контейнеризация**: Docker и Docker Compose — для запуска сервиса и БД.
- **Документация**: Swagger (docs/swagger.yaml/json) для API-документации.
- **Фронтенд**: Простые статические HTML-файлы с JavaScript (static/login.html, static/main.html) для интерфейса.

Зависимости управляются через `go.mod` и `go.sum`. Нет внешних ORM — прямые SQL-запросы.

## Структура проекта

Проект организован по принципам чистой архитектуры (Clean Architecture) с разделением на слои: presentation (handlers), business logic (services), data access (repositories), models и utils.

```
.
├── cmd/                  # Точка входа в приложение
│   ├── main.go           # Основной файл запуска
│   └── app/              # Конфигурация приложения
│       └── app.go        # Инициализация сервера, роутов и зависимостей
├── config/               # Конфигурационные файлы
│   └── config.yaml       # Настройки (порт, БД, JWT-secret и т.д.)
├── docs/                 # Документация API
│   ├── docs.go           # Генератор Swagger
│   ├── swagger.json      # JSON-описание API
│   └── swagger.yaml      # YAML-описание API
├── internal/             # Внутренняя логика (не экспортируется)
│   ├── auth/             # Авторизация
│   │   ├── jwt.go        # Генерация/проверка JWT-токенов
│   │   └── password.go   # Хэширование паролей (bcrypt?)
│   ├── config/           # Загрузка конфигурации
│   │   ├── config.go     # Парсер YAML
│   │   └── types.go      # Типы конфигурации
│   ├── dto/              # Data Transfer Objects
│   │   └── dto.go        # Структуры для запросов/ответов (Item, User и т.д.)
│   ├── handler/          # HTTP-обработчики (presentation layer)
│   │   ├── handler.go    # Базовый handler
│   │   ├── create.go     # POST /items
│   │   ├── get.go        # GET /items, /items/{id}
│   │   ├── update.go     # PUT /items/{id}
│   │   ├── delete.go     # DELETE /items/{id}
│   │   └── handler_test.go # Тесты handlers
│   ├── middleware/       # Middleware
│   │   └── middleware.go # Проверка JWT и ролей
│   ├── model/            # Модели БД
│   │   └── model.go      # Структуры Item, User, History
│   ├── repository/       # Доступ к БД (data layer)
│   │   ├── repository.go # Базовый репозиторий (DB connection)
│   │   ├── create.go     # INSERT в items/users
│   │   ├── get.go        # SELECT из items/history
│   │   ├── update.go     # UPDATE items
│   │   └── delete.go     # DELETE items
│   ├── service/          # Бизнес-логика
│   │   ├── service.go    # Базовый сервис
│   │   ├── create.go     # Логика создания товара
│   │   ├── get.go        # Получение товаров/истории
│   │   ├── update.go     # Обновление товара
│   │   └── delete.go     # Удаление товара
│   └── validator/        # Валидация
│       └── validator.go  # Проверка входных данных
├── migrations/           # Миграции БД
│   └── 20250928205834_create_tables.sql # Создание таблиц (items, users, history с триггерами)
├── static/               # Статические файлы для UI
│   ├── login.html        # Страница входа (выбор роли)
│   └── main.html         # Основная страница (CRUD + история)
├── Dockerfile            # Docker-образ для Go-приложения
├── docker-compose.yml    # Композиция: app + postgres
├── go.mod                # Go-модули
└── go.sum                # Зависимости
```

## Установка и запуск

### Предварительные требования
- Go
- Docker и Docker Compose
- PostgreSQL

### Через Docker (рекомендуется)
1. Клонируйте репозиторий.
2. Запустите: 
```bash
    git clone https://github.com/Komilov31/warehouse-control.git
```
```bash
    cd warehouse-control
```
```bash
    cp .env.example .env #редактируйте .env файл
```
```bash
    docker-compose up -d # поднимет PostgreSQL и Go-приложение.
```
 
3. Приложение доступно на `http://localhost:8080`.

4. UI: Откройте `http://localhost:8080/static/login.html`.

5. Миграции применятся автоматически (или вручную через `docker-compose exec app goose up`).

Swagger: `http://localhost:8080/swagger/index.html`.

## API и примеры curl-запросов

Сервер работает на порту 8080. Все защищенные эндпоинты требуют JWT-токен в заголовке `Authorization: Bearer <token>`. Роли проверяются в middleware.

### 1. Создание пользователя (доступно без токена)
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "admin",
    "role": "admin"
  }'
```
Ответ: `{"id":1,"name":"admin","role":"admin","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...","create_at":"2023-..."}`. Сохраните `token` для следующих запросов.

### 2. CRUD-операции для товаров (требует токен, роль проверяется)
#### Создание товара (POST /items) — admin/manager
```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "Ноутбук",
    "count": 10
  }'
```
Ответ: `{"id":1,"name":"Ноутбук","count":10,"created_at":"2023-..."}`. Триггер запишет в историю.

#### Получение списка товаров (GET /items) — все роли
```bash
curl -X GET http://localhost:8080/items \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```
Ответ: `[{"id":1,"name":"Ноутбук","count":10,"created_at":"2023-..."}]`. Viewer увидит только чтение.

#### Обновление товара (PUT /items/{id}) — admin/manager
```bash
curl -X PUT http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "user_id": 1,
    "name": "Ноутбук Pro",
    "count": 5
  }'
```
Ответ: `{"message":"updated"}`. Триггер запишет изменения в историю.

#### Удаление товара (DELETE /items/{id}) — admin
```bash
curl -X DELETE http://localhost:8080/items/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```
Ответ: `{"message":"Item deleted successfully"}`. Триггер запишет удаление.

### 4. Просмотр истории изменений (GET /items/{id}/history) — admin/viewer (зависит от роли)
```bash
curl -X GET http://localhost:8080/items/1/history \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```
Ответ: `[{"id":1,"item_id":1,"user_id":1,"action":"UPDATE","timestamp":"2023-...","changes":"quantity:10->5"}, ...]`.


## UI-интерфейс
- `http://localhost:8080/login.html`: Выбор пользователя/роли, вход, получение токена.
- `http://localhost:8080/main.html`: Таблица товаров (CRUD-формы, если права), колонка с историей по клику.

## Ограничения и замечания
- Триггеры в PostgreSQL для истории — антипаттерн: усложняют тестирование и миграции. В реальности используйте application-level логирование.
- Тестирование: Написаны unit-тесты в `internal/handler/handler_test.go`.

