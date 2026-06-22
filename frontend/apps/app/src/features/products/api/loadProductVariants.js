import { apiPaths } from "@smarterp/api/publicApi";
import { getJSON } from "@smarterp/api/http";

export function loadProductVariants(row) {
  if (Array.isArray(row.variants)) {
    return Promise.resolve(withParentProduct(row.variants, row));
  }
  return getJSON(apiPaths.variants, { product_id: row.id }).then((response) => withParentProduct(readVariants(response), row));
}

function readVariants(response) {
  return response.data || [];
}

function withParentProduct(variants, product) {
  return variants.map((variant) => ({ ...variant, product }));
}
