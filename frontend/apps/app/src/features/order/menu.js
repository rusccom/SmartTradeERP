export const orderPath = "/dashboard/orders";
export const orderSettingsSegment = "settings";

export const ordersMenu = {
  label: "Замовлення",
  lang: "uk",
  collapsible: true,
  items: [
    { key: "orders", path: orderPath, title: "Моє замовлення", end: true },
    { key: "orderSettings", path: `${orderPath}/${orderSettingsSegment}`, title: "Налаштування" },
  ],
};
