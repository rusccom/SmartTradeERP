import { PickupMap, polishPostalCode } from "./pickup-map.js";

export class CheckoutDelivery {
  constructor(form, dialog) {
    this.form = form;
    this.fields = form.elements;
    this.button = form.querySelector("[data-open-map]");
    this.map = new PickupMap(dialog, (point, postal) => this.choose(point, postal));
    this.button.addEventListener("click", () => this.open());
    form.querySelector(".delivery-methods").addEventListener("change", () => this.changeMethod());
    this.fields.country.addEventListener("change", () => this.resetPoint());
    this.fields.postal_code.addEventListener("input", () => this.postalChanged());
    this.fields.postal_code.addEventListener("blur", () => this.formatPostal());
    form.addEventListener("reset", () => queueMicrotask(() => this.resetPoint()));
    this.resetPoint();
  }

  resetPoint() {
    this.point = null;
    this.searchPostal = "";
    const country = this.fields.country.value;
    const available = !country || country === "PL";
    const option = this.form.querySelector("[data-pickup-option]");
    option.hidden = !available;
    option.querySelector("input").disabled = !available;
    if (!available) this.fields.delivery.value = "address";
    this.sync();
  }

  changeMethod() {
    if (this.fields.delivery.value === "pickup" && !this.fields.country.value) this.fields.country.value = "PL";
    this.sync();
    if (this.fields.delivery.value === "pickup") this.open();
  }

  sync() {
    const pickup = this.fields.delivery.value === "pickup";
    const polish = this.fields.country.value === "PL";
    this.form.querySelector("[data-address]").hidden = pickup;
    this.fields.address.disabled = pickup;
    this.fields.postal_code.required = polish && !pickup;
    this.fields.postal_code.pattern = polish && !pickup ? "[0-9]{2}-?[0-9]{3}" : ".*";
    this.fields.postal_code.placeholder = polish ? "00-000" : "";
    this.form.querySelector("[data-pickup-selection]").hidden = !pickup;
    this.showPoint();
  }

  postalChanged() {
    if (polishPostalCode(this.fields.postal_code.value) === this.searchPostal) return;
    this.point = null;
    this.showPoint();
  }

  formatPostal() {
    const input = this.fields.postal_code;
    if (this.fields.country.value === "PL") input.value = polishPostalCode(input.value) || input.value.trim();
  }

  showPoint() {
    const p = this.point;
    this.form.querySelector("[data-point-name]").textContent = p?.name || "";
    this.form.querySelector("[data-point-address]").textContent = p ? [p.street, p.postal_code + " " + p.city].join(", ") : "";
    this.button.textContent = p ? this.button.dataset.change : this.button.dataset.choose;
  }

  open() {
    this.map.open(this.fields.postal_code.value, this.point?.carrier);
  }

  choose(point, postal) {
    this.point = point;
    this.searchPostal = postal;
    this.fields.postal_code.value = postal;
    this.showPoint();
    queueMicrotask(() => this.button.focus());
  }

  orderDetails() {
    const { delivery, note, ...customer } = Object.fromEntries(new FormData(this.form));
    if (delivery === "pickup" && !this.point) {
      this.open();
      return null;
    }
    return { customer, note, pickup: delivery === "pickup" ? this.point : undefined };
  }
}
