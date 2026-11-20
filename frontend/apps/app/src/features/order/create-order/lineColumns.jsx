import ProductTableThumb from "../../products/ui/ProductTableThumb";
import { displayMoney } from "./entryModel";

export function lineColumns({ t, locale, code, quote, edit, locked }) {
  const cell = (id, render) => ({ id, label: t(`entry.${id}`), cell: render });
  const money = (amount) => displayMoney(amount, code, locale);
  return [
    cell("photo", (l) => <ProductTableThumb url={l.image_thumb_url} name={l.product_name} />),
    cell("id", (l) => <span className="order-item-id" title={l.variant_id}>{l.sku_code || l.variant_id}</span>),
    cell("product", (l) => <span className="orders-cell"><strong>{l.product_name}</strong>{l.variant_name !== l.product_name && <small>{l.variant_name}</small>}</span>),
    cell("price", (l, i) => money(quote?.items[i]?.unit_price ?? l.unit_price)),
    cell("qty", (l) => numericInput(l, "qty", { edit, locked, label: t("entry.qty"), min: quantityStepBase(l.qty), step: "1", max: "100000" })),
    cell("vat", (l) => numericInput(l, "vat_rate", { edit, locked, label: t("entry.vat"), min: "0", step: "0.01", max: "100" })),
    cell("lineTotal", (_, i) => money(quote?.items[i]?.total_amount)),
  ];
}

function quantityStepBase(qty) {
  // Native arrows step from min. Keep the entered thousandths so 1.125 steps to 2.125.
  const thousandths = Math.round(Number(qty) * 1000) % 1000;
  return thousandths > 0 ? thousandths / 1000 : 1;
}

function numericInput(line, field, options) {
  return <input className="order-number" type="number" value={line[field]} required={field === "qty"} min={options.min}
    max={options.max} step={options.step} disabled={options.locked} aria-label={`${options.label} · ${line.product_name}`}
    onChange={(event) => options.edit(line.key, field, event.target.value)} />;
}
