import { postJSON } from "../../../shared/api/http";
import { apiPaths } from "../../../shared/api/publicApi";

export function createBundle(payload) {
  return postJSON(apiPaths.products, { ...payload, is_composite: true });
}
