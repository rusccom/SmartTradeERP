import { useEffect, useState } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import FormField from "../../../shared/ui/form-modal/FormField";
import Modal from "../../../shared/ui/modal/Modal";
import { updateProduct } from "../api/updateProduct";
import { updateVariant } from "../api/updateVariant";

const productFields = [
  { name: "name", labelKey: "products.form.name", type: "text", required: true, autoFocus: true },
];

const variantFields = [
  { name: "name", labelKey: "products.form.variantName", type: "text", required: true, autoFocus: true },
  { name: "skuCode", labelKey: "products.form.skuCode", type: "text" },
  { name: "barcode", labelKey: "products.form.barcode", type: "text" },
  { name: "unit", labelKey: "products.form.unit", type: "text", required: true },
  { name: "price", labelKey: "products.form.price", type: "number", min: "0", step: "0.01", required: true },
];

const defaultVariantFields = variantFields.filter((field) => field.name !== "name");

function ProductCatalogEditModal({ onClose, onSaved, open, target }) {
  const { t } = useI18n();
  const [form, setForm] = useState(() => readForm(target));
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  useEffect(() => setForm(readForm(target)), [target]);
  if (!open || !target) return null;
  return (
    <Modal open={open} onClose={onClose} closeLabel={t("common.close")} title={readTitle(target, t)}>
      <form className="product-create-form" onSubmit={(event) => handleSubmit(event, { form, onClose, onSaved, setError, setSaving, target })}>
        {renderFields(readFields(target), form, setFormValue(setForm), t)}
        {error && <p className="product-create-error">{error}</p>}
        {renderActions(saving, onClose, t)}
      </form>
    </Modal>
  );
}

function readFields(target) {
  if (target.type === "variant") return variantFields;
  return readDefaultVariant(target.data) ? [...productFields, ...defaultVariantFields] : productFields;
}

function renderFields(fields, form, onChange, t) {
  return (
    <section className="product-create-card">
      <div className="product-create-grid">
        {fields.map((field) => <FormField key={field.name} field={field} value={form[field.name]} onChange={onChange} t={t} />)}
      </div>
    </section>
  );
}

function renderActions(saving, onClose, t) {
  return (
    <div className="product-create-actions">
      <button className="product-create-secondary" type="button" onClick={onClose}>{t("common.cancel")}</button>
      <button className="product-create-primary" type="submit" disabled={saving}>{readSaveLabel(saving, t)}</button>
    </div>
  );
}

async function handleSubmit(event, params) {
  event.preventDefault();
  params.setError("");
  params.setSaving(true);
  try {
    await saveTarget(params.target, params.form);
    params.onSaved?.();
    params.onClose();
  } catch (error) {
    params.setError(error.message);
  } finally {
    params.setSaving(false);
  }
}

async function saveTarget(target, form) {
  if (target.type === "variant") {
    await updateVariant(target.data.id, toVariantPayload(form, target.data));
    return;
  }
  await updateProduct(target.data.id, { name: form.name.trim() });
  const variant = readDefaultVariant(target.data);
  if (variant) await updateVariant(variant.id, toDefaultVariantPayload(form, variant));
}

function toDefaultVariantPayload(form, variant) {
  return toVariantPayload({ ...form, name: variant.name || "Default" }, variant);
}

function toVariantPayload(form, variant) {
  return {
    name: form.name.trim(),
    sku_code: form.skuCode.trim(),
    barcode: form.barcode.trim(),
    unit: form.unit.trim(),
    price: Number(form.price) || 0,
    option1: variant.option1 || "",
    option2: variant.option2 || "",
    option3: variant.option3 || "",
  };
}

function readForm(target) {
  if (target?.type === "variant") return readVariantForm(target.data);
  const variant = readDefaultVariant(target?.data);
  return { ...readVariantForm(variant), name: target?.data?.name || "" };
}

function readVariantForm(variant) {
  return {
    name: variant?.name || "",
    skuCode: variant?.sku_code || "",
    barcode: variant?.barcode || "",
    unit: variant?.unit || "pcs",
    price: String(variant?.price ?? "0"),
  };
}

function readDefaultVariant(product) {
  const variants = Array.isArray(product?.variants) ? product.variants : [];
  return variants.length === 1 ? variants[0] : null;
}

function readTitle(target, t) {
  return t(target.type === "variant" ? "products.variantEditModal.title" : "products.editModal.title");
}

function readSaveLabel(saving, t) {
  return saving ? t("products.form.saving") : t("products.form.save");
}

function setFormValue(setForm) {
  return (event) => setForm((prev) => ({ ...prev, [event.target.name]: event.target.value }));
}

export default ProductCatalogEditModal;
