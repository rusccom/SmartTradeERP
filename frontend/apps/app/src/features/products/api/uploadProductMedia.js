import { apiPaths } from "@smarterp/api/publicApi";
import { postJSON } from "@smarterp/api/http";
import { uploadDirectMedia } from "../../media/api/uploadDirectMedia";

export async function uploadProductMedia(productID, file) {
  const upload = await createProductMediaUpload(productID, file);
  await uploadDirectMedia(upload, file);
  return completeProductMediaUpload(productID, upload.id);
}

function createProductMediaUpload(productID, file) {
  return postJSON(apiPaths.productMedia(productID), {
    file_name: file.name,
    content_type: file.type,
    size_bytes: file.size,
  });
}

function completeProductMediaUpload(productID, mediaID) {
  return postJSON(apiPaths.productMediaComplete(productID, mediaID), {});
}
