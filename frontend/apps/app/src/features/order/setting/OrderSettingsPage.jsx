import { useRef } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";
import { useStatusDraft } from "./useStatusDraft";
import StatusDraftForm from "./StatusDraftForm";
import StatusDraftTable from "./StatusDraftTable";
import "./setting.css";

export default function OrderSettingsPage() {
  const { t, locale } = useI18n();
  const draft = useStatusDraft(locale);
  const input = useRef(null);
  function startEditing(status) {
    draft.edit(status);
    input.current?.focus();
  }
  return (
    <section className="order-settings" aria-labelledby="order-settings-title">
      <h1 id="order-settings-title">{t("orders.settings.title")}</h1>
      <div className="order-status-panel">
        <h2>{t("orders.statuses.title")}</h2>
        <p>{t("orders.statuses.draftHint")}</p>
        <StatusDraftForm draft={draft} inputRef={input} t={t} />
        <StatusDraftTable statuses={draft.statuses} onEdit={startEditing} />
      </div>
    </section>
  );
}
