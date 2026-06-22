function DataTableOpenLink({ children, onOpen, target }) {
  if (!onOpen) {
    return <span className="dt-open-text">{children}</span>;
  }
  return (
    <button className="dt-open-link" type="button" onClick={(event) => handleOpen(event, onOpen, target)}>
      {children}
    </button>
  );
}

function handleOpen(event, onOpen, target) {
  event.stopPropagation();
  onOpen(target);
}

export default DataTableOpenLink;
