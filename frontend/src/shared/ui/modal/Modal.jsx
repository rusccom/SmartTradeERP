import { useEffect, useId } from "react";

import "./modal.css";

function Modal({
  children,
  closeLabel,
  closeOnBackdrop = false,
  closeOnEscape = false,
  description,
  onClose,
  open,
  title,
}) {
  const titleId = useId();
  const descriptionId = useId();
  useEscapeClose({ closeOnEscape, onClose, open });
  useBodyScrollLock(open);
  if (!open) return null;
  return (
    <div className="ui-modal-backdrop" role="presentation" onMouseDown={createBackdropHandler({ closeOnBackdrop, onClose })}>
      <section className="ui-modal-surface" role="dialog" aria-modal="true" aria-labelledby={titleId} aria-describedby={readDescriptionId(description, descriptionId)}>
        <header className="ui-modal-header">
          <div className="ui-modal-copy">
            <h2 id={titleId} className="ui-modal-title">{title}</h2>
            {description && <p id={descriptionId} className="ui-modal-description">{description}</p>}
          </div>
          <button className="ui-modal-close" type="button" aria-label={closeLabel} onClick={onClose}>&times;</button>
        </header>
        <div className="ui-modal-body">{children}</div>
      </section>
    </div>
  );
}

function useEscapeClose({ closeOnEscape, onClose, open }) {
  useEffect(() => {
    if (!open || !closeOnEscape) return undefined;
    const handleKeyDown = (event) => event.key === "Escape" && onClose();
    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [closeOnEscape, onClose, open]);
}

function useBodyScrollLock(open) {
  useEffect(() => {
    if (!open) return undefined;
    const previous = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => { document.body.style.overflow = previous; };
  }, [open]);
}

function createBackdropHandler({ closeOnBackdrop, onClose }) {
  return (event) => closeOnBackdrop && event.target === event.currentTarget && onClose();
}

function readDescriptionId(description, descriptionId) {
  return description ? descriptionId : undefined;
}

export default Modal;
