const defaultForm = {
  name: "",
  unit: "pcs",
  price: "0",
  skuCode: "",
  barcode: "",
  variantMode: "single",
  variants: [],
};

const productSections = [
  {
    id: "product",
    titleKey: "products.form.sections.product",
    fields: [
      { name: "name", labelKey: "products.form.name", type: "text", required: true, autoFocus: true },
    ],
  },
];

const singleVariantSection = {
  id: "variant",
  titleKey: "products.form.sections.variant",
  fields: [
    { name: "unit", labelKey: "products.form.unit", type: "text", required: true },
    { name: "price", labelKey: "products.form.price", type: "number", min: "0", step: "0.01", required: true },
    { name: "skuCode", labelKey: "products.form.skuCode", type: "text" },
    { name: "barcode", labelKey: "products.form.barcode", type: "text" },
  ],
};

const variantFields = [
  { name: "name", labelKey: "products.form.variantName", type: "text", required: true },
  { name: "unit", labelKey: "products.form.unit", type: "text", required: true },
  { name: "price", labelKey: "products.form.price", type: "number", min: "0", step: "0.01", required: true },
  { name: "skuCode", labelKey: "products.form.skuCode", type: "text" },
  { name: "barcode", labelKey: "products.form.barcode", type: "text" },
  { name: "option1", labelKey: "products.form.option1", type: "text" },
  { name: "option2", labelKey: "products.form.option2", type: "text" },
  { name: "option3", labelKey: "products.form.option3", type: "text" },
];

let draftCounter = 0;

export function addProductVariant(form) {
  return { ...form, variants: [...form.variants, createVariantDraft()] };
}

export function changeProductVariantMode(form, mode) {
  if (mode === form.variantMode) return form;
  return mode === "multiple" ? toMultipleMode(form) : toSingleMode(form);
}

export function createProductFormState() {
  return { ...defaultForm, variants: [] };
}

export function patchProductForm(form, event) {
  const { checked, name, type, value } = event.target;
  return { ...form, [name]: type === "checkbox" ? checked : value };
}

export function patchProductVariant(form, variantId, event) {
  return { ...form, variants: mapVariants(form.variants, variantId, event) };
}

export function readPendingVariantPayloads(form, productId, startIndex) {
  return form.variants.slice(startIndex).map((item) => toVariantPayload(item, productId));
}

export function readProductSections(mode) {
  return mode === "multiple" ? productSections : [productSections[0], singleVariantSection];
}

export function readVariantFields() {
  return variantFields;
}

export function removeProductVariant(form, variantId) {
  if (form.variants.length <= 2) return form;
  return { ...form, variants: form.variants.filter((item) => item.id !== variantId) };
}

export function toCreateProductPayload(form) {
  const variant = readPrimaryVariant(form);
  return {
    name: form.name.trim(),
    unit: variant.unit.trim(),
    price: Number(variant.price) || 0,
    sku_code: variant.skuCode.trim(),
    barcode: variant.barcode.trim(),
  };
}

function createSeedVariants(form) {
  return [createVariantDraft(readSingleVariantSeed(form)), createVariantDraft()];
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

function readSingleVariantSeed(form) {
  return {
    name: form.name.trim(),
    skuCode: form.skuCode,
    barcode: form.barcode,
    unit: form.unit,
    price: form.price,
  };
}

function syncSingleVariant(form) {
  const first = form.variants[0];
  if (!first) return form;
  return { ...form, unit: first.unit, price: first.price, skuCode: first.skuCode, barcode: first.barcode };
}

function toMultipleMode(form) {
  return { ...form, variantMode: "multiple", variants: form.variants.length ? form.variants : createSeedVariants(form) };
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
