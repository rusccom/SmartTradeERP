import { useI18n } from "../../../shared/i18n/useI18n";
import PlaceholderPage from "../../../shared/ui/PlaceholderPage";

function SettingsPage() {
  const { t } = useI18n();
  return <PlaceholderPage title={t("dashboard.settings.title")} text={t("dashboard.settings.text")} />;
}

export default SettingsPage;
