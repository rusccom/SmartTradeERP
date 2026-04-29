import { useMemo, useState } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import { ServerListTable } from "../../../shared/ui/tables/list-table";
import { loadProductVariants } from "../api/loadProductVariants";
import { createProductsTablePreset } from "../model/productsTablePreset";
import ProductCreateModal from "../ui/ProductCreateModal";
import ProductEditModal from "../ui/ProductEditModal";

function ProductsPage() {
  const { t } = useI18n();
  const [createOpen, setCreateOpen] = useState(false);
  const [editProduct, setEditProduct] = useState(null);
  const preset = useMemo(() => createProductsTablePreset(t), [t]);
  return (
    <ServerListTable
      preset={preset}
      selectable={true}
      onRowOpen={setEditProduct}
      subRows={readSubRowsConfig()}
      primaryAction={createProductAction(t, setCreateOpen)}
    >
      {({ retry }) => (
        <>
          <ProductCreateModal open={createOpen} onClose={() => setCreateOpen(false)} onCreated={retry} />
          <ProductEditModal product={editProduct} open={Boolean(editProduct)} onClose={() => setEditProduct(null)} onSaved={retry} />
        </>
      )}
    </ServerListTable>
  );
}

function createProductAction(t, setCreateOpen) {
  return {
    label: t("products.addButton"),
    onClick: () => setCreateOpen(true),
  };
}

function readSubRowsConfig() {
  return {
    enabled: true,
    getRows: loadProductVariants,
    canExpand: hasMultipleVariants,
  };
}

function hasMultipleVariants(row) {
  return Array.isArray(row.variants) && row.variants.length > 1;
}

export default ProductsPage;
