import { feedback, quantity, readCart, setQuantity } from "./cart-store.js";
import { startCart } from "./cart-page.js";

function updateBadges() {
  const count = readCart().reduce((sum, item) => sum + Number(item.qty), 0);
  for (const badge of document.querySelectorAll("[data-sf-cart-count]")) badge.textContent = count;
}

document.addEventListener("click", event => {
  const button = event.target.closest("[data-sf-add]");
  if (!button) return;
  const item = { variant_id: button.dataset.sfAdd };
  setQuantity(item, quantity(item) + 1);
  feedback(button);
});

window.addEventListener("cartchange", updateBadges);
window.addEventListener("pageshow", updateBadges);
updateBadges();
const root = document.querySelector("[data-sf-cart]");
if (root) startCart(root);
