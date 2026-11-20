import { useI18n } from "@smarterp/i18n/useI18n";

export default function VatRateField({ country, value, onChange, onRemove, disabled }) {
  const { t } = useI18n();
  return <div className="vat-rate-row">
    <label><span>{country} · VAT (%)</span>
      <input type="number" value={value} min="0" max="100" step="0.01" required disabled={disabled}
        aria-label={t("vat.rateLabel", { country })} onChange={(event) => onChange(event.target.value)} />
    </label>
    <button type="button" className="vat-action" disabled={disabled} aria-label={t("vat.removeRate", { country })}
      onClick={onRemove}>{t("entry.remove")}</button>
  </div>;
}
