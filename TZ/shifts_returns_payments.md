# SmartERP — ТЗ: Смены, Возвраты, Способы оплаты

> Зависимость: реализуется поверх текущей архитектуры (`arch.md`).
> Порядок: Миграция → Payments → Returns → Shifts → Wiring.

---

## 0. Контекст

В системе реализованы документы (RECEIPT, SALE, WRITEOFF, INVENTORY, TRANSFER) и леджер с формулами AVCost.
В леджере уже зарезервированы reason `RETURN_IN` / `RETURN_OUT`, но ни один документ их не использует.
Нет понятий: кассовая смена, кассир на смене, способ оплаты, возврат товара.

**Цель**: добавить полноценный цикл работы кассира — открытие смены, продажи/возвраты с фиксацией оплаты, внесение/изъятие наличных, закрытие смены с отчётом.

---

## 1. Миграция БД

**Файл**: `backend/migrations/000002_shifts_returns_payments.sql`

### 1.1 Расширение documents.documents

- [ ] Добавить `'RETURN'` в CHECK constraint на `type`
- [ ] Добавить колонку `shift_id UUID REFERENCES documents.shifts(id)` (nullable)

### 1.2 Таблица `documents.shifts`

```sql
shifts (
  id            UUID PRIMARY KEY,
  tenant_id     UUID NOT NULL,
  user_id       UUID NOT NULL REFERENCES platform.tenant_users(id),
  warehouse_id  UUID NOT NULL REFERENCES catalog.warehouses(id),
  opened_at     TIMESTAMP NOT NULL DEFAULT now(),
  closed_at     TIMESTAMP,
  opening_cash  DECIMAL(14,4) NOT NULL DEFAULT 0,
  closing_cash  DECIMAL(14,4),
  status        VARCHAR NOT NULL DEFAULT 'open'
                CHECK (status IN ('open','closed')),
  created_at    TIMESTAMP NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP NOT NULL DEFAULT now()
)
```

- [ ] Partial unique index: `UNIQUE (tenant_id, user_id) WHERE status = 'open'` — гарантирует максимум 1 открытую смену на кассира
- [ ] Index: `(tenant_id, user_id, status)` — быстрый поиск текущей смены

### 1.3 Таблица `documents.shift_cash_ops`

```sql
shift_cash_ops (
  id          UUID PRIMARY KEY,
  shift_id    UUID NOT NULL REFERENCES documents.shifts(id),
  type        VARCHAR NOT NULL CHECK (type IN ('cash_in','cash_out')),
  amount      DECIMAL(14,4) NOT NULL CHECK (amount > 0),
  note        TEXT,
  created_at  TIMESTAMP NOT NULL DEFAULT now()
)
```

- [ ] Index: `(shift_id)`

### 1.4 Таблица `documents.document_payments`

```sql
document_payments (
  id            UUID PRIMARY KEY,
  document_id   UUID NOT NULL REFERENCES documents.documents(id) ON DELETE CASCADE,
  method        VARCHAR NOT NULL CHECK (method IN ('cash','card','transfer')),
  amount        DECIMAL(14,4) NOT NULL CHECK (amount > 0)
)
```

- [ ] Index: `(document_id)`
- [ ] Один документ может иметь несколько строк оплаты (смешанная оплата: часть наличными, часть картой)

**Done when:** миграция применена, все таблицы и индексы созданы, constraints работают.

---

## 2. Способы оплаты (Payments)

> Payments — расширение существующей фичи `documents/`, не отдельная фича.

### 2.1 Модели

Добавить в `documents/models.go`:

- [ ] Структура `PaymentInput` — `method` (cash/card/transfer), `amount`
- [ ] Структура `Payment` — `id`, `method`, `amount`
- [ ] Поле `Payments []PaymentInput` в `CreateRequest`
- [ ] Поле `Payments []Payment` в `Document`

### 2.2 Репозиторий

Новый файл `documents/repository_payments.go` (~60 строк):

- [ ] `ReplacePayments(ctx, tx, documentID, []PaymentInput)` — DELETE + INSERT (как ReplaceItems)
- [ ] `LoadPayments(ctx, documentID)` — SELECT для загрузки в ByID

### 2.3 Интеграция в service_core.go

- [ ] При создании документа (createDraftTx): после `ReplaceItems` → вызвать `ReplacePayments`
- [ ] При обновлении документа (updateDraftTx): аналогично
- [ ] При чтении документа (ByID): загрузить payments

### 2.4 Валидация

- [ ] `method` — только `cash`, `card`, `transfer`
- [ ] `amount > 0`
- [ ] Сумма payments должна совпадать с суммой документа (items total_amount)
- [ ] Payments обязательны для SALE и RETURN, опциональны для остальных типов

**Done when:** создание/редактирование/чтение документов с payments работает, валидация сумм проходит.

---

## 3. Возвраты (Returns)

> Возврат = свободный (без привязки к конкретной продаже).
> Кассир сам указывает товар, количество и цену возврата.

