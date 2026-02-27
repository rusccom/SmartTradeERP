/**
 * Dashboard page registry.
 * Adding a new page = adding one object here.
 */

const CATALOG_PAGES = [
  {
    key: "products",
    path: "/dashboard/products",
    title: "Products",
    columns: [
      { accessorKey: "name", header: "Name" },
      { accessorKey: "sku", header: "SKU" },
      { accessorKey: "price", header: "Price" },
    ],
    mock: [
      { id: 1, name: "Widget A", sku: "W-001", price: "$12.00" },
      { id: 2, name: "Widget B", sku: "W-002", price: "$18.50" },
      { id: 3, name: "Gadget X", sku: "G-010", price: "$45.00" },
    ],
  },
  {
    key: "groups",
    path: "/dashboard/groups",
    title: "Groups",
    columns: [
      { accessorKey: "name", header: "Name" },
      { accessorKey: "count", header: "Items" },
    ],
    mock: [
      { id: 1, name: "Electronics", count: 24 },
      { id: 2, name: "Accessories", count: 12 },
      { id: 3, name: "Packaging", count: 8 },
    ],
  },
  {
    key: "customers",
    path: "/dashboard/customers",
    title: "Customers",
    columns: [
      { accessorKey: "name", header: "Name" },
      { accessorKey: "email", header: "Email" },
      { accessorKey: "phone", header: "Phone" },
    ],
    mock: [
      { id: 1, name: "John Smith", email: "john@mail.com", phone: "+1-555-0101" },
      { id: 2, name: "Alice Brown", email: "alice@mail.com", phone: "+1-555-0202" },
      { id: 3, name: "Bob Wilson", email: "bob@mail.com", phone: "+1-555-0303" },
    ],
  },
];

const DOCUMENT_PAGES = [
  {
    key: "income",
    path: "/dashboard/docs/income",
    title: "Income",
    columns: [
      { accessorKey: "number", header: "Number" },
      { accessorKey: "date", header: "Date" },
      { accessorKey: "supplier", header: "Supplier" },
      { accessorKey: "total", header: "Total" },
    ],
    mock: [
      { id: 1, number: "INC-001", date: "2026-02-20", supplier: "Supplier A", total: "$1,200" },
      { id: 2, number: "INC-002", date: "2026-02-22", supplier: "Supplier B", total: "$850" },
    ],
  },
  {
    key: "expense",
    path: "/dashboard/docs/expense",
    title: "Expense",
    columns: [
      { accessorKey: "number", header: "Number" },
      { accessorKey: "date", header: "Date" },
      { accessorKey: "customer", header: "Customer" },
      { accessorKey: "total", header: "Total" },
    ],
    mock: [
      { id: 1, number: "EXP-001", date: "2026-02-21", customer: "John Smith", total: "$560" },
      { id: 2, number: "EXP-002", date: "2026-02-23", customer: "Alice Brown", total: "$320" },
    ],
  },
  {
    key: "transfer",
    path: "/dashboard/docs/transfer",
    title: "Transfer",
    columns: [
      { accessorKey: "number", header: "Number" },
      { accessorKey: "date", header: "Date" },
      { accessorKey: "from", header: "From" },
      { accessorKey: "to", header: "To" },
    ],
    mock: [
      { id: 1, number: "TRF-001", date: "2026-02-19", from: "Warehouse A", to: "Warehouse B" },
    ],
  },
  {
    key: "inventory",
    path: "/dashboard/docs/inventory",
    title: "Inventory",
    columns: [
      { accessorKey: "number", header: "Number" },
      { accessorKey: "date", header: "Date" },
      { accessorKey: "warehouse", header: "Warehouse" },
      { accessorKey: "status", header: "Status" },
    ],
    mock: [
      { id: 1, number: "INV-001", date: "2026-02-18", warehouse: "Main", status: "Completed" },
    ],
  },
  {
    key: "receipt",
    path: "/dashboard/docs/receipt",
    title: "Receipt",
    columns: [
      { accessorKey: "number", header: "Number" },
      { accessorKey: "date", header: "Date" },
      { accessorKey: "warehouse", header: "Warehouse" },
      { accessorKey: "total", header: "Total" },
    ],
    mock: [
      { id: 1, number: "REC-001", date: "2026-02-17", warehouse: "Main", total: "$2,400" },
    ],
  },
  {
    key: "writeoff",
    path: "/dashboard/docs/writeoff",
    title: "Write-off",
    columns: [
      { accessorKey: "number", header: "Number" },
      { accessorKey: "date", header: "Date" },
      { accessorKey: "reason", header: "Reason" },
      { accessorKey: "total", header: "Total" },
    ],
    mock: [
      { id: 1, number: "WO-001", date: "2026-02-16", reason: "Damaged", total: "$150" },
    ],
  },
];

export const MENU_SECTIONS = [
  { label: "Catalog", items: CATALOG_PAGES },
  { label: "Documents", items: DOCUMENT_PAGES },
];

export const ALL_PAGES = [...CATALOG_PAGES, ...DOCUMENT_PAGES];

export function findPageByKey(key) {
  return ALL_PAGES.find((p) => p.key === key);
}
