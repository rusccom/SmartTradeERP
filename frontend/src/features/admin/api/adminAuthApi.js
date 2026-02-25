import { apiPaths } from "../../../shared/api/client";
import { postJSON } from "../../../shared/api/http";

export function loginAdmin(payload) {
  return postJSON(apiPaths.adminLogin, payload);
}
