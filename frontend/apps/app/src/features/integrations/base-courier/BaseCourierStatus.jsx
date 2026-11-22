import { useI18n } from "@smarterp/i18n/useI18n";
import BaseCourierActions from "./BaseCourierActions";

export default function BaseCourierStatus({ connection, action, canManage }) {
  const { t, locale } = useI18n();
  const verifiedAt = new Date(connection.verified_at).toLocaleString(locale);
  return (
    <section className="base-courier-status" aria-labelledby="base-courier-connected">
      <h2 id="base-courier-connected">{t("baseCourier.current")}</h2>
      <p>{connection.email}</p>
      <p>{t("baseCourier.lastVerified")} <time dateTime={connection.verified_at}>{verifiedAt}</time></p>
      {canManage && <BaseCourierActions connection={connection} action={action} />}
    </section>
  );
}
