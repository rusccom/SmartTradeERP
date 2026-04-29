import { useEffect, useState } from "react";

import { loadBundleDetails } from "../api/loadBundleDetails";
import { loadBundles } from "../api/loadBundles";
import { loadComponentOptions } from "../api/loadComponentOptions";
import { saveBundleComponents } from "../api/saveBundleComponents";
import {
  addComponentRow,
  createInitialRows,
  patchComponentRow,
  removeComponentRow,
  toComponentPayload,
} from "../model/bundleComponents";
import BundleComponentsEditor from "../ui/BundleComponentsEditor";
import BundleList from "../ui/BundleList";
import "../ui/bundles.css";

const OPTION_PAGE_SIZE = 20;

function BundlesPage() {
  const state = useBundlePageState();
  return (
    <main className="bundles-page">
      <BundleList bundles={state.bundles} loading={state.loading} selectedID={state.selectedID} onSelect={state.selectBundle} onReload={state.reload} />
      <BundleComponentsEditor {...editorProps(state)} />
    </main>
  );
}

function editorProps(state) {
  return {
    bundle: state.bundle,
    canLoadMoreOptions: state.canLoadMoreOptions,
    componentSearch: state.componentSearch,
    error: state.error,
    loading: state.loading,
    onAdd: state.addRow,
    onChange: state.changeRow,
    onLoadMoreOptions: state.loadMoreOptions,
    onRemove: state.removeRow,
    onSave: state.saveRows,
    onSearchComponents: state.searchComponents,
    options: state.options,
    optionsLoading: state.optionsLoading,
    rows: state.rows,
    saving: state.saving,
  };
}

function useBundlePageState() {
  const [bundles, setBundles] = useState([]);
  const [bundle, setBundle] = useState(null);
  const [options, setOptions] = useState([]);
  const [optionsMeta, setOptionsMeta] = useState(null);
  const [rows, setRows] = useState([]);
  const [selectedID, setSelectedID] = useState("");
  const [componentSearch, setComponentSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [optionsLoading, setOptionsLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  useLoadBundles({ setBundles, setError, setLoading, setSelectedID });
  useLoadOptions({ componentSearch, setError, setOptions, setOptionsLoading, setOptionsMeta });
  useLoadSelectedBundle({ selectedID, setBundle, setError, setRows });
  return createState({
    bundles, bundle, componentSearch, error, loading, options, optionsLoading,
    optionsMeta, rows, saving, selectedID, setBundles, setComponentSearch,
    setError, setLoading, setOptions, setOptionsLoading, setOptionsMeta,
    setRows, setSaving, setSelectedID,
  });
}

function useLoadBundles(state) {
  useEffect(() => {
    const controller = new AbortController();
    loadInitialData(controller.signal, state);
    return () => controller.abort();
  }, [state.setBundles, state.setError, state.setLoading, state.setSelectedID]);
}

async function loadInitialData(signal, state) {
  state.setLoading(true);
  state.setError("");
  try {
    const data = await loadBundles(signal);
    state.setBundles(data.bundles);
    state.setSelectedID((current) => nextSelectedID(current, data.bundles));
  } catch (error) {
    if (error.name !== "AbortError") state.setError(error.message);
  } finally {
    state.setLoading(false);
  }
}

function nextSelectedID(current, bundles) {
  if (bundles.some((item) => item.variant_id === current)) {
    return current;
  }
  return bundles[0]?.variant_id || "";
}

function useLoadOptions(state) {
  useEffect(() => {
    const controller = new AbortController();
    loadOptionsPage(controller.signal, state, 1, false);
    return () => controller.abort();
  }, [state.componentSearch, state.setError, state.setOptions, state.setOptionsMeta, state.setOptionsLoading]);
}

async function loadOptionsPage(signal, state, page, append) {
  state.setOptionsLoading(true);
  try {
    const query = { page, perPage: OPTION_PAGE_SIZE, search: state.componentSearch, signal };
    const data = await loadComponentOptions(query);
    state.setOptions((options) => (append ? mergeOptions(options, data.options) : data.options));
    state.setOptionsMeta(data.meta || fallbackMeta(page, data.options));
  } catch (error) {
    if (error.name !== "AbortError") state.setError(error.message);
  } finally {
    state.setOptionsLoading(false);
  }
}

function useLoadSelectedBundle(state) {
  useEffect(() => {
    if (!state.selectedID) {
      state.setBundle(null);
      state.setRows([]);
      return;
    }
    const controller = new AbortController();
    loadSelectedBundle(controller.signal, state);
    return () => controller.abort();
  }, [state.selectedID, state.setBundle, state.setError, state.setRows]);
}

async function loadSelectedBundle(signal, state) {
  state.setError("");
  try {
    const data = await loadBundleDetails(state.selectedID, signal);
    state.setBundle(data);
    state.setRows(createInitialRows(data?.components || []));
  } catch (error) {
    if (error.name !== "AbortError") state.setError(error.message);
  }
}

function createState(values) {
  return {
    ...values,
    addRow: () => values.setRows((rows) => addComponentRow(rows, values.options)),
    canLoadMoreOptions: hasMoreOptions(values.options, values.optionsMeta),
    changeRow: (id, event) => values.setRows((rows) => patchComponentRow(rows, id, event, values.options)),
    loadMoreOptions: () => loadMoreOptions(values),
    reload: () => loadInitialData(new AbortController().signal, values),
    removeRow: (id) => values.setRows((rows) => removeComponentRow(rows, id)),
    saveRows: () => saveRows(values),
    searchComponents: values.setComponentSearch,
    selectBundle: values.setSelectedID,
  };
}

function loadMoreOptions(state) {
  if (!hasMoreOptions(state.options, state.optionsMeta)) return;
  const page = (state.optionsMeta?.page || 1) + 1;
  loadOptionsPage(undefined, state, page, true);
}

function hasMoreOptions(options, meta) {
  return Boolean(meta && options.length < (meta.total || 0));
}

function mergeOptions(current, next) {
  const index = new Map(current.map((item) => [item.id, item]));
  next.forEach((item) => index.set(item.id, item));
  return Array.from(index.values());
}

function fallbackMeta(page, options) {
  return { page, per_page: OPTION_PAGE_SIZE, total: options.length };
}

async function saveRows(state) {
  if (!state.selectedID) return;
  state.setSaving(true);
  state.setError("");
  try {
    await saveBundleComponents(state.selectedID, toComponentPayload(state.rows));
    const data = await loadBundleDetails(state.selectedID);
    state.setBundle(data);
    state.setRows(createInitialRows(data?.components || []));
  } catch (error) {
    state.setError(error.message);
  } finally {
    state.setSaving(false);
  }
}

export default BundlesPage;
