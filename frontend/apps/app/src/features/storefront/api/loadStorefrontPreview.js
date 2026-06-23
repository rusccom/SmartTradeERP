import { getJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export async function loadStorefrontPreview(signal) {
  const { data } = await getJSON(apiPaths.storefrontPreview, undefined, signal);
  return data;
}
