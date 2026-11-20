export const orderPath = "/dashboard/orders";
export const orderSettingsSegment = "settings";

export function createOrdersMenu(t) {
  return {
    label: t("orders.menu"),
    collapsible: true,
    items: [
      { key: "orders", path: orderPath, title: t("orders.list.title"), end: true },
      { key: "orderSettings", path: `${orderPath}/${orderSettingsSegment}`, title: t("dashboard.menu.settings") },
    ],
  };
}
