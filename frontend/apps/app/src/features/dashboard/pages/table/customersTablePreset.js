import { apiPaths } from "../../../../shared/api/publicApi";
import { createApiTablePreset } from "../../../../shared/model/data-table/createApiTablePreset";
import { customersTableSorting } from "./customersTableSorting";

export function createCustomersTablePreset(t) {
  return createApiTablePreset({
    id: "customers",
    path: apiPaths.customers,
    rowId: readCustomerId,
    columns: createColumns(t),
    capabilities: { sorting: true, search: true },
    sorting: customersTableSorting,
  });
}

function createColumns(t) {
  return [
    { accessorKey: "name", header: t("customers.columns.name") },
    { accessorKey: "phone", header: t("customers.columns.phone") },
    { accessorKey: "email", header: t("customers.columns.email") },
    createDefaultColumn(t),
  ];
}

function createDefaultColumn(t) {
  return {
    accessorKey: "is_default",
    header: t("customers.columns.default"),
    cell: (value) => (value ? t("common.yes") : "-"),
  };
}

function readCustomerId(row) {
  return row.id;
}
