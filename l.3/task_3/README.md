# CommentTree — Древовидные комментарии с навигацией и поиском

CommentTree — это веб-сервис для работы с древовидными комментариями, поддерживающий неограниченную вложенность, поиск, пагинацию и простой веб-интерфейс.

## Возможности

- **Древовидная структура комментариев** с неограниченной вложенностью
- **REST API** с полным набором CRUD операций
- **Полнотекстовый поиск** по содержимому комментариев
- **Пагинация и сортировка** для эффективной навигации
- **Веб-интерфейс** для удобного управления комментариями
- **Swagger документация** для API
- **Docker развертывание** для простой установки

## Требования

### Основные HTTP методы

- `POST /comments` — создание комментария (с указанием родительского)
- `GET /comments?parent={id}` — получение комментария и всех вложенных
- `DELETE /comments/{id}` — удаление комментария и всех вложенных
- `GET /comments/all` — получение всех комментариев
- `POST /comments/search` — полнотекстовый поиск по комментариям

### Веб-интерфейс

- Просмотр дерева комментариев с визуальной вложенностью
- Создание новых комментариев и ответов
- Удаление комментариев
- Поиск комментариев по ключевым словам

## Технологии

- **Backend**: Go, Gin Web Framework
- **Database**: PostgreSQL с полнотекстовым поиском
- **Frontend**: HTML/CSS/JavaScript
- **Containerization**: Docker & Docker Compose
- **Migrations**: Goose
- **Documentation**: Swagger

## Архитектура

Проект построен с использованием **Clean Architecture**:

- **Handler Layer** (`internal/handler/`) — HTTP обработчики, валидация запросов
- **Service Layer** (`internal/service/`) — бизнес-логика, оркестрация
- **Repository Layer** (`internal/repository/`) — работа с базой данных
- **Model Layer** (`internal/model/`) — структуры данных
- **DTO Layer** (`internal/dto/`) — объекты передачи данных

## Структура проекта

```
comment-tree/
├── cmd/                    # Точки входа приложения
│   ├── main.go            # Основная точка входа
│   └── app/
│       └── app.go         # Конфигурация приложения
├── internal/              # Внутренний код приложения
│   ├── config/           # Конфигурация
│   ├── dto/              # Data Transfer Objects
│   ├── handler/          # HTTP обработчики
│   ├── model/            # Модели данных
│   ├── repository/       # Работа с базой данных
│   └── service/          # Бизнес-логика
├── migrations/           # Миграции базы данных
├── static/               # Статические файлы фронтенда
├── config/               # Конфигурационные файлы
├── docs/                 # Swagger документация
├── docker-compose.yml    # Docker Compose конфигурация
├── Dockerfile           # Docker образ
└── go.mod               # Go модули
```


### Установка и запуск

### Переменные окружения

Создайте файл `.env` в корне проекта:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=comment-tree
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=migrations
```

1. **Клонируйте репозиторий:**
   ```bash
   git clone https://github.com/Komilov31/comment-tree
   cd comment-tree
   ```

2. **Запустите сервисы:**
   ```bash
   docker-compose up --build
   ```

3. **Откройте в браузере:**
   - Веб-интерфейс: http://localhost:8080
   - Swagger документация: http://localhost:8080/swagger/index.html


## API Документация

### Создание комментария

**POST** `/comments`

Создает новый комментарий или ответ на существующий.

**Request Body:**
```json
{
  "parent_id": 1,
  "text": "Текст комментария"
}
```

**Примеры cURL:**

### Создание корневого комментария
```bash
curl -X POST http://localhost:8080/comments \
  -H "Content-Type: application/json" \
  -d '{"text": "Это мой первый комментарий"}'
```
### Создание ответа на комментарий
```bash
curl -X POST http://localhost:8080/comments \
  -H "Content-Type: application/json" \
  -d '{"parent_id": 1, "text": "Это ответ на комментарий"}'
```

### Получение комментариев

**GET** `/comments/all`

Получает все комментарии в виде дерева.

**Пример cURL:**
```bash
curl http://localhost:8080/comments/all
```

**GET** `/comments?parent={id}&page={page}&limit={limit}`

Получает комментарии с пагинацией.

**Параметры запроса:**
- `parent` (объязательно): ID родительского комментария
- `page` (опционально): номер страницы (начиная с 1)
- `limit` (опционально): количество комментариев на странице

**Пример cURL:**
# Получение всех комментариев
```bash
curl "http://localhost:8080/comments"
```

### Получение комментариев с пагинацией
```bash
curl "http://localhost:8080/comments?page=1&limit=10"
```

### Получение ответов на конкретный комментарий
```bash
curl "http://localhost:8080/comments?parent=1"
```

### Поиск комментариев

**POST** `/comments/search`

Выполняет полнотекстовый поиск по комментариям.

**Request Body:**
```json
{
  "text": "искомый текст"
}
```

**Пример cURL:**
```bash
curl -X POST http://localhost:8080/comments/search \
  -H "Content-Type: application/json" \
  -d '{"text": "поисковый запрос"}'
```

### Удаление комментария

**DELETE** `/comments/{id}`

Удаляет комментарий и все его вложенные комментарии.

**Пример cURL:**
```bash
curl -X DELETE http://localhost:8080/comments/1
```

## Веб-интерфейс

Простой и интуитивный интерфейс включает:

- **Загрузка комментариев**: Кнопка "Load Comments" для загрузки всех комментариев
- **Визуальная вложенность**: Комментарии отображаются с отступами, показывающими иерархию
- **Создание комментариев**: Форма для создания новых комментариев с возможностью указания родителя
- **Ответы на комментарии**: Кнопка "Reply" для каждого комментария
- **Удаление**: Кнопка "Delete" для удаления комментариев
- **Поиск**: Поле поиска с кнопкой "Find" для поиска по тексту

##  Тестирование

Проект включает unit тесты для сервиса(бизнес-логики) и обработчиков:

```bash
# Запуск всех тестов
go test ./...

# Запуск тестов с покрытием
go test -cover ./...

# Запуск конкретного пакета
go test ./internal/service/...
```


