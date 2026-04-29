import { readVariantFields } from "../model/productForm";
import ProductVariantCard from "./ProductVariantCard";
import "./product-create.css";

function ProductVariantsBuilder({ canRemove, lockedVariantCount, onRemoveVariant, onVariantChange, t, variants }) {
  const fields = readVariantFields();
  return (
    <div className="product-variants-builder">
      <div className="product-variants-list">
        {variants.map((variant, index) => <ProductVariantCard key={variant.id} variant={variant} index={index} fields={fields} t={t} onChange={onVariantChange} onRemove={onRemoveVariant} locked={index < lockedVariantCount} canRemove={canRemove} />)}
      </div>
    </div>
  );
}

export default ProductVariantsBuilder;
