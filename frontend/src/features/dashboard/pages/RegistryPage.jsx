import { useParams } from "react-router-dom";

import { findPageByKey } from "../registry";
import ListPage from "../components/ListPage";

function RegistryPage({ pageKey }) {
  const params = useParams();
  const key = pageKey || params.pageKey;
  const page = findPageByKey(key);

  if (!page) {
    return (
      <section className="placeholder">
        <h2>Page not found</h2>
        <p>No configuration for "{key}".</p>
      </section>
    );
  }

  return (
    <ListPage
      title={page.title}
      columns={page.columns}
      data={page.mock}
    />
  );
}

export default RegistryPage;
