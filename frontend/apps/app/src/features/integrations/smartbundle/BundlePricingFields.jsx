import { useI18n } from "@smarterp/i18n/useI18n";

export default function BundlePricingFields({ settings, patch }) {
  const { t } = useI18n();
  return (
    <>
      <label>{t("smartBundle.pricing")}
        <select value={settings.mode} onChange={event => patch({ mode: event.target.value })}>
          <option value="fixed">{t("smartBundle.fixed")}</option>
          <option value="dynamic">{t("smartBundle.dynamic")}</option>
        </select>
      </label>
      <label>{t("smartBundle.bundleDiscount")}
        <input type="number" min="0" max="100" step="0.01" required
          value={settings.discount} onChange={event => patch({ discount: event.target.value })} />
      </label>
    </>
  );
}
