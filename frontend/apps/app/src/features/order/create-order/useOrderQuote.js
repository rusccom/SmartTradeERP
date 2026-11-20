import { useEffect, useState } from "react";
import { postJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";
import { lineInputs } from "./entryModel";

export default function useOrderQuote(form, card) {
  const [result, setResult] = useState({});
  const body = JSON.stringify({ items: lineInputs(form.items), revision: form.revision, delivery_cost: form.delivery_cost || null });
  const path = apiPaths.orderQuote(form.id);
  const enabled = !card?.invoice && form.items.length > 0 && form.items.every((l) => Number(l.qty) > 0);
  useEffect(() => {
    if (!enabled) return;
    const controller = new AbortController();
    const publish = (v) => { if (!controller.signal.aborted) setResult({ path, body, ...v }); };
    const timer = setTimeout(() => postJSON(path, JSON.parse(body), controller.signal)
      .then((data) => publish({ data })).catch((e) => publish({ error: e.message })), 180);
    return () => { clearTimeout(timer); controller.abort(); };
  }, [enabled, body, path]);
  if (card?.invoice) return { data: card };
  return enabled && result.body === body && result.path === path ? result : { loading: enabled };
}
