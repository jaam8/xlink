# Gateway routing documentation

## 1. Routes themselves

### User service

#### Common part

- /api/v1/user/create POST
- /api/v1/user/:id PATCH
- /api/v1/user/token/refresh POST

#### For staff

- /api/v1/user/staff/:id GET
- /api/v1/user/staff/get/by-tg-id/:tgId GET
- /api/v1/user/staff/:id DELETE
- /api/v1/user/staff/role/:id GET

#### For admins

- /api/v1/user/admin/create POST
- /api/v1/user/admin/update/:id PATCH
- /api/v1/user/admin/delete/:id DELETE
- /api/v1/user/admin/get/by-token POST <- POST потому что иначе токен передавался бы в url а это уязвимость
- /api/v1/user/admin/token/check POST <- POST по той же причине

### Shortener

#### Common part

- /api/v1/s/:shortLink GET
- /api/v1/s/crud/ POST <- create

#### For link owner

- /api/v1/s/crud/owner/:id PUT
- /api/v1/s/crud/owner/:id DELETE

#### For admins

- /api/v1/s/crud/admin/:id PUT
- /api/v1/s/crud/admin/:id DELETE

### Analytics (auth'd users only)

- /api/v1/analytics/by-country GET
- /api/v1/analytics/by-region GET
- /api/v1/analytics/by-browser GET
- /api/v1/analytics/by-os GET
- /api/v1/analytics/by-device-type GET
- /api/v1/analytics/by-hour GET
- /api/v1/analytics/by-date GET
- /api/v1/analytics/by-referrer GET

## 2. Explaination

### User service

#### Common part

- /api/v1/user/create POST

| Input               | Output                                 | Summary                                            |
|---------------------|----------------------------------------|----------------------------------------------------|
| `{"tg_id": *int64}` | `{"user_id": string, "token": string}` | Create user with `is_staff` & `is_admin` set to **false** |

- /api/v1/user/:id PATCH

| Input              | Output             | Summary                                                   |
|--------------------|--------------------|-----------------------------------------------------------|
| `{"tg_id": int64}` | `{"status": bool}` | Update user with `is_staff` & `is_admin` set to **false** |

- /api/v1/user/token/refresh POST

| Input                                  | Output                                 | Summary                                                                            |
|----------------------------------------|----------------------------------------|------------------------------------------------------------------------------------|
| `{"token": string, "user_id": string}` | `{"token": string}` | Refresh token with using `user_id` and `token` given in **body** (not auth header) |

#### For staff

- /api/v1/user/staff/:id GET

| Input                      | Output                                                                       | Summary                                    |
|----------------------------|------------------------------------------------------------------------------|--------------------------------------------|
| **id** - string (in query) | `{"user_id": string, "role": string, "tg_id": *string, "link_count": int32}` | Get user by **user_id** given **in query** |

- /api/v1/user/staff/get/by-tg-id/:tgId GET

| Input                        | Output                                                                       | Summary                                 |
|------------------------------|------------------------------------------------------------------------------|-----------------------------------------|
| **tg_id** - int64 (in query) | `{"user_id": string, "status": bool}` | Get user by **tg_id** given **in query** |

- /api/v1/user/staff/:id DELETE

| Input                      | Output                                                                      | Summary                                  |
|----------------------------|-----------------------------------------------------------------------------|------------------------------------------|
| **id** - string (in query) | `{"status": bool}` | Delete user by **id** given **in query** |

- /api/v1/user/staff/role/:id GET

| Input                      | Output                                                 | Summary                                           |
|----------------------------|--------------------------------------------------------|---------------------------------------------------|
| **id** - string (in query) | `{"role": string, "is_admin": bool, "is_staff": bool}` | Get user roles data by **id** given **in query**  |

#### For admins

- /api/v1/user/admin/create POST

| Input                                                   | Output                                 | Summary         |
|---------------------------------------------------------|----------------------------------------|-----------------|
| `{"tg_id": *int64, "is_admin": bool, "is_staff": bool}` | `{"user_id": string, "token": string}` | Create new user |

- /api/v1/user/admin/update/:id PATCH

