import { apiPaths } from "../../../shared/api/publicApi";
import { getJSON } from "../../../shared/api/http";

export function loadProductMedia(productID, signal) {
  return getJSON(apiPaths.productMedia(productID), {}, signal).then((response) => response.data || []);
}
