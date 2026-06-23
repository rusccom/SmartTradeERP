import ThemeCard from "./ThemeCard";

function ThemeGallery({ themes, themeId, onSelect, t }) {
  return (
    <div className="storefront-themes">
      {themes.map((theme) => (
        <ThemeCard
          key={theme.id}
          theme={theme}
          selected={theme.id === themeId}
          onSelect={onSelect}
          t={t}
        />
      ))}
    </div>
  );
}

export default ThemeGallery;
