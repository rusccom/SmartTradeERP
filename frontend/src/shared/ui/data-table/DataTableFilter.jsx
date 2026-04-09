import { useI18n } from "../../i18n/useI18n";

function DataTableFilter({ column }) {
  const { t } = useI18n();
  const variant = column.columnDef.meta?.filterVariant || "text";
  if (variant === "select") {
    return <SelectFilter column={column} allLabel={t("dataTable.filterAll")} />;
  }
  if (variant === "range") {
    return (
      <RangeFilter
        column={column}
        minLabel={t("dataTable.filterMin")}
        maxLabel={t("dataTable.filterMax")}
      />
    );
  }
  return <TextFilter column={column} />;
}

function TextFilter({ column }) {
  const value = column.getFilterValue() || "";
  return (
    <input
      className="dt-filter"
      type="text"
      value={value}
      onChange={(event) => column.setFilterValue(event.target.value)}
    />
  );
}

function SelectFilter({ column, allLabel }) {
  const value = column.getFilterValue() || "";
  const options = normalizeOptions(column.columnDef.meta?.filterOptions || []);
  return (
    <select className="dt-filter" value={value} onChange={(event) => column.setFilterValue(event.target.value)}>
      <option value="">{allLabel}</option>
      {options.map((option) => (
        <option key={option.value} value={option.value}>
          {option.label}
        </option>
      ))}
    </select>
  );
}

function RangeFilter({ column, minLabel, maxLabel }) {
  const value = Array.isArray(column.getFilterValue()) ? column.getFilterValue() : ["", ""];
  return (
    <div className="dt-filter-range">
      <input
        className="dt-filter"
        type="number"
        value={value[0] || ""}
        onChange={(event) => setRangeValue(column, event.target.value, value[1])}
        placeholder={minLabel}
      />
      <input
        className="dt-filter"
        type="number"
        value={value[1] || ""}
        onChange={(event) => setRangeValue(column, value[0], event.target.value)}
        placeholder={maxLabel}
      />
    </div>
  );
}

function setRangeValue(column, min, max) {
  column.setFilterValue([min, max]);
}

function normalizeOptions(options) {
  return options.map((option) => {
    if (typeof option === "string") {
      return { value: option, label: option };
    }
    return option;
  });
}

export default DataTableFilter;
