import { useI18n } from "../../../shared/i18n/useI18n";
import PlaceholderPage from "../../../shared/ui/PlaceholderPage";

function DocumentCardPage() {
  const { t } = useI18n();
  return <PlaceholderPage title={t("dashboard.documentCard.title")} text={t("dashboard.documentCard.text")} />;
}

export default DocumentCardPage;
