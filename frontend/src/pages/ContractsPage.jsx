import { useState, useEffect, useCallback } from 'react';
import { contractsAPI, employeesAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { ConfirmDialog } from '../components/ConfirmDialog.jsx';
import { Pagination } from '../components/Pagination.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';
import { usePermissions } from '../contexts/PermissionsContext.jsx';

const PAGE_SIZE = 20;

const emptyForm = {
  employee_id: '',
  contract_code: '',
  position_base_rate: 0,
  start_date: '',
  end_date: '',
  status: 'ACTIVE',
};

function formatDate(str) {
  if (!str) return '—';
  try { return new Date(str).toLocaleDateString('vi-VN'); } catch { return str; }
}

export function ContractsPage() {
  const toast = useToast();
  const { can } = usePermissions();
  const [contracts, setContracts] = useState([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(false);

  const [employees, setEmployees] = useState([]);

  // Modals
  const [createOpen, setCreateOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);

  const [form, setForm] = useState(emptyForm);
  const [selected, setSelected] = useState(null);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    loadMeta();
  }, []);

  useEffect(() => {
    const timer = setTimeout(() => {
      setPage(1);
      loadContracts(1, search);
    }, 350);
    return () => clearTimeout(timer);
  }, [search]);

  useEffect(() => {
    loadContracts(page, search);
  }, [page]);

  const loadMeta = async () => {
    try {
      const res = await employeesAPI.list({ page_size: 200 });
      setEmployees(res.data?.data?.items || []);
    } catch { /**/ }
  };

  const loadContracts = useCallback(async (p, s) => {
    setLoading(true);
    try {
      const params = { page_number: p, page_size: PAGE_SIZE };
      if (s) params.search = s;
      const res = await contractsAPI.list(params);
      const d = res.data?.data;
      setContracts(d?.items || []);
      setTotal(d?.total_items || 0);
    } catch {
      toast.error('Lỗi', 'Không thể tải danh sách hợp đồng');
    } finally {
      setLoading(false);
    }
  }, [toast]);

  const openCreate = () => {
    setForm(emptyForm);
    setCreateOpen(true);
  };

  const openEdit = (contract) => {
    setSelected(contract);
    setForm({
      employee_id: contract.employee_id || '',
      contract_code: contract.contract_code || '',
      position_base_rate: contract.position_base_rate || 0,
      start_date: contract.start_date ? contract.start_date.substring(0, 10) : '',
      end_date: contract.end_date ? contract.end_date.substring(0, 10) : '',
      status: contract.status || 'ACTIVE',
    });
    setEditOpen(true);
  };

  const openDelete = (contract) => {
    setSelected(contract);
    setDeleteOpen(true);
  };

  const handleCreate = async () => {
    if (!form.employee_id || !form.contract_code || !form.start_date) {
      toast.error('Thiếu thông tin', 'Vui lòng điền nhân viên, số hợp đồng và ngày bắt đầu');
      return;
    }
    setSaving(true);
    try {
      const data = { 
        ...form,
        position_base_rate: Number(form.position_base_rate)
      };
      if (!data.end_date) {
        delete data.end_date;
      }
      await contractsAPI.create(data);
      toast.success('Thành công', 'Đã thêm hợp đồng lao động mới');
      setCreateOpen(false);
      loadContracts(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể tạo hợp đồng');
    } finally {
      setSaving(false);
    }
  };

  const handleEdit = async () => {
    if (!form.contract_code || !form.start_date) {
      toast.error('Thiếu thông tin', 'Vui lòng điền số hợp đồng và ngày bắt đầu');
      return;
    }
    setSaving(true);
    try {
      const data = { 
        ...form,
        position_base_rate: Number(form.position_base_rate)
      };
      if (!data.end_date) {
        data.end_date = null;
      }
      await contractsAPI.update(selected.id, data);
      toast.success('Thành công', 'Đã cập nhật hợp đồng');
      setEditOpen(false);
      loadContracts(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể cập nhật');
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    setDeleting(true);
    try {
      await contractsAPI.delete(selected.id);
      toast.success('Đã xóa', 'Hợp đồng lao động đã bị xóa');
      setDeleteOpen(false);
      loadContracts(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể xóa hợp đồng');
    } finally {
      setDeleting(false);
    }
  };

  const handleFormChange = (e) => {
    setForm((p) => ({ ...p, [e.target.name]: e.target.value }));
  };

  const columns = [
    {
      key: 'contract_code',
      title: 'Số Hợp đồng',
      render: (val) => <span className="font-bold">{val}</span>,
    },
    {
      key: 'employee_id',
      title: 'Nhân viên',
      render: (val, row) => {
        const emp = employees.find(e => e.id === val) || row.employee;
        return emp ? <span className="font-bold">{emp.full_name || emp.fullName || 'N/A'}</span> : <span className="text-muted">—</span>;
      }
    },
    {
      key: 'position_base_rate',
      title: 'Lương cơ bản (P1)',
      render: (val) => val ? val.toLocaleString('vi-VN') + ' đ' : '0 đ',
    },
    {
      key: 'status',
      title: 'Trạng thái',
      render: (val) => {
        const statuses = {
          ACTIVE: <Badge color="green">Hoạt động</Badge>,
          EXPIRED: <Badge color="gray">Hết hạn</Badge>,
          TERMINATED: <Badge color="red">Chấm dứt</Badge>
        };
        return statuses[val] || <Badge color="gray">{val}</Badge>;
      }
    },
    {
      key: 'start_date',
      title: 'Từ ngày',
      render: (val) => formatDate(val),
    },
    {
      key: 'end_date',
      title: 'Đến ngày',
      render: (val, row) => !val || row.status === 'indefinite' 
        ? <Badge color="green">Không thời hạn</Badge> 
        : formatDate(val),
    },
    {
      key: 'id',
      title: 'Thao tác',
      style: { width: 140 },
      render: (_, row) => (
        <div className="td-actions">
          {can('userProfile::edit') && (
            <button className="btn btn-secondary btn-sm" onClick={() => openEdit(row)}>
              Sửa
            </button>
          )}
          {can('userProfile::delete') && (
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
          <h2 className="page-title">Hợp đồng lao động</h2>
          <p className="page-subtitle">Quản lý các hợp đồng lao động, mức lương cơ bản (P1) của nhân sự ({total} hợp đồng)</p>
        </div>
        {can('userProfile::edit') && (
          <button className="btn btn-primary" onClick={openCreate}>
            Thêm hợp đồng
          </button>
        )}
      </div>

      <div className="table-container">
        <div className="table-toolbar">
          <div className="search-bar">
            <input
              className="form-control"
              placeholder="Tìm kiếm theo số hợp đồng..."
              value={search}
              onChange={e => setSearch(e.target.value)}
            />
          </div>
          <div className="text-muted" style={{ fontSize: 13 }}>
            {total} kết quả
          </div>
        </div>
        <Table columns={columns} data={contracts} loading={loading} emptyMessage="Chưa có hợp đồng lao động nào" />
        <Pagination page={page} totalItems={total} pageSize={PAGE_SIZE} onPageChange={setPage} />
      </div>

      {/* Create Modal */}
      <ContractForm title="Thêm hợp đồng" onSubmit={handleCreate} isOpen={createOpen} onClose={() => setCreateOpen(false)} form={form} handleFormChange={handleFormChange} saving={saving} employees={employees} />

      {/* Edit Modal */}
      <ContractForm title="Sửa thông tin hợp đồng" onSubmit={handleEdit} isOpen={editOpen} onClose={() => setEditOpen(false)} form={form} handleFormChange={handleFormChange} saving={saving} employees={employees} />

      {/* Delete Confirm */}
      <ConfirmDialog
        isOpen={deleteOpen}
        onClose={() => setDeleteOpen(false)}
        onConfirm={handleDelete}
        loading={deleting}
        title="Xóa hợp đồng"
        message={`Bạn có chắc muốn xóa hợp đồng "${selected?.contract_code}"?`}
        confirmText="Xóa hợp đồng"
      />
    </div>
  );
}

function ContractForm({ title, onSubmit, isOpen, onClose, form, handleFormChange, saving, employees }) {
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
        <label className="form-label">Nhân viên <span className="required">*</span></label>
        <select 
          className="form-control" 
          name="employee_id" 
          value={form.employee_id} 
          onChange={handleFormChange}
          disabled={title.includes('Sửa')}
        >
          <option value="">-- Chọn nhân sự --</option>
          {employees.map(e => <option key={e.id} value={e.id}>{e.full_name} ({e.code})</option>)}
        </select>
      </div>
      <div className="form-group">
        <label className="form-label">Số hợp đồng <span className="required">*</span></label>
        <input 
          className="form-control" 
          name="contract_code" 
          value={form.contract_code} 
          onChange={handleFormChange} 
          placeholder="VD: HĐLĐ-2024-001" 
        />
      </div>
      <div className="form-group">
        <label className="form-label">Mức lương cơ bản (P1) <span className="required">*</span></label>
        <input 
          className="form-control" 
          type="number"
          name="position_base_rate" 
          value={form.position_base_rate} 
          onChange={handleFormChange} 
          placeholder="VD: 10000000" 
        />
      </div>
      <div className="form-group">
        <label className="form-label">Trạng thái</label>
        <select className="form-control" name="status" value={form.status} onChange={handleFormChange}>
          <option value="ACTIVE">Hoạt động (ACTIVE)</option>
          <option value="EXPIRED">Hết hạn (EXPIRED)</option>
          <option value="TERMINATED">Chấm dứt (TERMINATED)</option>
        </select>
      </div>
      <div className="form-row">
        <div className="form-group">
          <label className="form-label">Từ ngày <span className="required">*</span></label>
          <input className="form-control" type="date" name="start_date" value={form.start_date} onChange={handleFormChange} />
        </div>
        <div className="form-group">
          <label className="form-label">Đến ngày</label>
          <input className="form-control" type="date" name="end_date" value={form.end_date} onChange={handleFormChange} />
        </div>
      </div>
    </Modal>
  );
}

export default ContractsPage;
