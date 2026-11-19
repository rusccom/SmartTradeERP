export default function StorefrontFeedback({ store, t }) {
  if (store.error) {
    return <span role="alert" className="storefront-error">{store.error}</span>;
  }
  if (store.previewUrl) {
    return (
      <a className="storefront-preview-link" href={store.previewUrl} target="_blank" rel="noreferrer">
        {t("storefront.openPreview")}
      </a>
    );
  }
  if (store.notice) {
    return <span role="status" className="storefront-notice">{t(`storefront.${store.notice}`)}</span>;
  }
  return null;
}
