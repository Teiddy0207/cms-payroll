import { useState, useEffect } from 'react';
import { timekeepingAPI, departmentsAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Pagination } from '../components/Pagination.jsx';
import { useToast } from '../hooks/useToast.js';

const PAGE_SIZE = 10;

const parseJwt = (token) => {
	try {
		const base64Url = token.split('.')[1];
		const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
		const pad = base64.length % 4;
		const paddedBase64 = pad ? base64 + '='.repeat(4 - pad) : base64;
		const jsonPayload = decodeURIComponent(window.atob(paddedBase64).split('').map(function(c) {
			return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
		}).join(''));
		return JSON.parse(jsonPayload);
	} catch (e) {
		console.error("JWT parse error:", e);
		return null;
	}
};

export function AttendanceLogsPage() {
	const toast = useToast();
	const [logs, setLogs] = useState([]);
	const [loading, setLoading] = useState(false);
	const [searchQuery, setSearchQuery] = useState('');
	const [filterDate, setFilterDate] = useState(() => {
		const now = new Date();
		return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
	});
	const [selectedDepartment, setSelectedDepartment] = useState('ALL');
	const [departments, setDepartments] = useState([]);
	const [isAdmin, setIsAdmin] = useState(false);
	const [page, setPage] = useState(1);
	const [total, setTotal] = useState(0);

	useEffect(() => {
		const token = localStorage.getItem('auth_token');
		if (token) {
			const payload = parseJwt(token);
			if (payload) {
				const adminCheck = (payload.user_name && payload.user_name.toLowerCase().includes('admin')) || 
				                   (payload.email && payload.email.toLowerCase().includes('admin')) || 
				                   payload.user_id === '00000000-0000-0000-0000-000000000000' || 
				                   payload.user_id === 'ea6133fe-f3c8-421d-801e-31264ee47b43';
				setIsAdmin(adminCheck);
			}
		}
	}, []);

	useEffect(() => {
		if (isAdmin) {
			const loadDepartments = async () => {
				try {
					const res = await departmentsAPI.list({ page_number: 1, page_size: 100 });
					setDepartments(res.data?.data?.items || []);
				} catch (err) {
					console.error("loadDepartments error:", err);
				}
			};
			loadDepartments();
		}
	}, [isAdmin]);

	const loadLogs = async () => {
		setLoading(true);
		try {
			const params = {
				page_number: page,
				page_size: PAGE_SIZE,
				search: searchQuery,
				date: filterDate,
				department_id: selectedDepartment === 'ALL' ? '' : selectedDepartment
			};
			const res = await timekeepingAPI.getLogs(params);
			const d = res.data?.data;
			setLogs(d?.items || []);
			setTotal(d?.total_items || 0);
		} catch (err) {
			console.error("loadLogs error:", err);
			toast.error('Loi', 'Khong the tai nhat ky cham cong');
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		setPage(1);
	}, [searchQuery, filterDate, selectedDepartment]);

	useEffect(() => {
		loadLogs();
	}, [page, searchQuery, filterDate, selectedDepartment]);

	const formatDateTime = (dateTimeStr) => {
		if (!dateTimeStr) return '—';
		const d = new Date(dateTimeStr);
		return d.toLocaleString('vi-VN', {
			day: '2-digit',
			month: '2-digit',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		});
	};

	const columns = [
		{
			key: 'timestamp',
			title: 'Thoi gian quet',
			render: (val) => <span style={{ fontWeight: 600 }}>{formatDateTime(val)}</span>
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
			key: 'location_gps',
			title: 'Toa do GPS'
		},
		{
			key: 'device_id',
			title: 'Thiet bi',
			render: (val) => <span style={{ color: 'var(--text-secondary)' }}>{val}</span>
		}
	];

	return (
		<div className="page-container">
			<div className="page-header" style={{ marginBottom: 20 }}>
				<div>
					<h1 className="page-title">Nhat ky cham cong tho</h1>
					<p className="page-subtitle" style={{ color: 'var(--text-muted)', fontSize: 14 }}>
						Xem lich su quet van tay hoac nhan dang khuon mat thoi gian thuc cua nhan su
					</p>
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
						<label style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-secondary)', display: 'block', marginBottom: 6 }}>Loc theo ngay</label>
						<input
							type="date"
							className="form-input"
							value={filterDate}
							onChange={e => setFilterDate(e.target.value)}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-primary)', color: 'var(--text-primary)' }}
						/>
					</div>

					{isAdmin && (
						<div style={{ width: 180 }}>
							<label style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-secondary)', display: 'block', marginBottom: 6 }}>Phong ban</label>
							<select
								className="form-input"
								value={selectedDepartment}
								onChange={e => setSelectedDepartment(e.target.value)}
								style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-primary)', color: 'var(--text-primary)' }}
							>
								<option value="ALL">Tat ca phong ban</option>
								{departments.map(dept => (
									<option key={dept.id} value={dept.id}>
										{dept.name}
									</option>
								))}
							</select>
						</div>
					)}

					<div style={{ display: 'flex', alignItems: 'flex-end' }}>
						<button
							className="btn btn-secondary"
							onClick={() => {
								setSearchQuery('');
								setFilterDate('');
								setSelectedDepartment('ALL');
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
					data={logs}
					loading={loading}
					emptyMessage="Khong co du lieu nhat ky cham cong phu hop"
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
		</div>
	);
}

export default AttendanceLogsPage;
