import { ImageOff } from "lucide-react";
import { useState } from "react";

function ProductTableThumb({ url, name }) {
  const [failedUrl, setFailedUrl] = useState(null);
  if (url && failedUrl !== url) {
    return (
      <img
        className="product-table-thumb"
        src={url}
        alt={name}
        loading="lazy"
        onError={() => setFailedUrl(url)}
      />
    );
  }
  return (
    <span className="product-table-thumb product-table-thumb-empty" aria-hidden="true">
      <ImageOff size={16} strokeWidth={1.8} />
    </span>
  );
}

export default ProductTableThumb;