| Input                                                   | Output                                 | Summary         |
|---------------------------------------------------------|----------------------------------------|-----------------|
| `{"tg_id": *int64, "is_admin": bool, "is_staff": bool}` | `{"user_id": string, "token": string}` | Create new user |

- /api/v1/user/admin/delete/:id DELETE

| Input                      | Output             | Summary                                  |
|----------------------------|--------------------|------------------------------------------|
| **id** - string (in query) | `{"status": bool}` | Delete user by **id** given **in query** |

- /api/v1/user/admin/get/by-token POST <- POST потому что иначе токен передавался бы в url а это уязвимость

| Input               | Output                                | Summary                                                                                                    |
|---------------------|---------------------------------------|------------------------------------------------------------------------------------------------------------|
| `{"token": string}` | `{"user_id": string, "status": bool}` | Get user by **token** given in **body**. Use of **body** and not **query** is for "at least better" safety |

- /api/v1/user/admin/token/check POST <- POST по той же причине

| Input                                 | Output             | Summary                                                                |
|---------------------------------------|--------------------|------------------------------------------------------------------------|
| `{"user_id": string, token": string}` | `{"status": bool}` | Check if user with given **user_id** has given **token** (in **body**) |

### Shortener

#### Common part

- /api/v1/s/:shortLink GET

| Input                                                                                                                                                       | Output                 | Summary                                                                              |
|-------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------|--------------------------------------------------------------------------------------|
| **shortLink** - string (in **path**) <br> _referrer_ - string (in **header**: "HTTP_REFERER") <br> _visitorToken_ - string (in **cookie**: "xlinkVisitor") | Redirect to target url | Send given **short link** and request info, then redirect to target url in response. |

- /api/v1/s/crud/ POST <- create

| Input                                                              | Output                                                                                                                          | Summary                                                                                                                                         |
|--------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| `{"short_link": *string, "target_url": string"}` | `{"link_id": string, "user_id": string, "short_link": string, "target_url": string, "created_at": string, "expire_at": string}` | Simply create link with given data <br> **short_link** is generated if not specified <br> **user_id** is got from auth token (_GetUserIDByToken_) |

#### For link owner

- /api/v1/s/crud/owner/:id PUT

| Input                                                                                                                              | Output                                                                                                                          | Summary                                                                                                                                             |
|------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| **link_id** - string (in **path**) <br> `{"regenerate": bool, "short_link": *string, "target_url": string, "expire_at": *string}` | `{"link_id": string, "user_id": string, "short_link": string, "target_url": string, "created_at": string, "expire_at": string}` | Update link with given data <br> **short_link** is generated if **generated** == true <br>  **user_id** is got from auth token (_GetUserIDByToken_) |

- /api/v1/s/crud/owner/:id DELETE

| Input                               | Output             | Summary                                                                                                        |
|-------------------------------------|--------------------|----------------------------------------------------------------------------------------------------------------|
| **link_id** - string (in **path**) | `{"status": bool}` | Delete link with given **link_id** (in **path**) <br> **user_id** is got from auth token (_GetUserIDByToken_) |

#### For admins

- /api/v1/s/crud/admin/:id PUT

| Input                                                                                                                              | Output                                                                                                                          | Summary                                                                                                                                             |
|------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| **link_id** - string (in **path**) <br> `{"regenerate": bool, "short_link": *string, "target_url": string, "expire_at": *string}` | `{"link_id": string, "user_id": string, "short_link": string, "target_url": string, "created_at": string, "expire_at": string}` | Update link with given data <br> **short_link** is generated if **generated** == true <br>  **user_id** is got from auth token (_GetUserIDByToken_) |

- /api/v1/s/crud/admin/:id DELETE

| Input                               | Output             | Summary                                                                                                        |
|-------------------------------------|--------------------|----------------------------------------------------------------------------------------------------------------|
| **link_id** - string (in **path**) | `{"status": bool}` | Delete link with given **link_id** (in **path**) <br> **user_id** is got from auth token (_GetUserIDByToken_) |

### Analytics (auth'd users only)

- /api/v1/analytics/by-country GET

