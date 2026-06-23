import ProductBasicSection from "./ProductBasicSection";
import ProductInventorySection from "./ProductInventorySection";
import ProductMediaSection from "./ProductMediaSection";
import ProductPriceSection from "./ProductPriceSection";
import ProductSeoSection from "./ProductSeoSection";
import ProductVariantsSection from "./ProductVariantsSection";
import ProductVariantsBuilder from "./ProductVariantsBuilder";
import { decimalStep } from "../../currencies/model/formatMoney";
import { useCurrencies } from "../../currencies/model/useCurrencies";

function ProductCreateForm({ formId, state, t }) {
  const { defaultCurrency } = useCurrencies();
  const hasVariants = state.form.variantMode === "multiple";
  const priceStep = decimalStep(defaultCurrency);
  return (
    <form id={formId} className="product-create-form" data-density="compact" onSubmit={state.handleSubmit}>
      <ProductBasicSection form={state.form} onChange={state.handleChange} t={t} />
      <ProductMediaSection productId={state.productId} t={t} />
      <ProductPriceSection form={state.form} hasVariants={hasVariants} onChange={state.handleChange} priceStep={priceStep} t={t} />
      {!hasVariants && <ProductInventorySection form={state.form} onChange={state.handleChange} t={t} />}
      <ProductVariantsSection hasVariants={hasVariants} onAddVariant={state.handleAddVariant} disabled={state.hasPendingProduct} t={t}>
        {readBottomSlot(state, t, priceStep)}
      </ProductVariantsSection>
      <ProductSeoSection form={state.form} onChange={state.handleChange} t={t} />
    </form>
  );
}

function readBottomSlot(state, t, priceStep) {
  if (state.form.variantMode !== "multiple") return null;
  return (
    <ProductVariantsBuilder
      variants={state.form.variants}
      lockedVariantCount={state.lockedVariantCount}
      onRemoveVariant={state.handleRemoveVariant}
      onVariantChange={state.handleVariantChange}
      priceStep={priceStep}
      canRemove={true}
      t={t}
    />
  );
}

export default ProductCreateForm;
