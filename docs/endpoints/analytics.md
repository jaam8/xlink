# Analytics Endpoints Documentation

## 1. Endpoints

### GET /api/v1/analytics/by-country

**Description:** Получение статистики по странам.

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/analytics/by-country?short_link=shortLink&start_date=2023-01-01&end_date=2023-01-31" \
-H "Authorization: Bearer token"
```

**Headers:**

| Header          | Required | Description        |
|-----------------|----------|--------------------|
| `Authorization` | Yes      | Токен авторизации. |

**Query Parameters:**

| Parameter     | Type   | Required | Description                  |
|---------------|--------|----------|------------------------------|
| `short_link`  | string | Yes      | Короткая ссылка.             |
| `start_date`  | string | Yes      | Начальная дата (YYYY-MM-DD). |
| `end_date`    | string | Yes      | Конечная дата (YYYY-MM-DD).  |

**Response:**
- **200 OK**
```json
{
  "data": [
    {
      "country": "US",
      "clicks": 100,
      "unique_clicks": 80
    },
    {
      "country": "RU",
      "clicks": 50,
      "unique_clicks": 40
    }
  ]
}
```

**Response Codes:**
- `200 OK` — Статистика успешно получена.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `500 Internal Server Error` — Ошибка на сервере.

---

### GET /api/v1/analytics/by-region

**Description:** Получение статистики по регионам.

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/analytics/by-region?short_link=shortLink&start_date=2023-01-01&end_date=2023-01-31" \
-H "Authorization: Bearer token"
```

**Headers:**

| Header          | Required | Description        |
|-----------------|----------|--------------------|
| `Authorization` | Yes      | Токен авторизации. |

**Query Parameters:**

| Parameter     | Type   | Required | Description                  |
|---------------|--------|----------|------------------------------|
| `short_link`  | string | Yes      | Короткая ссылка.             |
| `start_date`  | string | Yes      | Начальная дата (YYYY-MM-DD). |
| `end_date`    | string | Yes      | Конечная дата (YYYY-MM-DD).  |

**Response:**
- **200 OK**
```json
{
  "data": [
    {
      "region": "California",
      "clicks": 60,
      "unique_clicks": 50
    },
    {
      "region": "Moscow",
      "clicks": 40,
      "unique_clicks": 30
    }
  ]
}
```

**Response Codes:**
- `200 OK` — Статистика успешно получена.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `500 Internal Server Error` — Ошибка на сервере.

---

### GET /api/v1/analytics/by-browser

**Description:** Получение статистики по браузерам.

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/analytics/by-browser?short_link=shortLink&start_date=2023-01-01&end_date=2023-01-31" \
-H "Authorization: Bearer token"
```

**Headers:**

| Header          | Required | Description        |
|-----------------|----------|--------------------|
| `Authorization` | Yes      | Токен авторизации. |

**Query Parameters:**

| Parameter     | Type   | Required | Description                  |
|---------------|--------|----------|------------------------------|
| `short_link`  | string | Yes      | Короткая ссылка.             |
| `start_date`  | string | Yes      | Начальная дата (YYYY-MM-DD). |
| `end_date`    | string | Yes      | Конечная дата (YYYY-MM-DD).  |

**Response:**
- **200 OK**
```json
{
  "data": [
    {
      "browser": "Chrome",
      "clicks": 120,
      "unique_clicks": 100
    },
    {
      "browser": "Firefox",
      "clicks": 30,
      "unique_clicks": 25
    }
  ]
}
```

**Response Codes:**
- `200 OK` — Статистика успешно получена.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `500 Internal Server Error` — Ошибка на сервере.

---

### GET /api/v1/analytics/by-os

**Description:** Получение статистики по операционным системам.

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/analytics/by-os?short_link=shortLink&start_date=2023-01-01&end_date=2023-01-31" \
-H "Authorization: Bearer token"
```

**Headers:**

| Header          | Required | Description        |
|-----------------|----------|--------------------|
| `Authorization` | Yes      | Токен авторизации. |

**Query Parameters:**

| Parameter     | Type   | Required | Description                  |
|---------------|--------|----------|------------------------------|
| `short_link`  | string | Yes      | Короткая ссылка.             |
| `start_date`  | string | Yes      | Начальная дата (YYYY-MM-DD). |
| `end_date`    | string | Yes      | Конечная дата (YYYY-MM-DD).  |

**Response:**
- **200 OK**
```json
{
  "data": [
    {
      "os": "Windows",
      "clicks": 80,
      "unique_clicks": 70
    },
    {
      "os": "MacOS",
      "clicks": 50,
      "unique_clicks": 40
    }
  ]
}
```

**Response Codes:**
- `200 OK` — Статистика успешно получена.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `500 Internal Server Error` — Ошибка на сервере.

