import { Link } from "react-router-dom";
import { useI18n } from "@smarterp/i18n/useI18n";
import { catalogSettingsPath } from "../settings/menu";
import VatCountryPicker from "./VatCountryPicker";
import VatRateField from "./VatRateField";
import useVatCountries from "./useVatCountries";
import "./vat.css";

export default function ProductVatRates({ rates, onChange, disabled }) {
  const { t, locale } = useI18n();
  const names = new Intl.DisplayNames([locale], { type: "region" });
  const countries = useVatCountries();
  const available = countries.enabled.filter(([code]) => !(code in rates));
  const remove = (code) => { const next = { ...rates }; delete next[code]; onChange(next); };
  return <section className="product-card" aria-label={t("vat.title")}>
    <div className="product-card__head"><h3 className="product-card__title">{t("vat.title")}</h3>
      <Link className="vat-manage" to={catalogSettingsPath} target="_blank" rel="noreferrer">{t("vat.manage")}</Link></div>
    {countries.loading && <p role="status">{t("entry.loading")}</p>}
    {countries.error && <p role="alert">{countries.error} <button type="button" onClick={countries.retry}>{t("entry.retry")}</button></p>}
    {countries.data?.length === 0 && <p className="vat-empty">{t("vat.noCountries")}</p>}
    {Object.keys(rates).sort().map((code) =>
      <VatRateField key={code} country={names.of(code)} value={rates[code]} disabled={disabled}
        onChange={(value) => onChange({ ...rates, [code]: value })} onRemove={() => remove(code)} />)}
    {available.length > 0 && <VatCountryPicker options={available} disabled={disabled}
      action={t("vat.addRate")} onAdd={(code) => onChange({ ...rates, [code]: "" })} />}
  </section>;
}
