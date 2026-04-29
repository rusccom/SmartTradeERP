function ProductCreateActions({ isSaving, onClose, state, t }) {
  return (
    <div className="product-create-actions">
      <button className="product-create-secondary" type="button" onClick={onClose}>{t("common.cancel")}</button>
      <button className="product-create-primary" type="submit" disabled={isSaving}>{readSubmitLabel(isSaving, state, t)}</button>
    </div>
  );
}

function readSubmitLabel(isSaving, state, t) {
  const key = isSaving ? state?.savingLabelKey : state?.submitLabelKey;
  return t(key || readDefaultLabelKey(isSaving));
}

function readDefaultLabelKey(isSaving) {
  return isSaving ? "products.form.saving" : "products.form.save";
}

export default ProductCreateActions;
