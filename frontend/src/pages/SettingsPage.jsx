import { useState, useEffect } from 'react';
import { settingsAPI } from '../api/client.js';
import { useToast } from '../hooks/useToast.js';

export function SettingsPage() {
  const toast = useToast();
  const [settings, setSettings] = useState([]);
  const [loading, setLoading] = useState(false);
  const [savingKey, setSavingKey] = useState(null);
  const [editValues, setEditValues] = useState({});

  useEffect(() => {
    loadSettings();
  }, []);

  const loadSettings = async () => {
    setLoading(true);
    try {
      const res = await settingsAPI.list();
      const items = res.data?.data || [];
      setSettings(items);
      
      // Initialize edit values state
      const vals = {};
      items.forEach(item => {
        vals[item.key] = item.value;
      });
      setEditValues(vals);
    } catch {
      toast.error('Lỗi', 'Không thể tải danh sách cấu hình hệ thống');
    } finally {
      setLoading(false);
    }
  };

  const handleValueChange = (key, value) => {
    setEditValues(p => ({ ...p, [key]: value }));
  };

  const handleSave = async (item) => {
    const newVal = editValues[item.key];
    if (newVal === undefined || newVal === item.value) return;

    setSavingKey(item.key);
    try {
      await settingsAPI.update(item.key, {
        value: String(newVal),
        description: item.description || '',
      });
      toast.success('Thành công', `Đã cập nhật cấu hình ${item.key}`);
      loadSettings(); // Reload to get fresh state
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể lưu cấu hình');
    } finally {
      setSavingKey(null);
    }
  };

  return (
    <div style={{ maxWidth: 800 }}>
      <div className="page-header">
        <div>
          <h2 className="page-title">Cấu hình hệ thống</h2>
          <p className="page-subtitle">Quản lý tham số tính lương toàn hệ thống</p>
        </div>
      </div>

      <div className="card">
        {loading ? (
          <div className="loading-center"><span className="spinner" /> Đang tải...</div>
        ) : settings.length > 0 ? (
          <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
            {settings.map((s) => (
              <div key={s.key} style={{ display: 'flex', flexDirection: 'column', gap: 8, paddingBottom: 20, borderBottom: '1px solid var(--border-color)' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: 16 }}>
                  <div>
                    <div style={{ fontWeight: 600, fontSize: 15, color: 'var(--text-primary)' }}>{s.key}</div>
                    <div style={{ fontSize: 13, color: 'var(--text-muted)', marginTop: 4 }}>{s.description || 'Chưa có mô tả'}</div>
                  </div>
                  <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                    <input 
                      className="form-control" 
                      value={editValues[s.key] ?? ''} 
                      onChange={e => handleValueChange(s.key, e.target.value)} 
                      style={{ width: 160, textAlign: 'right' }}
                    />
                    <button 
                      className="btn btn-primary btn-sm" 
                      onClick={() => handleSave(s)}
                      disabled={savingKey === s.key || editValues[s.key] === s.value}
                    >
                      {savingKey === s.key ? <span className="spinner spinner-sm" /> : 'Lưu'}
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="empty-state">
            <div className="empty-state-title">Không tìm thấy cấu hình nào</div>
          </div>
        )}
      </div>
    </div>
  );
}

export default SettingsPage;
