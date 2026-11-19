import { ordersMenu } from "../orders/menu";

export function createMenuSections(t) {
  return [
    menuGroup(t, "catalog", ["products", "customers", "warehouses"]),
    ordersMenu,
    menuGroup(t, "documents", ["documents"]),
    menuGroup(t, "onlineStore", ["storefront", "integrations"]),
    menuGroup(t, "system", ["reports", "settings"]),
  ];
}

function menuGroup(t, group, pages) {
  return {
    label: t(`dashboard.menu.${group}`),
    items: pages.map((key) => ({
      key,
      path: `/dashboard/${key}`,
      title: t(`dashboard.menu.${key}`),
    })),
  };
}
