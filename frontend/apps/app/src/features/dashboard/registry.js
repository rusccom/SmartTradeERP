const sections = [
  { group: "catalog", pages: ["products", "customers", "warehouses"] },
  { group: "documents", pages: ["documents"] },
  { group: "onlineStore", pages: ["storefront", "integrations"] },
  { group: "system", pages: ["reports", "settings"] },
];

export function createMenuSections(t) {
  return sections.map(({ group, pages }) => ({
    label: t(`dashboard.menu.${group}`),
    items: pages.map((key) => ({
      key,
      path: `/dashboard/${key}`,
      title: t(`dashboard.menu.${key}`),
    })),
  }));
}
