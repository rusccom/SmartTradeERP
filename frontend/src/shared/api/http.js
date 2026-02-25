const API_BASE_URL = normalizeBaseURL(import.meta.env.VITE_API_URL);

export async function postJSON(path, payload) {
  const response = await fetch(createURL(path), createPostOptions(payload));
  return parseEnvelope(response);
}

function createURL(path) {
  if (API_BASE_URL === "") {
    throw new Error("VITE_API_URL is required");
  }
  return `${API_BASE_URL}${path}`;
}

function createPostOptions(payload) {
  return {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  };
}

async function parseEnvelope(response) {
  const body = await parseBody(response);
  if (!response.ok || body.error) {
    throw new Error(getErrorMessage(body, response.status));
  }
  return body.data ?? null;
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
