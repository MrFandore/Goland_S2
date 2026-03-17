# Tech IP Sem2 – Практическое занятие №1

## Описание решения

Проект состоит из двух микросервисов:
- **Auth service** – отвечает за выдачу и проверку токенов (учебная реализация).
- **Tasks service** – CRUD для задач, перед каждой операцией проверяет токен через Auth.

Взаимодействие синхронное по HTTP. Во все запросы прокидывается `X-Request-ID` для сквозной трассировки, установлены таймауты (3 секунды) на вызов Auth.

## Границы ответственности

- **Auth**: только аутентификация/авторизация (в данном упрощении – проверка фиксированного токена).
- **Tasks**: управление задачами, хранение в памяти, проверка доступа делегируется Auth.

## Схема взаимодействия

```mermaid
sequenceDiagram
    participant Client
    participant Tasks
    participant Auth

    Client->>Tasks: Запрос с Authorization
    Tasks->>Auth: GET /v1/auth/verify (таймаут 3с, X-Request-ID)
    Auth-->>Tasks: 200 OK (валидный токен) / 401
    Tasks-->>Client: 200/201/204 или 401/503
```
# Запуск
## Предварительные требования
### Go 1.18+

Установите зависимости в корне проекта:
```markdown
go mod tidy
```
## Запуск Auth service
```bash
cd services/auth
export AUTH_PORT=8081
go run ./cmd/auth
```
## Запуск Tasks service
```bash
cd services/tasks
export TASKS_PORT=8082
export AUTH_BASE_URL=http://localhost:8081
go run ./cmd/tasks
```

## Переменные окружения
```markdown
Сервис	Переменная	Значение по умолчанию	Описание
Auth	AUTH_PORT	8081	Порт, на котором слушает Auth
Tasks	TASKS_PORT	8082	Порт Tasks
Tasks	AUTH_BASE_URL	http://localhost:8081	Базовый URL для вызова Auth
```
## Примеры запросов
### Получить токен
```bash
curl -s -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: req-001" \
  -d '{"username":"student","password":"student"}'
```
### Создать задачу (с токеном)
```bash
curl -i -X POST http://localhost:8082/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-002" \
  -d '{"title":"Выполнить ПЗ","description":"разделение монолита","due_date":"2026-01-15"}'
```
### Попытка без токена (ожидается 401)
```bash
curl -i http://localhost:8082/v1/tasks -H "X-Request-ID: req-003"
```
## Полная документация API – в файле docs/pz17_api.md.

# Скриншот 1: Запуск Auth service
## Файл: picture/skrin_1.png

Что видно на скриншоте:

Терминал PowerShell с выполненной командой go run ./cmd/auth.

Лог сервиса: Auth service starting on :8081.

![img.png](Фотокарточки/1.png)
# Скриншот 2: Запуск Tasks service
## Файл: picture/skrin_2.png
Что видно на скриншоте:

```powershell
$env:TASKS_PORT = 8082
$env:AUTH_BASE_URL = "http://localhost:8081"
go run ./cmd/tasks
```

Лог сервиса: Tasks service starting on :8082.

Подтверждение, что Tasks ожидает запросы.

![img.png](Фотокарточки/2.png)
# Скриншот 3: Получение токена через Auth service
## Файл: picture/skrin_3.png
Что видно на скриншоте:
```json
Команда invoke-restmethod для отправки POST-запроса на 
/v1/auth/login с телом {"username":"student","password":"student"} и 
заголовком X-Request-ID: req-login-001.
        
Ответ сервера: JSON-объект с полями access_token и 
token_type (значение demo-token и Bearer).
```

Подтверждение, что токен успешно получен.

![img.png](Фотокарточки/3.png)
# Скриншот 4: Запрос к Tasks без токена (401)
## Файл: picture/skrin_5.png
Что видно на скриншоте:

Команда invoke-restmethod без заголовка Authorization.

Ошибка: missing authorization header и код 401.

В логах Tasks service – запись с req-no-token и статусом 401.(Общий скрин 4 и 5 ошибки ниже)
![img.png](Фотокарточки/4.png)
## Что видно на скриншоте:

Команда invoke-restmethod с заголовком Authorization: Bearer wrong-token.

Ошибка: missing authorization header или unauthorized с кодом 401.

В логах Tasks service – запись с req-wrong-token и статусом 401.
![img.png](Фотокарточки/5.png)
## Скриншот общих ошибок в Tasks
![img_1.png](Фотокарточки/4_5.png)
# Скриншот 6: Auth service недоступен (503)
Что видно на скриншоте:

Auth service остановлен (Ctrl+C в окне 1).

Команда invoke-restmethod (или curl.exe) с правильным токеном, но Auth не отвечает.

Ошибка: The remote server returned an error: (503) Service Unavailable.

В логах Tasks service – запись с req-auth-down и статусом 503 (или 502, в зависимости от реализации).

Подтверждение, что Tasks корректно обрабатывает недоступность Auth.
![img.png](Фотокарточки/img.png)
# Скриншот 7: Получение списка задач после восстановления Auth
Что видно на скриншоте:

Auth service запущен заново.

Команда invoke-restmethod (или curl.exe) с заголовком Authorization: Bearer demo-token и X-Request-ID: req-list.

Ответ: массив задач в формате JSON.

В логах Auth и Tasks – соответствующие записи с request-id req-list.

![img.png](Фотокарточки/7.png)
