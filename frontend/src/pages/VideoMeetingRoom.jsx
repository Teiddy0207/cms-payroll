import React, { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { meetingsAPI, employeesAPI } from '../api/client.js';
import { useToast } from '../hooks/useToast.js';

// Remote Video Stream Helper Component
function RemoteVideo({ stream, name, style }) {
  const videoRef = useRef(null);

  useEffect(() => {
    if (videoRef.current && stream) {
      videoRef.current.srcObject = stream;
    }
  }, [stream]);

  return (
    <video
      ref={videoRef}
      autoPlay
      playsInline
      style={{
        width: '100%',
        height: '100%',
        objectFit: style?.objectFit || 'cover',
        borderRadius: style?.borderRadius || '14px'
      }}
    />
  );
}

export default function VideoMeetingRoom() {
  const { id } = useParams();
  const navigate = useNavigate();
  const toast = useToast();

  const [micOn, setMicOn] = useState(true);
  const [camOn, setCamOn] = useState(true);
  const [sharingScreen, setSharingScreen] = useState(false);
  const [remoteScreenPresenter, setRemoteScreenPresenter] = useState(null);
  const [meeting, setMeeting] = useState(null);
  const [employees, setEmployees] = useState([]);
  const [participants, setParticipants] = useState([]);
  const [remoteStreams, setRemoteStreams] = useState({});
  const [loading, setLoading] = useState(true);
  const [mediaError, setMediaError] = useState(null);
  const [copiedLink, setCopiedLink] = useState(false);

  const localVideoRef = useRef(null);
  const screenVideoRef = useRef(null);
  const localStreamRef = useRef(null);
  const screenStreamRef = useRef(null);
  
  const peerConnectionsRef = useRef({});
  const broadcastChannelRef = useRef(null);
  const eventSourceRef = useRef(null);

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

  const payload = parseToken();
  const currentUserID = payload?.user_id;

  const resolveUserName = (userId) => {
    if (!userId || !Array.isArray(employees)) return 'Thành viên';
    const emp = employees.find(e => e && (e.user_id === userId || e.id === userId));
    return emp ? emp.full_name : 'Thành viên';
  };

  const selfName = resolveUserName(currentUserID);

  // Load meeting & employees data
  useEffect(() => {
    const loadMeetingData = async () => {
      setLoading(true);
      try {
        const [meetingRes, empRes] = await Promise.all([
          meetingsAPI.get(id),
          employeesAPI.list({ page_size: 500 })
        ]);
        
        let mData = null;
        if (meetingRes && meetingRes.data) {
          mData = meetingRes.data.data || meetingRes.data;
        }
        setMeeting(mData);

        let eList = [];
        if (empRes && empRes.data && empRes.data.data) {
          if (Array.isArray(empRes.data.data.items)) {
            eList = empRes.data.data.items;
          } else if (Array.isArray(empRes.data.data)) {
            eList = empRes.data.data;
          }
        }
        setEmployees(eList);
      } catch (err) {
        console.error("Error loading video meeting room data:", err);
        setMediaError("Không thể tải thông tin cuộc họp nội bộ.");
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      loadMeetingData();
    }
  }, [id]);

  // Sync participants list
  useEffect(() => {
    if (!meeting) return;

    const resolveUserName = (userId) => {
      if (!userId || !Array.isArray(employees)) return 'Lãnh đạo / Hệ thống';
      const emp = employees.find(e => e && (e.user_id === userId || e.id === userId));
      return emp ? emp.full_name : 'Lãnh đạo / Hệ thống';
    };

    const hostIsSelf = meeting.host_id === currentUserID;
    const hostName = resolveUserName(meeting.host_id);
    
    const hostPart = {
      id: meeting.host_id,
      name: hostName + (hostIsSelf ? ' (Bạn - Host)' : ' (Host)'),
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

    const list = [hostPart];
    attendeeParts.forEach(att => {
      if (att.id !== hostPart.id) {
        list.push(att);
      }
    });

    setParticipants(list);
  }, [meeting, employees, currentUserID]);

  // Send WebRTC Signal helper
  const sendSignal = async (signalData) => {
    try {
      await meetingsAPI.signal(signalData);
    } catch (e) {
      console.warn("Failed to send WebRTC HTTP signal:", e);
    }
    if (broadcastChannelRef.current) {
      try {
        broadcastChannelRef.current.postMessage(signalData);
      } catch {
        // ignore
      }
    }
  };

  // Initialize WebRTC Local Stream & Signaling
  useEffect(() => {
    let isMounted = true;

    const removePeer = (targetUserId) => {
      if (peerConnectionsRef.current[targetUserId]) {
        try { peerConnectionsRef.current[targetUserId].close(); } catch {}
        delete peerConnectionsRef.current[targetUserId];
      }
      setRemoteStreams(prev => {
        const next = { ...prev };
        delete next[targetUserId];
        return next;
      });
      setRemoteScreenPresenter(prev => (prev?.userId === targetUserId ? null : prev));
    };

    const createPeerConnection = (targetUserId) => {
      if (peerConnectionsRef.current[targetUserId]) {
        return peerConnectionsRef.current[targetUserId];
      }

      const pc = new RTCPeerConnection({
        iceServers: [
          { urls: 'stun:stun.l.google.com:19302' },
          { urls: 'stun:stun1.l.google.com:19302' }
        ]
      });

      const currentTrack = (sharingScreen && screenStreamRef.current) 
        ? screenStreamRef.current.getVideoTracks()[0] 
        : localStreamRef.current?.getVideoTracks()[0];

      if (currentTrack) {
        pc.addTrack(currentTrack, localStreamRef.current || screenStreamRef.current);
      }
      if (localStreamRef.current?.getAudioTracks()[0]) {
        pc.addTrack(localStreamRef.current.getAudioTracks()[0], localStreamRef.current);
      }

      pc.ontrack = (event) => {
        if (event.streams && event.streams[0]) {
          console.log(`[WebRTC] Received remote stream for participant: ${targetUserId}`);
          setRemoteStreams(prev => ({ ...prev, [targetUserId]: event.streams[0] }));
        }
      };

      pc.onicecandidate = (event) => {
        if (event.candidate) {
          sendSignal({
            type: 'candidate',
            meeting_id: id,
            from_user_id: currentUserID,
            to_user_id: targetUserId,
            candidate: event.candidate
          });
        }
      };

      pc.onconnectionstatechange = () => {
        if (pc.connectionState === 'disconnected' || pc.connectionState === 'failed' || pc.connectionState === 'closed') {
          console.log(`[WebRTC] Peer ${targetUserId} connection state: ${pc.connectionState}`);
          removePeer(targetUserId);
        }
      };

      peerConnectionsRef.current[targetUserId] = pc;
      return pc;
    };

    const handleSignalPayload = async (data) => {
      if (!data || data.from_user_id === currentUserID) return;

      console.log("[WebRTC Signaling Received]", data.type, "from:", data.from_user_id);

      if (data.type === 'leave') {
        removePeer(data.from_user_id);
      } else if (data.type === 'join') {
        const pc = createPeerConnection(data.from_user_id);
        const offer = await pc.createOffer();
        await pc.setLocalDescription(offer);
        sendSignal({
          type: 'offer',
          meeting_id: id,
          from_user_id: currentUserID,
          to_user_id: data.from_user_id,
          offer: offer
        });
      } else if (data.type === 'offer' && (data.to_user_id === currentUserID || !data.to_user_id)) {
        const pc = createPeerConnection(data.from_user_id);
        await pc.setRemoteDescription(new RTCSessionDescription(data.offer));
        const answer = await pc.createAnswer();
        await pc.setLocalDescription(answer);
        sendSignal({
          type: 'answer',
          meeting_id: id,
          from_user_id: currentUserID,
          to_user_id: data.from_user_id,
          answer: answer
        });
      } else if (data.type === 'answer' && data.to_user_id === currentUserID) {
        const pc = peerConnectionsRef.current[data.from_user_id];
        if (pc) {
          await pc.setRemoteDescription(new RTCSessionDescription(data.answer));
        }
      } else if (data.type === 'candidate' && data.to_user_id === currentUserID) {
        const pc = peerConnectionsRef.current[data.from_user_id];
        if (pc) {
          try {
            await pc.addIceCandidate(new RTCIceCandidate(data.candidate));
          } catch (err) {
            console.warn("[WebRTC] Error adding ICE candidate:", err);
          }
        }
      } else if (data.type === 'screen_share_status') {
        if (data.sharing) {
          setRemoteScreenPresenter({ userId: data.from_user_id, name: data.name });
        } else {
          setRemoteScreenPresenter(null);
        }
      }
    };

    const initLocalStreamAndSignaling = async () => {
      try {
        if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
          const stream = await navigator.mediaDevices.getUserMedia({
            video: { width: { ideal: 1280 }, height: { ideal: 720 }, frameRate: { ideal: 30 } },
            audio: true
          });

          if (!isMounted) {
            stream.getTracks().forEach(t => t.stop());
            return;
          }

          localStreamRef.current = stream;
          if (localVideoRef.current) {
            localVideoRef.current.srcObject = stream;
          }
        }
      } catch (err) {
        console.warn("Camera/Microphone access not available or denied:", err);
        setMediaError("Không thể truy cập Camera/Microphone. Bạn vẫn có thể tham gia cuộc họp ở chế độ theo dõi.");
      }

      // Initialize SSE Notifications & Signaling Connection
      const token = localStorage.getItem('auth_token');
      if (token) {
        try {
          const sseUrl = `/api/v1/meetings/notifications/stream?token=${encodeURIComponent(token)}`;
          const es = new EventSource(sseUrl);
          eventSourceRef.current = es;

          es.onmessage = (event) => {
            try {
              const notif = JSON.parse(event.data);
              if (notif && notif.type === 'webrtc_signal' && notif.payload) {
                handleSignalPayload(notif.payload);
              }
            } catch (e) {
              console.warn("Error parsing SSE WebRTC signal:", e);
            }
          };

          es.onerror = (e) => {
            console.warn("SSE WebRTC signaling stream disconnected, reconnecting...", e);
          };
        } catch (e) {
          console.warn("Failed to create SSE signaling connection:", e);
        }
      }

      // BroadcastChannel Fallback for tab-to-tab on same browser
      try {
        const bc = new BroadcastChannel(`meeting_channel_${id}`);
        broadcastChannelRef.current = bc;
        bc.onmessage = (event) => {
          if (event.data) {
            handleSignalPayload(event.data);
          }
        };
      } catch (e) {
        console.warn("BroadcastChannel not supported:", e);
      }

      // Send JOIN Signal
      sendSignal({
        type: 'join',
        meeting_id: id,
        from_user_id: currentUserID,
        name: selfName
      });
    };

    if (id && currentUserID) {
      initLocalStreamAndSignaling();
    }

    const handleBeforeUnload = () => {
      try {
        sendSignal({
          type: 'leave',
          meeting_id: id,
          from_user_id: currentUserID
        });
      } catch {}
    };
    window.addEventListener('beforeunload', handleBeforeUnload);

    return () => {
      isMounted = false;
      window.removeEventListener('beforeunload', handleBeforeUnload);

      try {
        sendSignal({
          type: 'leave',
          meeting_id: id,
          from_user_id: currentUserID
        });
      } catch {}

      if (eventSourceRef.current) {
        try { eventSourceRef.current.close(); } catch {}
      }
      if (broadcastChannelRef.current) {
        try { broadcastChannelRef.current.close(); } catch {}
      }

      Object.values(peerConnectionsRef.current).forEach(pc => {
        try { pc.close(); } catch {}
      });
      peerConnectionsRef.current = {};

      if (localStreamRef.current) {
        localStreamRef.current.getTracks().forEach(track => track.stop());
      }
      if (screenStreamRef.current) {
        screenStreamRef.current.getTracks().forEach(track => track.stop());
      }
    };
  }, [id, currentUserID]);

  // Toggle Microphone
  const toggleMic = () => {
    if (localStreamRef.current) {
      const audioTrack = localStreamRef.current.getAudioTracks()[0];
      if (audioTrack) {
        audioTrack.enabled = !micOn;
        setMicOn(!micOn);
      }
    }
  };

  // Toggle Camera
  const toggleCam = () => {
    if (localStreamRef.current) {
      const videoTrack = localStreamRef.current.getVideoTracks()[0];
      if (videoTrack) {
        videoTrack.enabled = !camOn;
        setCamOn(!camOn);
      }
    }
  };

  // Toggle Screen Sharing
  const toggleScreenShare = async () => {
    if (sharingScreen) {
      if (screenStreamRef.current) {
        screenStreamRef.current.getTracks().forEach(t => t.stop());
        screenStreamRef.current = null;
      }
      setSharingScreen(false);

      const webcamTrack = localStreamRef.current?.getVideoTracks()[0];
      if (webcamTrack) {
        Object.values(peerConnectionsRef.current).forEach(pc => {
          const sender = pc.getSenders().find(s => s.track && s.track.kind === 'video');
          if (sender) {
            sender.replaceTrack(webcamTrack);
          }
        });
      }

      sendSignal({
        type: 'screen_share_status',
        meeting_id: id,
        from_user_id: currentUserID,
        sharing: false,
        name: selfName
      });
    } else {
      try {
        if (navigator.mediaDevices && navigator.mediaDevices.getDisplayMedia) {
          const stream = await navigator.mediaDevices.getDisplayMedia({ video: true });
          screenStreamRef.current = stream;
          setSharingScreen(true);

          if (screenVideoRef.current) {
            screenVideoRef.current.srcObject = stream;
          }

          const screenTrack = stream.getVideoTracks()[0];

          Object.values(peerConnectionsRef.current).forEach(pc => {
            const sender = pc.getSenders().find(s => s.track && s.track.kind === 'video');
            if (sender) {
              sender.replaceTrack(screenTrack);
            }
          });

          sendSignal({
            type: 'screen_share_status',
            meeting_id: id,
            from_user_id: currentUserID,
            sharing: true,
            name: selfName
          });

          screenTrack.onended = () => {
            setSharingScreen(false);
            screenStreamRef.current = null;

            const webcamTrack = localStreamRef.current?.getVideoTracks()[0];
            if (webcamTrack) {
              Object.values(peerConnectionsRef.current).forEach(pc => {
                const sender = pc.getSenders().find(s => s.track && s.track.kind === 'video');
                if (sender) {
                  sender.replaceTrack(webcamTrack);
                }
              });
            }

            sendSignal({
              type: 'screen_share_status',
              meeting_id: id,
              from_user_id: currentUserID,
              sharing: false,
              name: selfName
            });
          };
        }
      } catch (err) {
        console.warn("Screen share cancelled or failed:", err);
        setSharingScreen(false);
      }
    }
  };

  // Leave meeting confirm
  const handleConfirmLeaveMeeting = async () => {
    try {
      await sendSignal({
        type: 'leave',
        meeting_id: id,
        from_user_id: currentUserID
      });
    } catch {}
    navigate('/meetings');
  };

  // Copy meeting link
  const handleCopyLink = () => {
    const url = window.location.href;
    navigator.clipboard.writeText(url).then(() => {
      setCopiedLink(true);
      toast.success("Thành công", "Đã sao chép liên kết phòng họp nội bộ!");
      setTimeout(() => setCopiedLink(false), 3000);
    });
  };

  const isGoogleMeetUrl = meeting?.room_url && (meeting.room_url.includes('meet.google.com') || meeting.room_url.startsWith('http'));

  const activeRemotePresenterName = remoteScreenPresenter?.userId
    ? (participants.find(p => p.id === remoteScreenPresenter.userId)?.name || remoteScreenPresenter.name || 'Thành viên')
    : null;

  if (loading) {
    return (
      <div style={{ background: '#0f172a', minHeight: 'calc(100vh - 80px)', color: '#fff', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: '16px', borderRadius: '12px' }}>
        <div style={{ width: '48px', height: '48px', border: '4px solid #38bdf8', borderTopColor: 'transparent', borderRadius: '50%', animation: 'spin 1s linear infinite' }}></div>
        <h3 style={{ margin: 0, fontSize: '16px', color: '#94a3b8' }}>Đang kết nối phòng họp video nội bộ WebRTC...</h3>
      </div>
    );
  }

  if (!meeting) {
    return (
      <div style={{ background: '#0f172a', minHeight: 'calc(100vh - 80px)', color: '#fff', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: '16px', borderRadius: '12px' }}>
        <i className="fa-solid fa-calendar-xmark" style={{ fontSize: '48px', color: '#ef4444' }}></i>
        <h3 style={{ margin: 0, fontSize: '18px' }}>Không tìm thấy phòng họp hoặc cuộc họp đã kết thúc.</h3>
        <button
          onClick={() => navigate('/meetings')}
          style={{ background: '#0284c7', color: '#fff', border: 'none', padding: '10px 20px', borderRadius: '8px', fontWeight: '600', cursor: 'pointer' }}
        >
          Quay lại danh sách cuộc họp
        </button>
      </div>
    );
  }

  return (
    <div style={{ background: '#0f172a', minHeight: 'calc(100vh - 80px)', color: '#fff', padding: '20px', borderRadius: '16px', display: 'flex', flexDirection: 'column', justifyContent: 'space-between', boxShadow: '0 20px 40px rgba(0,0,0,0.5)', position: 'relative' }}>
      {/* Top Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid #1e293b', paddingBottom: '16px', flexWrap: 'wrap', gap: '12px' }}>
        <div>
          <h2 style={{ fontSize: '18px', fontWeight: '700', margin: 0, color: '#f8fafc', display: 'flex', alignItems: 'center', gap: '10px' }}>
            <i className="fa-solid fa-video" style={{ color: '#38bdf8' }}></i>
            Phòng họp Nội bộ: {meeting.title || 'Họp Lãnh Đạo & Điều Hành'}
          </h2>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginTop: '6px', fontSize: '12px' }}>
            <span style={{ color: '#10b981', background: 'rgba(16,185,129,0.15)', padding: '3px 10px', borderRadius: '12px', display: 'inline-flex', alignItems: 'center', gap: '6px', fontWeight: '600' }}>
              <i className="fa-solid fa-circle" style={{ fontSize: '8px', color: '#10b981' }}></i> WebRTC Live Connected
            </span>
            <span style={{ color: '#94a3b8' }}>
              <i className="fa-solid fa-users" style={{ marginRight: '6px' }}></i> {participants.length} Thành viên
            </span>
          </div>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
          {isGoogleMeetUrl && (
            <a
              href={meeting.room_url}
              target="_blank"
              rel="noopener noreferrer"
              style={{ background: '#1e293b', color: '#38bdf8', textDecoration: 'none', padding: '8px 14px', borderRadius: '8px', fontSize: '13px', fontWeight: '600', display: 'flex', alignItems: 'center', gap: '6px', border: '1px solid #334155' }}
            >
              <i className="fa-brands fa-google"></i> Mở Google Meet
            </a>
          )}
          <button
            onClick={handleCopyLink}
            style={{ background: '#1e293b', color: '#f8fafc', border: '1px solid #334155', padding: '8px 14px', borderRadius: '8px', fontSize: '13px', fontWeight: '600', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '6px' }}
          >
            <i className="fa-solid fa-link"></i> {copiedLink ? 'Đã sao chép!' : 'Sao chép Link'}
          </button>
          <button
            onClick={handleConfirmLeaveMeeting}
            style={{ background: '#ef4444', color: '#fff', border: 'none', padding: '9px 18px', borderRadius: '8px', fontWeight: '600', fontSize: '13px', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '6px', boxShadow: '0 4px 12px rgba(239, 68, 68, 0.3)' }}
          >
            <i className="fa-solid fa-phone-slash"></i> Rời phòng họp
          </button>
        </div>
      </div>

      {mediaError && (
        <div style={{ background: 'rgba(239, 68, 68, 0.15)', border: '1px solid #ef4444', color: '#fca5a5', padding: '10px 16px', borderRadius: '8px', margin: '12px 0 0 0', fontSize: '13px', display: 'flex', alignItems: 'center', gap: '8px' }}>
          <i className="fa-solid fa-triangle-exclamation" style={{ color: '#ef4444' }}></i>
          {mediaError}
        </div>
      )}

      {/* Main Video Stream Container */}
      <div style={{ margin: '16px 0', flex: 1, display: 'flex', flexDirection: 'column', gap: '16px' }}>
        {/* Local Screen Share Presenter Container */}
        {sharingScreen && (
          <div style={{ background: '#020617', borderRadius: '14px', border: '2px solid #38bdf8', height: '420px', overflow: 'hidden', position: 'relative', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <video
              ref={screenVideoRef}
              autoPlay
              playsInline
              style={{ width: '100%', height: '100%', objectFit: 'contain' }}
            />
            <div style={{ position: 'absolute', top: '14px', left: '14px', background: 'rgba(2, 132, 199, 0.85)', backdropFilter: 'blur(8px)', padding: '6px 14px', borderRadius: '20px', fontSize: '12px', fontWeight: '600', display: 'flex', alignItems: 'center', gap: '8px' }}>
              <i className="fa-solid fa-desktop"></i> Bạn đang trình chiếu màn hình trực tiếp
            </div>
          </div>
        )}

        {/* Remote Screen Share Presenter Container */}
        {!sharingScreen && remoteScreenPresenter && (
          <div style={{ background: '#020617', borderRadius: '14px', border: '2px solid #a855f7', height: '420px', overflow: 'hidden', position: 'relative', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            {remoteStreams[remoteScreenPresenter.userId] ? (
              <RemoteVideo
                stream={remoteStreams[remoteScreenPresenter.userId]}
                name={activeRemotePresenterName}
                style={{ objectFit: 'contain', borderRadius: '14px' }}
              />
            ) : (
              <div style={{ textAlign: 'center', color: '#94a3b8' }}>
                <i className="fa-solid fa-desktop" style={{ fontSize: '36px', marginBottom: '10px', color: '#a855f7' }}></i>
                <div>{activeRemotePresenterName} đang chia sẻ màn hình...</div>
              </div>
            )}
            <div style={{ position: 'absolute', top: '14px', left: '14px', background: 'rgba(168, 85, 247, 0.85)', backdropFilter: 'blur(8px)', padding: '6px 14px', borderRadius: '20px', fontSize: '12px', fontWeight: '600', display: 'flex', alignItems: 'center', gap: '8px' }}>
              <i className="fa-solid fa-desktop"></i> {activeRemotePresenterName} đang trình chiếu màn hình
            </div>
          </div>
        )}

        {/* Participant Video Grid */}
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))', gap: '16px', minHeight: '360px' }}>
          {/* Local Participant Box */}
          <div style={{ background: '#1e293b', borderRadius: '14px', border: '1px solid #334155', overflow: 'hidden', position: 'relative', display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '260px' }}>
            {camOn ? (
              <video
                ref={localVideoRef}
                autoPlay
                playsInline
                muted
                style={{ width: '100%', height: '100%', objectFit: 'cover', transform: 'scaleX(-1)' }}
              />
            ) : (
              <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '10px' }}>
                <div style={{ width: '70px', height: '70px', borderRadius: '50%', background: '#0284c7', color: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '24px', fontWeight: '700' }}>
                  {(selfName || 'U').charAt(0).toUpperCase()}
                </div>
                <div style={{ fontSize: '13px', color: '#94a3b8' }}>Đã tắt camera</div>
              </div>
            )}
            <div style={{ position: 'absolute', bottom: '12px', left: '12px', background: 'rgba(15, 23, 42, 0.8)', backdropFilter: 'blur(6px)', padding: '4px 12px', borderRadius: '8px', fontSize: '12px', fontWeight: '600', display: 'flex', alignItems: 'center', gap: '6px' }}>
              <span>{selfName} (Bạn - Host)</span>
              <i className={`fa-solid ${micOn ? 'fa-microphone' : 'fa-microphone-slash'}`} style={{ color: micOn ? '#10b981' : '#ef4444' }}></i>
            </div>
          </div>

          {/* Remote Participants Video Boxes */}
          {participants.filter(p => !p.isSelf).map((p) => {
            const rStream = remoteStreams[p.id];
            return (
              <div key={p.id} style={{ background: '#1e293b', borderRadius: '14px', border: '1px solid #334155', overflow: 'hidden', position: 'relative', display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '260px' }}>
                {rStream ? (
                  <RemoteVideo stream={rStream} name={p.name} />
                ) : (
                  <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '10px' }}>
                    <div style={{ width: '70px', height: '70px', borderRadius: '50%', background: '#64748b', color: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '24px', fontWeight: '700' }}>
                      {(p.name || 'U').charAt(0).toUpperCase()}
                    </div>
                    <div style={{ fontSize: '13px', color: '#94a3b8' }}>Đang đợi tín hiệu video...</div>
                  </div>
                )}
                <div style={{ position: 'absolute', bottom: '12px', left: '12px', background: 'rgba(15, 23, 42, 0.8)', backdropFilter: 'blur(6px)', padding: '4px 12px', borderRadius: '8px', fontSize: '12px', fontWeight: '600', display: 'flex', alignItems: 'center', gap: '6px' }}>
                  <span>{p.name}</span>
                  <i className="fa-solid fa-microphone" style={{ color: '#10b981' }}></i>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Bottom Control Bar */}
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '16px', borderTop: '1px solid #1e293b', paddingTop: '16px' }}>
        <button
          onClick={toggleMic}
          style={{
            background: micOn ? '#334155' : '#ef4444',
            color: '#fff',
            border: 'none',
            padding: '12px 24px',
            borderRadius: '10px',
            fontWeight: '600',
            fontSize: '13px',
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            transition: 'all 0.2s ease'
          }}
        >
          <i className={`fa-solid ${micOn ? 'fa-microphone' : 'fa-microphone-slash'}`}></i>
          {micOn ? 'Tắt Micro' : 'Bật Micro'}
        </button>

        <button
          onClick={toggleCam}
          style={{
            background: camOn ? '#334155' : '#ef4444',
            color: '#fff',
            border: 'none',
            padding: '12px 24px',
            borderRadius: '10px',
            fontWeight: '600',
            fontSize: '13px',
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            transition: 'all 0.2s ease'
          }}
        >
          <i className={`fa-solid ${camOn ? 'fa-video' : 'fa-video-slash'}`}></i>
          {camOn ? 'Tắt Camera' : 'Bật Camera'}
        </button>

        <button
          onClick={toggleScreenShare}
          style={{
            background: sharingScreen ? '#0284c7' : '#334155',
            color: '#fff',
            border: 'none',
            padding: '12px 24px',
            borderRadius: '10px',
            fontWeight: '600',
            fontSize: '13px',
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            transition: 'all 0.2s ease'
          }}
        >
          <i className="fa-solid fa-desktop"></i>
          {sharingScreen ? 'Dừng chia sẻ' : 'Chia sẻ màn hình'}
        </button>

        <button
          onClick={handleConfirmLeaveMeeting}
          style={{
            background: 'linear-gradient(135deg, #ef4444 0%, #dc2626 100%)',
            color: '#fff',
            border: 'none',
            padding: '12px 28px',
            borderRadius: '10px',
            fontWeight: '700',
            fontSize: '13px',
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            boxShadow: '0 4px 14px rgba(239, 68, 68, 0.4)'
          }}
        >
          <i className="fa-solid fa-phone-slash"></i>
          Rời phòng họp
        </button>
      </div>
    </div>
  );
}
