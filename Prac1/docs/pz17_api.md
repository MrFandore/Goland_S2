
# API Documentation

## Auth Service (порт 8081)

### `POST /v1/auth/login`
Получение токена.

**Request body:**
```json
{
  "username": "student",
  "password": "student"
}
```
Response 200:
```json
{
"access_token": "demo-token",
"token_type": "Bearer"
}
```
Errors: 400 (неверный формат), 401 (неверные учётные данные).

GET /v1/auth/verify

Проверка токена.

Headers:
```markdown
Authorization: Bearer <token>

X-Request-ID (опционально)
```
Response 200:
```json
{
  "valid": true,
  "subject": "student"
}
```
Response 401:
```json
{
  "valid": false,
  "error": "unauthorized"
}
```
Tasks Service (порт 8082)

Все запросы обязательно должны содержать заголовок:

Authorization: Bearer <token> (токен, полученный от Auth).

POST /v1/tasks

Создать задачу.

Request body:
```json
{
  "title": "Read lecture",
  "description": "Prepare notes",
  "due_date": "2026-01-10"
}
```
Response 201: созданная задача (см. модель Task).
```markdown
GET /v1/tasks
```
Получить список всех задач.

Response 200: массив задач.
```markdown
GET /v1/tasks/{id}
```
Получить задачу по идентификатору.

Response 200: задача.
404: задача не найдена.
```markdown
PATCH /v1/tasks/{id}
```
Частичное обновление задачи.
```markdown
Request body: любые поля из Task (кроме id, created_at, updated_at).
```
Response 200: обновлённая задача.

404: задача не найдена.
```markdown
DELETE /v1/tasks/{id}
```

Удалить задачу.

Response 204: успешно (без тела).

Ошибки (общие для всех эндпоинтов Tasks)
```markdown
400 Bad Request – неверные данные в запросе.

401 Unauthorized – отсутствует или невалиден токен.

404 Not Found – задача не найдена.

503 Service Unavailable – сервис авторизации недоступен.
```
