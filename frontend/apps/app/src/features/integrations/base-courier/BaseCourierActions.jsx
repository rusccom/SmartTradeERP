import { deleteJSON, postJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";
import { useI18n } from "@smarterp/i18n/useI18n";

export default function BaseCourierActions({ connection, action }) {
  const { t } = useI18n();
  const { email, version } = connection;
  function check() {
    action.run(() => postJSON(apiPaths.baseCourierCheck(version), {}), "baseCourier.verified");
  }
  function disconnect() {
    if (window.confirm(t("baseCourier.confirmDisconnect", { email }))) {
      action.run(() => deleteJSON(apiPaths.baseCourierConnection(version)), "baseCourier.disconnected");
    }
  }
  return (
    <div className="base-courier-actions">
      <button type="button" disabled={action.busy} onClick={check}>{t("baseCourier.check")}</button>
      <button type="button" disabled={action.busy} onClick={disconnect}>{t("baseCourier.disconnect")}</button>
    </div>
  );
}
