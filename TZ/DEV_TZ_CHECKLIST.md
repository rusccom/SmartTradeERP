# SmartERP — ТЗ чеклист для разработчика

> Архитектурный приоритет: качество и надёжность выше скорости разработки.
> Связанный документ: `arch.md` — финальная структура таблиц и алгоритмы.

---

## 0. Технологический стек

- [ ] Backend: `Go`
- [ ] Frontend: `React` (deploy: `Cloudflare Pages`)
- [ ] Database: `PostgreSQL` (managed, хостинг: `Aiven`)
- [ ] Очереди/фон: `Redis + workers` (recalculate, background jobs)
- [ ] Контракты: OpenAPI-first + автоген TS-клиента для фронта

### Принцип деплоя
- [ ] **Локальной сборки нет.** Вся разработка ведётся сразу на удалённых средах.
- [ ] Frontend деплоится на `Cloudflare Pages` (push → auto deploy).
- [ ] Backend деплоится на облачный хостинг (выбрать: Fly.io / Railway / VPS).
- [ ] БД — managed PostgreSQL на `Aiven`. Локальный Postgres не используется.
- [ ] Миграции накатываются на Aiven напрямую через CI/CD или вручную на старте.
- [ ] Переменные окружения (DB connection string, JWT secret и т.д.) — через secrets хостинга, не в коде.
- [ ] Ветки: `main` = production. Для staging — отдельная Aiven DB instance или отдельная schema.

---

## 1. Platform и роли

- [ ] Таблица `tenants` (id, name, status, plan, created_at)
- [ ] Таблица `tenant_users` (id, tenant_id, email, password_hash, role)
- [ ] Таблица `platform_admins` (id, email, password_hash)
- [ ] Auth: JWT access + refresh tokens
- [ ] Пароли: argon2id
- [ ] Middleware: каждый запрос → проверка tenant_id → изоляция данных
- [ ] Роли tenant_users: owner / manager / cashier
- [ ] При регистрации тенанта → автоматически создать 1 склад (is_default=true)
- [ ] `/api/admin/auth/login` — отдельный от клиентского
- [ ] `/api/client/auth/login` — клиентский вход
- [ ] Никогда не джойнить данные разных тенантов

**Done when:** admin видит список тенантов, клиент логинится и видит только свои данные.

---

## 2. Web интерфейс (только заглушки на данном этапе)

> **На текущем этапе фронтенд — заглушки.** Маршруты созданы, страницы — пустые компоненты с заголовком. Данные не подтягиваются, формы не работают. Детальное наполнение страниц — отдельный этап после стабилизации backend.

### 2.1 Public зона
- [ ] `/` — заглушка Landing page
- [ ] `/register` — заглушка формы регистрации
- [ ] `/login` — заглушка формы логина

### 2.2 Admin зона
- [ ] `/admin` — заглушка логина администратора
- [ ] `/admin/dashboard` — заглушка сводки
- [ ] `/admin/tenants` — заглушка списка клиентов

### 2.3 Client зона
- [ ] `/dashboard` — заглушка сводки
- [ ] `/dashboard/products` — заглушка товаров
- [ ] `/dashboard/bundles` — заглушка составных товаров
- [ ] `/dashboard/warehouses` — заглушка складов
- [ ] `/dashboard/documents` — заглушка списка документов
- [ ] `/dashboard/documents/:id` — заглушка карточки документа
- [ ] `/dashboard/reports` — заглушка отчётов
- [ ] `/dashboard/settings` — заглушка настроек

**Done when:** все маршруты рендерятся, навигация между зонами работает, каждая страница показывает заголовок и placeholder-текст. Backend не вызывается.

---

## 3. API слой

### 3.1 Общие правила
- [ ] Envelope: `{ data, error, meta: { page, per_page, total } }`
- [ ] Пагинация: cursor или offset
- [ ] Ошибки: `{ code, message, details }`
- [ ] Все client endpoints → middleware: auth → scope → tenant_id

### 3.2 Admin API
- [ ] `POST /api/admin/auth/login`
- [ ] `GET /api/admin/tenants` — список тенантов
- [ ] `GET /api/admin/tenants/:id` — детали тенанта
- [ ] `GET /api/admin/stats` — статистика платформы

### 3.3 Client API — Products
- [ ] `GET /api/client/products` — список продуктов (фильтр: ?is_composite=true/false)
- [ ] `POST /api/client/products` — создать продукт (автоматически создаёт 1 variant "Default")
- [ ] `GET /api/client/products/:id`
- [ ] `PUT /api/client/products/:id`
- [ ] `DELETE /api/client/products/:id` — только если нет движений в cost_ledger

