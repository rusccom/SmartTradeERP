import { Search } from "lucide-react";

function DataTableToolbar({ globalFilter, onGlobalFilterChange, searchable, toolbar, rowCount }) {
  return (
    <div className="dt-toolbar">
      <div className="dt-toolbar-left">{searchable && <SearchField value={globalFilter} onChange={onGlobalFilterChange} />}</div>
      <div className="dt-toolbar-right">
        <span className="dt-count">Всего: {rowCount}</span>
        {toolbar && <div className="dt-toolbar-actions">{toolbar}</div>}
      </div>
    </div>
  );
}

function SearchField({ value, onChange }) {
  return (
    <label className="dt-search-wrap">
      <Search className="dt-search-icon" size={14} />
      <input className="dt-search" type="text" value={value || ""} onChange={(event) => onChange(event.target.value)} placeholder="Поиск..." />
    </label>
  );
}

export default DataTableToolbar;
