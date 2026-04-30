import { Search } from "lucide-react";

function DataTableSearchField({ value, onChange, placeholder }) {
  return (
    <label className="dt-search-wrap">
      <Search className="dt-search-icon" size={14} />
      <input
        className="dt-search"
        type="text"
        value={value || ""}
        onChange={(event) => onChange(event.target.value)}
        placeholder={placeholder}
      />
    </label>
  );
}

export default DataTableSearchField;
