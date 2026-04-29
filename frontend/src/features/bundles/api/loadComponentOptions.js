import { apiPaths } from "../../../shared/api/client";
import { getJSON } from "../../../shared/api/http";

export async function loadComponentOptions(query = {}) {
  const params = componentParams(query);
  const { data, meta } = await getJSON(apiPaths.products, params, query.signal);
  return { options: flattenVariants(data || []), meta };
}

function componentParams(query) {
  return {
    include: "variants",
    is_composite: "false",
    page: query.page || 1,
    per_page: query.perPage || 20,
    search: query.search || "",
  };
}

function flattenVariants(products) {
  return products.flatMap((product) => readProductVariants(product));
}

function readProductVariants(product) {
  return (product.variants || []).map((variant) => toOption(product, variant));
}

function toOption(product, variant) {
  return {
    id: variant.id,
    label: variantLabel(product.name, variant.name),
    unit: variant.unit,
  };
}

function variantLabel(productName, variantName) {
  if (!variantName || variantName === "Default") return productName;
  return `${productName} / ${variantName}`;
}
