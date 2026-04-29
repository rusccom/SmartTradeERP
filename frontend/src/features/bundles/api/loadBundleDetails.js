import { apiPaths } from "../../../shared/api/client";
import { getJSON } from "../../../shared/api/http";

export async function loadBundleDetails(id, signal) {
  const { data } = await getJSON(apiPaths.bundleById(id), {}, signal);
  return data || null;
}
