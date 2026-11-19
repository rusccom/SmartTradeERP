import { cart, changeQuantity, feedback } from "./cart-store.js";
import { startCart } from "./cart-page.js";

function updateBadges() {
  if (!cart) return;
  const count = cart.lines.reduce((sum, item) => sum + Number(item.qty), 0);
  for (const badge of document.querySelectorAll("[data-sf-cart-count]")) badge.textContent = count;
}

document.addEventListener("click", event => {
  const button = event.target.closest("[data-sf-add]");
  if (!button) return;
  const item = { variant_id: button.dataset.sfAdd };
  changeQuantity(item, 1).then(() => feedback(button), () => {});
});

window.addEventListener("cartchange", updateBadges);
updateBadges();
const root = document.querySelector("[data-sf-cart]");
if (root) startCart(root);
