import { useMemo, useState } from 'react';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  LineChart,
  Line,
  ComposedChart,
  AreaChart,
  Area,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from 'recharts';
import {
  DEPARTMENTS,
  KPI_SUMMARY,
  CHART_COLORS,
  PAYROLL_BY_MONTH,
  HEADCOUNT_ATTENDANCE_BY_MONTH,
  CONTRACT_TYPE_DONUT,
  TOP_OT_EMPLOYEES,
  PAYROLL_BY_WEEKDAY,
  OT_BY_TIMESLOT_STACKED,
  SALARY_DISTRIBUTION_BY_LEVEL,
  P1_P2_BY_WEEK,
  ALLOWANCE_BY_EMPLOYEE_TYPE,
  scaleByDepartment,
} from '../mock/dashboardMock.js';

function ChartCard({ title, subtitle, span, children }) {
  return (
    <div className="chart-card" style={span ? { gridColumn: '1 / -1' } : undefined}>
      <div className="chart-card-header">
        <div className="chart-card-title">{title}</div>
        {subtitle && <div className="chart-card-subtitle">{subtitle}</div>}
      </div>
      <div className="chart-card-body">{children}</div>
    </div>
  );
}

export function DashboardPage() {
  const [selectedDept, setSelectedDept] = useState(0);

  const payrollByMonth = useMemo(
    () => scaleByDepartment(selectedDept, PAYROLL_BY_MONTH, ['total', 'target']),
    [selectedDept]
  );
  const headcountAttendance = useMemo(
    () => scaleByDepartment(selectedDept, HEADCOUNT_ATTENDANCE_BY_MONTH, ['headcount']),
    [selectedDept]
  );
  const contractType = useMemo(
    () => scaleByDepartment(selectedDept, CONTRACT_TYPE_DONUT, ['value']),
    [selectedDept]
  );

  const kpiCards = [
    { label: 'Tổng quỹ lương tháng', value: KPI_SUMMARY.totalPayroll, color: '#10b981' },
    { label: 'Tổng nhân sự', value: KPI_SUMMARY.totalEmployees.toLocaleString('vi-VN'), color: '#3b82f6' },
    { label: 'Tỷ lệ chấm công đúng giờ', value: KPI_SUMMARY.onTimeRate, color: '#06b6d4' },
    { label: 'Phòng ban đông nhất', value: KPI_SUMMARY.topDepartment, color: '#f59e0b' },
    { label: 'Vị trí P2 cao nhất', value: KPI_SUMMARY.topPositionP2, color: '#6366f1' },
    { label: 'Khung giờ OT cao điểm', value: KPI_SUMMARY.peakOtSlot, color: '#f43f5e' },
  ];

  return (
    <div className="dashboard-page">
      <div className="page-header" style={{ marginBottom: 20 }}>
        <h1 className="page-title">Dashboard tổng quan HRM &amp; Lương 3Ps</h1>
        <p className="page-subtitle" style={{ color: 'var(--text-muted)', fontSize: 14 }}>
          Số liệu minh họa (mock) — sẽ được thay bằng dữ liệu thực khi có API báo cáo
        </p>
      </div>

      {/* KPI row */}
      <div className="kpi-grid" style={{ marginBottom: 24 }}>
        {kpiCards.map((card) => (
          <div key={card.label} className="stat-card">
            <div className="stat-info">
              <h3 style={{ color: card.color, fontSize: 20 }}>{card.value}</h3>
              <p>{card.label}</p>
            </div>
          </div>
        ))}
      </div>

      <div className="dashboard-body-layout">
        {/* Department filter sidebar */}
        <aside className="dept-sidebar">
          <div className="dept-sidebar-title">TOP 9<br />PHÒNG BAN</div>
          {DEPARTMENTS.map((dept, i) => (
            <button
              key={dept.name}
              className={`dept-btn${selectedDept === i ? ' active' : ''}`}
              onClick={() => setSelectedDept(i)}
            >
              {dept.name}
            </button>
          ))}
        </aside>

        <div className="dashboard-sections">
          {/* Section A: Department performance */}
          <div className="dashboard-section-title">
            HIỆU SUẤT PHÒNG BAN <span className="text-muted" style={{ fontWeight: 500, fontSize: 12.5 }}>— Đang xem: {DEPARTMENTS[selectedDept].name}</span>
          </div>
          <div className="dashboard-charts-grid">
            <ChartCard title="Top 10 nhân sự có giờ OT cao nhất tháng" subtitle="Đơn vị: giờ">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={TOP_OT_EMPLOYEES} layout="vertical" margin={{ left: 10, right: 16 }}>
                  <CartesianGrid strokeDasharray="3 3" horizontal={false} stroke="var(--border-color)" />
                  <XAxis type="number" tick={{ fontSize: 11 }} />
                  <YAxis type="category" dataKey="name" width={130} tick={{ fontSize: 11 }} />
                  <Tooltip formatter={(v) => [`${v} giờ`, 'OT']} />
                  <Bar dataKey="hours" fill={CHART_COLORS[0]} radius={[0, 4, 4, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </ChartCard>

            <ChartCard title="Cơ cấu loại hợp đồng" subtitle="Toàn công ty">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={contractType}
                    dataKey="value"
                    nameKey="name"
                    innerRadius={55}
                    outerRadius={85}
                    paddingAngle={2}
                  >
                    {contractType.map((entry, i) => (
                      <Cell key={entry.name} fill={CHART_COLORS[i % CHART_COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                  <Legend wrapperStyle={{ fontSize: 12 }} />
                </PieChart>
              </ResponsiveContainer>
            </ChartCard>

            <ChartCard title="Quỹ lương theo tháng so với chỉ tiêu" subtitle="Đơn vị: triệu VNĐ">
              <ResponsiveContainer width="100%" height="100%">
                <ComposedChart data={payrollByMonth} margin={{ left: -10, right: 16 }}>
                  <CartesianGrid strokeDasharray="3 3" stroke="var(--border-color)" />
                  <XAxis dataKey="month" tick={{ fontSize: 11 }} />
                  <YAxis tick={{ fontSize: 11 }} />
                  <Tooltip />
                  <Legend wrapperStyle={{ fontSize: 12 }} />
                  <Bar dataKey="total" name="Thực chi" fill={CHART_COLORS[0]} radius={[4, 4, 0, 0]} />
                  <Line dataKey="target" name="Chỉ tiêu" stroke={CHART_COLORS[4]} strokeWidth={2} dot={{ r: 3 }} />
                </ComposedChart>
              </ResponsiveContainer>
            </ChartCard>

            <ChartCard title="Số lượng nhân sự &amp; tỷ lệ chấm công" subtitle="Theo tháng">
              <ResponsiveContainer width="100%" height="100%">
                <ComposedChart data={headcountAttendance} margin={{ left: -10, right: 16 }}>
                  <CartesianGrid strokeDasharray="3 3" stroke="var(--border-color)" />
                  <XAxis dataKey="month" tick={{ fontSize: 11 }} />
                  <YAxis yAxisId="left" tick={{ fontSize: 11 }} />
                  <YAxis yAxisId="right" orientation="right" tick={{ fontSize: 11 }} unit="%" />
                  <Tooltip />
                  <Legend wrapperStyle={{ fontSize: 12 }} />
                  <Bar yAxisId="left" dataKey="headcount" name="Số nhân sự" fill={CHART_COLORS[1]} radius={[4, 4, 0, 0]} />
                  <Line yAxisId="right" dataKey="rate" name="Tỷ lệ chấm công (%)" stroke={CHART_COLORS[0]} strokeWidth={2} dot={{ r: 3 }} />
                </ComposedChart>
              </ResponsiveContainer>
            </ChartCard>
          </div>

          {/* Section B: 3Ps payroll performance */}
          <div className="dashboard-section-title" style={{ marginTop: 28 }}>
            HIỆU SUẤT LƯƠNG 3PS
          </div>
          <div className="dashboard-charts-grid">
            <ChartCard title="Quỹ lương theo ngày trong tuần" subtitle="Đơn vị: triệu VNĐ">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={PAYROLL_BY_WEEKDAY} margin={{ left: -10, right: 16 }}>
                  <CartesianGrid strokeDasharray="3 3" stroke="var(--border-color)" />
                  <XAxis dataKey="day" tick={{ fontSize: 11 }} />
                  <YAxis tick={{ fontSize: 11 }} />
                  <Tooltip />
                  <Bar dataKey="value" fill={CHART_COLORS[2]} radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </ChartCard>

            <ChartCard title="OT theo khung giờ — Top 3 phòng ban" subtitle="Đơn vị: giờ công OT">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={OT_BY_TIMESLOT_STACKED} margin={{ left: -10, right: 16 }}>
                  <CartesianGrid strokeDasharray="3 3" stroke="var(--border-color)" />
                  <XAxis dataKey="slot" tick={{ fontSize: 11 }} />
                  <YAxis tick={{ fontSize: 11 }} />
                  <Tooltip />
                  <Legend wrapperStyle={{ fontSize: 12 }} />
                  <Bar dataKey="Kinh doanh" stackId="ot" fill={CHART_COLORS[0]} />
                  <Bar dataKey="Kỹ thuật" stackId="ot" fill={CHART_COLORS[1]} />
                  <Bar dataKey="Sản xuất" stackId="ot" fill={CHART_COLORS[2]} radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </ChartCard>

            <ChartCard title="Phân bổ mức lương theo cấp bậc" subtitle="Đơn vị: triệu VNĐ/tháng">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={SALARY_DISTRIBUTION_BY_LEVEL} margin={{ left: -10, right: 16 }}>
                  <CartesianGrid strokeDasharray="3 3" stroke="var(--border-color)" />
                  <XAxis dataKey="level" tick={{ fontSize: 10.5 }} />
                  <YAxis tick={{ fontSize: 11 }} />
                  <Tooltip />
                  <Legend wrapperStyle={{ fontSize: 12 }} />
                  <Bar dataKey="chinhThuc" name="Chính thức" fill={CHART_COLORS[0]} radius={[4, 4, 0, 0]} />
                  <Bar dataKey="thuViec" name="Thử việc" fill={CHART_COLORS[5]} radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </ChartCard>

            <ChartCard title="Lương P1 &amp; P2 trung bình theo tuần" subtitle="Đơn vị: triệu VNĐ">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={P1_P2_BY_WEEK} margin={{ left: -10, right: 16 }}>
                  <CartesianGrid strokeDasharray="3 3" stroke="var(--border-color)" />
                  <XAxis dataKey="week" tick={{ fontSize: 11 }} />
                  <YAxis tick={{ fontSize: 11 }} />
                  <Tooltip />
                  <Legend wrapperStyle={{ fontSize: 12 }} />
                  <Line dataKey="p1" name="Lương P1 (vị trí)" stroke={CHART_COLORS[0]} strokeWidth={2} dot={{ r: 3 }} />
                  <Line dataKey="p2" name="Lương P2 (năng lực)" stroke={CHART_COLORS[1]} strokeWidth={2} dot={{ r: 3 }} />
                </LineChart>
              </ResponsiveContainer>
            </ChartCard>

            <ChartCard
              title="Tỷ lệ nhận phụ cấp theo loại nhân viên"
              subtitle="Đơn vị: % nhân sự nhận phụ cấp"
              span
            >
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={ALLOWANCE_BY_EMPLOYEE_TYPE} margin={{ left: -10, right: 16 }}>
                  <CartesianGrid strokeDasharray="3 3" stroke="var(--border-color)" />
                  <XAxis dataKey="type" tick={{ fontSize: 11 }} />
                  <YAxis tick={{ fontSize: 11 }} unit="%" />
                  <Tooltip formatter={(v) => `${v}%`} />
                  <Legend wrapperStyle={{ fontSize: 12 }} />
                  <Bar dataKey="chinhThuc" name="Chính thức" fill={CHART_COLORS[0]} radius={[4, 4, 0, 0]} />
                  <Bar dataKey="thuViec" name="Thử việc" fill={CHART_COLORS[5]} radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </ChartCard>
          </div>
        </div>
      </div>

      <style>{`
        .kpi-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
          gap: 16px;
        }

        .dashboard-body-layout {
          display: flex;
          gap: 20px;
          align-items: flex-start;
        }

        .dept-sidebar {
          width: 190px;
          flex-shrink: 0;
          display: flex;
          flex-direction: column;
          gap: 6px;
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-radius: var(--radius-lg);
          padding: 14px;
          position: sticky;
          top: 20px;
        }

        .dept-sidebar-title {
          font-size: 11px;
          font-weight: 700;
          color: var(--text-muted);
          letter-spacing: 0.5px;
          margin-bottom: 8px;
          line-height: 1.4;
        }

        .dept-btn {
          text-align: left;
          padding: 9px 12px;
          border-radius: var(--radius-sm);
          border: 1px solid var(--border-color);
          background: var(--bg-secondary);
          color: var(--text-secondary);
          font-size: 12.5px;
          font-weight: 500;
          cursor: pointer;
          transition: var(--transition);
        }

        .dept-btn:hover {
          border-color: var(--border-focus);
          color: var(--text-primary);
        }

        .dept-btn.active {
          background: var(--accent);
          border-color: var(--accent);
          color: #ffffff;
          font-weight: 700;
        }

        .dashboard-sections {
          flex: 1;
          min-width: 0;
        }

        .dashboard-section-title {
          font-size: 13px;
          font-weight: 700;
          color: var(--text-accent);
          letter-spacing: 0.5px;
          text-transform: uppercase;
          background: var(--accent-light);
          border: 1px solid var(--border-color);
          border-radius: var(--radius-sm);
          padding: 8px 14px;
          margin-bottom: 14px;
        }

        .dashboard-charts-grid {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 16px;
        }

        .chart-card {
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-radius: var(--radius-lg);
          box-shadow: var(--shadow-sm);
          padding: 16px 18px;
          display: flex;
          flex-direction: column;
        }

        .chart-card-header {
          margin-bottom: 8px;
        }

        .chart-card-title {
          font-size: 13.5px;
          font-weight: 700;
          color: var(--text-primary);
        }

        .chart-card-subtitle {
          font-size: 11.5px;
          color: var(--text-muted);
          margin-top: 2px;
        }

        .chart-card-body {
          height: 260px;
        }

        @media (max-width: 991px) {
          .dashboard-body-layout {
            flex-direction: column;
          }
          .dept-sidebar {
            width: 100%;
            flex-direction: row;
            flex-wrap: wrap;
            position: static;
          }
          .dashboard-charts-grid {
            grid-template-columns: 1fr;
          }
        }
      `}</style>
    </div>
  );
}

export default DashboardPage;
