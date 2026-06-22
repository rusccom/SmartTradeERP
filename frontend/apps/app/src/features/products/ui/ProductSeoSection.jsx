import ProductSeoPreview from "./ProductSeoPreview";
import SeoField from "./SeoField";

const slugField = { name: "slug", labelKey: "products.form.slug", type: "text" };
const titleField = { name: "seoTitle", labelKey: "products.form.seoTitle", type: "text", max: 70 };
const descriptionField = { name: "seoDescription", labelKey: "products.form.seoDescription", type: "textarea", rows: 2, max: 160 };

function ProductSeoSection({ form, onChange, t }) {
  return (
    <section className="product-create-card product-seo-section">
      <h3 className="product-create-card-title">{t("products.form.sections.seo")}</h3>
      <ProductSeoPreview form={form} t={t} />
      <div className="product-create-grid">
        <SeoField field={slugField} value={form.slug} onChange={onChange} t={t} />
        <SeoField field={titleField} value={form.seoTitle} onChange={onChange} t={t} />
      </div>
      <SeoField field={descriptionField} value={form.seoDescription} onChange={onChange} t={t} />
    </section>
  );
}

export default ProductSeoSection;
