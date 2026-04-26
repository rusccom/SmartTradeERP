# AGENTS.md

Локальные правила для `backend/internal/api`.

- Эта папка — только публичный API-layer: регистрация маршрутов, группировка URL и связывание URL с handlers.
- Не добавляй сюда бизнес-логику, SQL, расчеты остатков/себестоимости, валидацию документов или работу с ledger.
- Логика остается в `backend/internal/features/*`: handler, service, repository и models конкретной фичи.
- `api` может импортировать features и shared, но features не должны импортировать `api`.
- Новые URL группируй по UX-доменам:
  - `catalog.go` — products, variants, warehouses, customers.
  - `operations.go` — documents, shifts, reports.
  - `admin.go` — admin/auth routes.
  - `routes.go` — общий вход `Register`, health и route helpers.
- Для универсальных read-model ответов используй `include`, но только для близких данных:
  - `products?include=variants,stock,warehouses`
  - `warehouses?include=stock`
- Не превращай один endpoint в выдачу всей ERP. Не добавляй в products документы, клиентов, смены и несвязанные отчеты.
- Если endpoint нужен только для экрана, сначала определи его домен: основной ресурс (`products`, `warehouses`) или аналитика (`reports`).
- Сохраняй текущие правила проекта: маленькие функции, файлы до 300 строк, без монолитной сборки логики в `api`.
