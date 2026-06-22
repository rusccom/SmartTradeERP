import { postJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export function createBundle(payload) {
  return postJSON(apiPaths.products, { ...payload, is_composite: true });
}
