export default function StatusDraftForm({ draft, inputRef }) {
  const { editor, duplicate } = draft;
  return (
    <form className="order-status-form" onSubmit={draft.apply}>
      <div className="order-status-field">
        <label htmlFor="order-status-name">Назва статусу</label>
        <input id="order-status-name" ref={inputRef} value={editor.name} required autoComplete="off"
          aria-invalid={duplicate} aria-describedby={duplicate ? "order-status-error" : undefined}
          onChange={(event) => draft.edit({ ...editor, name: event.target.value })} />
        {duplicate && <span id="order-status-error" role="alert">Статус із такою назвою вже є.</span>}
      </div>
      <div className="order-status-actions">
        <button className="order-status-primary" type="submit" disabled={!draft.canApply}>
          {editor.id ? "Застосувати" : "Додати статус"}
        </button>
        {editor.id && <button type="button" onClick={draft.reset}>Скасувати</button>}
      </div>
    </form>
  );
}
