import { apiPaths } from "../../../../shared/api/client";
import { getJSON } from "../../../../shared/api/http";
import { createTablePreset } from "../../../../shared/model/data-table/createTablePreset";

export const documentsTablePreset = createTablePreset({
  id: "documents",
  rowId: (row) => row.id,
  columns: [
    { accessorKey: "number", header: "Номер" },
    {
      accessorKey: "doc_type",
      header: "Тип",
      enableSorting: false,
      filterVariant: "select",
      filterOptions: [
        { value: "RECEIPT", label: "Приёмка" },
        { value: "SALE", label: "Продажа" },
        { value: "WRITEOFF", label: "Списание" },
        { value: "TRANSFER", label: "Перемещение" },
        { value: "RETURN", label: "Возврат" },
      ],
    },
    {
      accessorKey: "status",
      header: "Статус",
      enableSorting: false,
      filterVariant: "select",
      filterOptions: ["draft", "posted", "cancelled"],
    },
    { accessorKey: "total_cost", header: "Сумма" },
    { accessorKey: "date", header: "Дата" },
  ],
  capabilities: { sorting: true, search: true },
  mapStateToQuery: (state) => ({
    type: state.columnFilters.find((item) => item.id === "doc_type")?.value,
    // Skip auto-mapping for doc_type to avoid duplicate doc_type/type params.
    doc_type: undefined,
  }),
  fetchPage: async ({ query, signal }) => {
    const response = await getJSON(apiPaths.documents, query, signal);
    return {
      rows: (response.data || []).map((row) => ({ ...row, doc_type: row.type })),
      total: response.meta?.total || 0,
    };
  },
});
