import { Link, Outlet } from "react-router-dom";

const links = [
  { to: "/", label: "Landing" },
  { to: "/register", label: "Register" },
  { to: "/login", label: "Login" },
  { to: "/admin", label: "Admin" },
  { to: "/admin/dashboard", label: "Admin Dashboard" },
  { to: "/admin/tenants", label: "Admin Tenants" },
  { to: "/dashboard", label: "Client Dashboard" },
  { to: "/dashboard/products", label: "Products" },
  { to: "/dashboard/bundles", label: "Bundles" },
  { to: "/dashboard/warehouses", label: "Warehouses" },
  { to: "/dashboard/documents", label: "Documents" },
  { to: "/dashboard/reports", label: "Reports" },
  { to: "/dashboard/settings", label: "Settings" },
];

function AppFrame() {
  return (
    <div className="page">
      <header className="topbar">
        <h1>SmartERP</h1>
        <nav>
          {links.map((item) => (
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

export default AppFrame;
