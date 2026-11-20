import { useLocation, useNavigate, useParams } from "react-router-dom";
import { useI18n } from "@smarterp/i18n/useI18n";
import { apiPaths } from "@smarterp/api/publicApi";
import { orderPath } from "../menu";
import OrderEntryForm from "./OrderEntryForm";
import useOrderResource from "./useOrderResource";
import "./order-entry.css";

export default function OrderEntryPage() {
  const { id } = useParams();
  const { t } = useI18n();
  const navigate = useNavigate();
  const location = useLocation();
  const result = useOrderResource(id ? apiPaths.order(id) : "");
  const saved = (card, error) => {
    if (!id) navigate(`${orderPath}/${card.id}`, { replace: true, state: { error } });
  };
  if (result.error) return <section role="alert">{result.error} <button type="button" className="order-btn" onClick={result.retry}>{t("entry.retry")}</button></section>;
  if (id && !result.data) return <p role="status">{t("entry.loading")}</p>;
  return <OrderEntryForm key={id || "new"} initial={result.data} onPersist={saved} initialError={location.state?.error} />;
}
