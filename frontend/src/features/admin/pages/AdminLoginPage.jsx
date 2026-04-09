import { useState } from "react";
import { useNavigate } from "react-router-dom";

import { setAdminToken } from "../../../shared/auth/session";
import { useI18n } from "../../../shared/i18n/useI18n";
import { loginAdmin } from "../api/adminAuthApi";

const DEFAULT_ADMIN = {
  email: "owner@smarterp.local",
  password: "Owner#2026",
};

function AdminLoginPage() {
  const { t } = useI18n();
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
      setAdminToken(readToken(data, t));
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
        <h2>{t("admin.auth.title")}</h2>
        <p className="auth-text">{t("admin.auth.description")}</p>
        <label className="auth-label" htmlFor="admin-email">
          {t("auth.email")}
        </label>
        <input id="admin-email" name="email" className="auth-input" type="email" value={form.email} onChange={handleChange} required />
        <label className="auth-label" htmlFor="admin-password">
          {t("auth.password")}
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
          {t("admin.auth.defaultCredentials")}
        </button>
        {error && <p className="error-text">{error}</p>}
        <button className="primary-button" type="submit" disabled={isLoading}>
          {isLoading ? t("client.auth.login.loading") : t("admin.auth.button")}
        </button>
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

export default AdminLoginPage;
