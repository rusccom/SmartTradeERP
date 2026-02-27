# ТЗ: Компонент DataTable — единая таблица для SmartTradeERP

## 1. Контекст проекта

### Стек фронтенда
- **React 18** + **Vite** (JSX, не TypeScript)
- **TailwindCSS v4** (`@tailwindcss/vite`, `tailwindcss` v4.2.1)
- **Framer Motion** — анимации
- **Lucide React** — иконки
- **React Router DOM v7** — маршрутизация

### Стек бэкенда
- **Go 1.23** (стандартная библиотека, `net/http` с паттерн-роутингом Go 1.22+)
- **PostgreSQL** + pgx v5 — прямые SQL-запросы
- **JWT** — аутентификация
- **Multi-tenant** — изоляция по `tenant_id` из JWT

### Существующая структура фронтенда
```
frontend/src/
├── app/
│   ├── AppFrame.jsx          # dev-layout (не используется в роутинге)
│   ├── styles.css            # глобальные стили (landing-zone)
│   └── tailwind.css          # tailwind @theme токены
├── features/
│   ├── admin/
│   │   ├── layout/AdminLayout.jsx
│   │   └── pages/AdminTenantsPage.jsx  ← placeholder
│   ├── dashboard/
│   │   ├── layout/ClientLayout.jsx
│   │   └── pages/
│   │       ├── ProductsPage.jsx        ← placeholder
│   │       ├── DocumentsPage.jsx       ← placeholder
│   │       ├── WarehousesPage.jsx      ← placeholder
│   │       ├── BundlesPage.jsx         ← placeholder
│   │       ├── ReportsPage.jsx         ← placeholder
│   │       └── ...
│   └── public/ (landing, login, register)
├── shared/
│   ├── api/http.js           # postJSON helper
│   ├── auth/session.js       # токены
│   ├── router/               # routes, guards
│   ├── seo/RouteSeo.jsx
│   └── ui/
│       ├── PlaceholderPage.jsx
│       └── workspace-layout.css
└── main.jsx
```

### Существующая структура бэкенда
```
backend/internal/
├── shared/
│   ├── httpx/
│   │   ├── request.go        # ParsePagination, DecodeJSON
│   │   └── response.go       # Envelope { data, error, meta }
│   ├── auth/                 # JWT, middleware chain
│   ├── tenant/               # tenant.FromContext
│   └── db/                   # PostgreSQL pool + Store
└── features/
    ├── products/             # handler, service, repository, models, routes
    ├── variants/
    ├── documents/
    ├── warehouses/
    ├── shifts/
    ├── ledger/
    └── reports/
```

### Дизайн-система workspace-zone (где живут таблицы)
- Фон страницы: `#f4f6f9`
- Текст основной: `#0d2c5a`
- Текст вторичный: `#5a6a80`
- Акцент (hover, active): `#ff6a3d`
- Бордеры: `rgba(13, 44, 90, 0.12–0.15)`
- Фон карточек: `#fff`
- Border-radius карточек: `20px`
- Border-radius кнопок/ссылок: `999px` (pill) или `12px`
- Фон кнопок active: `#0d2c5a`, текст: `#fff`
- Шрифт: `"Inter", system-ui, sans-serif`
- Анимация появления: `translateY(6px) → 0` за 220ms

---

## 2. Задача

Создать **единый переиспользуемый компонент** `<DataTable>` в `src/shared/ui/data-table/`.

Основной режим проекта: `server` (production-first).
1. Страница подключает `preset` + `useServerDataTable`.
2. В `<DataTable>` передаются controlled state и данные из хука.
3. Сортировка, фильтрация, пагинация и стили централизованы в общем `DataTable` + `shared/model/data-table`.

---

## 3. Зависимость

> **ВАЖНО:** На этапе разработки **ничего не устанавливается локально**. Все `npm install` и запуски (`npm run dev`, `npm run build`) будут выполнены при деплое на сервер. Здесь мы только пишем код.

### Новый пакет (добавить в `package.json`)

В `dependencies` добавить `"@tanstack/react-table": "^8.21.3"` (вручную в `package.json`, без `npm install`).

Для базового внедрения достаточно `@tanstack/react-table`. Иконки — из уже установленного `lucide-react`.
`@tanstack/react-virtual` добавляется отдельно только при необходимости (см. раздел 17.2).

### Консистентность `shared/api/client` (JS-проект)

Проект на JSX/JS, поэтому конфиг `apiPaths` должен быть JS-модулем.

Production-правило:
1. Базовый файл: `src/shared/api/client.js`.
2. Если в ветке остался `src/shared/api/client.ts`, переименовать его в `src/shared/api/client.js`.
3. Убрать TS-аннотации из файла.
4. Импортировать без расширения: `import { apiPaths } from "../../../shared/api/client"`.

### Расширение `shared/api/http.js`

Существующий `http.js` содержит только `postJSON` **без авторизации**. Необходимо:
1. Добавить общий `resolveToken(path)` — авторизация по пути
2. Обновить `createPostOptions` — добавить `Authorization`
3. Добавить `getJSON` с поддержкой query params, signal и meta

