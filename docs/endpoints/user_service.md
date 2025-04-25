### User Service

#### Common Part

##### POST /api/v1/user/create

**Description:** Создание нового пользователя с `is_staff` и `is_admin`, установленными в `false`.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/user/create \
-H "Content-Type: application/json" \
-d '{"tg_id": 123456789}'
```

**Request Body:**

| Field   | Type   | Required | Description               |
|---------|--------|----------|---------------------------|
| `tg_id` | int64  | Yes      | Telegram ID пользователя. |

**Response:**
- **201 Created**
```json
{
  "user_id": "uuid",
  "token": "string"
}
```

**Response Codes:**
- `201 Created` — Пользователь успешно создан.
- `400 Bad Request` — Неверный формат данных.
- `422 Unprocessable Entity` — Некорректное тело запроса.
- `500 Internal Server Error` — Ошибка на сервере.

---

##### PATCH /api/v1/user/update/:id

**Description:** Обновление данных пользователя с `is_staff` и `is_admin`, установленными в `false`.

**Request:**
```bash
curl -X PATCH http://localhost:8080/api/v1/user/update/uuid \
-H "Content-Type: application/json" \
-d '{"tg_id": 987654321}'
```

**Request Body:**

| Field   | Type   | Required | Description               |
|---------|--------|----------|---------------------------|
| `tg_id` | int64  | Yes      | Новый Telegram ID.        |

**Response:**
- **200 OK**
```json
{
  "status": true
}
```

**Response Codes:**
- `200 OK` — Пользователь успешно обновлен.
- `400 Bad Request` — Неверный формат данных.
- `404 Not Found` — Пользователь не найден.
- `422 Unprocessable Entity` — Некорректное тело запроса.
- `500 Internal Server Error` — Ошибка на сервере.

---

##### POST /api/v1/user/refresh

**Description:** Обновление токена пользователя.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/user/refresh \
-H "Content-Type: application/json" \
-d '{"user_id": "uuid", "token": "old_token"}'
```

**Request Body:**

| Field     | Type   | Required | Description               |
|-----------|--------|----------|---------------------------|
| `user_id` | string | Yes      | UUID пользователя.        |
| `token`   | string | Yes      | Старый токен.             |

**Response:**
- **200 OK**
```json
{
  "token": "new_token"
}
```

**Response Codes:**
- `200 OK` — Токен успешно обновлен.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `422 Unprocessable Entity` — Некорректное тело запроса.
- `500 Internal Server Error` — Ошибка на сервере.

---

##### POST /api/v1/user/login

**Description:** Авторизация пользователя по API токену.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/user/login \
-H "Content-Type: application/json" \
-d '{"api_token": "string"}'
```

**Request Body:**

| Field       | Type   | Required | Description               |
|-------------|--------|----------|---------------------------|
| `api_token` | string | Yes      | API токен пользователя.   |

**Response:**
- **200 OK**
```json
{
  "id": "uuid",
  "telegram_id": 123456789
}
```

**Response Codes:**
- `200 OK` — Успешная авторизация.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `422 Unprocessable Entity` — Некорректное тело запроса.
- `500 Internal Server Error` — Ошибка на сервере.

---

#### For Staff

##### GET /api/v1/user/staff/get/:id

**Description:** Получение данных пользователя по ID.

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/user/staff/get/uuid \
-H "Authorization: Bearer token"
```

**Headers:**

| Header          | Required | Description               |
|-----------------|----------|---------------------------|
| `Authorization` | Yes      | Токен авторизации.        |

**Response:**
- **200 OK**
```json
{
  "user_id": "uuid",
  "role": "string",
  "tg_id": 123456789,
  "link_count": 10
}
```

**Response Codes:**
- `200 OK` — Данные пользователя успешно получены.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `404 Not Found` — Пользователь не найден.
- `422 Unprocessable Entity` — Некорректное тело запроса.
- `500 Internal Server Error` — Ошибка на сервере.

---

##### DELETE /api/v1/user/staff/delete/:id

**Description:** Удаление пользователя по ID (для сотрудников).

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/user/staff/delete/uuid \
-H "Authorization: Bearer token"
```

**Headers:**

| Header          | Required | Description        |
|-----------------|----------|--------------------|
| `Authorization` | Yes      | Токен авторизации. |

**Response Codes:**
- `204 No Content` — Пользователь успешно удален.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `404 Not Found` — Пользователь не найден.
- `422 Unprocessable Entity` — Некорректное тело запроса.
- `500 Internal Server Error` — Ошибка на сервере.

---

#### For Admins

##### POST /api/v1/user/admin/create

**Description:** Создание нового пользователя администратором.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/user/admin/create \
-H "Authorization: Bearer token" \
-H "Content-Type: application/json" \
-d '{"tg_id": 123456789, "is_staff": true, "is_admin": true}'
```

**Headers:**

| Header          | Required | Description               |
|-----------------|----------|---------------------------|
| `Authorization` | Yes      | Токен авторизации.        |

**Request Body:**

| Field      | Type   | Required | Description               |
|------------|--------|----------|---------------------------|
| `tg_id`    | int64  | Yes      | Telegram ID пользователя. |
| `is_staff` | bool   | Yes      | Флаг сотрудника.          |
| `is_admin` | bool   | Yes      | Флаг администратора.      |

**Response:**
- **201 Created**
```json
{
  "user_id": "uuid",
  "token": "string"
}
```

**Response Codes:**
- `201 Created` — Пользователь успешно создан.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `422 Unprocessable Entity` — Некорректное тело запроса.
- `500 Internal Server Error` — Ошибка на сервере.

---

##### PATCH /api/v1/user/admin/update/:id

**Description:** Обновление данных пользователя администратором.

**Request:**
```bash
curl -X PATCH http://localhost:8080/api/v1/user/admin/update/uuid \
-H "Authorization: Bearer token" \
-H "Content-Type: application/json" \
-d '{"tg_id": 123456789, "is_staff": true, "is_admin": true}'
```

**Headers:**

| Header          | Required | Description               |
|-----------------|----------|---------------------------|
| `Authorization` | Yes      | Токен авторизации.        |

**Request Body:**

| Field        | Type    | Required  | Description                |
|--------------|---------|-----------|----------------------------|
| `tg_id`      | int64   | No        | Новый Telegram ID.         |
| `is_staff`   | bool    | No        | Новый флаг сотрудника.     |
| `is_admin`   | bool    | No        | Новый флаг администратора. |

**Response:**
- **200 OK**
```json
{
  "status": true
}
```

**Response Codes:**
- `200 OK` — Пользователь успешно обновлен.
- `400 Bad Request` — Неверный формат данных.
- `401 Unauthorized` — Неверный токен.
- `404 Not Found` — Пользователь не найден.
- `422 Unprocessable Entity` — Некорректное тело запроса.
- `500 Internal Server Error` — Ошибка на сервере.

---
