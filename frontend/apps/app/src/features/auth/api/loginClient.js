import { apiPaths } from "@smarterp/api/publicApi";
import { postJSON } from "@smarterp/api/http";

export function loginClient(payload) {
  return postJSON(apiPaths.clientLogin, payload);
}
