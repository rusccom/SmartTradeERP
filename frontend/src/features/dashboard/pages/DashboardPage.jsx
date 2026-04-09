import { useI18n } from "../../../shared/i18n/useI18n";

function DashboardPage() {
  const { t } = useI18n();
  return (
    <section className="placeholder dashboard-stub">
      <h2>{t("dashboard.home.title")}</h2>
      <p>{t("dashboard.home.description")}</p>
    </section>
  );
}

export default DashboardPage;
