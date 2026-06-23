import { apiPaths } from "@smarterp/api/publicApi";

export function productMediaPaths(productId) {
  return {
    list: apiPaths.productMedia(productId),
    complete: (mediaID) => apiPaths.productMediaComplete(productId, mediaID),
    item: (mediaID) => apiPaths.productMediaItem(productId, mediaID),
    primary: (mediaID) => apiPaths.productMediaPrimary(productId, mediaID),
  };
}

export function variantMediaPaths(variantId) {
  return {
    list: apiPaths.variantMedia(variantId),
    complete: (mediaID) => apiPaths.variantMediaComplete(variantId, mediaID),
    item: (mediaID) => apiPaths.variantMediaItem(variantId, mediaID),
    primary: (mediaID) => apiPaths.variantMediaPrimary(variantId, mediaID),
  };
}
