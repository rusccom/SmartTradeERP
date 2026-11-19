import { cart } from "/_assets/cart-store.js";

const menu = document.querySelector(".menu-toggle");
const nav = document.getElementById(menu.getAttribute("aria-controls"));
const count = document.querySelector(".bag-count");
const total = document.querySelector("[data-cart-total]");

function showCartSummary() {
  total.textContent = cart?.total || "";
  count.hidden = !cart?.lines.length;
}

window.addEventListener("cartchange", showCartSummary);
showCartSummary();

function toggleMenu(open) {
  menu.setAttribute("aria-expanded", String(open));
}

menu.hidden = false;
menu.addEventListener("click", () => toggleMenu(menu.getAttribute("aria-expanded") === "false"));
nav.closest("header").addEventListener("keydown", event => {
  if (event.key !== "Escape" || menu.getAttribute("aria-expanded") !== "true") return;
  toggleMenu(false);
  menu.focus();
});

for (const link of nav.querySelectorAll("a")) {
  if (link.pathname === location.pathname) link.setAttribute("aria-current", "page");
}
