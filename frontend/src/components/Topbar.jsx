import React, { useState, useEffect, useMemo } from 'react';
import { useLocation } from 'react-router-dom';
import NotificationBell from './NotificationBell.jsx';
import { Button } from 'antd';
import { employeesAPI } from '../api/client.js';

const routeTitles = {
  '/dashboard': { title: 'Dashboard', subtitle: 'Tổng quan hệ thống' },
  '/meetings': { title: 'Lịch họp & Video Call', subtitle: 'Quản lý cuộc họp & phòng họp' },
  '/employees': { title: 'Quản lý nhân viên', subtitle: 'Danh sách & thông tin nhân viên' },
  '/contracts': { title: 'Hợp đồng lao động', subtitle: 'Quản lý hợp đồng' },
  '/departments': { title: 'Phòng ban', subtitle: 'Quản lý cơ cấu tổ chức' },
  '/job-positions': { title: 'Vị trí công việc', subtitle: 'Quản lý vị trí & tiêu chuẩn' },
  '/competencies': { title: 'Năng lực', subtitle: 'Quản lý năng lực nhân viên' },
  '/settings': { title: 'Cài đặt hệ thống', subtitle: 'Cấu hình và thông số' },
  '/attendance-logs': { title: 'Nhật ký chấm công', subtitle: 'Xem lịch sử quét vân tay hoặc nhận dạng khuôn mặt' },
};

function parseJWT(token) {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64).split('').map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)).join('')
    );
    return JSON.parse(jsonPayload);
  } catch {
    return null;
  }
}

export function Topbar({ onToggleMobileMenu }) {
  const location = useLocation();
  const info = routeTitles[location.pathname] || { title: 'CMS Payroll', subtitle: '' };
  const now = new Date();
  const dateStr = now.toLocaleDateString('vi-VN', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' });

  const [employeeName, setEmployeeName] = useState('');

  useEffect(() => {
    const token = localStorage.getItem('auth_token');
    if (!token) return;
    const claims = parseJWT(token);
    if (!claims) return;

    const fetchProfile = async () => {
      try {
        const res = await employeesAPI.list({ page_size: 500 });
        let list = [];
        if (res && res.data && res.data.data) {
          if (Array.isArray(res.data.data.items)) {
            list = res.data.data.items;
          } else if (Array.isArray(res.data.data)) {
            list = res.data.data;
          }
        }
        const currentUserID = claims.user_id;
        const currentUsername = claims.username;
        const emp = list.find(e => e && (e.user_id === currentUserID || e.id === currentUserID || e.employee_code === currentUsername));
        if (emp && emp.full_name) {
          setEmployeeName(emp.full_name);
        }
      } catch (e) {
        console.warn("Could not fetch employee profile for Topbar:", e);
      }
    };

    fetchProfile();
  }, []);

  const userInfo = useMemo(() => {
    const token = localStorage.getItem('auth_token');
    if (!token) return { displayName: 'Admin', email: '', initial: 'A' };
    const claims = parseJWT(token);
    if (!claims) return { displayName: 'Admin', email: '', initial: 'A' };
    const displayName = employeeName || claims.full_name || claims.name || claims.username || claims.email || 'Người dùng';
    const initial = displayName.charAt(0).toUpperCase();
    return { displayName, email: claims.email || '', initial };
  }, [employeeName]);

  return (
    <header className="topbar">
      <div className="topbar-left">
        <Button
          type="text"
          className="mobile-menu-btn"
          icon={<i className="fa-solid fa-bars" style={{ fontSize: '18px', color: 'var(--text-primary)' }}></i>}
          onClick={onToggleMobileMenu}
        />
      </div>
      <div className="topbar-right">
        <NotificationBell />
        <div className="topbar-user">
          <div className="topbar-avatar">{userInfo.initial}</div>
          <div className="topbar-user-info">
            <span className="topbar-user-name">{userInfo.displayName}</span>
            {userInfo.email ? <span className="topbar-user-role">{userInfo.email}</span> : null}
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
          background: rgba(255, 255, 255, 0.85);
          backdrop-filter: blur(12px);
          border-bottom: 1px solid var(--border-color);
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 0 24px;
          z-index: 90;
          gap: 16px;
        }

        .topbar-left {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .mobile-menu-btn {
          display: none;
          align-items: center;
          justify-content: center;
          width: 36px;
          height: 36px;
          padding: 0;
        }

        .topbar-right {
          display: flex;
          align-items: center;
          gap: 16px;
          flex-shrink: 0;
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
          background: rgba(15, 23, 42, 0.04);
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

        @media (max-width: 991px) {
          .topbar {
            left: 0;
            padding: 0 16px;
          }
          .mobile-menu-btn {
            display: inline-flex;
          }
        }

        @media (max-width: 640px) {
          .topbar-user-info { display: none; }
        }
      `}</style>
    </header>
  );
}

export default Topbar;
