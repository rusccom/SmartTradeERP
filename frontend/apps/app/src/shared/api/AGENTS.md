# AGENTS.md

Локальные правила для frontend API-слоя.

- Фронтенд работает только с публичным backend API: `/api/client/*` и `/api/admin/*`.
- Все публичные пути добавляй в `publicApi.js`; feature-код импортирует `apiPaths` отсюда.
- Не пиши строковые `/api/...` в feature-папках и не обходи `shared/api/http.js`.
- Не связывай frontend с `backend/internal/features/*`: это внутренняя реализация за `backend/internal/api`.
- Если нужен новый endpoint, сначала добавь его в `backend/internal/api`, обнови публичный контракт, затем подключай через `publicApi.js`.
