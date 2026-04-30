import { createContext } from "react";

export const defaultCurrencyState = Object.freeze({
  addCurrency: async () => null,
  currencies: [],
  defaultCurrency: null,
  error: "",
  formatMoney: formatRawMoney,
  loading: false,
  refresh: async () => null,
});

function formatRawMoney(value) {
  if (value === undefined || value === null || value === "") return "";
  return String(value);
}

const CurrencyContext = createContext(defaultCurrencyState);

export default CurrencyContext;
