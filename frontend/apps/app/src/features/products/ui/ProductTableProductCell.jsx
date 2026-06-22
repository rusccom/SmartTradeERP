import ProductTableMetaList from "./ProductTableMetaList";

function ProductTableProductCell({ openLink, row, t, value }) {
  const name = value || row.name || "";
  return (
    <div className="product-table-product">
      {renderName(name, row, openLink)}
      <ProductTableMetaList row={row} t={t} />
    </div>
  );
}

function renderName(name, row, openLink) {
  if (openLink) {
    return openLink(name, row.product || row);
  }
  return <span className="product-table-name">{name}</span>;
}

export default ProductTableProductCell;
