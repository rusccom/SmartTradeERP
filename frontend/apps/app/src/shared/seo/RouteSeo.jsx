import { useEffect } from "react";
import { useLocation } from "react-router-dom";

import { buildLandingPath, DEFAULT_LOCALE, isLandingPath, readRouteLocale, readSupportedLocales } from "../i18n/localeConfig";
import { messagesByLocale } from "../i18n/messages";
import { useI18n } from "../i18n/useI18n";

const supportedLocales = readSupportedLocales(messagesByLocale);

function RouteSeo() {
  const { pathname } = useLocation();
  const { t } = useI18n();

  useEffect(() => {
    const meta = readMeta(pathname, t);
    document.title = meta.title;
    upsertMeta("description", meta.description);
    upsertMeta("robots", meta.robots);
    upsertMeta("googlebot", meta.robots);
    syncCanonical(pathname);
    syncAlternates(pathname);
  }, [pathname, t]);

  return null;
}

function readMeta(pathname, t) {
  if (isLandingPath(pathname, supportedLocales)) {
    return {
      title: t("public.seo.landing.title"),
      description: t("public.seo.landing.description"),
      robots: "index,follow",
    };
  }
  return {
    title: t("workspace.private.title"),
    description: t("workspace.private.description"),
    robots: "noindex,nofollow",
  };
}

function upsertMeta(name, content) {
  const selector = `meta[name="${name}"]`;
  const element = document.head.querySelector(selector) ?? createMeta(name);
  element.setAttribute("content", content);
}

function createMeta(name) {
  const element = document.createElement("meta");
  element.setAttribute("name", name);
  document.head.appendChild(element);
  return element;
}

function syncCanonical(pathname) {
  const existing = document.head.querySelector('link[rel="canonical"]');
  if (!isLandingPath(pathname, supportedLocales)) {
    existing?.remove();
    return;
  }
  const locale = readRouteLocale(pathname, supportedLocales) || DEFAULT_LOCALE;
  const canonical = existing ?? createCanonical();
  canonical.setAttribute("href", `${window.location.origin}${buildLandingPath(locale)}`);
}

function createCanonical() {
  const element = document.createElement("link");
  element.setAttribute("rel", "canonical");
  document.head.appendChild(element);
  return element;
}

function syncAlternates(pathname) {
  clearAlternates();
  if (!isLandingPath(pathname, supportedLocales)) {
    return;
  }
  supportedLocales.forEach((locale) => createAlternate(locale, buildLandingPath(locale)));
  createAlternate("x-default", buildLandingPath(DEFAULT_LOCALE));
}

function createAlternate(hreflang, path) {
  const element = document.createElement("link");
  element.setAttribute("rel", "alternate");
  element.setAttribute("hreflang", hreflang);
  element.setAttribute("href", `${window.location.origin}${path}`);
  element.setAttribute("data-i18n-alt", "true");
  document.head.appendChild(element);
}

function clearAlternates() {
  document.head.querySelectorAll('link[data-i18n-alt="true"]').forEach((element) => element.remove());
}

export default RouteSeo;
