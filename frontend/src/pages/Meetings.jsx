import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { 
  Card, Button, Modal, Form, Input, DatePicker, 
  Select, Checkbox, Space, Tag, Typography, Row, Col, Tooltip, Empty 
} from 'antd';
import { 
  CalendarOutlined, VideoCameraOutlined, CheckOutlined, 
  UserOutlined, TeamOutlined, PlusOutlined, GoogleOutlined,
  ClockCircleOutlined, InfoCircleOutlined
} from '@ant-design/icons';
import { meetingsAPI, employeesAPI } from '../api/client.js';
import { useToast } from '../hooks/useToast.js';
import dayjs from 'dayjs';

const { Title, Text, Paragraph } = Typography;
const { RangePicker } = DatePicker;

export default function Meetings() {
  const navigate = useNavigate();
  const toast = useToast();
  const [form] = Form.useForm();
  
  const [meetings, setMeetings] = useState([]);
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [meetingsRes, employeesRes] = await Promise.all([
        meetingsAPI.list(),
        employeesAPI.list({ page_size: 200 })
      ]);
      
      console.log("Meetings API raw response:", meetingsRes);
      console.log("Employees API raw response:", employeesRes);
      
      // Defensively parse meetings list
      let mList = [];
      if (meetingsRes && meetingsRes.data) {
        if (Array.isArray(meetingsRes.data)) {
          mList = meetingsRes.data;
        } else if (Array.isArray(meetingsRes.data.data)) {
          mList = meetingsRes.data.data;
        }
      }
      setMeetings(mList);

      // Defensively parse employees list
      let eList = [];
      if (employeesRes && employeesRes.data && employeesRes.data.data) {
        if (Array.isArray(employeesRes.data.data.items)) {
          eList = employeesRes.data.data.items;
        } else if (Array.isArray(employeesRes.data.data)) {
          eList = employeesRes.data.data;
        }
      }
      console.log("Parsed employees list:", eList);
      setEmployees(eList);
    } catch (err) {
      console.error("Error fetching meetings or employees:", err);
      toast.error('Lỗi', 'Không thể tải danh sách cuộc họp');
      setMeetings([]);
      setEmployees([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleCreateMeeting = async (values) => {
    const { title, description, timeRange, attendeeIds, syncGoogle, enableWebRTC } = values;
    if (!timeRange || timeRange.length !== 2) {
      toast.error('Lỗi', 'Vui lòng chọn thời gian bắt đầu và kết thúc');
      return;
    }

    setSubmitting(true);
    try {
      await meetingsAPI.create({
        title,
        description: description || '',
        start_time: timeRange[0].toISOString(),
        end_time: timeRange[1].toISOString(),
        attendee_ids: attendeeIds || [],
        sync_google: !!syncGoogle,
        enable_webrtc: !!enableWebRTC
      });

      toast.success('Thành công', 'Đã tạo cuộc họp mới thành công');
      setShowModal(false);
      form.resetFields();
      fetchData();
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.error || 'Tạo cuộc họp thất bại');
    } finally {
      setSubmitting(false);
    }
  };

  const handleRSVP = async (meetingId, status) => {
    try {
      await meetingsAPI.rsvp(meetingId, status, 'Xác nhận tham gia');
      toast.success('Thành công', 'Đã cập nhật trạng thái tham gia cuộc họp');
      fetchData();
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.error || 'Không thể cập nhật trạng thái tham gia');
    }
  };

  const resolveUserName = (userId) => {
    if (!userId || !Array.isArray(employees)) return 'Lãnh đạo / Hệ thống';
    const emp = employees.find(e => e && e.user_id === userId);
    return emp ? emp.full_name : 'Lãnh đạo / Hệ thống';
  };

  const formatTime = (timeStr) => {
    if (!timeStr) return '—';
    try {
      const d = new Date(timeStr);
      return isNaN(d.getTime()) ? '—' : d.toLocaleString('vi-VN');
    } catch (e) {
      return '—';
    }
  };

  return (
    <div style={{ padding: '24px', maxWidth: '1200px', margin: '0 auto' }}>
      {/* Page Header */}
      <Row justify="space-between" align="middle" style={{ marginBottom: '32px' }}>
        <Col xs={24} md={18}>
          <Title level={3} style={{ margin: 0, fontWeight: 700 }}>
           Phòng họp
          </Title>
          <Paragraph style={{ color: '#64748b', margin: '6px 0 0 0', fontSize: '14px' }}>
          </Paragraph>
        </Col>
        <Col xs={24} md={6} style={{ textAlign: 'right', marginTop: '16px', mdStyle: { marginTop: 0 } }}>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            size="large"
            onClick={() => setShowModal(true)}
            style={{
              background: '#10b981',
              borderColor: '#10b981',
              borderRadius: '8px',
              fontWeight: 600,
              boxShadow: '0 4px 12px rgba(16, 185, 129, 0.25)'
            }}
          >
            Tạo cuộc họp mới
          </Button>
        </Col>
      </Row>

      {/* Main Grid */}
      {loading ? (
        <Card loading={true} style={{ borderRadius: '12px' }} />
      ) : (
        <div>
          {!Array.isArray(meetings) || meetings.length === 0 ? (
            <Empty 
              image={Empty.PRESENTED_IMAGE_SIMPLE}
              description={
                <span style={{ color: '#94a3b8' }}>Chưa có cuộc họp nào được lên lịch. Hãy tạo mới một cuộc họp!</span>
              }
              style={{
                padding: '60px 0',
                background: '#f8fafc',
                borderRadius: '12px',
                border: '1.5px dashed #cbd5e1'
              }}
            />
          ) : (
            <Row gutter={[20, 20]}>
              {meetings.map((m) => {
                if (!m) return null;
                const isScheduled = m.status === 'SCHEDULED';
                const isCancelled = m.status === 'CANCELLED';
                const isCompleted = m.status === 'COMPLETED';

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
                      actions={[
                        <Button
                          type="link"
                          icon={<VideoCameraOutlined />}
                          onClick={() => navigate(`/meetings/room/${m.id}`)}
                          style={{ color: '#0284c7', fontWeight: 600 }}
                        >
                          Vào phòng họp
                        </Button>,
                        <Button
                          type="link"
                          icon={<CheckOutlined />}
                          onClick={() => handleRSVP(m.id, 'ACCEPTED')}
                          style={{ color: '#10b981', fontWeight: 600 }}
                        >
                          Đồng ý
                        </Button>
                      ]}
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
                          ellipsis={{ rows: 2, expandable: true, symbol: 'Xem thêm' }}
                          style={{ fontSize: '13px', color: '#64748b', marginBottom: '16px', lineHeight: '1.5' }}
                        >
                          {m.description || 'Không có mô tả chi tiết.'}
                        </Paragraph>

                        {/* Meta Info */}
                        <div style={{ 
                          display: 'flex', 
                          flexDirection: 'column', 
                          gap: '10px', 
                          borderTop: '1px solid #f1f5f9', 
                          paddingTop: '14px', 
                          fontSize: '13px', 
                          color: '#475569' 
                        }}>
                          <div>
                            <ClockCircleOutlined style={{ marginRight: '8px', color: '#94a3b8' }} />
                            <strong>Thời gian:</strong> {formatTime(m.start_time)}
                          </div>
                          <div>
                            <UserOutlined style={{ marginRight: '8px', color: '#94a3b8' }} />
                            <strong>Người chủ trì:</strong> {resolveUserName(m.host_id)}
                          </div>

                          {/* Attendees */}
                          {m.attendees && Array.isArray(m.attendees) && m.attendees.length > 0 && (
                            <div style={{ display: 'flex', alignItems: 'flex-start', marginTop: '2px' }}>
                              <TeamOutlined style={{ marginRight: '8px', marginTop: '3px', color: '#94a3b8' }} />
                              <div style={{ flex: 1 }}>
                                <strong>Người tham gia:</strong>
                                <div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px', marginTop: '6px' }}>
                                  {m.attendees.map((att, idx) => {
                                    if (!att) return null;
                                    const isAccepted = att.rsvp_status === 'ACCEPTED';
                                    const isDeclined = att.rsvp_status === 'DECLINED';
                                    return (
                                      <Tag 
                                        key={att.user_id || idx} 
                                        color={isAccepted ? 'success' : isDeclined ? 'error' : 'default'}
                                        style={{ borderRadius: '12px', margin: 0, fontSize: '11px' }}
                                      >
                                        {resolveUserName(att.user_id)} ({att.rsvp_status || 'PENDING'})
                                      </Tag>
                                    );
                                  })}
                                </div>
                              </div>
                            </div>
                          )}
                        </div>
                      </div>
                    </Card>
                  </Col>
                );
              })}
            </Row>
          )}
        </div>
      )}

      <Modal
        title={
          <Title level={4} style={{ margin: 0 }}>
            Tạo Cuộc Họp
          </Title>
        }
        open={showModal}
        onCancel={() => {
          setShowModal(false);
          form.resetFields();
        }}
        footer={null}
        destroyOnClose
        width={550}
        style={{ borderRadius: '12px', overflow: 'hidden' }}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreateMeeting}
          initialValues={{
            syncGoogle: true,
            enableWebRTC: true,
            attendeeIds: []
          }}
          style={{ marginTop: '20px' }}
        >
          {/* Tiêu đề */}
          <Form.Item
            name="title"
            label={<strong>Tiêu đề cuộc họp</strong>}
            rules={[{ required: true, message: 'Vui lòng nhập tiêu đề cuộc họp!' }]}
          >
            <Input placeholder="Vd: Họp đánh giá năng lực P2" size="large" style={{ borderRadius: '6px' }} />
          </Form.Item>

          {/* Mô tả */}
          <Form.Item
            name="description"
            label={<strong>Mô tả nội dung</strong>}
          >
            <Input.TextArea 
              placeholder="Nội dung thảo luận chính, tài liệu đính kèm..." 
              rows={3} 
              style={{ borderRadius: '6px' }} 
            />
          </Form.Item>

          {/* Thời gian Bắt đầu & Kết thúc (RangePicker Antd) */}
          <Form.Item
            name="timeRange"
            label={<strong>Thời gian họp (Bắt đầu - Kết thúc)</strong>}
            rules={[{ required: true, message: 'Vui lòng chọn thời gian họp!' }]}
          >
            <RangePicker 
              showTime={{ format: 'HH:mm' }}
              format="DD/MM/YYYY HH:mm"
              size="large"
              placeholder={['Bắt đầu', 'Kết thúc']}
              style={{ width: '100%', borderRadius: '6px' }}
            />
          </Form.Item>

          {/* Debug Block */}
          <div style={{ background: '#f1f5f9', padding: '12px', borderRadius: '6px', marginBottom: '16px', fontSize: '12px' }}>
          </div>

          {/* Chọn Người tham gia */}
          <Form.Item
            name="attendeeIds"
            label={<strong>Chọn người tham gia họp</strong>}
            rules={[{ required: true, message: 'Vui lòng chọn ít nhất một người tham gia!' }]}
          >
            <Select
              mode="multiple"
              showSearch
              placeholder="Tìm kiếm và chọn nhân viên..."
              optionFilterProp="label"
              filterOption={(input, option) =>
                (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
              }
              size="large"
              style={{ width: '100%', borderRadius: '6px' }}
              options={employees
                .filter(emp => emp && emp.user_id)
                .map(emp => ({
                  value: emp.user_id,
                  label: `${emp.full_name} (${emp.code})`
                }))
              }
            />
          </Form.Item>

          {/* Cấu hình thêm */}
          <Form.Item label={<strong>Tùy chọn thiết lập</strong>}>
            <Space direction="vertical" style={{ width: '100%' }}>
              <Form.Item name="syncGoogle" valuePropName="checked" noStyle>
                <Checkbox>
                  Đồng bộ Google Calendar
                </Checkbox>
              </Form.Item>
              <Form.Item name="enableWebRTC" valuePropName="checked" noStyle>
                <Checkbox>
                  Tạo phòng họp Video
                </Checkbox>
              </Form.Item>
            </Space>
          </Form.Item>

          {/* Action Buttons */}
          <Form.Item style={{ margin: '24px 0 0 0', textAlign: 'right' }}>
            <Space size="middle">
              <Button 
                onClick={() => {
                  setShowModal(false);
                  form.resetFields();
                }}
                disabled={submitting}
                style={{ borderRadius: '6px' }}
              >
                Hủy
              </Button>
              <Button 
                type="primary" 
                htmlType="submit" 
                loading={submitting}
                style={{ 
                  background: '#10b981', 
                  borderColor: '#10b981', 
                  borderRadius: '6px',
                  fontWeight: 600
                }}
              >
                Tạo cuộc họp
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
