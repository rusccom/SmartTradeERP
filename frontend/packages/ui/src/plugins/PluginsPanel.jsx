import { useI18n } from "@smarterp/i18n/useI18n";
import { usePlugins } from "@smarterp/model/plugins/usePlugins";
import PluginControl from "./PluginControl";
import "./plugins.css";

export default function PluginsPanel({ tenantId = "" }) {
  const { t } = useI18n();
  const { access, busy, error, save } = usePlugins(tenantId);
  return (
    <section className="plugins-panel" aria-label={t("plugins.title")}>
      <h3>{t("plugins.title")}</h3>
      {error && <p role="alert">{error}</p>}
      {!error && !access && <p role="status">{t("plugins.loading")}</p>}
      {access?.plugins.length === 0 && <p>{t("plugins.empty")}</p>}
      {access?.plugins.map(plugin => (
        <PluginControl key={plugin.key} plugin={plugin} grant={Boolean(tenantId)}
          disabled={busy || !access.can_manage} onChange={save} />
      ))}
    </section>
  );
}
