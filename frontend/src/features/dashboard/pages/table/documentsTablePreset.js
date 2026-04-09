import { apiPaths } from "../../../../shared/api/client";
import { getJSON } from "../../../../shared/api/http";
import { createTablePreset } from "../../../../shared/model/data-table/createTablePreset";

export function createDocumentsTablePreset(t) {
  return createTablePreset({
    id: "documents",
    rowId: (row) => row.id,
    columns: [
      { accessorKey: "number", header: t("documents.columns.number") },
      {
        accessorKey: "doc_type",
        header: t("documents.columns.type"),
        enableSorting: false,
        filterVariant: "select",
        filterOptions: createDocumentTypeOptions(t),
        cell: (value) => readDocumentTypeLabel(value, t),
      },
      {
        accessorKey: "status",
        header: t("documents.columns.status"),
        enableSorting: false,
        filterVariant: "select",
        filterOptions: createStatusOptions(t),
        cell: (value) => readStatusLabel(value, t),
      },
      { accessorKey: "total_cost", header: t("documents.columns.totalCost") },
      { accessorKey: "date", header: t("documents.columns.date") },
    ],
    capabilities: { sorting: true, search: true },
    mapStateToQuery: (state) => ({
      type: state.columnFilters.find((item) => item.id === "doc_type")?.value,
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
}

function createDocumentTypeOptions(t) {
  return [
    { value: "RECEIPT", label: t("documents.types.receipt") },
    { value: "SALE", label: t("documents.types.sale") },
    { value: "WRITEOFF", label: t("documents.types.writeoff") },
    { value: "TRANSFER", label: t("documents.types.transfer") },
    { value: "RETURN", label: t("documents.types.return") },
  ];
}

function createStatusOptions(t) {
  return [
    { value: "draft", label: t("documents.status.draft") },
    { value: "posted", label: t("documents.status.posted") },
    { value: "cancelled", label: t("documents.status.cancelled") },
  ];
}

function readDocumentTypeLabel(value, t) {
  const labels = {
    RECEIPT: t("documents.types.receipt"),
    RETURN: t("documents.types.return"),
    SALE: t("documents.types.sale"),
    TRANSFER: t("documents.types.transfer"),
    WRITEOFF: t("documents.types.writeoff"),
  };
  return labels[value] || value;
}

function readStatusLabel(value, t) {
  const labels = {
    cancelled: t("documents.status.cancelled"),
    draft: t("documents.status.draft"),
    posted: t("documents.status.posted"),
  };
  return labels[value] || value;
}
