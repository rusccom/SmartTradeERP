const defaultForm = {
  name: "",
  unit: "pcs",
  price: "0",
  skuCode: "",
  barcode: "",
};

export function createBundleFormState() {
  return { ...defaultForm };
}

export function patchBundleForm(form, event) {
  const { name, value } = event.target;
  return { ...form, [name]: value };
}

export function readBundleSections(priceStep = "0.01") {
  return [{
    id: "bundle",
    titleKey: "bundles.form.sections.bundle",
    fields: bundleFields(priceStep),
  }];
}

export function toCreateBundlePayload(form) {
  return {
    name: form.name.trim(),
    unit: form.unit.trim(),
    price: Number(form.price) || 0,
    sku_code: form.skuCode.trim(),
    barcode: form.barcode.trim(),
  };
}

function bundleFields(priceStep) {
  return [
    { name: "name", labelKey: "bundles.form.name", type: "text", required: true, autoFocus: true },
    { name: "unit", labelKey: "bundles.form.unit", type: "text", required: true },
    { name: "price", labelKey: "bundles.form.price", type: "number", min: "0", step: priceStep, required: true },
    { name: "skuCode", labelKey: "bundles.form.skuCode", type: "text" },
    { name: "barcode", labelKey: "bundles.form.barcode", type: "text" },
  ];
}
