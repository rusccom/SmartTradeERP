function DataTableError({ message, onRetry }) {
  return (
    <div className="dt-error">
      <span>{message}</span>
      {onRetry && (
        <button className="dt-page-btn dt-error-btn" type="button" onClick={onRetry}>
          Повторить
        </button>
      )}
    </div>
  );
}

export default DataTableError;
