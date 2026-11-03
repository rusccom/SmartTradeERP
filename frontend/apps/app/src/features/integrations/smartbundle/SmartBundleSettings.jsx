import { useBundleSettings } from "./useBundleSettings";
import BundleSlotRule from "./BundleSlotRule";
import "./smartbundle.css";

export default function SmartBundleSettings({ bundle }) {
  const { editor, message, busy, patch, save } = useBundleSettings(bundle);
  const s = editor?.settings;
  return (
    <details className="smartbundle-settings"><summary>SmartBundle · Storefront settings</summary>
      {s && <form onSubmit={event => { event.preventDefault(); save(); }}>
        <label>Pricing <select value={s.mode} onChange={e => patch({ mode: e.target.value })}>
          <option value="fixed">Fixed · ERP variant price</option><option value="dynamic">Sum of selected components</option>
        </select></label>
        <label>Bundle discount (%) <input type="number" min="0" max="100" step="0.01" required value={s.discount} onChange={e => patch({ discount: e.target.value })} /></label>
        {editor.components.map(slot => <BundleSlotRule key={slot.key} slot={slot} rule={s.slots[slot.key]} mode={s.mode}
          onChange={rule => patch({ slots: { ...s.slots, [slot.key]: rule } })} />)}
        <button disabled={busy}>Save storefront settings</button>
      </form>}
      <p role="status">{message}</p>
    </details>
  );
}
