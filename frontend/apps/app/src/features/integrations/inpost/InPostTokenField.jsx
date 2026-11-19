import { useI18n } from "@smarterp/i18n/useI18n";

export default function InPostTokenField({ reuse }) {
  const { t } = useI18n();
  return (
    <div>
      <label>{t("inpost.token")}
        <input name="token" type="password" autoComplete="new-password" spellCheck={false} autoCapitalize="none"
          maxLength={4096} required={!reuse} aria-describedby="inpost-token-hint" />
      </label>
      <p id="inpost-token-hint">{t(reuse ? "inpost.reuseToken" : "inpost.tokenHint")}</p>
    </div>
  );
}
