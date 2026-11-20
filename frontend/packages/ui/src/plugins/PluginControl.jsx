import { useI18n } from "@smarterp/i18n/useI18n";

export default function PluginControl({ plugin, grant, disabled, onChange }) {
  const { t } = useI18n();
  const locked = !grant && !plugin.granted;
  const label = t(grant ? "plugins.grant" : "plugins.enable", { name: plugin.name });
  return (
    <div className="plugin-control">
      <label>
        <input type="checkbox" checked={grant ? plugin.granted : plugin.enabled}
          disabled={disabled || locked} onChange={event => onChange(plugin, event.target.checked)} />
        {label}
      </label>
      {locked && <small>{t("plugins.licenseRequired")}</small>}
      {grant && plugin.granted && <small>{t(plugin.enabled ? "plugins.enabledByStore" : "plugins.disabledByStore")}</small>}
    </div>
  );
}
