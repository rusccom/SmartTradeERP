import { ChevronLeft, ChevronRight } from "lucide-react";

import { useI18n } from "../../../i18n/useI18n";

function DataTablePagination({ table }) {
  const { t } = useI18n();
  const model = readPaginationModel(table, t);
  return (
    <div className="dt-pagination">
      <div className="dt-page-switcher">
        {model.canPaginate && <PageControls table={table} t={t} />}
        <span className="dt-page-range">{model.range}</span>
      </div>
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
    canPaginate: pageCount > 1,
    rowCount,
    range: readRange(pageIndex, pageSize, rowCount, t),
  };
}

function PageControls({ table, t }) {
  return (
    <div className="dt-page-stepper" role="navigation" aria-label={t("dataTable.pagination")}>
      <PagerButton label={t("dataTable.previousPage")} onClick={() => table.previousPage()} disabled={!table.getCanPreviousPage()}>
        <ChevronLeft size={19} aria-hidden="true" />
      </PagerButton>
      <PagerButton label={t("dataTable.nextPage")} onClick={() => table.nextPage()} disabled={!table.getCanNextPage()}>
        <ChevronRight size={19} aria-hidden="true" />
      </PagerButton>
    </div>
  );
}

function PagerButton({ onClick, disabled, label, children }) {
  return (
    <button className="dt-page-nav-button" type="button" onClick={onClick} disabled={disabled} aria-label={label} title={label}>
      {children}
    </button>
  );
}

function readRange(pageIndex, pageSize, rowCount, t) {
  if (rowCount <= 0) {
    return t("dataTable.rangeEmpty");
  }
  const start = Math.min(pageIndex * pageSize + 1, rowCount);
  const end = Math.min((pageIndex + 1) * pageSize, rowCount);
  return t("dataTable.range", { start, end });
}

export default DataTablePagination;
