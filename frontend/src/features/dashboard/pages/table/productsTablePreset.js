import { apiPaths } from "../../../../shared/api/client";
import { getJSON } from "../../../../shared/api/http";
import { createTablePreset } from "../../../../shared/model/data-table/createTablePreset";

export const productsTablePreset = createTablePreset({
  id: "products",
  rowId: (row) => row.id,
  columns: [
    { accessorKey: "name", header: "Название", enableFilter: true },
    {
      accessorKey: "is_composite",
      header: "Составной",
      enableSorting: false,
      filterVariant: "select",
      filterOptions: [
        { value: "true", label: "Да" },
        { value: "false", label: "Нет" },
      ],
      cell: (value) => (value ? "Да" : "Нет"),
    },
    { accessorKey: "updated_at", header: "Обновлён", enableSorting: false },
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
