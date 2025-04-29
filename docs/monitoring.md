## 📊 Monitoring
![](https://img.shields.io/badge/Prometheus-000000?style=for-the-badge&logo=prometheus&labelColor=000000)
![](https://img.shields.io/badge/Grafana-F2F4F9?style=for-the-badge&logo=grafana&logoColor=orange&labelColor=F2F4F9)

В проекте настроено наблюдение за HTTP-запросами в API-gateway 
c помощью связки **Prometheus** и **Grafana**. 
Это позволяет в реальном времени отслеживать количество 
запросов и задержки по каждому эндпоинту.

### Архитектура

- **Prometheus** собирает метрики с API Gateway
- **Grafana** визуализирует эти метрики
- Метрики доступны по эндпоинту `/metrics`, который хендлится через middleware

### Как это работает

1. В API Gateway добавлен middleware, собирающий:
    - количество запросов `http_requests_total`
    - длительность выполнения `http_duration_seconds`
2. Prometheus настроен на опрос `/metrics`
3. Grafana подключается к Prometheus и строит графики

### 📈 Дашборды Grafana

Дашборд Grafana находится в папке `build/grafana` и включают в себя:

- Панель **HTTP Requests per Second**: считает количество запросов по методам и маршрутам.
- Панель **HTTP Duration (p95)**: показывает 95-й перцентиль времени отклика.

### Как запустить
1. Запустите проект с помощью 
```bash
cd build/docker
docker compose up
```
2. Перейдите в Grafana по адресу [`http://localhost:3001`](http://localhost:3001)
3. Войдите с помощью логина и пароля `admin:admin`
4. Перейдите на [`http://localhost:3001/dashboard/import`](http://localhost:3001/dashboard/import)
5. Импортируйте дашборд из файла `build/grafana/dashboard.json`
6. Наслаждайтесь графиками!
