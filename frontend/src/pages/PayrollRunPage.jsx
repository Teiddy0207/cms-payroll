import { useState, useEffect, useCallback, useRef } from 'react';
import { calculatorAPI, departmentsAPI } from '../api/client.js';
import { Table, DatePicker, Select, Button, Progress, Modal, Tag, Card } from 'antd';
import { EyeOutlined } from '@ant-design/icons';
import { useToast } from '../hooks/useToast.js';
import dayjs from 'dayjs';

const { Option } = Select;

export function PayrollRunPage() {
	const toast = useToast();
	const [loading, setLoading] = useState(false);
	const [records, setRecords] = useState([]);
	const [departments, setDepartments] = useState([]);
	const [selectedPeriod, setSelectedPeriod] = useState(() => {
		const d = new Date();
		return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
	});
	const [selectedDepartment, setSelectedDepartment] = useState(undefined);

	const [jobStatus, setJobStatus] = useState(null);
	const [polling, setPolling] = useState(false);
	const pollInterval = useRef(null);

	const [selectedRow, setSelectedRow] = useState(null);
	const [detailOpen, setDetailOpen] = useState(false);

	const formatVND = (value) => {
		return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(value || 0);
	};

	const fetchDepartments = async () => {
		try {
			const res = await departmentsAPI.list({ page_size: 100 });
			setDepartments(res.data?.data?.items || []);
		} catch (err) {
			console.error("fetchDepartments error:", err);
		}
	};

	const fetchRecords = useCallback(async (period, departmentId) => {
		setLoading(true);
		try {
			const extraParams = {};
			if (departmentId) {
				extraParams.department_id = departmentId;
			}
			const res = await calculatorAPI.getSavedRecords(period, extraParams);
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
				fetchRecords(selectedPeriod, selectedDepartment);
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
	}, [fetchRecords, selectedPeriod, selectedDepartment, toast]);

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
		fetchDepartments();
	}, []);

	useEffect(() => {
		fetchRecords(selectedPeriod, selectedDepartment);
	}, [selectedPeriod, selectedDepartment, fetchRecords]);

	useEffect(() => {
		return () => {
			if (pollInterval.current) {
				clearInterval(pollInterval.current);
			}
		};
	}, []);

	const columns = [
		{
			dataIndex: 'fullName',
			title: 'Nhân sự',
			key: 'fullName',
			render: (_, row) => (
				<div>
					<div style={{ fontWeight: 600, color: 'var(--text-primary)' }}>{row.fullName}</div>
					<div style={{ fontSize: '12px', color: 'var(--text-muted)' }}>Mã NV: {row.employee_id.substring(0, 8)}</div>
				</div>
			)
		},
		{
			title: 'Lương P1',
			key: 'p1',
			render: (_, row) => formatVND(row.preview?.p1)
		},
		{
			title: 'Lương P2',
			key: 'p2',
			render: (_, row) => formatVND(row.preview?.p2)
		},
		{
			title: 'Tổng thu nhập',
			key: 'gross',
			render: (_, row) => formatVND(row.preview?.subtotal_p1_p2)
		},
		{
			title: 'Thuế TNCN',
			key: 'tax',
			render: (_, row) => formatVND(row.preview?.tax || (row.preview?.subtotal_p1_p2 * 0.1))
		},
		{
			title: 'Thực nhận',
			key: 'net',
			render: (_, row) => (
				<span style={{ fontWeight: 700, color: 'var(--accent)' }}>
					{formatVND(row.preview?.net_salary || (row.preview?.subtotal_p1_p2 * 0.9))}
				</span>
			)
		},
		{
			title: 'Trạng thái',
			key: 'status',
			dataIndex: 'status',
			render: (status) => {
				const isSuccess = status === 'SUCCESS' || status === 'APPROVED';
				return (
					<Tag color={isSuccess ? 'success' : 'warning'}>
						{status === 'SUCCESS' ? 'Đã lưu nháp' : status}
					</Tag>
				);
			}
		},
		{
			title: 'Thao tác',
			key: 'actions',
			render: (_, row) => (
				<Button
					type="primary"
					ghost
					size="small"
					icon={<EyeOutlined />}
					onClick={() => {
						setSelectedRow(row);
						setDetailOpen(true);
					}}
				>
					Chi tiết
				</Button>
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
					{/* Chọn Phòng ban */}
					<Select
						placeholder="Lọc theo phòng ban"
						allowClear
						style={{ width: 200 }}
						value={selectedDepartment}
						onChange={(val) => setSelectedDepartment(val)}
						disabled={polling}
					>
						{departments.map((dept) => (
							<Option key={dept.id} value={dept.id}>
								{dept.name}
							</Option>
						))}
					</Select>

					{/* Chọn Kỳ lương */}
					<DatePicker
						picker="month"
						allowClear={false}
						value={selectedPeriod ? dayjs(selectedPeriod, 'YYYY-MM') : null}
						onChange={(date) => {
							if (date) {
								setSelectedPeriod(date.format('YYYY-MM'));
							}
						}}
						disabled={polling}
						style={{ width: 140 }}
					/>

					{/* Tính lương */}
					<Button
						type="primary"
						onClick={handleCalculate}
						loading={loading || polling}
					>
						⚡ Tính lương cả công ty
					</Button>
				</div>
			</div>

			{polling && (
				<Card
					style={{ marginBottom: 24, background: 'rgba(16, 185, 129, 0.05)', borderColor: 'rgba(16, 185, 129, 0.2)' }}
					bodyStyle={{ padding: 24 }}
				>
					<div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
						<div style={{ fontWeight: 600, fontSize: 16 }}>Đang xử lý tính toán bảng lương...</div>
						<div style={{ fontWeight: 700, color: 'var(--accent)', fontSize: 18 }}>{progressPercent}%</div>
					</div>
					<Progress
						percent={progressPercent}
						status="active"
						strokeColor="#10b981"
						showInfo={false}
						style={{ marginBottom: 12 }}
					/>
					<div style={{ color: 'var(--text-muted)', fontSize: 13 }}>
						Tiến độ: {successCount + failedCount} / {totalEmployees} nhân sự ({successCount} thành công, {failedCount} thất bại)
					</div>
				</Card>
			)}

			<Card bodyStyle={{ padding: 0 }} style={{ overflow: 'hidden' }}>
				<Table
					columns={columns}
					dataSource={records}
					loading={loading && !polling}
					rowKey={(record) => record.id}
					pagination={{ pageSize: 10 }}
					locale={{ emptyText: 'Chưa có bảng lương chính thức cho kỳ này' }}
				/>
			</Card>

			<Modal
				open={detailOpen}
				onCancel={() => setDetailOpen(false)}
				title="Chi tiết bảng lương"
				footer={[
					<Button key="close" onClick={() => setDetailOpen(false)}>
						Đóng
					</Button>
				]}
			>
				{selectedRow && (
					<div style={{ display: 'flex', flexDirection: 'column', gap: 16, marginTop: 16 }}>
						<div style={{ borderBottom: '1px solid var(--border-color)', paddingBottom: 12 }}>
							<h3 style={{ fontSize: 16, fontWeight: 700, margin: 0 }}>{selectedRow.fullName}</h3>
							<p style={{ color: 'var(--text-muted)', fontSize: 13, margin: '4px 0 0' }}>
								Mã nhân sự: {selectedRow.employee_id}
							</p>
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
