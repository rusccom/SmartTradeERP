import { getJSON, putJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export async function loadSettings(id, signal) {
  const response = await getJSON(apiPaths.bundleStorefront(id), {}, signal);
  return response.data;
}

export function saveSettings(id, settings) {
  return putJSON(apiPaths.bundleStorefront(id), settings);
}
