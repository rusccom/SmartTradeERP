import { useLocation, useNavigate } from "react-router-dom";

import { buildLandingPath, isLandingPath } from "../i18n/localeConfig";
import { useI18n } from "../i18n/useI18n";

function LocaleSwitcher() {
  const { availableLocales, locale, setLocale, t } = useI18n();
  const location = useLocation();
  const navigate = useNavigate();
  if (availableLocales.length < 2) {
    return null;
  }

  function handleChange(nextLocale) {
    setLocale(nextLocale);
    if (isLandingPath(location.pathname, availableLocales)) {
      navigate(buildLandingPath(nextLocale));
    }
  }

  return (
    <label className="locale-switcher">
      <span className="locale-label">{t("locale.label")}</span>
      <select className="locale-select" value={locale} onChange={(event) => handleChange(event.target.value)}>
        {availableLocales.map((item) => (
          <option key={item} value={item}>
            {t(`locale.${item}`)}
          </option>
        ))}
      </select>
    </label>
  );
}

export default LocaleSwitcher;
