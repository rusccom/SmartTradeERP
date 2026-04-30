import { Link, Outlet, useNavigate } from "react-router-dom";
import { useCallback, useEffect, useMemo, useState } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import { clearClientToken } from "../../../shared/auth/session";
import CurrencyProvider from "../../currencies/model/CurrencyProvider";
import LocaleSwitcher from "../../../shared/ui/LocaleSwitcher";
import { createMenuSections } from "../registry";
import Sidebar from "./Sidebar";
import "../../../shared/ui/workspace-layout.css";
import "./sidebar.css";

function ClientLayout() {
  const { t } = useI18n();
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);
  const sections = useMemo(() => createMenuSections(t), [t]);
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
            aria-label={t("client.layout.openMenu")}
          >
            <span className="hamburger-icon" />
          </button>
          <Link to="/dashboard" className="workspace-brand">
            <span className="workspace-brand-mark" />
            <span>{t("client.layout.brand")}</span>
          </Link>
          <div className="flex items-center gap-3">
            <LocaleSwitcher />
            <button className="workspace-logout" type="button" onClick={handleLogout}>
              {t("client.signOut")}
            </button>
          </div>
        </header>
        <CurrencyProvider>
          <main className="workspace-content">
            <Outlet />
          </main>
        </CurrencyProvider>
      </div>

      <Sidebar open={menuOpen} sections={sections} onClose={closeMenu} />
    </div>
  );
}

export default ClientLayout;
