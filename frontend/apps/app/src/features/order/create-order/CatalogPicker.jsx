import { useMemo } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";
import { apiPaths } from "@smarterp/api/publicApi";
import { createApiTablePreset } from "@smarterp/model/data-table/createApiTablePreset";
import { ServerListTable } from "@smarterp/ui/tables/list-table";
import Modal from "@smarterp/ui/modal/Modal";
import ProductTableThumb from "../../products/ui/ProductTableThumb";
import { productSearchFilter } from "../../products/model/productSearchFilter";
import { catalogLine } from "./entryModel";

export default function CatalogPicker({ onClose, onPick }) {
  const { t } = useI18n();
  const preset = useMemo(() => createApiTablePreset({
    id: "order-catalog", path: apiPaths.products, rowId: (r) => r.id,
    capabilities: { sorting: true, search: true }, sorting: ["name"], search: productSearchFilter,
    mapStateToQuery: () => ({ include: "variants" }), columns: pickerColumns(t, onPick),
  }), [t, onPick]);
  return <Modal open onClose={onClose} closeOnEscape closeLabel={t("entry.close")} title={t("entry.addProduct")} size="lg">
    <div className="order-picker"><ServerListTable preset={preset} emptyText={t("entry.noProducts")} /></div>
  </Modal>;
}

function pickerColumns(t, onPick) {
  return [
    { accessorKey: "name", header: t("entry.product"), cell: (value, row) =>
      <span className="order-product"><ProductTableThumb name={value} url={row.image_thumb_url} />{value}</span> },
    { accessorKey: "variants", header: t("entry.variant"), cell: (variants, product) =>
      <div className="order-pick-options">{variants?.map((v) =>
        <button type="button" className="order-btn" key={v.id} disabled={v.price == null}
          onClick={() => onPick(catalogLine(product, v))}>{v.name || product.name}{v.sku_code && ` · ${v.sku_code}`}</button>)}</div> },
  ];
}
