import { useState, useEffect } from 'react';
import { timekeepingAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';

export function ExplanationRequestsPage() {
	const toast = useToast();
	const [requests, setRequests] = useState([]);
	const [loading, setLoading] = useState(false);
	const [modalOpen, setModalOpen] = useState(false);
	const [formData, setFormData] = useState({
		date: '',
		reason: ''
	});

	const loadRequests = async () => {
		setLoading(true);
		try {
			const res = await timekeepingAPI.getExplanations();
			setRequests(res.data?.data || []);
		} catch (err) {
			console.error("loadRequests error:", err);
			toast.error('Lỗi', 'Không thể tải danh sách đơn giải trình');
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		loadRequests();
	}, []);

	const handleOpenCreate = () => {
		setFormData({
			date: new Date().toISOString().split('T')[0],
			reason: ''
		});
		setModalOpen(true);
	};

	const handleSubmit = async (e) => {
		e.preventDefault();
		if (!formData.date || !formData.reason) {
			toast.warning('Chú ý', 'Vui lòng nhập đầy đủ thông tin');
			return;
		}

		setLoading(true);
		try {
			await timekeepingAPI.createExplanation({
				date: new Date(formData.date).toISOString(),
				reason: formData.reason
			});
			toast.success('Thành công', 'Đã gửi đơn giải trình bù công');
			setModalOpen(false);
			loadRequests();
		} catch (err) {
			console.error("create error:", err);
			toast.error('Lỗi', 'Không thể gửi đơn giải trình');
		} finally {
			setLoading(false);
		}
	};

	const handleApprove = async (id, status) => {
		if (!window.confirm(`Bạn có chắc chắn muốn ${status === 'APPROVED' ? 'duyệt' : 'từ chối'} đơn này?`)) return;
		setLoading(true);
		try {
			await timekeepingAPI.updateExplanationStatus(id, status);
			toast.success('Thành công', `Đã ${status === 'APPROVED' ? 'duyệt' : 'từ chối'} đơn giải trình`);
			loadRequests();
		} catch (err) {
			console.error("approve error:", err);
			toast.error('Lỗi', 'Không có quyền hoặc lỗi phê duyệt đơn');
		} finally {
			setLoading(false);
		}
	};

	const getStatusBadge = (status) => {
		switch (status) {
			case 'APPROVED':
				return <Badge variant="success">Đã duyệt</Badge>;
			case 'REJECTED':
				return <Badge variant="danger">Từ chối</Badge>;
			case 'PENDING':
				return <Badge variant="warning">Đang chờ</Badge>;
			default:
				return <Badge variant="secondary">{status}</Badge>;
		}
	};

	const formatDate = (dateStr) => {
		if (!dateStr) return '—';
		const d = new Date(dateStr);
		return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`;
	};

	const columns = [
		{
			key: 'fullName',
			title: 'Họ và tên',
			render: (val, row) => <span style={{ fontWeight: 600 }}>{val} ({row.employee_code})</span>
		},
		{
			key: 'date',
			title: 'Ngày xin bù công',
			render: (val) => formatDate(val)
		},
		{
			key: 'reason',
			title: 'Lý do giải trình'
		},
		{
			key: 'status',
			title: 'Trạng thái',
			render: (val) => getStatusBadge(val)
		},
		{
			key: 'approved_name',
			title: 'Người duyệt',
			render: (val) => val || <span style={{ color: 'var(--text-muted)' }}>—</span>
		},
		{
			key: 'created_at',
			title: 'Ngày gửi',
			render: (val) => formatDate(val)
		},
		{
			key: 'actions',
			title: 'Thao tác',
			render: (_, row) => (
				row.status === 'PENDING' ? (
					<div style={{ display: 'flex', gap: 8 }}>
						<button className="btn btn-sm btn-primary" onClick={() => handleApprove(row.id, 'APPROVED')}>Duyệt</button>
						<button className="btn btn-sm btn-danger" onClick={() => handleApprove(row.id, 'REJECTED')}>Từ chối</button>
					</div>
				) : <span style={{ color: 'var(--text-muted)' }}>Đã xử lý</span>
			)
		}
	];

	return (
		<div className="page-container">
			<div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
				<div>
					<h1 className="page-title">Giải trình bù công</h1>
					<p className="page-subtitle" style={{ color: 'var(--text-muted)', fontSize: 14 }}>
						Gửi giải trình cho những ngày đi muộn/về sớm và duyệt đơn của nhân viên cấp dưới
					</p>
				</div>
				<button className="btn btn-primary" onClick={handleOpenCreate}>
					➕ Tạo đơn giải trình
				</button>
			</div>

			<div className="card">
				<Table
					columns={columns}
					data={requests}
					loading={loading}
					emptyMessage="Chưa có đơn giải trình nào được tạo"
				/>
			</div>

			<Modal
				isOpen={modalOpen}
				onClose={() => setModalOpen(false)}
				title="Tạo đơn giải trình bù công"
				size="sm"
			>
				<form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
					<div className="form-group">
						<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Ngày cần giải trình <span style={{ color: 'var(--error)' }}>*</span></label>
						<input
							type="date"
							className="form-input"
							value={formData.date}
							onChange={e => setFormData(prev => ({ ...prev, date: e.target.value }))}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div className="form-group">
						<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Lý do giải trình <span style={{ color: 'var(--error)' }}>*</span></label>
						<textarea
							className="form-input"
							rows="3"
							placeholder="Nhập lý do đi muộn, về sớm hoặc quên quẹt thẻ..."
							value={formData.reason}
							onChange={e => setFormData(prev => ({ ...prev, reason: e.target.value }))}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 12 }}>
						<button type="button" className="btn btn-secondary" onClick={() => setModalOpen(false)}>Huỷ</button>
						<button type="submit" className="btn btn-primary">Gửi đơn</button>
					</div>
				</form>
			</Modal>
		</div>
	);
}

export default ExplanationRequestsPage;
