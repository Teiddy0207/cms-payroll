// Mock data for the HRM / 3Ps payroll dashboard.
// Replace with real API data once the reporting endpoints exist.

export const DEPARTMENTS = [
  { name: 'Kinh doanh', weight: 1.35 },
  { name: 'Kỹ thuật', weight: 1.15 },
  { name: 'Sản xuất', weight: 1.45 },
  { name: 'Vận hành', weight: 1.05 },
  { name: 'CSKH', weight: 0.95 },
  { name: 'Marketing', weight: 0.8 },
  { name: 'Nhân sự', weight: 0.55 },
  { name: 'Kế toán', weight: 0.5 },
  { name: 'Pháp chế', weight: 0.3 },
];

const TOTAL_WEIGHT = DEPARTMENTS.reduce((sum, d) => sum + d.weight, 0);

export const KPI_SUMMARY = {
  totalPayroll: '4.85 tỷ VNĐ',
  totalEmployees: 1245,
  onTimeRate: '96%',
  topDepartment: 'Sản xuất',
  topPositionP2: 'Trưởng phòng',
  peakOtSlot: '18h - 20h',
};

export const CHART_COLORS = ['#10b981', '#3b82f6', '#f59e0b', '#6366f1', '#f43f5e', '#64748b', '#06b6d4'];

export const PAYROLL_BY_MONTH = [
  { month: 'T1', total: 4180, target: 4300 },
  { month: 'T2', total: 4420, target: 4400 },
  { month: 'T3', total: 4560, target: 4500 },
  { month: 'T4', total: 4610, target: 4600 },
  { month: 'T5', total: 4780, target: 4700 },
  { month: 'T6', total: 4850, target: 4800 },
];

export const HEADCOUNT_ATTENDANCE_BY_MONTH = [
  { month: 'T1', headcount: 1180, rate: 91 },
  { month: 'T2', headcount: 1195, rate: 93 },
  { month: 'T3', headcount: 1205, rate: 92 },
  { month: 'T4', headcount: 1220, rate: 94 },
  { month: 'T5', headcount: 1235, rate: 95 },
  { month: 'T6', headcount: 1245, rate: 96 },
];

export const CONTRACT_TYPE_DONUT = [
  { name: 'Chính thức', value: 812 },
  { name: 'Thử việc', value: 268 },
  { name: 'Thời vụ', value: 165 },
];

// Reuses real seeded employee codes/names so the dashboard reads as project-native.
const OT_EMPLOYEE_POOL = [
  { code: '2510', name: 'Cao Xuân Hoan' },
  { code: '2529', name: 'Đặng Xuân Bình' },
  { code: '2506', name: 'Đào Hồng Quang' },
  { code: '2897', name: 'Trần Đại Hồng Dũng' },
  { code: '2825', name: 'Nguyễn Thụy Lưu' },
  { code: '3442', name: 'Hoàng Thị Đông' },
  { code: '3453', name: 'Phan Tá Nghĩa' },
  { code: '2519', name: 'Nguyễn Ngọc Linh' },
  { code: '2592', name: 'Trần Ngọc Huy' },
  { code: '2892', name: 'Nguyễn Minh Hà' },
];

const OT_HOURS_BASE = [42, 39, 37, 35, 33, 31, 29, 27, 25, 23];

export const TOP_OT_EMPLOYEES = OT_EMPLOYEE_POOL.map((emp, i) => ({
  name: emp.name,
  code: emp.code,
  hours: OT_HOURS_BASE[i],
}));

export const PAYROLL_BY_WEEKDAY = [
  { day: 'T2', value: 705 },
  { day: 'T3', value: 690 },
  { day: 'T4', value: 715 },
  { day: 'T5', value: 700 },
  { day: 'T6', value: 730 },
  { day: 'T7', value: 260 },
  { day: 'CN', value: 90 },
];

export const OT_BY_TIMESLOT_STACKED = [
  { slot: '17h-18h', 'Kinh doanh': 120, 'Kỹ thuật': 95, 'Sản xuất': 150 },
  { slot: '18h-19h', 'Kinh doanh': 160, 'Kỹ thuật': 130, 'Sản xuất': 210 },
  { slot: '19h-20h', 'Kinh doanh': 90, 'Kỹ thuật': 140, 'Sản xuất': 175 },
  { slot: '20h-21h', 'Kinh doanh': 40, 'Kỹ thuật': 60, 'Sản xuất': 80 },
];

export const SALARY_DISTRIBUTION_BY_LEVEL = [
  { level: 'Nhân viên', chinhThuc: 9.5, thuViec: 6.8 },
  { level: 'Chuyên viên', chinhThuc: 14.2, thuViec: 9.5 },
  { level: 'Trưởng nhóm', chinhThuc: 19.8, thuViec: null },
  { level: 'Trưởng phòng', chinhThuc: 28.5, thuViec: null },
  { level: 'Giám đốc', chinhThuc: 45.0, thuViec: null },
];

export const P1_P2_BY_WEEK = [
  { week: 'W1', p1: 12.1, p2: 6.2 },
  { week: 'W2', p1: 12.3, p2: 6.4 },
  { week: 'W3', p1: 12.4, p2: 6.5 },
  { week: 'W4', p1: 12.6, p2: 6.7 },
  { week: 'W5', p1: 12.7, p2: 6.9 },
  { week: 'W6', p1: 12.8, p2: 7.0 },
  { week: 'W7', p1: 12.9, p2: 7.2 },
  { week: 'W8', p1: 13.1, p2: 7.4 },
];

export const ALLOWANCE_BY_EMPLOYEE_TYPE = [
  { type: 'Ăn trưa', chinhThuc: 96, thuViec: 88 },
  { type: 'Xăng xe', chinhThuc: 82, thuViec: 60 },
  { type: 'Điện thoại', chinhThuc: 45, thuViec: 20 },
  { type: 'Nhà ở', chinhThuc: 18, thuViec: 5 },
  { type: 'Chuyên cần', chinhThuc: 90, thuViec: 75 },
];

export function getDepartmentFraction(deptIndex) {
  return DEPARTMENTS[deptIndex].weight / TOTAL_WEIGHT;
}

export function scaleByDepartment(deptIndex, base, keys) {
  const fraction = getDepartmentFraction(deptIndex);
  return base.map((row) => {
    const next = { ...row };
    keys.forEach((k) => {
      if (typeof row[k] === 'number') next[k] = Math.round(row[k] * fraction);
    });
    return next;
  });
}
