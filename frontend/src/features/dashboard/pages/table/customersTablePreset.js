import { apiPaths } from "../../../../shared/api/client";
import { createApiTablePreset } from "../../../../shared/model/data-table/createApiTablePreset";

export function createCustomersTablePreset(t) {
  return createApiTablePreset({
    id: "customers",
    path: apiPaths.customers,
    rowId: readCustomerId,
    columns: createColumns(t),
    capabilities: { sorting: true, search: true },
  });
}

function createColumns(t) {
  return [
    { accessorKey: "name", header: t("customers.columns.name") },
    { accessorKey: "phone", header: t("customers.columns.phone"), enableSorting: false },
    { accessorKey: "email", header: t("customers.columns.email"), enableSorting: false },
    createDefaultColumn(t),
  ];
}

function createDefaultColumn(t) {
  return {
    accessorKey: "is_default",
    header: t("customers.columns.default"),
    enableSorting: false,
    cell: (value) => (value ? t("common.yes") : "-"),
  };
}

function readCustomerId(row) {
  return row.id;
}
