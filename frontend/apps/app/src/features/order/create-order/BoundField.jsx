import OrderField from "./OrderField";

export default function BoundField({ editor, path, ...props }) {
  const value = path.split(".").reduce((node, key) => node?.[key], editor.form);
  return <OrderField {...props} value={value} disabled={editor.busy || Boolean(editor.card?.invoice)}
    onChange={(next) => editor.change(path, next)} />;
}