```js
import { getAdminToken, getClientToken } from "../auth/session";

// --- Авторизация (общая для GET и POST) ---

function resolveToken(path) {
  if (path.startsWith("/api/admin")) return getAdminToken() || "";
  if (path.startsWith("/api/client")) return getClientToken() || "";
  return getClientToken() || getAdminToken() || "";
}

function authHeaders(path) {
  const token = resolveToken(path);
  return token ? { Authorization: `Bearer ${token}` } : {};
}

// --- Обновить существующий postJSON ---

// Заменить createPostOptions на:
function createPostOptions(path, payload) {
  return {
    method: "POST",
    headers: { "Content-Type": "application/json", ...authHeaders(path) },
    body: JSON.stringify(payload),
  };
}

// Обновить вызов в postJSON:
export async function postJSON(path, payload) {
  const response = await fetch(createURL(path), createPostOptions(path, payload));
  return parseEnvelope(response);
}

// --- Новый getJSON ---

export async function getJSON(path, params, signal) {
  const url = buildURL(path, params);
  const response = await fetch(url, createGetOptions(path, signal));
  return parseEnvelopeWithMeta(response);
}

function buildURL(path, params) {
  const url = new URL(createURL(path));
  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== "") {
        url.searchParams.set(key, String(value));
      }
    });
  }
  return url.toString();
}

function createGetOptions(path, signal) {
  return { method: "GET", headers: authHeaders(path), signal };
}

async function parseEnvelopeWithMeta(response) {
  const body = await parseBody(response);
  if (!response.ok || body.error) {
    throw new Error(getErrorMessage(body, response.status));
  }
  return { data: body.data ?? null, meta: body.meta ?? null };
}
```

**Важно:** `buildURL` обязан переиспользовать уже существующий `createURL(path)` из `http.js`.
Дублировать логику base URL запрещено.

**Сигнатуры:**
- `postJSON(path, payload) → Promise<data>` — теперь с Authorization
- `getJSON(path, params?, signal?) → Promise<{ data, meta }>`
- `signal` — `AbortSignal` для отмены запроса (из `AbortController`)
- `Authorization: Bearer <token>` добавляется автоматически по пути для обоих методов

**БЛОКЕР:** если в `session.js` нет `getClientToken/getAdminToken`, добавить:

```js
export function getAdminToken() { return getToken(ADMIN_TOKEN_KEY); }
export function getClientToken() { return getToken(CLIENT_TOKEN_KEY); }
```

---

## 4. Файловая структура

### Фронтенд

```
src/shared/
├── ui/data-table/              ← РЕНДЕР (визуальные компоненты)
│   ├── DataTable.jsx           # публичный компонент (единственный экспорт)
│   ├── DataTableHeader.jsx     # <thead> + сортировка
│   ├── DataTableBody.jsx       # <tbody> + зебра + hover + expand
│   ├── DataTablePagination.jsx # навигация по страницам
│   ├── DataTableFilter.jsx     # фильтр per-column
│   ├── DataTableToolbar.jsx    # глобальный поиск + счётчик
│   ├── DataTableError.jsx      # блок ошибки + кнопка «Повторить»
│   └── data-table.css          # все стили
│
├── model/data-table/           ← ЛОГИКА (state, fetch, preset)
│   ├── createTablePreset.js
│   ├── useServerDataTable.js
│   ├── tableState.js
│   ├── toQueryParams.js
│   └── useDebounce.js
│
└── api/http.js                 ← ТРАНСПОРТ (+ getJSON)
```

### Бэкенд (расширения)

```
backend/internal/shared/httpx/
├── request.go     ← добавить ParseSort, ParseSearch, ParseFilters
└── response.go    ← без изменений
```

---

## 5. API компонента DataTable

### Props (server mode — controlled)

| Prop | Тип | Обяз. | Default | Описание |
|---|---|:---:|---|---|
| `columns` | `Array<ColumnDef>` | да | — | Описание колонок |
| `data` | `Array<Object>` | да | — | Массив строк данных |
| `getRowId` | `(row, index) => string` | да | — | Стабильный id строки (всегда из `preset.rowId`) |
| `searchable` | `boolean` | нет | `true` | Показывать глобальный поиск |
| `onRowClick` | `(row) => void` | нет | — | Клик по строке |
| `toolbar` | `ReactNode` | нет | — | Кнопки действий в тулбаре |
| `emptyText` | `string` | нет | `"Нет данных"` | Текст при пустой таблице |
| `rowCount` | `number` | да | — | Общее кол-во строк на сервере |
| `loading` | `boolean` | да | `false` | Состояние загрузки |
| `error` | `string \| null` | да | `null` | Текст ошибки |
| `onRetry` | `function` | нет | — | Callback повторной загрузки |
| `sorting` | `Array` | да | — | Controlled sorting state |
| `onSortingChange` | `function` | да | — | Setter для sorting |
| `columnFilters` | `Array` | да | — | Controlled filters state |
| `onColumnFiltersChange` | `function` | да | — | Setter для filters |
| `globalFilter` | `string` | да | — | Controlled search state |
| `onGlobalFilterChange` | `function` | да | — | Setter для search |
| `pagination` | `{ pageIndex, pageSize }` | да | — | Controlled pagination |
| `onPaginationChange` | `function` | да | — | Setter для pagination |
| `expandable` | `boolean` | нет | `false` | Включить раскрытие строк |
| `getSubRows` | `(row) => Promise<Row[]>` | нет | — | Подгрузка вложенных строк |

> **Правило:** `DataTable` — controlled-компонент. State приходит извне через `useServerDataTable`.
> Все `on[State]Change` callbacks получают `updater`, который может быть и значением, и функцией.

