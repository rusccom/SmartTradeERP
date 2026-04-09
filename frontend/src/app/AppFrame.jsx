import { Link, Outlet } from "react-router-dom";
import { useMemo } from "react";

import { hasAdminSession, hasClientSession } from "../shared/auth/session";
import { buildLandingPath } from "../shared/i18n/localeConfig";
import { useI18n } from "../shared/i18n/useI18n";

function AppFrame() {
  const { locale, t } = useI18n();
  const navLinks = useMemo(() => appendPrivateLinks(createBaseLinks(t, locale), t), [locale, t]);
  return (
    <div className="page">
      <header className="topbar">
        <h1>{t("brand.name")}</h1>
        <nav>
          {navLinks.map((item) => (
            <Link key={item.to} to={item.to} className="nav-link">
              {item.label}
            </Link>
          ))}
        </nav>
      </header>
      <main className="content">
        <Outlet />
      </main>
    </div>
  );
}

function createBaseLinks(t, locale) {
  return [
    { to: buildLandingPath(locale), label: t("app.nav.landing") },
    { to: "/login", label: t("app.nav.clientLogin") },
    { to: "/register", label: t("app.nav.clientRegister") },
    { to: "/admin", label: t("app.nav.adminLogin") },
  ];
}

function appendPrivateLinks(baseLinks, t) {
  const items = [...baseLinks];
  if (hasClientSession()) {
    items.push({ to: "/dashboard", label: t("app.nav.clientDashboard") });
  }
  if (hasAdminSession()) {
    items.push({ to: "/admin/dashboard", label: t("app.nav.adminDashboard") });
  }
  return items;
}

export default AppFrame;
