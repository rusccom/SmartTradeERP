const DEFAULT_SEARCH_FILTER = {
  debounceMs: 300,
  enabled: true,
  placeholderKey: "dataTable.searchPlaceholder",
  queryKey: "search",
  serialize: defaultSerialize,
};

export function createTableSearchFilter(config) {
  const normalized = normalizeSearchFilterConfig(config);
  return Object.freeze({ ...DEFAULT_SEARCH_FILTER, ...normalized });
}

export function serializeTableSearchFilter(filter, value) {
  const serializer = readSerializer(filter);
  return serializer(value);
}

export function readTableSearchDebounce(filter) {
  const delay = Number(filter?.debounceMs);
  return Number.isFinite(delay) && delay >= 0 ? delay : DEFAULT_SEARCH_FILTER.debounceMs;
}

function normalizeSearchFilterConfig(config) {
  if (typeof config === "string") {
    return { queryKey: config };
  }
  return config || {};
}

function readSerializer(filter) {
  return typeof filter?.serialize === "function" ? filter.serialize : DEFAULT_SEARCH_FILTER.serialize;
}

function defaultSerialize(value) {
  return String(value || "").trim();
}
