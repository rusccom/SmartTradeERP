import { getAdminToken, getClientToken } from "../auth/session";
import { assertPublicApiPath, isAdminApiPath, isClientApiPath } from "./publicApi";

const API_BASE_URL = normalizeBaseURL(import.meta.env.VITE_API_URL);

export async function postJSON(path, payload) {
  const response = await fetch(createURL(path), createPostOptions(path, payload));
  return parseEnvelope(response);
}

export async function putJSON(path, payload) {
  const response = await fetch(createURL(path), createWriteOptions(path, payload, "PUT"));
  return parseEnvelope(response);
}

export async function getJSON(path, params, signal) {
  const response = await fetch(buildURL(path, params), createGetOptions(path, signal));
  return parseEnvelopeWithMeta(response);
}

function createURL(path) {
  assertPublicApiPath(path);
  if (API_BASE_URL === "") {
    throw new Error("VITE_API_URL is required");
  }
  return `${API_BASE_URL}${path}`;
}

function resolveToken(path) {
  if (isAdminApiPath(path)) {
    return getAdminToken() || "";
  }
  if (isClientApiPath(path)) {
    return getClientToken() || "";
  }
  return getClientToken() || getAdminToken() || "";
}

function authHeaders(path) {
  const token = resolveToken(path);
  return token ? { Authorization: `Bearer ${token}` } : {};
}

function createPostOptions(path, payload) {
  return createWriteOptions(path, payload, "POST");
}

function createWriteOptions(path, payload, method) {
  return {
    method,
    headers: { "Content-Type": "application/json", ...authHeaders(path) },
    body: JSON.stringify(payload),
  };
}

function buildURL(path, params) {
  const url = new URL(createURL(path));
  Object.entries(params || {}).forEach(([key, value]) => {
    if (value === undefined || value === null || value === "") {
      return;
    }
    url.searchParams.set(key, String(value));
  });
  return url.toString();
}

function createGetOptions(path, signal) {
  return { method: "GET", headers: authHeaders(path), signal };
}

async function parseEnvelope(response) {
  const body = await parseBody(response);
  if (!response.ok || body.error) {
    throw new Error(getErrorMessage(body, response.status));
  }
  return body.data ?? null;
}

async function parseEnvelopeWithMeta(response) {
  const body = await parseBody(response);
  if (!response.ok || body.error) {
    throw new Error(getErrorMessage(body, response.status));
  }
  return { data: body.data ?? null, meta: body.meta ?? null };
}

async function parseBody(response) {
  const raw = await response.text();
  if (!raw) {
    return {};
  }
  try {
    return JSON.parse(raw);
  } catch {
    return {};
  }
}

function getErrorMessage(body, statusCode) {
  const message = body?.error?.message;
  if (message) {
    return message;
  }
  return `Request failed with status ${statusCode}`;
}

function normalizeBaseURL(value) {
  if (typeof value !== "string") {
    return "";
  }
  const trimmed = value.trim();
  if (trimmed === "") {
    return "";
  }
  return trimmed.replace(/\/+$/, "");
}
