const EMPTY_SORTING = Object.freeze({ columns: Object.freeze({}), enabled: false });

export function createTableSortingConfig(config) {
  const entries = normalizeEntries(config);
  if (entries.length === 0) {
    return EMPTY_SORTING;
  }
  return Object.freeze({ columns: Object.freeze(Object.fromEntries(entries)), enabled: true });
}

export function isColumnSortable(sorting, column) {
  if (column.enableSorting === false) {
    return false;
  }
  return Boolean(readColumnSortKey(sorting, column));
}

export function readSortQuery(sorting, item) {
  if (!item) {
    return null;
  }
  const sortBy = sorting?.columns?.[item.id] || "";
  return sortBy ? { sortBy, sortDir: item.desc ? "desc" : "asc" } : null;
}

function normalizeEntries(config) {
  const items = Array.isArray(config) ? config : config?.columns || [];
  return items.map(normalizeItem).filter(Boolean);
}

function normalizeItem(item) {
  if (typeof item === "string") {
    return [item, item];
  }
  if (!item?.id) {
    return null;
  }
  return [item.id, item.sortKey || item.id];
}

function readColumnSortKey(sorting, column) {
  const columnID = column.id || column.accessorKey;
  return columnID ? sorting?.columns?.[columnID] : "";
}
