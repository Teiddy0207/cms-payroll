import { useState, useEffect } from 'react';
import { timekeepingAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';

// ============================================================
// LeaveRequestsPage — Quản lý nghỉ phép thường niên
// ============================================================
// Quy tắc nghiệp vụ:
//  • 12 ngày phép / năm, tích lũy 1 ngày mỗi đầu tháng
//  • Phép không dùng cộng dồn sang tháng sau, reset đầu năm
//  • Nghỉ trong số dư phép → LEAVE_PAID (không trừ lương)
//  • Nghỉ vượt số dư → LEAVE_UNPAID (mất lương ngày vượt)
// ============================================================

export function LeaveRequestsPage() {
	const toast = useToast();
	const [requests, setRequests] = useState([]);
	const [balance, setBalance] = useState(null);
	const [loading, setLoading] = useState(false);
	const [modalOpen, setModalOpen] = useState(false);
	const [accrueModalOpen, setAccrueModalOpen] = useState(false);
	const [activeTab, setActiveTab] = useState('requests'); // 'requests' | 'balance'
	const currentYear = new Date().getFullYear();

	const [formData, setFormData] = useState({
		start_date: '',
		end_date: '',
		days_requested: '',
		reason: ''
	});

	const [accrueForm, setAccrueForm] = useState({
		month: new Date().getMonth() + 1,
		year: currentYear
	});

	// ----- Load data -----
	const loadRequests = async () => {
		setLoading(true);
		try {
			const res = await timekeepingAPI.getLeaveRequests();
			setRequests(res.data?.data || []);
		} catch (err) {
			console.error('loadLeaveRequests error:', err);
			toast.error('Lỗi', 'Không thể tải danh sách đơn nghỉ phép');
		} finally {
			setLoading(false);
		}
	};

	const loadBalance = async () => {
		try {
			const res = await timekeepingAPI.getLeaveBalance(currentYear);
			setBalance(res.data?.data || null);
		} catch (err) {
			console.error('loadLeaveBalance error:', err);
		}
	};

	useEffect(() => {
		loadRequests();
		loadBalance();
	}, []);

	// ----- Tạo đơn nghỉ phép -----
	const handleOpenCreate = () => {
		const today = new Date().toISOString().split('T')[0];
		setFormData({ start_date: today, end_date: today, days_requested: '', reason: '' });
		setModalOpen(true);
	};

	const calcDays = (start, end, override) => {
		if (override !== '' && override !== undefined) return parseFloat(override) || 0;
		if (!start || !end) return 0;
		const s = new Date(start), e = new Date(end);
		if (e < s) return 0;
		return Math.floor((e - s) / 86400000) + 1;
	};

	const handleSubmit = async (ev) => {
		ev.preventDefault();
		if (!formData.start_date || !formData.end_date || !formData.reason.trim()) {
			toast.warning('Chú ý', 'Vui lòng nhập đầy đủ thông tin');
			return;
		}
		const days = calcDays(formData.start_date, formData.end_date, formData.days_requested);
		if (days <= 0) {
			toast.warning('Chú ý', 'Số ngày không hợp lệ');
			return;
		}
		setLoading(true);
		try {
			await timekeepingAPI.createLeaveRequest({
				start_date: new Date(formData.start_date).toISOString(),
				end_date: new Date(formData.end_date).toISOString(),
				days_requested: days,
				reason: formData.reason
			});
			toast.success('Thành công', 'Đã gửi đơn xin nghỉ phép');
			setModalOpen(false);
			loadRequests();
			loadBalance();
		} catch (err) {
			console.error('createLeaveRequest error:', err);
			toast.error('Lỗi', err?.response?.data?.message || 'Không thể gửi đơn nghỉ phép');
		} finally {
			setLoading(false);
		}
	};

	// ----- Duyệt / Từ chối -----
	const handleUpdateStatus = async (id, status) => {
		const action = status === 'APPROVED' ? 'duyệt' : 'từ chối';
		if (!window.confirm(`Bạn có chắc chắn muốn ${action} đơn này?`)) return;
		setLoading(true);
		try {
			await timekeepingAPI.updateLeaveRequestStatus(id, status);
			toast.success('Thành công', `Đã ${action} đơn nghỉ phép`);
			loadRequests();
			loadBalance();
		} catch (err) {
			console.error('updateLeaveStatus error:', err);
			toast.error('Lỗi', err?.response?.data?.message || 'Không có quyền hoặc lỗi phê duyệt');
		} finally {
			setLoading(false);
		}
	};

	// ----- Admin: Cộng phép đầu tháng -----
	const handleAccrue = async (ev) => {
		ev.preventDefault();
		setLoading(true);
		try {
			const res = await timekeepingAPI.accrueLeave(accrueForm.month, accrueForm.year);
			toast.success('Thành công', res.data?.data?.message || 'Đã cộng phép thành công');
			setAccrueModalOpen(false);
			loadBalance();
		} catch (err) {
			console.error('accrueLeave error:', err);
			toast.error('Lỗi', err?.response?.data?.message || 'Không thể cộng phép');
		} finally {
			setLoading(false);
		}
	};

	// ----- Helpers -----
	const getStatusBadge = (status) => {
		switch (status) {
			case 'APPROVED': return <Badge variant="success">Đã duyệt</Badge>;
			case 'REJECTED': return <Badge variant="danger">Từ chối</Badge>;
			case 'PENDING':  return <Badge variant="warning">Đang chờ</Badge>;
			default:         return <Badge variant="secondary">{status}</Badge>;
		}
	};

	const fmtDate = (str) => {
		if (!str) return '—';
		return str; // backend trả YYYY-MM-DD
	};

	const fmtDays = (n) => {
		if (n === 0.5) return '0.5 ngày';
		return `${n} ngày`;
	};

	// ---- Tính số ngày preview khi điền form ----
	const previewDays = formData.days_requested !== ''
		? parseFloat(formData.days_requested) || 0
		: calcDays(formData.start_date, formData.end_date, '');

	// ---- Columns table ----
	const columns = [
		{
			key: 'full_name',
			title: 'Nhân viên',
			render: (val, row) => (
				<span>
					<span style={{ fontWeight: 600 }}>{val}</span>
					{row.employee_code && <span style={{ color: 'var(--text-muted)', fontSize: 12, marginLeft: 6 }}>({row.employee_code})</span>}
				</span>
			)
		},
		{
			key: 'start_date',
			title: 'Từ ngày',
			render: (val) => fmtDate(val)
		},
		{
			key: 'end_date',
			title: 'Đến ngày',
			render: (val) => fmtDate(val)
		},
		{
			key: 'days_requested',
			title: 'Số ngày',
			render: (val) => (
				<span style={{ fontWeight: 600, color: 'var(--primary)' }}>{fmtDays(val)}</span>
			)
		},
		{
			key: 'reason',
			title: 'Lý do',
			render: (val) => <span style={{ color: 'var(--text-muted)' }}>{val || '—'}</span>
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
			render: (_, row) =>
				row.status === 'PENDING' ? (
					<div style={{ display: 'flex', gap: 8 }}>
						<button
							className="btn btn-sm btn-primary"
							onClick={() => handleUpdateStatus(row.id, 'APPROVED')}
						>
							Duyệt
						</button>
						<button
							className="btn btn-sm btn-danger"
							onClick={() => handleUpdateStatus(row.id, 'REJECTED')}
						>
							Từ chối
						</button>
					</div>
				) : (
					<span style={{ color: 'var(--text-muted)', fontSize: 13 }}>Đã xử lý</span>
				)
		}
	];

	// ---- Balance card style ----
	const balanceCardStyle = {
		display: 'flex',
		gap: 16,
		flexWrap: 'wrap',
		marginBottom: 24
	};
	const balanceItemStyle = (color) => ({
		flex: '1 1 140px',
		background: 'var(--bg-card)',
		border: `1px solid var(--border-color)`,
		borderRadius: 'var(--radius-lg)',
		padding: '18px 24px',
		display: 'flex',
		flexDirection: 'column',
		gap: 6,
		borderTop: `3px solid ${color}`
	});
	const bigNumStyle = (color) => ({
		fontSize: 32,
		fontWeight: 700,
		color
	});

	return (
		<div className="page-container">
			{/* Header */}
			<div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: 24, flexWrap: 'wrap', gap: 12 }}>
				<div>
					<h1 className="page-title">Nghỉ phép thường niên</h1>
					<p style={{ color: 'var(--text-muted)', fontSize: 14, marginTop: 4 }}>
						12 ngày phép/năm · Tích 1 ngày mỗi đầu tháng · Không dùng cộng dồn sang tháng sau
					</p>
				</div>
				<div style={{ display: 'flex', gap: 10 }}>
					<button className="btn btn-secondary" onClick={() => setAccrueModalOpen(true)}>
						<i className="fa-solid fa-calendar-plus" style={{ marginRight: 6 }} />
						Cộng phép tháng
					</button>
					<button className="btn btn-primary" onClick={handleOpenCreate}>
						<i className="fa-solid fa-plus" style={{ marginRight: 6 }} />
						Xin nghỉ phép
					</button>
				</div>
			</div>

			{/* Balance summary cards */}
			{balance && (
				<div style={balanceCardStyle}>
					<div style={balanceItemStyle('#10b981')}>
						<span style={{ fontSize: 13, color: 'var(--text-muted)', fontWeight: 500 }}>Phép còn lại</span>
						<span style={bigNumStyle('#10b981')}>{balance.balance}</span>
						<span style={{ fontSize: 12, color: 'var(--text-muted)' }}>ngày · Năm {balance.year}</span>
					</div>
					<div style={balanceItemStyle('#6366f1')}>
						<span style={{ fontSize: 13, color: 'var(--text-muted)', fontWeight: 500 }}>Đã tích lũy</span>
						<span style={bigNumStyle('#6366f1')}>{balance.accrued_days}</span>
						<span style={{ fontSize: 12, color: 'var(--text-muted)' }}>/ 12 ngày tối đa</span>
					</div>
					<div style={balanceItemStyle('#f59e0b')}>
						<span style={{ fontSize: 13, color: 'var(--text-muted)', fontWeight: 500 }}>Đã sử dụng</span>
						<span style={bigNumStyle('#f59e0b')}>{balance.used_days}</span>
						<span style={{ fontSize: 12, color: 'var(--text-muted)' }}>ngày trong năm</span>
					</div>
					<div style={{ ...balanceItemStyle('#64748b'), justifyContent: 'center' }}>
						{/* Progress bar */}
						<span style={{ fontSize: 13, color: 'var(--text-muted)', fontWeight: 500 }}>Tiến độ năm</span>
						<div style={{ background: 'var(--bg-hover)', borderRadius: 99, height: 8, overflow: 'hidden', marginTop: 8 }}>
							<div style={{
								height: '100%',
								width: `${Math.min(100, (balance.accrued_days / 12) * 100)}%`,
								background: 'linear-gradient(90deg, #10b981, #6366f1)',
								borderRadius: 99,
								transition: 'width 0.5s ease'
							}} />
						</div>
						<span style={{ fontSize: 12, color: 'var(--text-muted)', marginTop: 4 }}>
							Tháng {balance.last_accrual_month || 0} / 12
						</span>
					</div>
				</div>
			)}

			{/* Table */}
			<div className="card">
				<div style={{ padding: '16px 20px', borderBottom: '1px solid var(--border-color)' }}>
					<span style={{ fontWeight: 600, fontSize: 15 }}>
						<i className="fa-solid fa-list" style={{ marginRight: 8, color: 'var(--primary)' }} />
						Danh sách đơn nghỉ phép
					</span>
				</div>
				<Table
					columns={columns}
					data={requests}
					loading={loading}
					emptyMessage="Chưa có đơn xin nghỉ phép nào"
				/>
			</div>

			{/* Modal: Tạo đơn nghỉ phép */}
			<Modal
				isOpen={modalOpen}
				onClose={() => setModalOpen(false)}
				title="Xin nghỉ phép thường niên"
				size="sm"
			>
				<form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
					{/* Balance hint */}
					{balance && (
						<div style={{
							background: balance.balance > 0 ? 'rgba(16,185,129,0.08)' : 'rgba(239,68,68,0.08)',
							border: `1px solid ${balance.balance > 0 ? '#10b981' : '#ef4444'}`,
							borderRadius: 'var(--radius-md)',
							padding: '10px 14px',
							fontSize: 13,
							color: balance.balance > 0 ? '#10b981' : '#ef4444',
							fontWeight: 500
						}}>
							<i className={`fa-solid ${balance.balance > 0 ? 'fa-circle-check' : 'fa-triangle-exclamation'}`} style={{ marginRight: 8 }} />
							Số dư phép hiện tại: <strong>{balance.balance} ngày</strong>
							{balance.balance === 0 && ' — Nếu nghỉ sẽ mất lương ngày vượt'}
						</div>
					)}

					<div className="form-group">
						<label style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>
							Ngày bắt đầu <span style={{ color: 'var(--error)' }}>*</span>
						</label>
						<input
							type="date"
							className="form-input"
							value={formData.start_date}
							onChange={e => setFormData(p => ({ ...p, start_date: e.target.value }))}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div className="form-group">
						<label style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>
							Ngày kết thúc <span style={{ color: 'var(--error)' }}>*</span>
						</label>
						<input
							type="date"
							className="form-input"
							value={formData.end_date}
							min={formData.start_date}
							onChange={e => setFormData(p => ({ ...p, end_date: e.target.value }))}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div className="form-group">
						<label style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>
							Số ngày nghỉ
							<span style={{ fontWeight: 400, color: 'var(--text-muted)', marginLeft: 8, fontSize: 13 }}>
								(để trống = tự tính, nhập 0.5 nếu nghỉ nửa ngày)
							</span>
						</label>
						<input
							type="number"
							step="0.5"
							min="0.5"
							max="30"
							className="form-input"
							placeholder={`Tự tính: ${previewDays} ngày`}
							value={formData.days_requested}
							onChange={e => setFormData(p => ({ ...p, days_requested: e.target.value }))}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
						{previewDays > 0 && (
							<p style={{ marginTop: 4, fontSize: 12, color: balance && previewDays > balance.balance ? '#f59e0b' : 'var(--text-muted)' }}>
								{previewDays > 0 && balance && previewDays > balance.balance
									? `⚠ Vượt ${previewDays - balance.balance} ngày — phần vượt sẽ bị trừ lương`
									: `✓ ${previewDays} ngày trong phép`}
							</p>
						)}
					</div>

					<div className="form-group">
						<label style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>
							Lý do nghỉ phép <span style={{ color: 'var(--error)' }}>*</span>
						</label>
						<textarea
							className="form-input"
							rows="3"
							placeholder="Nhập lý do xin nghỉ phép..."
							value={formData.reason}
							onChange={e => setFormData(p => ({ ...p, reason: e.target.value }))}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)', resize: 'vertical' }}
						/>
					</div>

					<div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 4 }}>
						<button type="button" className="btn btn-secondary" onClick={() => setModalOpen(false)}>Huỷ</button>
						<button type="submit" className="btn btn-primary" disabled={loading}>
							{loading ? 'Đang gửi...' : 'Gửi đơn'}
						</button>
					</div>
				</form>
			</Modal>

			{/* Modal: Cộng phép đầu tháng (admin) */}
			<Modal
				isOpen={accrueModalOpen}
				onClose={() => setAccrueModalOpen(false)}
				title="Cộng phép đầu tháng"
				size="sm"
			>
				<form onSubmit={handleAccrue} style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
					<div style={{
						background: 'rgba(99,102,241,0.08)',
						border: '1px solid #6366f1',
						borderRadius: 'var(--radius-md)',
						padding: '12px 16px',
						fontSize: 13,
						color: '#6366f1'
					}}>
						<i className="fa-solid fa-info-circle" style={{ marginRight: 8 }} />
						Thao tác này cộng <strong>1 ngày phép</strong> cho tất cả nhân viên trong tháng đã chọn.
						Idempotent: mỗi tháng chỉ cộng 1 lần.
					</div>

					<div style={{ display: 'flex', gap: 12 }}>
						<div className="form-group" style={{ flex: 1 }}>
							<label style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Tháng</label>
							<select
								className="form-input"
								value={accrueForm.month}
								onChange={e => setAccrueForm(p => ({ ...p, month: parseInt(e.target.value) }))}
								style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
							>
								{Array.from({ length: 12 }, (_, i) => (
									<option key={i + 1} value={i + 1}>Tháng {i + 1}</option>
								))}
							</select>
						</div>
						<div className="form-group" style={{ flex: 1 }}>
							<label style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Năm</label>
							<input
								type="number"
								className="form-input"
								value={accrueForm.year}
								min={2020}
								max={2100}
								onChange={e => setAccrueForm(p => ({ ...p, year: parseInt(e.target.value) }))}
								style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
							/>
						</div>
					</div>

					<div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12 }}>
						<button type="button" className="btn btn-secondary" onClick={() => setAccrueModalOpen(false)}>Huỷ</button>
						<button type="submit" className="btn btn-primary" disabled={loading}>
							{loading ? 'Đang xử lý...' : 'Cộng phép'}
						</button>
					</div>
				</form>
			</Modal>
		</div>
	);
}

export default LeaveRequestsPage;
