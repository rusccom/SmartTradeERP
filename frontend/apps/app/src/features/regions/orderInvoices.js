import { orderPath } from "../order/menu";
import KsefDelivery from "./pl/KsefDelivery";

export const invoiceRegions = {
  "pl-fa3-vat": { Delivery: KsefDelivery, settings: `${orderPath}/settings/poland`, label: "entry.plSettings" },
};
