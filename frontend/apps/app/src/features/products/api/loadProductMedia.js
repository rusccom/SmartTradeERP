import { apiPaths } from "@smarterp/api/publicApi";
import { getJSON } from "@smarterp/api/http";

export function loadProductMedia(productID, signal) {
  return getJSON(apiPaths.productMedia(productID), {}, signal).then((response) => response.data || []);
}
