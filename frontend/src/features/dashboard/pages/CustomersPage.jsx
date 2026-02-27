import { useServerDataTable } from "../../../shared/model/data-table/useServerDataTable";
import DataTable from "../../../shared/ui/data-table/DataTable";
import { customersTablePreset } from "./table/customersTablePreset";

function CustomersPage() {
  const { data, total, loading, error, retry, tableState } = useServerDataTable(customersTablePreset);
  return (
    <DataTable
      columns={customersTablePreset.columns}
      data={data}
      getRowId={customersTablePreset.rowId}
      searchable={customersTablePreset.capabilities.search === true}
      rowCount={total}
      loading={loading}
      error={error}
      onRetry={retry}
      toolbar={<button className="dt-page-btn" type="button">+ Добавить клиента</button>}
      {...tableState}
    />
  );
}

export default CustomersPage;
