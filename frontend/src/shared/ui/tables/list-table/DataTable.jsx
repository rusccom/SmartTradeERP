import { getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { useMemo } from "react";

import { useI18n } from "../../../i18n/useI18n";
import DataTableBody from "./DataTableBody";
import DataTableError from "./DataTableError";
import DataTableHeader from "./DataTableHeader";
import DataTablePagination from "./DataTablePagination";
import DataTableSelectionCheckbox from "./DataTableSelectionCheckbox";
import DataTableToolbar from "./DataTableToolbar";
import "./data-table-actions.css";
import "./data-table.css";
import "./data-table-state.css";

function DataTable(props) {
  const { t } = useI18n();
  const table = useDataTableInstance(props);
  const slots = readSlots(props);
  const emptyText = props.emptyText || t("dataTable.emptyText");
  const loadingClass = props.loading ? "dt-loading" : "";
  return (
    <section className="dt-wrapper">
      {slots.toolbar && <ToolbarBlock props={props} slots={slots} />}
      <ErrorBlock error={props.error} onRetry={props.onRetry} />
      <TableBlock
        table={table}
        props={props}
        loadingClass={loadingClass}
        emptyText={emptyText}
      />
      {slots.pagination && <DataTablePagination table={table} />}
    </section>
  );
}

function useDataTableInstance(props) {
  const mappedColumns = useMemo(() => readColumns(props), [props.columns, props.selectable]);
  return useReactTable({
    data: props.data,
    columns: mappedColumns,
    getRowId: props.getRowId,
    state: readControlledState(props),
    onSortingChange: props.onSortingChange,
    onGlobalFilterChange: props.onGlobalFilterChange,
    onPaginationChange: props.onPaginationChange,
    onRowSelectionChange: props.onRowSelectionChange,
    enableRowSelection: props.selectable === true,
    manualPagination: true,
    manualSorting: true,
    manualFiltering: true,
    rowCount: Number(props.rowCount) || 0,
    getCoreRowModel: getCoreRowModel(),
  });
}

function readControlledState(props) {
  return {
    sorting: props.sorting,
    globalFilter: props.globalFilter,
    pagination: props.pagination,
    rowSelection: props.rowSelection || {},
  };
}

function ToolbarBlock({ props, slots }) {
  return (
    <DataTableToolbar
      globalFilter={props.globalFilter}
      onGlobalFilterChange={props.onGlobalFilterChange}
      searchable={slots.search}
      actions={props.actions}
      showCount={slots.count}
      rowCount={props.rowCount}
    />
  );
}

function ErrorBlock({ error, onRetry }) {
  return error ? <DataTableError message={error} onRetry={onRetry} /> : null;
}

function TableBlock({ table, props, loadingClass, emptyText }) {
  const showHeader = table.getRowModel().rows.length > 0;
  return (
    <div className={`dt-table-scroll ${loadingClass}`.trim()}>
      <table className="dt-table">
        {showHeader && <DataTableHeader table={table} />}
        <DataTableBody
          table={table}
          onRowOpen={props.onRowOpen}
          emptyText={emptyText}
          subRows={props.subRows}
        />
      </table>
    </div>
  );
}

function readSlots(props) {
  const components = props.components || {};
  const search = props.searchable !== false && components.search !== false;
  const actions = Boolean(props.actions) && components.actions !== false;
  const count = components.count !== false;
  return {
    search,
    actions,
    count,
    toolbar: components.toolbar !== false && (search || actions || count),
    pagination: components.pagination !== false,
  };
}

function readColumns(props) {
  const columns = props.columns.map((column) => mapColumn(column));
  return props.selectable ? [createSelectionColumn(), ...columns] : columns;
}

function mapColumn(column) {
  const mapped = {
    id: column.id || column.accessorKey,
    accessorKey: column.accessorKey,
    header: column.header,
    size: column.size,
    enableSorting: column.enableSorting !== false,
    meta: {
      rawCell: column.cell,
      accessorKey: column.accessorKey,
    },
  };
  if (column.cell) {
    mapped.cell = (info) => column.cell(info.getValue(), info.row.original);
  }
  return mapped;
}

function createSelectionColumn() {
  return {
    id: "select",
    size: 42,
    enableSorting: false,
    header: ({ table }) => (
      <DataTableSelectionCheckbox
        checked={table.getIsAllPageRowsSelected()}
        indeterminate={table.getIsSomePageRowsSelected()}
        onChange={table.getToggleAllPageRowsSelectedHandler()}
        label="Select all rows"
      />
    ),
    cell: ({ row }) => (
      <DataTableSelectionCheckbox
        checked={row.getIsSelected()}
        disabled={!row.getCanSelect()}
        onChange={row.getToggleSelectedHandler()}
        label="Select row"
      />
    ),
  };
}

export default DataTable;