### ColumnDef (описание колонки)

| Поле | Тип | Обяз. | Описание |
|---|---|:---:|---|
| `accessorKey` | `string` | да | Ключ поля в объекте данных |
| `header` | `string` | да | Заголовок колонки |
| `cell` | `(value, row) => ReactNode` | нет | Кастомный рендер ячейки |
| `enableSorting` | `boolean` | нет | По умолчанию `true` |
| `enableFilter` | `boolean` | нет | По умолчанию `false` |
| `filterVariant` | `"text" \| "select" \| "range"` | нет | Тип фильтра |
| `filterOptions` | `Array<string \| { label, value }>` | нет | Опции для select-фильтра |
| `size` | `number` | нет | Ширина колонки в px |

---

## 6. Реализация UI-компонентов

### 6.1 `data-table.css`

> **Почему чистый CSS, а не Tailwind-классы:** DataTable — сложный компонент с десятками стилевых правил, состояниями (hover, active, disabled, sorted) и responsive-логикой. Inline Tailwind-классы сделали бы JSX нечитаемым. Отдельный CSS-файл с неймспейсом `.dt-*` изолирует стили и не конфликтует с Tailwind. Аналогичный подход уже используется в проекте (`workspace-layout.css`).

Все классы начинаются с `.dt-`.

| Класс | Стили |
|---|---|
| `.dt-wrapper` | `border-radius: 20px`, `background: #fff`, `border: 1px solid rgba(13,44,90,0.12)`, анимация `ws-rise` 220ms |
| `.dt-toolbar` | flex, `justify-content: space-between`, `padding: 16px 20px`, `border-bottom: 1px solid rgba(13,44,90,0.08)` |
| `.dt-search` | `border-radius: 12px`, `border: 1px solid rgba(13,44,90,0.15)`, `padding: 8px 12px 8px 36px`, `min-height: 38px`, `font-size: 13px`. Focus: `outline: 2px solid rgba(255,106,61,0.3)` |
| `.dt-search-wrap` | `position: relative` (иконка Search + input) |
| `.dt-count` | `color: #5a6a80`, `font-size: 13px` |
| `.dt-table-scroll` | `overflow-x: auto` |
| `.dt-table` | `width: 100%`, `border-collapse: collapse` |
| `.dt-th` | `padding: 10px 16px`, `font-size: 12px`, `font-weight: 600`, `text-transform: uppercase`, `letter-spacing: 0.04em`, `color: #5a6a80`, `background: #f8f9fb`, `border-bottom: 2px solid rgba(13,44,90,0.1)`, `white-space: nowrap`, `user-select: none` |
| `.dt-th--sortable` | `cursor: pointer`. Hover: `color: #0d2c5a` |
| `.dt-th--sorted` | `color: #0d2c5a` |
| `.dt-sort-icon` | `margin-left: 4px`, `width: 14px`, `height: 14px`, `opacity: 0.5`. Active: `opacity: 1`, `color: #ff6a3d` |
| `.dt-filter` | `width: 100%`, `margin-top: 6px`, `padding: 4px 8px`, `border: 1px solid rgba(13,44,90,0.12)`, `border-radius: 8px`, `font-size: 12px` |
| `.dt-td` | `padding: 10px 16px`, `font-size: 14px`, `color: #0d2c5a`, `border-bottom: 1px solid rgba(13,44,90,0.06)` |
| `.dt-row:nth-child(even)` | `background: #f8f9fb` |
| `.dt-row:hover` | `background: rgba(255,106,61,0.04)` |
| `.dt-row--clickable` | `cursor: pointer` |
| `.dt-row--sub` | `background: #f0f4f8`, вложенная строка (вариант), отступ через `padding-left` |
| `.dt-expand-btn` | `cursor: pointer`, `background: none`, `border: none`, `padding: 2px` |
| `.dt-empty` | `text-align: center`, `padding: 48px 20px`, `color: #5a6a80` |
| `.dt-pagination` | flex, `justify-content: space-between`, `align-items: center`, `padding: 12px 20px`, `border-top: 1px solid rgba(13,44,90,0.08)`, `font-size: 13px` |
| `.dt-page-btn` | `min-width: 32px`, `min-height: 32px`, `border-radius: 8px`, `border: 1px solid rgba(13,44,90,0.15)`, `background: #fff`, `color: #0d2c5a`. Hover: `border-color: rgba(255,106,61,0.65)`, `color: #ff6a3d`. Disabled: `opacity: 0.4` |
| `.dt-page-btn--active` | `background: #0d2c5a`, `color: #fff`, `border-color: transparent` |
| `.dt-page-size` | стили аналогичны `.dt-filter` |
| `.dt-page-info` | `color: #5a6a80` |
| `.dt-loading` | `opacity: 0.5`, `pointer-events: none` |
| `.dt-error` | `text-align: center`, `padding: 16px 20px`, `color: #dc2626` |
| `.dt-error-btn` | `margin-left: 12px` (стиль `.dt-page-btn`) |

**Responsive (max-width: 768px):** `.dt-toolbar` и `.dt-pagination` → `flex-direction: column`, `gap: 8px`.

### 6.2 `DataTableFilter.jsx`

**Props:** `column` (TanStack column object)

