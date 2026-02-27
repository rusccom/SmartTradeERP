import { apiPaths } from "../../../shared/api/client";
import { getJSON } from "../../../shared/api/http";
import { useServerDataTable } from "../../../shared/model/data-table/useServerDataTable";
import DataTable from "../../../shared/ui/data-table/DataTable";
import { productsTablePreset } from "./table/productsTablePreset";

function ProductsPage() {
  const { data, total, loading, error, retry, tableState } = useServerDataTable(productsTablePreset);
  return (
    <DataTable
      columns={productsTablePreset.columns}
      data={data}
      getRowId={productsTablePreset.rowId}
      searchable={productsTablePreset.capabilities.search === true}
      rowCount={total}
      loading={loading}
      error={error}
      onRetry={retry}
      expandable={true}
      getSubRows={loadProductVariants}
      toolbar={<button className="dt-page-btn" type="button">+ Добавить товар</button>}
      {...tableState}
    />
  );
}

function loadProductVariants(row) {
  return getJSON(apiPaths.variants, { product_id: row.id }).then((response) => response.data || []);
}

export default ProductsPage;

