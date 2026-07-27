import React, { useState, useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Card,
  Button,
  Tag,
  Typography,
  Space,
  Row,
  Col,
  Spin,
  Empty,
  Modal,
  Form,
  Input,
  DatePicker,
  Select,
  Tooltip,
  Popconfirm
} from 'antd';
import {
  VideoCameraOutlined,
  PlusOutlined,
  CalendarOutlined,
  UserOutlined,
  ClockCircleOutlined,
  GoogleOutlined,
  CheckOutlined,
  CloseOutlined,
  EditOutlined,
  DeleteOutlined
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
  } catch {
    return null;
  }
}

export default function Meetings() {
  const navigate = useNavigate();
  const toast = useToast();
  const [meetings, setMeetings] = useState([]);
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingMeeting, setEditingMeeting] = useState(null);
  const [form] = Form.useForm();
  const [submitting, setSubmitting] = useState(false);

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
      if (meetingsRes && meetingsRes.data) {
        if (Array.isArray(meetingsRes.data.data)) {
          mList = meetingsRes.data.data;
        } else if (Array.isArray(meetingsRes.data)) {
          mList = meetingsRes.data;
        }
      }
      setMeetings(mList);

      let eList = [];
      if (employeesRes && employeesRes.data && employeesRes.data.data) {
        if (Array.isArray(employeesRes.data.data.items)) {
          eList = employeesRes.data.data.items;
        } else if (Array.isArray(employeesRes.data.data)) {
          eList = employeesRes.data.data;
        }
      }
      setEmployees(eList);
    } catch (err) {
      console.error("Error fetching meetings or employees:", err);
      toast.error("Lỗi", "Không thể tải danh sách cuộc họp hoặc nhân sự");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();

    const handleMeetingUpdated = () => {
      fetchData();
    };
    window.addEventListener('meeting_updated', handleMeetingUpdated);
    return () => {
      window.removeEventListener('meeting_updated', handleMeetingUpdated);
    };
  }, []);

  const handleRSVP = async (meetingId, status) => {
    try {
      await meetingsAPI.rsvp(meetingId, status);
      toast.success("Thành công", `Đã cập nhật phản hồi: ${status === 'ACCEPTED' ? 'Đồng ý' : 'Từ chối'}`);
      fetchData();
    } catch (err) {
      toast.error("Lỗi", "Không thể cập nhật trạng thái tham gia cuộc họp");
    }
  };

  const handleCreateOrUpdateMeeting = async (values) => {
    setSubmitting(true);
    try {
      const payload = {
        title: values.title,
        description: values.description || '',
        start_time: values.timeRange[0].toISOString(),
        end_time: values.timeRange[1].toISOString(),
        attendee_ids: Array.isArray(values.attendees) ? values.attendees : []
      };

      if (editingMeeting) {
        await meetingsAPI.update(editingMeeting.id, payload);
        toast.success("Thành công", "Đã cập nhật thông tin cuộc họp!");
      } else {
        await meetingsAPI.create(payload);
        toast.success("Thành công", "Đã tạo cuộc họp và tự động gửi thông báo đến các thành viên!");
      }

      setIsModalOpen(false);
      setEditingMeeting(null);
      form.resetFields();
      fetchData();
    } catch (err) {
      console.error("Error saving meeting:", err);
      toast.error("Lỗi", "Không thể lưu thông tin cuộc họp");
    } finally {
      setSubmitting(false);
    }
  };

  const handleOpenCreateModal = () => {
    setEditingMeeting(null);
    form.resetFields();
    setIsModalOpen(true);
  };

  const handleOpenEditModal = (meeting) => {
    setEditingMeeting(meeting);
    const attendeeUserIDs = Array.isArray(meeting.attendees) 
      ? meeting.attendees.map(a => a.user_id || a.userID || a.id) 
      : [];

    form.setFieldsValue({
      title: meeting.title,
      description: meeting.description,
      timeRange: [dayjs(meeting.start_time), dayjs(meeting.end_time)],
      attendees: attendeeUserIDs
    });
    setIsModalOpen(true);
  };

  const resolveEmployeeName = (idOrUserId) => {
    if (!idOrUserId || !Array.isArray(employees)) return 'Thành viên / Lãnh đạo';
    const emp = employees.find(e => e && (e.user_id === idOrUserId || e.id === idOrUserId || e.employee_code === idOrUserId));
    return emp ? emp.full_name : 'Thành viên';
  };

  return (
    <div style={{ padding: '24px', background: '#f8fafc', minHeight: 'calc(100vh - 64px)' }}>
      {/* Page Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px', flexWrap: 'wrap', gap: '16px' }}>
        <div>
          <Title level={3} style={{ margin: 0, color: '#0f172a', fontWeight: 700 }}>
            <CalendarOutlined style={{ marginRight: '10px', color: '#0284c7' }} />
            Phòng họp Nội bộ & Lịch làm việc
          </Title>
          <Text type="secondary">
            Quản lý lịch họp điều hành, thông báo thành viên và kết nối phòng họp Video WebRTC P2P
          </Text>
        </div>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          size="large"
          onClick={handleOpenCreateModal}
          style={{
            background: 'linear-gradient(135deg, #0284c7 0%, #0369a1 100%)',
            border: 'none',
            borderRadius: '8px',
            boxShadow: '0 4px 12px rgba(2, 132, 199, 0.3)',
            fontWeight: 600
          }}
        >
          Tạo cuộc họp mới
        </Button>
      </div>

      {/* Main Content */}
      {loading ? (
        <div style={{ textAlign: 'center', padding: '60px 0' }}>
          <Spin size="large" />
          <div style={{ marginTop: '16px', color: '#64748b' }}>Đang tải danh sách cuộc họp nội bộ...</div>
        </div>
      ) : (
        <div>
          {meetings.length === 0 ? (
            <Empty
              image={Empty.PRESENTED_IMAGE_SIMPLE}
              description="Chưa có cuộc họp nào được khởi tạo"
              style={{ background: '#fff', padding: '48px', borderRadius: '12px', border: '1px border-dashed #cbd5e1' }}
            >
              <Button type="primary" onClick={handleOpenCreateModal}>
                Tạo cuộc họp đầu tiên
              </Button>
            </Empty>
          ) : (
            <Row gutter={[20, 20]}>
              {meetings.map((m) => {
                if (!m) return null;
                const isScheduled = m.status === 'SCHEDULED';
                const isCancelled = m.status === 'CANCELLED';
                const isHost = m.host_id === currentUserID;

                const myRSVP = Array.isArray(m.attendees) 
                  ? m.attendees.find(a => (a.user_id === currentUserID || a.userID === currentUserID))?.rsvp_status || 'PENDING'
                  : 'PENDING';

                const cardActions = [
                  <Button
                    type="link"
                    icon={<VideoCameraOutlined />}
                    onClick={() => navigate(`/meetings/room/${m.id}`)}
                    style={{ color: '#0284c7', fontWeight: 600 }}
                  >
                    Vào phòng
                  </Button>
                ];

                if (isHost) {
                  cardActions.push(
                    <Button
                      type="link"
                      icon={<EditOutlined />}
                      onClick={() => handleOpenEditModal(m)}
                      style={{ color: '#8b5cf6', fontWeight: 600 }}
                    >
                      Sửa
                    </Button>
                  );
                } else {
                  cardActions.push(
                    <Button
                      type="link"
                      icon={<CheckOutlined />}
                      onClick={() => handleRSVP(m.id, 'ACCEPTED')}
                      style={{ color: myRSVP === 'ACCEPTED' ? '#10b981' : '#64748b', fontWeight: 600 }}
                    >
                      {myRSVP === 'ACCEPTED' ? 'Đã đồng ý' : 'Đồng ý'}
                    </Button>,
                    <Button
                      type="link"
                      icon={<CloseOutlined />}
                      onClick={() => handleRSVP(m.id, 'DECLINED')}
                      style={{ color: myRSVP === 'DECLINED' ? '#ef4444' : '#94a3b8', fontWeight: 500 }}
                    >
                      Từ chối
                    </Button>
                  );
                }

                return (
                  <Col xs={24} sm={12} lg={8} key={m.id}>
                    <Card
                      hoverable
                      style={{
                        height: '100%',
                        display: 'flex',
                        flexDirection: 'column',
                        justifyContent: 'space-between',
                        borderRadius: '12px',
                        border: '1px solid #e2e8f0',
                        boxShadow: '0 2px 8px rgba(0,0,0,0.03)',
                        overflow: 'hidden'
                      }}
                      bodyStyle={{ padding: '20px', flex: 1 }}
                      actions={cardActions}
                    >
                      <div>
                        {/* Status Header */}
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '14px' }}>
                          <Tag color={isScheduled ? 'green' : isCancelled ? 'red' : 'blue'} style={{ borderRadius: '4px', fontWeight: 600 }}>
                            {m.status || 'SCHEDULED'}
                          </Tag>
                          {m.google_event_id && (
                            <Tooltip title="Đã đồng bộ sang Google Calendar">
                              <Tag color="processing" icon={<GoogleOutlined />} style={{ borderRadius: '4px', margin: 0 }}>
                                Google Synced
                              </Tag>
                            </Tooltip>
                          )}
                        </div>

                        {/* Title & Description */}
                        <Title level={5} style={{ margin: '0 0 10px 0', fontWeight: 700, color: '#1e293b' }}>
                          {m.title || 'Không có tiêu đề'}
                        </Title>
                        <Paragraph
                          ellipsis={{ rows: 2 }}
                          style={{ color: '#64748b', fontSize: '13px', marginBottom: '16px' }}
                        >
                          {m.description || 'Không có mô tả chi tiết.'}
                        </Paragraph>

                        {/* Details */}
                        <Space direction="vertical" size={8} style={{ width: '100%', fontSize: '13px', color: '#475569' }}>
                          <div>
                            <ClockCircleOutlined style={{ color: '#0284c7', marginRight: '8px' }} />
                            <span>Thời gian: {dayjs(m.start_time).format('HH:mm:ss DD/MM/YYYY')}</span>
                          </div>
                          <div>
                            <UserOutlined style={{ color: '#0284c7', marginRight: '8px' }} />
                            <span>Người chủ trì: <strong>{resolveEmployeeName(m.host_id)}</strong></span>
                          </div>
                          <div>
                            <UserOutlined style={{ color: '#0284c7', marginRight: '8px' }} />
                            <span>Thành viên mời:</span>
                            <div style={{ marginTop: '6px', display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
                              {Array.isArray(m.attendees) && m.attendees.length > 0 ? (
                                m.attendees.map((att, i) => {
                                  const targetId = att.user_id || att.userID || att.id;
                                  const rsvp = att.rsvp_status || att.RSVPStatus || 'PENDING';
                                  const color = rsvp === 'ACCEPTED' ? 'green' : rsvp === 'DECLINED' ? 'red' : 'orange';
                                  return (
                                    <Tag key={i} color={color} style={{ fontSize: '11px', borderRadius: '4px', margin: '2px' }}>
                                      {resolveEmployeeName(targetId)} ({rsvp})
                                    </Tag>
                                  );
                                })
                              ) : (
                                <Text type="secondary" style={{ fontSize: '12px', fontStyle: 'italic' }}>
                                  Chưa có thành viên phụ
                                </Text>
                              )}
                            </div>
                          </div>
                        </Space>
                      </div>
                    </Card>
                  </Col>
                );
              })}
            </Row>
          )}
        </div>
      )}

      {/* Modal Create or Edit Meeting */}
      <Modal
        title={editingMeeting ? "Chỉnh sửa thông tin cuộc họp" : "Tạo cuộc họp mới"}
        open={isModalOpen}
        onCancel={() => { setIsModalOpen(false); setEditingMeeting(null); }}
        footer={null}
        destroyOnClose
      >
        <Form form={form} layout="vertical" onFinish={handleCreateOrUpdateMeeting} style={{ marginTop: '16px' }}>
          <Form.Item
            name="title"
            label="Tiêu đề cuộc họp"
            rules={[{ required: true, message: 'Vui lòng nhập tiêu đề cuộc họp!' }]}
          >
            <Input placeholder="Ví dụ: Họp triển khai kế hoạch sản lượng tháng 8" />
          </Form.Item>

          <Form.Item name="description" label="Nội dung / Mô tả">
            <Input.TextArea rows={3} placeholder="Nhập nội dung thảo luận hoặc chương trình họp..." />
          </Form.Item>

          <Form.Item
            name="timeRange"
            label="Thời gian bắt đầu & Kết thúc"
            rules={[{ required: true, message: 'Vui lòng chọn thời gian họp!' }]}
          >
            <DatePicker.RangePicker showTime format="YYYY-MM-DD HH:mm" style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item name="attendees" label="Mời thành viên tham gia">
            <Select
              mode="multiple"
              placeholder="Chọn nhân sự tham dự..."
              optionFilterProp="children"
              style={{ width: '100%' }}
            >
              {employees.map(emp => (
                <Option key={emp.user_id || emp.id} value={emp.user_id || emp.id}>
                  {emp.full_name} ({emp.employee_code || 'NV'})
                </Option>
              ))}
            </Select>
          </Form.Item>

          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '8px', marginTop: '24px' }}>
            <Button onClick={() => { setIsModalOpen(false); setEditingMeeting(null); }}>Hủy</Button>
            <Button type="primary" htmlType="submit" loading={submitting}>
              {editingMeeting ? "Cập nhật cuộc họp" : "Tạo cuộc họp & Báo NATS"}
            </Button>
          </div>
        </Form>
      </Modal>
    </div>
  );
}
