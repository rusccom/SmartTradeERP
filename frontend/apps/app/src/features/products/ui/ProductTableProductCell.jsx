import ProductTableMetaList from "./ProductTableMetaList";
import ProductTableThumb from "./ProductTableThumb";

function ProductTableProductCell({ openLink, row, t, value }) {
  const name = value || row.name || "";
  return (
    <div className="product-table-product">
      <ProductTableThumb url={row.image_thumb_url} name={name} />
      <div className="product-table-product-main">
        {renderName(name, row, openLink)}
        <ProductTableMetaList row={row} t={t} />
      </div>
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
