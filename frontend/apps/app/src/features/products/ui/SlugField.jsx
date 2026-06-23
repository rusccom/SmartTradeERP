function SlugField({ value, onChange, t }) {
  return (
    <label className="shared-form-field">
      <span className="shared-form-field-label">{t("products.form.slug")}</span>
      <span className="product-slug-field">
        <span className="product-slug-prefix">products/</span>
        <input
          className="shared-form-field-input product-slug-input"
          type="text"
          name="slug"
          value={value}
          onChange={onChange}
          autoComplete="off"
        />
      </span>
    </label>
  );
}

export default SlugField;
