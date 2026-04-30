import { Plus } from "lucide-react";
import { useState } from "react";

import { initialCurrencyForm, patchCurrencyForm, toCurrencyPayload } from "../model/currencyForm";

function CurrencyCreateForm({ labels, onSubmit, options }) {
  const [form, setForm] = useState(initialCurrencyForm);
  const [submitting, setSubmitting] = useState(false);
  const selectedID = form.currencyID || options[0]?.id || "";
  return (
    <form className="currency-form" onSubmit={(e) => submitForm(e, formState(form, selectedID, setForm, setSubmitting, onSubmit))}>
      {renderCurrencySelect(labels, options, selectedID, setForm)}
      {renderSymbolField(labels, form, setForm)}
      <button className="currency-primary-btn" type="submit" disabled={submitting || !selectedID}>
        <Plus size={16} /> {submitting ? labels.saving : labels.save}
      </button>
    </form>
  );
}

function formState(form, selectedID, setForm, setSubmitting, onSubmit) {
  return { form, onSubmit, selectedID, setForm, setSubmitting };
}

async function submitForm(event, state) {
  event.preventDefault();
  state.setSubmitting(true);
  try {
    await state.onSubmit(toCurrencyPayload(state.form, state.selectedID));
    state.setForm(initialCurrencyForm);
  } finally {
    state.setSubmitting(false);
  }
}

function renderCurrencySelect(labels, options, selectedID, setForm) {
  return (
    <label>
      <span>{labels.currency}</span>
      <select name="currencyID" required value={selectedID} onChange={(e) => setForm((v) => patchCurrencyForm(v, e))}>
        {options.map((item) => <option key={item.id} value={item.id}>{item.code} - {item.name}</option>)}
      </select>
    </label>
  );
}

function renderSymbolField(labels, form, setForm) {
  return (
    <label>
      <span>{labels.symbol}</span>
      <input name="displaySymbol" maxLength="8" value={form.displaySymbol} onChange={(e) => setForm((v) => patchCurrencyForm(v, e))} />
    </label>
  );
}

export default CurrencyCreateForm;
