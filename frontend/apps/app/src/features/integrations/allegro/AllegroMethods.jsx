import { useI18n } from "@smarterp/i18n/useI18n";
import TokenSettings from "./token/TokenSettings";
import OAuthSettings from "./oauth/OAuthSettings";

export default function AllegroMethods({ data, action }) {
  const { t } = useI18n();
  const mode = data.connection?.method;
  if (!data.can_manage) return <p>{t("allegro.owner")}</p>;
  return <div className="allegro-methods">
    {mode !== "oauth" && <TokenSettings available={data.token_available} version={data.connection?.version} action={action} />}
    {mode !== "token" && <OAuthSettings available={data.oauth_available} connected={Boolean(mode)} action={action} />}
  </div>;
}
