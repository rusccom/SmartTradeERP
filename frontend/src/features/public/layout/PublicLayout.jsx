import { Link, Outlet } from "react-router-dom";

import "../styles/public-layout.css";

function PublicLayout() {
  return (
    <div className="public-shell">
      <header className="public-header">
        <Link to="/" className="public-brand" aria-label="SmartTrade ERP home">
          <span className="public-brand-mark" />
          <span className="public-brand-text">SmartTrade ERP</span>
        </Link>
        <nav className="public-actions">
          <Link to="/login" className="public-action-link">
            Sign in
          </Link>
          <Link to="/register" className="public-action-link public-action-link--primary">
            Register
          </Link>
        </nav>
      </header>
      <main className="public-content">
        <Outlet />
      </main>
    </div>
  );
}

export default PublicLayout;
