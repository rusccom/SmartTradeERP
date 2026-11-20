import { useId } from "react";

export default function DocumentEditableTable({ rows, columns, selection, onSelect, labels }) {
  const id = useId();
  const selected = new Set(selection);
  const all = rows.length > 0 && rows.every((row) => selected.has(row.key));
  const toggle = (key) => onSelect(selected.has(key) ? selection.filter((k) => k !== key) : [...selection, key]);
  return <div className="dt-scroll document-lines-scroll" role="region" aria-label={labels.table} tabIndex={0}>
    <table className="dt-table"><thead><tr>
      {onSelect && <th scope="col"><input type="checkbox" checked={all} disabled={!rows.length}
        aria-label={labels.all} ref={(el) => { if (el) el.indeterminate = !all && selected.size > 0; }}
        onChange={() => onSelect(all ? [] : rows.map((r) => r.key))} /></th>}
      {columns.map((c) => <th scope="col" key={c.id} id={`${id}-${c.id}`}>{c.label}</th>)}
    </tr></thead><tbody>{rows.map((row, index) => <tr key={row.key}>
      {onSelect && <td><input type="checkbox" aria-label={labels.row(row)} checked={selected.has(row.key)} onChange={() => toggle(row.key)} /></td>}
      {columns.map((c) => <td key={c.id} headers={`${id}-${c.id}`}>{c.cell(row, index)}</td>)}
    </tr>)}</tbody></table>
    {!rows.length && <p className="document-lines-empty">{labels.empty}</p>}
  </div>;
}
