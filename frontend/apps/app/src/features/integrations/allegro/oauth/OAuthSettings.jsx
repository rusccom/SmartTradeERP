import { useI18n } from "@smarterp/i18n/useI18n";
import { connectSeller } from "./api";

export default function OAuthSettings({ available, connected, action }) {
  const { t } = useI18n();
  return <section className="allegro-card" aria-labelledby="allegro-login-title">
    <h2 id="allegro-login-title">{t("allegro.byLogin")}</h2>
    <p id="allegro-login-hint">{t(available ? "allegro.loginHint" : "allegro.loginUnavailable")}</p>
    <button aria-describedby="allegro-login-hint" disabled={!available || action.busy} onClick={() => action.run(connectSeller)}>
      {t(connected ? "allegro.reconnect" : "allegro.signIn")}
    </button>
  </section>;
}
