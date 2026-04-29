import { apiPaths } from "../../../shared/api/publicApi";
import { createApiTablePreset } from "../../../shared/model/data-table/createApiTablePreset";

const capabilities = { sorting: true, search: true };

export function createProductsTablePreset(t) {
  return createApiTablePreset({
    id: "products",
    path: apiPaths.products,
    rowId: readProductId,
    columns: createColumns(t),
    capabilities,
  });
}

function createColumns(t) {
  return [
    { accessorKey: "name", header: t("products.columns.name") },
  ];
}

function readProductId(row) {
  return row.id;
}
