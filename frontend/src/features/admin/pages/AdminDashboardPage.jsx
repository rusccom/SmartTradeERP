import { Link, useNavigate } from "react-router-dom";

import { clearAdminToken } from "../../../shared/auth/session";

function AdminDashboardPage() {
  const navigate = useNavigate();

  function handleLogout() {
    clearAdminToken();
    navigate("/admin");
  }

  return (
    <section className="placeholder dashboard-stub">
      <h2>Admin Dashboard</h2>
      <p>Admin dashboard placeholder. Login is separate from the client area and has no registration.</p>
      <div className="dashboard-links">
        <Link className="nav-link" to="/admin/tenants">
          Tenant List
        </Link>
      </div>
      <button className="secondary-button" type="button" onClick={handleLogout}>
        Sign out from admin
      </button>
    </section>
  );
}

export default AdminDashboardPage;

