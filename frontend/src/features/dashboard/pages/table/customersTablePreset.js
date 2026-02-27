import { apiPaths } from "../../../../shared/api/client";
import { getJSON } from "../../../../shared/api/http";
import { createTablePreset } from "../../../../shared/model/data-table/createTablePreset";

export const customersTablePreset = createTablePreset({
  id: "customers",
  rowId: (row) => row.id,
  columns: [
    { accessorKey: "name", header: "Имя" },
    { accessorKey: "phone", header: "Телефон", enableSorting: false },
    { accessorKey: "email", header: "Email", enableSorting: false },
    {
      accessorKey: "is_default",
      header: "По умолчанию",
      enableSorting: false,
      cell: (value) => (value ? "Да" : "—"),
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
