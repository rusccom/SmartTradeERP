import { deleteJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export function deleteVariant(id) {
  return deleteJSON(apiPaths.variantById(id));
}
