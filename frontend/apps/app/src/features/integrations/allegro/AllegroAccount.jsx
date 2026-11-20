import { useI18n } from "@smarterp/i18n/useI18n";
import AllegroActions from "./AllegroActions";

export default function AllegroAccount({ account, canManage, action }) {
  const { t, locale } = useI18n();
  const { login, method, verified_at, expires_at, needs_reconnect } = account;
  const expired = needs_reconnect || (method === "token" && expires_at && Date.parse(expires_at) <= Date.now());
  const reauthorize = method === "token" ? "allegro.tokenExpired" : "allegro.needsReconnect";
  return <section className="allegro-card allegro-current">
    <h2>{login}</h2>
    <p className={expired ? "allegro-error" : "allegro-verified"}>{t(expired ? reauthorize : "allegro.verified")}</p>
    <dl><dt>{t("allegro.method")}</dt><dd>{t(method === "token" ? "allegro.tokenLabel" : "allegro.byLogin")}</dd>
      <dt>{t("allegro.verifiedAt")}</dt><dd><time dateTime={verified_at}>{new Date(verified_at).toLocaleString(locale)}</time></dd>
      {method === "token" && expires_at && <><dt>{t("allegro.expires")}</dt>
        <dd><time dateTime={expires_at}>{new Date(expires_at).toLocaleString(locale)}</time></dd></>}
    </dl>
    {canManage && <AllegroActions account={account} action={action} />}
  </section>;
}
