import { Search } from "lucide-react";

import { useI18n } from "../../i18n/useI18n";

function DataTableToolbar({ globalFilter, onGlobalFilterChange, searchable, toolbar, rowCount }) {
  const { t } = useI18n();
  return (
    <div className="dt-toolbar">
      <div className="dt-toolbar-left">
        {searchable && (
          <SearchField
            value={globalFilter}
            onChange={onGlobalFilterChange}
            placeholder={t("dataTable.searchPlaceholder")}
          />
        )}
      </div>
      <div className="dt-toolbar-right">
        <span className="dt-count">{t("dataTable.totalCount", { count: rowCount })}</span>
        {toolbar && <div className="dt-toolbar-actions">{toolbar}</div>}
      </div>
    </div>
  );
}

function SearchField({ value, onChange, placeholder }) {
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

export default DataTableToolbar;
