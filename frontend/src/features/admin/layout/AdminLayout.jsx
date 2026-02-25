import { Link, NavLink, Outlet, useNavigate } from "react-router-dom";

import { clearAdminToken, hasAdminSession } from "../../../shared/auth/session";
import "../../../shared/ui/workspace-layout.css";

const ADMIN_LINKS = [
  { to: "/admin/dashboard", label: "Dashboard" },
  { to: "/admin/tenants", label: "Tenants" },
];

function AdminLayout() {
  const navigate = useNavigate();
  const isAuthorized = hasAdminSession();

  function handleLogout() {
    clearAdminToken();
    navigate("/admin", { replace: true });
  }

  return (
    <div className="workspace-shell">
      <header className="workspace-header">
        <Link to={isAuthorized ? "/admin/dashboard" : "/admin"} className="workspace-brand">
          <span className="workspace-brand-mark workspace-brand-mark--admin" />
          <span>SmartTrade ERP Admin</span>
        </Link>
        {isAuthorized && (
          <nav className="workspace-nav">
            {ADMIN_LINKS.map((item) => (
              <NavLink key={item.to} to={item.to} className={readNavClass}>
                {item.label}
              </NavLink>
            ))}
          </nav>
        )}
        {isAuthorized && (
          <button className="workspace-logout" type="button" onClick={handleLogout}>
            Sign out
          </button>
        )}
      </header>
      <main className="workspace-content">
        <Outlet />
      </main>
    </div>
  );
}

function readNavClass({ isActive }) {
  return isActive ? "workspace-link is-active" : "workspace-link";
}

export default AdminLayout;
