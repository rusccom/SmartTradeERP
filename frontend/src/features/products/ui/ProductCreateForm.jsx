import ProductBasicSection from "./ProductBasicSection";
import ProductCreateActions from "./ProductCreateActions";
import ProductInventorySection from "./ProductInventorySection";
import ProductMediaSection from "./ProductMediaSection";
import ProductPriceSection from "./ProductPriceSection";
import ProductVariantsSection from "./ProductVariantsSection";
import ProductVariantsBuilder from "./ProductVariantsBuilder";

function ProductCreateForm({ state, t }) {
  const hasVariants = state.form.variantMode === "multiple";
  return (
    <form className="product-create-form" onSubmit={state.handleSubmit}>
      <ProductBasicSection form={state.form} onChange={state.handleChange} t={t} />
      <ProductMediaSection t={t} />
      <ProductPriceSection form={state.form} hasVariants={hasVariants} onChange={state.handleChange} t={t} />
      {!hasVariants && <ProductInventorySection form={state.form} onChange={state.handleChange} t={t} />}
      <ProductVariantsSection hasVariants={hasVariants} onAddVariant={state.handleAddVariant} disabled={state.hasPendingProduct} t={t}>
        {readBottomSlot(state, t)}
      </ProductVariantsSection>
      {state.error && <p className="product-create-error">{state.error}</p>}
      <ProductCreateActions isSaving={state.isSaving} onClose={state.handleClose} t={t} />
    </form>
  );
}

function readBottomSlot(state, t) {
  if (state.form.variantMode !== "multiple") return null;
  return (
    <ProductVariantsBuilder
      variants={state.form.variants}
      lockedVariantCount={state.lockedVariantCount}
      onRemoveVariant={state.handleRemoveVariant}
      onVariantChange={state.handleVariantChange}
      canRemove={true}
      t={t}
    />
  );
}

export default ProductCreateForm;
