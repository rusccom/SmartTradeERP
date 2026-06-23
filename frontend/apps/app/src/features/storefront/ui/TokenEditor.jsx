function TokenEditor({ tokenKeys, tokens, onChange, onReset, t }) {
  if (tokenKeys.length === 0) {
    return null;
  }
  return (
    <div className="storefront-tokens">
      <div className="storefront-tokens-head">
        <h3>{t("storefront.customize")}</h3>
        <button type="button" className="storefront-link" onClick={onReset}>
          {t("storefront.reset")}
        </button>
      </div>
      <div className="storefront-token-grid">
        {tokenKeys.map((key) => (
          <TokenField key={key} name={key} value={tokens[key] || ""} onChange={onChange} />
        ))}
      </div>
    </div>
  );
}

function TokenField({ name, value, onChange }) {
  return (
    <label className="storefront-token">
      <span className="storefront-token-label">{name}</span>
      <input
        type={name.startsWith("color-") ? "color" : "text"}
        value={value}
        onChange={(event) => onChange(name, event.target.value)}
      />
    </label>
  );
}

export default TokenEditor;
