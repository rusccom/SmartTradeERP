import { apiPaths } from "../../../shared/api/publicApi";
import { getJSON } from "../../../shared/api/http";

export function loadProductVariants(row) {
  return getJSON(apiPaths.variants, { product_id: row.id }).then(readVariants);
}

function readVariants(response) {
  return response.data || [];
}