### 3.1 Новый тип документа `RETURN`

Возврат использует существующую механику документов:
- Создаётся как `draft` → проводится → создаёт записи в леджере
- `warehouse_id` — склад, на который возвращается товар
- `shift_id` — опциональная привязка к смене

### 3.2 Леджер-логика

Модификации в `documents/service_posting_helpers.go`:

- [ ] В `docMeta()`: добавить `case "RETURN": return "IN", "RETURN_IN"`
- [ ] В `revenueForType()`: для `RETURN` возвращать `amount.Neg()` (отрицательная выручка = убыток/возврат денег)

При проведении RETURN:
- Тип записи: `IN` (товар возвращается на склад)
- Reason: `RETURN_IN`
- Revenue: отрицательная (минус от суммы возврата)
- AVCost пересчитывается как при обычном приходе (applyIn)
- Profit будет отрицательным (отражает убыток от возврата)

### 3.3 Composite возвраты

Модификация в `documents/service_posting.go`:

- [ ] В `buildEntriesForComposite()`: добавить ветку `if doc.Type == "RETURN"` → `buildCompositeReturnEntries()`
- [ ] `buildCompositeReturnEntries()` — зеркало `buildCompositeSaleEntries()`, но:
  - type = `IN`, reason = `RETURN_IN`
  - revenue = отрицательная (пропорционально распределена между компонентами)

### 3.4 Пример сценария

```
Товар "Кофе" продан за 150₸, себестоимость 80₸, прибыль = 70₸.
Возврат: кассир создаёт RETURN, товар "Кофе", qty=1, unit_price=150₸.
Леджер: type=IN, reason=RETURN_IN, qty=1, unit_price=150₸, revenue=-150₸.
AVCost пересчитан. Profit = -150₸ - cogs (отрицательный).
```

**Done when:** RETURN документ проводится, леджер корректно считает отрицательную прибыль, AVCost пересчитывается, composite возвраты работают.

---

## 4. Кассовые смены (Shifts)

> Новая фича: `backend/internal/features/shifts/`

### 4.1 Бизнес-правила

- [ ] Кассир открывает смену перед началом работы
- [ ] Максимум 1 открытая смена на кассира (enforced DB index + service check)
- [ ] При открытии указывается склад и начальная сумма наличных
- [ ] Все SALE и RETURN документы в рамках смены привязываются через `shift_id`
- [ ] Кассир может вносить/изъимать наличные (cash_in/cash_out)
- [ ] При закрытии система рассчитывает ожидаемый остаток наличных

### 4.2 API эндпоинты

| Метод | Путь | Тело/Параметры | Описание |
|-------|------|----------------|----------|
| POST | `/api/client/shifts/open` | `{ warehouse_id, opening_cash }` | Открыть смену |
| GET | `/api/client/shifts/current` | — | Текущая открытая смена кассира |
| POST | `/api/client/shifts/cash-op` | `{ type, amount, note }` | Внесение/изъятие наличных |
| POST | `/api/client/shifts/close` | — | Закрыть смену, получить отчёт |
| GET | `/api/client/shifts/{id}/report` | — | Отчёт по смене (любой) |

Все эндпоинты требуют auth scope `client`. UserID берётся из JWT claims.

### 4.3 Структура файлов

| Файл | Строк | Содержание |
|------|-------|------------|
| `models.go` | ~60 | OpenRequest, CashOpRequest, Shift, CashOp, ShiftReport |
| `errors.go` | ~10 | ErrShiftAlreadyOpen, ErrNoOpenShift, ErrInvalidCashOpType |
| `repository.go` | ~120 | CRUD + агрегации для отчёта |
| `service.go` | ~90 | Open, Current, CashOp, Close, Report |
| `handler.go` | ~100 | 5 HTTP handlers |
| `routes.go` | ~20 | RegisterRoutes |

### 4.4 Модели

```
OpenRequest:
  - warehouse_id  string
  - opening_cash  decimal

CashOpRequest:
  - type    string   (cash_in | cash_out)
  - amount  decimal
  - note    string

Shift:
  - id, user_id, warehouse_id
  - opened_at, closed_at
  - opening_cash, closing_cash
  - status (open | closed)

CashOp:
  - id, type, amount, note, created_at

ShiftReport:
  - shift         Shift
  - cash_ops      []CashOp
  - total_sales   decimal    (сумма всех SALE документов смены)
  - total_returns decimal    (сумма всех RETURN документов смены)
  - sales_cash    decimal    (продажи оплаченные наличными)
  - sales_card    decimal    (продажи оплаченные картой)
  - returns_cash  decimal    (возвраты наличными)
  - returns_card  decimal    (возвраты картой)
  - total_cash_in  decimal   (сумма внесений)
  - total_cash_out decimal   (сумма изъятий)
  - expected_cash  decimal   (расчётный остаток)
```

### 4.5 Формула ожидаемых наличных

