# Renderer Endpoints Documentation

## 1. Endpoints

### GET /api/v1/img/:shortLink/?param=&start_date=&end_date=

**Description:** Получение статистики по странам.

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/img/dsalakviqm/?param=browser&start_date=2025-01-01&end_date=2025-04-01&token=djalkdjaslkdaskldjaskl"
```

**Query Parameters:**

| Parameter    | Type   | Required | Description                           |
|--------------|--------|----------|---------------------------------------|
| `token`      | string | Yes      | Токен авторизации (в GET параметрах). |
| `short_link` | string | Yes      | Короткая ссылка (в самом пути).       |
| `start_date` | string | Yes      | Начальная дата (YYYY-MM-DD).          |
| `end_date`   | string | Yes      | Конечная дата (YYYY-MM-DD).           |

**Response:**
- **200 OK** - `html` страница с графиками

**Response Codes:**
- `200 OK` — Статистика успешно получена.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `500 Internal Server Error` — Ошибка на сервере.