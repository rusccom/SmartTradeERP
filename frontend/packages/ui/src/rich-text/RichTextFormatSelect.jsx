import { setBlock } from "./richTextCommands";

const OPTIONS = [
  { value: "P", labelKey: "rte.format.paragraph" },
  { value: "H1", labelKey: "rte.format.h1" },
  { value: "H2", labelKey: "rte.format.h2" },
  { value: "H3", labelKey: "rte.format.h3" },
  { value: "BLOCKQUOTE", labelKey: "rte.format.blockquote" },
];

function RichTextFormatSelect({ editor, t }) {
  return (
    <select
      className="rte-select"
      aria-label={t("rte.format.label")}
      defaultValue="P"
      onMouseDown={(event) => event.stopPropagation()}
      onChange={(event) => setBlock(editor, event.target.value)}
    >
      {OPTIONS.map((option) => (
        <option key={option.value} value={option.value}>{t(option.labelKey)}</option>
      ))}
    </select>
  );
}

export default RichTextFormatSelect;
