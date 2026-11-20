import { useI18n } from "@smarterp/i18n/useI18n";
import useDeviceAuthorization from "./useDeviceAuthorization";
import RefreshTokenForm from "./RefreshTokenForm";

export default function PrivateSettings({ version = "", action, reload }) {
  const { t } = useI18n();
  const { ticket, busy, error, start } = useDeviceAuthorization(version, reload);
  return <section className="allegro-card" aria-labelledby="allegro-renew-title">
    <h2 id="allegro-renew-title">{t("allegro.privateApp")}</h2>
    <p>{t(ticket ? "allegro.waitSeller" : "allegro.privateHint")}</p>
    {error && <p role="alert" className="allegro-error">{error}</p>}
    <div className="allegro-actions">
      {ticket && <a className="allegro-confirm" href={ticket.verification_url} target="_blank" rel="noreferrer">{t("allegro.confirmSeller")}</a>}
      <button className={ticket ? "" : "allegro-primary"} onClick={start} disabled={busy || action.busy}>
        {t(ticket ? "allegro.startAgain" : "allegro.enableRenewal")}
      </button>
    </div>
    {!ticket && <RefreshTokenForm version={version} action={action} />}
  </section>;
}
