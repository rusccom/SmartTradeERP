import { Link, Outlet } from "react-router-dom";

import { hasAdminSession, hasClientSession } from "../shared/auth/session";

const links = [
  { to: "/", label: "Landing" },
  { to: "/login", label: "Client Login" },
  { to: "/register", label: "Client Register" },
  { to: "/admin", label: "Admin Login" },
];

function AppFrame() {
  const navLinks = appendPrivateLinks(links);
  return (
    <div className="page">
      <header className="topbar">
        <h1>SmartERP</h1>
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

function appendPrivateLinks(baseLinks) {
  const items = [...baseLinks];
  if (hasClientSession()) {
    items.push({ to: "/dashboard", label: "Client Dashboard" });
  }
  if (hasAdminSession()) {
    items.push({ to: "/admin/dashboard", label: "Admin Dashboard" });
  }
  return items;
}

export default AppFrame;
