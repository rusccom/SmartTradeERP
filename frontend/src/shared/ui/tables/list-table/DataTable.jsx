import { getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { useMemo } from "react";

import { useI18n } from "../../../i18n/useI18n";
import DataTableBody from "./DataTableBody";
import DataTableError from "./DataTableError";
import DataTableHeader from "./DataTableHeader";
import DataTablePagination from "./DataTablePagination";
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
        slots={slots}
        loadingClass={loadingClass}
        emptyText={emptyText}
      />
      {slots.pagination && <DataTablePagination table={table} />}
    </section>
  );
}

function useDataTableInstance(props) {
  const mappedColumns = useMemo(() => mapColumns(props.columns), [props.columns]);
  return useReactTable({
    data: props.data,
    columns: mappedColumns,
    getRowId: props.getRowId,
    state: readControlledState(props),
    onSortingChange: props.onSortingChange,
    onColumnFiltersChange: props.onColumnFiltersChange,
    onGlobalFilterChange: props.onGlobalFilterChange,
    onPaginationChange: props.onPaginationChange,
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
    columnFilters: props.columnFilters,
    globalFilter: props.globalFilter,
    pagination: props.pagination,
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

function TableBlock({ table, props, slots, loadingClass, emptyText }) {
  return (
    <div className={`dt-table-scroll ${loadingClass}`.trim()}>
      <table className="dt-table">
        <DataTableHeader table={table} showFilters={slots.filters} />
        <DataTableBody
          table={table}
          onRowClick={props.onRowClick}
          emptyText={emptyText}
          expandable={props.expandable === true}
          getSubRows={props.getSubRows}
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
    filters: components.filters !== false,
    pagination: components.pagination !== false,
  };
}

function mapColumns(columns) {
  return columns.map((column) => mapColumn(column));
}

function mapColumn(column) {
  const mapped = {
    id: column.accessorKey,
    accessorKey: column.accessorKey,
    header: column.header,
    size: column.size,
    enableSorting: column.enableSorting !== false,
    enableColumnFilter: Boolean(column.enableFilter || column.filterVariant),
    meta: {
      filterVariant: column.filterVariant,
      filterOptions: column.filterOptions,
      rawCell: column.cell,
      accessorKey: column.accessorKey,
    },
  };
  if (column.cell) {
    mapped.cell = (info) => column.cell(info.getValue(), info.row.original);
  }
  return mapped;
}

export default DataTable;
