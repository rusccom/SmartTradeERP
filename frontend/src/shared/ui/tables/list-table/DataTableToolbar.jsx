import { useI18n } from "../../../i18n/useI18n";
import DataTableSearchField from "./DataTableSearchField";

function DataTableToolbar(props) {
  const { t } = useI18n();
  return (
    <div className="dt-toolbar">
      <div className="dt-toolbar-left">
        {props.searchable && (
          <DataTableSearchField
            value={props.globalFilter}
            onChange={props.onGlobalFilterChange}
            placeholder={readSearchPlaceholder(props.search, t)}
          />
        )}
      </div>
      <div className="dt-toolbar-right">
        {props.actions && <div className="dt-toolbar-actions">{props.actions}</div>}
      </div>
    </div>
  );
}

function readSearchPlaceholder(search, t) {
  if (search?.placeholder) {
    return search.placeholder;
  }
  return t(search?.placeholderKey || "dataTable.searchPlaceholder");
}

export default DataTableToolbar;