---

### GET /api/v1/analytics/by-device-type

**Description:** Получение статистики по типам устройств.

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/analytics/by-device-type?short_link=shortLink&start_date=2023-01-01&end_date=2023-01-31" \
-H "Authorization: Bearer token"
```

**Headers:**

| Header          | Required | Description        |
|-----------------|----------|--------------------|
| `Authorization` | Yes      | Токен авторизации. |

**Query Parameters:**

| Parameter     | Type   | Required | Description                  |
|---------------|--------|----------|------------------------------|
| `short_link`  | string | Yes      | Короткая ссылка.             |
| `start_date`  | string | Yes      | Начальная дата (YYYY-MM-DD). |
| `end_date`    | string | Yes      | Конечная дата (YYYY-MM-DD).  |

**Response:**
- **200 OK**
```json
{
  "data": [
    {
      "device_type": "Desktop",
      "clicks": 100,
      "unique_clicks": 90
    },
    {
      "device_type": "Mobile",
      "clicks": 60,
      "unique_clicks": 50
    }
  ]
}
```

**Response Codes:**
- `200 OK` — Статистика успешно получена.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `500 Internal Server Error` — Ошибка на сервере.

---

### GET /api/v1/analytics/by-hour

**Description:** Получение статистики по часам.

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/analytics/by-hour?short_link=shortLink&start_date=2023-01-01&end_date=2023-01-31" \
-H "Authorization: Bearer token"
```

**Headers:**

| Header          | Required | Description        |
|-----------------|----------|--------------------|
| `Authorization` | Yes      | Токен авторизации. |

**Query Parameters:**

| Parameter     | Type   | Required | Description                  |
|---------------|--------|----------|------------------------------|
| `short_link`  | string | Yes      | Короткая ссылка.             |
| `start_date`  | string | Yes      | Начальная дата (YYYY-MM-DD). |
| `end_date`    | string | Yes      | Конечная дата (YYYY-MM-DD).  |

**Response:**
- **200 OK**
```json
{
  "data": [
    {
      "hour": "10:00",
      "clicks": 30,
      "unique_clicks": 25
    },
    {
      "hour": "11:00",
      "clicks": 40,
      "unique_clicks": 35
    }
  ]
}
```

**Response Codes:**
- `200 OK` — Статистика успешно получена.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `500 Internal Server Error` — Ошибка на сервере.

---

### GET /api/v1/analytics/by-date

**Description:** Получение статистики по датам.

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/analytics/by-date?short_link=shortLink&start_date=2023-01-01&end_date=2023-01-31" \
-H "Authorization: Bearer token"
```

**Headers:**

| Header          | Required | Description        |
|-----------------|----------|--------------------|
| `Authorization` | Yes      | Токен авторизации. |

**Query Parameters:**

| Parameter     | Type   | Required | Description                  |
|---------------|--------|----------|------------------------------|
| `short_link`  | string | Yes      | Короткая ссылка.             |
| `start_date`  | string | Yes      | Начальная дата (YYYY-MM-DD). |
| `end_date`    | string | Yes      | Конечная дата (YYYY-MM-DD).  |

**Response:**
- **200 OK**
```json
{
  "data": [
    {
      "date": "2023-01-01",
      "clicks": 50,
      "unique_clicks": 40
    },
    {
      "date": "2023-01-02",
      "clicks": 70,
      "unique_clicks": 60
    }
  ]
}
```

**Response Codes:**
- `200 OK` — Статистика успешно получена.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `500 Internal Server Error` — Ошибка на сервере.

---

### GET /api/v1/analytics/by-referrer

**Description:** Получение статистики по реферерам.

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/analytics/by-referrer?short_link=shortLink&start_date=2023-01-01&end_date=2023-01-31" \
-H "Authorization: Bearer token"
```

**Headers:**

| Header          | Required | Description        |
|-----------------|----------|--------------------|
| `Authorization` | Yes      | Токен авторизации. |

**Query Parameters:**

| Parameter     | Type   | Required | Description                  |
|---------------|--------|----------|------------------------------|
| `short_link`  | string | Yes      | Короткая ссылка.             |
| `start_date`  | string | Yes      | Начальная дата (YYYY-MM-DD). |
| `end_date`    | string | Yes      | Конечная дата (YYYY-MM-DD).  |

**Response:**
- **200 OK**
```json
{
  "data": [
    {
      "referrer": "https://example.com",
      "clicks": 100,
      "unique_clicks": 80
    },
    {
      "referrer": "https://another.com",
      "clicks": 50,
      "unique_clicks": 40
    }
  ]
}
```

**Response Codes:**
- `200 OK` — Статистика успешно получена.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `500 Internal Server Error` — Ошибка на сервере.

`