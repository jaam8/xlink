# Gateway routing documentation

## 1. Routes themselves

### User service

#### Common part

- /api/v1/user/create POST
- /api/v1/user/:id PATCH
- /api/v1/user/token/refresh POST

#### For staff

- /api/v1/user/staff/:id GET
- /api/v1/user/staff/get/by-tg-id GET
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

| Input | Output | Summary |
|-------|--------|---------|
|       |        |         |



- /api/v1/user/:id PATCH
- /api/v1/user/token/refresh POST

#### For staff

- /api/v1/user/staff/:id GET
- /api/v1/user/staff/get/by-tg-id GET
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