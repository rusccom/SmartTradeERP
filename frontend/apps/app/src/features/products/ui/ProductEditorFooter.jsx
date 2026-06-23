function ProductEditorFooter({ formId, state, t }) {
  return (
    <div className="product-editor-footer">
      {state.error ? <p className="product-create-error">{state.error}</p> : null}
      <div className="shared-form-modal-actions">
        <button className="shared-form-modal-secondary" type="button" onClick={state.handleClose}>{t("common.cancel")}</button>
        <button className="shared-form-modal-primary" type="submit" form={formId} disabled={state.isSaving}>{readSubmitLabel(state, t)}</button>
      </div>
    </div>
  );
}

function readSubmitLabel(state, t) {
  const key = state.isSaving ? state.savingLabelKey : state.submitLabelKey;
  return t(key || readDefaultLabelKey(state.isSaving));
}

function readDefaultLabelKey(isSaving) {
  return isSaving ? "products.form.saving" : "products.form.save";
}

export default ProductEditorFooter;
