import { useCallback, useEffect, useMemo, useState } from "react";

import { loadCurrencies } from "../api/loadCurrencies";
import { createCurrency } from "../api/createCurrency";
import { setBaseCurrency } from "../api/setBaseCurrency";
import CurrencyContext from "./currencyContext";
import { formatMoneyValue } from "./formatMoney";

function CurrencyProvider({ children }) {
  const [currencies, setCurrencies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const defaultCurrency = useMemo(() => readDefaultCurrency(currencies), [currencies]);
  const formatMoney = useCallback((value) => formatMoneyValue(value, defaultCurrency), [defaultCurrency]);
  const refresh = useCallback(() => reloadCurrencies(setCurrencies, setError, setLoading), []);
  const addCurrency = useCallback((payload) => createAndReload(payload, refresh, setError), [refresh]);
  const changeBase = useCallback((payload) => setBaseAndReload(payload, refresh, setError), [refresh]);
  useInitialLoad(setCurrencies, setError, setLoading);
  const value = useProviderValue({
    addCurrency, currencies, defaultCurrency, error, formatMoney,
    loading, refresh, setBaseCurrency: changeBase,
  });
  return <CurrencyContext.Provider value={value}>{children}</CurrencyContext.Provider>;
}

function useInitialLoad(setCurrencies, setError, setLoading) {
  useEffect(() => {
    const controller = new AbortController();
    reloadCurrencies(setCurrencies, setError, setLoading, controller.signal);
    return () => controller.abort();
  }, [setCurrencies, setError, setLoading]);
}

function useProviderValue(value) {
  return useMemo(() => value, [
    value.addCurrency, value.currencies, value.defaultCurrency, value.error,
    value.formatMoney, value.loading, value.refresh, value.setBaseCurrency,
  ]);
}

function readDefaultCurrency(currencies) {
  return currencies.find((item) => item.is_base) || currencies[0] || null;
}

async function reloadCurrencies(setCurrencies, setError, setLoading, signal) {
  setLoading(true);
  setError("");
  try {
    setCurrencies(await loadCurrencies(signal));
  } catch (error) {
    if (error.name !== "AbortError") setError(error.message);
  } finally {
    if (!signal?.aborted) setLoading(false);
  }
}

async function createAndReload(payload, refresh, setError) {
  setError("");
  try {
    await createCurrency(payload);
    await refresh();
  } catch (error) {
    setError(error.message);
    throw error;
  }
}

async function setBaseAndReload(payload, refresh, setError) {
  setError("");
  try {
    await setBaseCurrency(payload);
    await refresh();
  } catch (error) {
    setError(error.message);
    throw error;
  }
}

export default CurrencyProvider;
