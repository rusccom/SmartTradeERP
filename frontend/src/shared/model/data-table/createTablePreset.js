import { FALLBACK_STATE } from "./tableState";

const DEFAULT_CAPABILITIES = { sorting: false, search: false };
const REQUIRED_FIELDS = ["id", "rowId", "columns", "fetchPage"];

export function createTablePreset(config) {
  validateConfig(config);
  const preset = {
    ...config,
    defaultState: mergeState(config.defaultState),
    capabilities: mergeCapabilities(config.capabilities),
  };
  return Object.freeze(preset);
}

function validateConfig(config) {
  REQUIRED_FIELDS.forEach((field) => {
    if (config?.[field] == null) {
      throw new Error(`Table preset requires "${field}"`);
    }
  });
  if (!Array.isArray(config.columns)) {
    throw new Error('Table preset field "columns" must be an array');
  }
}

function mergeState(defaultState) {
  const state = defaultState || {};
  return {
    ...FALLBACK_STATE,
    ...state,
    pagination: { ...FALLBACK_STATE.pagination, ...(state.pagination || {}) },
  };
}

function mergeCapabilities(capabilities) {
  return { ...DEFAULT_CAPABILITIES, ...(capabilities || {}) };
}
