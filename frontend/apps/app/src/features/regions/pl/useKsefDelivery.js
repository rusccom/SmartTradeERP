import { useCallback, useEffect, useRef, useState } from "react";
import { getJSON, postJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export default function useKsefDelivery(id) {
  const [data, setData] = useState(null);
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const active = useRef(false);
  const mounted = useRef(false);
  const run = useCallback((action, environment) => perform(action, environment, id, {
    active, mounted, setData, setError, setBusy,
  }), [id]);
  useEffect(() => { mounted.current = true; run("read"); return () => { mounted.current = false; }; }, [run]);
  useEffect(() => {
    if (data?.state !== "pending" && data?.state !== "sending" || busy || error) return;
    const timer = setTimeout(() => run("check"), 12000);
    return () => clearTimeout(timer);
  }, [data, busy, error, run]);
  return { data, error, busy, run };
}

async function perform(action, environment, id, s) {
  if (s.active.current) return;
  s.active.current = true;
  s.setBusy(true);
  s.setError("");
  try {
    const data = await request(action, id, environment);
    if (s.mounted.current) s.setData(data);
  } catch (e) {
    if (s.mounted.current) s.setError(e.message);
  } finally {
    s.active.current = false;
    if (s.mounted.current) s.setBusy(false);
  }
}

async function request(action, id, environment) {
  const path = apiPaths.invoiceDelivery(id);
  if (action === "send") return postJSON(path, { confirmed: true, environment });
  if (action === "check") return postJSON(apiPaths.checkInvoiceDelivery(id), {});
  return (await getJSON(path)).data;
}
