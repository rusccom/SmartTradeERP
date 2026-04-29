import { Search } from "lucide-react";

import { useI18n } from "../../../i18n/useI18n";

function DataTableToolbar(props) {
  const { t } = useI18n();
  return (
    <div className="dt-toolbar">
      <div className="dt-toolbar-left">
        {props.searchable && (
          <SearchField
            value={props.globalFilter}
            onChange={props.onGlobalFilterChange}
            placeholder={t("dataTable.searchPlaceholder")}
          />
        )}
      </div>
      <div className="dt-toolbar-right">
        {props.showCount && <span className="dt-count">{t("dataTable.totalCount", { count: props.rowCount })}</span>}
        {props.actions && <div className="dt-toolbar-actions">{props.actions}</div>}
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
