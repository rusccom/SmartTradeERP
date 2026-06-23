import { slugify } from "../model/slugify";

function ProductSeoPreview({ form, t }) {
  const title = pick(form.seoTitle, form.name) || t("products.form.seoPreviewTitle");
  const handle = slugify(form.slug) || slugify(form.name) || t("products.form.seoPreviewHandle");
  const description = pick(form.seoDescription) || t("products.form.seoPreviewDescription");
  return (
    <div className="product-seo-preview">
      <span className="product-seo-preview-caption">{t("products.form.seoPreview")}</span>
      <div className="product-seo-preview-url">{t("products.form.seoPreviewHost")} › products › {handle}</div>
      <div className="product-seo-preview-title">{title}</div>
      <div className="product-seo-preview-desc">{description}</div>
    </div>
  );
}

function pick(...values) {
  for (const value of values) {
    const trimmed = String(value || "").trim();
    if (trimmed) return trimmed;
  }
  return "";
}

export default ProductSeoPreview;
