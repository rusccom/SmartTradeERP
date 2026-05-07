import { Plus } from "lucide-react";
import { useEffect, useState } from "react";

import { useServerDataTable } from "../../../model/data-table/useServerDataTable";
import DataTable from "./DataTable";

function ServerListTable(props) {
  const table = useServerDataTable(props.preset);
  const [rowSelection, setRowSelection] = useState({});
  useEffect(() => setRowSelection({}), [table.data]);
  const api = readPublicApi(table, props, rowSelection, setRowSelection);
  return (
    <>
      <DataTable {...readDataTableProps(props, table, api, rowSelection, setRowSelection)} />
      {renderSlot(props.children, api)}
    </>
  );
}

function readDataTableProps(props, table, api, rowSelection, setRowSelection) {
  return {
    columns: props.preset.columns,
    data: table.data,
    getRowId: props.preset.rowId,
    searchable: readSearchable(props),
    search: props.preset.search,
    sortingConfig: props.preset.sortingConfig,
    rowCount: table.total,
    loading: table.loading,
    error: table.error,
    onRetry: table.retry,
    actions: renderActions(props, api),
    components: props.components,
    selectable: props.selectable === true,
    rowSelection,
    onRowSelectionChange: setRowSelection,
    ...readOptionalProps(props),
    ...table.tableState,
  };
}

function readOptionalProps(props) {
  return {
    onRowOpen: props.onRowOpen || props.onRowClick,
    emptyText: props.emptyText,
    subRows: readSubRowsConfig(props),
  };
}

function readSubRowsConfig(props) {
  if (props.subRows) return props.subRows;
  return {
    enabled: props.expandable === true,
    getRows: props.getSubRows,
    canExpand: props.canExpandRow,
    onRowOpen: props.onSubRowOpen,
  };
}

function readSearchable(props) {
  return props.searchable ?? isPresetSearchable(props.preset);
}

function isPresetSearchable(preset) {
  return preset.capabilities.search === true && preset.search?.enabled !== false;
}

function readPublicApi(table, props, rowSelection, setRowSelection) {
  const selectedRows = readSelectedRows(table.data, props.preset.rowId, rowSelection);
  return {
    data: table.data,
    total: table.total,
    loading: table.loading,
    error: table.error,
    selectedRows,
    selectedCount: selectedRows.length,
    selectedRow: selectedRows[0] || null,
    retry: table.retry,
    refresh: table.retry,
    clearSelection: () => setRowSelection({}),
  };
}

function readSelectedRows(rows, getRowId, rowSelection) {
  return rows.filter((row, index) => rowSelection[getRowId(row, index)]);
}

function renderActions(props, api) {
  const actions = renderSlot(props.actions, api);
  const primary = renderPrimaryAction(props.primaryAction, api);
  if (!actions && !primary) {
    return null;
  }
  return (
    <>
      {actions}
      {primary}
    </>
  );
}

function renderPrimaryAction(action, api) {
  const config = typeof action === "function" ? action(api) : action;
  if (!config?.label || typeof config.onClick !== "function") {
    return null;
  }
  return (
    <button className="dt-action-primary" type="button" onClick={config.onClick}>
      <Plus size={15} />
      <span>{config.label}</span>
    </button>
  );
}

function renderSlot(slot, api) {
  return typeof slot === "function" ? slot(api) : slot;
}

export default ServerListTable;
