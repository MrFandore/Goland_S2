# Практическое занятие №2
ФИО: Пряшников Дмитрий Максимович
Группа: ПИМО-01-25
## Описание решения

Проект состоит из двух микросервисов, взаимодействующих по **gRPC**:
- **Auth service** – gRPC сервер, реализующий метод `Verify` для проверки токена. В демо-режиме считает валидным только токен `"valid-token"`.
- **Tasks service** – HTTP API (REST) для управления задачами. Каждый защищённый запрос проходит через middleware, который вызывает gRPC метод `Verify` в Auth service с таймаутом 2 секунды.

Код генерируется из единого `.proto` контракта, что гарантирует строгую типизацию и согласованность данных между сервисами. При недоступности Auth или превышении таймаута Tasks возвращает клиенту соответствующий HTTP статус (503/504).

## Границы ответственности

- **Auth service**: только проверка токена. Не хранит состояние, не ведёт логов доступа (кроме отладочных). Возвращает gRPC статусы:
  - `codes.Unauthenticated` – невалидный или отсутствующий токен.
  - `codes.Internal` – внутренняя ошибка (не используется в демо).
- **Tasks service**:
  - Публичный эндпоинт `GET /health` для проверки доступности.
  - Защищённый эндпоинт `POST /tasks` для создания задачи (сохраняется в памяти).
  - В middleware извлекает токен из заголовка `Authorization`, вызывает Auth с таймаутом, при успехе сохраняет `subject` (идентификатор пользователя) в контексте запроса.
  - Маппинг gRPC ошибок в HTTP:
    - `Unauthenticated` → 401
    - `DeadlineExceeded` → 504
    - остальные (включая недоступность сервера) → 503

## Cхема взаимодействия

```
Клиент (curl/Postman)
        │
        │ POST /tasks (с токеном)
        ▼
┌─────────────────┐         gRPC (порт 50051)         ┌─────────────────┐
│  Tasks service  │ ────────────────────────────────▶ │  Auth service   │
│   (порт 8082)   │ ◀──────────────────────────────── │  (gRPC server)  │
└─────────────────┘    Verify(token) → (valid, sub)   └─────────────────┘
        │
        │ ответ 201/4xx/5xx
        ▼
      Клиент
```

**Последовательность обработки запроса**:
1. Tasks получает HTTP POST `/tasks` с заголовком `Authorization: Bearer <token>`.
2. Middleware создаёт контекст с таймаутом 2 секунды и вызывает gRPC метод `AuthService.Verify`.
3. Auth проверяет токен (в демо: `<token>` == `"valid-token"`).
4. Если токен верен, Auth возвращает `valid=true` и `subject=user-123`.
5. Tasks сохраняет задачу в памяти, привязывая её к `subject`, и возвращает клиенту `201 Created`.
6. При любой ошибке (невалидный токен, таймаут, недоступность Auth) Tasks возвращает соответствующий HTTP статус с сообщением об ошибке.

## Структура проекта

```
Prac2//
├── go.mod
├── pkg/
│   └── api/
│       └── auth/              
│           ├── auth.pb.go
│           └── auth_grpc.pb.go
│
├── proto/
│   └── auth.proto            
└── services/
├── auth/                  
│   ├── go.mod
│   ├── cmd/
│   │   └── auth/
│   │       └── main.go
│   └── internal/
│       └── server/
│           └── server.go 
└── tasks/                  
├── go.mod
├── cmd/
│   └── tasks/
│       └── main.go
└── internal/
├── client/
│   └── auth_client.go  
├── handlers/
│   └── tasks.go        
└── middleware/
└── auth.go         
```
## Описание файлов

| Файл/папка | Назначение |
|------------|------------|
| `go.mod` (корневой) | Корневой модуль проекта, используется для локальных замен (replace) |
| `pkg/api/auth/` | Автоматически сгенерированный код из proto-файла (структуры и gRPC интерфейсы) |
| `Prac2/proto/auth.proto` | Protobuf контракт сервиса аутентификации |
| `services/auth/cmd/auth/main.go` | Точка входа Auth service, запуск gRPC сервера |
| `services/auth/internal/server/server.go` | Бизнес-логика метода Verify |
| `services/tasks/cmd/tasks/main.go` | Точка входа Tasks service, подключение gRPC клиента |
| `services/tasks/internal/client/auth_client.go` | Обёртка для gRPC вызовов к Auth с таймаутом |
| `services/tasks/internal/middleware/auth.go` | Middleware для проверки JWT через gRPC |
| `services/tasks/internal/handlers/tasks.go` | HTTP обработчик создания задачи |

