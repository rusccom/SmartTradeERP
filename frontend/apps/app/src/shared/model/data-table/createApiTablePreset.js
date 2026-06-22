import { getJSON } from "../../api/http";
import { createTablePreset } from "./createTablePreset";

export function createApiTablePreset(config) {
  validateApiConfig(config);
  return createTablePreset({
    ...config,
    fetchPage: (params) => fetchApiPage(config, params),
  });
}

function validateApiConfig(config) {
  if (!config?.path) {
    throw new Error('API table preset requires "path"');
  }
}

async function fetchApiPage(config, params) {
  const response = await getJSON(config.path, params.query, params.signal);
  return {
    rows: mapRows(response.data, config),
    total: response.meta?.total || 0,
  };
}

function mapRows(data, config) {
  const rows = Array.isArray(data) ? data : [];
  return typeof config.mapRows === "function" ? config.mapRows(rows) : rows;
}
