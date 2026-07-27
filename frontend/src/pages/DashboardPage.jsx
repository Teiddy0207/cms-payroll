import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { departmentsAPI, jobPositionsAPI, employeesAPI, contractsAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Badge } from '../components/Badge.jsx';

export function DashboardPage() {
  const navigate = useNavigate();
  const [stats, setStats] = useState({ employees: 0, contracts: 0, positions: 0, departments: 0 });
  const [recentEmployees, setRecentEmployees] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [deptRes, posRes, empRes, ctRes] = await Promise.allSettled([
        departmentsAPI.list({ page_size: 1 }),
        jobPositionsAPI.list({ page_size: 1 }),
        employeesAPI.list({ page_size: 8 }),
        contractsAPI.list({ page_size: 1 }),
      ]);

      const deptTotal = deptRes.status === 'fulfilled' ? deptRes.value.data?.data?.total_items || 0 : 0;
      const posTotal = posRes.status === 'fulfilled' ? posRes.value.data?.data?.total_items || 0 : 0;
      const empData = empRes.status === 'fulfilled' ? empRes.value.data?.data : null;
      const ctTotal = ctRes.status === 'fulfilled' ? ctRes.value.data?.data?.total_items || 0 : 0;

      setStats({
        employees: empData?.total_items || 0,
        contracts: ctTotal,
        positions: posTotal,
        departments: deptTotal,
      });
      setRecentEmployees(empData?.items || []);
    } catch {
      // fallback
    } finally {
      setLoading(false);
    }
  };

  const statCards = [
    { label: 'Tổng nhân viên', value: stats.employees, color: '#4f8ef7', path: '/employees' },
    { label: 'Hợp đồng', value: stats.contracts, color: '#22c55e', path: '/contracts' },
    { label: 'Lịch họp & Video Call', value: 'Trực tuyến', color: '#10b981', path: '/meetings' },
    { label: 'Vị trí công việc', value: stats.positions, color: '#f59e0b', path: '/job-positions' },
    { label: 'Phòng ban', value: stats.departments, color: '#06b6d4', path: '/departments' },
  ];

  const columns = [
    {
      key: 'code',
      title: 'Mã NV',
      render: (val) => <span className="font-mono text-accent">{val || '—'}</span>
    },
    {
      key: 'full_name',
      title: 'Họ và tên',
      render: (val, row) => (
        <div className="avatar-cell">
          <div className="avatar-info">
            <div className="name">{val || '—'}</div>
            <div className="code">{row.phone || 'Chưa có SĐT'}</div>
          </div>
        </div>
      )
    },
    {
      key: 'gender',
      title: 'Giới tính',
      render: (val) => val
        ? <Badge color={val === 'male' ? 'blue' : 'cyan'}>{val === 'male' ? 'Nam' : 'Nữ'}</Badge>
        : <Badge color="gray">—</Badge>
    },
    {
      key: 'job_position',
      title: 'Vị trí',
      render: (val, row) => {
        const pos = val || row.position;
        return pos ? <Badge color="yellow">{pos.name || pos}</Badge> : <span className="text-muted">—</span>;
      }
    },
  ];

  return (
    <div>
      {/* Stats */}
      <div className="grid-4" style={{ marginBottom: 28 }}>
        {statCards.map((card) => (
          <div
            key={card.label}
            className="stat-card"
            style={{ cursor: 'pointer' }}
            onClick={() => navigate(card.path)}
          >
            <div className="stat-info">
              {loading ? (
                <div style={{ height: 32, display: 'flex', alignItems: 'center' }}>
                  <span className="spinner spinner-sm" />
                </div>
              ) : (
                <h3 style={{ color: card.color }}>{card.value.toLocaleString('vi-VN')}</h3>
              )}
              <p>{card.label}</p>
            </div>
          </div>
        ))}
      </div>

      {/* Recent Employees */}
      <div className="table-container">
        <div className="table-toolbar">
          <div>
            <div style={{ fontWeight: 600, fontSize: 15 }}>Nhân viên gần nhất</div>
            <div className="text-muted" style={{ fontSize: 12.5 }}>8 nhân viên được thêm gần đây nhất</div>
          </div>
          <button className="btn btn-primary btn-sm" onClick={() => navigate('/employees')}>
            Xem tất cả →
          </button>
        </div>
        <Table
          columns={columns}
          data={recentEmployees}
          loading={loading}
          emptyMessage="Chưa có nhân viên nào"
        />
      </div>
    </div>
  );
}

export default DashboardPage;