```
expected_cash = opening_cash
              + sales_cash        (наличные от продаж)
              - returns_cash      (наличные возвраты покупателям)
              + total_cash_in     (внесения)
              - total_cash_out    (изъятия)
```

### 4.6 Repository — ключевые методы

- [ ] `Insert(ctx, tx, tenantID, shift)` — INSERT в documents.shifts
- [ ] `FindOpen(ctx, tenantID, userID)` — SELECT WHERE status='open'
- [ ] `Close(ctx, tx, tenantID, shiftID, closingCash)` — UPDATE status='closed', closed_at=now()
- [ ] `ByID(ctx, tenantID, shiftID)` — SELECT по ID
- [ ] `InsertCashOp(ctx, tx, shiftID, op)` — INSERT в shift_cash_ops
- [ ] `CashOps(ctx, shiftID)` — все операции смены
- [ ] `ShiftTotals(ctx, tenantID, shiftID)` — агрегация: JOIN documents + document_payments WHERE shift_id, GROUP BY type + method

### 4.7 Service — логика

**Open:**
- [ ] Проверить нет открытой смены (FindOpen)
- [ ] WithTx: вставить новую смену
- [ ] Вернуть ID

**Current:**
- [ ] FindOpen, вернуть ErrNoOpenShift если пусто

**CashOp:**
- [ ] FindOpen, проверить тип (cash_in/cash_out)
- [ ] WithTx: InsertCashOp

**Close:**
- [ ] FindOpen
- [ ] WithTx: рассчитать ShiftTotals → expected_cash → Close(closingCash = expected_cash)
- [ ] Вернуть полный ShiftReport

**Report:**
- [ ] ByID, CashOps, ShiftTotals
- [ ] Собрать ShiftReport

### 4.8 Ошибки

```
ErrShiftAlreadyOpen   — "у вас уже есть открытая смена"
ErrNoOpenShift        — "нет открытой смены"
ErrShiftAlreadyClosed — "смена уже закрыта"
ErrInvalidCashOpType  — "тип операции должен быть cash_in или cash_out"
```

**Done when:** полный цикл смены работает: открыть → продавать/возвращать → cash-op → закрыть → отчёт корректный.

---

## 5. Wiring — подключение в server.go

**Файл**: `backend/internal/shared/app/server.go`

- [ ] Импортировать `smarterp/backend/internal/features/shifts`
- [ ] Создать `shiftsRepo := shifts.NewRepository(store)`
- [ ] Создать `shiftsService := shifts.NewService(store, shiftsRepo)`
- [ ] Вызвать `shifts.RegisterRoutes(mux, shiftsService, tokens)`

**Done when:** сервер запускается, все новые эндпоинты доступны.

---

## 6. Модификация существующих файлов — сводка

| Файл | Изменение |
|------|-----------|
| `documents/models.go` | +PaymentInput, +Payment, +ShiftID в CreateRequest и Document, +Payments |
| `documents/service_posting_helpers.go` | docMeta() +RETURN, revenueForType() +RETURN, +buildCompositeReturnEntries() |
| `documents/service_posting.go` | buildEntriesForComposite() +ветка RETURN |
| `documents/repository_core.go` | INSERT/SELECT +shift_id |
| `documents/service_core.go` | +ReplacePayments в create/update, +LoadPayments в read |
| `shared/app/server.go` | +shifts wiring |

Новые файлы:
| Файл | Строк |
|------|-------|
| `migrations/000002_shifts_returns_payments.sql` | ~50 |
| `documents/repository_payments.go` | ~60 |
| `shifts/models.go` | ~60 |
| `shifts/errors.go` | ~10 |
| `shifts/repository.go` | ~120 |
| `shifts/service.go` | ~90 |
| `shifts/handler.go` | ~100 |
| `shifts/routes.go` | ~20 |

---

## 7. Порядок реализации

```
1. Миграция (§1)              — подготовить БД
2. Payments (§2)              — repository_payments + интеграция в documents
3. Returns (§3)               — posting helpers + posting + models
4. Shifts (§4)                — новая фича целиком (6 файлов)
5. Wiring (§5)                — server.go
6. Добавить shift_id (§6)     — в documents models/repository
```

Шаги 2–4 могут выполняться параллельно (разные файлы), но шаг 1 — первый, шаг 5 — последний.

---

## 8. Верификация

- [ ] Запустить миграцию: `go run ./cmd/migrate`
- [ ] Открыть смену → получить ID смены
- [ ] Создать SALE с payments (cash + card) и shift_id → провести → проверить леджер
- [ ] Создать RETURN с payments (cash) и shift_id → провести → проверить леджер (revenue отрицательная, profit отрицательный)
- [ ] Cash-op: внести наличные, изъять наличные
- [ ] Закрыть смену → проверить отчёт (expected_cash = формула из §4.5)
- [ ] Composite возврат: создать RETURN с composite товаром → проверить что компоненты корректно вернулись на склад
- [ ] Отчёты: прибыль за период корректно учитывает возвраты (отрицательная прибыль)