### 3.4 Client API — Variants (отдельный ресурс, variant = реальный товар)
- [ ] `GET /api/client/variants` — список (фильтр: ?product_id=xxx)
- [ ] `POST /api/client/variants` — создать (product_id в body)
- [ ] `GET /api/client/variants/:id`
- [ ] `PUT /api/client/variants/:id`
- [ ] `DELETE /api/client/variants/:id` — только если нет движений в cost_ledger
- [ ] `GET /api/client/variants/:id/components` — компоненты (для composite)
- [ ] `PUT /api/client/variants/:id/components` — задать/обновить компоненты
- [ ] `GET /api/client/variants/:id/stock` — остатки (глобально + по складам)

### 3.5 Client API — Warehouses
- [ ] `GET /api/client/warehouses`
- [ ] `POST /api/client/warehouses`
- [ ] `PUT /api/client/warehouses/:id`
- [ ] `DELETE /api/client/warehouses/:id` — только если нет движений

### 3.6 Client API — Documents
- [ ] `GET /api/client/documents` — список (фильтр по type, date, status)
- [ ] `POST /api/client/documents` — создать (draft)
- [ ] `GET /api/client/documents/:id` — карточка с прибылью по строкам
- [ ] `PUT /api/client/documents/:id` — редактировать (включая ретро)
- [ ] `POST /api/client/documents/:id/post` — провести
- [ ] `POST /api/client/documents/:id/cancel` — отменить проведение
- [ ] `DELETE /api/client/documents/:id` — только draft

### 3.7 Client API — Reports
- [ ] `GET /api/client/reports/profit` — прибыль за период
- [ ] `GET /api/client/reports/stock` — остатки (глобально и по складам)
- [ ] `GET /api/client/reports/top-products` — топ по прибыли
- [ ] `GET /api/client/reports/movements` — движения по variant_id

**Done when:** все endpoints работают, tenant_id изоляция проверена, ошибки возвращаются в едином формате.

---

## 4. Backend паттерны (Go)

### 4.1 Context — обязательно везде
- [ ] Все методы принимают `context.Context` первым аргументом
- [ ] Context берётся из `r.Context()` в HTTP handler — привязан к соединению клиента
- [ ] Context пробрасывается через всю цепочку: handler → service → store → SQL
- [ ] Если клиент отключился или нажал отмену → ctx отменяется → все SQL-запросы отменяются → транзакция откатывается автоматически

### 4.2 WithTx — единый метод для транзакций (без копипасты)
- [ ] Один универсальный метод `Store.WithTx(ctx, fn)` — точка входа для ВСЕХ транзакций
- [ ] Разработчик пишет только бизнес-логику внутри `fn(tx)`, всё остальное автоматически
- [ ] Rollback через defer (безопасно — после Commit не делает ничего)
- [ ] Автоматическое поведение:
  - Клиент отключился → ctx отменён → Rollback
  - Ошибка в fn() → return err → Rollback через defer
  - Всё ок → Commit
  - Паника → defer Rollback ловит

```go
// store.go — один раз написали, используем везде
type Store struct {
    pool *pgxpool.Pool
}

func (s *Store) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
    tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    if err := fn(tx); err != nil {
        return err
    }
    return tx.Commit(ctx)
}
```

### 4.3 Правила использования
- [ ] Никогда не вызывать `pool.BeginTx()` напрямую — только через `WithTx`
- [ ] Внутри `fn` все SQL-запросы через `tx.ExecContext(ctx, ...)` / `tx.QueryContext(ctx, ...)`
- [ ] Никогда не игнорировать ctx — если метод принимает context, передавать его в каждый SQL-запрос
- [ ] Handler → Service → Store: ctx прокидывается на каждом уровне

### 4.4 Архитектура вызовов

```
HTTP Request
  └→ r.Context()                     ← автоматически от net/http
       └→ Service.Method(ctx, ...)   ← бизнес-логика
            └→ Store.WithTx(ctx, fn) ← единая точка транзакций
                 └→ tx.ExecContext(ctx, ...) ← каждый SQL
                      └→ PostgreSQL          ← отменяется если ctx.Done()
```

**Done when:** все SQL-операции проходят через WithTx, нигде в коде нет прямого BeginTx, при обрыве соединения транзакция гарантированно откатывается.

---

## 5. Data model: Catalog

