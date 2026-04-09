import { useI18n } from "../../i18n/useI18n";

function DataTableError({ message, onRetry }) {
  const { t } = useI18n();
  return (
    <div className="dt-error">
      <span>{message}</span>
      {onRetry && (
        <button className="dt-page-btn dt-error-btn" type="button" onClick={onRetry}>
          {t("dataTable.retry")}
        </button>
      )}
    </div>
  );
}

export default DataTableError;
