import FormField from "../../../shared/ui/form-modal/FormField";

const priceFields = [
  { name: "price", labelKey: "products.form.price", type: "number", min: "0", step: "0.01", required: true },
  { name: "unit", labelKey: "products.form.unit", type: "text", required: true },
];

const unitFields = [
  { name: "unit", labelKey: "products.form.unit", type: "text", required: true },
];

function ProductPriceSection({ form, hasVariants, onChange, t }) {
  const fields = hasVariants ? unitFields : priceFields;
  return (
    <section className="product-create-card">
      <h3 className="product-create-card-title">{t(readTitleKey(hasVariants))}</h3>
      <div className="product-create-grid">
        {fields.map((field) => <FormField key={field.name} field={field} value={form[field.name]} onChange={onChange} t={t} />)}
      </div>
    </section>
  );
}

function readTitleKey(hasVariants) {
  return hasVariants ? "products.form.sections.unit" : "products.form.sections.price";
}

export default ProductPriceSection;
