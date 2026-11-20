import { Settings } from "lucide-react";
import { Link } from "react-router-dom";
import { useI18n } from "@smarterp/i18n/useI18n";

export default function IntegrationTile({ provider }) {
  const { t } = useI18n();
  const { nameKey, icon: Icon, tone, path } = provider;
  const name = t(nameKey);
  const content = (
    <>
      <span className={`integration-icon integration-icon--${tone}`}><Icon aria-hidden="true" /></span>
      <span className="integration-name">{name}</span>
      <span className="integration-caption">{t(provider.caption || "integrations.planned")}</span>
    </>
  );
  if (!path) return <div className="integration-tile">{content}</div>;
  return (
    <Link to={path} className="integration-tile" aria-label={t("integrations.settings", { name })}>
      {content}<Settings className="integration-gear" size={20} aria-hidden="true" />
    </Link>
  );
}
