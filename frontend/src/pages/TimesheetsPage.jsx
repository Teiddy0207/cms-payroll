import { useState, useEffect, useCallback } from 'react';
import { timekeepingAPI, departmentsAPI, employeesAPI } from '../api/client.js';
import { Calendar, Select, Modal, Badge, Card, Button, Spin, Tag, Descriptions, DatePicker } from 'antd';
import { CalendarOutlined, InfoCircleOutlined } from '@ant-design/icons';
import { useToast } from '../hooks/useToast.js';
import dayjs from 'dayjs';

const { Option } = Select;
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

export function TimesheetsPage() {
	const toast = useToast();
	const [sheets, setSheets] = useState([]);
	const [employees, setEmployees] = useState([]);
	const [selectedEmployeeId, setSelectedEmployeeId] = useState(null);
	const [currentPeriod, setCurrentPeriod] = useState(dayjs());
	const [loading, setLoading] = useState(false);
	const [isAdmin, setIsAdmin] = useState(false);
	const [currentUserProfileId, setCurrentUserProfileId] = useState(null);
	const [selectedDateDetails, setSelectedDateDetails] = useState(null);
	const [detailsModalOpen, setDetailsModalOpen] = useState(false);

	const [calcModalOpen, setCalcModalOpen] = useState(false);
	const [calcDates, setCalcDates] = useState({
		start_date: '',
		end_date: ''
	});

	// Phân tích JWT để lấy thông tin tài khoản đăng nhập
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
				setCurrentUserProfileId(payload.user_profile_id || payload.user_id);
			}
		}
	}, []);

	// Load danh sách nhân viên nếu là Admin để có thể chọn xem lịch của người khác
	useEffect(() => {
		const loadEmployees = async () => {
			try {
				const res = await employeesAPI.list({ page_size: 100 });
				const list = res.data?.data?.items || [];
				setEmployees(list);
				if (list.length > 0 && !selectedEmployeeId) {
					// Chọn nhân viên đầu tiên làm mặc định nếu là Admin
					setSelectedEmployeeId(list[0].id);
				}
			} catch (err) {
				console.error("loadEmployees error:", err);
			}
		};

		if (isAdmin) {
			loadEmployees();
		}
	}, [isAdmin, selectedEmployeeId]);

	// Tải dữ liệu bảng công cho lịch theo Tháng và Nhân viên được chọn
	const loadSheets = useCallback(async () => {
		setLoading(true);
		try {
			const periodStr = currentPeriod.format('YYYY-MM');
			const params = {
				page_number: 1,
				page_size: 100, // Load tất cả các ngày trong tháng
				period: periodStr
			};
			
			// Nếu là admin và có chọn nhân viên cụ thể
			if (isAdmin && selectedEmployeeId) {
				params.search = employees.find(e => e.id === selectedEmployeeId)?.code || '';
			}

			const res = await timekeepingAPI.getSheets(params);
			setSheets(res.data?.data?.items || []);
		} catch (err) {
			console.error("loadSheets error:", err);
			toast.error('Lỗi', 'Không thể tải dữ liệu bảng công');
		} finally {
			setLoading(false);
		}
	}, [currentPeriod, selectedEmployeeId, isAdmin, employees, toast]);

	useEffect(() => {
		loadSheets();
	}, [currentPeriod, selectedEmployeeId, loadSheets]);

	// Tự động gán cấu hình ngày khi mở hộp thoại Tính công ngày
	const handleOpenCalc = () => {
		const startOfMonth = currentPeriod.startOf('month').format('YYYY-MM-DD');
		const endOfMonth = currentPeriod.endOf('month').format('YYYY-MM-DD');
		setCalcDates({
			start_date: startOfMonth,
			end_date: endOfMonth
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
			toast.success('Thành công', 'Đã tính toán xong bảng công ngày');
			setCalcModalOpen(false);
			loadSheets();
		} catch (err) {
			console.error("calculate error:", err);
			toast.error('Lỗi', 'Lỗi tính toán bảng công');
		} finally {
			setLoading(false);
		}
	};

	// Khi người dùng bấm trực tiếp vào ngày
	const handleCellClick = (value, daySheet) => {
		const dateStr = value.format('YYYY-MM-DD');
		if (daySheet) {
			setSelectedDateDetails(daySheet);
		} else {
			setSelectedDateDetails({
				date: dateStr,
				status: 'NO_DATA',
				actual_work_day: 0,
				ot_hours: 0
			});
		}
		setDetailsModalOpen(true);
	};

	// Xử lý render cho từng ô ngày trên Lịch
	const dateCellRender = (value) => {
		const dateStr = value.format('YYYY-MM-DD');
		const daySheet = sheets.find(s => s.date && s.date.substring(0, 10) === dateStr);

		let badgeStatus = 'default';
		let badgeText = '';

		if (daySheet) {
			switch (daySheet.status) {
				case 'NORMAL':
					badgeStatus = 'success';
					badgeText = 'Đủ công';
					break;
				case 'LATE':
					badgeStatus = 'warning';
					badgeText = 'Đi muộn';
					break;
				case 'EARLY':
					badgeStatus = 'warning';
					badgeText = 'Về sớm';
					break;
				case 'LATE_EARLY':
					badgeStatus = 'warning';
					badgeText = 'Muộn/Sớm';
					break;
				case 'ABSENT':
					badgeStatus = 'error';
					badgeText = 'Vắng mặt';
					break;
			}
		}

		return (
			<div 
				style={{ minHeight: '50px', width: '100%', padding: '4px' }} 
				onClick={() => handleCellClick(value, daySheet)}
			>
				{daySheet && (
					<>
						<Badge status={badgeStatus} text={badgeText} style={{ fontSize: '11px', display: 'block' }} />
						{daySheet.ot_hours > 0 && (
							<div style={{ fontSize: '10px', color: 'var(--accent)', fontWeight: 700, marginLeft: 8 }}>
								OT: +{daySheet.ot_hours.toFixed(1)}h
							</div>
						)}
					</>
				)}
			</div>
		);
	};

	const formatDateTime = (dateTimeStr) => {
		if (!dateTimeStr) return '—';
		// Loại bỏ hậu tố Z/UTC để trình duyệt hiểu là giờ Local đã lưu ở DB, tránh bị cộng lệch múi giờ (+7)
		const localStr = dateTimeStr.endsWith('Z') ? dateTimeStr.slice(0, -1) : dateTimeStr;
		return dayjs(localStr).format('HH:mm');
	};

	const getStatusTag = (status) => {
		switch (status) {
			case 'NORMAL':
				return <Tag color="success">Đúng giờ (Đủ công)</Tag>;
			case 'LATE':
				return <Tag color="warning">Đi muộn</Tag>;
			case 'EARLY':
				return <Tag color="warning">Về sớm</Tag>;
			case 'LATE_EARLY':
				return <Tag color="orange">Đi muộn & Về sớm</Tag>;
			case 'ABSENT':
				return <Tag color="error">Vắng mặt (Không công)</Tag>;
			default:
				return <Tag color="default">Chưa có dữ liệu</Tag>;
		}
	};

	return (
		<div className="page-container">
			<div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
				<div>
					<h1 className="page-title">Bảng chấm công tổng hợp (Dạng lịch)</h1>
					<p className="page-subtitle" style={{ color: 'var(--text-muted)', fontSize: 14 }}>
						Xem lịch biểu chấm công chi tiết. Click vào ngày bất kỳ để xem thời gian check-in/out thực tế.
					</p>
				</div>
				<div style={{ display: 'flex', gap: 12 }}>
					{isAdmin && (
						<Select
							showSearch
							placeholder="Chọn nhân viên để xem lịch"
							optionFilterProp="children"
							style={{ width: 260 }}
							value={selectedEmployeeId}
							onChange={(val) => setSelectedEmployeeId(val)}
						>
							{employees.map(emp => (
								<Option key={emp.id} value={emp.id}>
									{emp.full_name || emp.fullName} ({emp.code})
								</Option>
							))}
						</Select>
					)}
					<Button 
						type="primary" 
						onClick={handleOpenCalc}
						disabled={loading}
					>
						Tính công ngày
					</Button>
				</div>
			</div>

			<Card bodyStyle={{ padding: 24 }}>
				<Spin spinning={loading} tip="Đang tải dữ liệu lịch chấm công...">
					<Calendar
						value={currentPeriod}
						onPanelChange={(date) => setCurrentPeriod(date)}
						dateCellRender={dateCellRender}
						headerRender={({ value, onChange }) => {
							const current = value.clone();
							return (
								<div style={{ padding: 12, display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: 'var(--bg-primary)', borderRadius: '8px', marginBottom: 16 }}>
									<div style={{ fontWeight: 700, fontSize: 16, display: 'flex', alignItems: 'center', gap: 8 }}>
										<CalendarOutlined style={{ color: 'var(--accent)' }} />
										Tháng {current.format('MM / YYYY')}
									</div>
									<div style={{ display: 'flex', gap: 8 }}>
										<Button 
											size="small"
											onClick={() => {
												const prev = current.subtract(1, 'month');
												onChange(prev);
												setCurrentPeriod(prev);
											}}
										>
											Tháng trước
										</Button>
										<Button 
											size="small"
											onClick={() => {
												const now = dayjs();
												onChange(now);
												setCurrentPeriod(now);
											}}
										>
											Tháng này
										</Button>
										<Button 
											size="small"
											onClick={() => {
												const next = current.add(1, 'month');
												onChange(next);
												setCurrentPeriod(next);
											}}
										>
											Tháng sau
										</Button>
									</div>
								</div>
							);
						}}
					/>
				</Spin>
			</Card>

			{/* Modal chi tiết chấm công của ngày được click */}
			<Modal
				open={detailsModalOpen}
				onCancel={() => setDetailsModalOpen(false)}
				title={
					<div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
						<InfoCircleOutlined style={{ color: 'var(--accent)' }} />
						<span>Chi tiết ngày công: {selectedDateDetails ? dayjs(selectedDateDetails.date).format('DD/MM/YYYY') : ''}</span>
					</div>
				}
				footer={[
					<Button key="close" type="primary" onClick={() => setDetailsModalOpen(false)}>
						Đóng
					</Button>
				]}
			>
				{selectedDateDetails && (
					<div style={{ marginTop: 16 }}>
						{selectedDateDetails.status === 'NO_DATA' ? (
							<div style={{ textAlign: 'center', padding: '24px 0', color: 'var(--text-muted)' }}>
								Không tìm thấy dữ liệu chấm công cho ngày này. Có thể là ngày cuối tuần hoặc kỳ tính công chưa được chạy.
							</div>
						) : (
							<Descriptions bordered column={1} size="small">
								<Descriptions.Item label="Họ và tên">
									<strong style={{ color: 'var(--text-primary)' }}>{selectedDateDetails.full_name || selectedDateDetails.fullName}</strong>
								</Descriptions.Item>
								<Descriptions.Item label="Mã nhân sự">
									{selectedDateDetails.employee_code || selectedDateDetails.employeeCode}
								</Descriptions.Item>
								<Descriptions.Item label="Trạng thái">
									{getStatusTag(selectedDateDetails.status)}
								</Descriptions.Item>
								<Descriptions.Item label="Giờ vào (Check-in)">
									<span style={{ fontWeight: 600, color: '#16a34a' }}>
										{formatDateTime(selectedDateDetails.check_in || selectedDateDetails.checkIn)}
									</span>
								</Descriptions.Item>
								<Descriptions.Item label="Giờ ra (Check-out)">
									<span style={{ fontWeight: 600, color: '#dc2626' }}>
										{formatDateTime(selectedDateDetails.check_out || selectedDateDetails.checkOut)}
									</span>
								</Descriptions.Item>
								<Descriptions.Item label="Số ngày công ghi nhận">
									<span style={{ fontWeight: 700, color: (selectedDateDetails.actual_work_day || selectedDateDetails.actualWorkDay) > 0 ? '#16a34a' : '#dc2626' }}>
										{(selectedDateDetails.actual_work_day || selectedDateDetails.actualWorkDay || 0).toFixed(1)} công
									</span>
								</Descriptions.Item>
								<Descriptions.Item label="Giờ làm thêm (OT)">
									{(selectedDateDetails.ot_hours || selectedDateDetails.otHours) > 0 ? (
										<span style={{ fontWeight: 700, color: 'var(--accent)' }}>
											+{(selectedDateDetails.ot_hours || selectedDateDetails.otHours).toFixed(1)} giờ
										</span>
									) : '0'}
								</Descriptions.Item>
							</Descriptions>
						)}
					</div>
				)}
			</Modal>

			{/* Modal Tính công ngày */}
			<Modal
				open={calcModalOpen}
				onCancel={() => setCalcModalOpen(false)}
				title="Tính công ngày tự động"
				footer={null}
			>
				<form onSubmit={handleRunCalculation} style={{ display: 'flex', flexDirection: 'column', gap: 16, marginTop: 16 }}>
					<p style={{ color: 'var(--text-muted)', fontSize: 13, margin: 0 }}>
						Hệ thống sẽ quét toàn bộ logs chấm công thô, đơn xin bù công và đơn đăng ký OT đã được duyệt để tổng hợp số ngày công của nhân sự trong khoảng thời gian được chọn.
					</p>

					<div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
						<label style={{ fontWeight: 600 }}>Từ ngày <span style={{ color: 'red' }}>*</span></label>
						<DatePicker
							style={{ width: '100%' }}
							value={calcDates.start_date ? dayjs(calcDates.start_date, 'YYYY-MM-DD') : null}
							onChange={(date) => setCalcDates(prev => ({ ...prev, start_date: date ? date.format('YYYY-MM-DD') : '' }))}
						/>
					</div>

					<div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
						<label style={{ fontWeight: 600 }}>Đến ngày <span style={{ color: 'red' }}>*</span></label>
						<DatePicker
							style={{ width: '100%' }}
							value={calcDates.end_date ? dayjs(calcDates.end_date, 'YYYY-MM-DD') : null}
							onChange={(date) => setCalcDates(prev => ({ ...prev, end_date: date ? date.format('YYYY-MM-DD') : '' }))}
						/>
					</div>

					<div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 12 }}>
						<Button onClick={() => setCalcModalOpen(false)}>Hủy</Button>
						<Button type="primary" htmlType="submit" loading={loading}>
							Kích hoạt tính công
						</Button>
					</div>
				</form>
			</Modal>
		</div>
	);
}

export default TimesheetsPage;
