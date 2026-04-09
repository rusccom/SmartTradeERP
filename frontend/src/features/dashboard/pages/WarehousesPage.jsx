import { useI18n } from "../../../shared/i18n/useI18n";
import PlaceholderPage from "../../../shared/ui/PlaceholderPage";

function WarehousesPage() {
  const { t } = useI18n();
  return <PlaceholderPage title={t("dashboard.warehouses.title")} text={t("dashboard.warehouses.text")} />;
}

export default WarehousesPage;
