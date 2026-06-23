import { getJSON } from "@smarterp/api/http";

export function loadMedia(paths, signal) {
  return getJSON(paths.list, {}, signal).then((response) => response.data || []);
}
