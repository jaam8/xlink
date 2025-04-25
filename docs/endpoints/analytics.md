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
      "stats": [
        {
          "country": "US",
          "clicks": "1",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-20"
    },
    {
      "stats": [
        {
          "country": "DE",
          "clicks": "1",
          "unique_clicks": "0"
        },
        {
          "country": "US",
          "clicks": "2",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-21"
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
      "stats": [
        {
          "region": "California",
          "clicks": "1",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-20"
    },
    {
      "stats": [
        {
          "region": "Moscow",
          "clicks": "1",
          "unique_clicks": "0"
        },
        {
          "region": "California",
          "clicks": "2",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-21"
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
      "stats": [
        {
          "browser": "Chrome",
          "clicks": "1",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-20"
    },
    {
      "stats": [
        {
          "browser": "Opera",
          "clicks": "1",
          "unique_clicks": "0"
        },
        {
          "browser": "Chrome",
          "clicks": "2",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-21"
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
      "stats": [
        {
          "os": "Windows",
          "clicks": "1",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-20"
    },
    {
      "stats": [
        {
          "os": "Linux",
          "clicks": "1",
          "unique_clicks": "0"
        },
        {
          "os": "Windows",
          "clicks": "2",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-21"
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
      "stats": [
        {
          "device_type": "Desktop",
          "clicks": "1",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-20"
    },
    {
      "stats": [
        {
          "device_type": "Desktop",
          "clicks": "1",
          "unique_clicks": "0"
        },
        {
          "device_type": "Mobile",
          "clicks": "2",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-21"
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
      "stats": [
        {
          "hour": "10:00",
          "clicks": "1",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-20"
    },
    {
      "stats": [
        {
          "hour": "10:00",
          "clicks": "1",
          "unique_clicks": "0"
        },
        {
          "hour": "11:00",
          "clicks": "2",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-21"
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
      "stats": [
        {
          "date": "2025-04-20",
          "clicks": "1",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-20"
    },
    {
      "stats": [
        {
          "date": "2025-04-21",
          "clicks": "2",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-21"
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
      "stats": [
        {
          "referrer": "https://example.com",
          "clicks": "1",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-20"
    },
    {
      "stats": [
        {
          "referrer": "https://example.com",
          "clicks": "1",
          "unique_clicks": "0"
        },
        {
          "referrer": "https://example-2.com",
          "clicks": "2",
          "unique_clicks": "1"
        }
      ],
      "date": "2025-04-21"
    }
  ]
}
```

**Response Codes:**
- `200 OK` — Статистика успешно получена.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `500 Internal Server Error` — Ошибка на сервере.