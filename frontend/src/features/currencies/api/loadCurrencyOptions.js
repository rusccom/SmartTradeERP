import { getJSON } from "../../../shared/api/http";
import { apiPaths } from "../../../shared/api/publicApi";

export async function loadCurrencyOptions(signal) {
  const { data } = await getJSON(apiPaths.currencyOptions, { per_page: 100 }, signal);
  return Array.isArray(data) ? data : [];
}
