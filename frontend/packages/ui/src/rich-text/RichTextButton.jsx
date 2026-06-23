function RichTextButton({ active, disabled, label, onClick, children }) {
  return (
    <button
      type="button"
      className={active ? "rte-btn is-active" : "rte-btn"}
      title={label}
      aria-label={label}
      aria-pressed={Boolean(active)}
      disabled={disabled}
      onMouseDown={(event) => event.preventDefault()}
      onClick={onClick}
    >
      {children}
    </button>
  );
}

export default RichTextButton;
