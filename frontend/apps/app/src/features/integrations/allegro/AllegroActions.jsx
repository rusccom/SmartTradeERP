import { deleteJSON, postJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";
import { useI18n } from "@smarterp/i18n/useI18n";

export default function AllegroActions({ account, action }) {
  const { t } = useI18n();
  function disconnect() {
    if (window.confirm(t("allegro.confirm", { login: account.login }))) {
      action.run(() => deleteJSON(apiPaths.allegroConnection(account.version)));
    }
  }
  return <div className="allegro-actions" aria-busy={action.busy}>
    <button disabled={action.busy} onClick={() => action.run(() => postJSON(apiPaths.allegroCheck, {}))}>{t("allegro.check")}</button>
    <button disabled={action.busy} onClick={disconnect}>{t("allegro.disconnect")}</button>
  </div>;
}
