import { useMemo, useState } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import { useServerDataTable } from "../../../shared/model/data-table/useServerDataTable";
import DataTable from "../../../shared/ui/data-table/DataTable";
import { loadProductVariants } from "../api/loadProductVariants";
import { createProductsTablePreset } from "../model/productsTablePreset";
import ProductCreateModal from "../ui/ProductCreateModal";

function ProductsPage() {
  const { t } = useI18n();
  const [createOpen, setCreateOpen] = useState(false);
  const preset = useMemo(() => createProductsTablePreset(t), [t]);
  const { data, error, loading, retry, tableState, total } = useServerDataTable(preset);
  return (
    <>
      <DataTable columns={preset.columns} data={data} getRowId={preset.rowId} searchable={preset.capabilities.search === true} rowCount={total} loading={loading} error={error} onRetry={retry} expandable={true} getSubRows={loadProductVariants} toolbar={renderToolbar(t, setCreateOpen)} {...tableState} />
      <ProductCreateModal open={createOpen} onClose={() => setCreateOpen(false)} onCreated={retry} />
    </>
  );
}

function renderToolbar(t, setCreateOpen) {
  return <button className="dt-page-btn" type="button" onClick={() => setCreateOpen(true)}>+ {t("products.addButton")}</button>;
}

export default ProductsPage;
