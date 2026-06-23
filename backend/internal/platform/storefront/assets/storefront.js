(function () {
  "use strict";

  var KEY = "sf_cart_v1";
  var CART_URL = "/api/storefront/cart";
  var CHECKOUT_URL = "/api/storefront/checkout";

  function readCart() {
    try {
      return JSON.parse(localStorage.getItem(KEY)) || [];
    } catch (e) {
      return [];
    }
  }

  function writeCart(items) {
    localStorage.setItem(KEY, JSON.stringify(items));
    updateBadges(items);
  }

  function totalQty(items) {
    return items.reduce(function (sum, item) {
      return sum + Number(item.qty || 0);
    }, 0);
  }

  function updateBadges(items) {
    var count = totalQty(items);
    document.querySelectorAll("[data-sf-cart-count]").forEach(function (el) {
      el.textContent = count;
    });
  }

  function addItem(variantId, qty) {
    var items = readCart();
    var found = items.find(function (item) {
      return item.variant_id === variantId;
    });
    if (found) {
      found.qty = Number(found.qty) + qty;
    } else {
      items.push({ variant_id: variantId, qty: qty });
    }
    writeCart(items);
  }

  function setQty(variantId, qty) {
    var items = readCart()
      .map(function (item) {
        return item.variant_id === variantId ? { variant_id: variantId, qty: qty } : item;
      })
      .filter(function (item) {
        return item.qty > 0;
      });
    writeCart(items);
  }

  function removeItem(variantId) {
    writeCart(
      readCart().filter(function (item) {
        return item.variant_id !== variantId;
      })
    );
  }

  function currentQty(variantId) {
    var found = readCart().find(function (item) {
      return item.variant_id === variantId;
    });
    return found ? Number(found.qty) : 0;
  }

  function postJSON(url, body) {
    return fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    }).then(parseEnvelope);
  }

  function parseEnvelope(res) {
    return res.text().then(function (raw) {
      var body = raw ? JSON.parse(raw) : {};
      if (!res.ok || body.error) {
        throw new Error((body.error && body.error.message) || "Request failed");
      }
      return body.data;
    });
  }

  function bindAddButtons() {
    document.querySelectorAll("[data-sf-add]").forEach(function (btn) {
      btn.addEventListener("click", function () {
        addItem(btn.getAttribute("data-sf-add"), 1);
        flash(btn);
      });
    });
  }

  function flash(btn) {
    if (!btn.getAttribute("data-label")) {
      btn.setAttribute("data-label", btn.textContent);
    }
    btn.textContent = btn.getAttribute("data-added") || "Added";
    setTimeout(function () {
      btn.textContent = btn.getAttribute("data-label");
    }, 1200);
  }

  function escapeHTML(value) {
    return String(value).replace(/[&<>"']/g, function (c) {
      return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c];
    });
  }

  function qtyButton(action, id, label) {
    return '<button type="button" data-action="' + action + '" data-variant-id="' + escapeHTML(id) + '">' + label + "</button>";
  }

  function variantLabel(line) {
    if (line.variant_name && line.variant_name !== "Default") {
      return " — " + escapeHTML(line.variant_name);
    }
    return "";
  }

  function lineHTML(line) {
    return (
      '<div class="sf-line">' +
      '<span class="sf-line-name">' + escapeHTML(line.product_name) + variantLabel(line) + "</span>" +
      '<span class="sf-line-qty">' + qtyButton("dec", line.variant_id, "−") +
      "<span>" + escapeHTML(line.qty) + "</span>" + qtyButton("inc", line.variant_id, "+") + "</span>" +
      '<span class="sf-line-price">' + escapeHTML(line.line_total) + "</span>" +
      qtyButton("remove", line.variant_id, "×") +
      "</div>"
    );
  }

  function toggle(root, selector, show) {
    var el = root.querySelector(selector);
    if (el) {
      el.hidden = !show;
    }
  }

  function showMessage(root, message) {
    var box = root.querySelector("[data-sf-checkout-msg]");
    if (box) {
      box.textContent = message;
    }
  }

  function showEmpty(root) {
    toggle(root, "[data-sf-cart-empty]", true);
    toggle(root, "[data-sf-cart-summary]", false);
    var box = root.querySelector("[data-sf-cart-lines]");
    if (box) {
      box.innerHTML = "";
    }
  }

  function paintCart(root, data) {
    toggle(root, "[data-sf-cart-empty]", false);
    toggle(root, "[data-sf-cart-summary]", true);
    var box = root.querySelector("[data-sf-cart-lines]");
    if (box) {
      box.innerHTML = data.lines.map(lineHTML).join("");
    }
    var total = root.querySelector("[data-sf-cart-total]");
    if (total) {
      total.textContent = data.total;
    }
  }

  function renderCart(root) {
    var items = readCart();
    if (items.length === 0) {
      showEmpty(root);
      return;
    }
    postJSON(CART_URL, { items: items })
      .then(function (data) {
        paintCart(root, data);
      })
      .catch(function (err) {
        showMessage(root, err.message);
      });
  }

  function onCartClick(root, event) {
    var action = event.target.getAttribute("data-action");
    if (!action) {
      return;
    }
    var id = event.target.getAttribute("data-variant-id");
    if (action === "remove") {
      removeItem(id);
    } else if (action === "inc") {
      addItem(id, 1);
    } else if (action === "dec") {
      setQty(id, currentQty(id) - 1);
    }
    renderCart(root);
  }

  function fieldValue(form, name) {
    var el = form.querySelector('[name="' + name + '"]');
    return el ? el.value : "";
  }

  function customerFrom(form) {
    return {
      name: fieldValue(form, "name"),
      email: fieldValue(form, "email"),
      phone: fieldValue(form, "phone"),
      address: fieldValue(form, "address"),
    };
  }

  function confirmOrder(root, data) {
    writeCart([]);
    toggle(root, "[data-sf-cart-summary]", false);
    toggle(root, "[data-sf-cart-empty]", false);
    var box = root.querySelector("[data-sf-confirmation]");
    if (box) {
      box.hidden = false;
      box.textContent = (box.getAttribute("data-label") || "Order placed:") + " " + data.number;
    }
  }

  function onCheckout(root, form, event) {
    event.preventDefault();
    var items = readCart();
    if (items.length === 0) {
      return;
    }
    showMessage(root, "");
    postJSON(CHECKOUT_URL, { items: items, customer: customerFrom(form), note: fieldValue(form, "note") })
      .then(function (data) {
        confirmOrder(root, data);
      })
      .catch(function (err) {
        showMessage(root, err.message);
      });
  }

  function initCart(root) {
    renderCart(root);
    root.addEventListener("click", function (event) {
      onCartClick(root, event);
    });
    var form = root.querySelector("[data-sf-checkout]");
    if (form) {
      form.addEventListener("submit", function (event) {
        onCheckout(root, form, event);
      });
    }
  }

  function init() {
    updateBadges(readCart());
    bindAddButtons();
    var cart = document.querySelector("[data-sf-cart]");
    if (cart) {
      initCart(cart);
    }
  }

  if (document.readyState !== "loading") {
    init();
  } else {
    document.addEventListener("DOMContentLoaded", init);
  }
})();
