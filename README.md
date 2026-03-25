[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/xR-tWBKa)

# Roomly — Meeting Room Booking Service

## Стек
- Go, chi, pgx
- PostgreSQL
- Docker Compose

## Запуск

### Требования
- Docker
- Docker Compose

### Быстрый старт
```bash
cp .env.example .env
make up
make seed
```

Сервис будет доступен на `http://localhost:8080`

## Команды

| Команда | Описание |
|---|---|
| `make up` | Собрать и запустить все сервисы |
| `make down` | Остановить сервисы |
| `make down-clean` | Остановить и удалить volumes |
| `make seed` | Заполнить БД тестовыми данными |
| `make test` | Запустить unit тесты |
| `make test-e2e` | Запустить E2E тесты (требует запущенного `make up`) |
| `make migrate-down` | Откатить последнюю миграцию |

## Аутентификация

Для тестирования используй `/dummyLogin`:
```bash
# admin токен
curl -X POST http://localhost:8080/dummyLogin \
  -H "Content-Type: application/json" \
  -d '{"role": "admin"}'

# user токен  
curl -X POST http://localhost:8080/dummyLogin \
  -H "Content-Type: application/json" \
  -d '{"role": "user"}'
```

Фиксированные UUID:
- admin: `00000000-0000-0000-0000-000000000001`
- user: `00000000-0000-0000-0000-000000000002`

## Архитектура
```
cmd/roomly/        — точка входа
internal/
  handler/         — HTTP handlers
  service/         — бизнес-логика
  repository/      — работа с БД
  domain/          — модели и интерфейсы
  middleware/       — JWT, роли, логирование
  scheduler/        — фоновый воркер генерации слотов
migrations/        — SQL миграции
scripts/           — seed скрипты
tests/e2e/         — E2E тесты
```

## Генерация слотов

Используется стратегия **rolling window** — фоновый воркер генерирует слоты на 7 дней вперёд при старте и затем каждый час. Слоты нарезаются по 30 минут из расписания комнаты. UUID слотов стабильны благодаря уникальному индексу `(room_id, starts_at)`.

## Тестовые данные (после `make seed`)

| Роль | Email | UUID |
|---|---|---|
| admin | admin@test.com | `00000000-0000-0000-0000-000000000001` |
| user | user@test.com | `00000000-0000-0000-0000-000000000002` |

Комнаты:
- **Room A** — Small meeting room, 4 места, пн-пт 09:00-18:00
- **Room B** — Large conference room, 10 мест, пн-сб 08:00-20:00
