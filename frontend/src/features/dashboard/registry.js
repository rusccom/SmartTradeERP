export function createMenuSections(t) {
  return [
    {
      label: t("dashboard.menu.catalog"),
      items: [
        createItem("products", "/dashboard/products", t("dashboard.menu.products")),
        createItem("customers", "/dashboard/customers", t("dashboard.menu.customers")),
        createItem("warehouses", "/dashboard/warehouses", t("dashboard.menu.warehouses")),
        createItem("bundles", "/dashboard/bundles", t("dashboard.menu.bundles")),
      ],
    },
    {
      label: t("dashboard.menu.documents"),
      items: [createItem("documents", "/dashboard/documents", t("dashboard.menu.documents"))],
    },
    {
      label: t("dashboard.menu.system"),
      items: [
        createItem("reports", "/dashboard/reports", t("dashboard.menu.reports")),
        createItem("settings", "/dashboard/settings", t("dashboard.menu.settings")),
      ],
    },
  ];
}

function createItem(key, path, title) {
  return { key, path, title };
}
