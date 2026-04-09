import { useEffect, useMemo, useState } from "react";
import { useLocation } from "react-router-dom";

import { DEFAULT_LOCALE, readLandingLocale, readSupportedLocales } from "./localeConfig";
import { I18nContext } from "./i18nContext";
import { messagesByLocale } from "./messages";

const STORAGE_KEY = "smarttrade.locale";
const SUPPORTED_LOCALES = readSupportedLocales(messagesByLocale);

function I18nProvider({ children }) {
  const location = useLocation();
  const [storedLocale, setStoredLocale] = useState(readInitialLocale);
  const landingLocale = readLandingLocale(location.pathname, SUPPORTED_LOCALES);
  const locale = landingLocale || storedLocale;

  useEffect(() => syncLocale(locale), [locale]);

  const value = useMemo(
    () => createValue({ locale, setStoredLocale }),
    [locale],
  );

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
}

function readInitialLocale() {
  return readStoredLocale() || readBrowserLocale() || DEFAULT_LOCALE;
}

function readStoredLocale() {
  if (typeof window === "undefined") {
    return null;
  }
  return normalizeLocale(window.localStorage.getItem(STORAGE_KEY));
}

function readBrowserLocale() {
  if (typeof navigator === "undefined") {
    return null;
  }
  return normalizeLocale(navigator.language);
}

function normalizeLocale(value) {
  const locale = value?.split("-")[0]?.toLowerCase();
  return locale && SUPPORTED_LOCALES.includes(locale) ? locale : null;
}

function syncLocale(locale) {
  if (typeof document !== "undefined") {
    document.documentElement.lang = locale;
  }
  if (typeof window !== "undefined") {
    window.localStorage.setItem(STORAGE_KEY, locale);
  }
}

function createValue({ locale, setStoredLocale }) {
  const current = messagesByLocale[locale] || messagesByLocale[DEFAULT_LOCALE];
  const fallback = messagesByLocale[DEFAULT_LOCALE];
  return {
    availableLocales: SUPPORTED_LOCALES,
    locale,
    setLocale: (nextLocale) => handleLocaleChange(nextLocale, setStoredLocale),
    t: (key, values) => formatMessage(current[key] || fallback[key] || key, values),
  };
}

function handleLocaleChange(nextLocale, setStoredLocale) {
  const locale = normalizeLocale(nextLocale);
  if (!locale) {
    return;
  }
  setStoredLocale(locale);
}

function formatMessage(template, values) {
  if (!values) {
    return template;
  }
  return Object.entries(values).reduce(
    (result, [key, value]) => result.replaceAll(`{${key}}`, String(value)),
    template,
  );
}

export default I18nProvider;
