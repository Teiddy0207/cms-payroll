import { useState, useEffect } from 'react';
import { timekeepingAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';

export function TimesheetsPage() {
	const toast = useToast();
	const [sheets, setSheets] = useState([]);
	const [loading, setLoading] = useState(false);
	const [period, setPeriod] = useState(() => {
		const now = new Date();
		return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
	});
	const [calcModalOpen, setCalcModalOpen] = useState(false);
	const [calcDates, setCalcDates] = useState({
		start_date: '',
		end_date: ''
	});

	const loadSheets = async () => {
		setLoading(true);
		try {
			const res = await timekeepingAPI.getSheets(period);
			setSheets(res.data?.data || []);
		} catch (err) {
			console.error("loadSheets error:", err);
			toast.error('Lỗi', 'Không thể tải bảng công của tháng');
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		loadSheets();
	}, [period]);

	const handleOpenCalc = () => {
		const [year, month] = period.split('-');
		const lastDay = new Date(year, month, 0).getDate();
		setCalcDates({
			start_date: `${year}-${month}-01`,
			end_date: `${year}-${month}-${String(lastDay).padStart(2, '0')}`
		});
		setCalcModalOpen(true);
	};

	const handleRunCalculation = async (e) => {
		e.preventDefault();
		if (!calcDates.start_date || !calcDates.end_date) {
			toast.warning('Chú ý', 'Vui lòng chọn đầy đủ ngày bắt đầu và ngày kết thúc');
			return;
		}

		setLoading(true);
		try {
			await timekeepingAPI.calculateSheets(calcDates);
			toast.success('Thành công', 'Đã tính toán xong bảng công ngày!');
			setCalcModalOpen(false);
			loadSheets();
		} catch (err) {
			console.error("calculate error:", err);
			toast.error('Lỗi', 'Lỗi tính toán bảng công');
		} finally {
			setLoading(false);
		}
	};

	const formatDateTime = (dateTimeStr) => {
		if (!dateTimeStr) return '—';
		const d = new Date(dateTimeStr);
		return d.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' });
	};

	const formatDate = (dateStr) => {
		if (!dateStr) return '—';
		const d = new Date(dateStr);
		return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`;
	};

	const getStatusBadge = (status) => {
		switch (status) {
			case 'NORMAL':
				return <Badge variant="success">Đúng giờ</Badge>;
			case 'LATE':
				return <Badge variant="warning">Đi muộn</Badge>;
			case 'EARLY':
				return <Badge variant="warning">Về sớm</Badge>;
			case 'LATE_EARLY':
				return <Badge variant="danger">Muộn/Sớm</Badge>;
			case 'ABSENT':
				return <Badge variant="danger">Vắng mặt</Badge>;
			default:
				return <Badge variant="secondary">{status}</Badge>;
		}
	};

	const columns = [
		{
			key: 'date',
			title: 'Ngày làm việc',
			render: (val) => formatDate(val)
		},
		{
			key: 'employee_code',
			title: 'Mã NV',
			render: (val) => <span style={{ fontWeight: 600 }}>{val}</span>
		},
		{
			key: 'full_name',
			title: 'Họ và tên',
			render: (val) => <span style={{ fontWeight: 600 }}>{val}</span>
		},
		{
			key: 'department_name',
			title: 'Phòng ban'
		},
		{
			key: 'check_in',
			title: 'Giờ vào',
			render: (val) => formatDateTime(val)
		},
		{
			key: 'check_out',
			title: 'Giờ ra',
			render: (val) => formatDateTime(val)
		},
		{
			key: 'actual_work_day',
			title: 'Ngày công',
			render: (val) => <span style={{ fontWeight: 700, color: val > 0 ? 'var(--success)' : 'var(--error)' }}>{val.toFixed(1)}</span>
		},
		{
			key: 'ot_hours',
			title: 'Giờ OT',
			render: (val) => val > 0 ? <span style={{ color: 'var(--accent)', fontWeight: 700 }}>{val.toFixed(1)}h</span> : <span style={{ color: 'var(--text-muted)' }}>0</span>
		},
		{
			key: 'status',
			title: 'Trạng thái',
			render: (val) => getStatusBadge(val)
		}
	];

	return (
		<div className="page-container">
			<div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
				<div>
					<h1 className="page-title">Bảng chấm công tổng hợp</h1>
					<p className="page-subtitle" style={{ color: 'var(--text-muted)', fontSize: 14 }}>
						Theo dõi thời gian vào/ra, số ngày công thực tế và số giờ làm thêm (OT) của nhân sự
					</p>
				</div>
				<div style={{ display: 'flex', gap: 12 }}>
					<input
						type="month"
						className="form-input"
						value={period}
						onChange={e => setPeriod(e.target.value)}
						style={{ padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
					/>
					<button className="btn btn-primary" onClick={handleOpenCalc}>
						⚡ Tính công ngày
					</button>
				</div>
			</div>

			<div className="card">
				<Table
					columns={columns}
					data={sheets}
					loading={loading}
					emptyMessage="Chưa có dữ liệu bảng công cho chu kỳ này. Vui lòng bấm Tính công ngày."
				/>
			</div>

			<Modal
				isOpen={calcModalOpen}
				onClose={() => setCalcModalOpen(false)}
				title="Tính công ngày tự động"
				size="sm"
			>
				<form onSubmit={handleRunCalculation} style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
					<p style={{ color: 'var(--text-muted)', fontSize: 13 }}>
						Hệ thống sẽ quét toàn bộ logs chấm công thô, đơn xin bù công và đơn đăng ký OT đã được duyệt để tổng hợp số ngày công của nhân sự trong khoảng thời gian được chọn.
					</p>

					<div className="form-group">
						<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Từ ngày <span style={{ color: 'var(--error)' }}>*</span></label>
						<input
							type="date"
							className="form-input"
							value={calcDates.start_date}
							onChange={e => setCalcDates(prev => ({ ...prev, start_date: e.target.value }))}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div className="form-group">
						<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Đến ngày <span style={{ color: 'var(--error)' }}>*</span></label>
						<input
							type="date"
							className="form-input"
							value={calcDates.end_date}
							onChange={e => setCalcDates(prev => ({ ...prev, end_date: e.target.value }))}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 12 }}>
						<button type="button" className="btn btn-secondary" onClick={() => setCalcModalOpen(false)}>Huỷ</button>
						<button type="submit" className="btn btn-primary" disabled={loading}>
							{loading ? '🔄 Đang tính...' : 'Kích hoạt tính công'}
						</button>
					</div>
				</form>
			</Modal>
		</div>
	);
}

export default TimesheetsPage;
