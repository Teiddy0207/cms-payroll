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

  // Auto-recording states
  const [isRecording, setIsRecording] = useState(false);
  const [recordingTime, setRecordingTime] = useState(0);

  // AI Summary states (shown after host ends meeting)
  const [aiSummaryModal, setAiSummaryModal] = useState(false);
  const [aiSummary, setAiSummary] = useState(null);
  const [aiProcessing, setAiProcessing] = useState(false);

  const localVideoRef = useRef(null);
  const screenVideoRef = useRef(null);
  const localStreamRef = useRef(null);
  const screenStreamRef = useRef(null);

  // Recording refs
  const mediaRecorderRef = useRef(null);
  const audioChunksRef = useRef([]);
  const recTimerRef = useRef(null);
  
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
          employeesAPI.list({ page_size: 500 }).catch(() => null)
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

          // Auto-start recording from the audio track of the local stream
          try {
            const audioStream = new MediaStream(stream.getAudioTracks());
            const mimeType = MediaRecorder.isTypeSupported('audio/webm;codecs=opus')
              ? 'audio/webm;codecs=opus'
              : 'audio/webm';
            const recorder = new MediaRecorder(audioStream, { mimeType });
            audioChunksRef.current = [];

            recorder.ondataavailable = (e) => {
              if (e.data && e.data.size > 0) audioChunksRef.current.push(e.data);
            };

            recorder.start(1000);
            mediaRecorderRef.current = recorder;
            setIsRecording(true);

            // Tick recording timer every second
            recTimerRef.current = setInterval(() => {
              setRecordingTime(prev => prev + 1);
            }, 1000);
          } catch (recErr) {
            console.warn('Auto-recording failed to start:', recErr);
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

      // Stop recording on unmount
      if (recTimerRef.current) clearInterval(recTimerRef.current);
      if (mediaRecorderRef.current && mediaRecorderRef.current.state !== 'inactive') {
        try { mediaRecorderRef.current.stop(); } catch {}
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

  // Format recording time mm:ss
  const formatRecTime = (secs) => {
    const m = Math.floor(secs / 60).toString().padStart(2, '0');
    const s = (secs % 60).toString().padStart(2, '0');
    return `${m}:${s}`;
  };

  // Leave meeting (non-host: just leave, no AI)
  const handleConfirmLeaveMeeting = async () => {
    try {
      await sendSignal({ type: 'leave', meeting_id: id, from_user_id: currentUserID });
    } catch {}
    navigate('/meetings');
  };

  // End meeting (host only): stop recording → send to AI → show summary
  const handleEndMeeting = async () => {
    setAiProcessing(true);

    // Stop recording and collect chunks
    if (recTimerRef.current) { clearInterval(recTimerRef.current); recTimerRef.current = null; }
    setIsRecording(false);

    const finishAndAnalyze = async (chunks) => {
      try {
        await sendSignal({ type: 'leave', meeting_id: id, from_user_id: currentUserID });
      } catch {}

      if (chunks.length === 0) {
        toast.error('Không có dữ liệu ghi âm', 'Cuộc họp quá ngắn hoặc microphone chưa hoạt động.');
        setAiProcessing(false);
        navigate('/meetings');
        return;
      }

      try {
        const mimeType = chunks[0]?.type || 'audio/webm';
        const blob = new Blob(chunks, { type: mimeType });
        const ext = mimeType.includes('webm') ? 'webm' : 'ogg';
        const audioFile = new File([blob], `meeting_${id}_${Date.now()}.${ext}`, { type: mimeType });

        const formData = new FormData();
        formData.append('audio', audioFile);
        const token = localStorage.getItem('auth_token');

        const res = await fetch(`/api/v1/meetings/${id}/complete`, {
          method: 'POST',
          headers: { Authorization: `Bearer ${token}` },
          body: formData,
        });
        const json = await res.json();
        if (!res.ok || json.status !== 'success') throw new Error(json.error || 'Lỗi không xác định');

        setAiSummary(json.data);
        setAiSummaryModal(true);
      } catch (err) {
        console.error('AI analysis failed:', err);
        toast.error('Lỗi AI', err.message || 'Không thể phân tích cuộc họp.');
        navigate('/meetings');
      } finally {
        setAiProcessing(false);
      }
    };

    const recorder = mediaRecorderRef.current;
    if (recorder && recorder.state !== 'inactive') {
      // Wait for onstop to fire to collect final chunk
      recorder.onstop = () => {
        finishAndAnalyze([...audioChunksRef.current]);
      };
      recorder.stop();
    } else {
      finishAndAnalyze([...audioChunksRef.current]);
    }
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

  const isHost = meeting && (meeting.host_id === currentUserID);

  return (
    <>
    {/* ── AI Processing Overlay ────────────────────────────── */}
    {aiProcessing && (
      <div style={{
        position: 'fixed', inset: 0, zIndex: 9999,
        background: 'rgba(15,23,42,0.92)',
        backdropFilter: 'blur(8px)',
        display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: '20px'
      }}>
        <div style={{ width: '64px', height: '64px', border: '4px solid #7c3aed', borderTopColor: 'transparent', borderRadius: '50%', animation: 'spin 1s linear infinite' }} />
        <div style={{ textAlign: 'center' }}>
          <div style={{ fontSize: '20px', fontWeight: '700', color: '#f8fafc', marginBottom: '8px' }}>
            <i className="fa-solid fa-robot" style={{ color: '#a78bfa', marginRight: '10px' }}></i>
            AI đang phân tích cuộc họp...
          </div>
          <div style={{ fontSize: '14px', color: '#94a3b8' }}>Whisper đang chuyển đổi âm thanh · RL Agent đang tóm tắt</div>
        </div>
      </div>
    )}

    {/* ── AI Summary Modal ─────────────────────────────────── */}
    {aiSummaryModal && aiSummary && (
      <div style={{
        position: 'fixed', inset: 0, zIndex: 9998,
        background: 'rgba(15,23,42,0.85)', backdropFilter: 'blur(8px)',
        display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '20px'
      }}>
        <div style={{
          background: '#0f172a', border: '1px solid #334155',
          borderRadius: '20px', padding: '32px', width: '100%', maxWidth: '720px',
          maxHeight: '85vh', overflowY: 'auto',
          boxShadow: '0 40px 80px rgba(0,0,0,0.6)'
        }}>
          {/* Modal Header */}
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
              <div style={{
                width: '44px', height: '44px', borderRadius: '12px',
                background: 'linear-gradient(135deg, #7c3aed, #6d28d9)',
                display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '20px'
              }}>
                <i className="fa-solid fa-robot" style={{ color: '#fff' }}></i>
              </div>
              <div>
                <div style={{ fontSize: '18px', fontWeight: '700', color: '#f8fafc' }}>Tóm tắt cuộc họp bằng AI</div>
                <div style={{ fontSize: '12px', color: '#94a3b8' }}>{meeting.title}</div>
              </div>
            </div>
            <button
              onClick={() => { setAiSummaryModal(false); navigate('/meetings'); }}
              style={{ background: '#1e293b', border: '1px solid #334155', color: '#94a3b8', width: '36px', height: '36px', borderRadius: '8px', cursor: 'pointer', fontSize: '16px' }}
            >
              <i className="fa-solid fa-xmark"></i>
            </button>
          </div>

          {/* Scores */}
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px', marginBottom: '24px' }}>
            {aiSummary.efficiency_score != null && (
              <div style={{ background: '#1e293b', borderRadius: '12px', padding: '16px', textAlign: 'center', border: '1px solid #334155' }}>
                <div style={{ fontSize: '11px', color: '#64748b', textTransform: 'uppercase', letterSpacing: '1px', marginBottom: '6px' }}>Hiệu suất họp</div>
                <div style={{ fontSize: '28px', fontWeight: '800', color: '#10b981' }}>{aiSummary.efficiency_score}<span style={{ fontSize: '14px', color: '#64748b' }}>/100</span></div>
              </div>
            )}
            {aiSummary.sentiment && (
              <div style={{ background: '#1e293b', borderRadius: '12px', padding: '16px', textAlign: 'center', border: '1px solid #334155' }}>
                <div style={{ fontSize: '11px', color: '#64748b', textTransform: 'uppercase', letterSpacing: '1px', marginBottom: '6px' }}>Không khí</div>
                <div style={{ fontSize: '22px', fontWeight: '700', color: '#38bdf8', textTransform: 'capitalize' }}>
                  {aiSummary.sentiment === 'positive' ? '😊 Tích cực' : aiSummary.sentiment === 'negative' ? '😟 Tiêu cực' : '😐 Trung tính'}
                </div>
              </div>
            )}
          </div>

          {/* Summary */}
          {aiSummary.summary && (
            <div style={{ marginBottom: '20px' }}>
              <div style={{ fontSize: '13px', fontWeight: '700', color: '#7c3aed', marginBottom: '10px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                <i className="fa-solid fa-align-left"></i> NỘI DUNG TÓM TẮT
              </div>
              <div style={{ background: '#1e293b', borderRadius: '10px', padding: '16px', fontSize: '14px', color: '#cbd5e1', lineHeight: '1.8', border: '1px solid #334155' }}>
                {aiSummary.summary}
              </div>
            </div>
          )}

          {/* Key Decisions */}
          {Array.isArray(aiSummary.key_decisions) && aiSummary.key_decisions.length > 0 && (
            <div style={{ marginBottom: '20px' }}>
              <div style={{ fontSize: '13px', fontWeight: '700', color: '#f59e0b', marginBottom: '10px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                <i className="fa-solid fa-gavel"></i> QUYẾT ĐỊNH CHÍNH
              </div>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                {aiSummary.key_decisions.map((d, i) => (
                  <div key={i} style={{ background: '#1e293b', borderLeft: '3px solid #f59e0b', borderRadius: '0 8px 8px 0', padding: '10px 14px', fontSize: '13px', color: '#e2e8f0' }}>
                    <i className="fa-solid fa-check" style={{ color: '#f59e0b', marginRight: '8px' }}></i>{d}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Action Items */}
          {Array.isArray(aiSummary.action_items) && aiSummary.action_items.length > 0 && (
            <div style={{ marginBottom: '24px' }}>
              <div style={{ fontSize: '13px', fontWeight: '700', color: '#10b981', marginBottom: '10px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                <i className="fa-solid fa-list-check"></i> ACTION ITEMS
              </div>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                {aiSummary.action_items.map((item, i) => (
                  <div key={i} style={{ background: '#1e293b', borderLeft: '3px solid #10b981', borderRadius: '0 8px 8px 0', padding: '10px 14px', fontSize: '13px', color: '#e2e8f0' }}>
                    <i className="fa-solid fa-arrow-right" style={{ color: '#10b981', marginRight: '8px' }}></i>{item}
                  </div>
                ))}
              </div>
            </div>
          )}

          <button
            onClick={() => { setAiSummaryModal(false); navigate('/meetings'); }}
            style={{
              width: '100%', padding: '12px', borderRadius: '10px', border: 'none',
              background: 'linear-gradient(135deg, #7c3aed, #6d28d9)',
              color: '#fff', fontWeight: '700', fontSize: '14px', cursor: 'pointer'
            }}
          >
            <i className="fa-solid fa-check" style={{ marginRight: '8px' }}></i>Hoàn tất & Quay về danh sách
          </button>
        </div>
      </div>
    )}

    <div style={{ background: '#0f172a', minHeight: 'calc(100vh - 80px)', color: '#fff', padding: '20px', borderRadius: '16px', display: 'flex', flexDirection: 'column', justifyContent: 'space-between', boxShadow: '0 20px 40px rgba(0,0,0,0.5)', position: 'relative' }}>
      {/* Top Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid #1e293b', paddingBottom: '16px', flexWrap: 'wrap', gap: '12px' }}>
        <div>
          <h2 style={{ fontSize: '18px', fontWeight: '700', margin: 0, color: '#f8fafc', display: 'flex', alignItems: 'center', gap: '10px' }}>
            <i className="fa-solid fa-video" style={{ color: '#38bdf8' }}></i>
            Phòng họp Nội bộ: {meeting.title || 'Họp Lãnh Đạo & Điều Hành'}
          </h2>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginTop: '6px', fontSize: '12px', flexWrap: 'wrap' }}>
            <span style={{ color: '#10b981', background: 'rgba(16,185,129,0.15)', padding: '3px 10px', borderRadius: '12px', display: 'inline-flex', alignItems: 'center', gap: '6px', fontWeight: '600' }}>
              <i className="fa-solid fa-circle" style={{ fontSize: '8px', color: '#10b981' }}></i> WebRTC Live Connected
            </span>
            <span style={{ color: '#94a3b8' }}>
              <i className="fa-solid fa-users" style={{ marginRight: '6px' }}></i> {participants.length} Thành viên
            </span>
            {isRecording && (
              <span style={{
                color: '#ef4444', background: 'rgba(239,68,68,0.15)',
                padding: '3px 10px', borderRadius: '12px',
                display: 'inline-flex', alignItems: 'center', gap: '6px', fontWeight: '600',
                animation: 'pulse 1.5s ease-in-out infinite'
              }}>
                <i className="fa-solid fa-circle" style={{ fontSize: '8px' }}></i>
                Đang ghi âm {formatRecTime(recordingTime)}
              </span>
            )}
          </div>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '10px', flexWrap: 'wrap' }}>
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
            onClick={handleEndMeeting}
            disabled={aiProcessing}
            style={{
              background: 'linear-gradient(135deg, #8b5cf6, #6d28d9)',
              color: '#fff', border: '1px solid #a78bfa', padding: '9px 18px', borderRadius: '8px',
              fontWeight: '600', fontSize: '13px', cursor: 'pointer',
              display: 'flex', alignItems: 'center', gap: '6px',
              boxShadow: '0 4px 14px rgba(139, 92, 246, 0.5)',
              opacity: aiProcessing ? 0.7 : 1
            }}
          >
            <i className="fa-solid fa-robot"></i> Kết thúc &amp; AI Tóm tắt
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
          onClick={handleEndMeeting}
          disabled={aiProcessing}
          style={{
            background: 'linear-gradient(135deg, #8b5cf6 0%, #6d28d9 100%)',
            color: '#fff',
            border: '1px solid #a78bfa',
            padding: '12px 28px',
            borderRadius: '10px',
            fontWeight: '700',
            fontSize: '14px',
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            boxShadow: '0 4px 18px rgba(139, 92, 246, 0.5)',
            opacity: aiProcessing ? 0.7 : 1
          }}
        >
          <i className="fa-solid fa-robot" style={{ fontSize: '16px' }}></i>
          Kết thúc &amp; AI Tóm tắt
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
    </>
  );
}
