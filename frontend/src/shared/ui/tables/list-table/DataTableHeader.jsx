import { flexRender } from "@tanstack/react-table";
import { ArrowDown, ArrowUp, ArrowUpDown } from "lucide-react";

function DataTableHeader({ table }) {
  return (
    <thead>
      {table.getHeaderGroups().map((group) => (
        <tr key={group.id}>
          {group.headers.map((header) => (
            <HeaderCell key={header.id} header={header} />
          ))}
        </tr>
      ))}
    </thead>
  );
}

function HeaderCell({ header }) {
  if (header.isPlaceholder) {
    return <th className="dt-th" colSpan={header.colSpan} />;
  }
  const sorted = header.column.getIsSorted();
  const className = readHeaderClass(header.column.getCanSort(), sorted);
  const onClick = header.column.getCanSort() ? header.column.getToggleSortingHandler() : undefined;
  return (
    <th className={className} colSpan={header.colSpan} style={readColumnStyle(header)}>
      <SortTrigger sortable={header.column.getCanSort()} onClick={onClick}>
        <span>{flexRender(header.column.columnDef.header, header.getContext())}</span>
        <SortIcon direction={sorted} />
      </SortTrigger>
    </th>
  );
}

function readHeaderClass(sortable, sorted) {
  const sortableClass = sortable ? "dt-th--sortable" : "";
  const sortedClass = sorted ? "dt-th--sorted" : "";
  return `dt-th ${sortableClass} ${sortedClass}`.trim();
}

function readColumnStyle(header) {
  const size = header.column.columnDef.size;
  return size ? { width: size } : undefined;
}

function SortTrigger({ sortable, onClick, children }) {
  if (!sortable) {
    return <div className="dt-th-content">{children}</div>;
  }
  return (
    <button className="dt-th-btn" type="button" onClick={onClick}>
      <span className="dt-th-content">{children}</span>
    </button>
  );
}

function SortIcon({ direction }) {
  if (direction === "asc") {
    return <ArrowUp className="dt-sort-icon dt-sort-icon--active" />;
  }
  if (direction === "desc") {
    return <ArrowDown className="dt-sort-icon dt-sort-icon--active" />;
  }
  return <ArrowUpDown className="dt-sort-icon" />;
}

export default DataTableHeader;
