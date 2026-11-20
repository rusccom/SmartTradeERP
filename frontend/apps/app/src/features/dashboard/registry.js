import { createOrdersMenu } from "../order/menu";
import { catalogSettingsMenu } from "../products/settings/menu";

export function createMenuSections(t) {
  const catalog = menuGroup(t, "catalog", ["products", "customers", "warehouses"]);
  catalog.items.push(catalogSettingsMenu(t));
  return [
    catalog,
    createOrdersMenu(t),
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
