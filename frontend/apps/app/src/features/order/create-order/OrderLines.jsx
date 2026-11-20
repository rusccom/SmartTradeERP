import { useCallback, useState } from "react";
import { useI18n } from "@smarterp/i18n/useI18n";
import DocumentEditableTable from "@smarterp/ui/tables/document-editable-table/DocumentEditableTable";
import CatalogPicker from "./CatalogPicker";
import { lineColumns } from "./lineColumns";

export default function OrderLines({ editor, code }) {
  const { t, locale } = useI18n();
  const [picker, setPicker] = useState(false);
  const [selected, setSelected] = useState([]);
  const { items } = editor.form;
  const locked = editor.busy || Boolean(editor.card?.invoice);
  const pick = useCallback((line) => { editor.change("items", [...items, line]); setPicker(false); }, [editor.change, items]);
  const edit = (key, field, value) => editor.change("items", items.map((l) => l.key === key ? { ...l, [field]: value } : l));
  const remove = () => { editor.change("items", items.filter((l) => !selected.includes(l.key))); setSelected([]); };
  const columns = lineColumns({ t, locale, code, quote: editor.quote.data, edit, locked });
  const labels = { table: t("entry.items"), all: t("entry.allLines"), row: (l) => t("entry.selectLine", { name: l.product_name }), empty: t("entry.noLines") };
  return <section className="order-card order-lines"><div className="order-section-head"><h2>{labels.table}</h2>
    {!locked && <div className="order-actions">{selected.some((k) => items.some((i) => i.key === k)) && <button type="button" className="order-btn" onClick={remove}>{t("entry.remove")}</button>}
      <button type="button" className="order-btn" onClick={() => setPicker(true)}>{t("entry.addProduct")}</button></div>}</div>
    <DocumentEditableTable rows={items} columns={columns} labels={labels} selection={selected} onSelect={locked ? undefined : setSelected} />
    {editor.quote.error && <p role="alert" className="order-error">{editor.quote.error}</p>}
    {picker && <CatalogPicker onPick={pick} onClose={() => setPicker(false)} />}
  </section>;
}
