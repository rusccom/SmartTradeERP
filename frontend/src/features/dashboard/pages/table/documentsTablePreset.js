import { apiPaths } from "../../../../shared/api/publicApi";
import { createApiTablePreset } from "../../../../shared/model/data-table/createApiTablePreset";
import { documentsTableSorting } from "./documentsTableSorting";

export function createDocumentsTablePreset(t, formatMoney) {
  return createApiTablePreset({
    id: "documents",
    path: apiPaths.documents,
    rowId: readDocumentId,
    columns: createColumns(t),
    capabilities: { sorting: true, search: true },
    sorting: documentsTableSorting,
    mapRows: (rows) => mapDocumentRows(rows, formatMoney),
  });
}

function createColumns(t) {
  return [
    { accessorKey: "number", header: t("documents.columns.number") },
    createTypeColumn(t),
    createStatusColumn(t),
    { accessorKey: "total_cost_label", header: t("documents.columns.totalCost") },
    { accessorKey: "date", header: t("documents.columns.date") },
  ];
}

function createTypeColumn(t) {
  return {
    accessorKey: "doc_type",
    header: t("documents.columns.type"),
    cell: (value) => readDocumentTypeLabel(value, t),
  };
}

function createStatusColumn(t) {
  return {
    accessorKey: "status",
    header: t("documents.columns.status"),
    cell: (value) => readStatusLabel(value, t),
  };
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

function mapDocumentRows(rows, formatMoney) {
  return rows.map((row) => ({
    ...row,
    doc_type: row.type,
    total_cost_label: formatMoney(row.total_cost),
  }));
}

function readDocumentId(row) {
  return row.id;
}
