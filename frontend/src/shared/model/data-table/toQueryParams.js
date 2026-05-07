import { serializeTableSearchFilter } from "./filters/tableSearchFilter";
import { readSortQuery } from "./tableSorting";

const RESERVED_KEYS = new Set(["page", "per_page", "sort_by", "sort_dir", "search"]);

export function toQueryParams(state, preset) {
  const params = readBaseParams(state);
  appendSorting(params, state, preset);
  appendSearch(params, state, preset);
  const custom = readCustomParams(state, preset);
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
  const sorting = readSortQuery(preset.sortingConfig, state.sorting[0]);
  if (!sorting) {
    return;
  }
  params.sort_by = sorting.sortBy;
  params.sort_dir = sorting.sortDir;
}

function appendSearch(params, state, preset) {
  if (!isSearchEnabled(preset)) {
    return;
  }
  const value = serializeTableSearchFilter(preset.search, state.globalFilter);
  if (value === "") {
    return;
  }
  params[readSearchQueryKey(preset)] = value;
}

function isSearchEnabled(preset) {
  return preset.capabilities.search === true && preset.search?.enabled !== false;
}

function readSearchQueryKey(preset) {
  return preset.search?.queryKey || "search";
}

function readCustomParams(state, preset) {
  const source = preset.mapStateToQuery ? preset.mapStateToQuery(state) : {};
  const reserved = readReservedKeys(preset);
  return Object.entries(source || {}).reduce((acc, [key, value]) => {
    if (!reserved.has(key)) {
      acc[key] = value;
    }
    return acc;
  }, {});
}

function readReservedKeys(preset) {
  return new Set([...RESERVED_KEYS, readSearchQueryKey(preset)]);
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
