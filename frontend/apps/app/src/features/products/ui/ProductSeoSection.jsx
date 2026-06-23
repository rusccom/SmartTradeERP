import ProductSeoPreview from "./ProductSeoPreview";
import SeoField from "./SeoField";
import SlugField from "./SlugField";

const titleField = { name: "seoTitle", labelKey: "products.form.seoTitle", type: "text", max: 70 };
const descriptionField = { name: "seoDescription", labelKey: "products.form.seoDescription", type: "textarea", rows: 2, max: 160 };

function ProductSeoSection({ form, onChange, t }) {
  return (
    <section className="product-create-card product-seo-section">
      <h3 className="product-create-card-title">{t("products.form.sections.seo")}</h3>
      <SeoField field={titleField} value={form.seoTitle} onChange={onChange} t={t} />
      <SeoField field={descriptionField} value={form.seoDescription} onChange={onChange} t={t} />
      <SlugField value={form.slug} onChange={onChange} t={t} />
      <ProductSeoPreview form={form} t={t} />
    </section>
  );
}

export default ProductSeoSection;
