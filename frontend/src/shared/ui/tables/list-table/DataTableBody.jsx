import { flexRender } from "@tanstack/react-table";
import { ChevronDown, ChevronRight } from "lucide-react";
import { useRef, useState } from "react";

import { useI18n } from "../../../i18n/useI18n";
import DataTableOpenLink from "./DataTableOpenLink";

function DataTableBody({ table, onRowOpen, emptyText, subRows }) {
  const { t } = useI18n();
  const [expandedRows, setExpandedRows] = useState(() => new Set());
  const subRowsCache = useRef(new Map());
  const rows = table.getRowModel().rows;
  if (rows.length === 0) {
    return <DataTableEmpty colSpan={table.getVisibleLeafColumns().length} emptyText={emptyText} />;
  }
  return (
    <tbody>
      {rows.map((row) => (
        <BodyRows
          key={row.id}
          row={row}
          onRowOpen={onRowOpen}
          subRows={subRows}
          expandedRows={expandedRows}
          setExpandedRows={setExpandedRows}
          subRowsCache={subRowsCache}
          t={t}
        />
      ))}
    </tbody>
  );
}

function DataTableEmpty({ colSpan, emptyText }) {
  return (
    <tbody>
      <tr>
        <td className="dt-empty" colSpan={colSpan}>
          {emptyText}
        </td>
      </tr>
    </tbody>
  );
}

function BodyRows(props) {
  const { row, expandedRows, onRowOpen, subRowsCache, t } = props;
  const expanded = expandedRows.has(row.id);
  const loadedSubRows = subRowsCache.current.get(row.id) || [];
  return (
    <>
      <MainRow {...props} expanded={expanded} />
      {expanded && <SubRows row={row} rows={loadedSubRows} t={t} onRowOpen={onRowOpen} />}
    </>
  );
}

function MainRow(props) {
  const { row } = props;
  const firstDataColumnId = readFirstDataColumnId(row);
  return (
    <tr className="dt-row">
      {row.getVisibleCells().map((cell) => (
        <td key={cell.id} className="dt-td">
          <MainCell
            cell={cell}
            isFirstData={cell.column.id === firstDataColumnId}
            {...props}
          />
        </td>
      ))}
    </tr>
  );
}

function MainCell(props) {
  const { cell, isFirstData } = props;
  const content = renderMainCellValue(cell);
  if (!canShowExpand(props)) {
    return content;
  }
  return (
    <div className="dt-expand-cell">
      <button className="dt-expand-btn" type="button" onClick={(event) => handleExpand(createExpandParams(props, event))}>
        {props.expanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
      </button>
      <span>{content}</span>
    </div>
  );
}

function SubRows({ row, rows, onRowOpen, t }) {
  const firstDataColumnId = readFirstDataColumnId(row);
  return rows.map((subRow, index) => (
    <tr key={`${row.id}-sub-${index}`} className="dt-row dt-row--sub">
      {row.getVisibleCells().map((cell) => (
        <td key={`${cell.id}-sub-${index}`} className="dt-td">
          <SubCell columnDef={cell.column.columnDef} row={subRow} isFirstData={cell.column.id === firstDataColumnId} onRowOpen={onRowOpen} t={t} />
        </td>
      ))}
    </tr>
  ));
}

function SubCell({ columnDef, row, isFirstData, onRowOpen, t }) {
  const content = resolveSubCellValue(columnDef, row, onRowOpen, t);
  if (!isFirstData) {
    return content;
  }
  return <div className="dt-sub-indent">{content}</div>;
}

function canShowExpand({ isFirstData, row, subRows }) {
  if (!subRowsEnabled(subRows) || !isFirstData) return false;
  const canExpand = subRows?.canExpand;
  return typeof canExpand === "function" ? canExpand(row.original) : true;
}

function createExpandParams(props, event) {
  const { expandedRows, row, setExpandedRows, subRows, subRowsCache } = props;
  return { event, row, expandedRows, setExpandedRows, subRowsCache, getSubRows: subRows?.getRows };
}

function readFirstDataColumnId(row) {
  return row.getVisibleCells().find((cell) => cell.column.id !== "select")?.column.id;
}

function resolveSubCellValue(columnDef, row, onRowOpen, t) {
  const key = readAccessorKey(columnDef);
  const value = key ? row[key] : "";
  const renderer = columnDef.meta?.rawCell;
  if (typeof renderer === "function") {
    return renderer(value, row, createCellApi(columnDef, row, onRowOpen));
  }
  if (columnDef.meta?.openOnClick === true) {
    return renderOpenCell(value, row, onRowOpen);
  }
  if (value === undefined || value === null) {
    return formatCellValue(value);
  }
  if (typeof value === "boolean") {
    return value ? t("common.yes") : t("common.no");
  }
  return String(value);
}

function createCellApi(columnDef, row, onRowOpen) {
  return {
    openLink: columnDef.meta?.openOnClick ? createOpenLink(onRowOpen, row) : null,
  };
}

function createOpenLink(onRowOpen, row) {
  return (children, target = row) => <DataTableOpenLink onOpen={onRowOpen} target={target}>{children}</DataTableOpenLink>;
}

function renderOpenCell(value, row, onRowOpen) {
  return <DataTableOpenLink onOpen={onRowOpen} target={row}>{formatCellValue(value)}</DataTableOpenLink>;
}

function formatCellValue(value) {
  if (value === undefined || value === null) {
    return "";
  }
  return String(value);
}

function readAccessorKey(columnDef) {
  return columnDef.meta?.accessorKey || columnDef.accessorKey;
}

function renderMainCellValue(cell) {
  if (cell.column.columnDef.cell) {
    return flexRender(cell.column.columnDef.cell, cell.getContext());
  }
  const value = cell.getValue();
  return formatCellValue(value);
}

async function handleExpand({ event, row, expandedRows, setExpandedRows, subRowsCache, getSubRows }) {
  event.stopPropagation();
  if (expandedRows.has(row.id)) {
    setExpandedRows((prev) => removeFromSet(prev, row.id));
    return;
  }
  await loadSubRows(row, subRowsCache, getSubRows);
  setExpandedRows((prev) => addToSet(prev, row.id));
}

async function loadSubRows(row, subRowsCache, getSubRows) {
  if (!getSubRows || subRowsCache.current.has(row.id)) {
    return;
  }
  try {
    const payload = await getSubRows(row.original);
    subRowsCache.current.set(row.id, Array.isArray(payload) ? payload : []);
  } catch {
    subRowsCache.current.set(row.id, []);
  }
}

function subRowsEnabled(subRows) {
  return subRows?.enabled === true;
}

function addToSet(set, value) {
  const next = new Set(set);
  next.add(value);
  return next;
}

function removeFromSet(set, value) {
  const next = new Set(set);
  next.delete(value);
  return next;
}

export default DataTableBody;
