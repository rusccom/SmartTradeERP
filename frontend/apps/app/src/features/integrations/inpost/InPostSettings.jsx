import { useI18n } from "@smarterp/i18n/useI18n";
import InPostForm from "./InPostForm";
import InPostStatus from "./InPostStatus";
import useInPostSettings from "./useInPostSettings";

export default function InPostSettings() {
  const { t } = useI18n();
  const { access, error, saved, retry, onSaved, onStart } = useInPostSettings();
  if (error) return <div role="alert">{error} <button className="inpost-retry" onClick={retry}>{t("inpost.retry")}</button></div>;
  if (!access) return <p role="status">{t("inpost.loading")}</p>;
  const { connection, can_manage } = access;
  return (
    <>
      <InPostStatus connection={connection} />
      {can_manage
        ? <InPostForm key={connection?.revision || 0} connection={connection} onSaved={onSaved} onStart={onStart} />
        : <p>{t("inpost.owner")}</p>}
      {saved && <p role="status" className="inpost-success">{t("inpost.saved")}</p>}
    </>
  );
}
