function ProductTableProductCell({ row, value }) {
  const name = value || row.name || "";
  const details = readDetails(row);
  return (
    <div className="product-table-product">
      <span className="product-table-name">{name}</span>
      {details.length > 0 && (
        <span className="product-table-meta">
          {details.map((item) => <span key={item} className="product-table-meta-item">{item}</span>)}
        </span>
      )}
    </div>
  );
}

function readDetails(row) {
  return [
    row.sku_code ? `SKU: ${row.sku_code}` : "",
    row.barcode ? `Barcode: ${row.barcode}` : "",
  ].filter(Boolean);
}

export default ProductTableProductCell;
