import { useEffect, useState } from "react";
import { getJSON } from "@smarterp/api/http";

export default function useApiResource(path, refreshOnFocus = false) {
  const [result, setResult] = useState({});
  const [attempt, setAttempt] = useState(0);
  useEffect(() => {
    if (!refreshOnFocus) return;
    const reload = () => setAttempt((n) => n + 1);
    window.addEventListener("focus", reload);
    return () => window.removeEventListener("focus", reload);
  }, [refreshOnFocus]);
  useEffect(() => {
    if (!path) return;
    const controller = new AbortController();
    setResult({ path, loading: true });
    getJSON(path, undefined, controller.signal)
      .then(({ data }) => { if (!controller.signal.aborted) setResult({ path, data }); })
      .catch((e) => { if (!controller.signal.aborted) setResult({ path, error: e.message }); });
    return () => controller.abort();
  }, [path, attempt]);
  return { ...(result.path === path ? result : { loading: Boolean(path) }), retry: () => setAttempt((n) => n + 1) };
}
