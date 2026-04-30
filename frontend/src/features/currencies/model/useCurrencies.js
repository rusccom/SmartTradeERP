import { useContext } from "react";

import CurrencyContext from "./currencyContext";

export function useCurrencies() {
  return useContext(CurrencyContext);
}
