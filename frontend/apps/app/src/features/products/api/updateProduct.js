import { apiPaths } from "@smarterp/api/publicApi";
import { putJSON } from "@smarterp/api/http";

export function updateProduct(id, payload) {
  return putJSON(apiPaths.productById(id), { ...payload, is_composite: false });
}
