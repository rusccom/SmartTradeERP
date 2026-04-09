import { Navigate, useParams } from "react-router-dom";

import { buildLandingPath, DEFAULT_LOCALE, readSupportedLocales } from "../../../shared/i18n/localeConfig";
import { messagesByLocale } from "../../../shared/i18n/messages";
import LandingPage from "./LandingPage";

const supportedLocales = readSupportedLocales(messagesByLocale);

function LocalizedLandingPage() {
  const { locale } = useParams();
  const normalizedLocale = locale?.toLowerCase();
  if (!normalizedLocale || !supportedLocales.includes(normalizedLocale)) {
    return <Navigate to={buildLandingPath(DEFAULT_LOCALE)} replace />;
  }
  if (normalizedLocale === DEFAULT_LOCALE) {
    return <Navigate to={buildLandingPath(DEFAULT_LOCALE)} replace />;
  }
  return <LandingPage />;
}

export default LocalizedLandingPage;
