import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

export default function NotificationBell() {
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);
  const [notifications, setNotifications] = useState([]);

  useEffect(() => {
    const token = localStorage.getItem('auth_token');
    if (!token) return;

    const baseUrl = import.meta.env.VITE_API_BASE_URL || '/api/v1';
    const streamUrl = `${baseUrl}/meetings/notifications/stream?token=${encodeURIComponent(token)}`;
    const eventSource = new EventSource(streamUrl);

    eventSource.onmessage = (event) => {
      try {
        const notif = JSON.parse(event.data);
        if (notif && notif.id) {
          const newNotif = {
            id: notif.id,
            title: 'Lịch họp & Video Call',
            message: notif.text,
            time: notif.time,
            read: notif.read,
            meetingId: notif.meeting_id
          };
          setNotifications(prev => [newNotif, ...prev]);
          window.dispatchEvent(new CustomEvent('meeting_updated'));
        }
      } catch (err) {
        console.error("Failed to parse SSE notification:", err);
      }
    };

    eventSource.onerror = (err) => {
      console.warn("Notification SSE connection closed or failed, retrying...");
    };

    return () => {
      eventSource.close();
    };
  }, []);

  const unreadCount = notifications.filter(n => !n.read).length;

  return (
    <div style={{ position: 'relative' }}>
      <button
        onClick={() => setOpen(!open)}
        title="Thông báo cuộc họp real-time"
        style={{
          background: 'rgba(241, 245, 249, 0.9)',
          border: '1px solid #cbd5e1',
          borderRadius: '20px',
          padding: '6px 14px',
          fontSize: '13px',
          cursor: 'pointer',
          position: 'relative',
          display: 'flex',
          alignItems: 'center',
          gap: '6px',
          color: '#0f172a',
          fontWeight: '600',
          boxShadow: '0 2px 4px rgba(0,0,0,0.05)'
        }}
      >
        <i className="fa-solid fa-bell" style={{ color: '#0284c7', fontSize: '15px' }}></i>
        <span>Thông báo</span>
        {unreadCount > 0 && (
          <span
            style={{
              background: '#ef4444',
              color: '#fff',
              borderRadius: '10px',
              padding: '2px 8px',
              fontSize: '10.5px',
              fontWeight: '700',
              marginLeft: '2px'
            }}
          >
            {unreadCount}
          </span>
        )}
      </button>

      {open && (
        <div
          style={{
            position: 'absolute',
            right: 0,
            top: '40px',
            width: '320px',
            maxWidth: 'calc(100vw - 32px)',
            background: '#ffffff',
            borderRadius: '12px',
            boxShadow: '0 10px 25px rgba(0,0,0,0.15)',
            border: '1px solid #e2e8f0',
            zIndex: 1000,
            overflow: 'hidden'
          }}
        >
          <div style={{ padding: '12px 16px', background: '#f8fafc', borderBottom: '1px solid #e2e8f0', fontWeight: '700', fontSize: '13px', color: '#0f172a' }}>
            Thông báo hệ thống
          </div>
          <div style={{ maxHeight: '300px', overflowY: 'auto' }}>
            {notifications.map(n => (
              <div
                key={n.id}
                onClick={() => {
                  setNotifications(notifications.map(item => item.id === n.id ? { ...item, read: true } : item));
                  setOpen(false);
                  if (n.meetingId) {
                    navigate(`/meetings/room/${n.meetingId}`);
                  } else {
                    navigate('/meetings');
                  }
                }}
                style={{
                  padding: '12px 16px',
                  borderBottom: '1px solid #f1f5f9',
                  cursor: 'pointer',
                  background: n.read ? '#fff' : '#f0fdf4',
                  transition: 'background 0.2s'
                }}
              >
                <div style={{ fontSize: '13px', fontWeight: '600', color: '#0f172a', marginBottom: '2px', display: 'flex', alignItems: 'center', gap: '6px' }}>
                  <i className="fa-solid fa-calendar-days" style={{ color: '#0284c7' }}></i>
                  {n.title}
                </div>
                <div style={{ fontSize: '12px', color: '#475569', lineHeight: '1.3' }}>{n.message}</div>
                <div style={{ fontSize: '10px', color: '#94a3b8', marginTop: '4px' }}>{n.time}</div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
