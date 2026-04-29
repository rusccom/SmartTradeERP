import { useEffect, useState } from "react";

import { createProduct } from "../api/createProduct";
import { createVariant } from "../api/createVariant";
import {
  addProductVariant,
  createProductFormState,
  patchProductForm,
  patchProductVariant,
  readPendingVariantPayloads,
  removeProductVariant,
  toCreateProductPayload,
} from "./productForm";

export function useProductCreateForm({ open, onClose, onCreated, t }) {
  const [form, setForm] = useState(createProductFormState);
  const [error, setError] = useState("");
  const [isSaving, setSaving] = useState(false);
  const [createdProductId, setCreatedProductId] = useState("");
  const [savedVariantCount, setSavedVariantCount] = useState(0);
  useResetOnClose({ open, setForm, setError, setSaving, setCreatedProductId, setSavedVariantCount });
  return {
    form,
    error,
    isSaving,
    lockedVariantCount: savedVariantCount,
    hasPendingProduct: createdProductId !== "",
    handleAddVariant: () => setForm(addProductVariant),
    handleChange: (event) => setForm((prev) => patchProductForm(prev, event)),
    handleClose: onClose,
    handleRemoveVariant: (variantId) => setForm((prev) => createdProductId ? prev : removeProductVariant(prev, variantId)),
    handleSubmit: createSubmitHandler({ createdProductId, form, onClose, onCreated, savedVariantCount, setCreatedProductId, setError, setSavedVariantCount, setSaving, t }),
    handleVariantChange: (variantId, event) => setForm((prev) => patchProductVariant(prev, variantId, event)),
  };
}

function createSubmitHandler(params) {
  return async function handleSubmit(event) {
    event.preventDefault();
    await submitProduct(params);
  };
}

async function submitProduct(params) {
  params.setError("");
  params.setSaving(true);
  try {
    await saveProduct(params);
  } catch (error) {
    params.setError(error.partial ? params.t("products.form.partialError") : error.message);
  } finally {
    params.setSaving(false);
  }
}

async function saveProduct(params) {
  if (params.form.variantMode === "multiple") {
    await saveMultiVariantProduct(params);
    return;
  }
  await saveSingleVariantProduct(params);
}

async function saveMultiVariantProduct(params) {
  const productId = await ensureProductCreated(params);
  try {
    await savePendingVariants(params, productId);
  } catch (error) {
    throw markAsPartial(error);
  }
  finishMultiVariantSave(params);
}

async function savePendingVariants(params, productId) {
  let savedCount = params.savedVariantCount || 1;
  const payloads = readPendingVariantPayloads(params.form, productId, savedCount);
  for (const payload of payloads) {
    await createVariant(payload);
    savedCount += 1;
    params.setSavedVariantCount(savedCount);
  }
}

async function saveSingleVariantProduct(params) {
  await createProduct(toCreateProductPayload(params.form));
  params.onCreated?.();
  params.onClose();
}

async function ensureProductCreated(params) {
  if (params.createdProductId) return params.createdProductId;
  const data = await createProduct(toCreateProductPayload(params.form));
  const productId = readCreatedProductId(data);
  params.setCreatedProductId(productId);
  params.setSavedVariantCount(1);
  params.onCreated?.();
  return productId;
}

function finishMultiVariantSave(params) {
  params.setCreatedProductId("");
  params.setSavedVariantCount(0);
  params.onClose();
}

function markAsPartial(error) {
  const reason = error instanceof Error ? error : new Error("Failed to create variants");
  reason.partial = true;
  return reason;
}

function readCreatedProductId(data) {
  const productId = data?.id;
  if (productId) return productId;
  throw new Error("Product ID was not returned");
}

function resetForm(state) {
  state.setForm(createProductFormState());
  state.setError("");
  state.setSaving(false);
  state.setCreatedProductId("");
  state.setSavedVariantCount(0);
}

function useResetOnClose(state) {
  const { open, setCreatedProductId, setError, setForm, setSavedVariantCount, setSaving } = state;
  useEffect(() => {
    if (open) return;
    resetForm({ setForm, setError, setSaving, setCreatedProductId, setSavedVariantCount });
  }, [open, setCreatedProductId, setError, setForm, setSavedVariantCount, setSaving]);
}
