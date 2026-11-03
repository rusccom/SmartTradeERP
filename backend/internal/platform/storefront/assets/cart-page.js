import { lineKey, quantity, readCart, setQuantity, shopRequest, writeCart } from "./cart-store.js";

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
  const item = readCart().find(row => lineKey(row) === key);
  if (!item) return;
  const delta = button.dataset.action === "inc" ? 1 : -1;
  setQuantity(item, button.dataset.action === "remove" ? 0 : quantity(item) + delta);
}

function paint(root, data) {
  root.querySelector("[data-lines]").replaceChildren(...data.lines.map(paintLine));
  root.querySelector("[data-total]").textContent = data.total;
  root.querySelector(".submit").disabled = data.lines.some(line => !line.purchasable);
}

function showError(root, error) {
  const expired = error.message === "403";
  root.querySelector("[data-message]").textContent = expired
    ? message("Podgląd wygasł. Otwórz go ponownie w ERP.", "Preview expired. Reopen it in the ERP.")
    : message("Nie udało się potwierdzić koszyka. Odśwież stronę i spróbuj ponownie.", "Could not confirm your cart. Refresh and try again.");
}

function refresher(root) {
  let version = 0;
  return async function refresh() {
    const current = ++version;
    const items = readCart();
    root.querySelector(".submit").disabled = true;
    root.querySelector("[data-empty]").hidden = items.length > 0;
    root.querySelector("[data-summary]").hidden = !items.length;
    if (!items.length) return;
    try {
      const data = await shopRequest("cart", { items });
      if (version === current) paint(root, data);
    } catch (error) { if (version === current) showError(root, error); }
  };
}

async function checkout(root, event) {
  event.preventDefault();
  const form = event.currentTarget;
  if (form.inert || form.querySelector("button").disabled) return;
  const items = readCart();
  const { note, ...customer } = Object.fromEntries(new FormData(form));
  form.inert = true;
  try {
    const order = await shopRequest("checkout", { items, customer, note });
    consumeItems(items);
    form.reset();
    root.querySelector("[data-message]").textContent = confirmation(order);
  } catch (error) { showError(root, error); }
  finally { form.inert = false; }
}

function consumeItems(submitted) {
  const remaining = readCart().map(item => ({ ...item, qty: Number(item.qty) - quantity(item, submitted) }));
  writeCart(remaining.filter(item => item.qty > 0));
}

function confirmation(order) {
  if (order.status === "preview") return message("Zamówienie próbne sprawdzone. Suma: ", "Preview checkout complete. Total: ") + order.total;
  return message("Zamówienie przyjęte: ", "Order received: ") + order.number;
}

export function startCart(root) {
  const refresh = refresher(root);
  root.addEventListener("click", changeLine);
  root.querySelector("form").addEventListener("submit", event => checkout(root, event));
  window.addEventListener("cartchange", refresh);
  window.addEventListener("pageshow", refresh);
  refresh();
}
