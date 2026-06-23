function ProductVariantsSection({ children, disabled, hasVariants, onAddVariant, t }) {
  return (
    <section className="product-card">
      <header className="product-card__head">
        <div>
          <h3 className="product-card__title">{t("products.form.sections.variants")}</h3>
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
