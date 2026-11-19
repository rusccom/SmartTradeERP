import StorefrontFeedback from "./StorefrontFeedback";

export default function StorefrontActions({ store, t }) {
  const disabled = store.busy || !store.themeId;
  return (
    <div className="storefront-actions">
      <button type="button" className="storefront-btn" onClick={store.preview} disabled={disabled}>
        {store.busy ? t("storefront.busy") : t("storefront.preview")}
      </button>
      <button type="button" className="storefront-btn" onClick={store.save} disabled={disabled}>
        {store.busy ? t("storefront.busy") : t("storefront.save")}
      </button>
      <button type="button" className="storefront-btn is-primary" onClick={store.publish} disabled={disabled}>
        {store.busy ? t("storefront.busy") : t("storefront.publish")}
      </button>
      <StorefrontFeedback store={store} t={t} />
    </div>
  );
}
