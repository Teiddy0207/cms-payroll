import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

export default function Meetings() {
  const navigate = useNavigate();
  const [meetings, setMeetings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);

  // Form State
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [startTime, setStartTime] = useState('');
  const [endTime, setEndTime] = useState('');
  const [syncGoogle, setSyncGoogle] = useState(true);
  const [enableWebRTC, setEnableWebRTC] = useState(true);

  useEffect(() => {
    // Mock initial meetings for demo / immediate feedback
    const mockMeetings = [
      {
        id: 'meeting-1',
        title: 'Họp Chiến lược Quỹ lương Q3',
        description: 'Thống nhất cơ chế phụ cấp P1 & P2 cho các phòng ban kỹ thuật',
        start_time: new Date(Date.now() + 3600000).toISOString(),
        end_time: new Date(Date.now() + 7200000).toISOString(),
        host_name: 'Lãnh đạo (Giám đốc)',
        room_url: '/meetings/room-demo-123',
        status: 'SCHEDULED',
        google_event_id: 'gcal_89123789',
        attendees: [
          { user_name: 'Nguyễn Văn A (IT)', rsvp_status: 'ACCEPTED' },
          { user_name: 'Trần Thị B (Kế toán)', rsvp_status: 'PENDING' }
        ]
      }
    ];
    setMeetings(mockMeetings);
    setLoading(false);
  }, []);

  const handleCreateMeeting = (e) => {
    e.preventDefault();
    if (!title || !startTime || !endTime) return;

    const newMeeting = {
      id: `meeting-${Date.now()}`,
      title,
      description,
      start_time: new Date(startTime).toISOString(),
      end_time: new Date(endTime).toISOString(),
      host_name: 'Lãnh đạo',
      room_url: `/meetings/room-${Date.now()}`,
      status: 'SCHEDULED',
      google_event_id: syncGoogle ? `gcal_${Date.now()}` : null,
      attendees: [
        { user_name: 'Toàn bộ Nhân viên được chọn', rsvp_status: 'ACCEPTED' }
      ]
    };

    setMeetings([newMeeting, ...meetings]);
    setShowModal(false);
    // Reset form
    setTitle('');
    setDescription('');
    setStartTime('');
    setEndTime('');
  };

  const handleRSVP = (meetingId, status) => {
    setMeetings(meetings.map(m => {
      if (m.id === meetingId) {
        return {
          ...m,
          attendees: m.attendees.map(a => ({ ...a, rsvp_status: status }))
        };
      }
      return m;
    }));
  };

  return (
    <div style={{ padding: '24px', maxWidth: '1200px', margin: '0 auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <div>
          <h2 style={{ fontSize: '22px', fontWeight: '700', color: '#0f172a', margin: 0 }}>
            📅 Lịch họp Lãnh đạo & Phòng họp Video Call
          </h2>
          <p style={{ fontSize: '13px', color: '#64748b', margin: '4px 0 0 0' }}>
            Đặt lịch họp, tự động đồng bộ Google Calendar & tham gia phòng họp Video WebRTC trực tiếp trên Web.
          </p>
        </div>
        <button
          onClick={() => setShowModal(true)}
          style={{
            background: '#10b981',
            color: '#fff',
            border: 'none',
            padding: '10px 20px',
            borderRadius: '8px',
            fontWeight: '600',
            cursor: 'pointer',
            boxShadow: '0 4px 12px rgba(16, 185, 129, 0.25)'
          }}
        >
          ➕ Tạo cuộc họp mới
        </button>
      </div>

      {loading ? (
        <div>Đang tải danh sách cuộc họp...</div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(350px, 1fr))', gap: '20px' }}>
          {meetings.map((m) => (
            <div
              key={m.id}
              style={{
                background: '#ffffff',
                borderRadius: '12px',
                border: '1px solid #e2e8f0',
                padding: '20px',
                boxShadow: '0 2px 8px rgba(0,0,0,0.04)',
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'space-between'
              }}
            >
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
                  <span style={{ fontSize: '11px', fontWeight: '700', padding: '4px 8px', borderRadius: '4px', background: '#dcfce7', color: '#15803d' }}>
                    {m.status}
                  </span>
                  {m.google_event_id && (
                    <span style={{ fontSize: '11px', color: '#475569', display: 'flex', alignItems: 'center', gap: '4px' }}>
                      🟢 Google Schedule Synced
                    </span>
                  )}
                </div>

                <h3 style={{ fontSize: '16px', fontWeight: '700', color: '#0f172a', margin: '0 0 8px 0' }}>
                  {m.title}
                </h3>
                <p style={{ fontSize: '13px', color: '#475569', margin: '0 0 16px 0', lineHeight: '1.4' }}>
                  {m.description || 'Không có mô tả chi tiết.'}
                </p>

                <div style={{ fontSize: '12px', color: '#64748b', marginBottom: '12px' }}>
                  🕒 <strong>Thời gian:</strong> {new Date(m.start_time).toLocaleString('vi-VN')}
                </div>
              </div>

              <div style={{ borderTop: '1px solid #f1f5f9', paddingTop: '14px', marginTop: '14px', display: 'flex', gap: '10px' }}>
                <button
                  onClick={() => navigate(`/meetings/room/${m.id}`)}
                  style={{
                    flex: 1,
                    background: '#0284c7',
                    color: '#fff',
                    border: 'none',
                    padding: '8px 12px',
                    borderRadius: '6px',
                    fontWeight: '600',
                    fontSize: '13px',
                    cursor: 'pointer'
                  }}
                >
                  🎥 Vào Phòng Họp Video
                </button>
                <button
                  onClick={() => handleRSVP(m.id, 'ACCEPTED')}
                  style={{
                    background: '#ecfdf5',
                    color: '#047857',
                    border: '1px solid #a7f3d0',
                    padding: '8px 12px',
                    borderRadius: '6px',
                    fontWeight: '600',
                    fontSize: '12px',
                    cursor: 'pointer'
                  }}
                >
                  ✅ Đồng ý
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Modal Tạo Cuộc Họp */}
      {showModal && (
        <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 }}>
          <div style={{ background: '#fff', padding: '28px', borderRadius: '12px', width: '100%', maxWidth: '500px', boxShadow: '0 20px 25px -5px rgba(0,0,0,0.1)' }}>
            <h3 style={{ margin: '0 0 16px 0', fontSize: '18px', fontWeight: '700' }}>Tạo Cuộc Họp Lãnh Đạo Mới</h3>
            <form onSubmit={handleCreateMeeting}>
              <div style={{ marginBottom: '12px' }}>
                <label style={{ display: 'block', fontSize: '12px', fontWeight: '600', marginBottom: '4px' }}>Tiêu đề cuộc họp</label>
                <input
                  type="text"
                  required
                  value={title}
                  onChange={e => setTitle(e.target.value)}
                  style={{ width: '100%', padding: '8px 12px', borderRadius: '6px', border: '1px solid #cbd5e1' }}
                  placeholder="Vd: Họp đánh giá năng lực P2"
                />
              </div>

              <div style={{ marginBottom: '12px' }}>
                <label style={{ display: 'block', fontSize: '12px', fontWeight: '600', marginBottom: '4px' }}>Mô tả nội dung</label>
                <textarea
                  value={description}
                  onChange={e => setDescription(e.target.value)}
                  style={{ width: '100%', padding: '8px 12px', borderRadius: '6px', border: '1px solid #cbd5e1', height: '70px' }}
                  placeholder="Nội dung thảo luận chính..."
                />
              </div>

              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px', marginBottom: '12px' }}>
                <div>
                  <label style={{ display: 'block', fontSize: '12px', fontWeight: '600', marginBottom: '4px' }}>Bắt đầu</label>
                  <input
                    type="datetime-local"
                    required
                    value={startTime}
                    onChange={e => setStartTime(e.target.value)}
                    style={{ width: '100%', padding: '8px', borderRadius: '6px', border: '1px solid #cbd5e1' }}
                  />
                </div>
                <div>
                  <label style={{ display: 'block', fontSize: '12px', fontWeight: '600', marginBottom: '4px' }}>Kết thúc</label>
                  <input
                    type="datetime-local"
                    required
                    value={endTime}
                    onChange={e => setEndTime(e.target.value)}
                    style={{ width: '100%', padding: '8px', borderRadius: '6px', border: '1px solid #cbd5e1' }}
                  />
                </div>
              </div>

              <div style={{ marginBottom: '16px', display: 'flex', flexDirection: 'column', gap: '8px' }}>
                <label style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '13px', cursor: 'pointer' }}>
                  <input type="checkbox" checked={syncGoogle} onChange={e => setSyncGoogle(e.target.checked)} />
                  Đồng bộ Google Calendar (Google Schedule)
                </label>
                <label style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '13px', cursor: 'pointer' }}>
                  <input type="checkbox" checked={enableWebRTC} onChange={e => setEnableWebRTC(e.target.checked)} />
                  Tạo phòng họp Video WebRTC trực tuyến
                </label>
              </div>

              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px' }}>
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  style={{ padding: '8px 16px', borderRadius: '6px', border: '1px solid #cbd5e1', background: '#fff', cursor: 'pointer' }}
                >
                  Hủy
                </button>
                <button
                  type="submit"
                  style={{ padding: '8px 16px', borderRadius: '6px', border: 'none', background: '#10b981', color: '#fff', fontWeight: '600', cursor: 'pointer' }}
                >
                  Tạo cuộc họp
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
