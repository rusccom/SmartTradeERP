import { useI18n } from "../../../shared/i18n/useI18n";
import FormModal from "../../../shared/ui/form-modal/FormModal";
import { readProductSections } from "../model/productForm";
import { useProductCreateForm } from "../model/useProductCreateForm";
import ProductVariantModeSwitch from "./ProductVariantModeSwitch";
import ProductVariantsBuilder from "./ProductVariantsBuilder";

function ProductCreateModal({ onClose, onCreated, open }) {
  const { t } = useI18n();
  const state = useProductCreateForm({ open, onClose, onCreated, t });
  return (
    <FormModal
      t={t}
      form={state.form}
      error={state.error}
      open={open}
      onChange={state.handleChange}
      onClose={state.handleClose}
      onSubmit={state.handleSubmit}
      closeLabel={t("common.close")}
      cancelLabel={t("common.cancel")}
      description={t("products.createModal.description")}
      isSubmitting={state.isSaving}
      sections={readProductSections(state.form.variantMode)}
      submitLabel={t("products.form.save")}
      submittingLabel={t("products.form.saving")}
      topSlot={<ProductVariantModeSwitch mode={state.form.variantMode} onChange={state.handleModeChange} disabled={state.hasPendingProduct} t={t} />}
      bottomSlot={readBottomSlot(state, t)}
      title={t("products.createModal.title")}
    />
  );
}

function readBottomSlot(state, t) {
  if (state.form.variantMode !== "multiple") return null;
  return <ProductVariantsBuilder variants={state.form.variants} lockedVariantCount={state.lockedVariantCount} onAddVariant={state.handleAddVariant} onRemoveVariant={state.handleRemoveVariant} onVariantChange={state.handleVariantChange} canRemove={state.form.variants.length > 2} disabled={state.hasPendingProduct} t={t} />;
}

export default ProductCreateModal;
