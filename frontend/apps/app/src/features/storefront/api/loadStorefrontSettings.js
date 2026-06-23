import { getJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

export async function loadStorefrontSettings(signal) {
  const { data } = await getJSON(apiPaths.storefrontSettings, undefined, signal);
  return data;
}
