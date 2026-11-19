function ThemeCard({ theme, selected, onSelect, t }) {
  return (
    <button
      type="button"
      aria-pressed={selected}
      className={selected ? "theme-card is-selected" : "theme-card"}
      onClick={() => onSelect(theme.id)}
    >
      {theme.preview_url && <span className="theme-card-preview"><img src={theme.preview_url} alt="" loading="lazy" /></span>}
      <span className="theme-card-name">{theme.name}</span>
      <span className="theme-card-state">{selected ? t("storefront.selected") : t("storefront.use")}</span>
    </button>
  );
}

export default ThemeCard;
