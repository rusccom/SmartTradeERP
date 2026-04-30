export const initialCurrencyForm = Object.freeze({
  currencyID: "",
  displaySymbol: "",
});

export function patchCurrencyForm(form, event) {
  const { name, value } = event.target;
  return { ...form, [name]: value };
}

export function toCurrencyPayload(form, fallbackID) {
  return {
    currency_id: form.currencyID || fallbackID,
    display_symbol: form.displaySymbol.trim(),
    is_base: true,
  };
}
