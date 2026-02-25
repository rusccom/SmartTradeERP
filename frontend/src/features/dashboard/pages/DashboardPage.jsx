import { Link } from "react-router-dom";

function DashboardPage() {
  return (
    <section className="placeholder dashboard-stub">
      <h2>Client Dashboard</h2>
      <p>Client dashboard is isolated from public and admin zones after sign in.</p>
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
    </section>
  );
}

export default DashboardPage;

