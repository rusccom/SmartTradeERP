const RESERVED_KEYS = new Set(["page", "per_page", "sort_by", "sort_dir", "search"]);

export function toQueryParams(state, preset) {
  const params = readBaseParams(state);
  appendSorting(params, state, preset);
  appendSearch(params, state, preset);
  const custom = readCustomParams(state, preset);
  const overrides = new Set(Object.keys(custom));
  appendColumnFilters(params, state.columnFilters, overrides);
  return cleanParams({ ...params, ...custom });
}

function readBaseParams(state) {
  return {
    page: state.pagination.pageIndex + 1,
    per_page: state.pagination.pageSize,
  };
}

function appendSorting(params, state, preset) {
  if (preset.capabilities.sorting !== true || state.sorting.length === 0) {
    return;
  }
  params.sort_by = state.sorting[0].id;
  params.sort_dir = state.sorting[0].desc ? "desc" : "asc";
}

function appendSearch(params, state, preset) {
  const value = state.globalFilter.trim();
  if (preset.capabilities.search !== true || value === "") {
    return;
  }
  params.search = value;
}

function readCustomParams(state, preset) {
  const source = preset.mapStateToQuery ? preset.mapStateToQuery(state) : {};
  return Object.entries(source || {}).reduce((acc, [key, value]) => {
    if (!RESERVED_KEYS.has(key)) {
      acc[key] = value;
    }
    return acc;
  }, {});
}

function appendColumnFilters(params, columnFilters, overrides) {
  columnFilters.forEach((item) => {
    if (!item || overrides.has(item.id)) {
      return;
    }
    params[item.id] = normalizeFilterValue(item.value);
  });
}

function normalizeFilterValue(value) {
  if (Array.isArray(value)) {
    return value.filter((item) => item !== "" && item != null).join(",");
  }
  return value;
}

function cleanParams(params) {
  return Object.entries(params).reduce((acc, [key, value]) => {
    if (value === undefined || value === null || value === "") {
      return acc;
    }
    acc[key] = value;
    return acc;
  }, {});
}
