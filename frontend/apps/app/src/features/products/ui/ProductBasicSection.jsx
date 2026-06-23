import FormField from "@smarterp/ui/form-modal/FormField";

const nameField = {
  name: "name",
  labelKey: "products.form.name",
  type: "text",
  required: true,
  autoFocus: true,
};

function ProductBasicSection({ form, onChange, t }) {
  return (
    <section className="product-card">
      <h3 className="product-card__title">{t("products.form.sections.product")}</h3>
      <FormField field={nameField} value={form.name} onChange={onChange} t={t} />
    </section>
  );
}

export default ProductBasicSection;
