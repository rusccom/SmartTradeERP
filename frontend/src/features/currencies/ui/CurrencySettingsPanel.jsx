import { Plus, RefreshCw } from "lucide-react";
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
          {state.canCreate && renderAddButton(t, state.openForm)}
          {renderRefreshButton(t, state.refresh, state.loading)}
        </div>
      </header>
      <CurrencySummary
        currency={state.defaultCurrency}
        emptyLabel={t("currencies.empty")}
        formatMoney={state.formatMoney}
        loadingLabel={state.loading ? t("currencies.loading") : ""}
      />
      {state.showForm && <CurrencyCreateForm labels={currencyFormLabels(t)} onSubmit={state.handleSubmit} options={state.options} />}
      {state.error && <p className="currency-error">{state.error}</p>}
    </section>
  );
}

function usePanelState() {
  const currencies = useCurrencies();
  const options = useCurrencyOptions();
  const [formOpen, setFormOpen] = useState(false);
  const canCreate = !currencies.loading && currencies.currencies.length === 0;
  return {
    ...currencies,
    canCreate,
    error: currencies.error || options.error,
    handleSubmit: (payload) => submitCurrency(payload, currencies.addCurrency, setFormOpen),
    openForm: () => setFormOpen(true),
    options: options.items,
    showForm: canCreate && options.items.length > 0 && (formOpen || !options.loading),
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

async function submitCurrency(payload, addCurrency, setFormOpen) {
  await addCurrency(payload);
  setFormOpen(false);
}

function renderAddButton(t, onClick) {
  return <button type="button" onClick={onClick}><Plus size={16} /> {t("currencies.addButton")}</button>;
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

export default CurrencySettingsPanel;
