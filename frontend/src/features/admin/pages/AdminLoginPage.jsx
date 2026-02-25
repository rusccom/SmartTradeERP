import { useState } from "react";
import { useNavigate } from "react-router-dom";

import { loginAdmin } from "../api/adminAuthApi";
import { setAdminToken } from "../../../shared/auth/session";

const DEFAULT_ADMIN = {
  email: "owner@smarterp.local",
  password: "Owner#2026",
};

function AdminLoginPage() {
  const navigate = useNavigate();
  const [form, setForm] = useState(DEFAULT_ADMIN);
  const [error, setError] = useState("");
  const [isLoading, setLoading] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    setLoading(true);
    try {
      const data = await loginAdmin(form);
      setAdminToken(readToken(data));
      navigate("/admin/dashboard");
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setLoading(false);
    }
  }

  function handleChange(event) {
    const { name, value } = event.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  }

  function applyDefaultCredentials() {
    setForm(DEFAULT_ADMIN);
  }

  return (
    <section className="auth-shell">
      <form className="auth-card auth-form" onSubmit={handleSubmit}>
        <h2>Admin Sign In</h2>
        <p className="auth-text">Registration is disabled. Sign in with standard admin credentials only.</p>
        <label className="auth-label" htmlFor="admin-email">
          Login (email)
        </label>
        <input id="admin-email" name="email" className="auth-input" type="email" value={form.email} onChange={handleChange} required />
        <label className="auth-label" htmlFor="admin-password">
          Password
        </label>
        <input
          id="admin-password"
          name="password"
          className="auth-input"
          type="password"
          value={form.password}
          onChange={handleChange}
          required
        />
        <button className="secondary-button" type="button" onClick={applyDefaultCredentials}>
          Use standard credentials
        </button>
        {error && <p className="error-text">{error}</p>}
        <button className="primary-button" type="submit" disabled={isLoading}>
          {isLoading ? "Signing in..." : "Sign in to admin panel"}
        </button>
      </form>
    </section>
  );
}

function readToken(data) {
  const token = data?.access_token;
  if (token) {
    return token;
  }
  throw new Error("Server did not return an access token");
}

export default AdminLoginPage;

