import { useI18n } from "@smarterp/i18n/useI18n";
import TokenSettings from "./token/TokenSettings";
import OAuthSettings from "./oauth/OAuthSettings";
import PrivateSettings from "./personal/PrivateSettings";

export default function AllegroMethods({ data, action, reload }) {
  const { t } = useI18n();
  const mode = data.connection?.method;
  const automatic = mode === "token" && data.connection.auto_renew;
  const privateSetup = data.private_renewal_available && mode !== "oauth" && (!automatic || data.connection.needs_reconnect);
  if (!data.can_manage) return <p>{t("allegro.owner")}</p>;
  return <div className="allegro-methods">
    {privateSetup && <PrivateSettings version={data.connection?.version} action={action} reload={reload} />}
    {mode !== "oauth" && !automatic && <TokenSettings available={data.token_available} version={data.connection?.version} action={action} />}
    {mode !== "token" && <OAuthSettings available={data.oauth_available} connected={Boolean(mode)} action={action} />}
  </div>;
}
