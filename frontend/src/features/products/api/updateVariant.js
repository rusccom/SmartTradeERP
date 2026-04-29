import { apiPaths } from "../../../shared/api/publicApi";
import { putJSON } from "../../../shared/api/http";

export function updateVariant(id, payload) {
  return putJSON(apiPaths.variantById(id), payload);
}
