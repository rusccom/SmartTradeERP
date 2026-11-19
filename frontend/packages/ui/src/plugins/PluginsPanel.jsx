import { usePlugins } from "@smarterp/model/plugins/usePlugins";
import PluginControl from "./PluginControl";
import "./plugins.css";

export default function PluginsPanel({ tenantId = "" }) {
  const { access, busy, error, save } = usePlugins(tenantId);
  return (
    <section className="plugins-panel" aria-label="Plugins">
      <h3>Plugins</h3>
      {error && <p role="alert">{error}</p>}
      {!error && !access && <p role="status">Loading plugins…</p>}
      {access?.plugins.length === 0 && <p>No plugins available.</p>}
      {access?.plugins.map(plugin => (
        <PluginControl key={plugin.key} plugin={plugin} grant={Boolean(tenantId)}
          disabled={busy || !access.can_manage} onChange={save} />
      ))}
    </section>
  );
}
