import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import { loginClient } from "../api/clientAuthApi";
import { setClientToken } from "../../../shared/auth/session";

const initialForm = { email: "", password: "" };

function LoginPage() {
  const navigate = useNavigate();
  const [form, setForm] = useState(initialForm);
  const [error, setError] = useState("");
  const [isLoading, setLoading] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    setLoading(true);
    try {
      const data = await loginClient(form);
      setClientToken(readToken(data));
      navigate("/dashboard");
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

  return (
    <section className="auth-shell">
      <form className="auth-card auth-form" onSubmit={handleSubmit}>
        <h2>Client Sign In</h2>
        <p className="auth-text">Enter login and password to open the client dashboard.</p>
        <label className="auth-label" htmlFor="client-email">
          Login (email)
        </label>
        <input id="client-email" name="email" className="auth-input" type="email" value={form.email} onChange={handleChange} required />
        <label className="auth-label" htmlFor="client-password">
          Password
        </label>
        <input
          id="client-password"
          name="password"
          className="auth-input"
          type="password"
          value={form.password}
          onChange={handleChange}
          required
        />
        {error && <p className="error-text">{error}</p>}
        <button className="primary-button" type="submit" disabled={isLoading}>
          {isLoading ? "Signing in..." : "Sign in"}
        </button>
        <Link to="/register" className="text-link">
          No account? Register
        </Link>
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

export default LoginPage;

