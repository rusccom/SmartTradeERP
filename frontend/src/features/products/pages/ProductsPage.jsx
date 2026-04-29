import { useMemo, useState } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import { ServerListTable } from "../../../shared/ui/tables/list-table";
import { loadProductVariants } from "../api/loadProductVariants";
import { createProductsTablePreset } from "../model/productsTablePreset";
import ProductCreateModal from "../ui/ProductCreateModal";

function ProductsPage() {
  const { t } = useI18n();
  const [createOpen, setCreateOpen] = useState(false);
  const preset = useMemo(() => createProductsTablePreset(t), [t]);
  return (
    <ServerListTable
      preset={preset}
      expandable={true}
      getSubRows={loadProductVariants}
      primaryAction={createProductAction(t, setCreateOpen)}
    >
      {({ retry }) => (
        <ProductCreateModal open={createOpen} onClose={() => setCreateOpen(false)} onCreated={retry} />
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

export default ProductsPage;
