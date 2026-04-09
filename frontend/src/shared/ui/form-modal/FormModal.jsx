import Modal from "../modal/Modal";
import FormSection from "./FormSection";
import "./form-modal.css";

function FormModal(props) {
  return (
    <Modal open={props.open} title={props.title} description={props.description} closeLabel={props.closeLabel} onClose={props.onClose}>
      <form className="shared-form-modal" onSubmit={props.onSubmit}>
        {props.topSlot}
        {props.sections.map((section) => <FormSection key={section.id} form={props.form} onChange={props.onChange} section={section} t={props.t} />)}
        {props.bottomSlot}
        {props.error && <p className="shared-form-modal-error">{props.error}</p>}
        <div className="shared-form-modal-actions">
          <button className="shared-form-modal-secondary" type="button" onClick={props.onClose}>{props.cancelLabel}</button>
          <button className="shared-form-modal-primary" type="submit" disabled={props.isSubmitting}>{props.isSubmitting ? props.submittingLabel : props.submitLabel}</button>
        </div>
      </form>
    </Modal>
  );
}

export default FormModal;
