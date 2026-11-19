export default function PluginControl({ plugin, grant, disabled, onChange }) {
  const locked = !grant && !plugin.granted;
  const label = grant ? "Grant access" : "Enable";
  return (
    <div className="plugin-control">
      <label>
        <input type="checkbox" checked={grant ? plugin.granted : plugin.enabled}
          disabled={disabled || locked} onChange={event => onChange(plugin, event.target.checked)} />
        {label} {plugin.name}
      </label>
      {locked && <small>License required</small>}
      {grant && plugin.granted && <small>{plugin.enabled ? "Enabled by store" : "Disabled by store"}</small>}
    </div>
  );
}
