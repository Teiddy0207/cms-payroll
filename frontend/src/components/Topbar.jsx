import { useLocation } from 'react-router-dom';

const routeTitles = {
  '/dashboard': { title: 'Dashboard', subtitle: 'Tổng quan hệ thống' },
  '/employees': { title: 'Quản lý nhân viên', subtitle: 'Danh sách & thông tin nhân viên' },
  '/contracts': { title: 'Hợp đồng lao động', subtitle: 'Quản lý hợp đồng' },
  '/departments': { title: 'Phòng ban', subtitle: 'Quản lý cơ cấu tổ chức' },
  '/job-positions': { title: 'Vị trí công việc', subtitle: 'Quản lý vị trí & tiêu chuẩn' },
  '/job-standards': { title: 'Tiêu chuẩn nghề nghiệp', subtitle: 'Các tiêu chuẩn phụ cấp' },
  '/competencies': { title: 'Năng lực', subtitle: 'Quản lý năng lực nhân viên' },
  '/settings': { title: 'Cài đặt hệ thống', subtitle: 'Cấu hình và thông số' },
};

export function Topbar() {
  const location = useLocation();
  const info = routeTitles[location.pathname] || { title: 'CMS Payroll', subtitle: '' };
  const now = new Date();
  const dateStr = now.toLocaleDateString('vi-VN', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' });

  return (
    <header className="topbar">
      <div className="topbar-left">
        <h1 className="topbar-title">{info.title}</h1>
        {info.subtitle && <p className="topbar-subtitle">{info.subtitle}</p>}
      </div>
      <div className="topbar-right">
        <div className="topbar-date">{dateStr}</div>
        <div className="topbar-user">
          <div className="topbar-avatar">A</div>
          <div className="topbar-user-info">
            <span className="topbar-user-name">Admin</span>
            <span className="topbar-user-role">Quản trị viên</span>
          </div>
        </div>
      </div>

      <style>{`
        .topbar {
          position: fixed;
          top: 0;
          left: var(--sidebar-width);
          right: 0;
          height: var(--topbar-height);
          background: rgba(15, 17, 23, 0.92);
          backdrop-filter: blur(12px);
          border-bottom: 1px solid var(--border-color);
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 0 32px;
          z-index: 90;
          gap: 20px;
        }

        .topbar-left {
          display: flex;
          flex-direction: column;
          justify-content: center;
        }

        .topbar-title {
          font-size: 16px;
          font-weight: 700;
          color: var(--text-primary);
          letter-spacing: -0.2px;
          line-height: 1.2;
        }

        .topbar-subtitle {
          font-size: 11.5px;
          color: var(--text-muted);
          font-weight: 400;
          margin-top: 1px;
        }

        .topbar-right {
          display: flex;
          align-items: center;
          gap: 20px;
          flex-shrink: 0;
        }

        .topbar-date {
          font-size: 12px;
          color: var(--text-muted);
          font-weight: 500;
        }

        .topbar-user {
          display: flex;
          align-items: center;
          gap: 10px;
          cursor: pointer;
          padding: 5px 10px;
          border-radius: var(--radius-md);
          transition: var(--transition);
        }

        .topbar-user:hover {
          background: rgba(255, 255, 255, 0.05);
        }

        .topbar-avatar {
          width: 34px;
          height: 34px;
          background: var(--accent);
          border-radius: var(--radius-full);
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 14px;
          font-weight: 700;
          color: #fff;
          flex-shrink: 0;
        }

        .topbar-user-info {
          display: flex;
          flex-direction: column;
        }

        .topbar-user-name {
          font-size: 13px;
          font-weight: 600;
          color: var(--text-primary);
          line-height: 1.3;
        }

        .topbar-user-role {
          font-size: 11px;
          color: var(--text-muted);
        }

        @media (max-width: 768px) {
          .topbar {
            left: 0;
            padding: 0 16px;
          }
          .topbar-date { display: none; }
          .topbar-user-info { display: none; }
        }
      `}</style>
    </header>
  );
}

export default Topbar;
