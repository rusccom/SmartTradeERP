import { useI18n } from "@smarterp/i18n/useI18n";
import PlaceholderPage from "@smarterp/ui/PlaceholderPage";

function AdminTenantsPage() {
  const { t } = useI18n();
  return <PlaceholderPage title={t("admin.tenants.title")} text={t("admin.tenants.text")} />;
}

export default AdminTenantsPage;
