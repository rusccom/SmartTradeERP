import { useI18n } from "@smarterp/i18n/useI18n";

export default function BaseCourierFields({ email }) {
  const { t } = useI18n();
  return (
    <>
      <label htmlFor="base-courier-email">{t("baseCourier.email")}</label>
      <input id="base-courier-email" name="email" type="email" defaultValue={email || ""} maxLength={254}
        autoCapitalize="none" spellCheck={false} required />
      <label htmlFor="base-courier-key">{t("baseCourier.key")}</label>
      <input id="base-courier-key" name="api_key" type="password" maxLength={4096} autoComplete="off" required />
      <details>
        <summary>{t("baseCourier.keyTitle")}</summary>
        <p>{t("baseCourier.keyHelp")}</p>
      </details>
    </>
  );
}
