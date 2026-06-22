import { apiPaths } from "../../../shared/api/publicApi";
import { putJSON } from "../../../shared/api/http";

export function updateProduct(id, payload) {
  return putJSON(apiPaths.productById(id), { ...payload, is_composite: false });
}
