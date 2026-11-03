import { feedback, quantity, readCart, setQuantity, shopRequest } from "/_assets/cart-store.js";

const builder = document.querySelector(".builder");
const rows = Array.from(document.querySelectorAll("[data-builder-row]"));

function selected(row) {
  return { variant_id: row.querySelector("[data-variant]").value };
}

function syncRow(row, items) {
  const option = row.querySelector("[data-variant]").selectedOptions[0];
  const count = quantity(selected(row), items);
  row.querySelector("[data-unit-price]").textContent = option?.dataset.price || "—";
  row.querySelector("[data-quantity]").textContent = count;
  row.querySelector("[data-picked]").checked = count > 0;
  row.querySelector("[data-picked]").disabled = !count && (!option || option.disabled);
  row.querySelector('[data-step="-1"]').disabled = !count;
  row.querySelector('[data-step="1"]').disabled = !option || option.disabled;
}

function bindRow(row) {
  const select = row.querySelector("select");
  const prior = [...select.options].find(option => quantity({ variant_id: option.value }) > 0);
  const option = prior || [...select.options].find(option => !option.disabled);
  if (option) select.value = option.value;
  select.addEventListener("change", () => syncRow(row, readCart()));
  row.querySelector("[data-picked]").addEventListener("change", event => {
    setQuantity(selected(row), event.target.checked ? 1 : 0);
  });
  row.addEventListener("click", event => {
    const step = event.target.closest("[data-step]");
    if (step) setQuantity(selected(row), quantity(selected(row)) + Number(step.dataset.step));
  });
}

let request = 0;
async function syncBuilder() {
  const current = ++request;
  const items = readCart();
  for (const row of rows) syncRow(row, items);
  const total = builder.querySelector("[data-builder-total]");
  try {
    const data = await shopRequest("cart", { items });
    if (current === request) total.textContent = data.total;
  } catch { if (current === request) total.textContent = "—"; }
}

function kitItem(form) {
  const options = [...form.querySelectorAll("[data-slot]")];
  return { variant_id: form.dataset.kit, picks: Object.fromEntries(options.map(el => [el.dataset.slot, el.value])) };
}

function bindKit(form) {
  const submit = form.querySelector("button");
  const quote = kitQuoter(form, submit);
  form.addEventListener("change", quote);
  form.addEventListener("submit", event => {
    event.preventDefault();
    if (submit.disabled) return;
    const item = kitItem(form);
    setQuantity(item, quantity(item) + 1);
    feedback(submit);
  });
  quote();
}

function kitQuoter(form, button) {
  let revision = 0;
  return async function quote() {
    const current = ++revision;
    button.disabled = true;
    try {
      const data = await shopRequest("cart", { items: [{ ...kitItem(form), qty: 1 }] });
      if (current !== revision) return;
      const line = data.lines[0];
      button.disabled = !line?.purchasable;
      form.querySelector("[data-kit-price]").textContent = line?.unit_price || "—";
      form.querySelector("[data-kit-status]").textContent = button.disabled ? "Wybrany skład jest niedostępny." : "";
    } catch { if (current === revision) form.querySelector("[data-kit-status]").textContent = "Nie udało się sprawdzić zestawu."; }
  };
}

if (builder) {
  rows.forEach(bindRow);
  window.addEventListener("cartchange", syncBuilder);
  window.addEventListener("pageshow", syncBuilder);
  syncBuilder();
}
document.querySelectorAll("[data-kit]").forEach(bindKit);
