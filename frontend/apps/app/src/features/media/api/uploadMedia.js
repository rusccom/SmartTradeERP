import { postJSON } from "@smarterp/api/http";

import { uploadDirectMedia } from "./uploadDirectMedia";

export async function uploadMedia(paths, file) {
  const upload = await postJSON(paths.list, {
    file_name: file.name,
    content_type: file.type,
    size_bytes: file.size,
  });
  await uploadDirectMedia(upload, file);
  return postJSON(paths.complete(upload.id), {});
}
