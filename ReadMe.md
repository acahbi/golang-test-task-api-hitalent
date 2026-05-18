# HiTalent
API сервис для создания и редактирования подразделений и сотрудников на Go.

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

## Основные API:

Проверить сервис на health-метрику:

```text
curl -X GET http://localhost:8080/health
```

Создать подразделении (запись в department):

```text
curl -X POST http://localhost:8080/departments/ -H "Content-Type: application/json" -d "{\"name\": \"TestDepartmentCMD\", \"parent_id\": null}"
```

Создать сотрудника в подразделении (запись в employee):

```text
curl -X POST http://localhost:8080/departments/1/employees/ -H "Content-Type: application/json" -d "{\"full_name\": \"TestEmployeeCMD\", \"position\": \"manager\", \"hired_at\": \"2026-03-24T00:00:00Z\"}"
```

Получить подразделение (детали + сотрудники + поддерево):
```text
curl -X GET "http://localhost:8080/departments/1?depth=2&include_employees=false"
```

Обновить имя подразделения или переместить подразделение в другое (изменить parent):

```text
curl -X PATCH http://localhost:8080/departments/3 -H "Content-Type: application/json" -d "{\"name\": \"TestDepartmentMovedCMD\", \"parent_id\": 1}"
```

Удалить подразделение:

```text
curl -X DELETE "http://localhost:8080/departments/3?mode=cascade"
```