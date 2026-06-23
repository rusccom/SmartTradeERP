import { useI18n } from "@smarterp/i18n/useI18n";
import Modal from "@smarterp/ui/modal/Modal";
import { useProductCreateForm } from "../model/useProductCreateForm";
import ProductCreateForm from "./ProductCreateForm";
import ProductEditorFooter from "./ProductEditorFooter";
import "./product-create.css";

const FORM_ID = "product-create-form";

function ProductCreateModal({ onClose, onCreated, open }) {
  const { t } = useI18n();
  const state = useProductCreateForm({ open, onClose, onCreated, t });
  return (
    <Modal
      open={open}
      onClose={state.handleClose}
      closeLabel={t("common.close")}
      title={t("products.createModal.title")}
      size="lg"
      bodyTone="muted"
      footer={<ProductEditorFooter formId={FORM_ID} state={state} t={t} />}
    >
      <ProductCreateForm formId={FORM_ID} state={state} t={t} />
    </Modal>
  );
}

export default ProductCreateModal;
