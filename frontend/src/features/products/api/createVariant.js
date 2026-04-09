import { apiPaths } from "../../../shared/api/client";
import { postJSON } from "../../../shared/api/http";

export function createVariant(payload) {
  return postJSON(apiPaths.variants, payload);
}
