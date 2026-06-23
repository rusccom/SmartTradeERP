import { useMemo, useRef } from "react";
import { ImagePlus } from "lucide-react";

import { productMediaPaths, variantMediaPaths } from "../api/mediaPaths";
import { useMedia } from "../model/useMedia";
import MediaTile from "./MediaTile";
import "./media.css";

const ACCEPT = "image/jpeg,image/png,image/webp,image/gif";

function MediaManager({ kind, ownerId, t }) {
  const paths = useMemo(() => buildPaths(kind, ownerId), [kind, ownerId]);
  const media = useMedia(paths);
  const inputRef = useRef(null);
  const handlePick = (event) => {
    const file = event.target.files && event.target.files[0];
    event.target.value = "";
    if (file) media.upload(file);
  };
  return (
    <div className="media-manager">
      <div className="media-grid">
        {media.items.map((item) => (
          <MediaTile
            key={item.id}
            item={item}
            busy={media.busy}
            t={t}
            onPrimary={() => media.makePrimary(item.id)}
            onRemove={() => media.remove(item.id)}
          />
        ))}
        <button type="button" className="media-add" disabled={media.busy} onClick={() => inputRef.current && inputRef.current.click()}>
          <ImagePlus size={20} strokeWidth={1.8} aria-hidden="true" />
          <span>{t("products.form.mediaUpload")}</span>
        </button>
      </div>
      <input ref={inputRef} className="media-file-input" type="file" accept={ACCEPT} onChange={handlePick} />
      {media.error && <p className="media-error">{media.error}</p>}
    </div>
  );
}

function buildPaths(kind, ownerId) {
  return kind === "variant" ? variantMediaPaths(ownerId) : productMediaPaths(ownerId);
}

export default MediaManager;
