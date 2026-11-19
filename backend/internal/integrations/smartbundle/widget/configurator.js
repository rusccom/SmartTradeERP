import { changeQuantity, feedback, shopRequest } from "/_assets/cart-store.js";

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
      form.querySelector("[data-kit-status]").textContent = button.disabled ? form.dataset.unavailable : "";
    } catch { if (current === revision) form.querySelector("[data-kit-status]").textContent = form.dataset.quoteError; }
  };
}

function bindKits(root) {
  root.querySelectorAll("[data-kit]").forEach(bindKit);
}

document.addEventListener("storefront:variant-change", event => bindKits(event.target));
bindKits(document);
