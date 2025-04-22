# 🧾 configuration 
![.env](https://img.shields.io/badge/.env-black?style=for-the-badge&logo=dotenv&logoColor=white)
![.env](https://img.shields.io/badge/config.yaml-black?style=for-the-badge&logo=yaml&logoColor=white)

Documentation for environment variables to run all microservices.

---

### ![ClickHouse](https://img.shields.io/badge/ClickHouse-FFCC01?style=for-the-badge&logo=clickhouse&logoColor=white)

| variable                    | description              | example      |
|-----------------------------|--------------------------|--------------|
| `CLICKHOUSE_HOST`           | ClickHouse host          | `clickhouse` |
| `CLICKHOUSE_PORT`           | ClickHouse port          | `9000`       |
| `CLICKHOUSE_USER`           | username                 | `default`    |
| `CLICKHOUSE_PASSWORD`       | password                 | `""`         |
| `CLICKHOUSE_DB`             | database name            | `default`    |
| `CLICKHOUSE_MAX_OPEN_CONNS` | max open connections     | `5`          |
| `CLICKHOUSE_MAX_IDLE_CONNS` | max inactive connections | `5`          |

---

### ![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-000?style=for-the-badge&logo=apachekafka)

| variable                   | description                 | example      |
|----------------------------|-----------------------------|--------------|
| `KAFKA_HOST`               | kafka host                  | `kafka`      |
| `KAFKA_PORT`               | kafka port                  | `9092`       |
| `KAFKA_BROKERS`            | list brokers                | `kafka:9092` |
| `KAFKA_MIN_BYTES`          | min size message            | `10`         |
| `KAFKA_MAX_BYTES`          | max size message            | `1048576`    |
| `KAFKA_MAX_WAIT_MS`        | max wait time message (ms)  | `500`        |
| `KAFKA_COMMIT_INTERVAL_MS` | commit interval offset (ms) | `1000`       |

---

### ![Zookeeper](https://img.shields.io/badge/ZooKeeper-600a6e?style=for-the-badge&logo=apache&logoColor=white)

| variable         | description    | example     |
|------------------|----------------|-------------|
| `ZOOKEEPER_HOST` | Zookeeper host | `zookeeper` |
| `ZOOKEEPER_PORT` | Zookeeper port | `2181`      |

---

### ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)

| variable             | description     | example    |
|----------------------|-----------------|------------|
| `POSTGRES_HOST`      | postgres host   | `postgres` |
| `POSTGRES_PORT`      | postgres port   | `5432`     |
| `POSTGRES_DB`        | database name   | `xlink`    |
| `POSTGRES_USER`      | user            | `postgres` |
| `POSTGRES_PASSWORD`  | password        | `1234`     |
| `POSTGRES_MAX_CONNS` | max connections | `10`       |
| `POSTGRES_MIN_CONNS` | min connections | `5`        |

---

### ![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)

| variable                      | description                 | example |
|-------------------------------|-----------------------------|---------|
| `REDIS_HOST`                  | redis host                  | `redis` |
| `REDIS_PORT`                  | redis port                  | `6379`  |
| `REDIS_USER`                  | user                        | `""`    |
| `REDIS_PASSWORD`              | password                    | `""`    |
| `REDIS_USER_PASSWORD`         | ?password for certain user? | `""`    |
| `REDIS_MAX_RETRIES`           | amount of connect attempt   | `3`     |
| `REDIS_POOL_SIZE`             | pool connection size        | `10`    |
| `REDIS_DIAL_TIMEOUT_SECONDS`  | dial timeout (sec)          | `5`     |
| `REDIS_READ_TIMEOUT_SECONDS`  | read timeout (sec)          | `3`     |
| `REDIS_WRITE_TIMEOUT_SECONDS` | write timeout (sec)         | `3`     |

---

### ![Nginx](https://img.shields.io/badge/nginx-%23009639.svg?style=for-the-badge&logo=nginx&logoColor=white)

| variable           | description         | example |
|--------------------|---------------------|---------|
| `NGINX_PORT_HTTP`  | port for HTTP       | `81`    |
| `NGINX_PORT_GRPC`  | port for gRPC       | `50050` |
| `NGINX_ENTRY_PORT` | external entry port | `80`    |

---

### ![API gateway](https://img.shields.io/badge/API_gateway-FF4F8B.svg?style=for-the-badge&logo=amazonapigateway&logoColor=white)

| variable       | description      | example |
|----------------|------------------|---------|
| `GATEWAY_PORT` | API gateway port | `3000`  |

---

### ![user service](https://img.shields.io/badge/user_service-0E353D.svg?style=for-the-badge&logo=monkeytie)

| variable                                | description        | example                           |
|-----------------------------------------|--------------------|-----------------------------------|
| `USER_SERVICE_GRPC_PORT`                | gRPC port          | `50051`                           |
| `USER_SERVICE_REDIS_DB`                 | redis table        | `0`                               |
| `USER_SERVICE_TOKEN_LEN`                | token length       | `32`                              |
| `USER_SERVICE_CACHE_EXPIRATION_SECONDS` | ?? (sec)           | `3600`                            |
| `USER_SERVICE_MIGRATIONS_PATH`          | path to migrations | `file:///user_service/migrations` |

---

### ![shortener](https://img.shields.io/badge/shortener-blue.svg?style=for-the-badge&logo=data:image/svg%2b;base64,iVBORw0KGgoAAAANSUhEUgAAACoAAAAqCAYAAADFw8lbAAAAAXNSR0IArs4c6QAAA8tJREFUWIXN2VvIVFUYxvEfGmTlMU3xXJZpJ7MDVBQmRWlUdBOVXUSQXaTRCUqpiC4y0aKIssKkyMpIyqALrYwyMCnsCFZIJxU7aObZDLqYbp6J7Xwz38wHzZ4eWMzaa68173/2Wut937Wnt85qCE7F4djVYZYuGoJ5+A6VQtmLl3B2pwFhJv4M2O94FwuxDF8XoF/GUZ2CvDcQf2BG2k7GVEzJ9Xh8lH5foW/ZkJfH+LcYhUlYWTP1m3Fj+j+RttfKBt0Yw8fhphrA2rI6Y1bn+pyyIKfH4EKMaQJZLXMDWMHSskAXxeApmN8i6KaM/aHqtnqVADoWf2dXn9WDMUOwFgMxqAzQUdiXeu8ejDsyHgL6tQt0dJz6NpyOo9EnTr4VHcCWbD7Y3l3nSbgGc3AdJrdg4BK80WDdTccVLa7RZfm+XfiykbFbsKPBF+zErJr+/XF7nZBYW95P/xea9NuMofGpFdxdD3JJYcBzeUJTcTGeLtx7JU94CQ62+JSqLgoWN7i/ESdgYmL/9nrRaW46r8MIjMM9eDDTPwHD8EEPwOqVR2LvDDyAt/AYrkr7ldlElUS0Q3RSbnyPfri2gZEb0v/zHoB9lmmcjA1p24DZODfJx8Tsh1cL42bUQkq2UskUT2pi+ILCD+uuLA1IUX3xaJPlsjLBoa524dPU1zUBWJ9+a+vc+ylLaHAjQ1F/XI+nsCo/6i4c392gYTHyZKHerAzNWqter4rraZsOwxGpH+xB/jcCe1KfjnfaxPeveuG31EcnCdjXZMxf2XRjc/1NmxkP0Xrszg58tsm0P584XA0Ks8sEvbXgkIclzaoHuQ0j8VBN+9sYXhZs9XB1NQYl6hxI20G8mE10WYMfsbOR7/uvNT7TXw2fY2ruj4xnaOYRXk+21FZNxI8Fo2uxpsZn/twC7K+4tN2wfXAH3sypsRgKZ6bPtHiLZsCLO3HkrdXgOPpmsFvKPEl2p1ktRrMFuBPvZWnsxYdpP6Ys2An4okXgSvbBGmzN9Y5EuNK0oAngffHTVfVJ6lhd79PKhL2wwUa7P5Hv8Rw39idHOC3voPYnoAwsE3YgVhQgP057o4PfmMJav61M0KrmxfjNyfAbLYlFGFAIGF3U7hcQW/O5CWd2029y0sZfEgG7qN2gu/M5PDu9kaovJkZkrZau8ZnO5bne0GDqz8sptIKHOwEqobeSrOvEmhPsgYTlATlvVXJw7IguKoDNSdu4wpu98ws5xTOdgqxqSv5cqOatK5IyFv9gmN9pyKqG5uXXnpr1+Ume+v9SxzZxV130Dzcj6RC4AkJMAAAAAElFTkSuQmCC&color=white)

| variable                                    | description                     | example                        |
|---------------------------------------------|---------------------------------|--------------------------------|
| `SHORTENER_GRPC_PORT`                       | gRPC port                       | `50052`                        |
| `SHORTENER_REDIS_DB`                        | redis table                     | `1`                            |
| `SHORTENER_KAFKA_TOPIC`                     | kafka topic for sending events  | `redirect-events`              |
| `SHORTENER_KAFKA_NUM_PARTITIONS`            | amount of kafka partitions      | `1`                            |
| `SHORTENER_KAFKA_REPLICATION_FACTOR`        | amount of kafka replications    | `1`                            |
| `SHORTENER_UPSTREAM_NAMES`                  | upstream service name           | `shortener_upstream`           |
| `SHORTENER_UPSTREAM_PORTS`                  | upstream port                   | `50060`                        |
| `SHORTENER_TIMEOUTS`                        | timeout for internal calls (ms) | `10000`                        |
| `SHORTENER_EXPIRATION_SECONDS`              | ??                              | `500`                          |
| `SHORTENER_DEFAULT_LINK_EXPIRATION_MINUTES` | ?? (min)                        | `1440`                         |
| `SHORTENER_MIGRATIONS_PATH`                 | path to migrations              | `file:///shortener/migrations` |

---

### ![analytics](https://img.shields.io/badge/Analytics-1A1A1A.svg?style=for-the-badge&logo=googleanalytics&logoColor=white)

| variable                    | description                     | example                        |
|-----------------------------|---------------------------------|--------------------------------|
| `ANALYTICS_GRPC_PORT`       | gRPC port                       | `50053`                        |
| `ANALYTICS_KAFKA_TOPIC`     | kafka topic                     | `redirect-events`              |
| `ANALYTICS_KAFKA_GROUP_ID`  | id consumer-group in kafka      | `analytics-consumer`           |
| `ANALYTICS_REDIS_DB`        | redis table                     | `2`                            |
| `ANALYTICS_UPSTREAM_NAMES`  | upstream service name           | `analytics_upstream`           |
| `ANALYTICS_UPSTREAM_PORTS`  | upstream port                   | `50061`                        |
| `ANALYTICS_MIGRATIONS_PATH` | path to migrations              | `file:///analytics/migrations` |
| `ANALYTICS_TIMEOUTS`        | timeout for internal calls (ms) | `10000`                        |

---

> 🧠 All variables can be set local via `make copy_env` command

