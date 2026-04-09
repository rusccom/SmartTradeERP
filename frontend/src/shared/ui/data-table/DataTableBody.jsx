import { flexRender } from "@tanstack/react-table";
import { ChevronDown, ChevronRight } from "lucide-react";
import { useRef, useState } from "react";

import { useI18n } from "../../i18n/useI18n";

function DataTableBody({ table, onRowClick, emptyText, expandable, getSubRows }) {
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
          onRowClick={onRowClick}
          expandable={expandable}
          getSubRows={getSubRows}
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
  const { row, expandedRows, subRowsCache, t } = props;
  const expanded = expandedRows.has(row.id);
  const subRows = subRowsCache.current.get(row.id) || [];
  return (
    <>
      <MainRow {...props} expanded={expanded} />
      {expanded && <SubRows row={row} subRows={subRows} t={t} />}
    </>
  );
}

function MainRow({ row, onRowClick, expandable, getSubRows, expanded, expandedRows, setExpandedRows, subRowsCache }) {
  const className = readMainRowClass(onRowClick);
  const onClick = onRowClick ? () => onRowClick(row.original) : undefined;
  return (
    <tr className={className} onClick={onClick}>
      {row.getVisibleCells().map((cell, index) => (
        <td key={cell.id} className="dt-td">
          <MainCell
            cell={cell}
            isFirst={index === 0}
            expanded={expanded}
            expandable={expandable}
            onExpand={(event) => handleExpand({ event, row, expandedRows, setExpandedRows, subRowsCache, getSubRows })}
          />
        </td>
      ))}
    </tr>
  );
}

function MainCell({ cell, isFirst, expanded, expandable, onExpand }) {
  const content = renderMainCellValue(cell);
  if (!expandable || !isFirst) {
    return content;
  }
  return (
    <div className="dt-expand-cell">
      <button className="dt-expand-btn" type="button" onClick={onExpand}>
        {expanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
      </button>
      <span>{content}</span>
    </div>
  );
}

function SubRows({ row, subRows, t }) {
  return subRows.map((subRow, index) => (
    <tr key={`${row.id}-sub-${index}`} className="dt-row dt-row--sub">
      {row.getVisibleCells().map((cell, cellIndex) => (
        <td key={`${cell.id}-sub-${index}`} className="dt-td">
          <SubCell columnDef={cell.column.columnDef} row={subRow} isFirst={cellIndex === 0} t={t} />
        </td>
      ))}
    </tr>
  ));
}

function SubCell({ columnDef, row, isFirst, t }) {
  const content = resolveSubCellValue(columnDef, row, t);
  if (!isFirst) {
    return content;
  }
  return <div className="dt-sub-indent">{content}</div>;
}

function resolveSubCellValue(columnDef, row, t) {
  const key = readAccessorKey(columnDef);
  const value = key ? row[key] : "";
  const renderer = columnDef.meta?.rawCell;
  if (typeof renderer === "function") {
    return renderer(value, row);
  }
  if (value === undefined || value === null) {
    return "";
  }
  if (typeof value === "boolean") {
    return value ? t("common.yes") : t("common.no");
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
  if (value === undefined || value === null) {
    return "";
  }
  return String(value);
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

function readMainRowClass(onRowClick) {
  const clickable = onRowClick ? "dt-row--clickable" : "";
  return `dt-row ${clickable}`.trim();
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
