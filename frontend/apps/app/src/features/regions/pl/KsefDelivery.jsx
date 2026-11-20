import { useState } from "react";
import { Link } from "react-router-dom";
import { orderPath } from "../../order/menu";
import { CheckCircle2 } from "lucide-react";
import { useI18n } from "@smarterp/i18n/useI18n";
import { apiPaths } from "@smarterp/api/publicApi";
import ConfirmKsef from "./ConfirmKsef";
import useApiResource from "@smarterp/api/useApiResource";
import useKsefDelivery from "./useKsefDelivery";

export default function KsefDelivery({ invoice }) {
  const { t } = useI18n();
  const [confirm, setConfirm] = useState(false);
  const status = useKsefDelivery(invoice.id);
  const profile = useApiResource(apiPaths.polishSeller);
  const d = status.data;
  const env = d?.environment || profile.data?.seller?.environment;
  const accepted = d?.state === "accepted" && Boolean(d.ksef_number);
  const check = ["pending", "sending", "uncertain"].includes(d?.state);
  const action = () => check ? status.run("check") : setConfirm(true);
  const send = () => { setConfirm(false); status.run("send", env); };
  return <div className="order-ksef">
    {accepted ? <span className="order-ksef-accepted"><CheckCircle2 size={17}/>{t("entry.accepted", { env: t(`entry.${env}`) })} · {d.ksef_number}</span>
      : <button className="order-btn" type="button" disabled={status.busy || !d || !env} onClick={action}>{t(check ? "entry.checkKsef" : "entry.sendKsef")}</button>}
    {d && !accepted && <small role="status">{t(`entry.state.${d.state}`)}{d.message && ` · ${d.message}`}</small>}
    {(status.error || profile.error) && <button type="button" className="order-link" onClick={() => { status.run("read"); profile.retry(); }}>{status.error || profile.error} · {t("entry.retry")}</button>}
    {!accepted && !profile.loading && !profile.data?.seller?.token_configured && <Link className="order-link" to={`${orderPath}/settings/poland`}>{t("entry.plSettings")}</Link>}
    {confirm && <ConfirmKsef number={invoice.number} environment={env} onCancel={() => setConfirm(false)} onConfirm={send} />}
  </div>;
}
