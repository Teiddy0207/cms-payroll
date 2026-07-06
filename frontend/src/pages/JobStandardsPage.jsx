import { useState, useEffect, useCallback } from 'react';
import { jobStandardsAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { ConfirmDialog } from '../components/ConfirmDialog.jsx';
import { Pagination } from '../components/Pagination.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';

const PAGE_SIZE = 20;

const emptyForm = {
  standard_code: '',
  name: '',
  description: '',
  allowance_value: 0,
};

export function JobStandardsPage() {
  const toast = useToast();
  const [standards, setStandards] = useState([]);
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
      loadStandards(1, search);
    }, 350);
    return () => clearTimeout(timer);
  }, [search]);

  useEffect(() => {
    loadStandards(page, search);
  }, [page]);

  const loadStandards = useCallback(async (p, s) => {
    setLoading(true);
    try {
      const params = { page_number: p, page_size: PAGE_SIZE };
      if (s) params.search = s;
      const res = await jobStandardsAPI.list(params);
      const d = res.data?.data;
      setStandards(d?.items || []);
      setTotal(d?.total_items || 0);
    } catch {
      toast.error('Lỗi', 'Không thể tải danh sách tiêu chuẩn');
    } finally {
      setLoading(false);
    }
  }, [toast]);

  const openCreate = () => {
    setForm(emptyForm);
    setCreateOpen(true);
  };

  const openEdit = (std) => {
    setSelected(std);
    setForm({
      standard_code: std.standard_code || '',
      name: std.name || '',
      description: std.description || '',
      allowance_value: std.allowance_value || 0,
    });
    setEditOpen(true);
  };

  const openDelete = (std) => {
    setSelected(std);
    setDeleteOpen(true);
  };

  const handleCreate = async () => {
    if (!form.standard_code || !form.name || form.allowance_value <= 0) {
      toast.error('Thiếu thông tin', 'Vui lòng nhập mã, tên và trị số phụ cấp hợp lệ');
      return;
    }
    setSaving(true);
    try {
      await jobStandardsAPI.create({
        ...form,
        allowance_value: parseFloat(form.allowance_value),
      });
      toast.success('Thành công', 'Đã thêm tiêu chuẩn mới');
      setCreateOpen(false);
      loadStandards(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể tạo tiêu chuẩn');
    } finally {
      setSaving(false);
    }
  };

  const handleEdit = async () => {
    if (!form.name || form.allowance_value <= 0) {
      toast.error('Thiếu thông tin', 'Vui lòng nhập tên và trị số phụ cấp hợp lệ');
      return;
    }
    setSaving(true);
    try {
      await jobStandardsAPI.update(selected.id, {
        ...form,
        allowance_value: parseFloat(form.allowance_value),
      });
      toast.success('Thành công', 'Đã cập nhật tiêu chuẩn');
      setEditOpen(false);
      loadStandards(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể cập nhật');
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    setDeleting(true);
    try {
      await jobStandardsAPI.delete(selected.id);
      toast.success('Đã xóa', 'Tiêu chuẩn đã được xóa');
      setDeleteOpen(false);
      loadStandards(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể xóa tiêu chuẩn');
    } finally {
      setDeleting(false);
    }
  };

  const handleFormChange = (e) => {
    setForm((p) => ({ ...p, [e.target.name]: e.target.value }));
  };

  const columns = [
    {
      key: 'standard_code',
      title: 'Mã tiêu chuẩn',
      render: (val) => <Badge color="green">{val}</Badge>,
    },
    {
      key: 'name',
      title: 'Tên tiêu chuẩn',
      render: (val) => <span className="font-bold">{val}</span>,
    },
    {
      key: 'allowance_value',
      title: 'Điểm phụ cấp (allowance_value)',
      render: (val) => <span className="text-success font-bold">{val?.toLocaleString('vi-VN')}</span>,
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



  return (
    <div>
      <div className="page-header">
        <div>
          <h2 className="page-title">Tiêu chuẩn công việc (P1)</h2>
          <p className="page-subtitle">Quản lý danh mục các tiêu chuẩn áp dụng để tính lương cứng P1 ({total} tiêu chuẩn)</p>
        </div>
        <button className="btn btn-primary" onClick={openCreate}>
          Thêm tiêu chuẩn
        </button>
      </div>

      <div className="table-container">
        <div className="table-toolbar">
          <div className="search-bar">
            <input
              className="form-control"
              placeholder="Tìm kiếm tiêu chuẩn..."
              value={search}
              onChange={e => setSearch(e.target.value)}
            />
          </div>
          <div className="text-muted" style={{ fontSize: 13 }}>
            {total} kết quả
          </div>
        </div>
        <Table columns={columns} data={standards} loading={loading} emptyMessage="Chưa có tiêu chuẩn công việc nào" />
        <Pagination page={page} totalItems={total} pageSize={PAGE_SIZE} onPageChange={setPage} />
      </div>

      {/* Create Modal */}
      <StandardForm title="Thêm tiêu chuẩn" onSubmit={handleCreate} isOpen={createOpen} onClose={() => setCreateOpen(false)} form={form} handleFormChange={handleFormChange} saving={saving} />

      {/* Edit Modal */}
      <StandardForm title="Sửa thông tin tiêu chuẩn" onSubmit={handleEdit} isOpen={editOpen} onClose={() => setEditOpen(false)} form={form} handleFormChange={handleFormChange} saving={saving} />

      {/* Delete Confirm */}
      <ConfirmDialog
        isOpen={deleteOpen}
        onClose={() => setDeleteOpen(false)}
        onConfirm={handleDelete}
        loading={deleting}
        title="Xóa tiêu chuẩn"
        message={`Bạn có chắc muốn xóa tiêu chuẩn "${selected?.name}"?`}
        confirmText="Xóa tiêu chuẩn"
      />
    </div>
  );
}

function StandardForm({ title, onSubmit, isOpen, onClose, form, handleFormChange, saving }) {
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
      <div className="form-row">
        <div className="form-group">
          <label className="form-label">Mã tiêu chuẩn <span className="required">*</span></label>
          <input 
            className="form-control" 
            name="standard_code" 
            value={form.standard_code} 
            onChange={handleFormChange} 
            placeholder="VD: STD_IT_01" 
            disabled={title.includes('Sửa')}
          />
        </div>
        <div className="form-group">
          <label className="form-label">Điểm phụ cấp <span className="required">*</span></label>
          <input 
            className="form-control" 
            type="number" 
            name="allowance_value" 
            value={form.allowance_value} 
            onChange={handleFormChange} 
            placeholder="VD: 5" 
          />
        </div>
      </div>
      <div className="form-group">
        <label className="form-label">Tên tiêu chuẩn <span className="required">*</span></label>
        <input 
          className="form-control" 
          name="name" 
          value={form.name} 
          onChange={handleFormChange} 
          placeholder="VD: Tiêu chuẩn công nghệ thông tin cấp 1" 
        />
      </div>
      <div className="form-group">
        <label className="form-label">Mô tả</label>
        <textarea 
          className="form-control" 
          name="description" 
          value={form.description} 
          onChange={handleFormChange} 
          placeholder="Mô tả tiêu chuẩn..." 
          style={{ minHeight: 80 }}
        />
      </div>
    </Modal>
  );
}

export default JobStandardsPage;
