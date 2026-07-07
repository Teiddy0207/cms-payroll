import { useState, useEffect } from 'react';
import { timekeepingAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { Badge } from '../components/Badge.jsx';
import { Pagination } from '../components/Pagination.jsx';
import { useToast } from '../hooks/useToast.js';

const PAGE_SIZE = 10;

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

	const [searchQuery, setSearchQuery] = useState('');
	const [statusFilter, setStatusFilter] = useState('ALL');
	const [filterDate, setFilterDate] = useState('');
	const [page, setPage] = useState(1);
	const [total, setTotal] = useState(0);

	const loadSheets = async () => {
		setLoading(true);
		try {
			const params = {
				page_number: page,
				page_size: PAGE_SIZE,
				period: period,
				search: searchQuery,
				status: statusFilter,
				date: filterDate
			};
			const res = await timekeepingAPI.getSheets(params);
			const d = res.data?.data;
			setSheets(d?.items || []);
			setTotal(d?.total_items || 0);
		} catch (err) {
			console.error("loadSheets error:", err);
			toast.error('Loi', 'Khong the tai bang cong');
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		setPage(1);
	}, [period, searchQuery, statusFilter, filterDate]);

	useEffect(() => {
		loadSheets();
	}, [page, period, searchQuery, statusFilter, filterDate]);

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
			toast.warning('Chu y', 'Vui long chon day du ngay bat dau va ngay ket thu');
			return;
		}

		setLoading(true);
		try {
			await timekeepingAPI.calculateSheets(calcDates);
			toast.success('Thanh cong', 'Da tinh toan xong bang cong ngay');
			setCalcModalOpen(false);
			loadSheets();
		} catch (err) {
			console.error("calculate error:", err);
			toast.error('Loi', 'Loi tinh toan bang cong');
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
				return <Badge variant="success">Duong gio</Badge>;
			case 'LATE':
				return <Badge variant="warning">Di muon</Badge>;
			case 'EARLY':
				return <Badge variant="warning">Ve som</Badge>;
			case 'LATE_EARLY':
				return <Badge variant="danger">Muon hoac ve som</Badge>;
			case 'ABSENT':
				return <Badge variant="danger">Vang mat</Badge>;
			default:
				return <Badge variant="secondary">{status}</Badge>;
		}
	};

	const countPresent = sheets.filter(s => s.status !== 'ABSENT').length;
	const countLate = sheets.filter(s => s.status === 'LATE' || s.status === 'LATE_EARLY').length;
	const countAbsent = sheets.filter(s => s.status === 'ABSENT').length;

	const columns = [
		{
			key: 'date',
			title: 'Ngay lam viec',
			render: (val) => formatDate(val)
		},
		{
			key: 'employee_code',
			title: 'Ma NV',
			render: (val) => <span style={{ fontWeight: 600 }}>{val}</span>
		},
		{
			key: 'full_name',
			title: 'Ho va ten',
			render: (val) => <span style={{ fontWeight: 600 }}>{val}</span>
		},
		{
			key: 'department_name',
			title: 'Phong ban'
		},
		{
			key: 'check_in',
			title: 'Gio vao',
			render: (val) => formatDateTime(val)
		},
		{
			key: 'check_out',
			title: 'Gio ra',
			render: (val) => formatDateTime(val)
		},
		{
			key: 'actual_work_day',
			title: 'Ngay cong',
			render: (val) => <span style={{ fontWeight: 700, color: val > 0 ? 'var(--success)' : 'var(--error)' }}>{val.toFixed(1)}</span>
		},
		{
			key: 'ot_hours',
			title: 'Gio OT',
			render: (val) => val > 0 ? <span style={{ color: 'var(--accent)', fontWeight: 700 }}>{val.toFixed(1)}h</span> : <span style={{ color: 'var(--text-muted)' }}>0</span>
		},
		{
			key: 'status',
			title: 'Trang thai',
			render: (val) => getStatusBadge(val)
		}
	];

	return (
		<div className="page-container">
			<div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
				<div>
					<h1 className="page-title">Bang cham cong tong hop</h1>
					<p className="page-subtitle" style={{ color: 'var(--text-muted)', fontSize: 14 }}>
						Theo doi thoi gian vao ra, so ngay cong thuc te va thoi gian lam them gio cua nhan su
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
						Tinh cong ngay
					</button>
				</div>
			</div>

			<div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 16, marginBottom: 20 }}>
				<div className="card" style={{ padding: '16px 20px', background: 'var(--bg-card)', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)' }}>
					<div style={{ fontSize: 12, color: 'var(--text-muted)', fontWeight: 600, textTransform: 'uppercase', marginBottom: 4 }}>Tong so ban ghi loc duoc</div>
					<div style={{ fontSize: 24, fontWeight: 800, color: 'var(--text-primary)' }}>{total}</div>
				</div>
				<div className="card" style={{ padding: '16px 20px', background: 'var(--bg-card)', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)' }}>
					<div style={{ fontSize: 12, color: 'var(--text-muted)', fontWeight: 600, textTransform: 'uppercase', marginBottom: 4 }}>Co mat (trang nay)</div>
					<div style={{ fontSize: 24, fontWeight: 800, color: 'var(--success)' }}>{countPresent}</div>
				</div>
				<div className="card" style={{ padding: '16px 20px', background: 'var(--bg-card)', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)' }}>
					<div style={{ fontSize: 12, color: 'var(--text-muted)', fontWeight: 600, textTransform: 'uppercase', marginBottom: 4 }}>Di muon / Ve som (trang nay)</div>
					<div style={{ fontSize: 24, fontWeight: 800, color: 'var(--warning)' }}>{countLate}</div>
				</div>
				<div className="card" style={{ padding: '16px 20px', background: 'var(--bg-card)', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)' }}>
					<div style={{ fontSize: 12, color: 'var(--text-muted)', fontWeight: 600, textTransform: 'uppercase', marginBottom: 4 }}>Vang mat (trang nay)</div>
					<div style={{ fontSize: 24, fontWeight: 800, color: 'var(--error)' }}>{countAbsent}</div>
				</div>
			</div>

			<div className="card" style={{ padding: 16, background: 'var(--bg-card)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-md)', marginBottom: 20 }}>
				<div style={{ display: 'flex', gap: 16, flexWrap: 'wrap' }}>
					<div style={{ flex: 1, minWidth: 200 }}>
						<label style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-secondary)', display: 'block', marginBottom: 6 }}>Tim kiem nhan vien</label>
						<input
							type="text"
							placeholder="Nhap ten hoac ma nhan vien..."
							className="form-input"
							value={searchQuery}
							onChange={e => setSearchQuery(e.target.value)}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-primary)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div style={{ width: 180 }}>
						<label style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-secondary)', display: 'block', marginBottom: 6 }}>Loc theo ngay cu the</label>
						<input
							type="date"
							className="form-input"
							value={filterDate}
							onChange={e => setFilterDate(e.target.value)}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-primary)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div style={{ width: 180 }}>
						<label style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-secondary)', display: 'block', marginBottom: 6 }}>Trang thai cong</label>
						<select
							className="form-input"
							value={statusFilter}
							onChange={e => setStatusFilter(e.target.value)}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-primary)', color: 'var(--text-primary)' }}
						>
							<option value="ALL">Tat ca trang thai</option>
							<option value="NORMAL">Duong gio</option>
							<option value="LATE">Di muon</option>
							<option value="EARLY">Ve som</option>
							<option value="LATE_EARLY">Muon hoac ve som</option>
							<option value="ABSENT">Vang mat</option>
						</select>
					</div>

					<div style={{ display: 'flex', alignItems: 'flex-end' }}>
						<button
							className="btn btn-secondary"
							onClick={() => {
								setSearchQuery('');
								setFilterDate('');
								setStatusFilter('ALL');
							}}
							style={{ padding: '9px 16px', borderRadius: 'var(--radius-md)' }}
						>
							Xoa bo loc
						</button>
					</div>
				</div>
			</div>

			<div className="card">
				<Table
					columns={columns}
					data={sheets}
					loading={loading}
					emptyMessage="Khong co du lieu bang cong phu hop"
				/>
				{total > PAGE_SIZE && (
					<div style={{ display: 'flex', justifyContent: 'center', marginTop: 16 }}>
						<Pagination
							page={page}
							totalItems={total}
							pageSize={PAGE_SIZE}
							onPageChange={setPage}
						/>
					</div>
				)}
			</div>

			<Modal
				isOpen={calcModalOpen}
				onClose={() => setCalcModalOpen(false)}
				title="Tinh cong ngay tu dong"
				size="sm"
			>
				<form onSubmit={handleRunCalculation} style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
					<p style={{ color: 'var(--text-muted)', fontSize: 13 }}>
						He thong se quet toan bo logs cham cong tho, don xin bu cong va don dang ky OT da duoc duyet de tong hop so ngay cong cua nhan su trong khoang thoi gian duoc chon.
					</p>

					<div className="form-group">
						<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Tu ngay <span style={{ color: 'var(--error)' }}>*</span></label>
						<input
							type="date"
							className="form-input"
							value={calcDates.start_date}
							onChange={e => setCalcDates(prev => ({ ...prev, start_date: e.target.value }))}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div className="form-group">
						<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Den ngay <span style={{ color: 'var(--error)' }}>*</span></label>
						<input
							type="date"
							className="form-input"
							value={calcDates.end_date}
							onChange={e => setCalcDates(prev => ({ ...prev, end_date: e.target.value }))}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 12 }}>
						<button type="button" className="btn btn-secondary" onClick={() => setCalcModalOpen(false)}>Huy</button>
						<button type="submit" className="btn btn-primary" disabled={loading}>
							{loading ? 'Dang tinh...' : 'Kich hoat tinh cong'}
						</button>
					</div>
				</form>
			</Modal>
		</div>
	);
}

export default TimesheetsPage;
