function SeoField({ field, onChange, t, value }) {
  const length = String(value || "").length;
  const over = Boolean(field.max) && length > field.max;
  return (
    <label className="shared-form-field">
      <span className="seo-field-head">
        <span className="shared-form-field-label">{t(field.labelKey)}</span>
        {field.max ? <span className={counterClass(over)}>{length} / {field.max}</span> : null}
      </span>
      {renderControl(field, onChange, value)}
    </label>
  );
}

function renderControl(field, onChange, value) {
  if (field.type === "textarea") {
    return (
      <textarea
        className="shared-form-field-input seo-textarea"
        name={field.name}
        rows={field.rows || 2}
        value={value}
        onChange={onChange}
      />
    );
  }
  return (
    <input
      className="shared-form-field-input"
      type="text"
      name={field.name}
      value={value}
      onChange={onChange}
      autoComplete="off"
    />
  );
}

function counterClass(over) {
  return over ? "seo-counter is-over" : "seo-counter";
}

export default SeoField;
