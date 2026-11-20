import { useEffect, useRef, useState } from "react";
import { postJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";

function watch(ticket, connected, failed) {
  const controller = new AbortController();
  let timer;
  async function poll() {
    try {
      const result = await postJSON(apiPaths.allegroPrivatePoll, { id: ticket.id }, controller.signal);
      if (controller.signal.aborted) return;
      if (result.connected) connected();
      else timer = setTimeout(poll, result.retry_after * 1000);
    } catch (error) {
      if (!controller.signal.aborted) failed(error.message);
    }
  }
  timer = setTimeout(poll, ticket.interval * 1000);
  return () => { controller.abort(); clearTimeout(timer); };
}

export default function useDeviceAuthorization(version, reload) {
  const [state, setState] = useState({});
  const latest = useRef(reload);
  latest.current = reload;
  useEffect(() => {
    if (state.ticket) return watch(state.ticket, () => latest.current(), error => setState({ error }));
  }, [state.ticket]);
  async function start() {
    setState({ busy: true });
    try {
      const ticket = await postJSON(apiPaths.allegroPrivateDevice, { version });
      setState({ ticket });
    } catch (error) { setState({ error: error.message }); }
  }
  return { ...state, start };
}
