import { useI18n } from "@smarterp/i18n/useI18n";
import Modal from "@smarterp/ui/modal/Modal";

export default function ConfirmKsef({ number, environment, onCancel, onConfirm }) {
  const { t } = useI18n();
  return <Modal open closeOnEscape title={t("entry.confirmSend")} onClose={onCancel} closeLabel={t("entry.close")}>
    <p>{t("entry.sendQuestion", { number, env: t(`entry.${environment}`) })}</p>
    <div className="order-actions">
      <button autoFocus type="button" className="order-btn" onClick={onCancel}>{t("entry.no")}</button>
      <button type="button" className="order-btn order-primary" onClick={onConfirm}>{t("entry.yes")}</button>
    </div>
  </Modal>;
}
