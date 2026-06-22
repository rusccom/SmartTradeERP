import { useEffect, useState } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import FormModal from "../../../shared/ui/form-modal/FormModal";
import { decimalStep } from "../../currencies/model/formatMoney";
import { useCurrencies } from "../../currencies/model/useCurrencies";
import { createBundle } from "../api/createBundle";
import { createBundleFormState, patchBundleForm, readBundleSections, toCreateBundlePayload } from "../model/bundleForm";

function BundleCreateModal({ onClose, onCreated, open }) {
  const { t } = useI18n();
  const { defaultCurrency } = useCurrencies();
  const state = useBundleCreateState({ onClose, onCreated, open });
  return <FormModal {...readModalProps(t, state, open, decimalStep(defaultCurrency))} />;
}

function readModalProps(t, state, open, priceStep) {
  return {
    t,
    form: state.form,
    error: state.error,
    open,
    onChange: state.handleChange,
    onClose: state.handleClose,
    onSubmit: state.handleSubmit,
    closeLabel: t("common.close"),
    cancelLabel: t("common.cancel"),
    description: t("bundles.createModal.description"),
    isSubmitting: state.isSaving,
    sections: readBundleSections(priceStep),
    submitLabel: t("bundles.form.save"),
    submittingLabel: t("bundles.form.saving"),
    title: t("bundles.createModal.title"),
  };
}

function useBundleCreateState({ onClose, onCreated, open }) {
  const [form, setForm] = useState(createBundleFormState);
  const [error, setError] = useState("");
  const [isSaving, setSaving] = useState(false);
  useResetOnClose({ open, setError, setForm, setSaving });
  return {
    form,
    error,
    isSaving,
    handleChange: (event) => setForm((prev) => patchBundleForm(prev, event)),
    handleClose: onClose,
    handleSubmit: (event) => submitBundle(event, { form, onClose, onCreated, setError, setSaving }),
  };
}

async function submitBundle(event, state) {
  event.preventDefault();
  state.setError("");
  state.setSaving(true);
  try {
    await createBundle(toCreateBundlePayload(state.form));
    state.onCreated?.();
    state.onClose();
  } catch (error) {
    state.setError(error.message);
  } finally {
    state.setSaving(false);
  }
}

function useResetOnClose({ open, setError, setForm, setSaving }) {
  useEffect(() => {
    if (open) return;
    setForm(createBundleFormState());
    setError("");
    setSaving(false);
  }, [open, setError, setForm, setSaving]);
}

export default BundleCreateModal;
