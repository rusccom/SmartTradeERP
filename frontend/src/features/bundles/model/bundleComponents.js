let draftID = 0;

export function addComponentRow(rows, options) {
  return [...rows, createComponentRow(options)];
}

export function createComponentRows(items) {
  return (items || []).map((item) => ({
    id: nextDraftID(),
    componentVariantID: item.component_variant_id,
    qty: String(item.qty || "1"),
    label: componentLabel(item),
    unit: item.unit || "",
  }));
}

export function createInitialRows(items) {
  const rows = createComponentRows(items);
  return rows.length ? rows : [createComponentRow([])];
}

export function patchComponentRow(rows, rowID, event, options) {
  const { name, value } = event.target;
  return rows.map((row) => patchRow(row, rowID, name, value, options));
}

export function removeComponentRow(rows, rowID) {
  if (rows.length <= 1) return rows;
  return rows.filter((row) => row.id !== rowID);
}

export function toComponentPayload(rows) {
  return rows.map((row) => ({
    component_variant_id: row.componentVariantID,
    qty: Number(row.qty) || 0,
  }));
}

function createComponentRow(options) {
  const option = options[0] || {};
  return {
    id: nextDraftID(),
    componentVariantID: option.id || "",
    qty: "1",
    label: option.label || "",
    unit: option.unit || "",
  };
}

function patchRow(row, rowID, name, value, options) {
  if (row.id !== rowID) return row;
  if (name !== "componentVariantID") return { ...row, [name]: value };
  const option = options.find((item) => item.id === value);
  return { ...row, componentVariantID: value, label: option?.label || "", unit: option?.unit || "" };
}

function componentLabel(item) {
  if (!item.product_name) return item.component_variant_id || "";
  if (!item.variant_name || item.variant_name === "Default") return item.product_name;
  return `${item.product_name} / ${item.variant_name}`;
}

function nextDraftID() {
  draftID += 1;
  return `component-${draftID}`;
}
