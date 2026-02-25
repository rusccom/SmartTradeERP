import { Link } from "react-router-dom";

function AdminDashboardPage() {
  return (
    <section className="placeholder dashboard-stub">
      <h2>Admin Dashboard</h2>
      <p>Admin zone is isolated from client and public areas and opens only under `/admin`.</p>
      <div className="dashboard-links">
        <Link className="nav-link" to="/admin/tenants">
          Tenant list
        </Link>
      </div>
    </section>
  );
}

export default AdminDashboardPage;

