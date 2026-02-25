import { Link } from "react-router-dom";

function LandingPage() {
  return (
    <section className="auth-shell">
      <div className="auth-card">
        <h2>SmartERP</h2>
        <p className="auth-text">Basic public page. Choose client or admin sign in.</p>
        <Link to="/login" className="primary-button">
          Sign in or register
        </Link>
        <Link to="/admin" className="text-link">
          Admin sign in
        </Link>
      </div>
    </section>
  );
}

export default LandingPage;

