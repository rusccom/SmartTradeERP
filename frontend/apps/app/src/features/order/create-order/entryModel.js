const emptyAddress = { name: "", company: "", street: "", postal_code: "", city: "", region: "", country: "", tax_id: "" };
const emptyPickup = { name: "", point_id: "", street: "", postal_code: "", city: "", country: "" };

export function entryFrom(card) {
  const d = card?.details ?? {};
  return {
    id: card?.id, request_id: card?.id ?? crypto.randomUUID(), revision: card?.revision ?? 0,
    customer_name: card?.customer_name ?? "", customer_email: card?.customer_email ?? "", customer_phone: card?.customer_phone ?? "",
    billing_customer_id: card?.billing_customer_id ?? "", invoice_requested: Boolean(card?.invoice_requested),
    delivery_cost: card?.delivery_cost ?? "", items: (card?.items ?? []).map((l) => ({ ...l, key: l.id })),
    details: { ...d, nickname: d.nickname ?? "", delivery_method: d.delivery_method ?? "", delivery_kind: d.delivery_kind || "address",
      delivery_vat: d.delivery_vat ?? "", payment_method: d.payment_method ?? "", cod: Boolean(d.cod),
      shipping: { ...emptyAddress, ...d.shipping }, billing: { ...emptyAddress, ...d.billing }, pickup: { ...emptyPickup, ...d.pickup } },
  };
}

export function lineInputs(items) {
  return items.map(({ id, variant_id, qty, vat_rate, vat_mode, vat_country }) => ({ id, variant_id, qty,
    vat_mode, ...(vat_mode === "manual" ? { vat_rate, vat_country } : {}) }));
}

export function entryPayload(form) {
  const { id, items, details, delivery_cost, ...entry } = form;
  return { ...entry, items: lineInputs(items), delivery_cost: delivery_cost === "" ? null : delivery_cost,
    details: { ...details, pickup: details.delivery_kind === "pickup" ? details.pickup : null } };
}

export function patchEntry(source, path, value) {
  const [key, ...rest] = path.split(".");
  if (!rest.length) return { ...source, [key]: value };
  return { ...source, [key]: patchEntry(source[key], rest.join("."), value) };
}

export function catalogLine(product, variant) {
  return { key: crypto.randomUUID(), variant_id: variant.id, sku_code: variant.sku_code, unit: variant.unit,
    product_name: product.name, variant_name: variant.name, unit_price: variant.price, qty: "1", vat_rate: "", vat_mode: "catalog", vat_country: "",
    image_thumb_url: variant.image_thumb_url || product.image_thumb_url };
}

export function displayMoney(value, code, locale) {
  if (value == null || value === "") return "—";
  const format = code ? { style: "currency", currency: code, maximumFractionDigits: 4 } : { maximumFractionDigits: 4 };
  return new Intl.NumberFormat(locale, format).format(value);
}
