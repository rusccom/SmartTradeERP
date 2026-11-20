import { useEffect, useState } from "react";

import ProductTableCopyButton from "./ProductTableCopyButton";

function ProductTableMetaList({ row, t }) {
  const [copied, setCopied] = useState(null);
  const details = readDetails(row, t);
  useEffect(() => clearCopied(copied, setCopied), [copied]);
  if (details.length === 0) {
    return null;
  }
  return (
    <span className="product-table-meta">
      {details.map((item) => (
        <ProductTableCopyButton key={item.key} active={copied === item.key} item={item} onCopied={setCopied} t={t} />
      ))}
    </span>
  );
}

function readDetails(row, t) {
  return [
    row.sku_code ? { key: "sku", label: t("products.columns.sku"), value: row.sku_code } : null,
    row.barcode ? { key: "barcode", label: t("products.columns.barcode"), value: row.barcode } : null,
  ].filter(Boolean);
}

function clearCopied(copied, setCopied) {
  if (!copied) {
    return undefined;
  }
  const timer = window.setTimeout(() => setCopied(null), 1400);
  return () => window.clearTimeout(timer);
}

export default ProductTableMetaList;
