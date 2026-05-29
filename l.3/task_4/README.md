# Image Processor API

**Image Processor API** - это высокопроизводительный сервис для обработки изображений с поддержкой асинхронной обработки, распределенного хранения и веб-интерфейса. Сервис позволяет загружать изображения, применять различные виды обработки (ресайз, добавление водяного знака, генерация миниатюры) и управлять обработанными файлами.

## 🚀 Основные возможности

- **Асинхронная обработка изображений** с использованием Apache Kafka
- **Распределенное хранение** файлов в MinIO
- **RESTful API** с полной документацией Swagger
- **Веб-интерфейс** для удобного взаимодействия
- **Поддержка различных форматов** изображений (JPEG, PNG, GIF)

## 🛠 Технологии

### Backend
- **Go** - основной язык программирования
- **Gin** - высокопроизводительный HTTP фреймворк
- **PostgreSQL** - основная база данных
- **Apache Kafka** - очередь сообщений для асинхронной обработки
- **MinIO** - объектное хранилище для файлов
- **Goose** - миграции базы данных

### Frontend
- **HTML5/CSS3/JavaScript** - веб-интерфейс

### Инфраструктура
- **Docker & Docker Compose** - контейнеризация
- **Swagger** - API документация
- **Zerolog** - структурированное логирование


## 🏗 Архитектура

### Компоненты системы

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │   Image API     │    │   Kafka Queue   │
│   (Frontend)    │───▶│   (Gin HTTP)    │───▶│   (Messages)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │     MinIO       │    │   Workers       │
│   (Metadata)    │    │   (File Store)  │    │   (Processing)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Структура проекта

```
image-processor/
├── cmd/                    # Точки входа в приложение
│   ├── main.go            # Основная точка входа
│   └── app/               # Инициализация приложения
├── internal/              # Внутренний код приложения
│   ├── config/           # Конфигурация
│   ├── dto/              # Data Transfer Objects
│   ├── handler/          # HTTP обработчики
│   ├── model/            # Модели данных
│   ├── repository/       # Репозитории для работы с данными
│   ├── service/          # Бизнес-логика
│   └── kafka/            # Kafka клиент
├── config/               # Конфигурационные файлы
├── migrations/           # Миграции базы данных
├── static/               # Статические файлы (frontend)
└── docs/                 # Swagger документация
```

### Слои архитектуры

1. **Handler Layer** - HTTP обработчики, валидация запросов
2. **Service Layer** - Бизнес-логика, координация между компонентами
3. **Repository Layer** - Работа с базой данных и хранилищем
4. **Worker Pool** - Асинхронная обработка изображений

## 📡 API Endpoints

### 1. Загрузка изображения

**POST** `/upload`

Загружает изображение для обработки с метаданными.

**Параметры формы:**
- `image` (file) - файл изображения
- `metadata` (string) - JSON метаданные обработки

**Пример curl:**
```bash
curl -X POST http://localhost:8080/upload \
  -F "image=@image.jpg" \
  -F 'metadata={"file_name":"image.jpg","content_type":"image/jpeg","watermark_string":"Sample","task":"watermark"'
```

**Возможные задачи (task):**
- `resize` - изменить размер
- `watermark` - добавить водяной знак
- `miniature generating` - создать миниатюру

### 2. Получение обработанного изображения

**GET** `/image/{id}`

Скачивает обработанное изображение по ID.

**Пример curl:**
```bash
curl -X GET http://localhost:8080/image/550e8400-e29b-41d4-a716-446655440000 \
  -o processed_image.jpg
```

### 3. Получение информации об изображении

**GET** `/image/info/{id}`

Возвращает статус и метаданные изображения.

**Пример curl:**
```bash
curl -X GET http://localhost:8080/image/info/550e8400-e29b-41d4-a716-446655440000
```

**Пример ответа:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "completed",
  "create_at": "2024-01-15T10:30:00Z"
}
```

### 4. Удаление изображения

**DELETE** `/image/{id}`

Удаляет изображение и связанные данные.

**Пример curl:**
```bash
curl -X DELETE http://localhost:8080/image/550e8400-e29b-41d4-a716-446655440000
```

### 5. Главная страница

**GET** `/`

Возвращает веб-интерфейс для работы с API.

**Пример curl:**
```bash
curl -X GET http://localhost:8080/
```

### 6. Swagger документация

**GET** `/swagger/*`

Интерактивная документация API.

**Пример curl:**
```bash
curl -X GET http://localhost:8080/swagger/index.html
```

## 🚀 Установка и запуск


### 1. Клонирование репозитория

```bash
git clone https://github.com/Komilov31/image-processor.git
cd image-processor
```

### 2. Создание .env файла

```bash
cp .env.example .env
```

Заполните `.env` файл с необходимыми переменными окружения:

```env
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=image-processor

# MinIO
MINIO_USER=minioadmin
MINIO_PASSWORD=minioadmin

# Goose migration
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=/migrations
```

### 3. Запуск сервисов

```bash
docker-compose up -d
```

### 4. Проверка работоспособности

- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- MinIO Console: http://localhost:9001

## 💻 Использование

### Через веб-интерфейс

1. Откройте http://localhost:8080
2. Выберите изображение для загрузки
3. Укажите параметры обработки:
   - Тип обработки (resize, watermark, miniature)
   - Для resize: укажите ширину и высоту
   - Для watermark: введите текст водяного знака
   - Выберите формат изображения
4. Нажмите "Отправить на обработку"
5. Получите ID обработанного изображения
6. Используйте ID для скачивания или удаления

### Через API

#### Пример полной обработки изображения

```bash
# 1. Загрузка изображения
RESPONSE=$(curl -s -X POST http://localhost:8080/upload \
  -F "image=@sample.jpg" \
  -F 'metadata={"file_name":"sample.jpg","content_type":"image/jpeg","watermark_string":"CONFIDENTIAL","task":"watermark","resize":{"width":800,"height":600}}')

IMAGE_ID=$(echo $RESPONSE | jq -r '.id')
echo "Image ID: $IMAGE_ID"

# 2. Проверка статуса обработки
curl -X GET http://localhost:8080/image/info/$IMAGE_ID

# 3. Скачивание обработанного изображения
curl -X GET http://localhost:8080/image/$IMAGE_ID -o processed_sample.jpg

# 4. Удаление изображения (при необходимости)
curl -X DELETE http://localhost:8080/image/$IMAGE_ID
```

## 🧪 Тестирование

### Запуск тестов
```bash
go test ./...
```