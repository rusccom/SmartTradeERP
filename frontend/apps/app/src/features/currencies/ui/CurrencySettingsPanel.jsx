import { RefreshCw } from "lucide-react";
import { useEffect, useState } from "react";

import { useI18n } from "../../../shared/i18n/useI18n";
import { loadCurrencyOptions } from "../api/loadCurrencyOptions";
import { useCurrencies } from "../model/useCurrencies";
import CurrencyCreateForm from "./CurrencyCreateForm";
import CurrencySummary from "./CurrencySummary";
import "./currencies.css";

function CurrencySettingsPanel() {
  const { t } = useI18n();
  const state = usePanelState();
  return (
    <section className="currency-settings">
      <header className="currency-settings-head">
        <h2>{t("currencies.title")}</h2>
        <div className="currency-settings-actions">
          {renderRefreshButton(t, state.refresh, state.loading)}
        </div>
      </header>
      <CurrencySummary
        currency={state.defaultCurrency}
        emptyLabel={t("currencies.empty")}
        formatMoney={state.formatMoney}
        loadingLabel={state.loading ? t("currencies.loading") : ""}
      />
      {state.showForm && <CurrencyCreateForm {...currencyFormProps(t, state)} />}
      {state.error && <p className="currency-error">{state.error}</p>}
    </section>
  );
}

function usePanelState() {
  const currencies = useCurrencies();
  const options = useCurrencyOptions();
  return {
    ...currencies,
    error: currencies.error || options.error,
    handleSubmit: (payload) => currencies.setBaseCurrency(payload),
    options: options.items,
    showForm: !currencies.loading && options.items.length > 0,
  };
}

function useCurrencyOptions() {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  useLoadOptions(setItems, setLoading, setError);
  return { error, items, loading };
}

function useLoadOptions(setItems, setLoading, setError) {
  useEffect(() => {
    const controller = new AbortController();
    loadOptions(controller.signal, setItems, setLoading, setError);
    return () => controller.abort();
  }, [setError, setItems, setLoading]);
}

async function loadOptions(signal, setItems, setLoading, setError) {
  try {
    setItems(await loadCurrencyOptions(signal));
  } catch (error) {
    if (error.name !== "AbortError") setError(error.message);
  } finally {
    if (!signal.aborted) setLoading(false);
  }
}

function renderRefreshButton(t, refresh, loading) {
  return <button type="button" onClick={refresh} disabled={loading} title={t("currencies.reload")}><RefreshCw size={16} /></button>;
}

function currencyFormLabels(t) {
  return {
    currency: t("currencies.select"),
    save: t("currencies.save"),
    saving: t("currencies.saving"),
    symbol: t("currencies.symbol"),
  };
}

function currencyFormProps(t, state) {
  return {
    initialCurrencyID: state.defaultCurrency?.currency_id || "",
    initialSymbol: state.defaultCurrency?.display_symbol || "",
    labels: currencyFormLabels(t),
    onSubmit: state.handleSubmit,
    options: state.options,
  };
}

export default CurrencySettingsPanel;
