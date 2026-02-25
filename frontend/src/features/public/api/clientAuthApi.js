import { apiPaths } from "../../../shared/api/client";
import { postJSON } from "../../../shared/api/http";

export function loginClient(payload) {
  return postJSON(apiPaths.clientLogin, payload);
}

export function registerClient(payload) {
  return postJSON(apiPaths.clientRegister, payload);
}
