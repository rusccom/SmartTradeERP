import { useI18n } from "@smarterp/i18n/useI18n";
import PlaceholderPage from "@smarterp/ui/PlaceholderPage";

function ReportsPage() {
  const { t } = useI18n();
  return <PlaceholderPage title={t("dashboard.reports.title")} text={t("dashboard.reports.text")} />;
}

export default ReportsPage;
