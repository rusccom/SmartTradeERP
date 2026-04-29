import { useEffect, useRef } from "react";

function DataTableSelectionCheckbox({ checked, disabled, indeterminate, label, onChange }) {
  const ref = useRef(null);
  useEffect(() => {
    if (ref.current) {
      ref.current.indeterminate = Boolean(indeterminate) && !checked;
    }
  }, [checked, indeterminate]);
  return (
    <input
      ref={ref}
      aria-label={label}
      checked={checked}
      className="dt-select-checkbox"
      disabled={disabled}
      type="checkbox"
      onChange={onChange}
      onClick={(event) => event.stopPropagation()}
    />
  );
}

export default DataTableSelectionCheckbox;
