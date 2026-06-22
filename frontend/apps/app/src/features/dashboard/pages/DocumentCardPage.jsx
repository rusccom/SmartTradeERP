import { useI18n } from "@smarterp/i18n/useI18n";
import PlaceholderPage from "@smarterp/ui/PlaceholderPage";

function DocumentCardPage() {
  const { t } = useI18n();
  return <PlaceholderPage title={t("dashboard.documentCard.title")} text={t("dashboard.documentCard.text")} />;
}

export default DocumentCardPage;
