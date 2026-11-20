import { useState } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";
import { putJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";
import OrderField from "../../order/create-order/OrderField";

export default function SellerProfileForm({ initial, canManage }) {
  const { t } = useI18n();
  const [profile, setProfile] = useState(() => editable(initial));
  const [status, setStatus] = useState({});
  const field = (key, label, props = {}) => <OrderField key={key} label={t(`entry.${label}`)} value={profile[key]}
    onChange={(v) => { setProfile((p) => ({ ...p, [key]: v })); setStatus({}); }} {...props}/>;
  const save = (e) => { e.preventDefault(); saveProfile(profile, setProfile, setStatus); };
  return <form className="order-seller" onSubmit={save}><fieldset className="order-card" disabled={!canManage || status.busy}>
    <h2>{t("entry.sellerDetails")}</h2>{field("name", "sellerName", { required: true })}{field("nip", "nip", { required: true, pattern: "[0-9]{10}", maxLength: 10 })}
    {field("street", "street", { required: true })}{field("postal_code", "postal_code", { required: true })}{field("city", "city", { required: true })}
    <h2>{t("entry.ksefAccess")}</h2>{field("environment", "environment", { required: true, options: [["", "—"], ["test", t("entry.test")], ["production", t("entry.production")]] })}
    {field("token", "token", { type: "password", autoComplete: "new-password", disabled: profile.clear_token })}
    {profile.token_configured && <label className="order-check"><input type="checkbox" checked={profile.clear_token}
      onChange={(e) => setProfile({ ...profile, clear_token: e.target.checked, token: "" })}/>{t("entry.clearToken")}</label>}
    {canManage && <button type="submit" className="order-btn order-primary">{t("entry.saveSettings")}</button>}
  </fieldset>{!canManage && <p>{t("entry.ownerOnly")}</p>}{status.message && <p role="status" className={status.error ? "order-error" : "order-notice"}>{status.message}</p>}
    {status.saved && <p role="status">{t("entry.settingsSaved")}</p>}</form>;
}

function editable(s) {
  return { name: "", nip: "", street: "", postal_code: "", city: "", environment: "", revision: 0, ...s, token: "", clear_token: false };
}

async function saveProfile(p, setProfile, setStatus) {
  setStatus({ busy: true });
  try {
    const saved = await putJSON(apiPaths.polishSeller, p);
    setProfile(editable(saved));
    setStatus({ saved: true });
  } catch (e) { setStatus({ message: e.message, error: true }); }
}
