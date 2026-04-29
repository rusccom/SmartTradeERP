import { apiPaths } from "../../../shared/api/client";
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
    { accessorKey: "name", header: t("products.columns.name"), enableFilter: true },
    { accessorKey: "is_composite", header: t("products.columns.composite"), enableSorting: false, filterVariant: "select", filterOptions: createBooleanOptions(t), cell: (value) => readBooleanLabel(value, t) },
    { accessorKey: "updated_at", header: t("products.columns.updatedAt"), enableSorting: false },
  ];
}

function createBooleanOptions(t) {
  return [{ value: "true", label: t("common.yes") }, { value: "false", label: t("common.no") }];
}

function readBooleanLabel(value, t) {
  return value ? t("common.yes") : t("common.no");
}

function readProductId(row) {
  return row.id;
}
