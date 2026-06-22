import { apiPaths } from "@smarterp/api/publicApi";
import { postJSON } from "@smarterp/api/http";

export function loginAdmin(payload) {
  return postJSON(apiPaths.adminLogin, payload);
}
