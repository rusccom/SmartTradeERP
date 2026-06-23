import StarterKit from "@tiptap/starter-kit";
import Image from "@tiptap/extension-image";
import { Table } from "@tiptap/extension-table";
import TableRow from "@tiptap/extension-table-row";
import TableHeader from "@tiptap/extension-table-header";
import TableCell from "@tiptap/extension-table-cell";
import TextAlign from "@tiptap/extension-text-align";
import { TextStyle } from "@tiptap/extension-text-style";
import { Color } from "@tiptap/extension-color";
import Youtube from "@tiptap/extension-youtube";

const ALIGN_TYPES = ["heading", "paragraph"];

export function buildRichTextExtensions() {
  return [
    StarterKit.configure({ heading: { levels: [1, 2, 3] } }),
    Image.configure({ inline: false, allowBase64: false }),
    Table.configure({ resizable: true }),
    TableRow,
    TableHeader,
    TableCell,
    TextAlign.configure({ types: ALIGN_TYPES }),
    TextStyle,
    Color,
    Youtube.configure({ controls: true, nocookie: true, width: 640, height: 360 }),
  ];
}
