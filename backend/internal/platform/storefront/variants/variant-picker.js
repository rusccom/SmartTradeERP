const picker = document.querySelector("[data-variant-picker]");
const product = picker.closest("article");
const errorNotice = product.querySelector("[data-variant-error]");
const preview = document.querySelector('meta[name="sf-preview"]').content === "true";
let pending;

function setBusy(busy) {
  picker.setAttribute("aria-busy", String(busy));
  product.querySelector('[data-variant-part="purchase"]').inert = busy;
}

async function loadVariant(url, signal) {
  const headers = preview ? { "X-Storefront-Preview": "1" } : {};
  const response = await fetch(url, { headers, signal });
  if (!response.ok) throw new Error("Variant unavailable");
  return new DOMParser().parseFromString(await response.text(), "text/html");
}

function replaceParts(page) {
  const parts = Array.from(product.querySelectorAll("[data-variant-part]"), current => {
    const next = page.querySelector(`[data-variant-part="${current.dataset.variantPart}"]`);
    if (!next) throw new Error("Incomplete product");
    return [current, next];
  });
  for (const [current, next] of parts) current.replaceWith(next);
}

function updateMetadata(page) {
  for (const selector of ['meta[property="og:image"]', 'script[type="application/ld+json"]']) {
    const next = page.head.querySelector(selector);
    document.head.querySelector(selector)?.remove();
    if (next) document.head.append(next);
  }
}

function showVariant(page, link, url) {
  const selected = page.querySelector("[data-variant-picker] [aria-current]");
  if (selected?.getAttribute("href") !== link.getAttribute("href")) throw new Error("Variant unavailable");
  replaceParts(page);
  updateMetadata(page);
  picker.querySelector("[aria-current]").removeAttribute("aria-current");
  link.setAttribute("aria-current", "true");
  history.replaceState(history.state, "", url);
  product.dispatchEvent(new Event("storefront:variant-change", { bubbles: true }));
}

async function switchVariant(link, controller) {
  const url = new URL(location.href);
  url.searchParams.set("variant", new URL(link.href).searchParams.get("variant"));
  setBusy(true);
  try {
    const page = await loadVariant(url, controller.signal);
    if (!controller.signal.aborted) showVariant(page, link, url);
  } catch {
    if (!controller.signal.aborted) errorNotice.hidden = false;
  } finally {
    if (pending === controller) setBusy(false);
  }
}

picker.addEventListener("click", event => {
  const link = event.target.closest("a[href]");
  if (!link || event.button !== 0 || event.ctrlKey || event.metaKey || event.shiftKey || event.altKey) return;
  event.preventDefault();
  pending?.abort();
  errorNotice.hidden = true;
  if (link.hasAttribute("aria-current")) return;
  pending = new AbortController();
  switchVariant(link, pending);
});
