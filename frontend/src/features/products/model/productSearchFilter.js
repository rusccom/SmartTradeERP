import { createTableSearchFilter } from "../../../shared/model/data-table/filters/tableSearchFilter";

export const productSearchFilter = createTableSearchFilter({
  id: "products",
  placeholderKey: "products.search.placeholder",
  queryKey: "search",
  serialize: serializeProductSearch,
});

function serializeProductSearch(value) {
  return String(value || "").trim();
}
