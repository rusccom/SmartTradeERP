import useApiResource from "@smarterp/api/useApiResource";
import { apiPaths } from "@smarterp/api/publicApi";
import { useI18n } from "@smarterp/i18n/useI18n";
import BaseCourierForm from "./BaseCourierForm";
import BaseCourierStatus from "./BaseCourierStatus";
import useConnectionAction from "./useConnectionAction";

export default function BaseCourierSettings() {
  const { t } = useI18n();
  const { data, error, loading, retry } = useApiResource(apiPaths.baseCourierSettings);
  const action = useConnectionAction(retry);
  return (
    <div className="base-courier-connection">
      {loading && <p role="status">{t("baseCourier.loading")}</p>}
      {action.busy && <p role="status">{t("baseCourier.working")}</p>}
      {(error || action.error) && <div className="base-courier-error" role="alert">
        <p>{error || action.error}</p><button onClick={action.reload}>{t("baseCourier.reload")}</button>
      </div>}
      {action.notice && <p className="base-courier-notice" role="status">{t(action.notice)}</p>}
      {data?.connection && <BaseCourierStatus connection={data.connection} action={action} canManage={data.can_manage} />}
      {data?.can_manage && <BaseCourierForm key={data.connection?.version || "new"} connection={data.connection} action={action} />}
      {data && !data.can_manage && <p>{t("baseCourier.ownerOnly")}</p>}
    </div>
  );
}
