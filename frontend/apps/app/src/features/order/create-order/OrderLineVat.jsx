import { useI18n } from "@smarterp/i18n/useI18n";

export default function OrderLineVat({ line, quote, pending, disabled, onChange }) {
  const { t } = useI18n();
  const manual = line.vat_mode === "manual";
  const data = manual ? line : quote ?? line;
  const value = data.vat_rate ?? "";
  const status = value === "" ? t("vat.notSet") : manual && t("vat.manual");
  const origin = [data.vat_country, status].filter(Boolean).join(" · ");
  return <div className="order-line-vat">
    <input className="order-number" type="number" min="0" max="100" step="0.01" disabled={disabled}
      value={value} aria-label={`${t("entry.vat")} · ${line.product_name}`} onChange={(event) => onChange(event.target.value)} />
    <small>{pending && !manual ? t("entry.loading") : origin}</small>
    {manual && !disabled && <button type="button" onClick={() => onChange("")}>{t("vat.restore")}</button>}
  </div>;
}
