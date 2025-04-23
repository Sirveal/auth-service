# Часть сервиса аутентификации

Часть сервиса аутентификации, реализующий механизм access, refresh token с IP-based защитой.

## 🚀 Features

- **JWT Access Tokens** с алгоритмом SHA512
- **Безопасные Refresh Tokens**, хранящиеся в виде bcrypt хешей
- **IP-based Security** с email уведомлениями
- Поддержка **Docker**
- **PostgreSQL** для хранения данных
- **Тестовое покрытие** с моковыми компонентами

## 🛠 Технологии

- Go 1.24
- JWT
- PostgreSQL
- Docker & Docker Compose
- bcrypt
- Gin Web Framework

## 📋 Требования

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL (для локального запуска)

## 🚀 Быстрый старт

1. Клонируем репозиторий:
```bash
git clone https://github.com/yourusername/auth-service.git
cd auth-service
```

2. Запускаем сервисы через Docker:
```bash
docker-compose up -d
```

3. Тестируем API:
```bash
# Получаем токены по UUID
curl -X GET "http://localhost:8000/auth/token?uuid=<user-uuid>"

# Обновляем токены
curl -X POST "http://localhost:8000/auth/refresh" \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "your-refresh-token"}'
```

## 🔑 API Endpoints

### GET /auth/token
Получение новой пары access и refresh токенов.

**Параметры запроса:**
- `uuid` (reqобязательный): UUID пользователя

**Ответ:**
```json
{
    "access_token": "jwt-token",
    "refresh_token": "base64-encoded-token"
}
```

### POST /auth/refresh
Обновление существующей пары токенов.

**Тело запроса:**
```json
{
    "refresh_token": "your-refresh-token"
}
```

**Ответ:**
```json
{
    "access_token": "new-jwt-token",
    "refresh_token": "new-refresh-token"
}
```

## 🔒 Функции безопасности

- Access token в формате JWT с алгоритмом SHA512
- Refresh token хранятся в виде bcrypt хешей
- Токены связаны для предотвращения повторного использования
- Отслеживание IP адреса с email уведомлениями
- Защита от повторного использования токенов

## 🧪 Тестирование

Запуск тестов:
```bash
# Запуск с покрытием
go test ./pkg/service -cover -v
```

## 📁 Структура проекта

```
.
├── cmd/
│   └── main.go
├── pkg/
│   ├── handler/
│   ├── repository/
│   └── service/
├── schema/
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
├── configs/
│   └── config.yml
├── docker-compose.yml
└── Dockerfile
```

## ⚙️ Конфигурация

Настройка через `config.yml` и переменные окружения:

```yaml
port: "8000"
auth:
    signing_key: "your-signing-key"
    token_ttl: "15m"
    refresh_token_ttl: "720h"
db:
    username: "postgres"
    host: "localhost"
    port: "5432"
    dbname: "test"
    sslmode: "disable"
```

## 📝 Лицензия

Проект распространяется под лицензией MIT.