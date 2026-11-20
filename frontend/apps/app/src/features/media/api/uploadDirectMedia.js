import { translate } from "@smarterp/i18n/translate";

export async function uploadDirectMedia(upload, file) {
  const response = await fetch(upload.upload_url, {
    method: upload.method || "PUT",
    headers: upload.headers || {},
    body: file,
  });
  if (!response.ok) {
    throw new Error(translate("media.uploadFailed", { status: response.status }));
  }
}
