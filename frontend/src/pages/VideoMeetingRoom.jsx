import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

export default function VideoMeetingRoom() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [micOn, setMicOn] = useState(true);
  const [camOn, setCamOn] = useState(true);
  const [sharingScreen, setSharingScreen] = useState(false);
  const [participants, setParticipants] = useState([
    { id: 1, name: 'Bạn (Lãnh đạo)', role: 'Host', isSelf: true },
    { id: 2, name: 'Nguyễn Văn A (IT)', role: 'Attendee', isSelf: false },
    { id: 3, name: 'Trần Thị B (Kế toán)', role: 'Attendee', isSelf: false }
  ]);

  return (
    <div style={{ background: '#0f172a', minHeight: 'calc(100vh - 80px)', color: '#fff', padding: '20px', borderRadius: '12px', display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
      {/* Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid #1e293b', paddingBottom: '16px' }}>
        <div>
          <h2 style={{ fontSize: '18px', fontWeight: '700', margin: 0, color: '#f8fafc' }}>
            🎥 Phòng họp Trực tuyến WebRTC (Room ID: {id})
          </h2>
          <span style={{ fontSize: '12px', color: '#10b981', background: 'rgba(16,185,129,0.1)', padding: '2px 8px', borderRadius: '4px', marginTop: '4px', display: 'inline-block' }}>
            🔴 Live Streaming Connected
          </span>
        </div>
        <button
          onClick={() => navigate('/meetings')}
          style={{ background: '#ef4444', color: '#fff', border: 'none', padding: '8px 16px', borderRadius: '6px', fontWeight: '600', cursor: 'pointer' }}
        >
          📞 Rời phòng họp
        </button>
      </div>

      {/* Grid Video Feeds */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '16px', margin: '20px 0', flex: 1 }}>
        {participants.map(p => (
          <div
            key={p.id}
            style={{
              background: '#1e293b',
              borderRadius: '12px',
              position: 'relative',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              minHeight: '220px',
              border: '1px solid #334155',
              overflow: 'hidden'
            }}
          >
            {p.isSelf && !camOn ? (
              <div style={{ textAlign: 'center', color: '#94a3b8' }}>
                <div style={{ width: '64px', height: '64px', background: '#334155', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '24px', margin: '0 auto 8px auto' }}>
                  📷
                </div>
                Camera đã tắt
              </div>
            ) : sharingScreen && p.isSelf ? (
              <div style={{ textAlign: 'center', color: '#38bdf8', fontWeight: '600' }}>
                🖥️ Bạn đang chia sẻ màn hình máy tính...
              </div>
            ) : (
              <div style={{ textAlign: 'center' }}>
                <div style={{ width: '72px', height: '72px', background: '#0284c7', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '28px', color: '#fff', fontWeight: '700', margin: '0 auto 12px auto' }}>
                  {p.name.charAt(0)}
                </div>
                <div style={{ fontSize: '14px', fontWeight: '600' }}>{p.name}</div>
                <div style={{ fontSize: '11px', color: '#94a3b8' }}>{p.role}</div>
              </div>
            )}

            <div style={{ position: 'absolute', bottom: '12px', left: '12px', background: 'rgba(0,0,0,0.6)', padding: '4px 10px', borderRadius: '6px', fontSize: '12px', display: 'flex', alignItems: 'center', gap: '6px' }}>
              <span>{p.name}</span>
              {p.isSelf && (micOn ? <span>🎙️</span> : <span style={{ color: '#ef4444' }}>🔇</span>)}
            </div>
          </div>
        ))}
      </div>

      {/* Control Bar */}
      <div style={{ background: '#1e293b', padding: '16px', borderRadius: '12px', display: 'flex', justifyContent: 'center', gap: '16px', alignItems: 'center' }}>
        <button
          onClick={() => setMicOn(!micOn)}
          style={{ background: micOn ? '#334155' : '#ef4444', color: '#fff', border: 'none', padding: '10px 20px', borderRadius: '8px', fontWeight: '600', cursor: 'pointer' }}
        >
          {micOn ? '🎙️ Tắt Micro' : '🔇 Bật Micro'}
        </button>

        <button
          onClick={() => setCamOn(!camOn)}
          style={{ background: camOn ? '#334155' : '#ef4444', color: '#fff', border: 'none', padding: '10px 20px', borderRadius: '8px', fontWeight: '600', cursor: 'pointer' }}
        >
          {camOn ? '📷 Tắt Camera' : '📷 Bật Camera'}
        </button>

        <button
          onClick={() => setSharingScreen(!sharingScreen)}
          style={{ background: sharingScreen ? '#0284c7' : '#334155', color: '#fff', border: 'none', padding: '10px 20px', borderRadius: '8px', fontWeight: '600', cursor: 'pointer' }}
        >
          {sharingScreen ? '🖥️ Dừng chia sẻ' : '🖥️ Chia sẻ màn hình'}
        </button>
      </div>
    </div>
  );
}
