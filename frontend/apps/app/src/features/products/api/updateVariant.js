import { apiPaths } from "@smarterp/api/publicApi";
import { putJSON } from "@smarterp/api/http";

export function updateVariant(id, payload) {
  return putJSON(apiPaths.variantById(id), payload);
}
