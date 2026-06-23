import { useEffect, useState } from "react";

import { loadStorefrontSettings } from "../api/loadStorefrontSettings";
import { loadStorefrontThemes } from "../api/loadStorefrontThemes";
import { loadStorefrontPreview } from "../api/loadStorefrontPreview";
import { saveStorefrontDraft } from "../api/saveStorefrontDraft";
import { publishStorefront } from "../api/publishStorefront";

const INITIAL_META = { loading: true, busy: false, error: "", notice: "", previewUrl: "" };

export function useStorefront() {
  const [themes, setThemes] = useState([]);
  const [themeId, setThemeId] = useState("");
  const [tokens, setTokens] = useState({});
  const [sections, setSections] = useState([]);
  const [meta, setMeta] = useState(INITIAL_META);
  useEffect(() => loadInitial({ setThemes, setThemeId, setTokens, setSections, setMeta }), []);
  return buildApi({ themes, themeId, tokens, sections, meta, setThemeId, setTokens, setSections, setMeta });
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
    setters.setThemes(themes);
    applySettings(themes, settings, setters);
    setters.setMeta((meta) => ({ ...meta, loading: false }));
  } catch (error) {
    failInitial(error, setters.setMeta);
  }
}

function failInitial(error, setMeta) {
  if (error.name === "AbortError") {
    return;
  }
  setMeta((meta) => ({ ...meta, loading: false, error: error.message }));
}

function applySettings(themes, settings, setters) {
  const id = settings?.draft_theme_id || settings?.theme_id || themes[0]?.id || "classic";
  setters.setThemeId(id);
  setters.setTokens(effectiveTokens(themes, id, settings?.draft_tokens));
  setters.setSections(settings?.draft_sections || []);
}

function effectiveTokens(themes, id, overrides) {
  return { ...themeDefaults(themes, id), ...(overrides || {}) };
}

function themeDefaults(themes, id) {
  const theme = themes.find((item) => item.id === id);
  return theme?.tokens ? { ...theme.tokens } : {};
}

function buildApi(state) {
  return {
    loading: state.meta.loading,
    busy: state.meta.busy,
    error: state.meta.error,
    notice: state.meta.notice,
    previewUrl: state.meta.previewUrl,
    themes: state.themes,
    themeId: state.themeId,
    tokens: state.tokens,
    tokenKeys: Object.keys(themeDefaults(state.themes, state.themeId)),
    sections: state.sections,
    selectTheme: (id) => selectTheme(state, id),
    setToken: (key, value) => state.setTokens((current) => ({ ...current, [key]: value })),
    resetTokens: () => state.setTokens(themeDefaults(state.themes, state.themeId)),
    toggleSection: (key) => toggleSection(state.setSections, key),
    moveSection: (index, direction) => moveSection(state.setSections, index, direction),
    save: () => runAction(state, false),
    publish: () => runAction(state, true),
    preview: () => runPreview(state),
  };
}

function selectTheme(state, id) {
  state.setThemeId(id);
  state.setTokens(themeDefaults(state.themes, id));
}

function toggleSection(setSections, key) {
  setSections((list) => list.map((item) => (item.key === key ? { ...item, enabled: !item.enabled } : item)));
}

function moveSection(setSections, index, direction) {
  setSections((list) => reorder(list, index, index + direction));
}

function reorder(list, from, to) {
  if (to < 0 || to >= list.length) {
    return list;
  }
  const copy = [...list];
  const [item] = copy.splice(from, 1);
  copy.splice(to, 0, item);
  return copy;
}

function draftPayload(state) {
  return { theme_id: state.themeId, tokens: state.tokens, sections: state.sections };
}

async function runAction(state, alsoPublish) {
  startBusy(state.setMeta);
  try {
    await saveStorefrontDraft(draftPayload(state));
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
    await saveStorefrontDraft(draftPayload(state));
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
