import { useEffect, useMemo, useRef, useState } from "react";

import { applyStateChange } from "./tableState";
import { toQueryParams } from "./toQueryParams";
import { useDebounce } from "./useDebounce";

export function useServerDataTable(preset) {
  const table = useTableState(preset.defaultState);
  const fetchState = useFetchState();
  const requestState = useRequestState(table.state);
  useFetchEffect({ preset, requestState, retryToken: table.retryToken, ...fetchState });
  return buildHookResult({ ...fetchState, retry: table.retry, state: table.state, tableHandlers: table.handlers });
}

function useTableState(defaultState) {
  const [state, setState] = useState(() => defaultState);
  const [retryToken, setRetryToken] = useState(0);
  const handlers = useMemo(() => createTableHandlers(setState), []);
  const retry = () => setRetryToken((value) => value + 1);
  return { state, handlers, retry, retryToken };
}

function useFetchState() {
  const [data, setData] = useState([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  return { data, total, loading, error, setData, setTotal, setLoading, setError };
}

function useRequestState(state) {
  const debouncedGlobalFilter = useDebounce(state.globalFilter, 300);
  return useMemo(
    () => createQueryState(state, debouncedGlobalFilter),
    [state.pagination, state.sorting, debouncedGlobalFilter],
  );
}

function useFetchEffect(params) {
  const { preset, requestState, retryToken, setData, setTotal, setLoading, setError } = params;
  const abortRef = useRef(null);
  useEffect(() => {
    const controller = startRequest(abortRef);
    setLoading(true);
    runFetch(preset, requestState, controller.signal)
      .then((result) => handleFetchSuccess(result, { setData, setTotal, setError }))
      .catch((reason) => handleFetchError(reason, setError))
      .finally(() => finishRequest(controller, setLoading));
    return () => controller.abort();
  }, [preset, requestState, retryToken, setData, setTotal, setLoading, setError]);
}

function createQueryState(state, globalFilter) {
  return { ...state, globalFilter };
}

function createTableHandlers(setState) {
  return {
    onSortingChange: (updater) => patchState(setState, "sorting", updater),
    onGlobalFilterChange: (updater) => patchState(setState, "globalFilter", updater),
    onPaginationChange: (updater) => patchState(setState, "pagination", updater),
  };
}

function patchState(setState, key, updater) {
  setState((prev) => {
    const value = resolveUpdater(updater, prev[key]);
    return applyStateChange(prev, { [key]: value });
  });
}

function resolveUpdater(updater, currentValue) {
  return updater instanceof Function ? updater(currentValue) : updater;
}

function startRequest(abortRef) {
  abortRef.current?.abort();
  const controller = new AbortController();
  abortRef.current = controller;
  return controller;
}

async function runFetch(preset, state, signal) {
  const query = toQueryParams(state, preset);
  return preset.fetchPage({ query, signal });
}

function handleFetchSuccess(result, state) {
  state.setData(Array.isArray(result?.rows) ? result.rows : []);
  state.setTotal(Number(result?.total) || 0);
  state.setError(null);
}

function handleFetchError(error, setError) {
  if (error?.name === "AbortError") {
    return;
  }
  setError(readErrorMessage(error));
}

function finishRequest(controller, setLoading) {
  if (!controller.signal.aborted) {
    setLoading(false);
  }
}

function readErrorMessage(error) {
  return error?.message || "Failed to load table data";
}

function buildHookResult({ data, total, loading, error, retry, state, tableHandlers }) {
  return {
    data,
    total,
    loading,
    error,
    retry,
    tableState: {
      sorting: state.sorting,
      globalFilter: state.globalFilter,
      pagination: state.pagination,
      ...tableHandlers,
    },
  };
}