Логика по `column.columnDef.meta?.filterVariant`:
- `"text"` → `<input type="text">`, при вводе `column.setFilterValue(value)`
- `"select"` → `<select>` из `meta.filterOptions`, первый option пустой (все)
- `"range"` → два `<input type="number">` (min, max)

Debounce выполняется только в `useServerDataTable`, не в UI.

### 6.3 `DataTableHeader.jsx`

**Props:** `table` (TanStack table instance)

- `table.getHeaderGroups()` → `<tr>` для каждой группы
- Заголовок через `flexRender(header.column.columnDef.header, header.getContext())`
- Сортировка: `ArrowUp`, `ArrowDown`, `ArrowUpDown` из `lucide-react`
- Фильтр: `<DataTableFilter column={...} />` под заголовком если `column.getCanFilter()`

### 6.4 `DataTableBody.jsx`

**Props:** `table`, `onRowClick`, `emptyText`, `expandable`, `getSubRows`

- `table.getRowModel().rows` — видимые строки
- Пустой массив → `<td colSpan>` с `emptyText`
- `dt-row--clickable` если `onRowClick` задан
- **Expand-логика** (если `expandable=true`):
  - Первая ячейка содержит chevron-кнопку (ChevronRight / ChevronDown из lucide-react)
  - Клик → вызов `getSubRows(row.original)` → рендер вложенных строк с классом `dt-row--sub`
  - Вложенные строки рендерятся с отступом `padding-left: 32px` в первой ячейке
  - Состояние expanded хранится локально в компоненте (`useState` set of row ids)
  - Подгруженные subRows кешируются в `useRef` map по row id

### 6.5 `DataTablePagination.jsx`

**Props:** `table`

- Кнопки: `«`, `‹`, номера страниц (макс 5 вокруг текущей), `›`, `»`
- Select pageSize: `[10, 20, 50, 100]`
- Текст: `"1–20 из 342"`
- Методы: `table.setPageIndex()`, `table.previousPage()`, `table.nextPage()`, `table.getCanPreviousPage()`, `table.getCanNextPage()`

### 6.6 `DataTableToolbar.jsx`

**Props:** `table`, `globalFilter`, `onGlobalFilterChange`, `searchable`, `toolbar`, `rowCount`

- Слева: поле поиска с иконкой `Search` (если `searchable`). Без debounce в UI
- Справа: счётчик `"Всего: ${rowCount}"` + пользовательский `toolbar` (ReactNode)

### 6.7 `DataTableError.jsx`

**Props:** `message`, `onRetry`

```jsx
function DataTableError({ message, onRetry }) {
  return (
    <div className="dt-error">
      <span>{message}</span>
      {onRetry && (
        <button className="dt-page-btn dt-error-btn" onClick={onRetry}>
          Повторить
        </button>
      )}
    </div>
  );
}
```

### 6.8 `DataTable.jsx` — главный компонент

Единственный экспортируемый компонент. Все остальные — внутренние.

1. **Маппинг columns** — функция `mapColumns(columns)`:
   - `cell` → обёртка: `(info) => col.cell(info.getValue(), info.row.original)`
   - `enableColumnFilter` → `true` если `col.enableFilter` или `col.filterVariant`
   - `meta: { filterVariant, filterOptions }`

2. **TanStack table instance:**
   ```jsx
   const table = useReactTable({
     data,
     columns: mappedColumns,
     getRowId: resolvedGetRowId,
     state: { sorting, columnFilters, globalFilter, pagination },
     onSortingChange,
     onColumnFiltersChange,
     onGlobalFilterChange,
     onPaginationChange,
     manualPagination: true,
     manualSorting: true,
     manualFiltering: true,
     rowCount,
     getCoreRowModel: getCoreRowModel(),
   });
   ```

3. **Рендер:**
   ```jsx
   <div className="dt-wrapper">
     <DataTableToolbar ... />
     {error && <DataTableError message={error} onRetry={onRetry} />}
     <div className={`dt-table-scroll ${loading ? "dt-loading" : ""}`}>
       <table className="dt-table">
         <DataTableHeader table={table} />
         <DataTableBody table={table} expandable={expandable}
           getSubRows={getSubRows} onRowClick={onRowClick} emptyText={emptyText} />
       </table>
     </div>
     <DataTablePagination table={table} />
   </div>
   ```

---

## 7. Реализация model-слоя

### 7.1 `tableState.js`

```js
export const FALLBACK_STATE = {
  pagination: { pageIndex: 0, pageSize: 20 },
  sorting: [],
  globalFilter: "",
  columnFilters: [],
};

// Смена фильтра/поиска/сортировки → pageIndex = 0
export function applyStateChange(prev, patch) { ... }
```

Ограничения: ≤50 строк.

### 7.2 `toQueryParams.js`

Единственное место формирования базовых query-полей.

```js
export function toQueryParams(state, preset) { ... }

// Базовые поля (ВСЕГДА):
//   page       = state.pagination.pageIndex + 1
//   per_page   = state.pagination.pageSize
// Сортировка (если preset.capabilities.sorting === true):
//   sort_by    = state.sorting[0]?.id
//   sort_dir   = state.sorting[0]?.desc ? "desc" : "asc"
// Поиск (если preset.capabilities.search === true и строка не пустая):
//   search     = state.globalFilter
// Feature-override (из preset.mapStateToQuery, если задан):
//   const custom = preset.mapStateToQuery ? preset.mapStateToQuery(state) : {}
//   const overriddenKeys = new Set(Object.keys(custom))
// Авто-маппинг columnFilters (пропуская ключи из mapStateToQuery):
//   state.columnFilters.forEach(f => {
//     if (!overriddenKeys.has(f.id)) params[f.id] = f.value
//   })
//   Object.assign(params, custom)
// Очистить undefined/null/"" значения
```

