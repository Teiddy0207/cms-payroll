import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { authAPI } from '../api/client.js';
import { useToast } from '../hooks/useToast.js';

export function LoginPage() {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const toast = useToast();

  const handleFinish = async (values) => {
    if (!values.identifier || !values.password) {
      toast.error('Thiếu thông tin', 'Vui lòng nhập đầy đủ tên đăng nhập và mật khẩu');
      return;
    }
    setLoading(true);
    try {
      const res = await authAPI.login(values.identifier, values.password);
      const token = res.data?.data?.access_token || res.data?.access_token;
      if (token) {
        localStorage.setItem('auth_token', token);
        toast.success('Đăng nhập thành công', 'Chào mừng trở lại!');
        navigate('/dashboard');
      } else {
        toast.error('Lỗi', 'Không nhận được token từ server');
      }
    } catch (err) {
      const msg = err.response?.data?.message || 'Tên đăng nhập hoặc mật khẩu không đúng';
      toast.error('Đăng nhập thất bại', msg);
    } finally {
      setLoading(false);
    }
  };

  const handleFinishFailed = () => {
    toast.error('Thiếu thông tin', 'Vui lòng nhập đầy đủ tên đăng nhập và mật khẩu');
  };

  return (
    <div className="login-page">
      <div className="login-bg" />

      <div className="login-card">
        <div className="login-logo" style={{ justifyContent: 'center' }}>
          <div className="login-logo-text"><span>CMS</span> Payroll</div>
        </div>

        <Form
          id="login-form"
          name="login"
          layout="vertical"
          onFinish={handleFinish}
          onFinishFailed={handleFinishFailed}
          requiredMark={false}
          size="large"
        >
          <Form.Item
            label={<span className="form-label">Tên đăng nhập / Email</span>}
            name="identifier"
            rules={[{ required: true, message: 'Vui lòng nhập tên đăng nhập / Email' }]}
            style={{ marginBottom: 20 }}
          >
            <Input
              id="login-identifier"
              prefix={<UserOutlined style={{ color: 'rgba(16, 185, 129, 0.45)' }} />}
              placeholder="Nhập tên đăng nhập hoặc email"
              autoComplete="username"
              autoFocus
            />
          </Form.Item>

          <Form.Item
            label={<span className="form-label">Mật khẩu</span>}
            name="password"
            rules={[{ required: true, message: 'Vui lòng nhập mật khẩu' }]}
            style={{ marginBottom: 24 }}
          >
            <Input.Password
              id="login-password"
              prefix={<LockOutlined style={{ color: 'rgba(16, 185, 129, 0.45)' }} />}
              placeholder="Nhập mật khẩu"
              autoComplete="current-password"
            />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0 }}>
            <Button
              id="login-submit-btn"
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              style={{
                height: 44,
                fontSize: 15,
                fontWeight: 600,
                background: 'var(--accent)',
                borderColor: 'var(--accent)',
              }}
            >
              {loading ? 'Đang đăng nhập...' : 'Đăng nhập'}
            </Button>
          </Form.Item>
        </Form>

        <div className="login-footer">
          developed by Quang Anh &copy; 2026
        </div>
      </div>
    </div>
  );
}

export default LoginPage;

