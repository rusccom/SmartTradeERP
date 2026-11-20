import { useEffect, useState } from "react";
import { loadSettings, saveSettings } from "./api";

export function useBundleSettings(bundle) {
  const [editor, setEditor] = useState(null);
  const [feedback, setFeedback] = useState(null);
  const [busy, setBusy] = useState(false);
  useEffect(() => loadEditor(bundle.variant_id, setEditor, setFeedback), [bundle]);
  const patch = change => setEditor(old => ({ ...old, settings: { ...old.settings, ...change } }));
  const save = () => submit(bundle.variant_id, editor.settings, { setBusy, setFeedback });
  return { editor, feedback, busy, patch, save };
}

function loadEditor(id, setEditor, setFeedback) {
  const controller = new AbortController();
  setEditor(null);
  setFeedback(null);
  loadSettings(id, controller.signal).then(data => {
    if (!controller.signal.aborted) setEditor(data);
  }).catch(error => { if (!controller.signal.aborted) setFeedback({ error: error.message }); });
  return () => controller.abort();
}

async function submit(id, settings, { setBusy, setFeedback }) {
  setBusy(true);
  setFeedback(null);
  try {
    await saveSettings(id, settings);
    setFeedback({ saved: true });
  } catch (error) { setFeedback({ error: error.message }); }
  finally { setBusy(false); }
}
