import FormField from "../../../shared/ui/form-modal/FormField";

function ProductVariantCard({ canRemove, fields, index, locked, onChange, onRemove, t, variant }) {
  return (
    <article className={`product-variant-card ${locked ? "is-locked" : ""}`.trim()}>
      <header className="product-variant-card-head">
        <div>
          <h3 className="product-variant-card-title">{t("products.form.variantCardTitle", { index: index + 1 })}</h3>
          {locked && <p className="product-variant-card-note">{t("products.form.variantSavedHint")}</p>}
        </div>
        <button className="product-variant-remove" type="button" onClick={() => onRemove(variant.id)} disabled={!canRemove || locked}>{t("products.form.removeVariant")}</button>
      </header>
      <div className="product-variant-grid">
        {fields.map((field) => <FormField key={`${variant.id}-${field.name}`} field={field} value={variant[field.name]} onChange={(event) => onChange(variant.id, event)} t={t} />)}
      </div>
    </article>
  );
}

export default ProductVariantCard;
