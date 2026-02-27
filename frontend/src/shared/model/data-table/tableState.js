export const FALLBACK_STATE = {
  pagination: { pageIndex: 0, pageSize: 20 },
  sorting: [],
  globalFilter: "",
  columnFilters: [],
};

const PAGE_RESET_KEYS = ["sorting", "globalFilter", "columnFilters"];

export function applyStateChange(prev, patch) {
  const next = { ...prev, ...patch };
  return shouldResetPage(patch) ? resetPage(next) : next;
}

function shouldResetPage(patch) {
  return PAGE_RESET_KEYS.some((key) => Object.prototype.hasOwnProperty.call(patch, key));
}

function resetPage(state) {
  return {
    ...state,
    pagination: { ...state.pagination, pageIndex: 0 },
  };
}
