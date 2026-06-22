import { putJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export function setBaseCurrency(payload) {
  return putJSON(apiPaths.currencyBase, payload);
}
