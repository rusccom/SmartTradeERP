import { useI18n } from "@smarterp/i18n/useI18n";

export default function InPostStatus({ connection }) {
  const { t, locale } = useI18n();
  if (!connection) return <p className="inpost-status">{t("inpost.notConnected")}</p>;
  const { organization_name, organization_id, environment, verified_at, services } = connection;
  return (
    <section className="inpost-status inpost-status--verified" aria-label={organization_name}>
      <h2>{organization_name}</h2>
      <p>{t(`inpost.${environment}`)} · {t("inpost.organization")} {organization_id}</p>
      <p>{t("inpost.verified")} <time dateTime={verified_at}>{new Date(verified_at).toLocaleString(locale)}</time></p>
      <h3>{t("inpost.services")}</h3>
      {services.length
        ? <ul>{services.map(code => <li key={code}>{code}</li>)}</ul>
        : <p>{t("inpost.noServices")}</p>}
    </section>
  );
}
