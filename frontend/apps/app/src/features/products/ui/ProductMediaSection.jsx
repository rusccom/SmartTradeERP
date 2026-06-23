import MediaManager from "../../media/ui/MediaManager";

function ProductMediaSection({ productId, t }) {
  return (
    <section className="product-create-card">
      <h3 className="product-create-card-title">{t("products.form.sections.media")}</h3>
      {productId ? (
        <MediaManager kind="product" ownerId={productId} t={t} />
      ) : (
        <p className="product-media-after-save">{t("products.form.mediaAfterSave")}</p>
      )}
    </section>
  );
}

export default ProductMediaSection;
