import { useI18n } from "@smarterp/i18n/useI18n";

import { useStorefront } from "../model/useStorefront";
import StorefrontActions from "./StorefrontActions";
import StorefrontHeader from "./StorefrontHeader";
import ThemeGallery from "./ThemeGallery";
import "./storefront.css";

export default function StorefrontPanel() {
  const { t } = useI18n();
  const store = useStorefront();
  if (store.loading) {
    return <p className="storefront-loading">{t("storefront.loading")}</p>;
  }
  return (
    <section className="storefront">
      <StorefrontHeader url={store.shopUrl} t={t} />
      <ThemeGallery themes={store.themes} themeId={store.themeId} onSelect={store.selectTheme} t={t} />
      <StorefrontActions store={store} t={t} />
    </section>
  );
}
