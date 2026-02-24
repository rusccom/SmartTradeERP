# SmartERP — Architecture (финальная версия)

## 1. Обзор

SmartERP — multi-tenant SaaS для управленческого учёта малого бизнеса.
Фокус: себестоимость (AVCost), прибыль, остатки.

**Стек:**
- Backend: `Go`
- Frontend: `React` (deploy: `Cloudflare Pages`)
- Database: `PostgreSQL` (managed: `Aiven`)
- Приоритет: качество и надёжность выше скорости разработки

**Деплой:**
- Локальной сборки нет. Разработка ведётся сразу на удалённых средах.
- Frontend → Cloudflare Pages (push → auto deploy).
- Backend → облачный хостинг (Fly.io / Railway / VPS).
- БД → managed PostgreSQL на Aiven. Миграции накатываются на Aiven напрямую.

## 2. Слои системы

```
┌─────────────────────────────────────────────────────────────┐
│  PLATFORM (schema: platform)                                 │
│  tenants, tenant_users, platform_admins                      │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────┴──────────────────────────────────┐
│  WEB                                                         │
│  Public: /, /register, /login                                │
│  Admin:  /admin, /admin/dashboard/*                          │
│  Client: /dashboard/*                                        │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────┴──────────────────────────────────┐
│  API GATEWAY                                                 │
│  /api/admin/*    → Platform API (только owner)               │
│  /api/client/*   → ERP API (для клиентов, фильтр tenant_id)  │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────┴──────────────────────────────────┐
│  CATALOG (schema: catalog)                                   │
│  Шаблоны: products, variants, options, components            │
│  "ЧТО продаём и КАК устроено"                               │
└──────────────────────────┬──────────────────────────────────┘
                           │  variant_id, qty, price
┌──────────────────────────┴──────────────────────────────────┐
│  LEDGER (schema: ledger)                                     │
│  Факты: cost_ledger (IN/OUT, AVCost, COGS, profit)           │
│  "Чистая математика. Единственный источник истины."          │
└──────────────────────────────────────────────────────────────┘
```

## 3. Ключевые архитектурные решения

| # | Решение | Обоснование |
|---|---------|-------------|
| 1 | Catalog и Ledger — два отдельных слоя | Математика изолирована от бизнес-логики каталога |
| 2 | product_variant = реальный товар (SKU) | Отдельная таблица SKU не нужна, variant и есть SKU |
| 3 | Composite определяется наличием variant_components | Не ENUM типов, а данные определяют поведение |
| 4 | Рецепт и набор — одна механика | Оба раскладываются на компоненты при продаже |
| 5 | AVCost глобальный (не по складам) | Для малого бизнеса достаточно, упрощает перемещения |
| 6 | Snapshot компонентов в документе | Изменение рецепта не ломает старые документы |
| 7 | document_item_id в cost_ledger | Точная связь строки документа и движений в Ledger |
| 8 | Ретро-редактирование = снять + провести заново | Один алгоритм для всех случаев |
| 9 | Admin (/admin) отдельно от клиента | Безопасность, независимость, индустриальный стандарт |
| 10 | Multi-tenancy: shared DB + tenant_id | Для MVP, позже партиционирование |
| 11 | WithTx + context.Context — единый паттерн транзакций | Нет копипасты, автооткат при обрыве соединения |

## 4. Backend паттерн: Context + WithTx

Все SQL-операции проходят через один универсальный метод. Прямой `BeginTx` запрещён.

```
HTTP Request
  └→ r.Context()                     ← привязан к соединению клиента
       └→ Service.Method(ctx, ...)   ← бизнес-логика
            └→ Store.WithTx(ctx, fn) ← единая точка транзакций
                 └→ tx.ExecContext(ctx, ...)
                      └→ PostgreSQL          ← отменяется если ctx.Done()
```

```go
func (s *Store) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
    tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
    if err != nil { return err }
    defer tx.Rollback(ctx)
    if err := fn(tx); err != nil { return err }
    return tx.Commit(ctx)
}
```

Автоматическое поведение:
- Клиент отключился → ctx отменён → Rollback
- Ошибка в fn() → return err → Rollback через defer
- Паника → defer Rollback
- Всё ок → Commit

## 5. Финальная структура таблиц

### 5.1 Platform (schema: platform)

```sql
-- Клиенты-организации
tenants (
  id          UUID PRIMARY KEY,
  name        VARCHAR NOT NULL,
  status      VARCHAR NOT NULL DEFAULT 'trial',  -- trial/active/suspended
  plan        VARCHAR NOT NULL DEFAULT 'free',
  created_at  TIMESTAMP NOT NULL DEFAULT now()
)

-- Пользователи внутри тенанта
tenant_users (
  id            UUID PRIMARY KEY,
  tenant_id     UUID NOT NULL REFERENCES tenants(id),
  email         VARCHAR NOT NULL,
  password_hash VARCHAR NOT NULL,
  role          VARCHAR NOT NULL DEFAULT 'owner',  -- owner/manager/cashier
  is_active     BOOLEAN DEFAULT true,
  created_at    TIMESTAMP NOT NULL DEFAULT now(),
  UNIQUE(email)
)

-- Администраторы платформы (вы)
platform_admins (
  id            UUID PRIMARY KEY,
  email         VARCHAR NOT NULL UNIQUE,
  password_hash VARCHAR NOT NULL,
  created_at    TIMESTAMP NOT NULL DEFAULT now()
)
```

