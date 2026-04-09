import { Link } from "react-router-dom";

import { useI18n } from "../../../shared/i18n/useI18n";

function AdminDashboardPage() {
  const { t } = useI18n();
  return (
    <section className="placeholder dashboard-stub">
      <h2>{t("admin.dashboard.title")}</h2>
      <p>{t("admin.dashboard.description")}</p>
      <div className="dashboard-links">
        <Link className="nav-link" to="/admin/tenants">
          {t("admin.dashboard.link")}
        </Link>
      </div>
    </section>
  );
}

export default AdminDashboardPage;
