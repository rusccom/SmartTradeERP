import { putJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";
import { useI18n } from "@smarterp/i18n/useI18n";
import TokenField from "./TokenField";

function submitToken(event, action, version) {
  event.preventDefault();
  const form = event.currentTarget;
  const access_token = form.elements.accessToken.value.trim();
  action.run(() => putJSON(apiPaths.allegroToken, { access_token, version })).then(ok => { if (ok) form.reset(); });
}

export default function TokenSettings({ available, version = "", action }) {
  const { t } = useI18n();
  return <section className="allegro-card" aria-labelledby="allegro-token-title">
    <h2 id="allegro-token-title">{t(version ? "allegro.replaceToken" : "allegro.byToken")}</h2>
    <form onSubmit={event => submitToken(event, action, version)}>
      <TokenField disabled={action.busy || !available} />
      {!available && <p>{t("allegro.tokenUnavailable")}</p>}
      <button className="allegro-primary" disabled={action.busy || !available}>{t("allegro.saveToken")}</button>
    </form>
  </section>;
}
