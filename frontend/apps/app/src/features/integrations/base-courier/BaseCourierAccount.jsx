import { useI18n } from "@smarterp/i18n/useI18n";

const panel = "https://basecourier.com/panel";

export default function BaseCourierAccount() {
  const { t } = useI18n();
  return (
    <div className="base-courier-account">
      <p>{t("baseCourier.intro")}</p>
      <div className="base-courier-actions">
        <a className="base-courier-register" href={panel + "/rejestracja"} target="_blank" rel="noopener noreferrer">
          {t("baseCourier.register")}
        </a>
        <a href={panel + "/logowanie"} target="_blank" rel="noopener noreferrer">{t("baseCourier.signIn")}</a>
      </div>
      <details>
        <summary>{t("baseCourier.keyTitle")}</summary>
        <p>{t("baseCourier.keyHelp")}</p>
      </details>
    </div>
  );
}
