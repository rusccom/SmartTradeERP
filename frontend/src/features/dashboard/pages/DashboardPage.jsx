import { Link, useNavigate } from "react-router-dom";

import { clearClientToken } from "../../../shared/auth/session";

function DashboardPage() {
  const navigate = useNavigate();

  function handleLogout() {
    clearClientToken();
    navigate("/login");
  }

  return (
    <section className="placeholder dashboard-stub">
      <h2>Client Dashboard</h2>
      <p>Client dashboard placeholder. Sign in redirect and session persistence are active.</p>
      <div className="dashboard-links">
        <Link className="nav-link" to="/dashboard/products">
          Products
        </Link>
        <Link className="nav-link" to="/dashboard/documents">
          Documents
        </Link>
        <Link className="nav-link" to="/dashboard/reports">
          Reports
        </Link>
      </div>
      <button className="secondary-button" type="button" onClick={handleLogout}>
        Sign out from client
      </button>
    </section>
  );
}

export default DashboardPage;

