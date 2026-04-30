import { readVariantFields } from "../model/productForm";
import ProductVariantCard from "./ProductVariantCard";
import "./product-create.css";

function ProductVariantsBuilder(props) {
  const fields = readVariantFields(props.priceStep);
  return (
    <div className="product-variants-builder">
      <div className="product-variants-list">
        {props.variants.map((variant, index) => <ProductVariantCard key={variant.id} variant={variant} index={index} fields={fields} t={props.t} onChange={props.onVariantChange} onRemove={props.onRemoveVariant} locked={index < props.lockedVariantCount} canRemove={props.canRemove} />)}
      </div>
    </div>
  );
}

export default ProductVariantsBuilder;
