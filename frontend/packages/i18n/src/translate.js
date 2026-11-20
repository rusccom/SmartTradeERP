import { DEFAULT_LOCALE } from "./localeConfig";
import { messagesByLocale } from "./messages";

export function translate(key, values, locale = globalThis.document?.documentElement.lang) {
  const template = messagesByLocale[locale]?.[key] || messagesByLocale[DEFAULT_LOCALE][key] || key;
  if (!values) return template;
  return template.replace(/\{(\w+)\}/g, (placeholder, name) =>
    Object.hasOwn(values, name) ? String(values[name]) : placeholder);
}
