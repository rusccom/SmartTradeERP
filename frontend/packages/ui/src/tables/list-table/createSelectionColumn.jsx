import DataTableSelectionCheckbox from "./DataTableSelectionCheckbox";

export function createSelectionColumn(labels = {}) {
  return {
    id: "select",
    size: 42,
    enableSorting: false,
    header: ({ table }) => pageCheckbox(table, labels.all),
    cell: ({ row }) => rowCheckbox(row, labels.row),
  };
}

function pageCheckbox(table, label = "Select all rows") {
  return (
    <DataTableSelectionCheckbox
      checked={table.getIsAllPageRowsSelected()}
      disabled={table.getRowModel().rows.length === 0}
      indeterminate={table.getIsSomePageRowsSelected()}
      onChange={table.getToggleAllPageRowsSelectedHandler()}
      label={label}
    />
  );
}

function rowCheckbox(row, label) {
  return (
    <DataTableSelectionCheckbox
      checked={row.getIsSelected()}
      disabled={!row.getCanSelect()}
      onChange={row.getToggleSelectedHandler()}
      label={label ? label(row.original) : "Select row"}
    />
  );
}
