export default function StorefrontHeader({ url, t }) {
  return (
    <header className="storefront-head">
      <h2>{t("storefront.title")}</h2>
      {url && (
        <p className="storefront-address">
          {t("storefront.shopAddress")}{" "}
          <a href={url} target="_blank" rel="noreferrer">
            {url.replace("https://", "")} ↗
          </a>
        </p>
      )}
    </header>
  );
}
