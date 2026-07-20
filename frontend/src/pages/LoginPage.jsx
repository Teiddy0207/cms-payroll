import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authAPI } from '../api/client.js';
import { useToast } from '../hooks/useToast.js';

export function LoginPage() {
  const [form, setForm] = useState({ identifier: '', password: '' });
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const toast = useToast();

  const handleChange = (e) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!form.identifier || !form.password) {
      toast.error('Thiếu thông tin', 'Vui lòng nhập đầy đủ tên đăng nhập và mật khẩu');
      return;
    }
    setLoading(true);
    try {
      const res = await authAPI.login(form.identifier, form.password);
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

  return (
    <div className="login-page">
      <div className="login-bg" />
      {/* Decorative circles */}
      <div style={{
        position: 'absolute', width: 400, height: 400, borderRadius: '50%',
        background: 'radial-gradient(circle, rgba(79,142,247,0.06), transparent)',
        top: '-100px', right: '-100px', pointerEvents: 'none'
      }} />
      <div style={{
        position: 'absolute', width: 300, height: 300, borderRadius: '50%',
        background: 'radial-gradient(circle, rgba(79,142,247,0.04), transparent)',
        bottom: '-50px', left: '-50px', pointerEvents: 'none'
      }} />

      <div className="login-card">
        <div className="login-logo">
          <div className="login-logo-text"><span>CMS</span> Payroll</div>
        </div>
        <h1 className="login-title">Chào mừng trở lại</h1>
        <p className="login-subtitle">Đăng nhập để quản lý hệ thống tính lương</p>

        <form className="login-form" onSubmit={handleSubmit} id="login-form">
          <div className="form-group">
            <label className="form-label" htmlFor="login-identifier">
              Tên đăng nhập / Email
            </label>
            <input
              id="login-identifier"
              type="text"
              name="identifier"
              className="form-control"
              placeholder="Nhập tên đăng nhập hoặc email"
              value={form.identifier}
              onChange={handleChange}
              autoComplete="username"
              autoFocus
            />
          </div>
          <div className="form-group">
            <label className="form-label" htmlFor="login-password">
              Mật khẩu
            </label>
            <input
              id="login-password"
              type="password"
              name="password"
              className="form-control"
              placeholder="Nhập mật khẩu"
              value={form.password}
              onChange={handleChange}
              autoComplete="current-password"
            />
          </div>
          <button
            type="submit"
            className="btn btn-primary btn-lg"
            disabled={loading}
            id="login-submit-btn"
            style={{ width: '100%', justifyContent: 'center', marginTop: 4 }}
          >
            {loading ? (
              <><span className="spinner spinner-sm" style={{ borderTopColor: '#fff' }} /> Đang đăng nhập...</>
            ) : (
              'Đăng nhập'
            )}
          </button>
        </form>

        <div className="login-footer">
          Hệ thống quản lý tính lương &copy; 2024
        </div>
      </div>
    </div>
  );
}

export default LoginPage;
