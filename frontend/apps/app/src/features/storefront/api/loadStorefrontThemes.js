import { getJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export async function loadStorefrontThemes(signal) {
  const { data } = await getJSON(apiPaths.storefrontThemes, undefined, signal);
  return data || [];
}
