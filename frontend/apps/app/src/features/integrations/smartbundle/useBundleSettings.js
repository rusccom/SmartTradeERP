import { useEffect, useState } from "react";
import { loadSettings, saveSettings } from "./api";

export function useBundleSettings(bundle) {
  const [editor, setEditor] = useState(null);
  const [message, setMessage] = useState("");
  const [busy, setBusy] = useState(false);
  useEffect(() => {
    const controller = new AbortController();
    setEditor(null);
    setMessage("");
    loadSettings(bundle.variant_id, controller.signal).then(data => {
      if (!controller.signal.aborted) setEditor(data);
    }).catch(error => { if (!controller.signal.aborted) setMessage(error.message); });
    return () => controller.abort();
  }, [bundle]);
  const patch = change => setEditor(old => ({ ...old, settings: { ...old.settings, ...change } }));
  const save = () => submit(bundle.variant_id, editor.settings, { setBusy, setMessage });
  return { editor, message, busy, patch, save };
}

async function submit(id, settings, { setBusy, setMessage }) {
  setBusy(true);
  setMessage("");
  try {
    await saveSettings(id, settings);
    setMessage("Saved");
  } catch (error) { setMessage(error.message); }
  finally { setBusy(false); }
}
