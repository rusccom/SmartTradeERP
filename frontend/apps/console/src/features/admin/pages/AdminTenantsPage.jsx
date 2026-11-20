import { useState } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";
import PluginsPanel from "@smarterp/ui/plugins/PluginsPanel";
import TenantPicker from "./TenantPicker";

export default function AdminTenantsPage() {
  const { t } = useI18n();
  const [selection, setSelection] = useState({ id: "", page: 1 });
  return (
    <section className="placeholder">
      <h2>{t("admin.organizations.title")}</h2>
      <TenantPicker selection={selection} select={id => setSelection({ ...selection, id })}
        changePage={delta => setSelection({ page: selection.page + delta, id: "" })} />
      {selection.id && <PluginsPanel key={selection.id} tenantId={selection.id} />}
    </section>
  );
}