### 5.1 Таблицы
- [ ] `products` (id, tenant_id, name, is_composite)
- [ ] `product_options` (id, product_id, name, position)
- [ ] `product_option_values` (id, option_id, value, position)
- [ ] `product_variants` (id, product_id, name, sku_code, barcode, unit, price, option1, option2, option3)
- [ ] `variant_components` (id, variant_id, component_variant_id, qty)
- [ ] `warehouses` (id, tenant_id, name, is_default, address, is_active)

### 5.2 Инварианты
- [ ] У product ВСЕГДА минимум 1 variant
- [ ] is_composite=false → variant_components пусто, sku_code/unit обязательны
- [ ] is_composite=true → variant_components минимум 1 строка
- [ ] variant_components.component_variant_id → только варианты с is_composite=false
- [ ] Нет циклов в компонентах (composite не может содержать composite)
- [ ] Максимум 3 опции на продукт (option1, option2, option3)
- [ ] У тенанта минимум 1 склад с is_default=true

### 5.3 CRUD-сценарии
- [ ] Создание простого товара → product + 1 variant "Default"
- [ ] Создание товара с вариантами → product + N variants + options
- [ ] Создание composite → product(is_composite=true) + variant + components
- [ ] Редактирование рецепта → обновить variant_components (старые документы не затронуты)
- [ ] Удаление товара → только если нет движений в cost_ledger

**Done when:** CRUD работает, инварианты проверяются, ошибки понятные.

---

## 6. Data model: Ledger

### 6.1 Таблицы
- [ ] `cost_ledger` — полная структура из arch.md
- [ ] Индекс `UNIQUE(tenant_id, variant_id, sequence_num)`
- [ ] Индекс `(document_id)`
- [ ] Индекс `(document_item_id)`
- [ ] Индекс `(tenant_id, date, type)`
- [ ] Индекс `(tenant_id, variant_id, warehouse_id)`

### 6.2 Формулы
- [ ] IN: `new_avg = (old_qty × old_avg + in_qty × price) / (old_qty + in_qty)`
- [ ] OUT: `cogs = qty × current_avg`, `profit = revenue - cogs`
- [ ] AVCost не меняется при продаже
- [ ] Composite: revenue распределяется пропорционально cogs компонентов

### 6.3 sequence_num
- [ ] Отдельный на каждый `(tenant_id, variant_id)`
- [ ] Новая запись: `MAX(seq) + 1`
- [ ] Вставка в прошлое: сдвиг seq >= точки вставки на +1
- [ ] Сортировка по date + порядок создания внутри одной даты

### 6.4 Пересчёт
- [ ] `recalculate(tenant_id, variant_id, from_seq)` — последовательно от seq и далее
- [ ] Блокировка: `SELECT FOR UPDATE` на затронутые строки
- [ ] Транзакция: весь пересчёт в одной транзакции

### 6.5 Запросы (проверить что работают)
- [ ] Глобальный остаток: `SELECT running_qty FROM cost_ledger WHERE variant_id=? ORDER BY seq DESC LIMIT 1`
- [ ] Остаток по складу: `SUM(CASE WHEN type='IN' THEN qty ELSE -qty END) WHERE warehouse_id=?`
- [ ] Прибыль за период: `SUM(profit) WHERE date BETWEEN ... AND type='OUT'`
- [ ] Прибыль по строке документа: `SUM(profit) WHERE document_item_id=?`
- [ ] Прибыль по документу: `SUM(profit) WHERE document_id=?`

**Done when:** формулы корректны, пересчёт работает от любой точки, остатки сходятся.

---

## 7. Documents и операции учета

### 7.1 Таблицы
- [ ] `documents` (id, tenant_id, type, date, number, status, warehouse_id, source/target_warehouse_id, note)
- [ ] `document_items` (id, document_id, variant_id, qty, unit_price, total_amount)
- [ ] `document_item_components` (id, document_item_id, component_variant_id, qty_per_unit, qty_total)

### 7.2 Типы документов
- [ ] `RECEIPT` — приход от поставщика → IN, reason=PURCHASE
- [ ] `SALE` — продажа → OUT, reason=SALE
- [ ] `WRITEOFF` — списание (порча) → OUT, reason=WRITEOFF
- [ ] `INVENTORY` — инвентаризация → OUT(SHORTAGE) или IN(SURPLUS)
- [ ] `TRANSFER` — перемещение → OUT(TRANSFER_OUT) + IN(TRANSFER_IN)

