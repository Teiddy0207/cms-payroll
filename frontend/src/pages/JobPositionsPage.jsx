import { useState, useEffect, useCallback } from 'react';
import { jobPositionsAPI, jobStandardsAPI, departmentsAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { ConfirmDialog } from '../components/ConfirmDialog.jsx';
import { Pagination } from '../components/Pagination.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';

const PAGE_SIZE = 20;

const emptyForm = {
  code: '',
  name: '',
  description: '',
  department_id: '',
};

export function JobPositionsPage() {
  const toast = useToast();
  const [positions, setPositions] = useState([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(false);

  const [departments, setDepartments] = useState([]);
  const [allStandards, setAllStandards] = useState([]);

  // Modals
  const [createOpen, setCreateOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [standardsOpen, setStandardsOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);

  const [form, setForm] = useState(emptyForm);
  const [selected, setSelected] = useState(null);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);

  // Standards Management for selected Position
  const [assignedStandards, setAssignedStandards] = useState([]);
  const [standardsLoading, setStandardsLoading] = useState(false);
  const [selectedStandardToAssign, setSelectedStandardToAssign] = useState('');
  const [assigning, setAssigning] = useState(false);

  useEffect(() => {
    loadMeta();
  }, []);

  useEffect(() => {
    const timer = setTimeout(() => {
      setPage(1);
      loadPositions(1, search);
    }, 350);
    return () => clearTimeout(timer);
  }, [search]);

  useEffect(() => {
    loadPositions(page, search);
  }, [page]);

  const loadMeta = async () => {
    try {
      const [deptRes, stdRes] = await Promise.allSettled([
        departmentsAPI.list({ page_size: 200 }),
        jobStandardsAPI.list({ page_size: 200 }),
      ]);
      if (deptRes.status === 'fulfilled') setDepartments(deptRes.value.data?.data?.items || []);
      if (stdRes.status === 'fulfilled') setAllStandards(stdRes.value.data?.data?.items || []);
    } catch { /**/ }
  };

  const loadPositions = useCallback(async (p, s) => {
    setLoading(true);
    try {
      const params = { page_number: p, page_size: PAGE_SIZE };
      if (s) params.search = s;
      const res = await jobPositionsAPI.list(params);
      const d = res.data?.data;
      setPositions(d?.items || []);
      setTotal(d?.total_items || 0);
    } catch {
      toast.error('Lỗi', 'Không thể tải danh sách vị trí công việc');
    } finally {
      setLoading(false);
    }
  }, [toast]);

  const openCreate = () => {
    setForm(emptyForm);
    setCreateOpen(true);
  };

  const openEdit = (pos) => {
    setSelected(pos);
    setForm({
      code: pos.code || '',
      name: pos.name || '',
      description: pos.description || '',
      department_id: pos.department_id || '',
    });
    setEditOpen(true);
  };

  const openStandards = async (pos) => {
    setSelected(pos);
    setStandardsOpen(true);
    setStandardsLoading(true);
    setAssignedStandards([]);
    setSelectedStandardToAssign('');
    try {
      const res = await jobPositionsAPI.getStandards(pos.id);
      setAssignedStandards(res.data?.data || []);
    } catch {
      toast.error('Lỗi', 'Không thể tải tiêu chuẩn của vị trí');
    } finally {
      setStandardsLoading(false);
    }
  };

  const openDelete = (pos) => {
    setSelected(pos);
    setDeleteOpen(true);
  };

  const handleCreate = async () => {
    if (!form.code || !form.name) {
      toast.error('Thiếu thông tin', 'Vui lòng nhập mã và tên vị trí');
      return;
    }
    setSaving(true);
    try {
      await jobPositionsAPI.create(form);
      toast.success('Thành công', 'Đã tạo vị trí công việc mới');
      setCreateOpen(false);
      loadPositions(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể tạo vị trí');
    } finally {
      setSaving(false);
    }
  };

  const handleEdit = async () => {
    if (!form.name) {
      toast.error('Thiếu thông tin', 'Vui lòng nhập tên vị trí');
      return;
    }
    setSaving(true);
    try {
      await jobPositionsAPI.update(selected.id, form);
      toast.success('Thành công', 'Đã cập nhật vị trí công việc');
      setEditOpen(false);
      loadPositions(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể cập nhật');
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    setDeleting(true);
    try {
      await jobPositionsAPI.delete(selected.id);
      toast.success('Đã xóa', 'Vị trí công việc đã được xóa');
      setDeleteOpen(false);
      loadPositions(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể xóa vị trí');
    } finally {
      setDeleting(false);
    }
  };

  const handleAssignStandard = async () => {
    if (!selectedStandardToAssign) {
      toast.error('Cảnh báo', 'Vui lòng chọn một tiêu chuẩn');
      return;
    }
    setAssigning(true);
    try {
      await jobPositionsAPI.assignStandard(selected.id, { job_standard_id: selectedStandardToAssign });
      toast.success('Thành công', 'Đã gán tiêu chuẩn cho vị trí');
      setSelectedStandardToAssign('');
      
      // Reload assigned standards
      const res = await jobPositionsAPI.getStandards(selected.id);
      setAssignedStandards(res.data?.data || []);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể gán tiêu chuẩn');
    } finally {
      setAssigning(false);
    }
  };

  const handleRemoveStandard = async (standardId) => {
    try {
      await jobPositionsAPI.removeStandard(selected.id, standardId);
      toast.success('Đã xóa', 'Đã gỡ tiêu chuẩn khỏi vị trí');
      
      // Reload assigned standards
      const res = await jobPositionsAPI.getStandards(selected.id);
      setAssignedStandards(res.data?.data || []);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể gỡ tiêu chuẩn');
    }
  };

  const handleFormChange = (e) => {
    setForm((p) => ({ ...p, [e.target.name]: e.target.value }));
  };

  const columns = [
    {
      key: 'code',
      title: 'Mã vị trí',
      render: (val) => <Badge color="blue">{val}</Badge>,
    },
    {
      key: 'name',
      title: 'Tên vị trí',
      render: (val) => <span className="font-bold">{val}</span>,
    },
    {
      key: 'description',
      title: 'Mô tả',
      render: (val) => val || <span className="text-muted">—</span>,
    },
    {
      key: 'department_id',
      title: 'Phòng ban',
      render: (val) => {
        const dept = departments.find(d => d.id === val);
        return dept ? <Badge color="gray">{dept.name}</Badge> : <span className="text-muted">—</span>;
      }
    },
    {
      key: 'id',
      title: 'Thao tác',
      style: { width: 220 },
      render: (_, row) => (
        <div className="td-actions">
          <button className="btn btn-secondary btn-sm" onClick={() => openStandards(row)}>
            Tiêu chuẩn
          </button>
          <button className="btn btn-secondary btn-sm" onClick={() => openEdit(row)}>
            Sửa
          </button>
          <button className="btn btn-danger btn-sm" onClick={() => openDelete(row)}>
            Xóa
          </button>
        </div>
      )
    }
  ];

  // Available standards to assign (those not already assigned)
  const availableStandards = allStandards.filter(
    std => !assignedStandards.some(assigned => assigned.job_standard_id === std.id)
  );



  return (
    <div>
      <div className="page-header">
        <div>
          <h2 className="page-title">Vị trí công việc</h2>
          <p className="page-subtitle">Quản lý các chức danh & tiêu chuẩn lương cứng P1 ({total} vị trí)</p>
        </div>
        <button className="btn btn-primary" onClick={openCreate}>
          + Thêm vị trí
        </button>
      </div>

      <div className="table-container">
        <div className="table-toolbar">
          <div className="search-bar">
            <input
              className="form-control"
              placeholder="Tìm kiếm vị trí..."
              value={search}
              onChange={e => setSearch(e.target.value)}
            />
          </div>
          <div className="text-muted" style={{ fontSize: 13 }}>
            {total} kết quả
          </div>
        </div>
        <Table columns={columns} data={positions} loading={loading} emptyMessage="Chưa có vị trí công việc nào" />
        <Pagination page={page} totalItems={total} pageSize={PAGE_SIZE} onPageChange={setPage} />
      </div>

      {/* Create Modal */}
      <PositionForm title="Thêm vị trí" onSubmit={handleCreate} isOpen={createOpen} onClose={() => setCreateOpen(false)} form={form} handleFormChange={handleFormChange} saving={saving} departments={departments} />

      {/* Edit Modal */}
      <PositionForm title="Sửa thông tin vị trí" onSubmit={handleEdit} isOpen={editOpen} onClose={() => setEditOpen(false)} form={form} handleFormChange={handleFormChange} saving={saving} departments={departments} />

      {/* Standards Assignment Modal */}
      <Modal 
        isOpen={standardsOpen} 
        onClose={() => setStandardsOpen(false)} 
        title={`Gán tiêu chuẩn P1 — ${selected?.name || ''}`} 
        size="lg"
        footer={<button className="btn btn-secondary" onClick={() => setStandardsOpen(false)}>Đóng</button>}
      >
        <div className="form-group" style={{ marginBottom: 20 }}>
          <label className="form-label">Thêm tiêu chuẩn mới cho vị trí này</label>
          <div style={{ display: 'flex', gap: 10 }}>
            <select 
              className="form-control" 
              value={selectedStandardToAssign} 
              onChange={e => setSelectedStandardToAssign(e.target.value)}
              style={{ flex: 1 }}
            >
              <option value="">-- Chọn tiêu chuẩn công việc --</option>
              {availableStandards.map(s => (
                <option key={s.id} value={s.id}>
                  {s.standard_code} — {s.name} ({s.allowance_value.toLocaleString('vi-VN')}đ)
                </option>
              ))}
            </select>
            <button className="btn btn-primary" onClick={handleAssignStandard} disabled={assigning || !selectedStandardToAssign}>
              {assigning ? 'Đang gán...' : 'Gán'}
            </button>
          </div>
        </div>

        <div className="section-title" style={{ fontSize: '14px', fontWeight: 600, marginBottom: 12 }}>
          Tiêu chuẩn đang áp dụng ({assignedStandards.length})
        </div>

        {standardsLoading ? (
          <div className="loading-center"><span className="spinner" /> Đang tải...</div>
        ) : assignedStandards.length > 0 ? (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Mã</th>
                  <th>Tên tiêu chuẩn</th>
                  <th>Phụ cấp</th>
                  <th style={{ width: 80 }}>Thao tác</th>
                </tr>
              </thead>
              <tbody>
                {assignedStandards.map((s) => (
                  <tr key={s.job_standard_id || s.id}>
                    <td><Badge color="green">{s.standard_code}</Badge></td>
                    <td className="font-bold">{s.standard_name}</td>
                    <td className="text-success font-bold">
                      {s.allowance_value?.toLocaleString('vi-VN')}đ
                    </td>
                    <td>
                      <button 
                        className="btn btn-danger btn-sm" 
                        onClick={() => handleRemoveStandard(s.job_standard_id || s.id)}
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
            Vị trí này chưa được gán tiêu chuẩn lương cứng nào.
          </p>
        )}
      </Modal>

      {/* Delete Confirm */}
      <ConfirmDialog
        isOpen={deleteOpen}
        onClose={() => setDeleteOpen(false)}
        onConfirm={handleDelete}
        loading={deleting}
        title="Xóa vị trí công việc"
        message={`Bạn có chắc muốn xóa vị trí "${selected?.name}"? Các nhân sự đang ở vị trí này có thể bị ảnh hưởng.`}
        confirmText="Xóa vị trí"
      />
    </div>
  );
}

function PositionForm({ title, onSubmit, isOpen, onClose, form, handleFormChange, saving, departments }) {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={title}
      size="md"
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
      <div className="form-group">
        <label className="form-label">Mã vị trí <span className="required">*</span></label>
        <input 
          className="form-control" 
          name="code" 
          value={form.code} 
          onChange={handleFormChange} 
          placeholder="VD: DEV_BE" 
          disabled={title.includes('Sửa')}
        />
      </div>
      <div className="form-group">
        <label className="form-label">Tên vị trí <span className="required">*</span></label>
        <input 
          className="form-control" 
          name="name" 
          value={form.name} 
          onChange={handleFormChange} 
          placeholder="VD: Lập trình viên Backend" 
        />
      </div>
      <div className="form-group">
        <label className="form-label">Phòng ban</label>
        <select className="form-control" name="department_id" value={form.department_id} onChange={handleFormChange}>
          <option value="">-- Chọn phòng ban --</option>
          {departments.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
        </select>
      </div>
      <div className="form-group">
        <label className="form-label">Mô tả</label>
        <textarea 
          className="form-control" 
          name="description" 
          value={form.description} 
          onChange={handleFormChange} 
          placeholder="Mô tả công việc..." 
          style={{ minHeight: 80 }}
        />
      </div>
    </Modal>
  );
}

export default JobPositionsPage;
