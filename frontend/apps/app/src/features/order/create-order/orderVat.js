export function destination(order) {
  const details = order.details;
  return (details?.delivery_kind === "pickup" ? details.pickup?.country : details?.shipping?.country) || "";
}

export function resetDestinationRates(previous, next) {
  if (destination(previous) === destination(next)) return next;
  return { ...next, items: next.items.map((line) => ({ ...line, vat_mode: "catalog", vat_rate: "", vat_country: "" })) };
}

export function editOrderLine(line, field, value, form) {
  if (field !== "vat_rate") return { ...line, [field]: value };
  return { ...line, vat_rate: value, vat_mode: value === "" ? "catalog" : "manual", vat_country: destination(form) };
}