### 5.2 Catalog (schema: catalog)

```sql
-- Продукт = карточка/группа для витрины
products (
  id            UUID PRIMARY KEY,
  tenant_id     UUID NOT NULL,
  name          VARCHAR NOT NULL,
  is_composite  BOOLEAN NOT NULL DEFAULT false,
  created_at    TIMESTAMP NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP NOT NULL DEFAULT now()
)

-- Опции продукта (Цвет, Размер) — для UI
product_options (
  id          UUID PRIMARY KEY,
  product_id  UUID NOT NULL REFERENCES products(id),
  name        VARCHAR NOT NULL,       -- "Цвет"
  position    SMALLINT NOT NULL        -- 1, 2, 3
)

-- Значения опций (Красный, Синий) — для UI
product_option_values (
  id          UUID PRIMARY KEY,
  option_id   UUID NOT NULL REFERENCES product_options(id),
  value       VARCHAR NOT NULL,       -- "Красный"
  position    SMALLINT NOT NULL
)

-- Вариант = реальный товар (SKU). Всегда минимум 1 на продукт.
product_variants (
  id          UUID PRIMARY KEY,
  product_id  UUID NOT NULL REFERENCES products(id),
  name        VARCHAR,                -- "Красный / S" или "Default"
  sku_code    VARCHAR,                -- артикул "MLK-001"
  barcode     VARCHAR,                -- штрихкод
  unit        VARCHAR NOT NULL,       -- "л", "кг", "шт"
  price       DECIMAL(12,4),          -- рекомендованная цена продажи
  option1     VARCHAR,
  option2     VARCHAR,
  option3     VARCHAR,
  created_at  TIMESTAMP NOT NULL DEFAULT now()
)

-- Компоненты варианта (только для is_composite=true)
-- Snapshot-источник: при проведении документа копируется в document_item_components
variant_components (
  id                    UUID PRIMARY KEY,
  variant_id            UUID NOT NULL REFERENCES product_variants(id),
  component_variant_id  UUID NOT NULL REFERENCES product_variants(id),
  qty                   DECIMAL(12,3) NOT NULL,
  CONSTRAINT no_self_reference CHECK (variant_id != component_variant_id)
)

-- Склады
warehouses (
  id          UUID PRIMARY KEY,
  tenant_id   UUID NOT NULL,
  name        VARCHAR NOT NULL,
  is_default  BOOLEAN NOT NULL DEFAULT false,
  address     VARCHAR,
  is_active   BOOLEAN NOT NULL DEFAULT true,
  created_at  TIMESTAMP NOT NULL DEFAULT now()
)
```

### 5.3 Documents (schema: documents)

```sql
-- Документ (приход, продажа, списание, инвентаризация, перемещение)
documents (
  id                    UUID PRIMARY KEY,
  tenant_id             UUID NOT NULL,
  type                  VARCHAR NOT NULL,  -- RECEIPT/SALE/WRITEOFF/INVENTORY/TRANSFER
  date                  DATE NOT NULL,
  number                VARCHAR,
  status                VARCHAR NOT NULL DEFAULT 'draft',  -- draft/posted/cancelled
  warehouse_id          UUID REFERENCES warehouses(id),    -- для RECEIPT/SALE/WRITEOFF/INVENTORY
  source_warehouse_id   UUID REFERENCES warehouses(id),    -- для TRANSFER
  target_warehouse_id   UUID REFERENCES warehouses(id),    -- для TRANSFER
  note                  TEXT,
  created_at            TIMESTAMP NOT NULL DEFAULT now(),
  updated_at            TIMESTAMP NOT NULL DEFAULT now()
)

-- Строки документа (что продали/приняли — коммерческий факт)
document_items (
  id            UUID PRIMARY KEY,
  document_id   UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  variant_id    UUID NOT NULL REFERENCES product_variants(id),
  qty           DECIMAL(12,3) NOT NULL,
  unit_price    DECIMAL(12,4) NOT NULL,
  total_amount  DECIMAL(14,4) NOT NULL  -- qty × unit_price
)

-- Snapshot компонентов НА МОМЕНТ проведения (только для composite)
-- Гарантирует: изменение рецепта не ломает старые документы
document_item_components (
  id                    UUID PRIMARY KEY,
  document_item_id      UUID NOT NULL REFERENCES document_items(id) ON DELETE CASCADE,
  component_variant_id  UUID NOT NULL REFERENCES product_variants(id),
  qty_per_unit          DECIMAL(12,3) NOT NULL,  -- из рецепта на момент проведения
  qty_total             DECIMAL(12,3) NOT NULL   -- qty_per_unit × document_item.qty
)
```

