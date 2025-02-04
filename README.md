# Jira Service

Этот сервис предназначен для интеграции с Jira API. Он предоставляет REST API для получения информации о задачах.

## 📌 Возможности
- Получение информации о конкретной задаче.
- Получение списка задач.
- Логирование запросов и ошибок.
- Гибкая конфигурация через переменные окружения.

## 🚀 Технологии
- **Язык**: Go
- **Фреймворк**: Fiber (gofiber.io)
- **HTTP-клиент**: resty (go-resty/resty)
- **Логирование**: zerolog
- **Конфигурация**: godotenv
- **Тестирование**: testify
- **Контейнеризация**: Docker

## 📂 Структура проекта
```
├── cmd/                     # Основной исполняемый файл
│   ├── main.go              # Точка входа в приложение
├── config/                  # Конфигурационные файлы
│   ├── config.go            # Логика загрузки конфигурации
├── internal/                # Внутренние модули сервиса
│   ├── adapter/             # Взаимодействие с API внешних систем
│   │   ├── client.go        # HTTP-клиент для работы с Jira API
│   │   ├── models.go        # Определение структур данных
│   ├── handler/             # HTTP-обработчики
│   │   ├── handlers.go      # Основные обработчики запросов
│   ├── service/             # Бизнес-логика
│   │   ├── service.go       # Jira-сервис
├── test/                    # Тесты
│   ├── client_test.go       # Тест HTTP-клиента Jira
│   ├── config_test.go       # Тест конфигурации
│   ├── handler_test.go      # Тест HTTP-обработчиков
│   ├── integration_test.go  # Интеграционные тесты
│   ├── service_test.go      # Тест сервисного слоя
├── .env                     # Файл с переменными окружения
├── .gitignore               # Файл игнорирования в Git
├── Dockerfile               # Docker-контейнеризация
├── go.mod                   # Файл зависимостей Go
├── go.sum                   # Контрольные суммы зависимостей
├── README.md                # Описание проекта
```

## ⚙️ Установка и запуск
### 🔧 Настройка переменных окружения
Перед запуском сервиса создайте файл `.env` с настройками:
```env
SERVER_PORT=8080
JIRA_BASE_URL=https://jira.example.com
JIRA_API_URL=/rest/api/latest/
JIRA_API_TOKEN=your_api_token
```

### 🏃‍♂️ Локальный запуск
```sh
go run cmd/main.go
```

### 🐳 Запуск в Docker
```sh
docker build -t jira-service .
docker run -p 8080:8080 --env-file .env jira-service
```

## 🛠 Тестирование
### ✅ Запуск юнит-тестов
```sh
go test ./test/... -coverpkg=./... -coverprofile=./coverage/coverage.out
```

### 📊 Анализ покрытия кода тестами
```sh
go tool cover -func=./coverage/coverage.out
```

### 🏆 Генерация HTML-отчета
```sh
go tool cover -html=./coverage/coverage.out -o ./coverage/report.html
```

## 📄 API эндпоинты
| Метод  | URL                  | Описание                |
|--------|----------------------|-------------------------|
| GET    | `/issue/:id`         | Получение задачи по ID  |
| GET    | `/issues?issueIDs=1,2,3` | Получение списка задач |

## ✨ Авторы
- **Виктория Пилипейко** — Разработка и проектирование сервиса
