function FormField({ field, onChange, t, value }) {
  if (field.type === "checkbox") return renderCheckbox(field, onChange, t, value);
  if (field.type === "textarea") return renderTextarea(field, onChange, t, value);
  if (field.type === "select") return renderSelect(field, onChange, t, value);
  return renderInput(field, onChange, t, value);
}

function renderCheckbox(field, onChange, t, value) {
  return <label className="shared-form-checkbox"><input name={field.name} type="checkbox" checked={Boolean(value)} onChange={onChange} /><span>{t(field.labelKey)}</span></label>;
}

function renderInput(field, onChange, t, value) {
  return <label className="shared-form-field"><span className="shared-form-field-label">{t(field.labelKey)}</span><input className="shared-form-field-input" type={field.type} {...buildSharedProps(field, onChange, value)} /></label>;
}

function renderTextarea(field, onChange, t, value) {
  return <label className="shared-form-field"><span className="shared-form-field-label">{t(field.labelKey)}</span><textarea className="shared-form-field-input shared-form-field-textarea" {...buildSharedProps(field, onChange, value)} rows={field.rows || 4} /></label>;
}

function renderSelect(field, onChange, t, value) {
  return <label className="shared-form-field"><span className="shared-form-field-label">{t(field.labelKey)}</span><select className="shared-form-field-input" {...buildSharedProps(field, onChange, value)}>{renderOptions(field.options, t)}</select></label>;
}

function buildSharedProps(field, onChange, value) {
  return {
    autoComplete: field.autoComplete || "off",
    autoFocus: field.autoFocus === true,
    min: field.min,
    name: field.name,
    onChange,
    placeholder: field.placeholder,
    required: field.required === true,
    step: field.step,
    value,
  };
}

function renderOptions(options, t) {
  return options?.map((option) => <option key={option.value} value={option.value}>{t(option.labelKey || option.label)}</option>);
}

export default FormField;
