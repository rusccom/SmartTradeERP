function ErpTable({ columns, data, sortKey, sortDir, onSort, onRowClick }) {
  return (
    <div className="erp-table-scroll">
      <table className="erp-table">
        <thead>
          <tr>
            {columns.map((col) => (
              <th
                key={col.accessorKey}
                className="erp-th"
                onClick={() => onSort(col.accessorKey)}
              >
                {col.header}
                <SortIndicator
                  active={sortKey === col.accessorKey}
                  dir={sortDir}
                />
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.length === 0 && (
            <tr>
              <td className="erp-td erp-empty" colSpan={columns.length}>
                No data
              </td>
            </tr>
          )}
          {data.map((row) => (
            <tr
              key={row.id}
              className="erp-row"
              onClick={() => onRowClick(row)}
            >
              {columns.map((col) => (
                <td key={col.accessorKey} className="erp-td">
                  {row[col.accessorKey] ?? "—"}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function SortIndicator({ active, dir }) {
  if (!active) return <span className="erp-sort-icon"> ↕</span>;
  return (
    <span className="erp-sort-icon"> {dir === "asc" ? "↑" : "↓"}</span>
  );
}

export default ErpTable;
