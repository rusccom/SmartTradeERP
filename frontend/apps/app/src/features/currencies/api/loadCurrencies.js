import { getJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export async function loadCurrencies(signal) {
  const { data } = await getJSON(apiPaths.currencies, { per_page: 100 }, signal);
  return Array.isArray(data) ? data : [];
}
