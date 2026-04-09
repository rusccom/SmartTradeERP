import { apiPaths } from "../../../../shared/api/client";
import { getJSON } from "../../../../shared/api/http";
import { createTablePreset } from "../../../../shared/model/data-table/createTablePreset";

export function createCustomersTablePreset(t) {
  return createTablePreset({
    id: "customers",
    rowId: (row) => row.id,
    columns: [
      { accessorKey: "name", header: t("customers.columns.name") },
      { accessorKey: "phone", header: t("customers.columns.phone"), enableSorting: false },
      { accessorKey: "email", header: t("customers.columns.email"), enableSorting: false },
      {
        accessorKey: "is_default",
        header: t("customers.columns.default"),
        enableSorting: false,
        cell: (value) => (value ? t("common.yes") : "-"),
      },
    ],
    defaultState: {
      pagination: { pageIndex: 0, pageSize: 20 },
      sorting: [],
      globalFilter: "",
      columnFilters: [],
    },
    capabilities: { sorting: true, search: true },
    fetchPage: async ({ query, signal }) => {
      const { data, meta } = await getJSON(apiPaths.customers, query, signal);
      return { rows: data || [], total: meta?.total || 0 };
    },
  });
}
