import { useEffect, useState } from "react";

import { loadStorefrontSettings } from "../api/loadStorefrontSettings";
import { loadStorefrontThemes } from "../api/loadStorefrontThemes";
import { loadStorefrontPreview } from "../api/loadStorefrontPreview";
import { saveStorefrontDraft } from "../api/saveStorefrontDraft";
import { publishStorefront } from "../api/publishStorefront";

const INITIAL_META = { loading: true, busy: false, error: "", notice: "", previewUrl: "", shopUrl: "" };

export function useStorefront() {
  const [themes, setThemes] = useState([]);
  const [themeId, setThemeId] = useState("");
  const [meta, setMeta] = useState(INITIAL_META);
  useEffect(() => loadInitial({ setThemes, setThemeId, setMeta }), []);
  return buildApi({ themes, themeId, meta, setThemeId, setMeta });
}

function loadInitial(setters) {
  const controller = new AbortController();
  fetchInitial(controller.signal, setters);
  return () => controller.abort();
}

async function fetchInitial(signal, setters) {
  try {
    const [themes, settings] = await Promise.all([
      loadStorefrontThemes(signal),
      loadStorefrontSettings(signal),
    ]);
    setters.setThemes(themePreviews(themes, settings?.shop_url));
    applySettings(themes, settings, setters);
    setters.setMeta((meta) => ({ ...meta, loading: false }));
  } catch (error) {
    failInitial(error, setters.setMeta);
  }
}

function themePreviews(themes, shopUrl) {
  return themes.map(theme => ({
    ...theme,
    preview_url: theme.preview_url && shopUrl ? new URL(theme.preview_url, shopUrl).href : "",
  }));
}

function failInitial(error, setMeta) {
  if (error.name === "AbortError") {
    return;
  }
  setMeta((meta) => ({ ...meta, loading: false, error: error.message }));
}

function applySettings(themes, settings, setters) {
  const id = settings?.draft_theme_id || settings?.theme_id || themes[0]?.id || "";
  setters.setThemeId(id);
  setters.setMeta((meta) => ({ ...meta, shopUrl: settings?.shop_url || "" }));
}

function buildApi(state) {
  return {
    loading: state.meta.loading,
    busy: state.meta.busy,
    error: state.meta.error,
    notice: state.meta.notice,
    previewUrl: state.meta.previewUrl,
    shopUrl: state.meta.shopUrl,
    themes: state.themes,
    themeId: state.themeId,
    selectTheme: state.setThemeId,
    save: () => runAction(state, false),
    publish: () => runAction(state, true),
    preview: () => runPreview(state),
  };
}

async function runAction(state, alsoPublish) {
  startBusy(state.setMeta);
  try {
    await saveStorefrontDraft({ theme_id: state.themeId });
    if (alsoPublish) {
      await publishStorefront();
    }
    state.setMeta((meta) => ({ ...meta, busy: false, notice: alsoPublish ? "published" : "draftSaved" }));
  } catch (error) {
    state.setMeta((meta) => ({ ...meta, busy: false, error: error.message }));
  }
}

async function runPreview(state) {
  startBusy(state.setMeta);
  try {
    await saveStorefrontDraft({ theme_id: state.themeId });
    finishPreview(state.setMeta, await loadStorefrontPreview());
  } catch (error) {
    state.setMeta((meta) => ({ ...meta, busy: false, error: error.message }));
  }
}

function finishPreview(setMeta, data) {
  if (data?.url) {
    setMeta((meta) => ({ ...meta, busy: false, previewUrl: data.url }));
    return;
  }
  setMeta((meta) => ({ ...meta, busy: false, notice: "previewNoHost" }));
}

function startBusy(setMeta) {
  setMeta((meta) => ({ ...meta, busy: true, error: "", notice: "", previewUrl: "" }));
}
