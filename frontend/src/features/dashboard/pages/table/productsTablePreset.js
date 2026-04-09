import { apiPaths } from "../../../../shared/api/client";
import { getJSON } from "../../../../shared/api/http";
import { createTablePreset } from "../../../../shared/model/data-table/createTablePreset";

export function createProductsTablePreset(t) {
  return createTablePreset({
    id: "products",
    rowId: (row) => row.id,
    columns: [
      { accessorKey: "name", header: t("products.columns.name"), enableFilter: true },
      {
        accessorKey: "is_composite",
        header: t("products.columns.composite"),
        enableSorting: false,
        filterVariant: "select",
        filterOptions: [
          { value: "true", label: t("common.yes") },
          { value: "false", label: t("common.no") },
        ],
        cell: (value) => (value ? t("common.yes") : t("common.no")),
      },
      { accessorKey: "updated_at", header: t("products.columns.updatedAt"), enableSorting: false },
    ],
    defaultState: {
      pagination: { pageIndex: 0, pageSize: 20 },
      sorting: [],
      globalFilter: "",
      columnFilters: [],
    },
    capabilities: { sorting: true, search: true },
    fetchPage: async ({ query, signal }) => {
      const { data, meta } = await getJSON(apiPaths.products, query, signal);
      return { rows: data || [], total: meta?.total || 0 };
    },
  });
}
