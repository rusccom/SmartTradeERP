import { Plus } from "lucide-react";

import { useServerDataTable } from "../../../model/data-table/useServerDataTable";
import DataTable from "./DataTable";

function ServerListTable(props) {
  const table = useServerDataTable(props.preset);
  const api = readPublicApi(table);
  return (
    <>
      <DataTable {...readDataTableProps(props, table, api)} />
      {renderSlot(props.children, api)}
    </>
  );
}

function readDataTableProps(props, table, api) {
  return {
    columns: props.preset.columns,
    data: table.data,
    getRowId: props.preset.rowId,
    searchable: readSearchable(props),
    rowCount: table.total,
    loading: table.loading,
    error: table.error,
    onRetry: table.retry,
    actions: renderActions(props, api),
    components: props.components,
    ...readOptionalProps(props),
    ...table.tableState,
  };
}

function readOptionalProps(props) {
  return {
    onRowClick: props.onRowClick,
    emptyText: props.emptyText,
    expandable: props.expandable,
    getSubRows: props.getSubRows,
  };
}

function readSearchable(props) {
  return props.searchable ?? (props.preset.capabilities.search === true);
}

function readPublicApi(table) {
  return {
    data: table.data,
    total: table.total,
    loading: table.loading,
    error: table.error,
    retry: table.retry,
    refresh: table.retry,
  };
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
