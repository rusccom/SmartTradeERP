import { useEffect, useState } from "react";
import { getJSON, putJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export function usePlugins(tenantId = "") {
  const path = tenantId ? apiPaths.adminTenantPlugins(tenantId) : apiPaths.clientPlugins;
  const [access, setAccess] = useState(null);
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  useEffect(() => loadAccess(path, setAccess, setError), [path]);
  const save = (plugin, checked) => saveAccess(
    { tenantId, path, plugin, checked }, access, { setAccess, setError, setBusy },
  );
  return { access, busy, error, save };
}

async function saveAccess(change, previous, { setAccess, setError, setBusy }) {
  const { tenantId, path, plugin, checked } = change;
  const updated = tenantId ? { ...plugin, granted: checked, enabled: checked && plugin.enabled } : { ...plugin, enabled: checked };
  setAccess({ ...previous, plugins: previous.plugins.map(item => item.key === plugin.key ? updated : item) });
  setBusy(true);
  setError("");
  let saved = false;
  try {
    await changeAccess(tenantId, plugin.key, checked);
    saved = true;
    setAccess((await getJSON(path)).data);
  } catch (error) {
    if (!saved) setAccess(previous);
    setError(error.message);
  } finally {
    setBusy(false);
  }
}

function loadAccess(path, setAccess, setError) {
  const controller = new AbortController();
  setAccess(null);
  setError("");
  getJSON(path, undefined, controller.signal)
    .then(({ data }) => { if (!controller.signal.aborted) setAccess(data); })
    .catch(error => { if (!controller.signal.aborted) setError(error.message); });
  return () => controller.abort();
}

function changeAccess(tenantId, key, checked) {
  if (tenantId) return putJSON(apiPaths.adminTenantPlugin(tenantId, key), { granted: checked });
  return putJSON(apiPaths.clientPlugin(key), { enabled: checked });
}
