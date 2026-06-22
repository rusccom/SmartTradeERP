import { apiPaths } from "@smarterp/api/publicApi";
import { postJSON } from "@smarterp/api/http";

export function createVariant(payload) {
  return postJSON(apiPaths.variants, payload);
}
