import { useState, useCallback } from 'react';
import {
  Upload, Button, Card, Tabs, Typography, Space, Tag, Alert,
  Progress, Input, Select, Slider, Divider, Row, Col, Tooltip,
  message as antMessage,
} from 'antd';
import {
  FilePdfOutlined, MergeCellsOutlined, CompressOutlined,
  HighlightOutlined, RetweetOutlined, DownloadOutlined,
  InboxOutlined, DeleteOutlined, CheckCircleOutlined,
  LoadingOutlined, FileProtectOutlined,
} from '@ant-design/icons';
import { pdfAPI } from '../api/client.js';

const { Title, Text, Paragraph } = Typography;
const { Dragger } = Upload;
const { Option } = Select;

/* ─────────── helper ─────────── */
const fmt = (bytes) => {
  if (!bytes) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / k ** i).toFixed(1)} ${sizes[i]}`;
};

const OperationBadge = ({ icon, label, color }) => (
  <Tag icon={icon} color={color} style={{ fontSize: 13, padding: '4px 12px', borderRadius: 20 }}>
    {label}
  </Tag>
);

/* ─────────── Result Card ─────────── */
function ResultCard({ result, loading }) {
  const [downloading, setDownloading] = useState(false);

  if (loading) {
    return (
      <Card style={cardStyle} className="pdf-result-card">
        <Space direction="vertical" align="center" style={{ width: '100%', padding: '24px 0' }}>
          <LoadingOutlined style={{ fontSize: 40, color: '#10b981' }} spin />
          <Text style={{ color: '#94a3b8' }}>Đang xử lý file...</Text>
          <Progress percent={null} status="active" strokeColor="#10b981" style={{ width: 300 }} />
        </Space>
      </Card>
    );
  }
  if (!result) return null;

  const downloadLink = result.output_path
    ? pdfAPI.downloadUrl(result.output_path)
    : null;

  const handleDownload = async () => {
    if (!result.output_path) return;
    setDownloading(true);
    try {
      const res = await pdfAPI.downloadFile(result.output_path);
      const url = window.URL.createObjectURL(new Blob([res.data], { type: 'application/pdf' }));
      const link = document.createElement('a');
      link.href = url;
      const fileName = result.output_path.split(/[/\\]/).pop() || 'result.pdf';
      link.setAttribute('download', fileName);
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
    } catch (err) {
      console.error("Download error:", err);
      if (downloadLink) {
        window.open(downloadLink, '_blank');
      } else {
        antMessage.error('Không thể tải file');
      }
    } finally {
      setDownloading(false);
    }
  };

  return (
    <Card
      style={{ ...cardStyle, border: '1px solid rgba(16,185,129,0.4)', background: 'rgba(16,185,129,0.05)' }}
      className="pdf-result-card"
    >
      <Space direction="vertical" style={{ width: '100%' }}>
        <Space>
          <CheckCircleOutlined style={{ color: '#10b981', fontSize: 22 }} />
          <Text strong style={{ color: '#10b981', fontSize: 16 }}>Xử lý thành công!</Text>
        </Space>
        <Paragraph style={{ color: '#94a3b8', marginBottom: 8 }}>{result.message}</Paragraph>
        {result.output_path && (
          <Text code style={{ color: '#64748b', fontSize: 11 }}>
            📁 {result.output_path}
          </Text>
        )}
        {result.output_path && (
          <Button
            type="primary"
            icon={<DownloadOutlined />}
            size="large"
            loading={downloading}
            onClick={handleDownload}
            style={{ background: '#10b981', borderColor: '#10b981', marginTop: 8 }}
          >
            Tải về file kết quả
          </Button>
        )}
      </Space>
    </Card>
  );
}

/* ─────────── Upload Zone ─────────── */
function UploadZone({ multiple, files, onChange, accept = '.pdf' }) {
  const props = {
    name: 'file',
    multiple,
    accept,
    beforeUpload: (file, fileList) => {
      onChange(multiple ? fileList : [file]);
      return false; // prevent auto-upload
    },
    onRemove: (file) => {
      onChange(files.filter(f => f.uid !== file.uid));
    },
    fileList: files,
    showUploadList: {
      showRemoveIcon: true,
      removeIcon: <DeleteOutlined style={{ color: '#ef4444' }} />,
    },
  };

  return (
    <Dragger {...props} style={draggerStyle}>
      <p className="ant-upload-drag-icon">
        <InboxOutlined style={{ color: '#10b981', fontSize: 40 }} />
      </p>
      <p style={{ color: '#e2e8f0', fontSize: 15, fontWeight: 600, marginBottom: 4 }}>
        Kéo thả file PDF vào đây
      </p>
      <p style={{ color: '#64748b', fontSize: 13 }}>
        {multiple ? 'Chọn nhiều file (Ctrl+Click)' : 'Chỉ chọn 1 file PDF'} — tối đa 32 MB mỗi file
      </p>
    </Dragger>
  );
}

/* ─────────── TAB 1: Gộp PDF ─────────── */
function MergeTab() {
  const [files, setFiles] = useState([]);
  const [outputName, setOutputName] = useState('merged_output');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);

  const handleMerge = async () => {
    if (files.length < 2) { antMessage.warning('Cần ít nhất 2 file PDF để gộp'); return; }
    setLoading(true); setResult(null);
    try {
      const nativeFiles = files.map(f => f.originFileObj || f);
      const res = await pdfAPI.merge(nativeFiles, outputName);
      setResult(res.data);
    } catch (e) {
      antMessage.error(e?.response?.data?.error || 'Gộp PDF thất bại');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Space direction="vertical" size={20} style={{ width: '100%' }}>
      <Alert
        message="Gộp nhiều file PDF thành một tài liệu duy nhất theo đúng thứ tự bạn chọn."
        type="info" showIcon
        style={{ background: 'rgba(59,130,246,0.08)', border: '1px solid rgba(59,130,246,0.2)', color: '#93c5fd' }}
      />
      <UploadZone multiple files={files} onChange={setFiles} />
      {files.length > 0 && (
        <div style={{ display: 'flex', gap: 12, alignItems: 'flex-end', flexWrap: 'wrap' }}>
          <div style={{ flex: 1, minWidth: 200 }}>
            <Text style={{ color: '#94a3b8', fontSize: 13, display: 'block', marginBottom: 6 }}>Tên file kết quả</Text>
            <Input
              value={outputName}
              onChange={e => setOutputName(e.target.value)}
              prefix={<FilePdfOutlined style={{ color: '#10b981' }} />}
              suffix={<Text style={{ color: '#64748b' }}>.pdf</Text>}
              placeholder="merged_output"
              style={inputStyle}
            />
          </div>
          <Button
            type="primary" icon={<MergeCellsOutlined />} size="large"
            onClick={handleMerge} loading={loading}
            style={{ background: '#10b981', borderColor: '#10b981', height: 42, paddingInline: 28 }}
          >
            Gộp {files.length} file
          </Button>
        </div>
      )}
      <ResultCard result={result} loading={loading} />
    </Space>
  );
}

/* ─────────── TAB 2: Nén PDF ─────────── */
function CompressTab() {
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);

  const handleCompress = async () => {
    if (!files.length) { antMessage.warning('Chưa chọn file PDF'); return; }
    setLoading(true); setResult(null);
    try {
      const file = files[0].originFileObj || files[0];
      const res = await pdfAPI.compress(file);
      setResult(res.data);
    } catch (e) {
      antMessage.error(e?.response?.data?.error || 'Nén PDF thất bại');
    } finally {
      setLoading(false);
    }
  };

  const sizeInfo = files[0]?.size
    ? `Kích thước: ${fmt(files[0].size)}`
    : null;

  return (
    <Space direction="vertical" size={20} style={{ width: '100%' }}>
      <Alert
        message="Tối ưu hóa và giảm dung lượng file PDF mà vẫn giữ nguyên chất lượng hiển thị."
        type="info" showIcon
        style={{ background: 'rgba(59,130,246,0.08)', border: '1px solid rgba(59,130,246,0.2)' }}
      />
      <UploadZone multiple={false} files={files} onChange={setFiles} />
      {files.length > 0 && (
        <Space>
          {sizeInfo && <Tag color="blue">{sizeInfo}</Tag>}
          <Button
            type="primary" icon={<CompressOutlined />} size="large"
            onClick={handleCompress} loading={loading}
            style={{ background: '#10b981', borderColor: '#10b981', height: 42, paddingInline: 28 }}
          >
            Nén PDF
          </Button>
        </Space>
      )}
      <ResultCard result={result} loading={loading} />
    </Space>
  );
}

/* ─────────── TAB 3: Watermark ─────────── */
function WatermarkTab() {
  const [files, setFiles] = useState([]);
  const [text, setText] = useState('BẢO MẬT - NỘI BỘ');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);

  const presets = ['BẢO MẬT', 'NỘI BỘ', 'DRAFT', 'CONFIDENTIAL', 'KHÔNG SAO CHÉP'];

  const handleWatermark = async () => {
    if (!files.length) { antMessage.warning('Chưa chọn file PDF'); return; }
    if (!text.trim()) { antMessage.warning('Nhập nội dung watermark'); return; }
    setLoading(true); setResult(null);
    try {
      const file = files[0].originFileObj || files[0];
      const res = await pdfAPI.watermark(file, text);
      setResult(res.data);
    } catch (e) {
      antMessage.error(e?.response?.data?.error || 'Thêm watermark thất bại');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Space direction="vertical" size={20} style={{ width: '100%' }}>
      <Alert
        message="Đóng dấu chìm (watermark) lên toàn bộ trang PDF. Watermark hiển thị xuyên suốt, không thể xóa."
        type="info" showIcon
        style={{ background: 'rgba(59,130,246,0.08)', border: '1px solid rgba(59,130,246,0.2)' }}
      />
      <UploadZone multiple={false} files={files} onChange={setFiles} />
      <div>
        <Text style={{ color: '#94a3b8', fontSize: 13, display: 'block', marginBottom: 8 }}>
          Nội dung Watermark
        </Text>
        <Input
          value={text}
          onChange={e => setText(e.target.value)}
          placeholder="Nhập nội dung watermark..."
          prefix={<HighlightOutlined style={{ color: '#10b981' }} />}
          style={{ ...inputStyle, marginBottom: 10 }}
          maxLength={50}
        />
        <Space wrap>
          <Text style={{ color: '#64748b', fontSize: 12 }}>Mẫu nhanh:</Text>
          {presets.map(p => (
            <Tag
              key={p}
              onClick={() => setText(p)}
              style={{ cursor: 'pointer', borderRadius: 20, padding: '2px 12px' }}
              color={text === p ? 'green' : 'default'}
            >
              {p}
            </Tag>
          ))}
        </Space>
      </div>
      {/* Preview */}
      {text && (
        <div style={{
          background: 'rgba(255,255,255,0.03)', border: '1px dashed rgba(255,255,255,0.1)',
          borderRadius: 12, padding: '32px 20px', textAlign: 'center', position: 'relative', overflow: 'hidden',
        }}>
          <div style={{
            position: 'absolute', inset: 0, display: 'flex', alignItems: 'center',
            justifyContent: 'center', transform: 'rotate(-45deg)',
          }}>
            <span style={{ fontSize: 36, fontWeight: 900, color: 'rgba(16,185,129,0.15)', whiteSpace: 'nowrap', letterSpacing: 4 }}>
              {text}
            </span>
          </div>
          <FilePdfOutlined style={{ fontSize: 40, color: '#475569' }} />
          <div style={{ color: '#475569', marginTop: 8, fontSize: 12 }}>Xem trước watermark</div>
        </div>
      )}
      {files.length > 0 && (
        <Button
          type="primary" icon={<HighlightOutlined />} size="large"
          onClick={handleWatermark} loading={loading}
          style={{ background: '#10b981', borderColor: '#10b981', height: 42, paddingInline: 28 }}
        >
          Thêm Watermark
        </Button>
      )}
      <ResultCard result={result} loading={loading} />
    </Space>
  );
}

/* ─────────── TAB 4: Xoay Trang ─────────── */
function RotateTab() {
  const [files, setFiles] = useState([]);
  const [angle, setAngle] = useState(90);
  const [pages, setPages] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);

  const handleRotate = async () => {
    if (!files.length) { antMessage.warning('Chưa chọn file PDF'); return; }
    setLoading(true); setResult(null);
    try {
      const file = files[0].originFileObj || files[0];
      const res = await pdfAPI.rotate(file, angle, pages);
      setResult(res.data);
    } catch (e) {
      antMessage.error(e?.response?.data?.error || 'Xoay trang thất bại');
    } finally {
      setLoading(false);
    }
  };

  const angleOptions = [
    { value: 90, label: '90° (Xoay phải)', icon: '↻' },
    { value: 180, label: '180° (Lộn ngược)', icon: '↕' },
    { value: 270, label: '270° (Xoay trái)', icon: '↺' },
  ];

  return (
    <Space direction="vertical" size={20} style={{ width: '100%' }}>
      <Alert
        message='Xoay trang PDF theo góc bạn chọn. Để "Trang cần xoay" trống = xoay tất cả các trang.'
        type="info" showIcon
        style={{ background: 'rgba(59,130,246,0.08)', border: '1px solid rgba(59,130,246,0.2)' }}
      />
      <UploadZone multiple={false} files={files} onChange={setFiles} />
      <Row gutter={16}>
        <Col xs={24} md={12}>
          <Text style={{ color: '#94a3b8', fontSize: 13, display: 'block', marginBottom: 8 }}>Góc xoay</Text>
          <div style={{ display: 'flex', gap: 8 }}>
            {angleOptions.map(opt => (
              <div
                key={opt.value}
                onClick={() => setAngle(opt.value)}
                style={{
                  flex: 1, padding: '12px 8px', borderRadius: 10, textAlign: 'center', cursor: 'pointer',
                  border: `2px solid ${angle === opt.value ? '#10b981' : 'rgba(255,255,255,0.08)'}`,
                  background: angle === opt.value ? 'rgba(16,185,129,0.1)' : 'rgba(255,255,255,0.02)',
                  transition: 'all 0.2s',
                }}
              >
                <div style={{ fontSize: 22, color: angle === opt.value ? '#10b981' : '#475569' }}>{opt.icon}</div>
                <div style={{ fontSize: 11, color: angle === opt.value ? '#10b981' : '#64748b', marginTop: 4 }}>
                  {opt.label.split(' ')[0]}
                </div>
              </div>
            ))}
          </div>
        </Col>
        <Col xs={24} md={12}>
          <Text style={{ color: '#94a3b8', fontSize: 13, display: 'block', marginBottom: 8 }}>
            Trang cần xoay <Text style={{ color: '#64748b', fontSize: 11 }}>(để trống = tất cả)</Text>
          </Text>
          <Input
            value={pages}
            onChange={e => setPages(e.target.value)}
            placeholder="VD: 1,3,5 hoặc 1-3"
            prefix={<RetweetOutlined style={{ color: '#10b981' }} />}
            style={inputStyle}
          />
        </Col>
      </Row>
      {files.length > 0 && (
        <Button
          type="primary" icon={<RetweetOutlined />} size="large"
          onClick={handleRotate} loading={loading}
          style={{ background: '#10b981', borderColor: '#10b981', height: 42, paddingInline: 28 }}
        >
          Xoay {angle}°
        </Button>
      )}
      <ResultCard result={result} loading={loading} />
    </Space>
  );
}

/* ─────────── Styles ─────────── */
const cardStyle = {
  background: 'rgba(255,255,255,0.03)',
  border: '1px solid rgba(255,255,255,0.07)',
  borderRadius: 14,
};

const draggerStyle = {
  background: 'rgba(16,185,129,0.04)',
  border: '1.5px dashed rgba(16,185,129,0.35)',
  borderRadius: 12,
};

const inputStyle = {
  background: 'rgba(255,255,255,0.05)',
  border: '1px solid rgba(255,255,255,0.1)',
  color: '#e2e8f0',
  borderRadius: 8,
};

/* ─────────── Main Page ─────────── */
export default function PDFToolsPage() {
  const tabs = [
    {
      key: 'merge',
      label: (
        <Space>
          <MergeCellsOutlined />
          Gộp PDF
        </Space>
      ),
      children: <MergeTab />,
    },
    {
      key: 'compress',
      label: (
        <Space>
          <CompressOutlined />
          Nén PDF
        </Space>
      ),
      children: <CompressTab />,
    },
    {
      key: 'watermark',
      label: (
        <Space>
          <HighlightOutlined />
          Watermark
        </Space>
      ),
      children: <WatermarkTab />,
    },
    {
      key: 'rotate',
      label: (
        <Space>
          <RetweetOutlined />
          Xoay Trang
        </Space>
      ),
      children: <RotateTab />,
    },
  ];

  return (
    <div style={{ padding: '0 0 40px' }}>
      {/* Header */}
      <div style={{
        display: 'flex', alignItems: 'flex-start', gap: 16,
        marginBottom: 32, padding: '4px 0',
      }}>
        <div style={{
          width: 52, height: 52, borderRadius: 14, flexShrink: 0,
          background: 'linear-gradient(135deg, #10b981, #059669)',
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          boxShadow: '0 0 20px rgba(16,185,129,0.35)',
        }}>
          <FileProtectOutlined style={{ fontSize: 26, color: '#fff' }} />
        </div>
        <div>
          <Title level={3} style={{ color: '#e2e8f0', margin: 0, fontSize: 22, fontWeight: 700 }}>
            Công cụ xử lý PDF
          </Title>
          <Text style={{ color: '#64748b', fontSize: 14 }}>
            Gộp, nén, đóng dấu watermark và xoay trang PDF ngay trên trình duyệt
          </Text>
        </div>
      </div>

      {/* Feature overview */}
      <Row gutter={12} style={{ marginBottom: 28 }}>
        {[
          { icon: <MergeCellsOutlined />, label: 'Gộp PDF', desc: 'Kết hợp nhiều file', color: '#3b82f6' },
          { icon: <CompressOutlined />, label: 'Nén PDF', desc: 'Giảm dung lượng', color: '#8b5cf6' },
          { icon: <HighlightOutlined />, label: 'Watermark', desc: 'Đóng dấu bảo mật', color: '#f59e0b' },
          { icon: <RetweetOutlined />, label: 'Xoay trang', desc: 'Xoay 90°/180°/270°', color: '#ec4899' },
        ].map(f => (
          <Col xs={12} md={6} key={f.label}>
            <Card
              size="small"
              style={{
                background: `rgba(${f.color === '#3b82f6' ? '59,130,246' : f.color === '#8b5cf6' ? '139,92,246' : f.color === '#f59e0b' ? '245,158,11' : '236,72,153'},0.07)`,
                border: `1px solid ${f.color}30`,
                borderRadius: 12, textAlign: 'center', cursor: 'pointer',
              }}
              hoverable
            >
              <div style={{ fontSize: 24, color: f.color, marginBottom: 4 }}>{f.icon}</div>
              <div style={{ color: '#e2e8f0', fontWeight: 600, fontSize: 13 }}>{f.label}</div>
              <div style={{ color: '#64748b', fontSize: 11 }}>{f.desc}</div>
            </Card>
          </Col>
        ))}
      </Row>

      {/* Main tabs */}
      <Card style={cardStyle}>
        <Tabs
          items={tabs}
          size="large"
          tabBarStyle={{ borderBottom: '1px solid rgba(255,255,255,0.07)', marginBottom: 24 }}
          style={{ color: '#e2e8f0' }}
        />
      </Card>

      <style>{`
        .ant-upload-drag {
          background: rgba(16,185,129,0.04) !important;
          border-color: rgba(16,185,129,0.35) !important;
        }
        .ant-upload-drag:hover {
          border-color: #10b981 !important;
        }
        .ant-tabs-tab {
          color: #64748b !important;
        }
        .ant-tabs-tab-active .ant-tabs-tab-btn {
          color: #10b981 !important;
        }
        .ant-tabs-ink-bar {
          background: #10b981 !important;
        }
        .ant-card-body {
          padding: 20px;
        }
        .ant-upload-list-item-name {
          color: #94a3b8 !important;
        }
        .ant-upload-list-item {
          color: #94a3b8 !important;
        }
        .ant-input {
          background: rgba(255,255,255,0.05) !important;
          border-color: rgba(255,255,255,0.1) !important;
          color: #e2e8f0 !important;
        }
        .ant-input:focus, .ant-input:hover {
          border-color: #10b981 !important;
        }
        .ant-input-affix-wrapper {
          background: rgba(255,255,255,0.05) !important;
          border-color: rgba(255,255,255,0.1) !important;
        }
        .ant-input-affix-wrapper:focus-within, .ant-input-affix-wrapper:hover {
          border-color: #10b981 !important;
        }
        .ant-input-affix-wrapper .ant-input {
          background: transparent !important;
        }
        .ant-tag {
          cursor: pointer;
        }
      `}</style>
    </div>
  );
}
