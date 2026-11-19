import { useState } from "react";
import { useOutletContext } from "react-router-dom";

const blankStatus = { id: null, name: "" };

export function useStatusDraft() {
  const [statuses, setStatuses] = useOutletContext();
  const [editor, edit] = useState(blankStatus);
  const name = editor.name.trim();
  const duplicate = statuses.some((status) => status.id !== editor.id && sameName(status.name, name));
  const canApply = Boolean(name) && !duplicate;
  const reset = () => edit(blankStatus);

  function apply(event) {
    event.preventDefault();
    if (!canApply) return;
    const status = { id: editor.id ?? crypto.randomUUID(), name };
    setStatuses((current) => editor.id ? renameStatus(current, status) : [...current, status]);
    reset();
  }

  return { statuses, editor, edit, reset, duplicate, canApply, apply };
}

function sameName(left, right) {
  return left.localeCompare(right, "uk", { sensitivity: "accent" }) === 0;
}

function renameStatus(statuses, updated) {
  return statuses.map((status) => status.id === updated.id ? updated : status);
}
