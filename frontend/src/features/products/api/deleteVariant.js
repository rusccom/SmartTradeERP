import { deleteJSON } from "../../../shared/api/http";
import { apiPaths } from "../../../shared/api/publicApi";

export function deleteVariant(id) {
  return deleteJSON(apiPaths.variantById(id));
}
