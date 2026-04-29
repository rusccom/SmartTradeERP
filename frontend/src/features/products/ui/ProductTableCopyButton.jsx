function ProductTableCopyButton({ active, item, onCopied, t }) {
  return (
    <span className="product-table-copy-wrap">
      <button className="product-table-meta-copy" type="button" onClick={(event) => copyMeta(event, item, onCopied)}>
        {item.label}: {item.value}
      </button>
      {active && <span className="product-table-copy-popup">{readCopyLabel(item, t)}</span>}
    </span>
  );
}

async function copyMeta(event, item, onCopied) {
  event.stopPropagation();
  await copyText(item.value);
  onCopied(item.key);
}

async function copyText(value) {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(value);
    return;
  }
  window.prompt("", value);
}

function readCopyLabel(item, t) {
  const key = item.key === "sku" ? "products.copy.sku" : "products.copy.barcode";
  return t(key);
}

export default ProductTableCopyButton;
