import { useI18n } from "../../../shared/i18n/useI18n";
import PlaceholderPage from "../../../shared/ui/PlaceholderPage";

function BundlesPage() {
  const { t } = useI18n();
  return <PlaceholderPage title={t("dashboard.bundles.title")} text={t("dashboard.bundles.text")} />;
}

export default BundlesPage;
