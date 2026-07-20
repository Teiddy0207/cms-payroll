import { NavLink, useNavigate } from 'react-router-dom';
import { authAPI } from '../api/client.js';
import { useToast } from '../hooks/useToast.js';

const navSections = [
  {
    title: 'Thống kê',
    items: [
      { path: '/dashboard', label: 'Dashboard' }
    ]
  },
  {
    title: 'Lương',
    items: [
      // { path: '/payroll', label: 'Lương tham khảo' },
      { path: '/payroll-run', label: 'Tính lương tháng' },
      { path: '/payroll-formulas', label: 'Công thức lương' }
    ]
  },
  {
    title: 'Chấm công',
    items: [
      { path: '/face-scan', label: 'Chấm công khuôn mặt' },
      { path: '/timesheets', label: 'Bảng chấm công' },
      { path: '/attendance-logs', label: 'Nhật ký chấm công' },
      { path: '/explanation-requests', label: 'Giải trình bù công' },
      { path: '/ot-requests', label: 'Đăng ký làm thêm (OT)' }
    ]
  },
  {
    title: 'Quản lý nhân sự',
    items: [
      { path: '/departments', label: 'Phòng ban' },
      { path: '/employees', label: 'Nhân viên' }
    ]
  },
  {
    title: 'Đánh giá',
    items: [
      { path: '/competencies', label: 'Năng lực' },
      { path: '/job-positions', label: 'Vị trí công việc' }
    ]
  },
  {
    title: 'Hệ thống',
    items: [
      { path: '/contracts', label: 'Hợp đồng' },
      { path: '/settings', label: 'Cài đặt' }
    ]
  },
  {
    title: 'Hỗ trợ',
    items: [
      { path: '/pdf-tools', label: 'Công cụ PDF'}
    ]
  }
];

export function Sidebar() {
  const navigate = useNavigate();
  const toast = useToast();

  const handleLogout = async () => {
    try {
      await authAPI.logout();
    } catch {
      // ignore
    }
    localStorage.removeItem('auth_token');
    toast.success('Đã đăng xuất', 'Hẹn gặp lại!');
    navigate('/login');
  };

  return (
    <aside className="sidebar">
      <div className="sidebar-logo">
        <div className="sidebar-logo-text">
          <span>CMS</span> Payroll
        </div>
      </div>

      <nav className="sidebar-nav">
        {navSections.map((section) => (
          <div key={section.title} className="nav-section" style={{ marginBottom: 12 }}>
            <div className="nav-section-label">{section.title}</div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
              {section.items.map((item) => (
                <NavLink
                  key={item.path}
                  to={item.path}
                  className={({ isActive }) => `nav-item${isActive ? ' active' : ''}`}
                >
                  <span className="nav-label">{item.label}</span>
                </NavLink>
              ))}
            </div>
          </div>
        ))}
      </nav>

      <div className="sidebar-footer">
        <button className="sidebar-logout" onClick={handleLogout}>
          <span>Đăng xuất</span>
        </button>
      </div>

      <style>{`
        .sidebar {
          position: fixed;
          left: 0;
          top: 0;
          bottom: 0;
          width: var(--sidebar-width);
          background: #082920;
          border-right: 1px solid rgba(255, 255, 255, 0.06);
          display: flex;
          flex-direction: column;
          z-index: 100;
          overflow: hidden;
        }

        .sidebar-logo {
          display: flex;
          align-items: center;
          gap: 12px;
          padding: 20px 20px 18px;
          border-bottom: 1px solid rgba(255, 255, 255, 0.06);
          flex-shrink: 0;
        }

        .sidebar-logo-icon {
          width: 38px;
          height: 38px;
          background: var(--accent);
          border-radius: var(--radius-md);
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 18px;
          box-shadow: 0 0 16px rgba(16, 185, 129, 0.3);
          flex-shrink: 0;
        }

        .sidebar-logo-text {
          font-size: 17px;
          font-weight: 800;
          color: #ffffff;
          letter-spacing: -0.5px;
        }

        .sidebar-logo-text span {
          color: var(--accent);
        }

        .sidebar-nav {
          flex: 1;
          padding: 16px 10px;
          overflow-y: auto;
          display: flex;
          flex-direction: column;
          gap: 2px;
        }

        .nav-section-label {
          font-size: 10.5px;
          font-weight: 700;
          color: #638c82;
          letter-spacing: 1px;
          padding: 4px 10px 10px;
          text-transform: uppercase;
        }

        .nav-item {
          display: flex;
          align-items: center;
          gap: 10px;
          padding: 10px 12px;
          border-radius: var(--radius-md);
          color: #bce3db;
          text-decoration: none;
          font-size: 13.5px;
          font-weight: 500;
          transition: var(--transition);
          position: relative;
        }

        .nav-item:hover {
          background: rgba(255, 255, 255, 0.05);
          color: #ffffff;
        }

        .nav-item.active {
          background: rgba(16, 185, 129, 0.15);
          color: #10b981;
        }

        .nav-item.active::before {
          content: '';
          position: absolute;
          left: 0;
          top: 20%;
          bottom: 20%;
          width: 3px;
          background: #10b981;
          border-radius: 0 2px 2px 0;
        }

        .nav-icon {
          width: 22px;
          text-align: center;
          font-size: 15px;
          flex-shrink: 0;
        }

        .nav-label {
          flex: 1;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .sidebar-footer {
          padding: 14px 10px;
          border-top: 1px solid rgba(255, 255, 255, 0.06);
          flex-shrink: 0;
        }

        .sidebar-logout {
          display: flex;
          align-items: center;
          gap: 10px;
          width: 100%;
          padding: 10px 12px;
          border-radius: var(--radius-md);
          background: transparent;
          border: none;
          color: #638c82;
          font-size: 13.5px;
          font-weight: 500;
          cursor: pointer;
          transition: var(--transition);
        }

        .sidebar-logout:hover {
          background: rgba(220, 38, 38, 0.15);
          color: #ef4444;
        }
      `}</style>
    </aside>
  );
}

export default Sidebar;
