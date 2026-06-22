export async function uploadDirectMedia(upload, file) {
  const response = await fetch(upload.upload_url, {
    method: upload.method || "PUT",
    headers: upload.headers || {},
    body: file,
  });
  if (!response.ok) {
    throw new Error(`Upload failed with status ${response.status}`);
  }
}
