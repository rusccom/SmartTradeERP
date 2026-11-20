import { putJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";
import { useI18n } from "@smarterp/i18n/useI18n";

export default function RefreshTokenForm({ version, action }) {
  const { t } = useI18n();
  function submit(event) {
    event.preventDefault();
    const form = event.currentTarget;
    const refresh_token = form.elements.refreshToken.value.trim();
    action.run(() => putJSON(apiPaths.allegroPrivateToken, { refresh_token, version })).then(ok => { if (ok) form.reset(); });
  }
  return <details className="allegro-refresh-form"><summary>{t("allegro.useRefresh")}</summary>
    <form onSubmit={submit}><div className="allegro-field">
      <label htmlFor="allegro-refresh-token">{t("allegro.refreshLabel")}</label>
      <input id="allegro-refresh-token" name="refreshToken" type="password" autoComplete="off" autoCapitalize="none"
        spellCheck={false} required maxLength={16384} disabled={action.busy} />
    </div><button disabled={action.busy}>{t("allegro.saveRefresh")}</button></form>
  </details>;
}
