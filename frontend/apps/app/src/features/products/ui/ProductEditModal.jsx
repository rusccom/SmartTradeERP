import { useI18n } from "../../../shared/i18n/useI18n";
import Modal from "../../../shared/ui/modal/Modal";
import { useProductEditForm } from "../model/useProductEditForm";
import ProductCreateForm from "./ProductCreateForm";
import "./product-create.css";

function ProductEditModal({ onClose, onSaved, open, product }) {
  const { t } = useI18n();
  const state = useProductEditForm({ open, onClose, onSaved, product });
  if (!open || !product) return null;
  return (
    <Modal open={open} onClose={state.handleClose} closeLabel={t("common.close")} title={t("products.editModal.title")}>
      <ProductCreateForm state={state} t={t} />
    </Modal>
  );
}

export default ProductEditModal;
