import { useI18n } from "@smarterp/i18n/useI18n";
import Modal from "@smarterp/ui/modal/Modal";
import { useProductCreateForm } from "../model/useProductCreateForm";
import ProductCreateForm from "./ProductCreateForm";
import "./product-create.css";

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
    >
      <ProductCreateForm state={state} t={t} />
    </Modal>
  );
}

export default ProductCreateModal;
