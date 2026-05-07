import { putJSON } from "../../../shared/api/http";
import { apiPaths } from "../../../shared/api/publicApi";

export function setBaseCurrency(payload) {
  return putJSON(apiPaths.currencyBase, payload);
}
