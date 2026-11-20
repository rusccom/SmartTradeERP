import { useRef, useState } from "react";
import { postJSON, putJSON } from "@smarterp/api/http";
import { apiPaths } from "@smarterp/api/publicApi";
import { entryFrom, entryPayload, patchEntry } from "./entryModel";
import useOrderQuote from "./useOrderQuote";

export default function useOrderEditor(initial, onPersist) {
  const [card, setCard] = useState(initial);
  const [form, setForm] = useState(() => entryFrom(initial));
  const [dirty, setDirty] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const lock = useRef(false);
  const quote = useOrderQuote(form, card);
  function change(path, value) {
    setForm((old) => patchEntry(old, path, value));
    setError("");
    setDirty(true);
  }
  const state = { card, form, dirty, setCard, setForm, setDirty, setBusy, setError, lock, onPersist };
  return { card, form, quote, busy, dirty, error, change, setCard, setError, save: (format) => saveOrder(state, format) };
}

async function persist(s) {
  if (s.card && !s.dirty) return s.card;
  const body = entryPayload(s.form);
  const card = s.card ? await putJSON(apiPaths.order(s.card.id), body) : await postJSON(apiPaths.orders, body);
  s.setCard(card);
  s.setForm(entryFrom(card));
  s.setDirty(false);
  return card;
}

async function saveOrder(s, format) {
  if (s.lock.current) return;
  s.lock.current = true;
  s.setBusy(true);
  s.setError("");
  let saved;
  let error = "";
  try {
    saved = await persist(s);
    if (format) saved = { ...saved, invoice: await postJSON(apiPaths.orderInvoice(saved.id), { format, revision: saved.revision }) };
    s.setCard(saved);
  } catch (e) {
    error = e.message;
    s.setError(error);
  } finally {
    s.lock.current = false;
    s.setBusy(false);
  }
  if (saved) s.onPersist(saved, error);
}
