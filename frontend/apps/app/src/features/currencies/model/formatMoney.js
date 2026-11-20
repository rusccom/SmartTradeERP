export function formatMoneyValue(value, currency, locale) {
  const number = Number(value);
  if (!Number.isFinite(number)) {
    return value === undefined || value === null ? "" : String(value);
  }
  if (!currency) {
    return formatDecimal(number, 2, locale);
  }
  const amount = formatDecimal(number, currency.decimal_places, locale);
  const marker = currency.display_symbol || currency.symbol || currency.code;
  return marker ? `${amount} ${marker}` : amount;
}

export function decimalStep(currency) {
  const places = Number(currency?.decimal_places ?? 2);
  if (!Number.isFinite(places) || places <= 0) {
    return "1";
  }
  return `0.${"0".repeat(Math.max(places - 1, 0))}1`;
}

function formatDecimal(number, places, locale) {
  const fractionDigits = normalizedPlaces(places);
  return number.toLocaleString(locale, {
    maximumFractionDigits: fractionDigits,
    minimumFractionDigits: fractionDigits,
  });
}

function normalizedPlaces(value) {
  const parsed = Number(value);
  if (!Number.isFinite(parsed)) {
    return 2;
  }
  return Math.min(Math.max(parsed, 0), 4);
}
