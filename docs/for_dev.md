# Документация по Makefile

Этот документ описывает команды из Makefile и содержит примеры для копирования команд.

## Генерация gRPC кода

Генерирует gRPC код для сервисов: user_service, shortener, analytics в директорию `common/gen/`
```bash
make generate_user_service
```
```bash
make generate_shortener
```
```bash
make generate_analytics
```

## Работа с env файлами

### yaml_to_env
Преобразует YAML конфигурацию в переменные окружения.
```bash
make yaml_to_env
```

### copy_env
Копирует пример файла окружения (.env.example) в настоящий .env.
```bash
make copy_env
```

### update_env_example
Обновляет файл .env.example на основе текущего .env.
```bash
make update_env_example
```

## Сборка и запуск

### build-all
Собирает Docker образы для сервисов: token_service, shortener, tg_bot.
```bash
make build-all
```

### env_for_build
Генерирует .env и копирует его в директорию build/docker для сборки.
```bash
make env_for_build
```

## Линтинг, тестинг и запуск через Docker Compose

Линтинг и тестинг настроены через .gitlab-ci.yml. Для локального запуска проекта можно воспользоваться Docker Compose:
```bash
docker compose up
```
`