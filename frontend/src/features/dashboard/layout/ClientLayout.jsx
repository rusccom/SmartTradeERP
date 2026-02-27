import { Link, Outlet, useNavigate } from "react-router-dom";
import { useCallback, useEffect, useState } from "react";

import { clearClientToken } from "../../../shared/auth/session";
import { MENU_SECTIONS } from "../registry";
import Sidebar from "./Sidebar";
import "../../../shared/ui/workspace-layout.css";
import "./sidebar.css";

function ClientLayout() {
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  const toggleMenu = useCallback(() => setMenuOpen((v) => !v), []);
  const closeMenu = useCallback(() => setMenuOpen(false), []);

  useEffect(() => {
    if (!menuOpen) return;
    const onKey = (e) => e.key === "Escape" && closeMenu();
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [menuOpen, closeMenu]);

  function handleLogout() {
    clearClientToken();
    navigate("/login", { replace: true });
  }

  return (
    <div className="workspace-zone">
      <div className="workspace-shell">
        <header className="workspace-header">
          <button
            className="hamburger-btn"
            type="button"
            onClick={toggleMenu}
            aria-label="Toggle menu"
          >
            <span className="hamburger-icon" />
          </button>
          <Link to="/dashboard" className="workspace-brand">
            <span className="workspace-brand-mark" />
            <span>SmartTrade ERP</span>
          </Link>
          <button
            className="workspace-logout"
            type="button"
            onClick={handleLogout}
          >
            Sign out
          </button>
        </header>
        <main className="workspace-content">
          <Outlet />
        </main>
      </div>

      <Sidebar
        open={menuOpen}
        sections={MENU_SECTIONS}
        onClose={closeMenu}
        onLogout={handleLogout}
      />
    </div>
  );
}

export default ClientLayout;
