import { postJSON } from "../../../shared/api/http";
import { apiPaths } from "../../../shared/api/publicApi";

export function createCurrency(payload) {
  return postJSON(apiPaths.currencies, payload);
}
