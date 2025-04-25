# xlink

![Go Version](https://img.shields.io/badge/go-1.24-blue)

xlink — сокращатель ссылок с аналитикой

## О проекте

Проект построен на базе Go и использует gRPC для коммуникации между сервисами.   
В проекте реализованы следующие сервисы:
- **User Service** — управление пользователями и авторизация.
- **Shortener** — создание и редирект коротких ссылок.
- **Analytics** — сбор статистики и аналитика.
- **API Gateway** — маршрутизация запросов к микросервисам.
- **Telegram Bot** — взаимодействие с пользователем через Telegram.

## Документация

- [Конфигурация проекта](./docs/configs.md) - Описание переменных окружения и настроек.
- [Для разработчиков](./docs/for_dev.md) - Инструкции по Makefile и окружению.
- [Endpoints: Shortener](./docs/endpoints/shortener.md) - Документация по API эндпоинтам сервиса `shortener`.
- [Endpoints: Analytics](./docs/endpoints/analytics.md) - Документация по API эндпоинтам сервиса `analytics`.
- [Endpoints: User Service](./docs/endpoints/user_service.md) - Документация по API эндпоинтам сервиса `user_service`.

## Запуск проекта

### Локальный запуск

1. Клонируйте репозиторий:
   ```bash
   git clone https://gitlab.crja72.ru/golang/2025/spring/course/projects/go10/xlink.git
   ```

2. Откройте `configs/config.yaml` и по желанию измените, затем выполните:
   ```bash
   make yaml_to_env
   ```

3. Запустите проект с Docker Compose:
   ```bash
   docker compose up
   ```
   
## Миграции

Миграции стартую автоматически при запуске проекта.  
Для каждого микросервиса миграции находятся по следующим путям:
- **user_service:** `user_service/migrations`
- **shortener:** `shortener/migrations`
- **analytics:** `analytics/migrations`

## Структура проекта

```
xlink
├── analytics        # микросервис аналитики
│   ├── api          # proto файл
│   ├── cmd          # main.go
│   ├── internal     
│   │   ├── config   # загрузка конфигов
│   │   ├── models   # доменные модели
│   │   ├── ports    # контракты и адаптеры (Kafka, ClickHouse, Redis, shortener)
│   │   ├── server   # HTTP/gRPC-сервер
│   │   └── service  # бизнес-логика и хелперы
│   ├── migrations   # SQL-скрипты миграций БД
│   └── tests        # тесты
├── build            # файлы для сборки и деплоя в Docker
│   ├── docker       # Dockerfile'ы и docker-compose
│   └── entry_nginx  # конфигурация Nginx для API Gateway
├── common           # общие библиотеки и утилиты
│   ├── callers      # retry/timeout обёртки
│   ├── clickhouse   # клиент ClickHouse
│   ├── gen          # сгенерированные gRPC-код
│   ├── grpc         # интерсепторы и пул коннектов
│   ├── kafka        # клиент Kafka
│   ├── logger       # логгер
│   ├── postgres     # клиент Postgres
│   └── redis        # клиент Redis
├── configs          # примеры и шаблоны конфигов
├── docs             # документация по API и девопсу
│   ├── configs.md   # описание конфигов
│   ├── endpoints    # спецификации эндпойнтов (analytics, shortener, user_service)
│   └── for_dev.md   # общие инструкции для девелоперов
├── gateway          # API Gateway 
│   ├── cmd          # main.go для gateway
│   └── internal     # хэндлеры, middleware, порты и схемы
├── scripts          # утилиты (yaml_to_env)
├── shortener        # микросервис сокращателя ссылок
│   ├── api          # proto файл
│   ├── cmd          # main.go
│   ├── internal     
│   │   ├── config   # загрузка конфигов
│   │   ├── models   # доменные модели
│   │   ├── ports    # контракты и адаптеры (Kafka, ClickHouse, Redis, user_service)
│   │   └── service  # бизнес-логика и хелперы
│   ├── migrations   # SQL-миграции для links
│   └── tests        # тесты
├── tg_bot           # Telegram-бот для взаимодействия с микросервисами
│   ├── cmd          # main.go
│   └── internal     # handler, модели, адаптеры портов
└── user_service     # микросервис авторизации
    ├── api          # proto файл
    ├── cmd          # main.go
    ├── internal     
    │   ├── config   # загрузка конфигов
    │   ├── ports    # контракты и адаптеры (Kafka, Postgres, Redis, shortener)
    │   ├── runner   # инициализация сервисов
    │   ├── service  # бизнес-логика
    │   └── utils    # хелперы
    ├── migrations   # SQL-миграции для users
    └── tests        # тесты
```