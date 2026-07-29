import { NavLink, useNavigate } from 'react-router-dom';
import { Drawer } from 'antd';
import { authAPI } from '../api/client.js';
import { useToast } from '../hooks/useToast.js';
import { usePermissions } from '../contexts/PermissionsContext.jsx';

const navSections = [
  {
    title: 'Thống kê',
    items: [
      { path: '/dashboard', label: 'Dashboard', icon: 'fa-solid fa-chart-line', permission: 'dashboard::read' }
    ]
  },
  {
    title: 'Điều hành',
    items: [
      { path: '/meetings', label: 'Lịch họp & Video Call', icon: 'fa-solid fa-video' }
    ]
  },
  {
    title: 'Lương',
    items: [
      { path: '/payroll-run', label: 'Tính lương tháng', icon: 'fa-solid fa-calculator', permission: 'salary::read' },
      { path: '/payroll-formulas', label: 'Công thức lương', icon: 'fa-solid fa-square-root-variable', permission: 'formulaDynamic::read' }
    ]
  },
  {
    title: 'Chấm công',
    items: [
      { path: '/face-scan', label: 'Chấm công khuôn mặt', icon: 'fa-solid fa-camera', permission: 'dailyTimekeeping::edit' },
      { path: '/timesheets', label: 'Bảng chấm công', icon: 'fa-solid fa-calendar-days', permission: 'timekeepingSheet::read' },
      { path: '/attendance-logs', label: 'Nhật ký chấm công', icon: 'fa-solid fa-list-check', permission: 'dailyTimekeeping::read' },
      { path: '/explanation-requests', label: 'Giải trình bù công', icon: 'fa-solid fa-file-pen' },
      { path: '/ot-requests', label: 'Đăng ký làm thêm (OT)', icon: 'fa-solid fa-clock' }
    ]
  },
  {
    title: 'Quản lý nhân sự',
    items: [
      { path: '/departments', label: 'Phòng ban', icon: 'fa-solid fa-building', permission: 'department::read' },
      { path: '/employees', label: 'Nhân viên', icon: 'fa-solid fa-users', permission: 'userProfile::read' }
    ]
  },
  {
    title: 'Đánh giá',
    items: [
      { path: '/competencies', label: 'Năng lực', icon: 'fa-solid fa-award', permission: 'jobCapability::read' },
      { path: '/job-positions', label: 'Vị trí công việc', icon: 'fa-solid fa-briefcase', permission: 'jobPosition::read' }
    ]
  },
  {
    title: 'Hệ thống',
    items: [
      { path: '/contracts', label: 'Hợp đồng', icon: 'fa-solid fa-file-contract', permission: 'userProfile::read' },
      { path: '/roles-permissions', label: 'Phân quyền', icon: 'fa-solid fa-user-shield', permission: 'role::read' },
      { path: '/settings', label: 'Cài đặt', icon: 'fa-solid fa-gear', permission: 'setting::read' }
    ]
  },
  {
    title: 'Hỗ trợ',
    items: [
      { path: '/pdf-tools', label: 'Công cụ PDF', icon: 'fa-solid fa-file-pdf', permission: 'storage::edit' }
    ]
  }
];

export function Sidebar({ mobileOpen, onClose }) {
  const { can, loaded } = usePermissions();

  const hasPermission = (perm) => {
    if (!perm) return true;
    if (!loaded) return false; // ẩn menu cho đến khi load xong quyền hạn
    if (perm === 'dashboard::read') return true;
    return can(perm);
  };

  const handleLogout = async () => {
    try {
      await authAPI.logout();
    } catch {
      // ignore
    }
    localStorage.removeItem('auth_token');
    toast.success('Đã đăng xuất', 'Hẹn gặp lại!');
    if (onClose) onClose();
    navigate('/login');
  };

  const handleNavClick = () => {
    if (onClose) onClose();
  };

  const sidebarContent = (
    <div className="sidebar-inner">
      <div className="sidebar-logo">
        <div className="sidebar-logo-text">
          <span>CMS</span> Payroll
        </div>
      </div>

      <nav className="sidebar-nav">
        {navSections.map((section) => {
          const visibleItems = section.items.filter(item => hasPermission(item.permission));
          if (visibleItems.length === 0) return null;

          return (
            <div key={section.title} className="nav-section" style={{ marginBottom: 12 }}>
              <div className="nav-section-label">{section.title}</div>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                {visibleItems.map((item) => (
                  <NavLink
                    key={item.path}
                    to={item.path}
                    onClick={handleNavClick}
                    className={({ isActive }) => `nav-item${isActive ? ' active' : ''}`}
                  >
                    {item.icon && <i className={item.icon} style={{ width: '20px', fontSize: '14px', marginRight: '6px' }}></i>}
                    <span className="nav-label">{item.label}</span>
                  </NavLink>
                ))}
              </div>
            </div>
          );
        })}
      </nav>

      <div className="sidebar-footer">
        <button className="sidebar-logout" onClick={handleLogout}>
          <span>Đăng xuất</span>
        </button>
      </div>
    </div>
  );

  return (
    <>
      {/* Desktop Sidebar */}
      <aside className="sidebar desktop-sidebar">
        {sidebarContent}
      </aside>

      {/* Mobile Drawer Sidebar */}
      <Drawer
        open={mobileOpen}
        onClose={onClose}
        placement="left"
        width={260}
        closable={false}
        styles={{ body: { padding: 0, background: '#082920' } }}
        className="mobile-sidebar-drawer"
      >
        <div className="sidebar mobile-sidebar-content">
          {sidebarContent}
        </div>
      </Drawer>

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

        .mobile-sidebar-content {
          position: relative;
          width: 100%;
          height: 100%;
        }

        .sidebar-inner {
          display: flex;
          flex-direction: column;
          height: 100%;
          width: 100%;
        }

        @media (max-width: 991px) {
          .desktop-sidebar {
            display: none !important;
          }
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
    </>
  );
}

export default Sidebar;

