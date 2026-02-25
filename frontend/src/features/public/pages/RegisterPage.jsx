import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import { registerClient } from "../api/clientAuthApi";
import { setClientToken } from "../../../shared/auth/session";

const initialForm = { tenant_name: "", email: "", password: "" };

function RegisterPage() {
  const navigate = useNavigate();
  const [form, setForm] = useState(initialForm);
  const [error, setError] = useState("");
  const [isLoading, setLoading] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    setLoading(true);
    try {
      const data = await registerClient(form);
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
        <h2>Client Registration</h2>
        <p className="auth-text">Create a company account and the client dashboard will open.</p>
        <label className="auth-label" htmlFor="tenant-name">
          Company name
        </label>
        <input
          id="tenant-name"
          name="tenant_name"
          className="auth-input"
          type="text"
          value={form.tenant_name}
          onChange={handleChange}
          required
        />
        <label className="auth-label" htmlFor="register-email">
          Login (email)
        </label>
        <input
          id="register-email"
          name="email"
          className="auth-input"
          type="email"
          value={form.email}
          onChange={handleChange}
          required
        />
        <label className="auth-label" htmlFor="register-password">
          Password
        </label>
        <input
          id="register-password"
          name="password"
          className="auth-input"
          type="password"
          value={form.password}
          onChange={handleChange}
          required
        />
        {error && <p className="error-text">{error}</p>}
        <button className="primary-button" type="submit" disabled={isLoading}>
          {isLoading ? "Creating account..." : "Register"}
        </button>
        <Link to="/login" className="text-link">
          Already registered? Sign in
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

export default RegisterPage;

