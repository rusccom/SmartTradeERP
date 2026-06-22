import { ImagePlus } from "lucide-react";

function ProductMediaSection({ t }) {
  return (
    <section className="product-create-card">
      <h3 className="product-create-card-title">{t("products.form.sections.media")}</h3>
      <div className="product-media-placeholder">
        <ImagePlus aria-hidden="true" size={24} strokeWidth={1.8} />
        <div className="product-media-actions">
          <span className="product-media-button">{t("products.form.mediaUpload")}</span>
          <span className="product-media-button">{t("products.form.mediaSelect")}</span>
        </div>
        <p>{t("products.form.mediaHint")}</p>
      </div>
    </section>
  );
}

export default ProductMediaSection;
