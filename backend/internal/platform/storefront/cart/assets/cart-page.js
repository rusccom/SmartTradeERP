import { cart, changeQuantity, lineKey, placeOrder, refreshCart, reportCartError } from "./cart-store.js";

const polish = document.documentElement.lang === "pl";
const message = (pl, en) => polish ? pl : en;

function paintLine(line) {
  const node = document.querySelector("#cart-line").content.firstElementChild.cloneNode(true);
  node.dataset.key = lineKey(line);
  node.querySelector("[data-name]").textContent = line.product_name || message("Produkt niedostępny", "Unavailable item");
  node.querySelector("[data-variant]").textContent = line.variant_name === "Default" ? "" : line.variant_name;
  node.querySelector("[data-qty]").textContent = line.qty;
  node.querySelector("[data-price]").textContent = line.line_total || "—";
  node.querySelector("[data-unavailable]").hidden = line.purchasable;
  for (const c of line.components || []) {
    const li = document.createElement("li");
    li.textContent = c.qty + " × " + c.name;
    node.querySelector("[data-components]").append(li);
  }
  return node;
}

function changeLine(event) {
  const button = event.target.closest("[data-action]");
  if (!button) return;
  const key = button.closest("[data-key]").dataset.key;
  const item = cart?.lines.find(row => lineKey(row) === key);
  if (!item) return;
  const delta = button.dataset.action === "inc" ? 1 : -1;
  const { variant_id, picks } = item;
  const remove = button.dataset.action === "remove";
  changeQuantity({ variant_id, picks }, remove ? 0 : delta, !remove);
}

function paint(root, data) {
  root.querySelector("[data-lines]").replaceChildren(...data.lines.map(paintLine));
  root.querySelector("[data-total]").textContent = data.total;
  root.querySelector(".submit").disabled = data.lines.some(line => !line.purchasable);
}

function render(root) {
  if (!cart) return;
  root.querySelector("[data-empty]").hidden = cart.lines.length > 0;
  root.querySelector("[data-summary]").hidden = !cart.lines.length;
  paint(root, cart);
}

async function checkout(root, event) {
  event.preventDefault();
  const form = event.currentTarget;
  if (form.inert || form.querySelector("button").disabled) return;
  const { note, ...customer } = Object.fromEntries(new FormData(form));
  form.inert = true;
  try {
    const order = await placeOrder({ customer, note });
    form.reset();
    root.querySelector("[data-message]").textContent = confirmation(order);
  } catch (error) { reportCartError(error); await refreshCart().catch(reportCartError); }
  finally { form.inert = false; }
}

function confirmation(order) {
  if (order.status === "preview") return message("Zamówienie próbne sprawdzone. Suma: ", "Preview checkout complete. Total: ") + order.total;
  return message("Zamówienie przyjęte: ", "Order received: ") + order.number;
}

export function startCart(root) {
  root.addEventListener("click", changeLine);
  root.querySelector("form").addEventListener("submit", event => checkout(root, event));
  window.addEventListener("cartchange", () => render(root));
  render(root);
}
