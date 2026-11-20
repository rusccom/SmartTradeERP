import { NavLink, Outlet } from "react-router-dom";
import { useI18n } from "@smarterp/i18n/useI18n";
import "../vat/vat.css";

export default function CatalogSettingsPage() {
  const { t } = useI18n();
  return <section className="catalog-settings">
    <h1>{t("dashboard.menu.settings")}</h1>
    <div className="catalog-settings-grid">
      <nav aria-label={t("dashboard.menu.settings")}>
        <NavLink to="products">{t("dashboard.menu.products")}</NavLink>
      </nav>
      <Outlet />
    </div>
  </section>;
}
