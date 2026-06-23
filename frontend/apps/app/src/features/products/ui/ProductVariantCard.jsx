import FormField from "@smarterp/ui/form-modal/FormField";
import MediaManager from "../../media/ui/MediaManager";

function ProductVariantCard({ canRemove, fields, index, locked, onChange, onRemove, t, variant }) {
  return (
    <article className={`product-variant-card ${locked ? "is-locked" : ""}`.trim()}>
      <header className="product-variant-card-head">
        <div>
          <h4 className="product-variant-card-title">{t("products.form.variantCardTitle", { index: index + 1 })}</h4>
          {locked && <p className="product-variant-card-note">{t("products.form.variantSavedHint")}</p>}
        </div>
        <button className="product-variant-remove" type="button" onClick={() => onRemove(variant.id)} disabled={!canRemove || locked}>{t("products.form.removeVariant")}</button>
      </header>
      <div className="product-variant-grid">
        {fields.map((field) => <FormField key={`${variant.id}-${field.name}`} field={field} value={variant[field.name]} onChange={(event) => onChange(variant.id, event)} t={t} />)}
      </div>
      {isPersistedVariant(variant.id) && (
        <div className="product-variant-media">
          <span className="product-variant-media-title">{t("products.form.sections.media")}</span>
          <MediaManager kind="variant" ownerId={variant.id} t={t} />
        </div>
      )}
    </article>
  );
}

function isPersistedVariant(id) {
  return Boolean(id) && !String(id).startsWith("variant-");
}

export default ProductVariantCard;
