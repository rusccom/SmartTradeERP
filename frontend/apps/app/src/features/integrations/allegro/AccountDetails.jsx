import { Fragment } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";

export default function AccountDetails({ account }) {
  const { t, locale } = useI18n();
  const dates = [["allegro.verifiedAt", account.verified_at], ["allegro.refreshedAt", account.refreshed_at]];
  if (account.method === "token") dates.push(["allegro.expires", account.expires_at]);
  const method = account.auto_renew ? "allegro.privateApp" : "allegro.tokenLabel";
  return <dl>
    <dt>{t("allegro.method")}</dt><dd>{t(account.method === "oauth" ? "allegro.byLogin" : method)}
      {account.auto_renew && !account.needs_reconnect && <div className="allegro-verified">{t("allegro.autoRenew")}</div>}
    </dd>
    {dates.filter(([, value]) => value).map(([key, value]) => <Fragment key={key}>
      <dt>{t(key)}</dt><dd><time dateTime={value}>{new Date(value).toLocaleString(locale)}</time></dd>
    </Fragment>)}
  </dl>;
}
