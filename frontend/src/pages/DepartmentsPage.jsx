import { useState, useEffect, useCallback } from 'react';
import { departmentsAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { ConfirmDialog } from '../components/ConfirmDialog.jsx';
import { Pagination } from '../components/Pagination.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';
import { usePermissions } from '../contexts/PermissionsContext.jsx';

const PAGE_SIZE = 20;

const emptyForm = {
  code: '',
  name: '',
  description: '',
};

export function DepartmentsPage() {
  const toast = useToast();
  const { can } = usePermissions();
  const [departments, setDepartments] = useState([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(false);

  // Modals
  const [createOpen, setCreateOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);

  const [form, setForm] = useState(emptyForm);
  const [selected, setSelected] = useState(null);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => {
      setPage(1);
      loadDepartments(1, search);
    }, 350);
    return () => clearTimeout(timer);
  }, [search]);

  useEffect(() => {
    loadDepartments(page, search);
  }, [page]);

  const loadDepartments = useCallback(async (p, s) => {
    setLoading(true);
    try {
      const params = { page_number: p, page_size: PAGE_SIZE };
      if (s) params.search = s;
      const res = await departmentsAPI.list(params);
      const d = res.data?.data;
      setDepartments(d?.items || []);
      setTotal(d?.total_items || 0);
    } catch {
      toast.error('Lỗi', 'Không thể tải danh sách phòng ban');
    } finally {
      setLoading(false);
    }
  }, [toast]);

  const openCreate = () => {
    setForm(emptyForm);
    setCreateOpen(true);
  };

  const openEdit = (dept) => {
    setSelected(dept);
    setForm({
      code: dept.code || '',
      name: dept.name || '',
      description: dept.description || '',
    });
    setEditOpen(true);
  };

  const openDelete = (dept) => {
    setSelected(dept);
    setDeleteOpen(true);
  };

  const handleCreate = async () => {
    if (!form.code || !form.name) {
      toast.error('Thiếu thông tin', 'Vui lòng nhập mã và tên phòng ban');
      return;
    }
    setSaving(true);
    try {
      await departmentsAPI.create(form);
      toast.success('Thành công', 'Đã thêm phòng ban mới');
      setCreateOpen(false);
      loadDepartments(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể tạo phòng ban');
    } finally {
      setSaving(false);
    }
  };

  const handleEdit = async () => {
    if (!form.name) {
      toast.error('Thiếu thông tin', 'Vui lòng nhập tên phòng ban');
      return;
    }
    setSaving(true);
    try {
      await departmentsAPI.update(selected.id, form);
      toast.success('Thành công', 'Đã cập nhật phòng ban');
      setEditOpen(false);
      loadDepartments(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể cập nhật');
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    setDeleting(true);
    try {
      await departmentsAPI.delete(selected.id);
      toast.success('Đã xóa', 'Phòng ban đã được xóa');
      setDeleteOpen(false);
      loadDepartments(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể xóa phòng ban');
    } finally {
      setDeleting(false);
    }
  };

  const handleFormChange = (e) => {
    setForm((p) => ({ ...p, [e.target.name]: e.target.value }));
  };

  const columns = [
    {
      key: 'code',
      title: 'Mã phòng ban',
      render: (val) => <Badge color="cyan">{val}</Badge>,
    },
    {
      key: 'name',
      title: 'Tên phòng ban',
      render: (val) => <span className="font-bold">{val}</span>,
    },
    {
      key: 'description',
      title: 'Mô tả',
      render: (val) => val || <span className="text-muted">—</span>,
    },
    {
      key: 'id',
      title: 'Thao tác',
      style: { width: 140 },
      render: (_, row) => (
        <div className="td-actions">
          {can('department::edit') && (
            <button className="btn btn-secondary btn-sm" onClick={() => openEdit(row)}>
              Sửa
            </button>
          )}
          {can('department::delete') && (
            <button className="btn btn-danger btn-sm" onClick={() => openDelete(row)}>
              Xóa
            </button>
          )}
        </div>
      )
    }
  ];



  return (
    <div>
      <div className="page-header">
        <div>
          <h2 className="page-title">Phòng ban</h2>
          <p className="page-subtitle">Quản lý cơ cấu phòng ban, đơn vị của công ty ({total} phòng ban)</p>
        </div>
        {can('department::edit') && (
          <button className="btn btn-primary" onClick={openCreate}>
            Thêm phòng ban
          </button>
        )}
      </div>

      <div className="table-container">
        <div className="table-toolbar">
          <div className="search-bar">
            <input
              className="form-control"
              placeholder="Tìm kiếm phòng ban..."
              value={search}
              onChange={e => setSearch(e.target.value)}
            />
          </div>
          <div className="text-muted" style={{ fontSize: 13 }}>
            {total} kết quả
          </div>
        </div>
        <Table columns={columns} data={departments} loading={loading} emptyMessage="Chưa có phòng ban nào" />
        <Pagination page={page} totalItems={total} pageSize={PAGE_SIZE} onPageChange={setPage} />
      </div>

      {/* Create Modal */}
      <DepartmentForm title="Thêm phòng ban" onSubmit={handleCreate} isOpen={createOpen} onClose={() => setCreateOpen(false)} form={form} handleFormChange={handleFormChange} saving={saving} />

      {/* Edit Modal */}
      <DepartmentForm title="Sửa thông tin phòng ban" onSubmit={handleEdit} isOpen={editOpen} onClose={() => setEditOpen(false)} form={form} handleFormChange={handleFormChange} saving={saving} />

      {/* Delete Confirm */}
      <ConfirmDialog
        isOpen={deleteOpen}
        onClose={() => setDeleteOpen(false)}
        onConfirm={handleDelete}
        loading={deleting}
        title="Xóa phòng ban"
        message={`Bạn có chắc muốn xóa phòng ban "${selected?.name}"?`}
        confirmText="Xóa phòng ban"
      />
    </div>
  );
}

function DepartmentForm({ title, onSubmit, isOpen, onClose, form, handleFormChange, saving }) {
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
        <label className="form-label">Mã phòng ban <span className="required">*</span></label>
        <input 
          className="form-control" 
          name="code" 
          value={form.code} 
          onChange={handleFormChange} 
          placeholder="VD: PHONG_KHOA_HOC" 
          disabled={title.includes('Sửa')}
        />
      </div>
      <div className="form-group">
        <label className="form-label">Tên phòng ban <span className="required">*</span></label>
        <input 
          className="form-control" 
          name="name" 
          value={form.name} 
          onChange={handleFormChange} 
          placeholder="VD: Phòng Khoa học Công nghệ" 
        />
      </div>
      <div className="form-group">
        <label className="form-label">Mô tả</label>
        <textarea 
          className="form-control" 
          name="description" 
          value={form.description} 
          onChange={handleFormChange} 
          placeholder="Mô tả về phòng ban..." 
          style={{ minHeight: 80 }}
        />
      </div>
    </Modal>
  );
}

export default DepartmentsPage;
