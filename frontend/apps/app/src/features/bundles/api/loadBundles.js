import { apiPaths } from "@smarterp/api/publicApi";
import { getJSON } from "@smarterp/api/http";

export async function loadBundles(signal) {
  const params = { page: 1, per_page: 100 };
  const { data, meta } = await getJSON(apiPaths.bundles, params, signal);
  return { bundles: data || [], total: meta?.total || 0 };
}
