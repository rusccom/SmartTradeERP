export const DEFAULT_LOCALE = "en";

export function readSupportedLocales(messagesByLocale) {
  return Object.keys(messagesByLocale);
}

export function isSupportedLocale(locale, supportedLocales) {
  return supportedLocales.includes(locale);
}

export function readRouteLocale(pathname, supportedLocales) {
  const match = pathname.match(/^\/([a-z]{2})(?:\/)?$/i);
  const locale = match?.[1]?.toLowerCase();
  return locale && isSupportedLocale(locale, supportedLocales) ? locale : null;
}

export function readLandingLocale(pathname, supportedLocales) {
  if (pathname === "/") {
    return DEFAULT_LOCALE;
  }
  return readRouteLocale(pathname, supportedLocales);
}

export function isLandingPath(pathname, supportedLocales) {
  return readLandingLocale(pathname, supportedLocales) !== null;
}

export function buildLandingPath(locale) {
  return locale === DEFAULT_LOCALE ? "/" : `/${locale}`;
}
