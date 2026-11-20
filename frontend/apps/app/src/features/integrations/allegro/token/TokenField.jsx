import { useI18n } from "@smarterp/i18n/useI18n";

export default function TokenField({ disabled }) {
  const { t } = useI18n();
  return <div className="allegro-field">
    <label htmlFor="allegro-access-token">{t("allegro.tokenLabel")}</label>
    <input type="password" id="allegro-access-token" name="accessToken" autoComplete="off"
      autoCapitalize="none" spellCheck={false} required maxLength={16384} disabled={disabled}
      aria-describedby="allegro-token-hint" />
    <p id="allegro-token-hint">{t("allegro.tokenHint")}</p>
  </div>;
}
