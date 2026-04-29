function ProductCreateActions({ isSaving, onClose, t }) {
  return (
    <div className="product-create-actions">
      <button className="product-create-secondary" type="button" onClick={onClose}>{t("common.cancel")}</button>
      <button className="product-create-primary" type="submit" disabled={isSaving}>{readSubmitLabel(isSaving, t)}</button>
    </div>
  );
}

function readSubmitLabel(isSaving, t) {
  return isSaving ? t("products.form.saving") : t("products.form.save");
}

export default ProductCreateActions;
