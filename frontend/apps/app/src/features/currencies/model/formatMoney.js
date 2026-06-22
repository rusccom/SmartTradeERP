export function formatMoneyValue(value, currency) {
  const number = Number(value);
  if (!Number.isFinite(number)) {
    return value === undefined || value === null ? "" : String(value);
  }
  if (!currency) {
    return formatDecimal(number, 2);
  }
  const amount = formatDecimal(number, currency.decimal_places);
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

function formatDecimal(number, places) {
  const fractionDigits = normalizedPlaces(places);
  return number.toLocaleString(undefined, {
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
