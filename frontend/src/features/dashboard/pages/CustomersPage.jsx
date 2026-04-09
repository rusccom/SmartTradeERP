import { useMemo } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import { useServerDataTable } from "../../../shared/model/data-table/useServerDataTable";
import DataTable from "../../../shared/ui/data-table/DataTable";
import { createCustomersTablePreset } from "./table/customersTablePreset";

function CustomersPage() {
  const { t } = useI18n();
  const preset = useMemo(() => createCustomersTablePreset(t), [t]);
  const { data, total, loading, error, retry, tableState } = useServerDataTable(preset);
  return (
    <DataTable
      columns={preset.columns}
      data={data}
      getRowId={preset.rowId}
      searchable={preset.capabilities.search === true}
      rowCount={total}
      loading={loading}
      error={error}
      onRetry={retry}
      toolbar={<button className="dt-page-btn" type="button">+ {t("customers.addButton")}</button>}
      {...tableState}
    />
  );
}

export default CustomersPage;
