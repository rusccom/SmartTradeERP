import { Link, NavLink, Outlet, useNavigate } from "react-router-dom";

import { clearClientToken } from "../../../shared/auth/session";
import "../../../shared/ui/workspace-layout.css";

const CLIENT_LINKS = [
  { to: "/dashboard", label: "Dashboard" },
  { to: "/dashboard/products", label: "Products" },
  { to: "/dashboard/documents", label: "Documents" },
  { to: "/dashboard/reports", label: "Reports" },
];

function ClientLayout() {
  const navigate = useNavigate();

  function handleLogout() {
    clearClientToken();
    navigate("/login", { replace: true });
  }

  return (
    <div className="workspace-zone">
      <div className="workspace-shell">
        <header className="workspace-header">
          <Link to="/dashboard" className="workspace-brand">
            <span className="workspace-brand-mark" />
            <span>SmartTrade ERP</span>
          </Link>
          <nav className="workspace-nav">
            {CLIENT_LINKS.map((item) => (
              <NavLink key={item.to} to={item.to} className={readNavClass}>
                {item.label}
              </NavLink>
            ))}
          </nav>
          <button className="workspace-logout" type="button" onClick={handleLogout}>
            Sign out
          </button>
        </header>
        <main className="workspace-content">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

function readNavClass({ isActive }) {
  return isActive ? "workspace-link is-active" : "workspace-link";
}

export default ClientLayout;
