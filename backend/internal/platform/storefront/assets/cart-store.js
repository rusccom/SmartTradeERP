const meta = name => document.querySelector('meta[name="' + name + '"]')?.content;
export const preview = meta("sf-preview") === "true";
const storageKey = ["sf-cart", meta("sf-cart-scope"), preview ? "preview" : "shop"].join("/");

export function lineKey(item) {
  const choices = Object.entries(item.picks || {}).sort(([a], [b]) => a.localeCompare(b));
  return JSON.stringify([item.variant_id, choices]);
}

export function readCart() {
  try {
    const items = JSON.parse(localStorage.getItem(storageKey));
    return Array.isArray(items) ? items : [];
  } catch { return []; }
}

export function writeCart(items) {
  localStorage.setItem(storageKey, JSON.stringify(items));
  window.dispatchEvent(new Event("cartchange"));
}

export function quantity(item, items = readCart()) {
  return Number(items.find(row => lineKey(row) === lineKey(item))?.qty || 0);
}

export function setQuantity(item, qty) {
  const items = readCart();
  const index = items.findIndex(row => lineKey(row) === lineKey(item));
  if (qty <= 0 && index >= 0) items.splice(index, 1);
  if (qty > 0 && index >= 0) items[index].qty = qty;
  if (qty > 0 && index < 0) items.push({ ...item, qty });
  writeCart(items);
}

export async function shopRequest(action, payload) {
  const headers = { "Content-Type": "application/json" };
  if (preview) headers["X-Storefront-Preview"] = "1";
  const res = await fetch("/api/storefront/" + action, { method: "POST", headers, body: JSON.stringify(payload) });
  if (!res.ok) throw new Error(String(res.status));
  return (await res.json()).data;
}

export function feedback(button) {
  button.dataset.label ||= button.textContent;
  button.textContent = button.dataset.added || "Added ✓";
  window.setTimeout(() => { button.textContent = button.dataset.label; }, 1000);
}

window.addEventListener("storage", event => {
  if (event.key === storageKey || event.key === null) window.dispatchEvent(new Event("cartchange"));
});