**Авто-маппинг:** `columnFilters` маппятся на query params по `accessorKey` автоматически, **но ключи, переопределённые в `mapStateToQuery`, исключаются из авто-маппинга** (чтобы избежать дублей вроде `doc_type=SALE&type=SALE`). `mapStateToQuery` нужен только если query param key отличается от `accessorKey` (например, колонка `doc_type` → query param `type`).

`mapStateToQuery` запрещено возвращать зарезервированные поля: `page`, `per_page`, `sort_by`, `sort_dir`, `search`.

Ограничения: ≤40 строк.

### 7.3 `createTablePreset.js`

Фабрика preset-объекта.

```js
export function createTablePreset(config) { ... }
// 1. Проверить обязательные поля (id, rowId, columns, fetchPage)
// 2. Мёрж defaultState с FALLBACK_STATE
// 3. Мёрж capabilities с { sorting: false, search: false }
// 4. Object.freeze(result)
```

Ограничения: ≤40 строк.

### 7.4 `useDebounce.js`

```js
export function useDebounce(value, delay = 300) { ... }
// useState + useEffect + setTimeout + cleanup
```

Ограничения: ≤40 строк.

### 7.5 `useServerDataTable.js`

Основной хук для server-режима.

```js
export function useServerDataTable(preset) { ... }

// Возвращает:
// {
//   data: Array,           — текущие строки (keep previous при загрузке)
//   total: number,         — общее количество строк
//   loading: boolean,
//   error: string | null,
//   retry: () => void,     — повторить последний запрос
//   tableState: {          — spread в <DataTable>
//     sorting, onSortingChange,
//     columnFilters, onColumnFiltersChange,
//     globalFilter, onGlobalFilterChange,
//     pagination, onPaginationChange,
//   },
// }
```

Внутренняя логика:
1. `useState` для sorting, columnFilters, globalFilter, pagination — из `preset.defaultState`
2. `useState` для data, total, loading, error
3. `useRef` для AbortController
4. `useDebounce` для columnFilters (300ms) и globalFilter (300ms)
5. `useEffect`:
   - Abort предыдущий запрос
   - `query = toQueryParams(requestState, preset)`
   - `preset.fetchPage({ query, signal })`
   - Success → `setData(rows)`, `setTotal(total)`, `setError(null)`
   - AbortError → игнорировать
   - Другая ошибка → `setError(message)`, НЕ очищать data (keep previous)
6. Handlers применяют `applyStateChange` (авто-сброс pageIndex)
7. Handlers принимают updater (значение ИЛИ функция)

Ограничения: ≤100 строк. Разнести на функции: `createQueryState`, `runFetch`, `handleFetchSuccess`, `handleFetchError`, `createTableHandlers`.

---

## 8. Контракт `tablePreset`

Каждая фича создаёт один preset-файл:

```
src/features/dashboard/pages/table/productsTablePreset.js
src/features/dashboard/pages/table/documentsTablePreset.js
src/features/admin/pages/table/adminTenantsTablePreset.js
```

### Поля preset

| Поле | Тип | Обяз. | Описание |
|---|---|:---:|---|
| `id` | `string` | да | Стабильный id таблицы |
| `rowId` | `(row) => string` | да | Уникальный id строки |
| `columns` | `Array<ColumnDef>` | да | Колонки (раздел 5) |
| `fetchPage` | `({ query, signal }) => Promise<{ rows, total }>` | да | Загрузка страницы |
| `defaultState` | `TableState` | нет | Мёржится с `FALLBACK_STATE` |
| `capabilities` | `{ sorting?, search? }` | нет | Поддержка backend. Default: `{ sorting: false, search: false }` |
| `mapStateToQuery` | `(state) => object` | нет | Override авто-маппинга: только если accessorKey ≠ query param key |

### Пример: productsTablePreset

```js
import { createTablePreset } from "../../../../shared/model/data-table/createTablePreset";
import { apiPaths } from "../../../../shared/api/client";
import { getJSON } from "../../../../shared/api/http";

export const productsTablePreset = createTablePreset({
  id: "products",
  rowId: (row) => row.id,
  columns: [
    { accessorKey: "name", header: "Название", enableFilter: true },
    {
      accessorKey: "is_composite", header: "Составной",
      filterVariant: "select",
      filterOptions: [{ value: "true", label: "Да" }, { value: "false", label: "Нет" }],
      cell: (value) => (value ? "Да" : "Нет"),
    },
    { accessorKey: "updated_at", header: "Обновлён" },
  ],
  defaultState: {
    pagination: { pageIndex: 0, pageSize: 20 },
    sorting: [],
    globalFilter: "",
    columnFilters: [],
  },
  capabilities: { sorting: true, search: true },
  // mapStateToQuery не нужен: accessorKey "is_composite" = query param "is_composite" (авто-маппинг)
  fetchPage: async ({ query, signal }) => {
    const { data, meta } = await getJSON(apiPaths.products, query, signal);
    return { rows: data ?? [], total: meta?.total ?? 0 };
  },
});
```

### Пример: documentsTablePreset

