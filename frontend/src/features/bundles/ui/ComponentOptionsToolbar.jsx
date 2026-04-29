import { ChevronDown, Search } from "lucide-react";

function ComponentOptionsToolbar(props) {
  const { canLoadMore, loading, onLoadMore, onSearch, search } = props;
  return (
    <div className="bundle-component-toolbar">
      <label className="bundle-component-search">
        <Search size={15} />
        <input value={search} onChange={(event) => onSearch(event.target.value)} placeholder="Search components" />
      </label>
      <button type="button" onClick={onLoadMore} disabled={loading || !canLoadMore}>
        <ChevronDown size={16} /> Load more
      </button>
    </div>
  );
}

export default ComponentOptionsToolbar;
