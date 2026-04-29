import { useMemo, useState } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import { ServerListTable } from "../../../shared/ui/tables/list-table";
import { loadProductVariants } from "../api/loadProductVariants";
import { createProductsTablePreset } from "../model/productsTablePreset";
import ProductCatalogEditModal from "../ui/ProductCatalogEditModal";
import ProductCreateModal from "../ui/ProductCreateModal";

function ProductsPage() {
  const { t } = useI18n();
  const [createOpen, setCreateOpen] = useState(false);
  const [editTarget, setEditTarget] = useState(null);
  const preset = useMemo(() => createProductsTablePreset(t), [t]);
  return (
    <ServerListTable
      preset={preset}
      selectable={true}
      onRowOpen={(row) => setEditTarget({ type: "product", data: row })}
      subRows={readSubRowsConfig(setEditTarget)}
      primaryAction={createProductAction(t, setCreateOpen)}
    >
      {({ retry }) => (
        <>
          <ProductCreateModal open={createOpen} onClose={() => setCreateOpen(false)} onCreated={retry} />
          <ProductCatalogEditModal target={editTarget} open={Boolean(editTarget)} onClose={() => setEditTarget(null)} onSaved={retry} />
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

function readSubRowsConfig(setEditTarget) {
  return {
    enabled: true,
    getRows: loadProductVariants,
    canExpand: hasMultipleVariants,
    onRowOpen: (row) => setEditTarget({ type: "variant", data: row }),
  };
}

function hasMultipleVariants(row) {
  return Array.isArray(row.variants) && row.variants.length > 1;
}

export default ProductsPage;
