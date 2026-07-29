import { useState, useEffect, useCallback } from 'react';
import { jobPositionsAPI, departmentsAPI, settingsAPI } from '../api/client.js';
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
  department_id: '',
  e_score: 0,
  c_score: 0,
  r_score: 0,
  we_weight: 0,
  wc_weight: 0,
  wr_weight: 0,
  salary_spread: 0,
  is_benchmark: false,
  market_salary: 0,
  search_keyword: '',
};

export function JobPositionsPage() {
  const toast = useToast();
  const { can } = usePermissions();
  const [positions, setPositions] = useState([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(false);

  const [departments, setDepartments] = useState([]);

  // Modals
  const [createOpen, setCreateOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);

  const [form, setForm] = useState(emptyForm);
  const [selected, setSelected] = useState(null);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);

  const [kFactor, setKFactor] = useState(4000000);
  const [kCalculating, setKCalculating] = useState(false);
  const [scrapingLogs, setScrapingLogs] = useState([]);
  const [scrapingActive, setScrapingActive] = useState(false);
  const [scrapedJobs, setScrapedJobs] = useState([]);
  const [scrapingSource, setScrapingSource] = useState('TopCV');

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
      const [deptRes, settingsRes] = await Promise.allSettled([
        departmentsAPI.list({ page_size: 200 }),
        settingsAPI.list(),
      ]);
      if (deptRes.status === 'fulfilled') setDepartments(deptRes.value.data?.data?.items || []);
      if (settingsRes.status === 'fulfilled') {
        const settings = settingsRes.value.data?.data || [];
        const kSetting = settings.find(s => s.key === 'payroll_k_factor');
        if (kSetting && kSetting.value) {
          setKFactor(parseFloat(kSetting.value) || 4000000);
        }
      }
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
      e_score: pos.e_score ?? 0,
      c_score: pos.c_score ?? 0,
      r_score: pos.r_score ?? 0,
      we_weight: (pos.we_weight ?? 0) * 100,
      wc_weight: (pos.wc_weight ?? 0) * 100,
      wr_weight: (pos.wr_weight ?? 0) * 100,
      salary_spread: (pos.salary_spread ?? 0) * 100,
      is_benchmark: pos.is_benchmark || false,
      market_salary: pos.market_salary ?? 0,
      search_keyword: pos.search_keyword || '',
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

  const getPayload = () => {
    return {
      ...form,
      e_score: parseFloat(form.e_score) || 0,
      c_score: parseFloat(form.c_score) || 0,
      r_score: parseFloat(form.r_score) || 0,
      we_weight: (parseFloat(form.we_weight) || 0) / 100,
      wc_weight: (parseFloat(form.wc_weight) || 0) / 100,
      wr_weight: (parseFloat(form.wr_weight) || 0) / 100,
      salary_spread: (parseFloat(form.salary_spread) || 0) / 100,
      is_benchmark: !!form.is_benchmark,
      market_salary: parseFloat(form.market_salary) || 0,
      search_keyword: form.search_keyword || '',
    };
  };

  const handleCreate = async () => {
    if (!form.code || !form.name) {
      toast.error('Thiếu thông tin', 'Vui lòng nhập mã và tên vị trí');
      return;
    }
    setSaving(true);
    try {
      await jobPositionsAPI.create(getPayload());
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
      await jobPositionsAPI.update(selected.id, getPayload());
      toast.success('Thành công', 'Đã cập nhật vị trí công việc');
      setEditOpen(false);
      loadPositions(page, search);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Không thể cập nhật');
    } finally {
      setSaving(false);
    }
  };

  const handleCalculateK = async () => {
    setKCalculating(true);
    try {
      const res = await jobPositionsAPI.calculateK();
      const newK = res.data?.data?.new_k_factor || 4000000;
      toast.success('Thành công', `Đã tính toán lại hệ số K bằng thuật toán OLS thành công! Hệ số K mới: ${Math.round(newK).toLocaleString('vi-VN')}đ/điểm`);
      setKFactor(newK);
      loadPositions(page, search);
    } catch (err) {
      toast.error('Lỗi tính toán', err.response?.data?.message || 'Có lỗi xảy ra');
    } finally {
      setKCalculating(false);
    }
  };

  const handleScrapeMarketSalary = async () => {
    if (!form.search_keyword) {
      toast.error('Lỗi', 'Vui lòng nhập từ khóa quét lương tuyển dụng');
      return;
    }
    setScrapingActive(true);
    setScrapingLogs(['[Hệ thống] Khởi tạo yêu cầu quét dữ liệu...']);
    setScrapedJobs([]);
    try {
      const res = await jobPositionsAPI.scrapeMarketSalary(selected?.id, {
        source: scrapingSource,
        keyword: form.search_keyword
      });
      
      const serverLogs = res.data?.data?.logs || [];
      const averageSalary = res.data?.data?.average_salary || 0;
      const jobs = res.data?.data?.jobs || [];
      
      let currentLogIndex = 0;
      const interval = setInterval(() => {
        if (currentLogIndex < serverLogs.length) {
          setScrapingLogs(prev => [...prev, serverLogs[currentLogIndex]]);
          currentLogIndex++;
        } else {
          clearInterval(interval);
          setForm(prev => ({
            ...prev,
            market_salary: averageSalary
          }));
          setScrapedJobs(jobs);
          toast.success('Quét thành công', `Hệ thống đã tính được lương trung vị thị trường: ${Math.round(averageSalary).toLocaleString('vi-VN')}đ`);
        }
      }, 350);
      
    } catch (err) {
      setScrapingLogs(prev => [...prev, '[Lỗi] Không thể kết nối với Scraper API.']);
      toast.error('Lỗi quét dữ liệu', err.response?.data?.message || 'Có lỗi xảy ra');
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
    const { name, type, checked, value } = e.target;
    setForm((p) => ({ ...p, [name]: type === 'checkbox' ? checked : value }));
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
      render: (val, row) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span className="font-bold">{val}</span>
          {row.is_benchmark && <Badge color="green">Mấu chốt</Badge>}
        </div>
      ),
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
      key: 'min_salary',
      title: 'Khung lương P1 (Min - Max)',
      render: (_, row) => {
        if (row.min_salary && row.max_salary) {
          return (
            <span className="font-bold text-success" style={{ fontSize: 13 }}>
              {Math.round(row.min_salary).toLocaleString('vi-VN')}đ - {Math.round(row.max_salary).toLocaleString('vi-VN')}đ
            </span>
          );
        }
        return <span className="text-muted">—</span>;
      }
    },
    {
      key: 'id',
      title: 'Thao tác',
      style: { width: 220 },
      render: (_, row) => (
        <div className="td-actions">
          {can('jobPosition::edit') && (
            <button className="btn btn-secondary btn-sm" onClick={() => openEdit(row)}>
              Sửa
            </button>
          )}
          {can('jobPosition::delete') && (
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
          <h2 className="page-title">Vị trí công việc</h2>
          <p className="page-subtitle">Quản lý các chức danh & tiêu chuẩn lương cứng P1 ({total} vị trí)</p>
        </div>
        {can('jobPosition::edit') && (
          <button className="btn btn-primary" onClick={openCreate}>
            Thêm vị trí
          </button>
        )}
      </div>

      <div style={{
        background: 'var(--bg-secondary)',
        border: '1px solid var(--border-color)',
        borderRadius: 'var(--radius-md)',
        padding: '16px 24px',
        marginBottom: '24px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        gap: '20px'
      }}>
        <div>
          <div style={{ fontSize: '13px', color: 'var(--text-muted)', fontWeight: 500 }}>Đường cong tiền lương (Salary Curve)</div>
          <div style={{ fontSize: '18px', fontWeight: 700, color: 'var(--text-primary)', marginTop: '4px' }}>
            Hệ số K hiện tại: <span style={{ color: 'var(--accent)' }}>{Math.round(kFactor).toLocaleString('vi-VN')}đ</span> / điểm giá trị
          </div>
          <div style={{ fontSize: '12px', color: 'var(--text-muted)', marginTop: '4px' }}>
            Được sử dụng làm đơn giá tính lương trung vị (Midpoint = Job Score x K). Hệ số K được tự động tối ưu hóa qua hồi quy OLS từ các vị trí mấu chốt.
          </div>
        </div>
        <button 
          className="btn btn-secondary" 
          onClick={handleCalculateK}
          disabled={kCalculating}
          style={{ whiteSpace: 'nowrap' }}
        >
          {kCalculating ? 'Đang tính toán...' : 'Tính toán lại hệ số K'}
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
      <PositionForm 
        title="Thêm vị trí" 
        onSubmit={handleCreate} 
        isOpen={createOpen} 
        onClose={() => setCreateOpen(false)} 
        form={form} 
        setForm={setForm}
        handleFormChange={handleFormChange} 
        saving={saving} 
        departments={departments} 
        kFactor={kFactor} 
        scrapingActive={scrapingActive}
        scrapingLogs={scrapingLogs}
        scrapingSource={scrapingSource}
        setScrapingSource={setScrapingSource}
        handleScrapeMarketSalary={handleScrapeMarketSalary}
      />

      {/* Edit Modal */}
      <PositionForm 
        title="Sửa thông tin vị trí" 
        onSubmit={handleEdit} 
        isOpen={editOpen} 
        onClose={() => setEditOpen(false)} 
        form={form} 
        setForm={setForm}
        handleFormChange={handleFormChange} 
        saving={saving} 
        departments={departments} 
        kFactor={kFactor} 
        scrapingActive={scrapingActive}
        scrapingLogs={scrapingLogs}
        scrapingSource={scrapingSource}
        setScrapingSource={setScrapingSource}
        handleScrapeMarketSalary={handleScrapeMarketSalary}
      />



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

function PositionForm({ 
  title, 
  onSubmit, 
  isOpen, 
  onClose, 
  form, 
  setForm,
  handleFormChange, 
  saving, 
  departments, 
  kFactor,
  scrapingActive,
  scrapingLogs,
  scrapingSource,
  setScrapingSource,
  handleScrapeMarketSalary
}) {
  const eVal = parseFloat(form.e_score) || 0;
  const cVal = parseFloat(form.c_score) || 0;
  const rVal = parseFloat(form.r_score) || 0;
  const weVal = (parseFloat(form.we_weight) || 0) / 100;
  const wcVal = (parseFloat(form.wc_weight) || 0) / 100;
  const wrVal = (parseFloat(form.wr_weight) || 0) / 100;
  const spreadVal = (parseFloat(form.salary_spread) || 0) / 100;

  const totalWeight = weVal + wcVal + wrVal;
  const jobScore = (eVal * weVal) + (cVal * wcVal) + (rVal * wrVal);
  const midpoint = jobScore * kFactor;
  let minSal = midpoint;
  let maxSal = midpoint;
  if (spreadVal > 0) {
    minSal = midpoint / (1 + spreadVal / 2);
    maxSal = minSal * (1 + spreadVal);
  }

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
          style={{ minHeight: 60 }}
        />
      </div>

      <div className="section-title" style={{ fontSize: '13px', fontWeight: 600, marginTop: 20, marginBottom: 12, borderTop: '1px solid #eee', paddingTop: 15, color: '#495057' }}>
        Định giá lương P1 (Job Evaluation)
      </div>
      
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '10px 15px', marginBottom: 15 }}>
        <div className="form-group" style={{ marginBottom: 0 }}>
          <label className="form-label">Chuyên môn (E)</label>
          <input 
            type="number" 
            className="form-control" 
            name="e_score" 
            min="0" 
            max="5" 
            step="0.1"
            value={form.e_score} 
            onChange={handleFormChange} 
          />
        </div>
        <div className="form-group" style={{ marginBottom: 0 }}>
          <label className="form-label">Trọng số E (%)</label>
          <input 
            type="number" 
            className="form-control" 
            name="we_weight" 
            min="0" 
            max="100" 
            value={form.we_weight} 
            onChange={handleFormChange} 
          />
        </div>

        <div className="form-group" style={{ marginBottom: 0 }}>
          <label className="form-label">Phức tạp (C)</label>
          <input 
            type="number" 
            className="form-control" 
            name="c_score" 
            min="0" 
            max="5" 
            step="0.1"
            value={form.c_score} 
            onChange={handleFormChange} 
          />
        </div>
        <div className="form-group" style={{ marginBottom: 0 }}>
          <label className="form-label">Trọng số C (%)</label>
          <input 
            type="number" 
            className="form-control" 
            name="wc_weight" 
            min="0" 
            max="100" 
            value={form.wc_weight} 
            onChange={handleFormChange} 
          />
        </div>

        <div className="form-group" style={{ marginBottom: 0 }}>
          <label className="form-label">Trách nhiệm (R)</label>
          <input 
            type="number" 
            className="form-control" 
            name="r_score" 
            min="0" 
            max="5" 
            step="0.1"
            value={form.r_score} 
            onChange={handleFormChange} 
          />
        </div>
        <div className="form-group" style={{ marginBottom: 0 }}>
          <label className="form-label">Trọng số R (%)</label>
          <input 
            type="number" 
            className="form-control" 
            name="wr_weight" 
            min="0" 
            max="100" 
            value={form.wr_weight} 
            onChange={handleFormChange} 
          />
        </div>

        <div className="form-group" style={{ gridColumn: 'span 2', marginBottom: 0 }}>
          <label className="form-label">Độ rộng dải lương (Spread Rp %)</label>
          <input 
            type="number" 
            className="form-control" 
            name="salary_spread" 
            min="0" 
            max="100" 
            value={form.salary_spread} 
            onChange={handleFormChange} 
            placeholder="VD: 40"
          />
        </div>
      </div>

      <div className="section-title" style={{ fontSize: '13px', fontWeight: 600, marginTop: 20, marginBottom: 12, borderTop: '1px solid #eee', paddingTop: 15, color: '#495057' }}>
        Định vị Lương thị trường (Market Benchmarking)
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', marginBottom: '15px' }}>
        <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer', fontSize: '13.5px', fontWeight: 500 }}>
          <input 
            type="checkbox" 
            name="is_benchmark" 
            checked={!!form.is_benchmark}
            onChange={(e) => setForm(p => ({ ...p, is_benchmark: e.target.checked }))}
            style={{ width: '16px', height: '16px' }}
          />
          Là vị trí mấu chốt (Benchmark Job)
        </label>

        {form.is_benchmark && (
          <div style={{ padding: '12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: '8px', display: 'flex', flexDirection: 'column', gap: '10px' }}>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '10px' }}>
              <div className="form-group" style={{ marginBottom: 0 }}>
                <label className="form-label">Nguồn khảo sát</label>
                <select 
                  className="form-control" 
                  value={scrapingSource} 
                  onChange={(e) => setScrapingSource(e.target.value)}
                >
                  <option value="TopCV">TopCV.vn</option>
                  <option value="VietnamWorks">VietnamWorks.com</option>
                  <option value="ITviec">ITviec.com</option>
                </select>
              </div>

              <div className="form-group" style={{ marginBottom: 0 }}>
                <label className="form-label">Từ khóa cào dữ liệu</label>
                <input 
                  type="text" 
                  className="form-control" 
                  name="search_keyword" 
                  value={form.search_keyword} 
                  onChange={handleFormChange}
                  placeholder="VD: Golang Developer"
                />
              </div>
            </div>

            <div style={{ display: 'flex', gap: '10px', alignItems: 'flex-end' }}>
              <div className="form-group" style={{ flex: 1, marginBottom: 0 }}>
                <label className="form-label">Lương thị trường (VND)</label>
                <input 
                  type="number" 
                  className="form-control" 
                  name="market_salary" 
                  value={form.market_salary} 
                  onChange={handleFormChange}
                  placeholder="Điền tự động hoặc nhập tay"
                />
              </div>
              <button 
                type="button" 
                className="btn btn-secondary" 
                onClick={handleScrapeMarketSalary}
                style={{ height: '40px', padding: '0 16px' }}
              >
                Quét dữ liệu lương
              </button>
            </div>

            {scrapingActive && (
              <div style={{
                background: '#0f172a',
                color: '#4ade80',
                fontFamily: 'monospace',
                fontSize: '11px',
                padding: '10px',
                borderRadius: '6px',
                maxHeight: '120px',
                overflowY: 'auto',
                border: '1px solid #334155',
                marginTop: '8px'
              }}>
                {scrapingLogs.map((log, idx) => (
                  <div key={idx} style={{ marginBottom: '4px', lineHeight: '1.4' }}>{log}</div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>

      <div style={{ backgroundColor: 'var(--bg-secondary)', padding: 12, borderRadius: 8, marginTop: 15, border: '1px solid var(--border-color)' }}>
        <div style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-primary)', marginBottom: 8 }}>Xem trước tính toán P1 (Hệ số K: {kFactor.toLocaleString('vi-VN')}đ)</div>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px 12px', fontSize: 12 }}>
          <div>Tổng trọng số: <strong style={{ color: Math.abs(totalWeight - 1) > 0.001 ? '#dc3545' : '#28a745' }}>{(totalWeight * 100).toFixed(0)}%</strong></div>
          <div>Job Score (S): <strong className="font-bold">{jobScore.toFixed(2)}</strong></div>
          <div>Lương trung vị (Mid): <strong>{Math.round(midpoint).toLocaleString('vi-VN')}đ</strong></div>
          <div>Độ rộng dải (Rp): <strong>{(spreadVal * 100).toFixed(0)}%</strong></div>
        </div>
        <div style={{ marginTop: 10, paddingTop: 8, borderTop: '1px dashed #ced4da', display: 'flex', justifyContent: 'space-between', fontSize: 12 }}>
          <div>Sàn lương (Min): <strong style={{ color: '#28a745' }}>{Math.round(minSal).toLocaleString('vi-VN')}đ</strong></div>
          <div>Trần lương (Max): <strong style={{ color: '#28a745' }}>{Math.round(maxSal).toLocaleString('vi-VN')}đ</strong></div>
        </div>
        {Math.abs(totalWeight - 1) > 0.001 && totalWeight > 0 && (
          <div style={{ fontSize: 11, color: '#dc3545', marginTop: 8 }}>
            ⚠️ Tổng trọng số các yếu tố phải bằng 100% (hiện tại là {(totalWeight * 100).toFixed(0)}%)
          </div>
        )}
      </div>
    </Modal>
  );
}

export default JobPositionsPage;
