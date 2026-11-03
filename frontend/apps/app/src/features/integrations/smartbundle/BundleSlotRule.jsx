export default function BundleSlotRule({ slot, rule, mode, onChange }) {
  const toggle = (id, checked) => onChange({ ...rule, options: checked ? [...rule.options, id] : rule.options.filter(value => value !== id) });
  return (
    <fieldset><legend>{slot.name}</legend>
      {slot.options.map(option => <label key={option.id}>
        <input type="checkbox" disabled={option.id === slot.key} checked={option.id === slot.key || rule.options.includes(option.id)}
          onChange={event => toggle(option.id, event.target.checked)} /> {option.name}
      </label>)}
      {mode === "dynamic" && <label>Component discount (%) <input type="number" min="0" max="100" step="0.01" required
        value={rule.discount} onChange={e => onChange({ ...rule, discount: e.target.value })} /></label>}
    </fieldset>
  );
}
