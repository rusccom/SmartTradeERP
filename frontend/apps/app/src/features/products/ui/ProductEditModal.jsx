import { useI18n } from "@smarterp/i18n/useI18n";
import Modal from "@smarterp/ui/modal/Modal";
import { useProductEditForm } from "../model/useProductEditForm";
import ProductCreateForm from "./ProductCreateForm";
import ProductEditorFooter from "./ProductEditorFooter";
import "./product-create.css";

const FORM_ID = "product-edit-form";

function ProductEditModal({ onClose, onSaved, open, product }) {
  const { t } = useI18n();
  const state = useProductEditForm({ open, onClose, onSaved, product });
  if (!open || !product) return null;
  return (
    <Modal open={open} onClose={state.handleClose} closeLabel={t("common.close")} title={t("products.editModal.title")} description={product.name} size="lg" bodyTone="muted" footer={<ProductEditorFooter formId={FORM_ID} state={state} t={t} />}>
      <ProductCreateForm formId={FORM_ID} state={state} t={t} />
    </Modal>
  );
}

export default ProductEditModal;
