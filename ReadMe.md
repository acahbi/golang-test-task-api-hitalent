# HiTalent
API сервис для создания и редактирования отделов и сотрудников на Go.

## Требования

- Go `1.23+`
- PostgreSQL
- Docker и Docker Compose

## Запуск через Docker Compose

```bash
docker compose up --build
```

После запуска сервис будет доступен по адресу `http://localhost:8080`.

## Swagger

Был добавлен Swagger для запуска curl через UI:

```text
http://localhost:8080/swagger/
```

OpenAPI JSON:

```text
http://localhost:8080/swagger/openapi.json
```
