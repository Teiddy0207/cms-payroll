import { useState, useEffect } from 'react';
import { timekeepingAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';

export function OTRequestsPage() {
	const toast = useToast();
	const [requests, setRequests] = useState([]);
	const [loading, setLoading] = useState(false);
	const [modalOpen, setModalOpen] = useState(false);
	const [formData, setFormData] = useState({
		date: '',
		hours_requested: 2,
		is_night_ot: false,
		is_holiday_ot: false
	});

	const loadRequests = async () => {
		setLoading(true);
		try {
			const res = await timekeepingAPI.getOTRequests();
			setRequests(res.data?.data || []);
		} catch (err) {
			console.error("loadRequests error:", err);
			toast.error('Lỗi', 'Không thể tải danh sách đăng ký OT');
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
			hours_requested: 2,
			is_night_ot: false,
			is_holiday_ot: false
		});
		setModalOpen(true);
	};

	const handleSubmit = async (e) => {
		e.preventDefault();
		if (!formData.date || formData.hours_requested <= 0) {
			toast.warning('Chú ý', 'Vui lòng nhập đầy đủ thông tin hợp lệ');
			return;
		}

		setLoading(true);
		try {
			await timekeepingAPI.createOTRequest({
				date: new Date(formData.date).toISOString(),
				hours_requested: parseFloat(formData.hours_requested),
				is_night_ot: formData.is_night_ot,
				is_holiday_ot: formData.is_holiday_ot
			});
			toast.success('Thành công', 'Đã gửi đơn đăng ký làm thêm giờ');
			setModalOpen(false);
			loadRequests();
		} catch (err) {
			console.error("create error:", err);
			toast.error('Lỗi', 'Không thể gửi đơn đăng ký OT');
		} finally {
			setLoading(false);
		}
	};

	const handleApprove = async (id, status) => {
		if (!window.confirm(`Bạn có chắc chắn muốn ${status === 'APPROVED' ? 'duyệt' : 'từ chối'} đơn này?`)) return;
		setLoading(true);
		try {
			await timekeepingAPI.updateOTRequestStatus(id, status);
			toast.success('Thành công', `Đã ${status === 'APPROVED' ? 'duyệt' : 'từ chối'} đơn đăng ký OT`);
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
			title: 'Ngày làm thêm',
			render: (val) => formatDate(val)
		},
		{
			key: 'hours_requested',
			title: 'Số giờ đăng ký',
			render: (val) => <span style={{ fontWeight: 700, color: 'var(--accent)' }}>{val.toFixed(1)}h</span>
		},
		{
			key: 'is_night_ot',
			title: 'Làm đêm',
			render: (val) => val ? <Badge variant="danger">Ca đêm</Badge> : <span style={{ color: 'var(--text-muted)' }}>Không</span>
		},
		{
			key: 'is_holiday_ot',
			title: 'Ngày nghỉ/lễ',
			render: (val) => val ? <Badge variant="success">Ngày nghỉ</Badge> : <span style={{ color: 'var(--text-muted)' }}>Không</span>
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
					<h1 className="page-title">Đăng ký làm thêm giờ (OT)</h1>
					<p className="page-subtitle" style={{ color: 'var(--text-muted)', fontSize: 14 }}>
						Gửi đơn đăng ký số giờ làm thêm và duyệt đơn làm thêm của nhân viên cấp dưới
					</p>
				</div>
				<button className="btn btn-primary" onClick={handleOpenCreate}>
					Dang ky lam them
				</button>
			</div>

			<div className="card">
				<Table
					columns={columns}
					data={requests}
					loading={loading}
					emptyMessage="Chưa có đơn đăng ký làm thêm nào được tạo"
				/>
			</div>

			<Modal
				isOpen={modalOpen}
				onClose={() => setModalOpen(false)}
				title="Đăng ký làm thêm giờ (OT)"
				size="sm"
			>
				<form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
					<div className="form-group">
						<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Ngày làm thêm <span style={{ color: 'var(--error)' }}>*</span></label>
						<input
							type="date"
							className="form-input"
							value={formData.date}
							onChange={e => setFormData(prev => ({ ...prev, date: e.target.value }))}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div className="form-group">
						<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Số giờ đăng ký làm thêm <span style={{ color: 'var(--error)' }}>*</span></label>
						<input
							type="number"
							className="form-input"
							step="0.5"
							min="0.5"
							max="12"
							value={formData.hours_requested}
							onChange={e => setFormData(prev => ({ ...prev, hours_requested: e.target.value }))}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div style={{ display: 'flex', gap: 20 }}>
						<label style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 14, fontWeight: 500, cursor: 'pointer' }}>
							<input
								type="checkbox"
								checked={formData.is_night_ot}
								onChange={e => setFormData(prev => ({ ...prev, is_night_ot: e.target.checked }))}
								style={{ width: 16, height: 16 }}
							/>
							Làm đêm (Ca đêm)
						</label>
						<label style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 14, fontWeight: 500, cursor: 'pointer' }}>
							<input
								type="checkbox"
								checked={formData.is_holiday_ot}
								onChange={e => setFormData(prev => ({ ...prev, is_holiday_ot: e.target.checked }))}
								style={{ width: 16, height: 16 }}
							/>
							Làm ngày nghỉ/lễ
						</label>
					</div>

					<div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 12 }}>
						<button type="button" className="btn btn-secondary" onClick={() => setModalOpen(false)}>Huỷ</button>
						<button type="submit" className="btn btn-primary">Gửi đơn đăng ký</button>
					</div>
				</form>
			</Modal>
		</div>
	);
}

export default OTRequestsPage;
