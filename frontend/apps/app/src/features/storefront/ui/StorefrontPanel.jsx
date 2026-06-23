import { useI18n } from "@smarterp/i18n/useI18n";

import { useStorefront } from "../model/useStorefront";
import ThemeGallery from "./ThemeGallery";
import TokenEditor from "./TokenEditor";
import SectionsEditor from "./SectionsEditor";
import "./storefront.css";

function StorefrontPanel() {
  const { t } = useI18n();
  const store = useStorefront();
  if (store.loading) {
    return <p className="storefront-loading">{t("storefront.loading")}</p>;
  }
  return (
    <section className="storefront">
      <header className="storefront-head">
        <h2>{t("storefront.title")}</h2>
        <p className="storefront-sub">{t("storefront.subtitle")}</p>
      </header>
      <ThemeGallery themes={store.themes} themeId={store.themeId} onSelect={store.selectTheme} t={t} />
      <TokenEditor
        tokenKeys={store.tokenKeys}
        tokens={store.tokens}
        onChange={store.setToken}
        onReset={store.resetTokens}
        t={t}
      />
      <SectionsEditor sections={store.sections} onToggle={store.toggleSection} onMove={store.moveSection} t={t} />
      <StorefrontActions store={store} t={t} />
    </section>
  );
}

function StorefrontActions({ store, t }) {
  return (
    <div className="storefront-actions">
      <button type="button" className="storefront-btn" onClick={store.preview} disabled={store.busy}>
        {store.busy ? t("storefront.busy") : t("storefront.preview")}
      </button>
      <button type="button" className="storefront-btn" onClick={store.save} disabled={store.busy}>
        {store.busy ? t("storefront.busy") : t("storefront.save")}
      </button>
      <button type="button" className="storefront-btn is-primary" onClick={store.publish} disabled={store.busy}>
        {store.busy ? t("storefront.busy") : t("storefront.publish")}
      </button>
      <StorefrontFeedback store={store} t={t} />
    </div>
  );
}

function StorefrontFeedback({ store, t }) {
  if (store.error) {
    return <span className="storefront-error">{store.error}</span>;
  }
  if (store.previewUrl) {
    return (
      <a className="storefront-preview-link" href={store.previewUrl} target="_blank" rel="noreferrer">
        {t("storefront.openPreview")}
      </a>
    );
  }
  if (store.notice) {
    return <span className="storefront-notice">{t(`storefront.${store.notice}`)}</span>;
  }
  return null;
}

export default StorefrontPanel;
