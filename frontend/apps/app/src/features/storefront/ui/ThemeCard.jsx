function ThemeCard({ theme, selected, onSelect, t }) {
  return (
    <button
      type="button"
      className={selected ? "theme-card is-selected" : "theme-card"}
      onClick={() => onSelect(theme.id)}
    >
      <span className="theme-card-preview" style={previewStyle(theme.tokens)} />
      <span className="theme-card-name">{theme.name}</span>
      <span className="theme-card-state">{selected ? t("storefront.selected") : t("storefront.use")}</span>
    </button>
  );
}

function previewStyle(tokens) {
  return {
    background: tokens?.["color-bg"] || "#ffffff",
    color: tokens?.["color-fg"] || "#111827",
    borderColor: tokens?.["color-border"] || "#e5e7eb",
  };
}

export default ThemeCard;
