import { useState } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";

export default function VatCountryPicker({ options, onAdd, action, disabled }) {
  const { t } = useI18n();
  const [selected, setSelected] = useState("");
  const code = options.some(([key]) => key === selected) ? selected : "";
  return <div className="vat-picker">
    <select aria-label={t("entry.country")} value={code} disabled={disabled || !options.length}
      onChange={(event) => setSelected(event.target.value)}>
      <option value="">{t("vat.select")}</option>
      {options.map(([key, name]) => <option value={key} key={key}>{name}</option>)}
    </select>
    <button type="button" className="vat-action" disabled={disabled || !code} onClick={() => onAdd(code)}>{action}</button>
  </div>;
}
