import { Link } from "react-router-dom";
import { useI18n } from "@smarterp/i18n/useI18n";
import { apiPaths } from "@smarterp/api/publicApi";
import { orderPath } from "../../order/menu";
import useOrderResource from "../../order/create-order/useOrderResource";
import SellerProfileForm from "./SellerProfileForm";
import "../../order/create-order/order-entry.css";

export default function PolishOrderSettings() {
  const { t } = useI18n();
  const profile = useOrderResource(apiPaths.polishSeller);
  return <section className="order-entry"><Link className="order-back" to={`${orderPath}/settings`}>{t("orders.settings.title")}</Link>
    <h1>{t("entry.plSettings")}</h1>
    {profile.loading && <p role="status">{t("entry.loading")}</p>}
    {profile.error && <p role="alert">{profile.error} <button className="order-link" onClick={profile.retry}>{t("entry.retry")}</button></p>}
    {profile.data && <SellerProfileForm initial={profile.data.seller} canManage={profile.data.can_manage} />}
  </section>;
}
