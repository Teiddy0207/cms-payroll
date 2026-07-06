import { useState, useEffect, useCallback, useRef } from 'react';
import { calculatorAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';

export function PayrollRunPage() {
	const toast = useToast();
	const [loading, setLoading] = useState(false);
	const [records, setRecords] = useState([]);
	const [selectedPeriod, setSelectedPeriod] = useState(() => {
		const d = new Date();
		return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
	});

	const [jobStatus, setJobStatus] = useState(null);
	const [polling, setPolling] = useState(false);
	const pollInterval = useRef(null);

	const [selectedRow, setSelectedRow] = useState(null);
	const [detailOpen, setDetailOpen] = useState(false);

	const formatVND = (value) => {
		return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(value || 0);
	};

	const fetchRecords = useCallback(async (period) => {
		setLoading(true);
		try {
			const res = await calculatorAPI.getSavedRecords(period);
			setRecords(res.data?.data || []);
		} catch (err) {
			console.error("fetchRecords error:", err);
			toast.error('Lỗi', 'Không thể tải danh sách bảng lương đã lưu');
		} finally {
			setLoading(false);
		}
	}, [toast]);

	const checkJobStatus = useCallback(async (jobId) => {
		if (!jobId) {
			console.warn("checkJobStatus: jobId is empty");
			return;
		}
		try {
			const res = await calculatorAPI.getJobStatus(jobId);
			const status = res.data?.data;
			setJobStatus(status);

			if (status?.status === 'SUCCESS') {
				setPolling(false);
				if (pollInterval.current) {
					clearInterval(pollInterval.current);
					pollInterval.current = null;
				}
				toast.success('Thành công', 'Đã hoàn tất tính toán lương cả công ty');
				fetchRecords(selectedPeriod);
			} else if (status?.status === 'FAILED') {
				setPolling(false);
				if (pollInterval.current) {
					clearInterval(pollInterval.current);
					pollInterval.current = null;
				}
				toast.error('Thất bại', status?.error || 'Có lỗi xảy ra trong quá trình tính toán');
			}
		} catch (err) {
			console.error("checkJobStatus error:", err);
		}
	}, [fetchRecords, selectedPeriod, toast]);

	const startPolling = useCallback((jobId) => {
		if (!jobId) {
			console.warn("startPolling: cannot start without jobId");
			return;
		}
		setPolling(true);
		if (pollInterval.current) {
			clearInterval(pollInterval.current);
		}
		pollInterval.current = setInterval(() => {
			checkJobStatus(jobId);
		}, 1500);
	}, [checkJobStatus]);

	const handleCalculate = async () => {
		setLoading(true);
		try {
			const res = await calculatorAPI.runCalculation({ period: selectedPeriod });
			const jobId = res.data?.job_id;
			if (jobId) {
				toast.success('Khởi chạy', 'Hệ thống đang tiến hành tính toán lương ở nền');
				startPolling(jobId);
			} else {
				toast.error('Lỗi', 'Không nhận được mã công việc từ máy chủ');
			}
		} catch (err) {
			console.error("handleCalculate error:", err);
			if (err.response?.status === 409) {
				toast.warning('Bận', 'Hệ thống đang chạy một tiến trình tính lương khác cho kỳ này');
			} else {
				toast.error('Lỗi', 'Không thể bắt đầu tính toán lương');
			}
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		fetchRecords(selectedPeriod);
	}, [selectedPeriod, fetchRecords]);

	useEffect(() => {
		return () => {
			if (pollInterval.current) {
				clearInterval(pollInterval.current);
			}
		};
	}, []);

	const columns = [
		{
			key: 'fullName',
			title: 'Nhân sự',
			render: (_, row) => (
				<div>
					<div style={{ fontWeight: 600, color: 'var(--text-primary)' }}>{row.fullName}</div>
					<div style={{ fontSize: '12px', color: 'var(--text-muted)' }}>Mã NV: {row.employee_id.substring(0, 8)}</div>
				</div>
			)
		},
		{
			key: 'p1',
			title: 'Lương P1',
			render: (_, row) => formatVND(row.preview?.p1)
		},
		{
			key: 'p2',
			title: 'Lương P2',
			render: (_, row) => formatVND(row.preview?.p2)
		},
		{
			key: 'gross',
			title: 'Tổng thu nhập',
			render: (_, row) => formatVND(row.preview?.subtotal_p1_p2)
		},
		{
			key: 'tax',
			title: 'Thuế TNCN',
			render: (_, row) => formatVND(row.preview?.tax || (row.preview?.subtotal_p1_p2 * 0.1))
		},
		{
			key: 'net',
			title: 'Thực nhận',
			render: (_, row) => (
				<span style={{ fontWeight: 700, color: 'var(--accent)' }}>
					{formatVND(row.preview?.net_salary || (row.preview?.subtotal_p1_p2 * 0.9))}
				</span>
			)
		},
		{
			key: 'status',
			title: 'Trạng thái',
			render: (status) => (
				<Badge variant={status === 'SUCCESS' || status === 'APPROVED' ? 'success' : 'warning'}>
					{status === 'SUCCESS' ? 'Đã lưu nháp' : status}
				</Badge>
			)
		},
		{
			key: 'actions',
			title: 'Thao tác',
			render: (_, row) => (
				<button
					className="btn btn-sm btn-secondary"
					onClick={() => {
						setSelectedRow(row);
						setDetailOpen(true);
					}}
				>
					Chi tiết
				</button>
			)
		}
	];

	const successCount = parseInt(jobStatus?.success_count || 0);
	const failedCount = parseInt(jobStatus?.failed_count || 0);
	const totalEmployees = parseInt(jobStatus?.total_employees || 0);
	const progressPercent = totalEmployees > 0 ? Math.round(((successCount + failedCount) / totalEmployees) * 100) : 0;

	return (
		<div className="page-container">
			<div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
				<div>
					<h1 className="page-title">Tính lương tháng</h1>
					<p className="page-subtitle" style={{ color: 'var(--text-muted)', fontSize: 14 }}>
						Khởi chạy job tính lương hàng loạt cho toàn công ty và lưu trữ kết quả
					</p>
				</div>
				<div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
					<input
						type="month"
						className="form-input"
						value={selectedPeriod}
						onChange={(e) => setSelectedPeriod(e.target.value)}
						disabled={polling}
						style={{ padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
					/>
					<button
						className="btn btn-primary"
						onClick={handleCalculate}
						disabled={loading || polling}
						style={{ display: 'flex', alignItems: 'center', gap: 8 }}
					>
						<span>⚡ Tính lương cả công ty</span>
					</button>
				</div>
			</div>

			{polling && (
				<div className="card" style={{ marginBottom: 24, padding: 24, background: 'rgba(79, 142, 247, 0.05)', borderColor: 'rgba(79, 142, 247, 0.2)' }}>
					<div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
						<div style={{ fontWeight: 600, fontSize: 16 }}>Đang xử lý tính toán bảng lương...</div>
						<div style={{ fontWeight: 700, color: 'var(--accent)', fontSize: 18 }}>{progressPercent}%</div>
					</div>
					<div className="progress-bar-bg" style={{ width: '100%', height: 8, background: 'rgba(255,255,255,0.1)', borderRadius: 4, overflow: 'hidden', marginBottom: 12 }}>
						<div className="progress-bar-fill" style={{ width: `${progressPercent}%`, height: '100%', background: 'var(--accent)', transition: 'width 0.3s ease' }} />
					</div>
					<div style={{ color: 'var(--text-muted)', fontSize: 13 }}>
						Tiến độ: {successCount + failedCount} / {totalEmployees} nhân sự ({successCount} thành công, {failedCount} thất bại)
					</div>
				</div>
			)}

			<div className="card">
				<Table
					columns={columns}
					data={records}
					loading={loading && !polling}
					emptyMessage="Chưa có bảng lương chính thức cho kỳ này"
				/>
			</div>

			<Modal
				isOpen={detailOpen}
				onClose={() => setDetailOpen(false)}
				title="Chi tiết bảng lương"
				size="md"
			>
				{selectedRow && (
					<div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
						<div style={{ borderBottom: '1px solid var(--border-color)', paddingBottom: 12 }}>
							<h3 style={{ fontSize: 16, fontWeight: 700 }}>{selectedRow.fullName}</h3>
							<p style={{ color: 'var(--text-muted)', fontSize: 13 }}>Mã nhân sự: {selectedRow.employee_id}</p>
						</div>

						<div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
							<div style={{ display: 'flex', justifyContent: 'space-between' }}>
								<span style={{ color: 'var(--text-secondary)' }}>Lương vị trí P1:</span>
								<span style={{ fontWeight: 600 }}>{formatVND(selectedRow.preview?.p1)}</span>
							</div>
							<div style={{ display: 'flex', justifyContent: 'space-between' }}>
								<span style={{ color: 'var(--text-secondary)' }}>Lương năng lực P2:</span>
								<span style={{ fontWeight: 600 }}>{formatVND(selectedRow.preview?.p2)}</span>
							</div>
							<div style={{ display: 'flex', justifyContent: 'space-between', borderTop: '1px solid var(--border-color)', paddingTop: 8 }}>
								<span style={{ color: 'var(--text-primary)', fontWeight: 600 }}>Tổng thu nhập chịu thuế:</span>
								<span style={{ fontWeight: 700, color: 'var(--text-primary)' }}>{formatVND(selectedRow.preview?.subtotal_p1_p2)}</span>
							</div>
							<div style={{ display: 'flex', justifyContent: 'space-between' }}>
								<span style={{ color: 'var(--text-secondary)' }}>Thuế thu nhập cá nhân:</span>
								<span style={{ fontWeight: 600, color: 'var(--error)' }}>-{formatVND(selectedRow.preview?.tax || (selectedRow.preview?.subtotal_p1_p2 * 0.1))}</span>
							</div>
							<div style={{ display: 'flex', justifyContent: 'space-between', borderTop: '2px solid var(--accent)', paddingTop: 12 }}>
								<span style={{ color: 'var(--accent)', fontWeight: 700, fontSize: 16 }}>Thực nhận chuyển khoản:</span>
								<span style={{ fontWeight: 800, color: 'var(--accent)', fontSize: 18 }}>
									{formatVND(selectedRow.preview?.net_salary || (selectedRow.preview?.subtotal_p1_p2 * 0.9))}
								</span>
							</div>
						</div>

						{selectedRow.preview?.breakdown && Object.keys(selectedRow.preview.breakdown).length > 0 && (
							<div style={{ marginTop: 12, borderTop: '1px solid var(--border-color)', paddingTop: 12 }}>
								<h4 style={{ fontSize: 14, fontWeight: 700, marginBottom: 8 }}>Thành phần chi tiết:</h4>
								<div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 }}>
									{Object.entries(selectedRow.preview.breakdown).map(([k, v]) => (
										<div key={k} style={{ padding: '8px 12px', background: 'var(--bg-primary)', borderRadius: 'var(--radius-sm)', fontSize: 13 }}>
											<div style={{ color: 'var(--text-muted)', fontSize: 11 }}>{k}</div>
											<div style={{ fontWeight: 600 }}>{formatVND(v)}</div>
										</div>
									))}
								</div>
							</div>
						)}
					</div>
				)}
			</Modal>
		</div>
	);
}

export default PayrollRunPage;
