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
[//]: # (- **Bot** — взаимодействие с Telegram.)

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

