import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import { setClientToken } from "../../../shared/auth/session";
import { useI18n } from "../../../shared/i18n/useI18n";
import { registerClient } from "../api/clientAuthApi";

const initialForm = { tenant_name: "", email: "", password: "" };

function RegisterPage() {
  const { t } = useI18n();
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
      setClientToken(readToken(data, t));
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
        <h2>{t("client.auth.register.title")}</h2>
        <p className="auth-text">{t("client.auth.register.description")}</p>
        <label className="auth-label" htmlFor="tenant-name">
          {t("auth.companyName")}
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
          {t("auth.email")}
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
          {t("auth.password")}
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
          {isLoading ? t("client.auth.register.loading") : t("client.auth.register.button")}
        </button>
        <Link to="/login" className="text-link">
          {t("client.auth.register.link")}
        </Link>
      </form>
    </section>
  );
}

function readToken(data, t) {
  const token = data?.access_token;
  if (token) {
    return token;
  }
  throw new Error(t("auth.error.missingToken"));
}

export default RegisterPage;
