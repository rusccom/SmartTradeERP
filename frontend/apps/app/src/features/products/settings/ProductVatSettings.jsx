import { useState } from "react";
import { apiPaths } from "@smarterp/api/publicApi";
import { deleteJSON, putJSON } from "@smarterp/api/http";
import { useI18n } from "@smarterp/i18n/useI18n";
import VatCountryPicker from "../vat/VatCountryPicker";
import useVatCountries from "../vat/useVatCountries";

export default function ProductVatSettings() {
  const { t } = useI18n();
  const countries = useVatCountries();
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const available = countries.options.filter(([code]) => !countries.data?.includes(code));
  const change = (code, remove) => saveCountry({ code, remove, setError, setBusy, retry: countries.retry });
  return <section className="vat-settings-panel" aria-label={t("vat.title")}>
    <h2>{t("vat.title")}</h2><p className="vat-empty">{t("vat.settingsHint")}</p>
    <VatCountryPicker options={available} onAdd={(code) => change(code, false)}
      disabled={busy || !countries.data} action={t("vat.addCountry")} />
    {(error || countries.error) && <p role="alert">{error || countries.error} <button type="button" onClick={countries.retry}>{t("entry.retry")}</button></p>}
    {countries.loading && <p role="status">{t("entry.loading")}</p>}
    {countries.data?.length === 0 && <p className="vat-empty">{t("vat.noCountries")}</p>}
    <ul>{countries.enabled.map(([code, name]) => <li key={code}><span>{name} <small>{code}</small></span>
      <button type="button" className="vat-action" disabled={busy} aria-label={t("vat.removeCountry", { country: name })}
        onClick={() => change(code, true)}>{t("entry.remove")}</button></li>)}</ul>
  </section>;
}

async function saveCountry({ code, remove, setError, setBusy, retry }) {
  setBusy(true);
  setError("");
  try {
    const path = apiPaths.vatCountry(code);
    if (remove) await deleteJSON(path);
    else await putJSON(path, {});
    retry();
  } catch (error) {
    setError(error.message);
  } finally {
    setBusy(false);
  }
}