```js
export const documentsTablePreset = createTablePreset({
  id: "documents",
  rowId: (row) => row.id,
  columns: [
    { accessorKey: "number", header: "Номер" },
    {
      accessorKey: "doc_type", header: "Тип",
      filterVariant: "select",
      filterOptions: [
        { value: "RECEIPT", label: "Приёмка" },
        { value: "SALE", label: "Продажа" },
        { value: "WRITEOFF", label: "Списание" },
        { value: "TRANSFER", label: "Перемещение" },
        { value: "RETURN", label: "Возврат" },
      ],
    },
    { accessorKey: "status", header: "Статус" },
    { accessorKey: "total_cost", header: "Сумма" },
    { accessorKey: "date", header: "Дата" },
  ],
  capabilities: { sorting: true, search: true },
  // mapStateToQuery: doc_type → type (accessorKey ≠ query param)
  // Ключ "doc_type" исключается из авто-маппинга автоматически (см. toQueryParams).
  // status маппится автоматически (accessorKey = query param).
  mapStateToQuery: (state) => ({
    type: state.columnFilters.find(f => f.id === "doc_type")?.value,
  }),
  fetchPage: async ({ query, signal }) => {
    const { data, meta } = await getJSON(apiPaths.documents, query, signal);
    return { rows: data ?? [], total: meta?.total ?? 0 };
  },
});
```

### Что запрещено

- Дублировать query-логику в page-компонентах
- Дублировать debounce/cancel/error обработчики в каждой фиче
- Писать отдельные "локальные таблицы" вместо общего `DataTable`
- Возвращать из `mapStateToQuery` зарезервированные поля

---

## 9. Использование на странице

```jsx
import DataTable from "../../../shared/ui/data-table/DataTable";
import { useServerDataTable } from "../../../shared/model/data-table/useServerDataTable";
import { productsTablePreset } from "./table/productsTablePreset";

function ProductsPage() {
  const { data, total, loading, error, retry, tableState } =
    useServerDataTable(productsTablePreset);

  return (
    <DataTable
      columns={productsTablePreset.columns}
      data={data}
      getRowId={productsTablePreset.rowId}
      searchable={productsTablePreset.capabilities.search === true}
      rowCount={total}
      loading={loading}
      error={error}
      onRetry={retry}
      expandable={true}
      getSubRows={(row) =>
        getJSON(apiPaths.variants, { product_id: row.id })
          .then(r => r.data ?? [])
      }
      {...tableState}
      toolbar={<button>+ Добавить товар</button>}
    />
  );
}
```

---

## 10. Expand-строки (варианты товара)

### Назначение

Товары могут иметь варианты (размер, цвет и т.д.). При клике на chevron в строке товара — подгружаются варианты и рендерятся как вложенные строки.

### Механизм

1. `DataTable` получает `expandable={true}` и `getSubRows={(row) => Promise<Row[]>}`
2. `DataTableBody` управляет expand-состоянием локально (`useState` set of expanded row ids)
3. При клике на chevron:
   - Если строка уже раскрыта → свернуть (убрать из set)
   - Если нет → вызвать `getSubRows(row.original)`, закешировать результат в `useRef`, добавить в set
4. Вложенные строки рендерятся сразу после родительской `<tr>` с классом `dt-row--sub`
5. Колонки вложенных строк маппятся на те же `columns` родительской таблицы (если у варианта нет поля — ячейка пустая)

### API для подгрузки вариантов (уже существует)

```
GET /api/client/variants?product_id={id}
```

---

## 11. Backend: API-контракт списка

### 11.1 Формат ответа (уже реализован)

```json
{
  "data": [ ... ],
  "error": null,
  "meta": { "page": 1, "per_page": 20, "total": 142 }
}
```

### 11.2 Универсальные query-параметры

Каждый list-эндпоинт должен поддерживать единый набор:

| Параметр | Тип | Default | Описание |
|---|---|---|---|
| `page` | int | 1 | Номер страницы |
| `per_page` | int | 20 | Записей на страницу (max 100) |
| `sort_by` | string | `created_at` | Поле сортировки (валидируется по whitelist) |
| `sort_dir` | string | `desc` | `asc` / `desc` |
| `search` | string | `""` | Поиск по текстовым полям (ILIKE) |
| `filter.*` | string | — | Фильтры по конкретным полям |

### 11.3 Уже поддерживаемые фильтры

- `products`: `is_composite`, `page`, `per_page`
- `documents`: `type`, `status`, `date`, `page`, `per_page`
- `variants`: `product_id`, `page`, `per_page`
- `admin/tenants`: `page`, `per_page`
- `reports/*` — агрегатные endpoints, не list-контракт (нет `page/per_page`, нет `meta.total`)

### 11.4 Правило capabilities

- Products и Documents: после обновления handlers в Этапе A → `capabilities: { sorting: true, search: true }`
- Для endpoints, которые ещё не поддерживают sort/search (warehouses, admin/tenants) → `capabilities.sorting = false`, `capabilities.search = false`
- `toQueryParams` отправляет `sort_by/sort_dir` и `search` **только** если соответствующий capability = `true`

---

## 12. Backend: расширение `httpx/request.go`

Добавить три функции-парсера рядом с существующим `ParsePagination`:

### 12.1 `ParseSort`

