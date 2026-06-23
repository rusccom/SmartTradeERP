import { useCallback, useEffect, useState } from "react";

import { deleteMedia } from "../api/deleteMedia";
import { loadMedia } from "../api/loadMedia";
import { setPrimaryMedia } from "../api/setPrimaryMedia";
import { uploadMedia } from "../api/uploadMedia";

export function useMedia(paths) {
  const [items, setItems] = useState([]);
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const reload = useCallback((signal) => {
    loadMedia(paths, signal).then(setItems).catch((err) => ignoreAbort(err, setError));
  }, [paths]);
  useEffect(() => {
    const controller = new AbortController();
    reload(controller.signal);
    return () => controller.abort();
  }, [reload]);
  const run = useCallback(async (task) => {
    setBusy(true);
    setError("");
    try {
      await task();
      reload();
    } catch (err) {
      setError(err.message);
    } finally {
      setBusy(false);
    }
  }, [reload]);
  return {
    items,
    error,
    busy,
    upload: (file) => run(() => uploadMedia(paths, file)),
    remove: (mediaID) => run(() => deleteMedia(paths, mediaID)),
    makePrimary: (mediaID) => run(() => setPrimaryMedia(paths, mediaID)),
  };
}

function ignoreAbort(err, setError) {
  if (err.name === "AbortError") return;
  setError(err.message);
}
