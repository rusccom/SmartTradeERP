import { ExternalLink } from "lucide-react";
import { useI18n } from "@smarterp/i18n/useI18n";

export default function InPostToolLink({ tool }) {
  const { t } = useI18n();
  const { key, icon: Icon, url } = tool;
  return (
    <a className="inpost-tool" href={url} target="_blank" rel="noopener noreferrer">
      <Icon className="inpost-tool-icon" size={22} aria-hidden="true" />
      <span>
        <strong>{t(`integrations.inpost.${key}`)}</strong>
        <span className="integration-caption">{t(`integrations.inpost.${key}Hint`)}</span>
      </span>
      <ExternalLink className="inpost-external" size={16} aria-hidden="true" />
    </a>
  );
}
