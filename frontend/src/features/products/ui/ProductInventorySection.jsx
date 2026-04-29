import FormField from "../../../shared/ui/form-modal/FormField";

const inventoryFields = [
  { name: "skuCode", labelKey: "products.form.skuCode", type: "text" },
  { name: "barcode", labelKey: "products.form.barcode", type: "text" },
];

function ProductInventorySection({ form, onChange, t }) {
  return (
    <section className="product-create-card">
      <h3 className="product-create-card-title">{t("products.form.sections.inventory")}</h3>
      <div className="product-create-grid">
        {inventoryFields.map((field) => <FormField key={field.name} field={field} value={form[field.name]} onChange={onChange} t={t} />)}
      </div>
    </section>
  );
}

export default ProductInventorySection;