### 7.3 Lifecycle: draft → posted → cancelled
- [ ] `draft` — можно свободно редактировать, нет записей в cost_ledger
- [ ] `posted` — записи в cost_ledger созданы, AVCost пересчитан
- [ ] `cancelled` — записи из cost_ledger удалены, цепочка пересчитана

### 7.4 Проведение документа
- [ ] Для каждой строки document_items:
  - [ ] variant.is_composite = false → 1 запись в cost_ledger
  - [ ] variant.is_composite = true → snapshot компонентов в document_item_components → N записей в cost_ledger → revenue распределить пропорционально cogs
- [ ] Записать document_item_id в каждую строку cost_ledger

### 7.5 Ретро-редактирование проведённого документа
- [ ] Алгоритм "снять + провести заново" (одна транзакция):
  1. [ ] Собрать все document_item_id
  2. [ ] Запомнить earliest_seq и affected_variants
  3. [ ] DELETE FROM cost_ledger WHERE document_item_id IN (...)
  4. [ ] DELETE FROM document_item_components WHERE document_item_id IN (...)
  5. [ ] Применить правки к document_items
  6. [ ] Провести заново
  7. [ ] Пересчитать цепочки от earliest_seq

### 7.6 Отображение документа (UI)
- [ ] Таблица строк: товар, кол-во, цена, выручка, себестоимость, прибыль
- [ ] Для composite строки — прибыль = сумма прибыли компонентов
- [ ] Итого по документу внизу
- [ ] Для draft — показывать расчётную прибыль (по текущему AVCost)

### 7.7 Инвентаризация (UI flow)
- [ ] Выбрать склад
- [ ] Система подтягивает все товары с учётным остатком на этом складе
- [ ] Пользователь вводит факт
- [ ] Система считает разницу и создаёт: SHORTAGE (OUT) или SURPLUS (IN)
- [ ] Цена для SHORTAGE/SURPLUS = текущий running_avg

### 7.8 Перемещение (TRANSFER)
- [ ] Выбрать source + target warehouse
- [ ] В cost_ledger: 2 строки в 1 транзакции (OUT + IN), обе с price=current_avg
- [ ] AVCost НЕ меняется (глобальный), running_qty НЕ меняется

**Done when:** все типы документов проводятся корректно, ретро-редактирование не ломает математику, snapshot рецептов работает.

---

## 8. Reports и аналитика

- [ ] Прибыль за период (с фильтром по складу, товару)
- [ ] Остатки (глобальные + по складам)
- [ ] Топ товаров по прибыли
- [ ] История движений по товару
- [ ] Потери (WRITEOFF + SHORTAGE за период)
- [ ] Admin: кол-во тенантов, активных, новых за период
- [ ] Экспорт CSV

**Done when:** все отчёты сходятся с данными cost_ledger, нет расхождений.

---

## 9. Нефункциональные требования

- [ ] API response < 200ms для CRUD, < 500ms для отчётов
- [ ] Пересчёт цепочки < 2s для 10,000 записей на variant
- [ ] Все операции с Ledger в транзакциях
- [ ] tenant_id в WHERE каждого запроса (без исключений)
- [ ] Structured logs (JSON)
- [ ] Error tracking (Sentry или аналог)
- [ ] DB backups: Aiven managed backups (автоматические)
- [ ] CI/CD pipeline: push → lint → test → build → deploy (без локальной сборки)
- [ ] Unit tests для формул AVCost
- [ ] Integration tests для пересчёта цепочки (подключение к test schema на Aiven)
- [ ] Все secrets (DB URL, JWT secret) — через переменные окружения хостинга, не в коде
- [ ] SSL/TLS соединение с Aiven PostgreSQL обязательно

---

## 10. Этапы реализации (roadmap)

### MVP-1: Ядро
- [ ] Инфраструктура: Aiven PostgreSQL + Cloudflare Pages + backend хостинг + CI/CD
- [ ] Platform: tenants, auth, admin dashboard
- [ ] Catalog: products, variants, warehouses
- [ ] Ledger: cost_ledger + формулы + пересчёт
- [ ] Documents: RECEIPT, SALE (draft → posted)
- [ ] Базовые отчёты: остатки, прибыль

### MVP-2: Полный учёт
- [ ] Composite товары (рецепты, наборы) + snapshot
- [ ] Ретро-редактирование проведённых документов
- [ ] WRITEOFF, INVENTORY, TRANSFER
- [ ] Прибыль по строкам документа
- [ ] Остатки по складам

### v2: Масштаб
- [ ] Партиционирование cost_ledger
- [ ] Снапшоты (ежемесячные) для быстрого пересчёта
- [ ] Биллинг / тарифы
