import { useRef, useState } from "react";

export default function useAllegroAction(reload) {
  const lock = useRef(false);
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  async function run(action) {
    if (lock.current) return false;
    lock.current = true;
    setBusy(true);
    setError("");
    try { await action(); return true; }
    catch (reason) { setError(reason.message); return false; }
    finally { lock.current = false; setBusy(false); reload(); }
  }
  return { busy, error, run };
}
