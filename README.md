# xlink

![Go Version](https://img.shields.io/badge/go-1.24-blue)

xlink — сокращатель ссылок с аналитикой

## О проекте

Проект построен на базе Go и использует gRPC для коммуникации между сервисами  
В проекте реализованы следующие сервисы:
- **User Service** — управление пользователями и авторизация
- **Shortener** — создание и редирект коротких ссылок
- **Analytics** — сбор статистики и аналитика
- **Renderer** — генерация графиков и отчетов
- **API Gateway** — маршрутизация запросов к микросервисам
- **Telegram Bot** — взаимодействие с пользователем через Telegram

## Документация

- [Конфигурация проекта](./docs/configs.md) - Описание переменных окружения и настроек
- [Для разработчиков](./docs/for_dev.md) - Инструкции по Makefile и окружению
- [Endpoints: Shortener](./docs/endpoints/shortener.md) - Документация по API эндпоинтам сервиса `shortener`
- [Endpoints: Analytics](./docs/endpoints/analytics.md) - Документация по API эндпоинтам сервиса `analytics`
- [Endpoints: User Service](./docs/endpoints/user_service.md) - Документация по API эндпоинтам сервиса `user_service`
- [Endpoints: Renderer](./docs/endpoints/renderer.md) - Документация по API эндпоинтам сервиса `renderer`
- [Architecture](./docs/architecture.md) - Архитектура проекта

## Требования
- [Go 1.24](https://go.dev/)
- [Docker](https://www.docker.com/) 
- [Git](https://git-scm.com/)

## Видео как работает бот
- [Тут видео](./img/xlink_bot.mp4)

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

> [!IMPORTANT]  
> Обязательно добавтье токен Telegram-бота в конфиг, иначе бот не запустится

3. Запустите проект с Docker Compose:
   ```bash
   cd build/docker
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
│   ├── api          
│   ├── cmd          
│   ├── internal     
│   │   ├── config   
│   │   ├── models   
│   │   ├── ports    
│   │   ├── server   
│   │   └── service  
│   ├── migrations   
│   └── tests        
├── build            # файлы для сборки и деплоя в Docker
│   ├── docker       
│   └── entry_nginx  
├── common           # общие библиотеки и утилиты
│   ├── callers      
│   ├── clickhouse   
│   ├── gen          
│   ├── grpc         
│   ├── kafka        
│   ├── logger       
│   ├── postgres     
│   └── redis        
├── configs          # примеры и шаблоны конфигов
├── docs             # документация по API и девопсу
│   ├── endpoints    
│   ├── architecture.md    
│   ├── configs.md   
│   └── for_dev.md   
├── gateway          # API Gateway 
│   ├── cmd          
│   ├── internal     
│   │   ├── configs   
│   │   ├── handlers   
│   │   ├── ports 
│   │   ├── schemas 
│   │   └── server
│   └── web 
├── img              # изображения для документации
├── renderer         # сервис для генерации графиков
│   ├── cmd
│   └── internal
│       ├── config
│       ├── handlers
│       ├── ports
│       ├── services
│       └── statistics_data
├── scripts          # утилиты (yaml_to_env)
├── shortener        # микросервис сокращателя ссылок
│   ├── api          
│   ├── cmd          
│   ├── internal     
│   │   ├── config   
│   │   ├── models   
│   │   ├── ports    
│   │   └── service  
│   ├── migrations   
│   └── tests        
├── tg_bot            # tg-бот для взаимодействия с микросервисами
│   ├── cmd          
│   └── internal     
│       ├── config  
│       ├── handler  
│       ├── models   
│       └── ports    
└── user_service      # микросервис авторизации
    ├── api          
    ├── cmd          
    ├── internal     
    │   ├── config   
    │   ├── ports    
    │   ├── runner   
    │   ├── service  
    │   └── utils    
    ├── migrations   
    └── tests        
```
