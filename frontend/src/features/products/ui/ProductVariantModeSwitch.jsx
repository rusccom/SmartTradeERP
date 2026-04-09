function ProductVariantModeSwitch({ disabled, mode, onChange, t }) {
  return (
    <section className="product-mode-switch">
      <div className="product-mode-switch-head">
        <p className="product-mode-switch-label">{t("products.form.modeLabel")}</p>
        <p className="product-mode-switch-text">{t("products.form.modeHint")}</p>
      </div>
      <div className="product-mode-switch-actions">
        <ModeButton value="single" mode={mode} disabled={disabled} onChange={onChange} t={t} />
        <ModeButton value="multiple" mode={mode} disabled={disabled} onChange={onChange} t={t} />
      </div>
    </section>
  );
}

function ModeButton({ disabled, mode, onChange, t, value }) {
  const active = mode === value ? "is-active" : "";
  return <button className={`product-mode-btn ${active}`.trim()} type="button" onClick={() => onChange(value)} disabled={disabled}>{t(`products.form.mode.${value}`)}</button>;
}

export default ProductVariantModeSwitch;
