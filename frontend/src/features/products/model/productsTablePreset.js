import { createElement } from "react";

import { apiPaths } from "../../../shared/api/publicApi";
import { createApiTablePreset } from "../../../shared/model/data-table/createApiTablePreset";
import ProductTableProductCell from "../ui/ProductTableProductCell";

const capabilities = { sorting: true, search: true };

export function createProductsTablePreset(t) {
  return createApiTablePreset({
    id: "products",
    path: apiPaths.products,
    rowId: readProductId,
    columns: createColumns(t),
    capabilities,
    search: { queryKey: "search" },
    mapRows,
    mapStateToQuery: () => ({ include: "variants,stock" }),
  });
}

function createColumns(t) {
  return [
    {
      accessorKey: "name",
      header: t("products.columns.name"),
      openOnClick: true,
      cell: (value, row, api) => createProductCell(t, value, row, api),
    },
    { accessorKey: "price_label", header: t("products.columns.price"), enableSorting: false, cell: readPriceCell },
    { accessorKey: "stock_label", header: t("products.columns.stock"), enableSorting: false, cell: readStockCell },
  ];
}

function createProductCell(t, value, row, api) {
  return createElement(ProductTableProductCell, { openLink: api.openLink, row, t, value });
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
  const unit = single?.unit || readSharedValue(variants, "unit");
  return {
    sku_code: single?.sku_code || "",
    barcode: single?.barcode || "",
    unit,
    price_label: readPriceLabel(variants),
    stock_label: formatStock(row.global_qty, unit),
    global_qty: formatDecimal(row.global_qty),
  };
}

function readPriceCell(value, row) {
  return value || formatDecimal(row.price);
}

function readStockCell(value, row) {
  return value || formatStock(row.global_qty, row.unit);
}

function readPriceLabel(variants) {
  const prices = variants.map((item) => Number(item.price)).filter(Number.isFinite);
  if (prices.length === 0) return "";
  const min = Math.min(...prices);
  const max = Math.max(...prices);
  return min === max ? formatDecimal(min) : `${formatDecimal(min)} - ${formatDecimal(max)}`;
}

function formatStock(quantity, unit) {
  const amount = formatDecimal(quantity);
  if (!amount) return "";
  return unit ? `${amount} ${unit}` : amount;
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
