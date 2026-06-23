import { useCallback } from "react";

import { productMediaPaths } from "../../media/api/mediaPaths";
import { uploadMedia } from "../../media/api/uploadMedia";

const ACCEPT = "image/jpeg,image/png,image/webp,image/gif";

export function useDescriptionImageUpload(productId) {
  return useCallback(() => {
    if (!productId) return Promise.resolve(null);
    return pickAndUpload(productMediaPaths(productId));
  }, [productId]);
}

function pickAndUpload(paths) {
  return new Promise((resolve) => {
    const input = document.createElement("input");
    input.type = "file";
    input.accept = ACCEPT;
    input.onchange = () => resolve(uploadPicked(paths, input.files));
    input.click();
  });
}

async function uploadPicked(paths, files) {
  const file = files && files[0];
  if (!file) return null;
  const media = await uploadMedia(paths, file);
  return media?.url || null;
}
