import { apiPaths } from "../../../shared/api/publicApi";
import { createApiTablePreset } from "../../../shared/model/data-table/createApiTablePreset";

const capabilities = { sorting: true, search: true };

export function createProductsTablePreset(t) {
  return createApiTablePreset({
    id: "products",
    path: apiPaths.products,
    rowId: readProductId,
    columns: createColumns(t),
    capabilities,
    mapRows,
    mapStateToQuery: () => ({ include: "variants,stock" }),
  });
}

function createColumns(t) {
  return [
    { accessorKey: "name", header: t("products.columns.name") },
    { accessorKey: "sku_code", header: t("products.columns.sku"), enableSorting: false },
    { accessorKey: "barcode", header: t("products.columns.barcode"), enableSorting: false },
    { accessorKey: "unit", header: t("products.columns.unit"), enableSorting: false },
    { accessorKey: "price", header: t("products.columns.price"), enableSorting: false },
    { accessorKey: "global_qty", header: t("products.columns.quantity"), enableSorting: false },
    { accessorKey: "variant_count", header: t("products.columns.variants"), enableSorting: false },
  ];
}

function readProductId(row) {
  return row.id;
}

function mapRows(rows) {
  return rows.map((row) => ({ ...row, ...readProductDisplay(row) }));
}

function readProductDisplay(row) {
  const variants = readVariants(row);
  const single = variants.length === 1 ? variants[0] : null;
  return {
    sku_code: single?.sku_code || "",
    barcode: single?.barcode || "",
    unit: single?.unit || readSharedValue(variants, "unit"),
    price: single ? formatDecimal(single.price) : "",
    global_qty: formatDecimal(row.global_qty),
    variant_count: variants.length,
  };
}

function readSharedValue(variants, key) {
  const first = variants[0]?.[key] || "";
  return variants.every((item) => item[key] === first) ? first : "";
}

function readVariants(row) {
  return Array.isArray(row.variants) ? row.variants : [];
}

function formatDecimal(value) {
  if (value === undefined || value === null || value === "") return "";
  const number = Number(value);
  return Number.isFinite(number) ? String(number) : String(value);
}
