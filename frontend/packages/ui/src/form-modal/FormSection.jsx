import FormField from "./FormField";

function FormSection({ form, onChange, section, t }) {
  return (
    <fieldset className="shared-form-section">
      <legend className="shared-form-section-title">{t(section.titleKey)}</legend>
      <div className="shared-form-section-grid">
        {section.fields.map((field) => <FormField key={field.name} field={field} value={form[field.name]} onChange={onChange} t={t} />)}
      </div>
    </fieldset>
  );
}

export default FormSection;
