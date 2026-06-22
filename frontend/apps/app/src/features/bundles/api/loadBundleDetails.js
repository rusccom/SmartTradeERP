import { apiPaths } from "@smarterp/api/publicApi";
import { getJSON } from "@smarterp/api/http";

export async function loadBundleDetails(id, signal) {
  const { data } = await getJSON(apiPaths.bundleById(id), {}, signal);
  return data || null;
}
