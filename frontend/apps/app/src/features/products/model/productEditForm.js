export function createProductEditForm(product) {
  const variants = readVariants(product).map(toFormVariant);
  const first = variants[0] || createEmptyVariant();
  if (variants.length > 1) {
    return createMultiForm(product, variants);
  }
  return createSingleForm(product, first, variants);
}

export function isDraftVariant(variant) {
  return String(variant.id || "").startsWith("variant-");
}

export function toProductUpdatePayload(form) {
  return { name: form.name.trim() };
}

export function toVariantCreatePayload(productId, variant) {
  return { product_id: productId, ...toVariantUpdatePayload(variant) };
}

export function toVariantUpdatePayload(variant) {
  return {
    name: readVariantName(variant),
    sku_code: variant.skuCode.trim(),
    barcode: variant.barcode.trim(),
    unit: variant.unit.trim(),
    price: Number(variant.price) || 0,
    option1: readOption(variant.option1),
    option2: readOption(variant.option2),
    option3: readOption(variant.option3),
  };
}

export function readSingleVariant(form) {
  const variant = form.variants[0] || createEmptyVariant();
  return { ...variant, name: variant.name || "Default", skuCode: form.skuCode, barcode: form.barcode, unit: form.unit, price: form.price };
}

function createMultiForm(product, variants) {
  return {
    name: product?.name || "",
    unit: readSharedUnit(variants),
    price: variants[0]?.price || "0",
    skuCode: "",
    barcode: "",
    variantMode: "multiple",
    variants,
  };
}

function createSingleForm(product, variant, variants) {
  return {
    name: product?.name || "",
    unit: variant.unit,
    price: variant.price,
    skuCode: variant.skuCode,
    barcode: variant.barcode,
    variantMode: "single",
    variants,
  };
}

function createEmptyVariant() {
  return { id: "", name: "Default", skuCode: "", barcode: "", unit: "pcs", price: "0", option1: "", option2: "", option3: "" };
}

function toFormVariant(variant) {
  return {
    id: variant.id,
    name: variant.name || "Default",
    skuCode: variant.sku_code || "",
    barcode: variant.barcode || "",
    unit: variant.unit || "pcs",
    price: String(variant.price ?? "0"),
    option1: variant.option1 || "",
    option2: variant.option2 || "",
    option3: variant.option3 || "",
  };
}

function readVariantName(variant) {
  return (variant.name || "Default").trim() || "Default";
}

function readOption(value) {
  return String(value || "").trim();
}

function readSharedUnit(variants) {
  return variants[0]?.unit || "pcs";
}

function readVariants(product) {
  return Array.isArray(product?.variants) ? product.variants : [];
}
