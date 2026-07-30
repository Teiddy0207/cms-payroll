import React, { useState, useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Card, Button, Tag, Typography, Space, Row, Col,
  Spin, Empty, Modal, Form, Input, DatePicker, Select,
  Tooltip, Divider, Progress, Alert
} from 'antd';
import {
  VideoCameraOutlined, PlusOutlined, CalendarOutlined,
  UserOutlined, ClockCircleOutlined, GoogleOutlined,
  CheckOutlined, CloseOutlined, EditOutlined,
  RobotOutlined, AudioOutlined, FileTextOutlined,
  BulbOutlined, ThunderboltOutlined, SmileOutlined,
  FrownOutlined, MehOutlined, UploadOutlined, EyeOutlined
} from '@ant-design/icons';
import { meetingsAPI, employeesAPI } from '../api/client.js';
import { useToast } from '../hooks/useToast.js';
import dayjs from 'dayjs';

const { Title, Text, Paragraph } = Typography;
const { Option } = Select;

function parseJWT(token) {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64).split('').map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)).join('')
    );
    return JSON.parse(jsonPayload);
  } catch { return null; }
}

function SentimentBadge({ value }) {
  const map = {
    POSITIVE: { color: '#10b981', bg: '#ecfdf5', icon: <SmileOutlined />,  label: 'Tích cực' },
    NEGATIVE: { color: '#ef4444', bg: '#fef2f2', icon: <FrownOutlined />,  label: 'Tiêu cực' },
    NEUTRAL:  { color: '#f59e0b', bg: '#fffbeb', icon: <MehOutlined />,    label: 'Trung lập' },
  };
  const s = map[value] || map.NEUTRAL;
  return (
    <span style={{
      display: 'inline-flex', alignItems: 'center', gap: 6,
      padding: '4px 14px', borderRadius: 20,
      background: s.bg, color: s.color, fontWeight: 600, fontSize: 14
    }}>
      {s.icon} {s.label}
    </span>
  );
}

function parseScore(scoreStr) {
  if (!scoreStr) return 0;
  const n = parseFloat(scoreStr);
  return isNaN(n) ? 0 : n;
}

