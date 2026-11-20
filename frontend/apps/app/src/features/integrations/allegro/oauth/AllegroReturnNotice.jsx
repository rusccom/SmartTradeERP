import { useSearchParams } from "react-router-dom";
import { useI18n } from "@smarterp/i18n/useI18n";

const messages = {
  cancelled: "allegro.cancelled", expired: "allegro.expired", permissions: "allegro.permissions",
  different: "allegro.different", setup: "allegro.setup", unavailable: "allegro.unavailable",
  invalid_token: "allegro.invalidGrant", changed: "allegro.changed", rate_limit: "allegro.rateLimit",
};

export default function AllegroReturnNotice() {
  const { t } = useI18n();
  const [params, setParams] = useSearchParams();
  const result = params.get("result");
  if (!result || result === "connected") return null;
  const message = Object.hasOwn(messages, result) ? messages[result] : "allegro.failed";
  return <div className="allegro-notice" role="status">
    <span>{t(message)}</span>
    <button onClick={() => setParams({}, { replace: true })}>{t("allegro.dismiss")}</button>
  </div>;
}