### 5.4 Ledger (schema: ledger)

```sql
-- Журнал движений — ЕДИНСТВЕННЫЙ источник истины для математики
cost_ledger (
  id                BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  tenant_id         UUID NOT NULL,
  variant_id        UUID NOT NULL,          -- только физические варианты (is_composite=false)
  document_id       UUID NOT NULL,
  document_item_id  UUID NOT NULL,          -- точная связь со строкой документа
  warehouse_id      UUID NOT NULL,
  date              DATE NOT NULL,
  sequence_num      BIGINT NOT NULL,        -- порядок внутри (tenant_id, variant_id)
  type              VARCHAR(3) NOT NULL,    -- IN / OUT
  reason            VARCHAR NOT NULL,       -- PURCHASE/SALE/WRITEOFF/SHORTAGE/SURPLUS
                                            -- RETURN_IN/RETURN_OUT/TRANSFER_IN/TRANSFER_OUT
  qty               DECIMAL(12,3) NOT NULL,
  unit_price        DECIMAL(12,4) NOT NULL,
  total_amount      DECIMAL(14,4) NOT NULL, -- qty × unit_price
  -- вычисленные поля (пересчитываются)
  running_qty       DECIMAL(12,3) NOT NULL, -- глобальный остаток ПОСЛЕ операции
  running_avg       DECIMAL(12,4) NOT NULL, -- глобальный AVCost ПОСЛЕ операции
  cogs              DECIMAL(14,4),          -- себестоимость (только OUT)
  revenue           DECIMAL(14,4),          -- выручка (только OUT при SALE)
  profit            DECIMAL(14,4),          -- revenue - cogs (только OUT при SALE)
  created_at        TIMESTAMP NOT NULL DEFAULT now(),
  updated_at        TIMESTAMP NOT NULL DEFAULT now(),

  UNIQUE(tenant_id, variant_id, sequence_num)
)

-- Индексы
CREATE INDEX idx_cl_doc ON cost_ledger (document_id);
CREATE INDEX idx_cl_doc_item ON cost_ledger (document_item_id);
CREATE INDEX idx_cl_date ON cost_ledger (tenant_id, date, type);
CREATE INDEX idx_cl_warehouse ON cost_ledger (tenant_id, variant_id, warehouse_id);
```

## 6. Формулы Ledger

```
При ПРИХОДЕ (IN):
  new_avg = (old_qty × old_avg + in_qty × in_price) / (old_qty + in_qty)
  new_qty = old_qty + in_qty

При ПРОДАЖЕ (OUT):
  cogs    = out_qty × current_avg
  profit  = revenue - cogs
  new_qty = old_qty - out_qty
  new_avg = current_avg  (не меняется при продаже)

Composite (распределение revenue):
  share_i   = cogs_i / SUM(all cogs)
  revenue_i = total_sale_price × share_i
  profit_i  = revenue_i - cogs_i
```

## 7. Алгоритм проведения документа

```
1. Для каждой строки document_items:
   a. variant.is_composite?
      НЕТ → cost_ledger: INSERT (variant_id, qty, price)
      ДА  → скопировать variant_components → document_item_components
            → для каждого компонента: cost_ledger: INSERT
            → распределить revenue пропорционально cogs
2. Пересчитать цепочки затронутых variant_id
```

## 8. Алгоритм ретро-редактирования

```
BEGIN TRANSACTION
1. Собрать все document_item_id из документа
2. Запомнить earliest_sequence и affected_variants
3. DELETE FROM cost_ledger WHERE document_item_id IN (...)
4. DELETE FROM document_item_components WHERE document_item_id IN (...)
5. Применить правки к document_items
6. Провести документ заново (п.7)
7. Пересчитать цепочки affected ∪ new variants от earliest_seq
COMMIT
```

## 9. Остатки по складам

```sql
-- Глобальный остаток и AVCost (из последней строки цепочки)
SELECT running_qty, running_avg FROM cost_ledger
WHERE variant_id = ? ORDER BY sequence_num DESC LIMIT 1;

-- Остаток на конкретном складе
SELECT SUM(CASE WHEN type='IN' THEN qty ELSE -qty END) as warehouse_qty
FROM cost_ledger
WHERE variant_id = ? AND warehouse_id = ?;
```

## 10. Масштабирование (план на рост)

| Этап | Тенантов | Стратегия |
|------|----------|-----------|
| MVP | 0–50 | Одна таблица, индексы по (tenant_id, variant_id, seq) |
| Рост | 50–500 | Партиционирование cost_ledger по tenant_id |
| Масштаб | 500+ | Schema-per-tenant |
| Enterprise | 5000+ | Отдельная БД для крупных клиентов |
