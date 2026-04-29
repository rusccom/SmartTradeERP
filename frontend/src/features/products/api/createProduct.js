import { apiPaths } from "../../../shared/api/publicApi";
import { postJSON } from "../../../shared/api/http";

export function createProduct(payload) {
  return postJSON(apiPaths.products, payload);
}