## Запуск проекта

### 1. Сгенерируйте код из proto
Откройте терминал в корневой папке `Go_S2` и выполните:
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative Prac2/proto/auth.proto
```

### 2. Запустите Auth сервис (gRPC сервер)
```bash
cd Prac2/services/auth
export AUTH_GRPC_PORT=50051   # Linux/macOS
# или для PowerShell (Windows):
# $env:AUTH_GRPC_PORT="50051"
go mod tidy
go run ./cmd/auth
```

### 3. Запустите Tasks сервис (HTTP API)
В отдельном терминале:
```bash
cd Prac2/services/tasks
export TASKS_PORT=8082
export AUTH_GRPC_ADDR=localhost:50051
# для PowerShell:
# $env:TASKS_PORT="8082"
# $env:AUTH_GRPC_ADDR="localhost:50051"
go mod tidy
go run ./cmd/tasks
```

После запуска сервер Tasks доступен по адресу: http://localhost:8082

## Доступные эндпоинты (Tasks Service)

### `POST /tasks`
Создание задачи (требуется валидный токен).

**Заголовок:** `Authorization: Bearer <token>`  
**Тело запроса:** `{"title":"Название задачи"}`

**Успешный ответ (201 Created):**
```json
{
  "id": "1",
  "title": "Название задачи",
  "subject": "user-123"
}
```

**Ошибки:**
- `401 Unauthorized` – отсутствует или неверный токен
- `503 Service Unavailable` – Auth сервис недоступен
- `504 Gateway Timeout` – превышено время ожидания ответа от Auth

### `GET /health`
Проверка работоспособности (не требует токена).  
Ответ: `OK`

## Примеры тестирования

### Тестовые данные
- **Валидный токен:** `valid-token` (зашит в сервисе Auth для демо)
- **Невалидный токен:** любой другой (например, `wrong-token`)

### Через командную строку (curl)

**Успешный запрос:**
```bash
curl -X POST http://localhost:8082/tasks \
  -H "Authorization: Bearer valid-token" \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn gRPC"}'
```

**Невалидный токен (401):**
```bash
curl -X POST http://localhost:8082/tasks \
  -H "Authorization: Bearer wrong-token" \
  -H "Content-Type: application/json" \
  -d '{"title":"Hack"}'
```

**Auth сервис недоступен (503/504):**  
Остановите Auth сервис (Ctrl+C) и повторите первый запрос.

**Проверка здоровья:**
```bash
curl http://localhost:8082/health
```

### Через PowerShell (Windows)

```
powershell
# Успешный запрос
$body = @{title="Learn gRPC"} | ConvertTo-Json
Invoke-WebRequest -Uri http://localhost:8082/tasks -Method POST `
  -Headers @{Authorization="Bearer valid-token"} `
  -ContentType "application/json" -Body $body

# Невалидный токен
Invoke-WebRequest -Uri http://localhost:8082/tasks -Method POST `
  -Headers @{Authorization="Bearer wrong-token"} `
  -ContentType "application/json" -Body $body
```

## Особенности проекта

- **gRPC взаимодействие** между сервисами Auth и Tasks
- **Deadline/таймаут** при вызове gRPC (2 секунды) для предотвращения зависаний
- **Обработка ошибок** и преобразование gRPC статусов в HTTP коды:
  - `codes.Unauthenticated` → 401
  - `codes.DeadlineExceeded` → 504
  - остальные ошибки → 503
- **Stateless аутентификация** через токен (проверка в Auth)
- **Локальная разработка** с использованием replace в go.mod для общих пакетов

## Требования

- Установленный Go версии 1.21 или выше
- Protocol Buffers компилятор (`protoc`) и плагины `protoc-gen-go`, `protoc-gen-go-grpc`
- (Опционально) Git
