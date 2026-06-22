import { apiPaths } from "@smarterp/api/publicApi";
import { postJSON } from "@smarterp/api/http";

export function createProduct(payload) {
  return postJSON(apiPaths.products, { ...payload, is_composite: false });
}
