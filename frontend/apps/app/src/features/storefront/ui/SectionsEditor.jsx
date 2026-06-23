function SectionsEditor({ sections, onToggle, onMove, t }) {
  if (!sections || sections.length === 0) {
    return null;
  }
  return (
    <div className="storefront-sections">
      <h3>{t("storefront.sections")}</h3>
      <ul className="storefront-section-list">
        {sections.map((section, index) => (
          <SectionRow
            key={section.key}
            section={section}
            index={index}
            total={sections.length}
            onToggle={onToggle}
            onMove={onMove}
            t={t}
          />
        ))}
      </ul>
    </div>
  );
}

function SectionRow({ section, index, total, onToggle, onMove, t }) {
  return (
    <li className="storefront-section">
      <label className="storefront-section-label">
        <input type="checkbox" checked={section.enabled} onChange={() => onToggle(section.key)} />
        {t(`storefront.section.${section.key}`)}
      </label>
      <span className="storefront-section-move">
        <button type="button" onClick={() => onMove(index, -1)} disabled={index === 0} aria-label={t("storefront.moveUp")}>
          ↑
        </button>
        <button type="button" onClick={() => onMove(index, 1)} disabled={index === total - 1} aria-label={t("storefront.moveDown")}>
          ↓
        </button>
      </span>
    </li>
  );
}

export default SectionsEditor;
