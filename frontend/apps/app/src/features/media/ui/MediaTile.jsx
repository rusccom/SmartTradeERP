function MediaTile({ item, busy, t, onPrimary, onRemove }) {
  return (
    <figure className={`media-tile ${item.is_primary ? "is-primary" : ""}`.trim()}>
      <img className="media-tile-img" src={item.thumb_url || item.url} alt={item.file_name || ""} loading="lazy" />
      {item.is_primary && <span className="media-tile-badge">{t("products.form.mediaCover")}</span>}
      <div className="media-tile-actions">
        {!item.is_primary && (
          <button type="button" className="media-tile-btn" disabled={busy} onClick={onPrimary}>
            {t("products.form.mediaMakeCover")}
          </button>
        )}
        <button type="button" className="media-tile-btn is-danger" disabled={busy} onClick={onRemove}>
          {t("products.form.mediaDelete")}
        </button>
      </div>
    </figure>
  );
}

export default MediaTile;
