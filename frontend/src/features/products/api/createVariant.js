import { apiPaths } from "../../../shared/api/publicApi";
import { postJSON } from "../../../shared/api/http";

export function createVariant(payload) {
  return postJSON(apiPaths.variants, payload);
}
