import { cart, changeQuantity, feedback, quantity, shopRequest } from "/_assets/cart-store.js";

const builder = document.querySelector(".builder");
const rows = Array.from(document.querySelectorAll("[data-builder-row]"));

function selected(row) {
  return { variant_id: row.querySelector("[data-variant]").value };
}

function syncRow(row) {
  const option = row.querySelector("[data-variant]").selectedOptions[0];
  const count = quantity(selected(row));
  const available = option?.dataset.available === "true";
  const image = row.querySelector("img");
  if (image && option?.dataset.image) image.src = option.dataset.image;
  row.querySelector("[data-unit-price]").textContent = option?.dataset.price || "—";
  row.querySelector("[data-quantity]").textContent = count;
  row.querySelector("[data-picked]").checked = count > 0;
  row.querySelector("[data-picked]").disabled = !count && !available;
  row.querySelector('[data-step="-1"]').disabled = !count;
  row.querySelector('[data-step="1"]').disabled = !available;
}

function bindRow(row) {
  const select = row.querySelector("select");
  const prior = [...select.options].find(option => quantity({ variant_id: option.value }) > 0);
  const option = prior || [...select.options].find(option => option.dataset.available === "true");
  if (option) select.value = option.value;
  select.addEventListener("change", () => syncRow(row));
  row.querySelector("[data-picked]").addEventListener("change", event => {
    changeQuantity(selected(row), event.target.checked ? 1 : 0, false).catch(() => syncRow(row));
  });
  row.addEventListener("click", event => {
    const step = event.target.closest("[data-step]");
    if (step) changeQuantity(selected(row), Number(step.dataset.step));
  });
}

function syncBuilder() {
  for (const row of rows) syncRow(row);
  builder.querySelector("[data-builder-total]").textContent = cart?.total || "—";
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
    changeQuantity(item, 1).then(() => feedback(submit), () => {});
  });
  quote();
}

function kitQuoter(form, button) {
  let revision = 0;
  return async function quote() {
    const current = ++revision;
    button.disabled = true;
    try {
      const data = await shopRequest("quote", "POST", { items: [{ ...kitItem(form), qty: 1 }] });
      if (current !== revision) return;
      const line = data.lines[0];
      button.disabled = !line?.purchasable;
      form.querySelector("[data-kit-price]").textContent = line?.unit_price || "—";
      form.querySelector("[data-kit-status]").textContent = button.disabled ? "Wybrany skład jest niedostępny." : "";
    } catch { if (current === revision) form.querySelector("[data-kit-status]").textContent = "Nie udało się sprawdzić zestawu."; }
  };
}

if (builder) {
  window.addEventListener("cartchange", () => rows.forEach(bindRow), { once: true });
  window.addEventListener("cartchange", syncBuilder);
  syncBuilder();
}
document.querySelectorAll("[data-kit]").forEach(bindKit);
