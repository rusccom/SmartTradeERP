const preview = document.querySelector('meta[name="sf-preview"]')?.content === "true";
const updates = "BroadcastChannel" in window ? new BroadcastChannel(preview ? "sf-draft-cart" : "sf-cart") : null;
export let cart = null;
let pending = Promise.resolve();

export function lineKey(item) {
  return JSON.stringify([item.variant_id, Object.entries(item.picks || {}).sort(([a], [b]) => a.localeCompare(b))]);
}

export function quantity(item) {
  return Number(cart?.lines.find(line => lineKey(line) === lineKey(item))?.qty || 0);
}

export async function shopRequest(action, method, payload) {
  const headers = { "Content-Type": "application/json" };
  if (preview) headers["X-Storefront-Preview"] = "1";
  const res = await fetch("/api/storefront/" + action, { method, headers, cache: "no-store", body: JSON.stringify(payload) });
  const result = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(result.error?.code || String(res.status));
  return result.data;
}

function display(next) {
  if (cart && next.version < cart.version) return;
  cart = next;
  window.dispatchEvent(new Event("cartchange"));
}

export async function refreshCart() {
  display(await shopRequest("cart", "GET"));
}

function editedItems(item, amount, relative) {
  const items = cart.lines.map(({ variant_id, picks, qty }) => ({ variant_id, picks, qty }));
  const index = items.findIndex(row => lineKey(row) === lineKey(item));
  const qty = relative ? quantity(item) + amount : amount;
  if (qty <= 0 && index >= 0) items.splice(index, 1);
  if (qty > 0 && index < 0) items.push({ ...item, qty });
  if (qty > 0 && index >= 0) items[index].qty = qty;
  return items;
}

async function change(item, amount, relative) {
  if (!cart) await refreshCart();
  for (let attempt = 0; ; attempt++) {
    const body = { version: cart.version, items: editedItems(item, amount, relative) };
    try {
      display(await shopRequest("cart", "PUT", body));
      updates?.postMessage("changed");
      return;
    } catch (error) {
      if (error.message !== "cart_changed" || attempt > 0) throw error;
      await refreshCart();
    }
  }
}

export function changeQuantity(item, amount, relative = true) {
  const action = pending.then(() => change(item, amount, relative));
  pending = action.catch(reportCartError);
  return action;
}

export async function placeOrder(details) {
  await pending;
  const version = cart.version;
  const result = await shopRequest("checkout", "POST", { ...details, version });
  display({ ...cart, version: version + 1, lines: [], total: "" });
  updates?.postMessage("changed");
  return result;
}

export function reportCartError(error) {
  let el = document.querySelector("[data-message], [data-cart-error]");
  if (!el) {
    el = document.createElement("p");
    el.dataset.cartError = "";
    el.setAttribute("role", "alert");
    document.querySelector("main").prepend(el);
  }
  const pl = document.documentElement.lang === "pl";
  if (error.message === "403") {
    el.textContent = pl ? "Podgląd wygasł. Otwórz go ponownie w ERP." : "Preview expired. Reopen it in the ERP.";
    return;
  }
  el.textContent = pl ? "Sprawdź koszyk i spróbuj ponownie." : "Please check your cart and try again.";
}

export function feedback(button) {
  button.dataset.label ||= button.textContent;
  button.textContent = button.dataset.added || "Added ✓";
  window.setTimeout(() => { button.textContent = button.dataset.label; }, 1000);
}

const reload = () => refreshCart().catch(reportCartError);
updates?.addEventListener("message", reload);
window.addEventListener("pageshow", reload);
