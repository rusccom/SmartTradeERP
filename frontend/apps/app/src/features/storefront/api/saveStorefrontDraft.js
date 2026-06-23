import { putJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export function saveStorefrontDraft(payload) {
  return putJSON(apiPaths.storefrontDraft, payload);
}
