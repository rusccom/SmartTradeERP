import { ChevronDown, Search } from "lucide-react";
import { useI18n } from "@smarterp/i18n/useI18n";

function ComponentOptionsToolbar(props) {
  const { t } = useI18n();
  const { canLoadMore, loading, onLoadMore, onSearch, search } = props;
  return (
    <div className="bundle-component-toolbar">
      <label className="bundle-component-search">
        <Search size={15} />
        <input value={search} onChange={(event) => onSearch(event.target.value)}
          aria-label={t("bundles.editor.search")} placeholder={t("bundles.editor.search")} />
      </label>
      <button type="button" onClick={onLoadMore} disabled={loading || !canLoadMore}>
        <ChevronDown size={16} /> {t("bundles.editor.loadMore")}
      </button>
    </div>
  );
}

export default ComponentOptionsToolbar;