```go
// ParseSort парсит sort_by и sort_dir, валидирует по allowedFields.
// Возвращает (sortBy, sortDir). Если поле невалидно — возвращает fallback.
func ParseSort(r *http.Request, allowedFields []string, fallbackField string) (string, string) {
    sortBy := r.URL.Query().Get("sort_by")
    sortDir := r.URL.Query().Get("sort_dir")

    if !contains(allowedFields, sortBy) {
        sortBy = fallbackField
    }
    if sortDir != "asc" && sortDir != "desc" {
        sortDir = "desc"
    }
    return sortBy, sortDir
}
```

### 12.2 `ParseSearch`

```go
// ParseSearch возвращает очищенную строку поиска (trimmed, max 200 символов).
func ParseSearch(r *http.Request) string {
    s := strings.TrimSpace(r.URL.Query().Get("search"))
    if len(s) > 200 {
        return s[:200]
    }
    return s
}
```

### 12.3 `ParseFilters`

```go
// ParseFilters возвращает map[string]string из query-параметров,
// пропуская ключи не из allowedKeys.
func ParseFilters(r *http.Request, allowedKeys []string) map[string]string {
    result := make(map[string]string)
    for _, key := range allowedKeys {
        val := r.URL.Query().Get(key)
        if val != "" {
            result[key] = val
        }
    }
    return result
}
```

### 12.4 Вспомогательная функция

```go
func contains(list []string, value string) bool {
    for _, item := range list {
        if item == value {
            return true
        }
    }
    return false
}
```

> **Безопасность:** каждый хэндлер передаёт whitelist допустимых полей — защита от SQL-инъекций через произвольные имена колонок.

---

## 13. Backend: обновление хэндлеров

### 13.0 Общий тип `ListQuery` (shared/httpx)

Чтобы не нарушать правило ≤5 параметров, все list-параметры передаются в одном struct:

```go
// httpx/request.go

type ListQuery struct {
    Page    int
    PerPage int
    SortBy  string
    SortDir string
    Search  string
    Filters map[string]string
}

// ParseListQuery собирает все query-параметры в единый объект.
func ParseListQuery(r *http.Request, sortCfg SortConfig, filterKeys []string) ListQuery {
    page, perPage := ParsePagination(r)
    sortBy, sortDir := ParseSort(r, sortCfg.Allowed, sortCfg.Fallback)
    search := ParseSearch(r)
    filters := ParseFilters(r, filterKeys)
    return ListQuery{page, perPage, sortBy, sortDir, search, filters}
}

type SortConfig struct {
    Allowed  []string
    Fallback string
}
```

### 13.1 Products handler (эталонная реализация)

```go
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    tenantID := tenant.FromContext(r.Context())
    q := httpx.ParseListQuery(r, httpx.SortConfig{
        Allowed: []string{"name", "created_at"}, Fallback: "created_at",
    }, []string{"is_composite"})

    data, total, err := h.service.List(r.Context(), tenantID, q)
    if err != nil {
        httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to list products", err.Error())
        return
    }
    meta := &httpx.Meta{Page: q.Page, PerPage: q.PerPage, Total: total}
    httpx.WriteData(w, http.StatusOK, data, meta)
}
```

> **Сигнатура service:** `List(ctx, tenantID, q httpx.ListQuery) → ([]Product, int, error)` — 3 параметра, правило ≤5 соблюдено.

### 13.2 Products repository (динамический SQL)

Обновить `load` и `count`, принимая `httpx.ListQuery`:
- `q.Search` → добавить `AND name ILIKE '%' || $N || '%'` (если search не пустой)
- `q.SortBy` / `q.SortDir` → заменить `ORDER BY created_at DESC` на `ORDER BY <SortBy> <SortDir>`
- `q.Filters["is_composite"]` → сохранить текущую логику `AND is_composite=$N`

### 13.3 Documents handler (аналогично)

```go
q := httpx.ParseListQuery(r, httpx.SortConfig{
    Allowed: []string{"date", "number", "total_cost"}, Fallback: "date",
}, []string{"type", "status"})
```
- `q.Search` → поиск по `number`

### 13.4 Documents: обновить модель и SQL для `total_cost`

Текущий `ListItem` не содержит `total_cost`, а SQL `load()` не считает сумму. Необходимо:

1. Добавить поле в модель:
   ```go
   type ListItem struct {
       ID        string `json:"id"`
       Type      string `json:"type"`
       Date      string `json:"date"`
       Number    string `json:"number"`
       Status    string `json:"status"`
       Note      string `json:"note"`
       TotalCost string `json:"total_cost"`
   }
   ```

2. Обновить SQL в `load()` — добавить подзапрос суммы:
   ```sql
   SELECT d.id::text, d.type, d.date::text, COALESCE(d.number,''), d.status, COALESCE(d.note,''),
          COALESCE((SELECT SUM(di.qty * di.unit_price) FROM documents.document_items di WHERE di.document_id = d.id), 0)::text
   FROM documents.documents d WHERE d.tenant_id=$1
   ```

3. Обновить `scanList` — добавить `Scan` для `TotalCost`

4. Текущая сортировка `ORDER BY date DESC, created_at DESC` → заменить на динамическую через `q.SortBy`/`q.SortDir` (fallback: `date DESC`)

---

## 14. Правила кода (ОБЯЗАТЕЛЬНО)

