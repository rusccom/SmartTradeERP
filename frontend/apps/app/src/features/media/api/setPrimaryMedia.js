import { postJSON } from "@smarterp/api/http";

export function setPrimaryMedia(paths, mediaID) {
  return postJSON(paths.primary(mediaID), {});
}
