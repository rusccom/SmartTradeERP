import { getJSON, putJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export async function loadConnection(signal) {
  const { data } = await getJSON(apiPaths.inpostSettings, {}, signal);
  return data;
}

export function verifyAndSave(payload) {
  return putJSON(apiPaths.inpostSettings, payload);
}
