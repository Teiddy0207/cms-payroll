import Modal from './Modal.jsx';

export function ConfirmDialog({ isOpen, onClose, onConfirm, title, message, confirmText = 'Xác nhận', loading = false }) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} title="" size="sm">
      <div className="confirm-dialog">
        <h3 className="confirm-title">{title || 'Xác nhận thao tác'}</h3>
        <p className="confirm-message">{message || 'Bạn có chắc chắn muốn thực hiện thao tác này?'}</p>
      </div>
      <div className="modal-footer" style={{ borderTop: 'none', paddingTop: 0 }}>
        <button className="btn btn-secondary" onClick={onClose} disabled={loading}>
          Hủy
        </button>
        <button className="btn btn-danger" onClick={onConfirm} disabled={loading}>
          {loading ? <><span className="spinner spinner-sm" /> Đang xử lý...</> : confirmText}
        </button>
      </div>
    </Modal>
  );
}

export default ConfirmDialog;