export default function Meetings() {
  const navigate = useNavigate();
  const toast    = useToast();

  const [meetings, setMeetings]           = useState([]);
  const [employees, setEmployees]         = useState([]);
  const [loading, setLoading]             = useState(true);
  const [isModalOpen, setIsModalOpen]     = useState(false);
  const [editingMeeting, setEditingMeeting] = useState(null);
  const [form]                            = Form.useForm();
  const [submitting, setSubmitting]       = useState(false);

  // AI Summary states
  const [completeModal, setCompleteModal] = useState(false);
  const [summaryModal, setSummaryModal]   = useState(false);
  const [targetMeeting, setTargetMeeting] = useState(null);
  const [audioFile, setAudioFile]         = useState(null);
  const [completing, setCompleting]       = useState(false);
  const [currentSummary, setCurrentSummary] = useState(null);

  const currentUserID = useMemo(() => {
    const token = localStorage.getItem('auth_token');
    if (!token) return null;
    const claims = parseJWT(token);
    return claims?.user_id || null;
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [meetingsRes, employeesRes] = await Promise.all([
        meetingsAPI.list(),
        employeesAPI.list({ page_size: 500 })
      ]);
      let mList = [];
      if (meetingsRes?.data) {
        if (Array.isArray(meetingsRes.data.data)) mList = meetingsRes.data.data;
        else if (Array.isArray(meetingsRes.data))  mList = meetingsRes.data;
      }
      setMeetings(mList);

      let eList = [];
      if (employeesRes?.data?.data) {
        if (Array.isArray(employeesRes.data.data.items)) eList = employeesRes.data.data.items;
        else if (Array.isArray(employeesRes.data.data))  eList = employeesRes.data.data;
      }
      setEmployees(eList);
    } catch (err) {
      console.error('fetchData error:', err);
      toast.error('Lỗi', 'Không thể tải danh sách cuộc họp hoặc nhân sự');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    const handler = () => fetchData();
    window.addEventListener('meeting_updated', handler);
    return () => window.removeEventListener('meeting_updated', handler);
  }, []);

  const handleRSVP = async (meetingId, status) => {
    try {
      await meetingsAPI.rsvp(meetingId, status);
      toast.success('Thành công', `Đã cập nhật: ${status === 'ACCEPTED' ? 'Đồng ý' : 'Từ chối'}`);
      fetchData();
    } catch {
      toast.error('Lỗi', 'Không thể cập nhật trạng thái');
    }
  };

  const handleCreateOrUpdateMeeting = async (values) => {
    setSubmitting(true);
    try {
      const payload = {
        title: values.title,
        description: values.description || '',
        start_time: values.timeRange[0].toISOString(),
        end_time:   values.timeRange[1].toISOString(),
        attendee_ids: Array.isArray(values.attendees) ? values.attendees : []
      };
      if (editingMeeting) {
        await meetingsAPI.update(editingMeeting.id, payload);
        toast.success('Thành công', 'Đã cập nhật thông tin cuộc họp!');
      } else {
        await meetingsAPI.create(payload);
        toast.success('Thành công', 'Đã tạo cuộc họp và gửi thông báo!');
      }
      setIsModalOpen(false); setEditingMeeting(null); form.resetFields(); fetchData();
    } catch {
      toast.error('Lỗi', 'Không thể lưu thông tin cuộc họp');
    } finally {
      setSubmitting(false);
    }
  };

  const handleOpenCreateModal = () => { setEditingMeeting(null); form.resetFields(); setIsModalOpen(true); };
  const handleOpenEditModal = (m) => {
    setEditingMeeting(m);
    form.setFieldsValue({
      title: m.title, description: m.description,
      timeRange: [dayjs(m.start_time), dayjs(m.end_time)],
      attendees: (m.attendees || []).map(a => a.user_id || a.id)
    });
    setIsModalOpen(true);
  };

  // ── Complete meeting + AI ─────────────────────────────────
  const openCompleteModal = (m) => { setTargetMeeting(m); setAudioFile(null); setCompleteModal(true); };

  const handleCompleteMeeting = async () => {
    if (!audioFile) { toast.error('Thiếu file', 'Vui lòng chọn file audio'); return; }
    setCompleting(true);
    try {
      const formData = new FormData();
      formData.append('audio', audioFile);
      const token = localStorage.getItem('auth_token');
      const res = await fetch(`/api/v1/meetings/${targetMeeting.id}/complete`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
      });
      const json = await res.json();
      if (!res.ok || json.status !== 'success') throw new Error(json.error || 'Lỗi không xác định');
      toast.success('Hoàn thành!', 'AI đã tóm tắt cuộc họp thành công');
      setCompleteModal(false);
      setCurrentSummary(json.data);
      setSummaryModal(true);
      fetchData();
    } catch (err) {
      toast.error('Lỗi', err.message || 'Không thể kết thúc cuộc họp');
    } finally {
      setCompleting(false);
    }
  };

  const handleViewSummary = async (m) => {
    try {
      const token = localStorage.getItem('auth_token');
      const res = await fetch(`/api/v1/meetings/${m.id}/summary`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      const json = await res.json();
      if (json.status === 'success' && json.data) {
        setCurrentSummary(json.data); setSummaryModal(true);
      } else {
        toast.error('Chưa có tóm tắt', 'Cuộc họp này chưa có bản tóm tắt AI');
      }
    } catch {
      toast.error('Lỗi', 'Không thể tải tóm tắt');
    }
  };

  const resolveEmployeeName = (id) => {
    if (!id || !Array.isArray(employees)) return 'Thành viên';
    const emp = employees.find(e => e && (e.user_id === id || e.id === id));
    return emp ? emp.full_name : 'Thành viên';
  };

  // ────────────────────────────────────────────────────────────
  return (
    <div style={{ padding: '24px', background: '#f8fafc', minHeight: 'calc(100vh - 64px)' }}>

      {/* Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24, flexWrap: 'wrap', gap: 16 }}>
        <div>
          <Title level={3} style={{ margin: 0, color: '#0f172a', fontWeight: 700 }}>
            <CalendarOutlined style={{ marginRight: 10, color: '#0284c7' }} />
            Phòng họp Nội bộ &amp; Lịch làm việc
          </Title>
          <Text type="secondary">Quản lý lịch họp, thông báo thành viên và Video WebRTC P2P</Text>
        </div>
        <Button type="primary" icon={<PlusOutlined />} size="large" onClick={handleOpenCreateModal}
          style={{ background: 'linear-gradient(135deg, #0284c7 0%, #0369a1 100%)', border: 'none', borderRadius: 8, fontWeight: 600 }}>
          Tạo cuộc họp mới
        </Button>
      </div>

      {/* Cards */}
      {loading ? (
        <div style={{ textAlign: 'center', padding: '60px 0' }}>
          <Spin size="large" />
          <div style={{ marginTop: 16, color: '#64748b' }}>Đang tải...</div>
        </div>
      ) : meetings.length === 0 ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="Chưa có cuộc họp nào"
          style={{ background: '#fff', padding: 48, borderRadius: 12 }}>
          <Button type="primary" onClick={handleOpenCreateModal}>Tạo cuộc họp đầu tiên</Button>
        </Empty>
      ) : (
        <Row gutter={[20, 20]}>
          {meetings.map((m) => {
            if (!m) return null;
            const isScheduled = m.status === 'SCHEDULED';
            const isCompleted = m.status === 'COMPLETED';
            const isCancelled = m.status === 'CANCELLED';
            const isHost      = m.host_id === currentUserID;
            const myRSVP      = (m.attendees || []).find(a => a.user_id === currentUserID)?.rsvp_status || 'PENDING';
            const statusColor = isCompleted ? 'blue' : isScheduled ? 'green' : isCancelled ? 'red' : 'default';

            return (
              <Col xs={24} sm={12} lg={8} key={m.id}>
                <Card hoverable style={{ height: '100%', borderRadius: 12, border: '1px solid #e2e8f0', boxShadow: '0 2px 8px rgba(0,0,0,0.04)', overflow: 'hidden' }}
                  bodyStyle={{ padding: 20 }}>

                  {/* Status row */}
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 14 }}>
                    <Tag color={statusColor} style={{ borderRadius: 4, fontWeight: 600 }}>
                      {isCompleted ? '✅ COMPLETED' : m.status || 'SCHEDULED'}
                    </Tag>
                    {m.google_event_id && (
                      <Tooltip title="Đã đồng bộ Google Calendar">
                        <Tag color="processing" icon={<GoogleOutlined />} style={{ margin: 0 }}>Google Synced</Tag>
                      </Tooltip>
                    )}
                  </div>

                  {/* Title */}
                  <Title level={5} style={{ margin: '0 0 8px', fontWeight: 700, color: '#1e293b' }}>
                    {m.title || 'Không có tiêu đề'}
                  </Title>
                  <Paragraph ellipsis={{ rows: 2 }} style={{ color: '#64748b', fontSize: 13, marginBottom: 14 }}>
                    {m.description || 'Không có mô tả.'}
                  </Paragraph>

                  {/* Info */}
                  <Space direction="vertical" size={6} style={{ width: '100%', fontSize: 13, color: '#475569', marginBottom: 14 }}>
                    <div>
                      <ClockCircleOutlined style={{ color: '#0284c7', marginRight: 8 }} />
                      {dayjs(m.start_time).format('HH:mm DD/MM/YYYY')}
                      <span style={{ color: '#94a3b8' }}> → </span>
                      {dayjs(m.end_time).format('HH:mm DD/MM/YYYY')}
                    </div>
                    <div>
                      <UserOutlined style={{ color: '#0284c7', marginRight: 8 }} />
                      Chủ trì: <strong>{resolveEmployeeName(m.host_id)}</strong>
                    </div>
                    <div style={{ display: 'flex', flexWrap: 'wrap', gap: 4, marginTop: 2 }}>
                      {(m.attendees || []).map((att, i) => {
                        const rsvp  = att.rsvp_status || 'PENDING';
                        const color = rsvp === 'ACCEPTED' ? 'green' : rsvp === 'DECLINED' ? 'red' : 'orange';
                        return (
                          <Tag key={i} color={color} style={{ fontSize: 11, borderRadius: 4 }}>
                            {resolveEmployeeName(att.user_id || att.id)} ({rsvp})
                          </Tag>
                        );
                      })}
                    </div>
                  </Space>

                  {/* AI Summary badge */}
                  {isCompleted && (
                    <div style={{ background: 'linear-gradient(135deg, #ede9fe, #dbeafe)', borderRadius: 8, padding: '8px 12px', marginBottom: 14, display: 'flex', alignItems: 'center', gap: 8 }}>
                      <RobotOutlined style={{ color: '#7c3aed' }} />
                      <span style={{ color: '#7c3aed', fontSize: 13, fontWeight: 600 }}>Bản tóm tắt AI đã sẵn sàng</span>
                    </div>
                  )}

                  <Divider style={{ margin: '12px 0' }} />

                  {/* Action buttons */}
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                    <Button icon={<VideoCameraOutlined />} onClick={() => navigate(`/meetings/room/${m.id}`)}
                      style={{ color: '#0284c7', borderColor: '#bae6fd', fontWeight: 600 }}>
                      Vào phòng
                    </Button>

                    {isHost && isScheduled && (
                      <Button icon={<AudioOutlined />} onClick={() => openCompleteModal(m)}
                        style={{ background: 'linear-gradient(135deg, #7c3aed, #6d28d9)', color: '#fff', border: 'none', fontWeight: 600 }}>
                        Kết thúc &amp; Tóm tắt AI
                      </Button>
                    )}

                    {isCompleted && (
                      <Button icon={<EyeOutlined />} onClick={() => handleViewSummary(m)}
                        style={{ color: '#7c3aed', borderColor: '#ddd6fe', fontWeight: 600 }}>
                        Xem tóm tắt
                      </Button>
                    )}

                    {isHost && (
                      <Button icon={<EditOutlined />} onClick={() => handleOpenEditModal(m)} style={{ color: '#8b5cf6' }}>Sửa</Button>
                    )}

                    {!isHost && isScheduled && (
                      <>
                        <Button icon={<CheckOutlined />} onClick={() => handleRSVP(m.id, 'ACCEPTED')}
                          style={{ color: myRSVP === 'ACCEPTED' ? '#10b981' : '#64748b', fontWeight: 600 }}>
                          {myRSVP === 'ACCEPTED' ? 'Đã đồng ý' : 'Đồng ý'}
                        </Button>
                        <Button icon={<CloseOutlined />} onClick={() => handleRSVP(m.id, 'DECLINED')}
                          style={{ color: myRSVP === 'DECLINED' ? '#ef4444' : '#94a3b8' }}>
                          Từ chối
                        </Button>
                      </>
                    )}
                  </div>
                </Card>
              </Col>
            );
          })}
        </Row>
      )}

      {/* ── Modal: Tạo / Sửa cuộc họp ─────────────────────── */}
      <Modal title={editingMeeting ? 'Chỉnh sửa thông tin cuộc họp' : 'Tạo cuộc họp mới'}
        open={isModalOpen} onCancel={() => { setIsModalOpen(false); setEditingMeeting(null); }} footer={null} destroyOnClose>
        <Form form={form} layout="vertical" onFinish={handleCreateOrUpdateMeeting} style={{ marginTop: 16 }}>
          <Form.Item name="title" label="Tiêu đề cuộc họp" rules={[{ required: true, message: 'Vui lòng nhập tiêu đề!' }]}>
            <Input placeholder="Ví dụ: Họp triển khai kế hoạch tháng 8" />
          </Form.Item>
          <Form.Item name="description" label="Nội dung / Mô tả">
            <Input.TextArea rows={3} placeholder="Nhập nội dung thảo luận..." />
          </Form.Item>
          <Form.Item name="timeRange" label="Thời gian bắt đầu & Kết thúc" rules={[{ required: true, message: 'Vui lòng chọn thời gian!' }]}>
            <DatePicker.RangePicker showTime format="YYYY-MM-DD HH:mm" style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="attendees" label="Mời thành viên tham gia">
            <Select mode="multiple" placeholder="Chọn nhân sự tham dự..." optionFilterProp="children" style={{ width: '100%' }}>
              {employees.map(emp => (
                <Option key={emp.user_id || emp.id} value={emp.user_id || emp.id}>
                  {emp.full_name} ({emp.employee_code || 'NV'})
                </Option>
              ))}
            </Select>
          </Form.Item>
          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8, marginTop: 24 }}>
            <Button onClick={() => { setIsModalOpen(false); setEditingMeeting(null); }}>Hủy</Button>
            <Button type="primary" htmlType="submit" loading={submitting}>
              {editingMeeting ? 'Cập nhật cuộc họp' : 'Tạo cuộc họp & Báo NATS'}
            </Button>
          </div>
        </Form>
      </Modal>

      {/* ── Modal: Upload audio → Kết thúc họp ────────────── */}
      <Modal
        title={<div style={{ display: 'flex', alignItems: 'center', gap: 10 }}><RobotOutlined style={{ color: '#7c3aed', fontSize: 20 }} /><span>Kết thúc &amp; Tóm tắt AI</span></div>}
        open={completeModal} onCancel={() => { setCompleteModal(false); setAudioFile(null); }}
        footer={null} destroyOnClose width={520}>
        <div style={{ padding: '16px 0' }}>
          <Alert type="info" showIcon message="AI sẽ tự động phân tích"
            description="Upload file ghi âm cuộc họp. Hệ thống dùng Whisper (STT) và RL Agent để tạo bản tóm tắt, quyết định chính và action items."
            style={{ marginBottom: 24, borderRadius: 8 }} />

          <div style={{ border: '2px dashed #c4b5fd', borderRadius: 12, padding: 32, textAlign: 'center', background: '#faf5ff', marginBottom: 24 }}>
            <AudioOutlined style={{ fontSize: 40, color: '#7c3aed', marginBottom: 12, display: 'block' }} />
            <input type="file" accept=".wav,.mp3,.m4a,.ogg,.webm" style={{ display: 'none' }} id="audio-upload-input"
              onChange={(e) => setAudioFile(e.target.files?.[0] || null)} />
            <label htmlFor="audio-upload-input">
              <Button icon={<UploadOutlined />} style={{ borderColor: '#7c3aed', color: '#7c3aed', fontWeight: 600 }}>
                Chọn file audio
              </Button>
            </label>
            {audioFile ? (
              <div style={{ marginTop: 12 }}>
                <Tag color="purple" style={{ fontSize: 13, padding: '4px 12px' }}>
                  🎵 {audioFile.name} ({(audioFile.size / 1024 / 1024).toFixed(2)} MB)
                </Tag>
              </div>
            ) : (
              <div style={{ marginTop: 10, color: '#94a3b8', fontSize: 13 }}>Hỗ trợ: .wav · .mp3 · .m4a · .ogg · .webm</div>
            )}
          </div>

          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
            <Button onClick={() => { setCompleteModal(false); setAudioFile(null); }}>Hủy</Button>
            <Button type="primary" icon={<RobotOutlined />} loading={completing} disabled={!audioFile}
              onClick={handleCompleteMeeting}
              style={{ background: completing ? undefined : 'linear-gradient(135deg, #7c3aed, #6d28d9)', border: 'none', fontWeight: 600 }}>
              {completing ? 'AI đang phân tích...' : 'Kết thúc & Tóm tắt ngay'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* ── Modal: Xem AI Summary ──────────────────────────── */}
      <Modal
        title={<div style={{ display: 'flex', alignItems: 'center', gap: 10 }}><RobotOutlined style={{ color: '#7c3aed', fontSize: 20 }} /><span style={{ fontWeight: 700, fontSize: 16 }}>Tóm tắt cuộc họp bằng AI</span></div>}
        open={summaryModal} onCancel={() => setSummaryModal(false)}
        footer={<Button onClick={() => setSummaryModal(false)}>Đóng</Button>}
        destroyOnClose width={680}>
        {currentSummary && (
          <div style={{ padding: '8px 0' }}>

            {/* Efficiency + Sentiment */}
            <Row gutter={16} style={{ marginBottom: 24 }}>
              <Col span={12}>
                <div style={{ background: 'linear-gradient(135deg, #ede9fe, #dbeafe)', borderRadius: 12, padding: '16px 20px' }}>
                  <div style={{ color: '#7c3aed', fontSize: 12, fontWeight: 600, marginBottom: 8, textTransform: 'uppercase', letterSpacing: 1 }}>
                    <ThunderboltOutlined /> Hiệu quả cuộc họp
                  </div>
                  <div style={{ fontSize: 28, fontWeight: 800, color: '#4f46e5', marginBottom: 6 }}>
                    {currentSummary.efficiency_score || '—'}
                  </div>
                  <Progress percent={parseScore(currentSummary.efficiency_score)} showInfo={false}
                    strokeColor={{ '0%': '#7c3aed', '100%': '#3b82f6' }} trailColor="#e0e7ff" />
                </div>
              </Col>
              <Col span={12}>
                <div style={{ background: '#f8fafc', borderRadius: 12, padding: '16px 20px', border: '1px solid #e2e8f0', height: '100%', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
                  <div style={{ color: '#64748b', fontSize: 12, fontWeight: 600, marginBottom: 12, textTransform: 'uppercase', letterSpacing: 1 }}>
                    <SmileOutlined /> Không khí cuộc họp
                  </div>
                  <SentimentBadge value={currentSummary.sentiment} />
                </div>
              </Col>
            </Row>

            {/* Summary */}
            <div style={{ marginBottom: 20 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 10 }}>
                <FileTextOutlined style={{ color: '#0284c7', fontSize: 16 }} />
                <span style={{ fontWeight: 700, color: '#1e293b', fontSize: 14 }}>Nội dung tóm tắt</span>
              </div>
              <div style={{ background: '#f0f9ff', borderRadius: 10, padding: '14px 16px', border: '1px solid #bae6fd', color: '#0f172a', lineHeight: 1.8, fontSize: 14 }}>
                {currentSummary.summary || 'Không có dữ liệu.'}
              </div>
            </div>

            {/* Key Decisions */}
            <div style={{ marginBottom: 20 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 10 }}>
                <BulbOutlined style={{ color: '#f59e0b', fontSize: 16 }} />
                <span style={{ fontWeight: 700, color: '#1e293b', fontSize: 14 }}>Quyết định chính</span>
              </div>
              <div style={{ background: '#fffbeb', borderRadius: 10, padding: '14px 16px', border: '1px solid #fde68a', color: '#0f172a', lineHeight: 1.8, fontSize: 14 }}>
                {currentSummary.key_decisions || 'Không có quyết định cụ thể.'}
              </div>
            </div>

            {/* Action Items */}
            <div>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 10 }}>
                <ThunderboltOutlined style={{ color: '#10b981', fontSize: 16 }} />
                <span style={{ fontWeight: 700, color: '#1e293b', fontSize: 14 }}>Việc cần làm (Action Items)</span>
              </div>
              <div style={{ background: '#f0fdf4', borderRadius: 10, padding: '14px 16px', border: '1px solid #bbf7d0', color: '#0f172a', lineHeight: 1.8, fontSize: 14 }}>
                {currentSummary.action_items || 'Không có action items cụ thể.'}
              </div>
            </div>

            <div style={{ marginTop: 20, color: '#94a3b8', fontSize: 12, textAlign: 'right', fontStyle: 'italic' }}>
              ✨ Được phân tích bởi Whisper STT + RL Extractive Summarization (REINFORCE)
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}
