```markdown
# Практика 2 - gRPC микросервисы с аутентификацией

## Структура проекта
```
Go_S2/
├── go.mod
├── pkg/
│   └── api/
│       └── auth/                # сгенерированный код из proto
│           ├── auth.pb.go
│           └── auth_grpc.pb.go
└── Prac2/
├── proto/
│   └── auth.proto            # контракт gRPC сервиса
└── services/
├── auth/                  # gRPC сервер (Auth Service)
│   ├── go.mod
│   ├── cmd/
│   │   └── auth/
│   │       └── main.go
│   └── internal/
│       └── server/
│           └── server.go  # реализация Verify метода
└── tasks/                  # HTTP API (Tasks Service)
├── go.mod
├── cmd/
│   └── tasks/
│       └── main.go
└── internal/
├── client/
│   └── auth_client.go  # gRPC клиент к Auth
├── handlers/
│   └── tasks.go        # HTTP обработчики
└── middleware/
└── auth.go          # middleware проверки токена
```

## Описание файлов
- **go.mod** (корневой) - корневой модуль проекта, используется для локальных замен (replace)
- **pkg/api/auth/** - автоматически сгенерированный код из proto-файла (структуры и gRPC интерфейсы)
- **Prac2/proto/auth.proto** - Protobuf контракт сервиса аутентификации
- **services/auth/...** - реализация gRPC сервера Auth:
  - `cmd/auth/main.go` - точка входа, запуск gRPC сервера
  - `internal/server/server.go` - бизнес-логика метода Verify (проверка токена)
- **services/tasks/...** - HTTP API сервиса задач:
  - `cmd/tasks/main.go` - запуск HTTP сервера, подключение gRPC клиента
  - `internal/client/auth_client.go` - обёртка для gRPC вызовов к Auth с таймаутом
  - `internal/middleware/auth.go` - middleware, проверяющий JWT через gRPC
  - `internal/handlers/tasks.go` - обработчик создания задачи (пример защищённого ресурса)

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

```
POST /tasks
```
Создание задачи (требуется валидный токен в заголовке Authorization).

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
- `401 Unauthorized` - отсутствует или неверный токен
- `503 Service Unavailable` - Auth сервис недоступен
- `504 Gateway Timeout` - превышено время ожидания ответа от Auth

```
GET /health
```
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

```powershell
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

## Решение проблем

**Ошибка генерации proto:**
- Проверьте, что `protoc` доступен в PATH
- Убедитесь, что установлены плагины: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest` и `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
- Проверьте путь к proto-файлу: он должен быть указан относительно текущей директории

**Ошибка подключения к Auth:**
- Убедитесь, что Auth сервис запущен и слушает порт 50051
- Проверьте переменную `AUTH_GRPC_ADDR` в Tasks: `localhost:50051` (без протокола)

**Зависание запросов при недоступном Auth:**
- В коде клиента установлен таймаут 2 секунды, после которого возвращается ошибка 504
- Если этого не происходит, проверьте, что контекст с таймаутом действительно используется

**Проблемы с go mod:**
- В каждом сервисе выполните `go mod tidy` для загрузки зависимостей
- Убедитесь, что replace в go.mod указывает на правильный корневой модуль (например, `github.com/yourusername/Go_S2 => ../../../`)
```