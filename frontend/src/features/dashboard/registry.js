const CATALOG_PAGES = [
  { key: "products", path: "/dashboard/products", title: "Products" },
  { key: "customers", path: "/dashboard/customers", title: "Customers" },
  { key: "warehouses", path: "/dashboard/warehouses", title: "Warehouses" },
  { key: "bundles", path: "/dashboard/bundles", title: "Bundles" },
];

const DOCUMENT_PAGES = [
  { key: "documents", path: "/dashboard/documents", title: "Documents" },
];

const SYSTEM_PAGES = [
  { key: "reports", path: "/dashboard/reports", title: "Reports" },
  { key: "settings", path: "/dashboard/settings", title: "Settings" },
];

export const MENU_SECTIONS = [
  { label: "Catalog", items: CATALOG_PAGES },
  { label: "Documents", items: DOCUMENT_PAGES },
  { label: "System", items: SYSTEM_PAGES },
];
