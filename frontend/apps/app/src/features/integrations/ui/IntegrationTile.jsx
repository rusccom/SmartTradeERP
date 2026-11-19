import { ArrowRight } from "lucide-react";
import { Link } from "react-router-dom";
import { useI18n } from "@smarterp/i18n/useI18n";

export default function IntegrationTile({ provider }) {
  const { t } = useI18n();
  const { name, icon: Icon, tone, path } = provider;
  const content = (
    <>
      <span className={`integration-icon integration-icon--${tone}`}><Icon aria-hidden="true" /></span>
      <span className="integration-name">{name}</span>
      <span className="integration-caption">{t(path ? "integrations.inpost.manager" : "integrations.planned")}</span>
    </>
  );
  if (!path) return <div className="integration-tile">{content}</div>;
  return (
    <Link to={path} className="integration-tile">
      {content}<ArrowRight className="integration-arrow" size={18} aria-hidden="true" />
    </Link>
  );
}
