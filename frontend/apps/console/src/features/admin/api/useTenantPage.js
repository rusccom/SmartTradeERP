import { useEffect, useState } from "react";
import { getJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export function useTenantPage(page) {
  const [state, setState] = useState({ items: [], total: 0, error: "", loading: true });
  useEffect(() => {
    const controller = new AbortController();
    setState({ items: [], total: 0, error: "", loading: true });
    getJSON(apiPaths.adminTenants, { page, per_page: 50 }, controller.signal)
      .then(({ data, meta }) => {
        if (!controller.signal.aborted) setState({ items: data ?? [], total: meta?.total ?? 0, error: "", loading: false });
      })
      .catch(error => {
        if (!controller.signal.aborted) setState({ items: [], total: 0, error: error.message, loading: false });
      });
    return () => controller.abort();
  }, [page]);
  return state;
}
