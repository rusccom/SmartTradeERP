import { useEffect, useState } from "react";
import { postJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";
import { lineInputs } from "./entryModel";
import { destination } from "./orderVat";

export default function useOrderQuote(form, card) {
  const [result, setResult] = useState({});
  const body = quoteRequest(form);
  const saved = card && body === quoteRequest(card);
  const path = apiPaths.orderQuote(form.id);
  const enabled = !saved && !card?.invoice && form.items.length > 0 && form.items.every((l) => Number(l.qty) > 0);
  useEffect(() => {
    if (!enabled) return;
    const controller = new AbortController();
    const publish = (v) => { if (!controller.signal.aborted) setResult({ path, body, ...v }); };
    const timer = setTimeout(() => postJSON(path, JSON.parse(body), controller.signal)
      .then((data) => publish({ data })).catch((e) => publish({ error: e.message })), 180);
    return () => { clearTimeout(timer); controller.abort(); };
  }, [enabled, body, path]);
  if (saved || card?.invoice) return { data: card };
  return enabled && result.body === body && result.path === path ? result : { loading: enabled };
}

function quoteRequest(order) {
  const delivery = order.delivery_cost === "" ? null : order.delivery_cost;
  const pickup = order.details?.delivery_kind === "pickup";
  const details = { delivery_kind: pickup ? "pickup" : "address", [pickup ? "pickup" : "shipping"]: { country: destination(order) } };
  return JSON.stringify({ items: lineInputs(order.items), revision: order.revision, delivery_cost: delivery, details });
}
