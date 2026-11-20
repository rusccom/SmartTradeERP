import { apiPaths } from "@smarterp/api/publicApi";
import useApiResource from "@smarterp/api/useApiResource";
import { useI18n } from "@smarterp/i18n/useI18n";
import AllegroMethods from "./AllegroMethods";
import AllegroAccount from "./AllegroAccount";
import useAllegroAction from "./useAllegroAction";

export default function AllegroSettings() {
  const { t } = useI18n();
  const { data, loading, error, retry } = useApiResource(apiPaths.allegro);
  const action = useAllegroAction(retry);
  return <>
    {loading && <p role="status">{t("allegro.loading")}</p>}
    {error && <p role="alert" className="allegro-error">{error} <button onClick={retry}>{t("allegro.retry")}</button></p>}
    {action.error && <p role="alert" className="allegro-notice allegro-error">{action.error}</p>}
    {data?.connection && <AllegroAccount account={data.connection} canManage={data.can_manage} action={action} />}
    {data && <AllegroMethods data={data} action={action} reload={retry} />}
  </>;
}
