import { useState } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";
import InPostFields from "./InPostFields";
import { verifyAndSave } from "./api";

export default function InPostForm({ connection, onSaved, onStart }) {
  const { t } = useI18n();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");
  const submit = event => submitForm(event, { connection, onSaved, onStart, setPending, setError });
  return (
    <form className="inpost-form" method="post" autoComplete="off" onSubmit={submit}>
      <fieldset disabled={pending}>
        <legend>{t("inpost.credentials")}</legend>
        <InPostFields connection={connection} />
        {error && <p className="inpost-error" role="alert">{error}</p>}
        <button className="inpost-save" type="submit">{t(pending ? "inpost.busy" : "inpost.save")}</button>
      </fieldset>
    </form>
  );
}

async function submitForm(event, { connection, onSaved, onStart, setPending, setError }) {
  event.preventDefault();
  const form = event.currentTarget;
  const fields = new FormData(form);
  onStart();
  setPending(true);
  setError("");
  try {
    const payload = { ...Object.fromEntries(fields), revision: connection?.revision || 0 };
    const verified = await verifyAndSave(payload);
    form.elements.namedItem("token").value = "";
    onSaved(verified);
  } catch (error) {
    setError(error.message);
  } finally {
    setPending(false);
  }
}
