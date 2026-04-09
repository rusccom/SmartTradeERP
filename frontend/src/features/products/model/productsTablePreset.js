import { apiPaths } from "../../../shared/api/client";
import { getJSON } from "../../../shared/api/http";
import { createTablePreset } from "../../../shared/model/data-table/createTablePreset";

const defaultState = {
  pagination: { pageIndex: 0, pageSize: 20 },
  sorting: [],
  globalFilter: "",
  columnFilters: [],
};

const capabilities = { sorting: true, search: true };

export function createProductsTablePreset(t) {
  return createTablePreset({
    id: "products",
    rowId: readProductId,
    columns: createColumns(t),
    defaultState,
    capabilities,
    fetchPage: fetchProductsPage,
  });
}

function createColumns(t) {
  return [
    { accessorKey: "name", header: t("products.columns.name"), enableFilter: true },
    { accessorKey: "is_composite", header: t("products.columns.composite"), enableSorting: false, filterVariant: "select", filterOptions: createBooleanOptions(t), cell: (value) => readBooleanLabel(value, t) },
    { accessorKey: "updated_at", header: t("products.columns.updatedAt"), enableSorting: false },
  ];
}

async function fetchProductsPage({ query, signal }) {
  const { data, meta } = await getJSON(apiPaths.products, query, signal);
  return { rows: data || [], total: meta?.total || 0 };
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
