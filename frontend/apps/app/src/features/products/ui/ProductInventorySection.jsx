import FormField from "@smarterp/ui/form-modal/FormField";

const inventoryFields = [
  { name: "skuCode", labelKey: "products.form.skuCode", type: "text" },
  { name: "barcode", labelKey: "products.form.barcode", type: "text" },
];

function ProductInventorySection({ form, onChange, t }) {
  return (
    <section className="product-card">
      <h3 className="product-card__title">{t("products.form.sections.inventory")}</h3>
      <div className="product-create-grid">
        {inventoryFields.map((field) => <FormField key={field.name} field={field} value={form[field.name]} onChange={onChange} t={t} />)}
      </div>
    </section>
  );
}

export default ProductInventorySection;