1. **≤5 параметров** у функции
2. **≤300 строк** на файл
3. **≤20 строк** на метод/функцию
4. **≤3 уровня** вложенности
5. **Один компонент — один файл**
6. **Feature-based** структура
7. **Никаких тестов**, docs-файлов, `.md` если не просят
8. **JSX**, не TypeScript

---

## 15. Порядок внедрения

### Этап A — Foundation (server mode)

| Шаг | Что делать | Где |
|---|---|---|
| 1 | Добавить `getClientToken/getAdminToken` в `session.js` | frontend |
| 2 | Добавить `getJSON` в `http.js` (через `createURL`) | frontend |
| 3 | ~~Мигрировать `client.ts` → `client.js`~~ — **уже выполнено** | frontend |
| 4 | Добавить `ParseSort`, `ParseSearch`, `ParseFilters`, `ListQuery`, `ParseListQuery` в `request.go` | backend |
| 5 | Обновить Products handler + service + repository (sort/search/filters через `ListQuery`) | backend |
| 6 | Обновить Documents handler + service + repository (sort/search/filters + `total_cost` в модели и SQL) | backend |
| 7 | Добавить `"@tanstack/react-table": "^8.21.3"` в `dependencies` в `package.json` (без `npm install`) | frontend |
| 8 | Создать `shared/model/data-table/` (tableState, toQueryParams, createTablePreset, useDebounce, useServerDataTable) | frontend |
| 9 | Создать `shared/ui/data-table/` (DataTable + все подкомпоненты + CSS) | frontend |
| 10 | Создать `productsTablePreset.js`, подключить ProductsPage | frontend |
| 11 | Создать `documentsTablePreset.js`, подключить DocumentsPage | frontend |

> **Проверка** (`npm run dev`, `npm run build`) выполняется **на сервере при деплое**, не локально.

### Этап B — Масштабирование

| Шаг | Что делать |
|---|---|
| 1 | Перенести AdminTenants и Variants на preset-подход |
| 2 | Backend: привести `GET /api/client/warehouses` к list-контракту (page, per_page, meta.total) |
| 3 | Backend: добавить list-endpoints для отчётов (новые URL, не ломая текущие агрегатные) |
| 4 | Подключить Warehouses и Reports к DataTable |
| 5 | Включить URL-sync состояния таблицы (query params роутера) для deep-link |
| 6 | Добавить виртуализацию (`@tanstack/react-virtual` — добавить в `package.json`) по необходимости |
| 7 | Включить сортировку/поиск в presets после backend-готовности |

---

## 16. Backend-блокеры для Этапа B

1. `GET /api/client/warehouses` — привести к list-контракту: `page`, `per_page`, `meta.total`, envelope `{ data, meta, error }`
2. Для отчётов в DataTable — добавить отдельные list-endpoints (не ломая текущие агрегатные `reports/*`) с контрактом `{ data: Array<Row>, meta: { page, per_page, total } }`
3. До выполнения пунктов выше Warehouses и Reports **не переводить** на DataTable

---

## 17. Производительность

### 17.1 Правила
- Production → всегда `server` mode
- Более 200 видимых строк → виртуализация

### 17.2 Виртуализация (опционально)

Добавить `"@tanstack/react-virtual": "^3.x"` в `dependencies` в `package.json` (без `npm install`).
- `estimateSize: 40`, `overscan: 8`
- Отдельный scroll-container с фиксированной высотой

### 17.3 UX
- Во время загрузки **не очищать** предыдущие строки (keep previous data)
- При ошибке — встроенный блок ошибки + кнопка «Повторить»
- Отменять предыдущий запрос через `AbortController` при новом state

---

## 18. Чеклист проверки

### Этап A
- [ ] Таблица рендерится на `/dashboard/products` и `/dashboard/documents` с реальным API
- [ ] Пагинация, фильтры работают согласно `preset.capabilities`
- [ ] Responsive, hover, зебра работают
- [ ] Нет синтаксических ошибок в коде (проверка при деплое на сервере)
- [ ] `getJSON` работает с query params и AbortSignal, не дублирует base URL
- [ ] В `session.js` есть `getClientToken/getAdminToken`
- [ ] В page-компонентах нет дублирования table-state логики
- [ ] Backend: `ParseSort/ParseSearch/ParseFilters` добавлены в `request.go`
- [ ] Backend: Products и Documents handlers используют `ParseListQuery`
- [ ] Backend: Documents `ListItem` содержит `total_cost`, SQL считает сумму items
- [ ] Backend: Documents сортировка по `date` (не `created_at`) как fallback
- [ ] Backend: SQL-запросы валидируют поля по whitelist
- [ ] Запросы отменяются при быстрых изменениях фильтров (AbortController)
- [ ] Debounce 300ms через `useDebounce.js`
- [ ] Смена фильтра → pageIndex сбрасывается на 0
- [ ] Controlled handlers корректно обрабатывают updater (значение/функция)
- [ ] `getRowId` задан из `preset.rowId`
- [ ] Таблица не мигает пустотой при перезагрузке (keep previous)
- [ ] Ошибка API отображается с кнопкой «Повторить»
- [ ] Expand работает для товаров с вариантами (lazy-fetch + кеш)

### Этап B
- [ ] AdminTenants и Variants на preset-подходе
- [ ] `GET /api/client/warehouses` поддерживает list-контракт
- [ ] Report-list endpoints возвращают `{ data, meta, error }`
- [ ] Warehouses и Reports подключены к DataTable только после backend-готовности
