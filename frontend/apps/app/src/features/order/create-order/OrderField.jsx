export default function OrderField({ label, value, onChange, options, ...input }) {
  const control = { ...input, value: value ?? "", onChange: (event) => onChange(event.target.value) };
  return <label className="order-field">
    <span>{label}</span>
    {options ? <select {...control}>{options.map(([code, name]) => <option key={code} value={code}>{name}</option>)}</select> : <input {...control} />}
  </label>;
}
