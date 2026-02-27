import { useServerDataTable } from "../../../shared/model/data-table/useServerDataTable";
import DataTable from "../../../shared/ui/data-table/DataTable";
import { documentsTablePreset } from "./table/documentsTablePreset";

function DocumentsPage() {
  const { data, total, loading, error, retry, tableState } = useServerDataTable(documentsTablePreset);
  return (
    <DataTable
      columns={documentsTablePreset.columns}
      data={data}
      getRowId={documentsTablePreset.rowId}
      searchable={documentsTablePreset.capabilities.search === true}
      rowCount={total}
      loading={loading}
      error={error}
      onRetry={retry}
      toolbar={<button className="dt-page-btn" type="button">+ Новый документ</button>}
      {...tableState}
    />
  );
}

export default DocumentsPage;

