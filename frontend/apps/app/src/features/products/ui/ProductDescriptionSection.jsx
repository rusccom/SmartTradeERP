import RichTextEditor from "@smarterp/ui/rich-text/RichTextEditor";

import { useDescriptionImageUpload } from "./useDescriptionImageUpload";

function ProductDescriptionSection({ form, onChange, productId, t }) {
  const requestImage = useDescriptionImageUpload(productId);
  return (
    <section className="product-card">
      <h3 className="product-card__title">{t("products.form.sections.description")}</h3>
      <RichTextEditor
        initialContent={form.descriptionHtml}
        documentKey={productId || "new"}
        imageDisabled={!productId}
        onRequestImage={requestImage}
        onHtmlChange={(value) => onChange({ target: { name: "descriptionHtml", value } })}
        t={t}
      />
    </section>
  );
}

export default ProductDescriptionSection;
