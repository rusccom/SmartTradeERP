import { apiPaths } from "@smarterp/api/publicApi";
import { putJSON } from "@smarterp/api/http";

export function saveBundleComponents(id, components) {
  return putJSON(apiPaths.bundleComponents(id), { components });
}