| Input                                                                                                                                                            | Output                                                                                                              | Summary                                            |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------|----------------------------------------------------|
| **link_owner** - string (**calculated** from token) <br> **short_link** - string (in **query**) <br> **start_date** - YYYY-MM-DD, e.g. 2025-04-23 (in **query**) | `{"data": [{"date": string, "stats": [{"country": string, "clicks": uint64, "unique_clicks": uint64}, ...]}, ...]}` | Get link stats, aggregated by clickers' countries  |

- /api/v1/analytics/by-region GET

| Input                                                                                                                                                            | Output                                                                                                             | Summary                                         |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------|-------------------------------------------------|
| **link_owner** - string (**calculated** from token) <br> **short_link** - string (in **query**) <br> **start_date** - YYYY-MM-DD, e.g. 2025-04-23 (in **query**) | `{"data": [{"date": string, "stats": [{"region": string, "clicks": uint64, "unique_clicks": uint64}, ...]}, ...]}` | Get link stats, aggregated by clickers' regions |

- /api/v1/analytics/by-browser GET

| Input                                                                                                                                                            | Output                                                                                                              | Summary                                             |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------|
| **link_owner** - string (**calculated** from token) <br> **short_link** - string (in **query**) <br> **start_date** - YYYY-MM-DD, e.g. 2025-04-23 (in **query**) | `{"data": [{"date": string, "stats": [{"browser": string, "clicks": uint64, "unique_clicks": uint64}, ...]}, ...]}` | Get link stats, aggregated by clickers' countries   |

- /api/v1/analytics/by-os GET

| Input                                                                                                                                                            | Output                                                                                                         | Summary                                    |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------|--------------------------------------------|
| **link_owner** - string (**calculated** from token) <br> **short_link** - string (in **query**) <br> **start_date** - YYYY-MM-DD, e.g. 2025-04-23 (in **query**) | `{"data": [{"date": string, "stats": [{"os": string, "clicks": uint64, "unique_clicks": uint64}, ...]}, ...]}` | Get link stats, aggregated by clickers' OS |

- /api/v1/analytics/by-device-type GET

| Input                                                                                                                                                            | Output                                                                                                                  | Summary                                             |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------|
| **link_owner** - string (**calculated** from token) <br> **short_link** - string (in **query**) <br> **start_date** - YYYY-MM-DD, e.g. 2025-04-23 (in **query**) | `{"data": [{"date": string, "stats": [{"device_type": string, "clicks": uint64, "unique_clicks": uint64}, ...]}, ...]}` | Get link stats, aggregated by clickers' device type |

- /api/v1/analytics/by-hour GET

| Input                                                                                                                                                            | Output                                                                                                           | Summary                                         |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------|-------------------------------------------------|
| **link_owner** - string (**calculated** from token) <br> **short_link** - string (in **query**) <br> **start_date** - YYYY-MM-DD, e.g. 2025-04-23 (in **query**) | `{"data": [{"date": string, "stats": [{"hour": uint32, "clicks": uint64, "unique_clicks": uint64}, ...]}, ...]}` | Get link stats, aggregated by click hour (0-23) |

- /api/v1/analytics/by-date GET

| Input                                                                                                                                                            | Output                                                                                                           | Summary                                               |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------|
| **link_owner** - string (**calculated** from token) <br> **short_link** - string (in **query**) <br> **start_date** - YYYY-MM-DD, e.g. 2025-04-23 (in **query**) | `{"data": [{"date": string, "stats": [{"date": string, "clicks": uint64, "unique_clicks": uint64}, ...]}, ...]}` | Get link stats, aggregated by click date (YYYY-MM-DD) |

- /api/v1/analytics/by-referrer GET

| Input                                                                                                                                                            | Output                                                                                                               | Summary                                               |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------|
| **link_owner** - string (**calculated** from token) <br> **short_link** - string (in **query**) <br> **start_date** - YYYY-MM-DD, e.g. 2025-04-23 (in **query**) | `{"data": [{"date": string, "stats": [{"referrer": string, "clicks": uint64, "unique_clicks": uint64}, ...]}, ...]}` | Get link stats, aggregated by clickers' HTTP_REFERRER |
