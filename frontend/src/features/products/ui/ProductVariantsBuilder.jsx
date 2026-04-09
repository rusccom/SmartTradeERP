import { readVariantFields } from "../model/productForm";
import ProductVariantCard from "./ProductVariantCard";
import "./product-create.css";

function ProductVariantsBuilder({ canRemove, disabled, lockedVariantCount, onAddVariant, onRemoveVariant, onVariantChange, t, variants }) {
  const fields = readVariantFields();
  return (
    <section className="product-variants-builder">
      <header className="product-variants-head">
        <div>
          <h2 className="product-variants-title">{t("products.form.sections.variants")}</h2>
          <p className="product-variants-text">{t("products.form.multipleHint")}</p>
        </div>
        <button className="product-variants-add" type="button" onClick={onAddVariant} disabled={disabled}>{t("products.form.addVariant")}</button>
      </header>
      <div className="product-variants-list">
        {variants.map((variant, index) => <ProductVariantCard key={variant.id} variant={variant} index={index} fields={fields} t={t} onChange={onVariantChange} onRemove={onRemoveVariant} locked={index < lockedVariantCount} canRemove={canRemove} />)}
      </div>
    </section>
  );
}

export default ProductVariantsBuilder;
