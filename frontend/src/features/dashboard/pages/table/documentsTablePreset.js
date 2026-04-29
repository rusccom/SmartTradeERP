import { apiPaths } from "../../../../shared/api/client";
import { createApiTablePreset } from "../../../../shared/model/data-table/createApiTablePreset";

export function createDocumentsTablePreset(t) {
  return createApiTablePreset({
    id: "documents",
    path: apiPaths.documents,
    rowId: readDocumentId,
    columns: createColumns(t),
    capabilities: { sorting: true, search: true },
    mapStateToQuery: mapDocumentStateToQuery,
    mapRows: mapDocumentRows,
  });
}

function createColumns(t) {
  return [
    { accessorKey: "number", header: t("documents.columns.number") },
    createTypeColumn(t),
    createStatusColumn(t),
    { accessorKey: "total_cost", header: t("documents.columns.totalCost") },
    { accessorKey: "date", header: t("documents.columns.date") },
  ];
}

function createTypeColumn(t) {
  return {
    accessorKey: "doc_type",
    header: t("documents.columns.type"),
    enableSorting: false,
    filterVariant: "select",
    filterOptions: createDocumentTypeOptions(t),
    cell: (value) => readDocumentTypeLabel(value, t),
  };
}

function createStatusColumn(t) {
  return {
    accessorKey: "status",
    header: t("documents.columns.status"),
    enableSorting: false,
    filterVariant: "select",
    filterOptions: createStatusOptions(t),
    cell: (value) => readStatusLabel(value, t),
  };
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

function mapDocumentStateToQuery(state) {
  return {
    type: state.columnFilters.find((item) => item.id === "doc_type")?.value,
    doc_type: undefined,
  };
}

function mapDocumentRows(rows) {
  return rows.map((row) => ({ ...row, doc_type: row.type }));
}

function readDocumentId(row) {
  return row.id;
}
