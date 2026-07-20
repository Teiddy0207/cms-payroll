import { useState, useEffect, useCallback } from 'react';
import { competenciesAPI } from '../api/client.js';
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
  point_value: 0,
};

export function CompetenciesPage() {
  const toast = useToast();
  const [competencies, setCompetencies] = useState([]);
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
      loadCompetencies(1, search);
    }, 350);
    return () => clearTimeout(timer);
  }, [search]);

  useEffect(() => {
    loadCompetencies(page, search);
  }, [page]);

  const loadCompetencies = useCallback(async (p, s) => {
    setLoading(true);
    try {
      const params = { page_number: p, page_size: PAGE_SIZE };
      if (s) params.search = s;
      const res = await competenciesAPI.list(params);
      const d = res.data?.data;
      setCompetencies(d?.items || []);
      setTotal(d?.total_items || 0);
    } catch {
      toast.error('Lỗi', 'Không thể tải danh sách năng lực');
    } finally {
      setLoading(false);
    }
  }, [toast]);

  const openCreate = () => {
    setForm(emptyForm);
    setCreateOpen(true);
  };

  const openEdit = (comp) => {
    setSelected(comp);
    setForm({
      code: comp.code || '',
      name: comp.name || '',
      description: comp.description || '',
      point_value: comp.point_value || 0,
    });
    setEditOpen(true);
  };

  const openDelete = (comp) => {
    setSelected(comp);
    setDeleteOpen(true);
  };

  const handleCreate = async () => {
    if (!form.code || !form.name || form.point_value <= 0) {
      toast.error('Thiếu thông tin', 'Vui lòng nhập mã, tên và mức phụ cấp hợp lệ');
      return;
    }
    setSaving(true);
    try {
      await competenciesAPI.create({
        ...form,
        point_value: parseInt(form.point_value),
      });
      toast.success('Thành công', 'Đã thêm năng lực mới vào từ điển');
      setCreateOpen(false);
      loadCompetencies(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể tạo năng lực');
    } finally {
      setSaving(false);
    }
  };

  const handleEdit = async () => {
    if (!form.name || form.point_value <= 0) {
      toast.error('Thiếu thông tin', 'Vui lòng nhập tên và mức phụ cấp hợp lệ');
      return;
    }
    setSaving(true);
    try {
      await competenciesAPI.update(selected.id, {
        ...form,
        point_value: parseInt(form.point_value),
      });
      toast.success('Thành công', 'Đã cập nhật năng lực');
      setEditOpen(false);
      loadCompetencies(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể cập nhật');
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    setDeleting(true);
    try {
      await competenciesAPI.delete(selected.id);
      toast.success('Đã xóa', 'Năng lực đã được xóa khỏi từ điển');
      setDeleteOpen(false);
      loadCompetencies(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể xóa năng lực');
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
      title: 'Mã năng lực',
      render: (val) => <Badge color="purple">{val}</Badge>,
    },
    {
      key: 'name',
      title: 'Tên năng lực',
      render: (val) => <span className="font-bold">{val}</span>,
    },
    {
      key: 'point_value',
      title: 'Mức phụ cấp (VND)',
      render: (val) => <span style={{ color: 'var(--purple)', fontWeight: 'bold' }}>{val?.toLocaleString('vi-VN')}đ</span>,
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
          <h2 className="page-title">Từ điển năng lực (P2)</h2>
          <p className="page-subtitle">Quản lý từ điển năng lực dùng làm căn cứ tính lương năng lực P2 ({total} năng lực)</p>
        </div>
        <button className="btn btn-primary" onClick={openCreate}>
          Thêm năng lực
        </button>
      </div>

      <div className="table-container">
        <div className="table-toolbar">
          <div className="search-bar">
            <input
              className="form-control"
              placeholder="Tìm kiếm năng lực..."
              value={search}
              onChange={e => setSearch(e.target.value)}
            />
          </div>
          <div className="text-muted" style={{ fontSize: 13 }}>
            {total} kết quả
          </div>
        </div>
        <Table columns={columns} data={competencies} loading={loading} emptyMessage="Chưa có năng lực nào" />
        <Pagination page={page} totalItems={total} pageSize={PAGE_SIZE} onPageChange={setPage} />
      </div>

      {/* Create Modal */}
      <CompetencyForm title="Thêm năng lực" onSubmit={handleCreate} isOpen={createOpen} onClose={() => setCreateOpen(false)} form={form} handleFormChange={handleFormChange} saving={saving} />

      {/* Edit Modal */}
      <CompetencyForm title="Sửa thông tin năng lực" onSubmit={handleEdit} isOpen={editOpen} onClose={() => setEditOpen(false)} form={form} handleFormChange={handleFormChange} saving={saving} />

      {/* Delete Confirm */}
      <ConfirmDialog
        isOpen={deleteOpen}
        onClose={() => setDeleteOpen(false)}
        onConfirm={handleDelete}
        loading={deleting}
        title="Xóa năng lực"
        message={`Bạn có chắc muốn xóa năng lực "${selected?.name}"?`}
        confirmText="Xóa năng lực"
      />
    </div>
  );
}

function CompetencyForm({ title, onSubmit, isOpen, onClose, form, handleFormChange, saving }) {
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
          <label className="form-label">Mã năng lực <span className="required">*</span></label>
          <input 
            className="form-control" 
            name="code" 
            value={form.code} 
            onChange={handleFormChange} 
            placeholder="VD: ENG_IELTS_7" 
            disabled={title.includes('Sửa')}
          />
        </div>
        <div className="form-group">
          <label className="form-label">Mức phụ cấp (VND) <span className="required">*</span></label>
          <input 
            className="form-control" 
            type="number" 
            name="point_value" 
            value={form.point_value} 
            onChange={handleFormChange} 
            placeholder="VD: 500000" 
          />
        </div>
      </div>
      <div className="form-group">
        <label className="form-label">Tên năng lực <span className="required">*</span></label>
        <input 
          className="form-control" 
          name="name" 
          value={form.name} 
          onChange={handleFormChange} 
          placeholder="VD: Kỹ năng Tiếng Anh IELTS 7.0" 
        />
      </div>
      <div className="form-group">
        <label className="form-label">Mô tả</label>
        <textarea 
          className="form-control" 
          name="description" 
          value={form.description} 
          onChange={handleFormChange} 
          placeholder="Mô tả năng lực..." 
          style={{ minHeight: 80 }}
        />
      </div>
    </Modal>
  );
}

export default CompetenciesPage;
