function RichTextHtmlView({ value, onChange, t }) {
  return (
    <div className="rte-html">
      <textarea
        className="rte-html-area"
        spellCheck={false}
        aria-label={t("rte.html.label")}
        value={value}
        onChange={(event) => onChange(event.target.value)}
      />
    </div>
  );
}

export default RichTextHtmlView;
