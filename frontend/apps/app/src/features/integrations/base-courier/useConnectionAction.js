import { useRef, useState } from "react";

export default function useConnectionAction(refresh) {
  const lock = useRef(false);
  const [status, setStatus] = useState({});
  async function run(request, notice) {
    if (lock.current) return;
    lock.current = true;
    setStatus({ busy: true });
    try {
      await request();
      setStatus({ notice });
      refresh();
    } catch (error) {
      setStatus({ error: error.message });
    } finally {
      lock.current = false;
    }
  }
  function reload() { setStatus({}); refresh(); }
  return { ...status, run, reload };
}
