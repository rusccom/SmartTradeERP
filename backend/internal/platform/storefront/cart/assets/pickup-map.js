const baseCourier = "https://api.basecourier.com";

export function polishPostalCode(value) {
  const match = /^(\d{2})-?(\d{3})$/.exec(value.trim());
  return match ? match[1] + "-" + match[2] : "";
}

function cityName(value) {
  try { return decodeURIComponent(value); } catch { return value; }
}

function streetName(address, postal, city) {
  const suffix = /[\s,|]+(\d{2}-?\d{3})[\s,|]+(.+)$/.exec(address);
  if (suffix && polishPostalCode(suffix[1]) === postal && suffix[2].localeCompare(city, "pl", { sensitivity: "base" }) === 0) {
    return address.slice(0, suffix.index).trim();
  }
  return address;
}

function pointFromMap(value, carrier, network) {
  if (value?.apiType !== carrier) return null;
  const id = String(value.pointName ?? "").trim();
  const city = cityName(String(value.city ?? "")).trim();
  const postal = polishPostalCode(String(value.postalCode ?? ""));
  const street = streetName(String(value.address ?? "").trim(), postal, city);
  if (!id || !street || !city || !postal) return null;
  const type = String(value.pointApiName || "").trim();
  const label = network.toLowerCase().includes(type.toLowerCase()) ? network : network + " · " + type;
  const name = label + " " + id;
  return { carrier, country: "PL", point_id: id, name, street, postal_code: postal, city };
}

export class PickupMap {
  constructor(dialog, onChoose) {
    this.dialog = dialog;
    this.onChoose = onChoose;
    this.postal = dialog.querySelector("[data-map-postal]");
    this.carrier = dialog.querySelector("[data-map-carrier]");
    this.container = dialog.querySelector("[data-map-frame]");
    this.hint = dialog.querySelector("[data-map-hint]");
    this.error = dialog.querySelector("[data-map-error]");
    this.retry = dialog.querySelector("[data-map-retry]");
    this.postal.addEventListener("input", () => this.render());
    this.carrier.addEventListener("change", () => this.render());
    this.retry.addEventListener("click", () => this.render());
    dialog.querySelector("[data-close-map]").addEventListener("click", () => dialog.close());
    dialog.addEventListener("close", () => this.clear());
    window.addEventListener("message", event => this.receive(event));
  }

  open(postal, carrier) {
    this.postal.value = polishPostalCode(postal) || postal.slice(0, 6);
    this.carrier.value = carrier || "inpost";
    this.dialog.showModal();
    this.render();
    if (!polishPostalCode(this.postal.value)) this.postal.focus();
  }

  clear() {
    this.container.replaceChildren();
    this.frame = null;
    this.candidate = null;
    this.error.hidden = true;
  }

  render() {
    this.clear();
    const postal = polishPostalCode(this.postal.value);
    this.hint.hidden = Boolean(postal);
    this.retry.hidden = !postal;
    if (!postal) return;
    const type = this.carrier.value === "inpost" ? "paczkomaty" : this.carrier.value;
    const params = new URLSearchParams({ api_type: type, postalCode: postal, country: "PL" });
    this.frame = document.createElement("iframe");
    this.frame.title = this.container.dataset.title + " – " + this.carrier.selectedOptions[0].text;
    this.frame.referrerPolicy = "no-referrer";
    this.frame.src = baseCourier + "/pudo-map?" + params;
    this.container.append(this.frame);
  }

  receive(event) {
    if (!this.dialog.open || !this.frame || event.origin !== baseCourier || event.source !== this.frame.contentWindow) return;
    const { type, value } = event.data || {};
    if (type === "POINT_INFO_CHANGE") {
      this.candidate = pointFromMap(value, this.carrier.value, this.carrier.selectedOptions[0].text);
      this.error.hidden = true;
    }
    if (type !== "SELECT_CHANGE") return;
    if (!this.candidate || String(value?.name) !== this.candidate.point_id) {
      this.error.hidden = false;
      return;
    }
    this.onChoose(this.candidate, polishPostalCode(this.postal.value));
    this.dialog.close();
  }
}
