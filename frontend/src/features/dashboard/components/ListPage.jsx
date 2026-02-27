import { useMemo, useState } from "react";

import ErpTable from "./ErpTable";
import ModalShell from "./ModalShell";
import "./dashboard-components.css";

function ListPage({ title, columns, data }) {
  const [search, setSearch] = useState("");
  const [sortKey, setSortKey] = useState(null);
  const [sortDir, setSortDir] = useState("asc");
  const [selectedRow, setSelectedRow] = useState(null);

  const filtered = useFilteredData(data, columns, search);
  const sorted = useSortedData(filtered, sortKey, sortDir);

  function handleSort(key) {
    if (sortKey === key) {
      setSortDir((d) => (d === "asc" ? "desc" : "asc"));
    } else {
      setSortKey(key);
      setSortDir("asc");
    }
  }

  return (
    <section className="placeholder">
      <div className="list-page-header">
        <h2>{title}</h2>
        <input
          className="list-page-search"
          type="text"
          placeholder="Search..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>

      <ErpTable
        columns={columns}
        data={sorted}
        sortKey={sortKey}
        sortDir={sortDir}
        onSort={handleSort}
        onRowClick={setSelectedRow}
      />

      {selectedRow && (
        <ModalShell
          title="Details"
          onClose={() => setSelectedRow(null)}
        >
          <p>Details will be here.</p>
          <pre>{JSON.stringify(selectedRow, null, 2)}</pre>
        </ModalShell>
      )}
    </section>
  );
}

function useFilteredData(data, columns, search) {
  return useMemo(() => {
    if (!search.trim()) return data;
    const q = search.toLowerCase();
    return data.filter((row) =>
      columns.some((col) => {
        const val = row[col.accessorKey];
        return val != null && String(val).toLowerCase().includes(q);
      }),
    );
  }, [data, columns, search]);
}

function useSortedData(data, sortKey, sortDir) {
  return useMemo(() => {
    if (!sortKey) return data;
    return [...data].sort((a, b) => {
      const aVal = a[sortKey] ?? "";
      const bVal = b[sortKey] ?? "";
      const cmp = String(aVal).localeCompare(String(bVal));
      return sortDir === "asc" ? cmp : -cmp;
    });
  }, [data, sortKey, sortDir]);
}

export default ListPage;
