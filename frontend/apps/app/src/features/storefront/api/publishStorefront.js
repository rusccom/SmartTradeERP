import { postJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export function publishStorefront() {
  return postJSON(apiPaths.storefrontPublish, {});
}
