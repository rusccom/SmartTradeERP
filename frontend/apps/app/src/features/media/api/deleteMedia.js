import { deleteJSON } from "@smarterp/api/http";

export function deleteMedia(paths, mediaID) {
  return deleteJSON(paths.item(mediaID));
}
