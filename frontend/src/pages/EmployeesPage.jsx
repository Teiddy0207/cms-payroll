import { useState, useEffect, useCallback } from 'react';
import {
  employeesAPI, jobPositionsAPI, departmentsAPI, competenciesAPI, calculatorAPI
} from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { ConfirmDialog } from '../components/ConfirmDialog.jsx';
import { Pagination } from '../components/Pagination.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';

const PAGE_SIZE = 20;

const emptyForm = {
  code: '', full_name: '', phone: '', gender: '', date_of_birth: '',
  position_id: '', department_id: ''
};

function getInitials(name) {
  if (!name) return '?';
  const words = name.trim().split(' ');
  return words[words.length - 1].charAt(0).toUpperCase();
}

function formatDate(str) {
  if (!str) return '—';
  try { return new Date(str).toLocaleDateString('vi-VN'); } catch { return str; }
}

export function EmployeesPage() {
  const toast = useToast();
  const [employees, setEmployees] = useState([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(false);

  const [positions, setPositions] = useState([]);
  const [departments, setDepartments] = useState([]);
  const [allCompetencies, setAllCompetencies] = useState([]);

  // Modals
  const [createOpen, setCreateOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [detailOpen, setDetailOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [salaryOpen, setSalaryOpen] = useState(false);

  const [form, setForm] = useState(emptyForm);
  const [selected, setSelected] = useState(null);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);

  // Detail data
  const [detailData, setDetailData] = useState(null);
  const [empCompetencies, setEmpCompetencies] = useState([]);
  const [detailLoading, setDetailLoading] = useState(false);

  // Salary
  const [salaryPeriod, setSalaryPeriod] = useState(() => {
    const d = new Date();
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
  });
  const [salaryData, setSalaryData] = useState(null);
  const [salaryLoading, setSalaryLoading] = useState(false);

  // Competency management
  const [competenciesOpen, setCompetenciesOpen] = useState(false);
  const [selectedCompetencyToAssign, setSelectedCompetencyToAssign] = useState('');
  const [competenciesLoading, setCompetenciesLoading] = useState(false);
  const [assigning, setAssigning] = useState(false);

  useEffect(() => {
    loadMeta();
  }, []);

  useEffect(() => {
    const timer = setTimeout(() => { setPage(1); loadEmployees(1, search); }, 350);
    return () => clearTimeout(timer);
  }, [search]);

  useEffect(() => {
    loadEmployees(page, search);
  }, [page]);

  const loadMeta = async () => {
    try {
      const [posRes, deptRes, compRes] = await Promise.allSettled([
        jobPositionsAPI.list({ page_size: 200 }),
        departmentsAPI.list({ page_size: 200 }),
        competenciesAPI.list({ page_size: 200 }),
      ]);
      if (posRes.status === 'fulfilled') setPositions(posRes.value.data?.data?.items || []);
      if (deptRes.status === 'fulfilled') setDepartments(deptRes.value.data?.data?.items || []);
      if (compRes.status === 'fulfilled') setAllCompetencies(compRes.value.data?.data?.items || []);
    } catch { /**/ }
  };

  const loadEmployees = useCallback(async (p, s) => {
    setLoading(true);
    try {
      const params = { page_number: p, page_size: PAGE_SIZE };
      if (s) params.search = s;
      const res = await employeesAPI.list(params);
      const d = res.data?.data;
      setEmployees(d?.items || []);
      setTotal(d?.total_items || 0);
    } catch {
      toast.error('Lỗi', 'Không thể tải danh sách nhân viên');
    } finally {
      setLoading(false);
    }
  }, []);

  const openCreate = () => { setForm(emptyForm); setCreateOpen(true); };

  const openEdit = (emp) => {
    setSelected(emp);
    setForm({
      code: emp.code || '',
      full_name: emp.full_name || '',
      phone: emp.phone || '',
      gender: emp.gender || '',
      date_of_birth: emp.date_of_birth ? emp.date_of_birth.substring(0, 10) : '',
      position_id: emp.position_id || emp.job_position?.id || '',
      department_id: emp.department_id || emp.department?.id || '',
    });
    setEditOpen(true);
  };

  const openDetail = async (emp) => {
    setSelected(emp);
    setDetailOpen(true);
    setDetailLoading(true);
    setDetailData(null);
    setEmpCompetencies([]);
    try {
      const [detRes, compRes] = await Promise.allSettled([
        employeesAPI.get(emp.id),
        employeesAPI.getCompetencies(emp.id),
      ]);
      if (detRes.status === 'fulfilled') setDetailData(detRes.value.data?.data);
      if (compRes.status === 'fulfilled') setEmpCompetencies(compRes.value.data?.data?.items || compRes.value.data?.data || []);
    } catch { /**/ } finally {
      setDetailLoading(false);
    }
  };

  const openDelete = (emp) => { setSelected(emp); setDeleteOpen(true); };

  const openSalary = async (emp) => {
    setSelected(emp);
    setSalaryOpen(true);
    setSalaryData(null);
    setSalaryLoading(true);
    try {
      const res = await calculatorAPI.preview(emp.id, salaryPeriod);
      setSalaryData(res.data?.data);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể tải dữ liệu lương');
    } finally {
      setSalaryLoading(false);
    }
  };

  const openManageCompetencies = async (emp) => {
    setSelected(emp);
    setCompetenciesOpen(true);
    setCompetenciesLoading(true);
    setSelectedCompetencyToAssign('');
    try {
      const res = await employeesAPI.getCompetencies(emp.id);
      setEmpCompetencies(res.data?.data?.items || res.data?.data || []);
    } catch {
      toast.error('Lỗi', 'Không thể tải danh sách năng lực của nhân viên');
    } finally {
      setCompetenciesLoading(false);
    }
  };

  const handleAssignCompetency = async () => {
    if (!selectedCompetencyToAssign) return;
    setAssigning(true);
    try {
      await employeesAPI.assignCompetency(selected.id, {
        competency_ids: [selectedCompetencyToAssign]
      });
      toast.success('Thành công', 'Đã gán năng lực cho nhân viên');
      setSelectedCompetencyToAssign('');
      const res = await employeesAPI.getCompetencies(selected.id);
      setEmpCompetencies(res.data?.data?.items || res.data?.data || []);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể gán năng lực');
    } finally {
      setAssigning(false);
    }
  };

  const handleRemoveCompetency = async (competencyId) => {
    try {
      await employeesAPI.removeCompetency(selected.id, competencyId);
      toast.success('Thành công', 'Đã gỡ năng lực khỏi nhân viên');
      const res = await employeesAPI.getCompetencies(selected.id);
      setEmpCompetencies(res.data?.data?.items || res.data?.data || []);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể gỡ năng lực');
    }
  };

  const loadSalary = async () => {
    if (!selected) return;
    setSalaryLoading(true);
    setSalaryData(null);
    try {
      const res = await calculatorAPI.preview(selected.id, salaryPeriod);
      setSalaryData(res.data?.data);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể tải dữ liệu lương');
    } finally {
      setSalaryLoading(false);
    }
  };

  const handleCreate = async () => {
    if (!form.full_name) { toast.error('Thiếu thông tin', 'Vui lòng nhập họ tên'); return; }
    setSaving(true);
    try {
      await employeesAPI.create(form);
      toast.success('Thành công', 'Đã thêm nhân viên mới');
      setCreateOpen(false);
      loadEmployees(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể thêm nhân viên');
    } finally {
      setSaving(false);
    }
  };

  const handleEdit = async () => {
    setSaving(true);
    try {
      await employeesAPI.update(selected.id, form);
      toast.success('Thành công', 'Đã cập nhật thông tin nhân viên');
      setEditOpen(false);
      loadEmployees(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể cập nhật');
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    setDeleting(true);
    try {
      await employeesAPI.delete(selected.id);
      toast.success('Đã xóa', 'Nhân viên đã được xóa khỏi hệ thống');
      setDeleteOpen(false);
      loadEmployees(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể xóa nhân viên');
    } finally {
      setDeleting(false);
    }
  };

  const handleFormChange = (e) => {
    setForm((p) => ({ ...p, [e.target.name]: e.target.value }));
  };

  const columns = [
    {
      key: 'full_name',
      title: 'Nhân viên',
      render: (val, row) => (
        <div className="avatar-cell">
          <div className="avatar-info">
            <div className="name">{val || '—'}</div>
            <div className="code">{row.code || 'Chưa có mã'}</div>
          </div>
        </div>
      )
    },
    {
      key: 'phone',
      title: 'Số điện thoại',
      render: (val) => val || <span className="text-muted">—</span>
    },
    {
      key: 'gender',
      title: 'Giới tính',
      render: (val) => val
        ? <Badge color={val === 'male' ? 'blue' : 'cyan'}>{val === 'male' ? 'Nam' : 'Nữ'}</Badge>
        : <Badge color="gray">—</Badge>
    },
    {
      key: 'job_position',
      title: 'Vị trí',
      render: (val, row) => {
        const pos = val || positions.find(p => p.id === row.position_id);
        return pos ? <Badge color="yellow">{pos.name}</Badge> : <span className="text-muted">—</span>;
      }
    },
    {
      key: 'id',
      title: 'Thao tác',
      style: { width: 260 },
      render: (_, row) => (
        <div className="td-actions">
          <button className="btn btn-secondary btn-sm" id={`emp-detail-${row.id}`} onClick={() => openDetail(row)}>
            Chi tiết
          </button>
          <button className="btn btn-secondary btn-sm" id={`emp-comp-${row.id}`} onClick={() => openManageCompetencies(row)}>
            Năng lực
          </button>
          <button className="btn btn-secondary btn-sm" id={`emp-edit-${row.id}`} onClick={() => openEdit(row)}>
            Sửa
          </button>
          <button className="btn btn-danger btn-sm" id={`emp-delete-${row.id}`} onClick={() => openDelete(row)}>
            Xóa
          </button>
        </div>
      )
    }
  ];

  const availableCompetencies = allCompetencies.filter(
    comp => !empCompetencies.some(assigned => (assigned.competency_id || assigned.id) === comp.id)
  );



  return (
    <div>
      {/* Toolbar */}
      <div className="page-header">
        <div>
          <h2 className="page-title">Nhân viên</h2>
          <p className="page-subtitle">Quản lý thông tin nhân sự ({total} nhân viên)</p>
        </div>
        <button className="btn btn-primary" id="add-employee-btn" onClick={openCreate}>
          + Thêm nhân viên
        </button>
      </div>

      <div className="table-container">
        <div className="table-toolbar">
          <div className="search-bar">
            <input
              className="form-control"
              placeholder="Tìm kiếm nhân viên..."
              value={search}
              onChange={e => setSearch(e.target.value)}
              id="employee-search"
            />
          </div>
          <div className="text-muted" style={{ fontSize: 13 }}>
            {total} kết quả
          </div>
        </div>
        <Table columns={columns} data={employees} loading={loading} emptyMessage="Chưa có nhân viên nào" />
        <Pagination page={page} totalItems={total} pageSize={PAGE_SIZE} onPageChange={setPage} />
      </div>

      {/* Create Modal */}
      <EmployeeForm title="Thêm nhân viên" onSubmit={handleCreate} isOpen={createOpen} onClose={() => setCreateOpen(false)} form={form} handleFormChange={handleFormChange} saving={saving} positions={positions} departments={departments} />

      {/* Edit Modal */}
      <EmployeeForm title="Sửa thông tin nhân sự" onSubmit={handleEdit} isOpen={editOpen} onClose={() => setEditOpen(false)} form={form} handleFormChange={handleFormChange} saving={saving} positions={positions} departments={departments} />

      {/* Detail Modal */}
      <Modal isOpen={detailOpen} onClose={() => setDetailOpen(false)} title="Chi tiết nhân viên" size="lg"
        footer={
          <>
            <button className="btn btn-secondary" onClick={() => setDetailOpen(false)}>Đóng</button>
            <button className="btn btn-primary" onClick={() => { setDetailOpen(false); openSalary(selected); }}>
              Xem lương
            </button>
          </>
        }
      >
        {detailLoading ? (
          <div className="loading-center"><span className="spinner" /> Đang tải...</div>
        ) : detailData ? (
          <div>
            <div style={{ display: 'flex', alignItems: 'center', gap: 18, marginBottom: 24 }}>
              <div>
                <div style={{ fontSize: 18, fontWeight: 700 }}>{detailData.full_name}</div>
                <div className="text-muted" style={{ fontSize: 13 }}>{detailData.code}</div>
              </div>
            </div>
            <div className="detail-grid">
              <div className="detail-item">
                <label>Số điện thoại</label>
                <div className="value">{detailData.phone || '—'}</div>
              </div>
              <div className="detail-item">
                <label>Giới tính</label>
                <div className="value">
                  {detailData.gender === 'male' ? 'Nam' : detailData.gender === 'female' ? 'Nữ' : '—'}
                </div>
              </div>
              <div className="detail-item">
                <label>Ngày sinh</label>
                <div className="value">{formatDate(detailData.date_of_birth)}</div>
              </div>
              <div className="detail-item">
                <label>Vị trí công việc</label>
                <div className="value">
                  {detailData.job_position?.name || '—'}
                </div>
              </div>
            </div>
            <div className="divider" />
            <div className="section-title">Năng lực đã gán</div>
            {empCompetencies.length > 0 ? (
              <div className="tag-list">
                {empCompetencies.map((c) => (
                  <span key={c.id || c.competency_id} className="tag">
                    {c.name || c.competency?.name || 'N/A'} ({c.point_value || c.competency?.point_value || '?'} điểm)
                  </span>
                ))}
              </div>
            ) : (
              <p className="text-muted" style={{ fontSize: 13 }}>Chưa có năng lực nào được gán</p>
            )}
          </div>
        ) : (
          <p className="text-muted">Không thể tải thông tin</p>
        )}
      </Modal>

      {/* Salary Modal */}
      <Modal isOpen={salaryOpen} onClose={() => setSalaryOpen(false)} title={`Xem lương - ${selected?.full_name || ''}`} size="md"
        footer={<button className="btn btn-secondary" onClick={() => setSalaryOpen(false)}>Đóng</button>}
      >
        <div className="form-group" style={{ marginBottom: 16 }}>
          <label className="form-label">Kỳ lương</label>
          <div className="inline-edit">
            <input
              type="month"
              className="form-control"
              value={salaryPeriod}
              onChange={e => setSalaryPeriod(e.target.value)}
              style={{ width: 160 }}
            />
            <button className="btn btn-primary btn-sm" onClick={loadSalary} disabled={salaryLoading}>
              {salaryLoading ? <span className="spinner spinner-sm" /> : 'Tính toán'}
            </button>
          </div>
        </div>
        {salaryLoading ? (
          <div className="loading-center"><span className="spinner" /> Đang tính lương...</div>
        ) : salaryData ? (
          <div>
            <div className="grid-2" style={{ marginBottom: 16 }}>
              <div className="salary-card">
                <div className="salary-label">Lương P1 (Cơ bản)</div>
                <div className="salary-value accent">
                  {(salaryData.p1 || salaryData.salary_p1 || 0).toLocaleString('vi-VN')}đ
                </div>
              </div>
              <div className="salary-card">
                <div className="salary-label">Lương P2 (Phụ cấp)</div>
                <div className="salary-value accent">
                  {(salaryData.p2 || salaryData.salary_p2 || 0).toLocaleString('vi-VN')}đ
                </div>
              </div>
            </div>
            <div className="salary-card highlight">
              <div className="salary-label">Tổng lương</div>
              <div className="salary-value" style={{ fontSize: 30 }}>
                {(salaryData.total || salaryData.total_salary || (salaryData.p1 || 0) + (salaryData.p2 || 0) || 0).toLocaleString('vi-VN')}đ
              </div>
            </div>
            {salaryData.breakdown && (
              <div style={{ marginTop: 16 }}>
                <div className="section-title">Chi tiết</div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                  {Object.entries(salaryData.breakdown).map(([k, v]) => (
                    <div key={k} style={{ display: 'flex', justifyContent: 'space-between', padding: '8px 12px', background: 'var(--bg-input)', borderRadius: 'var(--radius-sm)' }}>
                      <span className="text-secondary">{k}</span>
                      <span className="font-bold">{typeof v === 'number' ? v.toLocaleString('vi-VN') + 'đ' : v}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        ) : (
          <div className="empty-state">
            <div className="empty-state-title">Chọn kỳ lương và bấm Tính toán</div>
          </div>
        )}
      </Modal>

      {/* Competencies Assignment Modal */}
      <Modal 
        isOpen={competenciesOpen} 
        onClose={() => setCompetenciesOpen(false)} 
        title={`Gán năng lực P2 — ${selected?.full_name || ''}`} 
        size="lg"
        footer={<button className="btn btn-secondary" onClick={() => setCompetenciesOpen(false)}>Đóng</button>}
      >
        <div className="form-group" style={{ marginBottom: 20 }}>
          <label className="form-label">Thêm năng lực mới cho nhân viên này</label>
          <div style={{ display: 'flex', gap: 10 }}>
            <select 
              className="form-control" 
              value={selectedCompetencyToAssign} 
              onChange={e => setSelectedCompetencyToAssign(e.target.value)}
              style={{ flex: 1 }}
            >
              <option value="">-- Chọn năng lực --</option>
              {availableCompetencies.map(c => (
                <option key={c.id} value={c.id}>
                  {c.code} — {c.name} ({c.point_value} điểm)
                </option>
              ))}
            </select>
            <button className="btn btn-primary" onClick={handleAssignCompetency} disabled={assigning || !selectedCompetencyToAssign}>
              {assigning ? 'Đang gán...' : 'Gán'}
            </button>
          </div>
        </div>

        <div className="section-title" style={{ fontSize: '14px', fontWeight: 600, marginBottom: 12 }}>
          Năng lực đang áp dụng ({empCompetencies.length})
        </div>

        {competenciesLoading ? (
          <div className="loading-center"><span className="spinner" /> Đang tải...</div>
        ) : empCompetencies.length > 0 ? (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Mã</th>
                  <th>Tên năng lực</th>
                  <th>Điểm số</th>
                  <th style={{ width: 80 }}>Thao tác</th>
                </tr>
              </thead>
              <tbody>
                {empCompetencies.map((c) => (
                  <tr key={c.competency_id || c.id}>
                    <td><Badge color="purple">{c.competency_code || c.code || 'N/A'}</Badge></td>
                    <td className="font-bold">{c.competency_name || c.name || 'N/A'}</td>
                    <td style={{ color: 'var(--purple)', fontWeight: 'bold' }}>
                      {(c.point_value || 0)} điểm
                    </td>
                    <td>
                      <button 
                        className="btn btn-danger btn-sm" 
                        onClick={() => handleRemoveCompetency(c.competency_id || c.id)}
                      >
                        Gỡ
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="text-muted" style={{ fontSize: 13, textAlign: 'center', padding: '20px 0' }}>
            Nhân viên này chưa được gán năng lực nào từ từ điển.
          </p>
        )}
      </Modal>

      {/* Delete Confirm */}
      <ConfirmDialog
        isOpen={deleteOpen}
        onClose={() => setDeleteOpen(false)}
        onConfirm={handleDelete}
        loading={deleting}
        title="Xóa nhân viên"
        message={`Bạn có chắc muốn xóa nhân viên "${selected?.full_name}"? Hành động này không thể hoàn tác.`}
        confirmText="Xóa nhân viên"
      />
    </div>
  );
}

function EmployeeForm({ title, onSubmit, isOpen, onClose, form, handleFormChange, saving, positions, departments }) {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={title}
      size="lg"
      footer={
        <>
          <button className="btn btn-secondary" onClick={onClose}>
            Hủy
          </button>
          <button className="btn btn-primary" onClick={onSubmit} disabled={saving}>
            {saving ? <><span className="spinner spinner-sm" /> Đang lưu...</> : 'Lưu'}
          </button>
        </>
      }
    >
      <div className="form-row">
        <div className="form-group">
          <label className="form-label">Mã nhân viên</label>
          <input className="form-control" name="code" value={form.code} onChange={handleFormChange} placeholder="NV001" />
        </div>
        <div className="form-group">
          <label className="form-label">Họ và tên <span className="required">*</span></label>
          <input className="form-control" name="full_name" value={form.full_name} onChange={handleFormChange} placeholder="Nguyễn Văn A" />
        </div>
      </div>
      <div className="form-row">
        <div className="form-group">
          <label className="form-label">Số điện thoại</label>
          <input className="form-control" name="phone" value={form.phone} onChange={handleFormChange} placeholder="0912345678" />
        </div>
        <div className="form-group">
          <label className="form-label">Giới tính</label>
          <select className="form-control" name="gender" value={form.gender} onChange={handleFormChange}>
            <option value="">-- Chọn giới tính --</option>
            <option value="male">Nam</option>
            <option value="female">Nữ</option>
          </select>
        </div>
      </div>
      <div className="form-row">
        <div className="form-group">
          <label className="form-label">Ngày sinh</label>
          <input className="form-control" type="date" name="date_of_birth" value={form.date_of_birth} onChange={handleFormChange} />
        </div>
        <div className="form-group">
          <label className="form-label">Vị trí công việc</label>
          <select className="form-control" name="position_id" value={form.position_id} onChange={handleFormChange}>
            <option value="">-- Chọn vị trí --</option>
            {positions.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
          </select>
        </div>
      </div>
      <div className="form-group">
        <label className="form-label">Phòng ban</label>
        <select className="form-control" name="department_id" value={form.department_id} onChange={handleFormChange}>
          <option value="">-- Chọn phòng ban --</option>
          {departments.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
        </select>
      </div>
    </Modal>
  );
}

export default EmployeesPage;
