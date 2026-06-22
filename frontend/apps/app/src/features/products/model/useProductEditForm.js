import { useEffect, useState } from "react";

import { createVariant } from "../api/createVariant";
import { deleteVariant } from "../api/deleteVariant";
import { updateProduct } from "../api/updateProduct";
import { updateVariant } from "../api/updateVariant";
import { addProductVariant, patchProductForm, patchProductVariant, removeProductVariant } from "./productForm";
import {
  createProductEditForm,
  isDraftVariant,
  readSingleVariant,
  toProductUpdatePayload,
  toVariantCreatePayload,
  toVariantUpdatePayload,
} from "./productEditForm";

export function useProductEditForm({ open, onClose, onSaved, product }) {
  const state = useEditState(product);
  useProductFormLoad({ open, product, ...state });
  return createEditApi({ onClose, onSaved, product, ...state });
}

function useEditState(product) {
  const [form, setForm] = useState(() => createProductEditForm(product));
  const [error, setError] = useState("");
  const [isSaving, setSaving] = useState(false);
  const [removedVariantIds, setRemovedVariantIds] = useState([]);
  return { error, form, isSaving, removedVariantIds, setError, setForm, setRemovedVariantIds, setSaving };
}

function createEditApi(params) {
  const { error, form, isSaving, onClose, setForm } = params;
  return {
    form,
    error,
    isSaving,
    lockedVariantCount: 0,
    hasPendingProduct: false,
    submitLabelKey: "products.form.saveEdit",
    savingLabelKey: "products.form.savingEdit",
    handleAddVariant: () => setForm(addProductVariant),
    handleChange: (event) => setForm((prev) => patchProductForm(prev, event)),
    handleClose: onClose,
    handleRemoveVariant: (id) => removeVariant({ id, setForm, setRemovedVariantIds: params.setRemovedVariantIds }),
    handleSubmit: createSubmitHandler(params),
    handleVariantChange: (id, event) => setForm((prev) => patchProductVariant(prev, id, event)),
  };
}

function useProductFormLoad(params) {
  const { open, product, setError, setForm, setRemovedVariantIds, setSaving } = params;
  useEffect(() => {
    if (!open || !product) return;
    setForm(createProductEditForm(product));
    setError("");
    setSaving(false);
    setRemovedVariantIds([]);
  }, [open, product, setError, setForm, setRemovedVariantIds, setSaving]);
}

function createSubmitHandler(params) {
  return async function handleSubmit(event) {
    event.preventDefault();
    await submitProductEdit(params);
  };
}

async function submitProductEdit(params) {
  params.setError("");
  params.setSaving(true);
  try {
    await saveProductEdit(params);
    params.onSaved?.();
    params.onClose();
  } catch (error) {
    params.setError(error.message);
  } finally {
    params.setSaving(false);
  }
}

async function saveProductEdit(params) {
  const productId = params.product?.id;
  if (!productId) throw new Error("Product is not selected");
  await updateProduct(productId, toProductUpdatePayload(params.form));
  await saveVariantChanges({ productId, form: params.form, removedVariantIds: params.removedVariantIds });
}

async function saveVariantChanges(params) {
  const variants = readVariantsForSave(params.form);
  for (const variant of variants) {
    await saveVariant(params.productId, variant);
  }
  for (const id of params.removedVariantIds) {
    await deleteVariant(id);
  }
}

async function saveVariant(productId, variant) {
  if (isNewVariant(variant)) {
    await createVariant(toVariantCreatePayload(productId, variant));
    return;
  }
  await updateVariant(variant.id, toVariantUpdatePayload(variant));
}

function removeVariant({ id, setForm, setRemovedVariantIds }) {
  setForm((prev) => {
    const variant = prev.variants.find((item) => item.id === id);
    if (shouldDeleteVariant(prev, variant)) {
      setRemovedVariantIds((ids) => appendOnce(ids, id));
    }
    return removeProductVariant(prev, id);
  });
}

function shouldDeleteVariant(form, variant) {
  return Boolean(variant) && !isDraftVariant(variant) && form.variants.length > 1;
}

function readVariantsForSave(form) {
  return form.variantMode === "multiple" ? form.variants : [readSingleVariant(form)];
}

function isNewVariant(variant) {
  return !variant.id || isDraftVariant(variant);
}

function appendOnce(items, value) {
  return items.includes(value) ? items : [...items, value];
}
