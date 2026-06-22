import FormField from "../../../shared/ui/form-modal/FormField";

const nameField = {
  name: "name",
  labelKey: "products.form.name",
  type: "text",
  required: true,
  autoFocus: true,
};

function ProductBasicSection({ form, onChange, t }) {
  return (
    <section className="product-create-card">
      <h3 className="product-create-card-title">{t("products.form.sections.product")}</h3>
      <FormField field={nameField} value={form.name} onChange={onChange} t={t} />
    </section>
  );
}

export default ProductBasicSection;
