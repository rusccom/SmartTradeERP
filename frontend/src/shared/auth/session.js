const ADMIN_TOKEN_KEY = "smarterp_admin_access_token";
const CLIENT_TOKEN_KEY = "smarterp_client_access_token";

export function setAdminToken(token) {
  saveToken(ADMIN_TOKEN_KEY, token);
}

export function setClientToken(token) {
  saveToken(CLIENT_TOKEN_KEY, token);
}

export function clearAdminToken() {
  clearToken(ADMIN_TOKEN_KEY);
}

export function clearClientToken() {
  clearToken(CLIENT_TOKEN_KEY);
}

export function hasAdminSession() {
  return hasToken(ADMIN_TOKEN_KEY);
}

export function hasClientSession() {
  return hasToken(CLIENT_TOKEN_KEY);
}

export function hasAnySession() {
  return hasAdminSession() || hasClientSession();
}

export function getAdminToken() {
  return getToken(ADMIN_TOKEN_KEY);
}

export function getClientToken() {
  return getToken(CLIENT_TOKEN_KEY);
}

export function resolveHomeRoute(pathname = "") {
  if (pathname.startsWith("/admin") && hasAdminSession()) {
    return "/admin/dashboard";
  }
  if (hasClientSession()) {
    return "/dashboard";
  }
  if (hasAdminSession()) {
    return "/admin/dashboard";
  }
  return "/";
}

function hasToken(key) {
  return getToken(key) !== "";
}

function getToken(key) {
  const storage = getStorage();
  if (!storage) {
    return "";
  }
  return storage.getItem(key) ?? "";
}

function saveToken(key, token) {
  const storage = getStorage();
  if (!storage) {
    return;
  }
  storage.setItem(key, token);
}

function clearToken(key) {
  const storage = getStorage();
  if (!storage) {
    return;
  }
  storage.removeItem(key);
}

function getStorage() {
  if (typeof window === "undefined") {
    return null;
  }
  return window.localStorage;
}
