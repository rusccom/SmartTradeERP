import { Link, NavLink, Outlet, useNavigate } from "react-router-dom";
import { useMemo } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import { clearAdminToken, hasAdminSession } from "../../../shared/auth/session";
import LocaleSwitcher from "../../../shared/ui/LocaleSwitcher";
import "../../../shared/ui/workspace-layout.css";

function AdminLayout() {
  const { t } = useI18n();
  const navigate = useNavigate();
  const isAuthorized = hasAdminSession();
  const links = useMemo(() => createAdminLinks(t), [t]);

  function handleLogout() {
    clearAdminToken();
    navigate("/admin", { replace: true });
  }

  return (
    <div className="workspace-zone">
      <div className="workspace-shell">
        <header className="workspace-header">
          <Link to={isAuthorized ? "/admin/dashboard" : "/admin"} className="workspace-brand">
            <span className="workspace-brand-mark workspace-brand-mark--admin" />
            <span>{t("workspace.brandAdmin")}</span>
          </Link>
          {isAuthorized && (
            <nav className="workspace-nav">
              {links.map((item) => (
                <NavLink key={item.to} to={item.to} className={readNavClass}>
                  {item.label}
                </NavLink>
              ))}
            </nav>
          )}
          <div className="flex items-center gap-3">
            <LocaleSwitcher />
            {isAuthorized && (
              <button className="workspace-logout" type="button" onClick={handleLogout}>
                {t("admin.signOut")}
              </button>
            )}
          </div>
        </header>
        <main className="workspace-content">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

function createAdminLinks(t) {
  return [
    { to: "/admin/dashboard", label: t("admin.nav.dashboard") },
    { to: "/admin/tenants", label: t("admin.nav.tenants") },
  ];
}

function readNavClass({ isActive }) {
  return isActive ? "workspace-link is-active" : "workspace-link";
}

export default AdminLayout;
