import { useEffect, useState } from "react";
import { loadConnection } from "./api";

export default function useInPostSettings() {
  const [access, setAccess] = useState(null);
  const [error, setError] = useState("");
  const [saved, setSaved] = useState(false);
  const [attempt, setAttempt] = useState(0);
  useEffect(() => {
    const controller = new AbortController();
    loadConnection(controller.signal).then(data => {
      if (!controller.signal.aborted) setAccess(data);
    }).catch(reason => { if (!controller.signal.aborted) setError(reason.message); });
    return () => controller.abort();
  }, [attempt]);
  function retry() { setError(""); setAttempt(value => value + 1); }
  function onSaved(connection) {
    setAccess(current => ({ ...current, connection }));
    setSaved(true);
  }
  return { access, error, saved, retry, onSaved, onStart: () => setSaved(false) };
}
