import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { meetingsAPI, employeesAPI } from '../api/client.js';

export default function VideoMeetingRoom() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [micOn, setMicOn] = useState(true);
  const [camOn, setCamOn] = useState(true);
  const [sharingScreen, setSharingScreen] = useState(false);
  const [meeting, setMeeting] = useState(null);
  const [employees, setEmployees] = useState([]);
  const [participants, setParticipants] = useState([]);
  const [loading, setLoading] = useState(true);

  const parseToken = () => {
    try {
      const token = localStorage.getItem('auth_token');
      if (!token) return null;
      const base64Url = token.split('.')[1];
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
      }).join(''));
      return JSON.parse(jsonPayload);
    } catch {
      return null;
    }
  };

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        const [meetingRes, employeesRes] = await Promise.all([
          meetingsAPI.get(id),
          employeesAPI.list({ page_size: 200 })
        ]);
        
        if (meetingRes && meetingRes.data) {
          setMeeting(meetingRes.data.data || meetingRes.data);
        }
        
        let eList = [];
        if (employeesRes && employeesRes.data && employeesRes.data.data) {
          if (Array.isArray(employeesRes.data.data.items)) {
            eList = employeesRes.data.data.items;
          }
        }
        setEmployees(eList);
      } catch (err) {
        console.error("Error loading meeting room details:", err);
      } finally {
        setLoading(false);
      }
    };
    loadData();
  }, [id]);

  useEffect(() => {
    if (!meeting) return;

    const payload = parseToken();
    const currentUserID = payload?.user_id;

    const resolveUserName = (userId) => {
      if (!userId || !Array.isArray(employees)) return 'Lãnh đạo / Hệ thống';
      const emp = employees.find(e => e && e.user_id === userId);
      return emp ? emp.full_name : 'Lãnh đạo / Hệ thống';
    };

    const hostIsSelf = meeting.host_id === currentUserID;
    const hostName = resolveUserName(meeting.host_id);
    
    const hostPart = {
      id: meeting.host_id,
      name: hostName + (hostIsSelf ? ' (Bạn - Lãnh đạo)' : ' (Lãnh đạo)'),
      role: 'Host',
      isSelf: hostIsSelf
    };

    const attendeeParts = Array.isArray(meeting.attendees) ? meeting.attendees.map((att, idx) => {
      const isSelf = att.user_id === currentUserID;
      const name = resolveUserName(att.user_id) + (isSelf ? ' (Bạn)' : '');
      return {
        id: att.user_id || idx,
        name: name,
        role: 'Attendee',
        isSelf: isSelf
      };
    }) : [];

    // Filter out duplicates (if host is also listed in attendees)
    const list = [hostPart];
    attendeeParts.forEach(att => {
      if (att.id !== hostPart.id) {
        list.push(att);
      }
    });

    setParticipants(list);
  }, [meeting, employees]);

  if (loading) {
    return (
      <div style={{ background: '#0f172a', minHeight: 'calc(100vh - 80px)', color: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <h3>Đang kết nối phòng họp</h3>
      </div>
    );
  }

  if (!meeting) {
    return (
      <div style={{ background: '#0f172a', minHeight: 'calc(100vh - 80px)', color: '#fff', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: '16px' }}>
        <h3>Không tìm thấy phòng họp hoặc phòng họp đã kết thúc.</h3>
        <button
          onClick={() => navigate('/meetings')}
          style={{ background: '#0284c7', color: '#fff', border: 'none', padding: '8px 16px', borderRadius: '6px', cursor: 'pointer' }}
        >
          Quay lại danh sách
        </button>
      </div>
    );
  }

  return (
    <div style={{ background: '#0f172a', minHeight: 'calc(100vh - 80px)', color: '#fff', padding: '20px', borderRadius: '12px', display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
      {/* Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid #1e293b', paddingBottom: '16px' }}>
        <div>
          <h2 style={{ fontSize: '18px', fontWeight: '700', margin: 0, color: '#f8fafc', display: 'flex', alignItems: 'center' }}>
            <i className="fa-solid fa-video" style={{ marginRight: '8px', color: '#38bdf8' }}></i>
            Phòng họp: {meeting.title || 'Họp Lãnh Đạo'}
          </h2>
          <span style={{ fontSize: '12px', color: '#10b981', background: 'rgba(16,185,129,0.1)', padding: '2px 8px', borderRadius: '4px', marginTop: '4px', display: 'inline-flex', alignItems: 'center', gap: '6px' }}>
            <i className="fa-solid fa-circle" style={{ fontSize: '8px', color: '#10b981' }}></i> Live Streaming Connected
          </span>
        </div>
        <button
          onClick={() => navigate('/meetings')}
          style={{ background: '#ef4444', color: '#fff', border: 'none', padding: '8px 16px', borderRadius: '6px', fontWeight: '600', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '6px' }}
        >
          <i className="fa-solid fa-phone-slash"></i> Rời phòng họp
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
                <div style={{ width: '64px', height: '64px', background: '#334155', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 8px auto' }}>
                  <i className="fa-solid fa-video-slash" style={{ fontSize: '20px' }}></i>
                </div>
                Camera đã tắt
              </div>
            ) : sharingScreen && p.isSelf ? (
              <div style={{ textAlign: 'center', color: '#38bdf8' }}>
                <i className="fa-solid fa-desktop" style={{ fontSize: '36px', display: 'block', margin: '0 auto 12px auto' }}></i>
                <div style={{ fontWeight: '600' }}>Bạn đang chia sẻ màn hình...</div>
              </div>
            ) : (
              <div style={{ textAlign: 'center' }}>
                <div style={{ width: '72px', height: '72px', background: '#0284c7', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', margin: '0 auto 12px auto' }}>
                  <i className="fa-solid fa-user" style={{ fontSize: '28px' }}></i>
                </div>
                <div style={{ fontSize: '14px', fontWeight: '600' }}>{p.name}</div>
                <div style={{ fontSize: '11px', color: '#94a3b8' }}>{p.role}</div>
              </div>
            )}

            <div style={{ position: 'absolute', bottom: '12px', left: '12px', background: 'rgba(0,0,0,0.6)', padding: '4px 10px', borderRadius: '6px', fontSize: '12px', display: 'flex', alignItems: 'center', gap: '6px' }}>
              <span>{p.name}</span>
              {p.isSelf && (micOn ? <i className="fa-solid fa-microphone" style={{ color: '#10b981' }}></i> : <i className="fa-solid fa-microphone-slash" style={{ color: '#ef4444' }}></i>)}
            </div>
          </div>
        ))}
      </div>

      {/* Control Bar */}
      <div style={{ background: '#1e293b', padding: '16px', borderRadius: '12px', display: 'flex', justifyContent: 'center', gap: '16px', alignItems: 'center' }}>
        <button
          onClick={() => setMicOn(!micOn)}
          style={{ background: micOn ? '#334155' : '#ef4444', color: '#fff', border: 'none', padding: '10px 20px', borderRadius: '8px', fontWeight: '600', cursor: 'pointer', display: 'flex', alignItems: 'center' }}
        >
          {micOn ? (
            <>
              <i className="fa-solid fa-microphone-slash" style={{ marginRight: '8px' }}></i> Tắt Micro
            </>
          ) : (
            <>
              <i className="fa-solid fa-microphone" style={{ marginRight: '8px' }}></i> Bật Micro
            </>
          )}
        </button>

        <button
          onClick={() => setCamOn(!camOn)}
          style={{ background: camOn ? '#334155' : '#ef4444', color: '#fff', border: 'none', padding: '10px 20px', borderRadius: '8px', fontWeight: '600', cursor: 'pointer', display: 'flex', alignItems: 'center' }}
        >
          {camOn ? (
            <>
              <i className="fa-solid fa-video-slash" style={{ marginRight: '8px' }}></i> Tắt Camera
            </>
          ) : (
            <>
              <i className="fa-solid fa-video" style={{ marginRight: '8px' }}></i> Bật Camera
            </>
          )}
        </button>

        <button
          onClick={() => setSharingScreen(!sharingScreen)}
          style={{ background: sharingScreen ? '#0284c7' : '#334155', color: '#fff', border: 'none', padding: '10px 20px', borderRadius: '8px', fontWeight: '600', cursor: 'pointer', display: 'flex', alignItems: 'center' }}
        >
          {sharingScreen ? (
            <>
              <i className="fa-solid fa-circle-stop" style={{ marginRight: '8px' }}></i> Dừng chia sẻ
            </>
          ) : (
            <>
              <i className="fa-solid fa-desktop" style={{ marginRight: '8px' }}></i> Chia sẻ màn hình
            </>
          )}
        </button>
      </div>
    </div>
  );
}
