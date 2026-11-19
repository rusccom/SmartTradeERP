import { useTenantPage } from "../api/useTenantPage";

export default function TenantPicker({ selection, select, changePage }) {
  const { items, error, loading, total } = useTenantPage(selection.page);
  if (error) return <p role="alert">{error}</p>;
  if (loading) return <p role="status">Loading organizations…</p>;
  if (!items.length) return <p>No organizations found.</p>;
  return (
    <div className="plugins-panel">
      <label htmlFor="plugin-organization">Organization</label>
      <select id="plugin-organization" value={selection.id} onChange={event => select(event.target.value)}>
        <option value="">Select an organization</option>
        {items.map(item => <option value={item.id} key={item.id}>{item.name}</option>)}
      </select>
      {selection.page > 1 && <button type="button" onClick={() => changePage(-1)}>Previous</button>}
      {selection.page * 50 < total && <button type="button" onClick={() => changePage(1)}>Next</button>}
    </div>
  );
}
