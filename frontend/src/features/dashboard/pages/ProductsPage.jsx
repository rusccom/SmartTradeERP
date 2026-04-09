import { useMemo } from "react";

import { apiPaths } from "../../../shared/api/client";
import { getJSON } from "../../../shared/api/http";
import { useI18n } from "../../../shared/i18n/useI18n";
import { useServerDataTable } from "../../../shared/model/data-table/useServerDataTable";
import DataTable from "../../../shared/ui/data-table/DataTable";
import { createProductsTablePreset } from "./table/productsTablePreset";

function ProductsPage() {
  const { t } = useI18n();
  const preset = useMemo(() => createProductsTablePreset(t), [t]);
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
      expandable={true}
      getSubRows={loadProductVariants}
      toolbar={<button className="dt-page-btn" type="button">+ {t("products.addButton")}</button>}
      {...tableState}
    />
  );
}

function loadProductVariants(row) {
  return getJSON(apiPaths.variants, { product_id: row.id }).then((response) => response.data || []);
}

export default ProductsPage;
