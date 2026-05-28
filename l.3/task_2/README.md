# URL Shortener API

## Описание проекта

Это сервис сокращения URL-адресов с аналитикой, разработанный на языке Go. Сервис позволяет создавать короткие ссылки, перенаправлять на оригинальные URL и собирать статистику по переходам, включая данные по датам, месяцам и user agent.

Проект использует следующие технологии:
- **Go** - основной язык программирования
- **wbf** - веб-фреймворк(обертки над многими библиотеками)
- **PostgreSQL** - база данных для хранения URL и статистики
- **Redis** - кэш для быстрого доступа
- **Swagger** - документация API

## Функциональность

- Создание коротких URL из длинных ссылок
- Перенаправление по коротким URL на оригинальные
- Сбор аналитики по переходам (время, user agent)
- Агрегация статистики по датам, месяцам и user agent
- Веб-интерфейс для взаимодействия
- Документация API через Swagger

## Запуск проекта

### Предварительные требования

- Docker и Docker Compose
- Файл `.env` с переменными окружения (пример ниже)

### Переменные окружения (.env)

Создайте файл `.env` в корне проекта:

```bash
  cp .env.example .env
```

### Команды для запуска

1. Клонируйте репозиторий и перейдите в директорию проекта.
```bash
  git clone https://github.com/Komilov31/url-shortner
```
2. Настройте файл `.env`.
3. Запустите сервисы:

```bash
docker-compose up --build
```

Сервис будет доступен по адресу: `http://localhost:8080`

Swagger документация: `http://localhost:8080/swagger/index.html`

## API Эндпоинты

### 1. Получить главную страницу
**GET /**

Возвращает HTML страницу интерфейса.

```bash
curl -X GET "http://localhost:8080/"
```

### 2. Создать короткий URL
**POST /shorten**

Создает короткий URL из предоставленного оригинального.

Тело запроса:
```json
{
  "url": "https://example.com"
}
```

```bash
curl -X POST "http://localhost:8080/shorten" \
     -H "Content-Type: application/json" \
     -d '{
       "url": "https://example.com"
     }'
```

Ответ:
```json
{
  "short_url": "abc123",
  "url": "https://example.com"
}
```

### 3. Перенаправление по короткому URL
**GET /s/{short_url}**

Перенаправляет на оригинальный URL.

```bash
curl -L -X GET "http://localhost:8080/s/abc123"
```

### 4. Получить аналитику для короткого URL
**GET /analytics/{short_url}**

Возвращает детальную аналитику по переходам для указанного короткого URL.

```bash
curl -X GET "http://localhost:8080/analytics/abc123"
```

Ответ:
```json
{
  "short_url": "abc123",
  "url": "https://example.com",
  "redirect_count": 10,
  "request_time": ["2023-10-01T12:00:00Z", "2023-10-01T13:00:00Z"],
  "user_agent": ["Mozilla/5.0 ...", "curl/7.68.0"]
}
```

### 5. Агрегированная аналитика по датам
**GET /analytics/date**

Возвращает статистику переходов, сгруппированную по датам.

```bash
curl -X GET "http://localhost:8080/analytics/date"
```

### 6. Агрегированная аналитика по месяцам
**GET /analytics/month**

Возвращает статистику переходов, сгруппированную по месяцам.

```bash
curl -X GET "http://localhost:8080/analytics/month"
```

### 7. Агрегированная аналитика по пользовательским агентам
**GET /analytics/user_agent**

Возвращает статистику переходов, сгруппированную по пользовательским агентам.

```bash
curl -X GET "http://localhost:8080/analytics/user_agent"
```

## Структура проекта

```
.
├── cmd/
│   ├── app/
│   │   └── app.go          # Настройка приложения и маршрутов
│   └── main.go             # Точка входа
├── config/
│ 
│   └── config.yaml         # YAML конфиг
├── docs/
│   ├── docs.go             # Swagger генерация
│   ├── swagger.json        # JSON спецификация
│   └── swagger.yaml        # YAML спецификация
├── internal/
│   ├── cache/redis/        # Redis кэш
│   ├── config/             # Получение конфигов из yaml и .env
│   ├── dto/                # Data Transfer Objects
│   ├── handler/            # HTTP обработчики
│   ├── model/              # Модели данных
│   ├── repository/         # Репозиторий (БД)
│   └── service/            # Бизнес-логика
├── migrations/             # Миграции БД
├── static/                 # Статические файлы (HTML, CSS, JS)
├── docker-compose.yml      # Docker Compose
├── Dockerfile              # Docker образ
├── go.mod                  # Go модули
└── go.sum                  # Checksums
```

