import { NavLink, Outlet } from "react-router-dom";

import { useI18n } from "../../../shared/i18n/useI18n";
import "../ui/product-catalog.css";

function ProductCatalogPage() {
  const { t } = useI18n();
  return (
    <section className="product-catalog">
      <nav className="product-catalog-tabs" aria-label={t("products.tabs.label")}>
        <NavLink end to="/dashboard/products" className={readTabClass}>
          {t("products.tabs.products")}
        </NavLink>
        <NavLink to="/dashboard/products/bundles" className={readTabClass}>
          {t("products.tabs.bundles")}
        </NavLink>
      </nav>
      <div className="product-catalog-panel">
        <Outlet />
      </div>
    </section>
  );
}

function readTabClass({ isActive }) {
  return isActive ? "product-catalog-tab is-active" : "product-catalog-tab";
}

export default ProductCatalogPage;
