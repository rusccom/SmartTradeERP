import { apiPaths } from "../../../shared/api/client";
import { putJSON } from "../../../shared/api/http";

export function saveBundleComponents(id, components) {
  return putJSON(apiPaths.bundleComponents(id), { components });
}
