import { useState } from "react";
import PluginsPanel from "@smarterp/ui/plugins/PluginsPanel";
import TenantPicker from "./TenantPicker";

export default function AdminTenantsPage() {
  const [selection, setSelection] = useState({ id: "", page: 1 });
  return (
    <section className="placeholder">
      <h2>Organizations</h2>
      <TenantPicker selection={selection} select={id => setSelection({ ...selection, id })}
        changePage={delta => setSelection({ page: selection.page + delta, id: "" })} />
      {selection.id && <PluginsPanel key={selection.id} tenantId={selection.id} />}
    </section>
  );
}
