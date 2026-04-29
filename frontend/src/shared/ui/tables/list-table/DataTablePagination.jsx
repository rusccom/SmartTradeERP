import { useI18n } from "../../../i18n/useI18n";

function DataTablePagination({ table }) {
  const { t } = useI18n();
  const model = readPaginationModel(table, t);
  return (
    <div className="dt-pagination">
      <span className="dt-page-info">{model.range}</span>
      <PageControls table={table} model={model} />
      <PageSizeSelect table={table} pageSize={model.pageSize} t={t} />
    </div>
  );
}

function readPaginationModel(table, t) {
  const pagination = table.getState().pagination;
  const pageIndex = pagination.pageIndex;
  const pageSize = pagination.pageSize;
  const rowCount = table.getRowCount();
  const pageCount = table.getPageCount();
  return {
    pageIndex,
    pageSize,
    pageCount,
    range: readRange(pageIndex, pageSize, rowCount, t),
    pages: readVisiblePages(pageIndex, pageCount),
  };
}

function PageControls({ table, model }) {
  return (
    <div className="dt-page-controls">
      <PagerButton onClick={() => table.setPageIndex(0)} disabled={!table.getCanPreviousPage()} title="<<" />
      <PagerButton onClick={() => table.previousPage()} disabled={!table.getCanPreviousPage()} title="<" />
      {model.pages.map((page) => (
        <PageNumber key={page} table={table} page={page} pageIndex={model.pageIndex} />
      ))}
      <PagerButton onClick={() => table.nextPage()} disabled={!table.getCanNextPage()} title=">" />
      <PagerButton onClick={() => table.setPageIndex(readLastPageIndex(model.pageCount))} disabled={!table.getCanNextPage()} title=">>" />
    </div>
  );
}

function PageNumber({ table, page, pageIndex }) {
  return (
    <button className={readPageButtonClass(pageIndex, page)} type="button" onClick={() => table.setPageIndex(page)}>
      {page + 1}
    </button>
  );
}

function PagerButton({ onClick, disabled, title }) {
  return (
    <button className="dt-page-btn" type="button" onClick={onClick} disabled={disabled}>
      {title}
    </button>
  );
}

function PageSizeSelect({ table, pageSize, t }) {
  return (
    <select className="dt-page-size" value={pageSize} onChange={(event) => table.setPageSize(Number(event.target.value))}>
      {[10, 20, 50, 100].map((size) => (
        <option key={size} value={size}>
          {t("dataTable.pageSize", { size })}
        </option>
      ))}
    </select>
  );
}

function readRange(pageIndex, pageSize, rowCount, t) {
  if (rowCount <= 0) {
    return t("dataTable.rangeEmpty");
  }
  const start = pageIndex * pageSize + 1;
  const end = Math.min((pageIndex + 1) * pageSize, rowCount);
  return t("dataTable.range", { start, end, count: rowCount });
}

function readVisiblePages(pageIndex, pageCount) {
  if (pageCount <= 0) {
    return [];
  }
  const size = 5;
  const start = Math.max(0, Math.min(pageIndex - 2, pageCount - size));
  const end = Math.min(pageCount, start + size);
  return Array.from({ length: end - start }, (_, index) => start + index);
}

function readLastPageIndex(pageCount) {
  return pageCount > 0 ? pageCount - 1 : 0;
}

function readPageButtonClass(pageIndex, page) {
  const active = pageIndex === page ? "dt-page-btn--active" : "";
  return `dt-page-btn ${active}`.trim();
}

export default DataTablePagination;
