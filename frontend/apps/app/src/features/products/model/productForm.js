const defaultForm = {
  name: "",
  unit: "pcs",
  price: "0",
  skuCode: "",
  barcode: "",
  variantMode: "single",
  variants: [],
};

function variantFields(priceStep = "0.01") {
  return [
    { name: "name", labelKey: "products.form.variantName", type: "text", required: true },
    { name: "skuCode", labelKey: "products.form.skuCode", type: "text" },
    { name: "barcode", labelKey: "products.form.barcode", type: "text" },
    { name: "price", labelKey: "products.form.price", type: "number", min: "0", step: priceStep, required: true },
  ];
}

let draftCounter = 0;

export function addProductVariant(form) {
  if (form.variantMode !== "multiple") return toMultipleMode(form);
  return { ...form, variants: [...form.variants, createVariantDraft(readNewVariantSeed(form))] };
}

export function createProductFormState() {
  return { ...defaultForm, variants: [] };
}

export function patchProductForm(form, event) {
  const { checked, name, type, value } = event.target;
  const next = { ...form, [name]: type === "checkbox" ? checked : value };
  return syncCommonVariantFields(next, name);
}

export function patchProductVariant(form, variantId, event) {
  return { ...form, variants: mapVariants(form.variants, variantId, event) };
}

export function readPendingVariantPayloads(form, productId, startIndex) {
  return form.variants.slice(startIndex).map((item) => toVariantPayload(item, productId));
}

export function readVariantFields(priceStep) {
  return variantFields(priceStep);
}

export function removeProductVariant(form, variantId) {
  if (form.variants.length <= 1) return toSingleMode(form);
  return { ...form, variants: form.variants.filter((item) => item.id !== variantId) };
}

export function toCreateProductPayload(form) {
  const variant = readPrimaryVariant(form);
  return {
    name: form.name.trim(),
    variant_name: readPrimaryVariantName(form),
    unit: variant.unit.trim(),
    price: Number(variant.price) || 0,
    sku_code: variant.skuCode.trim(),
    barcode: variant.barcode.trim(),
  };
}

function createSeedVariants(form) {
  return [createVariantDraft(readSingleVariantSeed(form))];
}

function createVariantDraft(seed = {}) {
  return {
    id: readDraftKey(),
    name: "",
    skuCode: "",
    barcode: "",
    unit: "pcs",
    price: "0",
    option1: "",
    option2: "",
    option3: "",
    ...seed,
  };
}

function mapVariants(variants, variantId, event) {
  return variants.map((item) => item.id === variantId ? patchVariantItem(item, event) : item);
}

function patchVariantItem(item, event) {
  const { checked, name, type, value } = event.target;
  return { ...item, [name]: type === "checkbox" ? checked : value };
}

function readDraftKey() {
  draftCounter += 1;
  return `variant-${draftCounter}`;
}

function readPrimaryVariant(form) {
  return form.variantMode === "multiple" ? form.variants[0] || createVariantDraft() : form;
}

function readPrimaryVariantName(form) {
  if (form.variantMode !== "multiple") return "Default";
  return readPrimaryVariant(form).name.trim() || "Default";
}

function readSingleVariantSeed(form) {
  return {
    name: "",
    skuCode: form.skuCode,
    barcode: form.barcode,
    unit: form.unit,
    price: form.price,
  };
}

function readNewVariantSeed(form) {
  return { unit: form.unit, price: form.price };
}

function syncSingleVariant(form) {
  const first = form.variants[0];
  if (!first) return form;
  return { ...form, unit: first.unit, price: first.price, skuCode: first.skuCode, barcode: first.barcode };
}

function syncCommonVariantFields(form, name) {
  if (name !== "unit" || form.variantMode !== "multiple") return form;
  return { ...form, variants: form.variants.map((item) => ({ ...item, unit: form.unit })) };
}

function toMultipleMode(form) {
  if (!form.variants.length) {
    return { ...form, variantMode: "multiple", variants: createSeedVariants(form) };
  }
  return { ...form, variantMode: "multiple", variants: [...form.variants, createVariantDraft(readNewVariantSeed(form))] };
}

function toSingleMode(form) {
  return { ...syncSingleVariant(form), variantMode: "single" };
}

function toVariantPayload(item, productId) {
  return {
    product_id: productId,
    name: item.name.trim(),
    sku_code: item.skuCode.trim(),
    barcode: item.barcode.trim(),
    unit: item.unit.trim(),
    price: Number(item.price) || 0,
    option1: item.option1.trim(),
    option2: item.option2.trim(),
    option3: item.option3.trim(),
  };
}
