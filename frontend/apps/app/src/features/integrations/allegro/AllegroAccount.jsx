import { useI18n } from "@smarterp/i18n/useI18n";
import AllegroActions from "./AllegroActions";
import AccountDetails from "./AccountDetails";

function accessMessage(account) {
  if (account.needs_reconnect) return account.auto_renew ? "allegro.needsReconnect" : "allegro.tokenExpired";
  if (!account.auto_renew && account.expires_at && Date.parse(account.expires_at) <= Date.now()) return "allegro.tokenExpired";
  return "allegro.verified";
}

export default function AllegroAccount({ account, canManage, action }) {
  const { t } = useI18n();
  const message = accessMessage(account);
  return <section className="allegro-card allegro-current">
    <h2>{account.login}</h2>
    <p className={message === "allegro.verified" ? "allegro-verified" : "allegro-error"}>{t(message)}</p>
    <AccountDetails account={account} />
    {canManage && <AllegroActions account={account} action={action} />}
  </section>;
}
