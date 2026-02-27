export default function DeleteConfirmModal({
  isOpen,
  title = "Confirm Deletion",
  message = "Are you sure you want to delete this item?",
  confirmLabel = "Delete",
  cancelLabel = "Cancel",
  isSubmitting = false,
  onConfirm = () => {},
  onCancel = () => {},
}) {
  if (!isOpen) {
    return null;
  }

  return (
    <div className="confirm-modal-overlay" role="presentation" onClick={isSubmitting ? undefined : onCancel}>
      <section
        className="confirm-modal"
        role="dialog"
        aria-modal="true"
        aria-label={title}
        onClick={(e) => e.stopPropagation()}
      >
        <h2 className="confirm-modal-title">{title}</h2>
        <p className="confirm-modal-message">{message}</p>
        <div className="confirm-modal-actions">
          <button className="confirm-modal-btn confirm-modal-btn-cancel" type="button" onClick={onCancel} disabled={isSubmitting}>
            {cancelLabel}
          </button>
          <button className="confirm-modal-btn confirm-modal-btn-danger" type="button" onClick={onConfirm} disabled={isSubmitting}>
            {isSubmitting ? "Deleting..." : confirmLabel}
          </button>
        </div>
      </section>
    </div>
  );
}
