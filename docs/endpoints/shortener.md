### Shortener

#### Common Part

##### GET /l/:shortLink

**Description:** Редирект по короткой ссылке.

**Request:**
```bash
curl -X GET http://localhost:8080/l/shortLink
```

**Response:**
- **302 Found** — Редирект на целевой URL.
- **404 Not Found** — Ссылка не найдена.
- **400 Bad Request** — Пустая ссылка.

---

##### POST /api/v1/link/create

**Description:** Создание новой короткой ссылки.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/link/create \
-H "Authorization: Bearer token" \
-H "Content-Type: application/json" \
-d '{"short_link": "customLink", "target_url": "https://example.com"}'
```

**Headers:**

| Header          | Required | Description               |
|-----------------|----------|---------------------------|
| `Authorization` | Yes      | Токен авторизации.        |

**Request Body:**

| Field        | Type   | Required | Description               |
|--------------|--------|----------|---------------------------|
| `short_link` | string | No       | Кастомная короткая ссылка.|
| `target_url` | string | Yes      | Целевой URL.             |

**Response:**
- **201 Created**
```json
{
  "link_id": "uuid",
  "user_id": "uuid",
  "short_link": "customLink",
  "target_url": "https://example.com",
  "created_at": "2023-01-01T00:00:00Z",
  "expire_at": "2023-12-31T23:59:59Z"
}
```

**Response Codes:**
- `201 Created` — Ссылка успешно создана.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `422 Unprocessable Entity` — Некорректное тело запроса.
- `500 Internal Server Error` — Ошибка на сервере.

---

#### For Authenticated Users

##### GET /api/v1/link/my-links/

**Description:** Получения списка своих ссылок.

**Headers:**

| Header          | Required | Description               |
|-----------------|----------|---------------------------|
| `Authorization` | Yes      | Токен авторизации.        |

**Response:**
- **200 Ok**
```json
[
  "some-short-link-part1",
  "some-short-link-part2"
]
```

**Response Codes:**
- `200 Ok` — Получен список ссылок.
- `400 Bad Request` — Неверный формат данных или ошибка на сервере.
- `401 Unauthorized` — Неверный токен.

---

#### For Link Owner

##### PUT /api/v1/link/update/:id

**Description:** Обновление короткой ссылки.

**Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/link/update/uuid \
-H "Authorization: Bearer token" \
-H "Content-Type: application/json" \
-d '{"regenerate": false, "short_link": "newLink", "target_url": "https://newexample.com", "expire_at": "2024-01-01T00:00:00Z"}'
```

**Headers:**

| Header          | Required | Description               |
|-----------------|----------|---------------------------|
| `Authorization` | Yes      | Токен авторизации.        |

**Request Body:**

| Field        | Type   | Required | Description                |
|--------------|--------|----------|----------------------------|
| `regenerate` | bool   | Yes      | Генерировать новую ссылку. |
| `short_link` | string | No       | Новая короткая ссылка.     |
| `target_url` | string | No       | Новый целевой URL.         |
| `expire_at`  | string | Yes      | Дата истечения (ISO 8601). |

**Response:**
- **200 OK**
```json
{
  "link_id": "uuid",
  "user_id": "uuid",
  "short_link": "newLink",
  "target_url": "https://newexample.com",
  "created_at": "2023-01-01T00:00:00Z",
  "expire_at": "2024-01-01T00:00:00Z"
}
```

**Response Codes:**
- `200 OK` — Ссылка успешно обновлена.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `404 Not Found` — Ссылка не найдена.
- `422 Unprocessable Entity` — Некорректное тело запроса.
- `500 Internal Server Error` — Ошибка на сервере.

---

##### DELETE /api/v1/link/delete/:shortLink

**Description:** Удаление короткой ссылки владельцем.

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/link/delete/shortLink \
-H "Authorization: Bearer token"
```

**Headers:**

| Header          | Required | Description               |
|-----------------|----------|---------------------------|
| `Authorization` | Yes      | Токен авторизации.        |

**Response Codes:**
- `204 No Content` — Ссылка успешно удалена.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `404 Not Found` — Ссылка не найдена.
- `500 Internal Server Error` — Ошибка на сервере.

---

#### For Admins

##### GET /api/v1/link/admin/links/:userId

**Description:** Получения списка ссылок нужного пользователя администратором.

**Headers:**

| Header          | Required | Description               |
|-----------------|----------|---------------------------|
| `Authorization` | Yes      | Токен авторизации.        |

**Response:**
- **200 Ok**
```json
[
  "some-short-link-part1",
  "some-short-link-part2"
]
```

**Response Codes:**
- `200 Ok` — Получен список ссылок.
- `400 Bad Request` — Неверный формат данных или ошибка на сервере.
- `401 Unauthorized` — Неверный токен.
- `404 Not Found` — Пользователь не найден.

##### DELETE /api/v1/link/admin/delete/:id

**Description:** Удаление короткой ссылки администратором.

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/link/admin/delete/uuid \
-H "Authorization: Bearer token"
```

**Headers:**

| Header          | Required | Description               |
|-----------------|----------|---------------------------|
| `Authorization` | Yes      | Токен авторизации.        |

**Response Codes:**
- `204 No Content` — Ссылка успешно удалена.
- `400 Bad Request` — Неверный формат данных или ошибка на сервере.
- `401 Unauthorized` — Неверный токен.
- `404 Not Found` — Ссылка не найдена.

---