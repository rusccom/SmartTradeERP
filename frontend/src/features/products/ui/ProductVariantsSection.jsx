function ProductVariantsSection({ children, disabled, hasVariants, onAddVariant, t }) {
  return (
    <section className="product-variants-section">
      <header className="product-variants-section-head">
        <div>
          <h3 className="product-variants-section-title">{t("products.form.sections.variants")}</h3>
          <p className="product-variants-text">{readHint(hasVariants, t)}</p>
        </div>
        <button className="product-variants-add" type="button" onClick={onAddVariant} disabled={disabled}>
          {t("products.form.addVariant")}
        </button>
      </header>
      {children}
    </section>
  );
}

function readHint(hasVariants, t) {
  return t(hasVariants ? "products.form.multipleHint" : "products.form.variantsEmptyHint");
}

export default ProductVariantsSection;
