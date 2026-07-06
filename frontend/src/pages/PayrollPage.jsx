import { useState, useEffect, useCallback } from 'react';
import { employeesAPI, calculatorAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';

export function PayrollPage() {
  const toast = useToast();
  const [loading, setLoading] = useState(false);
  const [payrollData, setPayrollData] = useState([]);
  const [selectedPeriod, setSelectedPeriod] = useState(() => {
    const d = new Date();
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
  });

  // Modal breakdown detail
  const [selectedRow, setSelectedRow] = useState(null);
  const [detailOpen, setDetailOpen] = useState(false);

  // Load all employees and calculate payroll for selected period
  const calculatePayroll = useCallback(async (period) => {
    setLoading(true);
    setPayrollData([]);
    try {
      // 1. Fetch all employees (max 200 for company preview)
      const empRes = await employeesAPI.list({ page_size: 200 });
      const empList = empRes.data?.data?.items || [];

      if (empList.length === 0) {
        setLoading(false);
        return;
      }

      // 2. Fetch preview for all employees in parallel
      const promises = empList.map(async (emp) => {
        try {
          const res = await calculatorAPI.preview(emp.id, period);
          return {
            employee: emp,
            success: true,
            data: res.data?.data || null,
            errorMsg: ''
          };
        } catch (err) {
          return {
            employee: emp,
            success: false,
            data: null,
            errorMsg: err.response?.data?.message || 'Thiếu thông tin hợp đồng/năng lực/vị trí'
          };
        }
      });

      const results = await Promise.all(promises);
      setPayrollData(results);
      toast.success('Thành công', `Đã tính toán bảng lương cho ${results.length} nhân sự`);
    } catch {
      toast.error('Lỗi', 'Không thể tải danh sách nhân sự');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    calculatePayroll(selectedPeriod);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedPeriod]);

  // Export CSV
  const handleExportCSV = () => {
    if (payrollData.length === 0) return;
    let csvContent = 'data:text/csv;charset=utf-8,\uFEFF';
    csvContent += 'Mã NV,Họ tên,Vị trí,Lương P1 (Cơ bản),Lương P2 (Năng lực),Tổng lương,Trạng thái\n';

    payrollData.forEach((row) => {
      const code = row.employee.code || '';
      const name = row.employee.full_name || '';
      const position = row.employee.job_position?.name || '';
      let p1 = 0;
      let p2 = 0;
      let total = 0;
      let statusStr = '';

      if (row.success && row.data) {
        p1 = row.data.p1 || row.data.salary_p1 || 0;
        p2 = row.data.p2 || row.data.salary_p2 || 0;
        total = row.data.total || row.data.total_salary || (p1 + p2);
        statusStr = 'Thành công';
      } else {
        statusStr = row.errorMsg || 'Lỗi';
      }

      csvContent += `"${code}","${name}","${position}",${p1},${p2},${total},"${statusStr}"\n`;
    });

    const encodedUri = encodeURI(csvContent);
    const link = document.createElement('a');
    link.setAttribute('href', encodedUri);
    link.setAttribute('download', `BangLuong_${selectedPeriod}.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const columns = [
    {
      key: 'code',
      title: 'Mã NV',
      render: (_, row) => <Badge color="cyan">{row.employee.code}</Badge>,
    },
    {
      key: 'name',
      title: 'Họ và tên',
      render: (_, row) => <span className="font-bold">{row.employee.full_name}</span>,
    },
    {
      key: 'position',
      title: 'Vị trí',
      render: (_, row) => row.employee.job_position?.name 
        ? <Badge color="yellow">{row.employee.job_position.name}</Badge>
        : <span className="text-muted">—</span>,
    },
    {
      key: 'p1',
      title: 'Lương P1 (Cơ bản)',
      render: (_, row) => {
        if (!row.success) return <span className="text-muted">0 đ</span>;
        const p1 = row.data?.p1 || row.data?.salary_p1 || 0;
        return <span className="font-bold">{p1.toLocaleString('vi-VN')} đ</span>;
      }
    },
    {
      key: 'p2',
      title: 'Lương P2 (Năng lực)',
      render: (_, row) => {
        if (!row.success) return <span className="text-muted">0 đ</span>;
        const p2 = row.data?.p2 || row.data?.salary_p2 || 0;
        return <span className="font-bold text-success">{p2.toLocaleString('vi-VN')} đ</span>;
      }
    },
    {
      key: 'total',
      title: 'Tổng lương (P1 + P2)',
      render: (_, row) => {
        if (!row.success) return <span className="text-muted" style={{ fontSize: 13 }}>0 đ (Chưa cấu hình)</span>;
        const p1 = row.data?.p1 || row.data?.salary_p1 || 0;
        const p2 = row.data?.p2 || row.data?.salary_p2 || 0;
        const total = row.data?.total || row.data?.total_salary || (p1 + p2);
        return <span className="font-bold text-primary" style={{ fontSize: 15 }}>{total.toLocaleString('vi-VN')} đ</span>;
      }
    },
    {
      key: 'action',
      title: 'Chi tiết',
      style: { width: 100 },
      render: (_, row) => (
        <button 
          className="btn btn-secondary btn-sm" 
          disabled={!row.success}
          onClick={() => {
            setSelectedRow(row);
            setDetailOpen(true);
          }}
        >
          Xem
        </button>
      )
    }
  ];

  return (
    <div>
      <div className="page-header">
        <div>
          <h2 className="page-title">Bảng tính lương công ty</h2>
          <p className="page-subtitle">Tính toán và hiển thị tổng hợp lương P1, P2 toàn thể nhân sự trong công ty</p>
        </div>
        <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
          <input 
            type="month" 
            className="form-control" 
            value={selectedPeriod} 
            onChange={e => setSelectedPeriod(e.target.value)} 
            style={{ width: 160 }}
          />
          <button className="btn btn-secondary" onClick={() => calculatePayroll(selectedPeriod)} disabled={loading}>
            Tính lại
          </button>
          <button className="btn btn-primary" onClick={handleExportCSV} disabled={loading || payrollData.length === 0}>
            Xuất file CSV
          </button>
        </div>
      </div>

      <div className="table-container">
        <Table columns={columns} data={payrollData} loading={loading} emptyMessage="Chưa có dữ liệu tính lương" />
      </div>

      {/* Detail Breakdown Modal */}
      <Modal 
        isOpen={detailOpen} 
        onClose={() => setDetailOpen(false)} 
        title={`Chi tiết lương — ${selectedRow?.employee?.full_name || ''}`}
        size="md"
        footer={<button className="btn btn-secondary" onClick={() => setDetailOpen(false)}>Đóng</button>}
      >
        {selectedRow?.data ? (
          <div>
            <div className="grid-2" style={{ marginBottom: 16 }}>
              <div className="salary-card">
                <div className="salary-label">Lương P1 (Cơ bản)</div>
                <div className="salary-value accent">
                  {(selectedRow.data.p1 || selectedRow.data.salary_p1 || 0).toLocaleString('vi-VN')} đ
                </div>
              </div>
              <div className="salary-card">
                <div className="salary-label">Lương P2 (Phụ cấp)</div>
                <div className="salary-value accent">
                  {(selectedRow.data.p2 || selectedRow.data.salary_p2 || 0).toLocaleString('vi-VN')} đ
                </div>
              </div>
            </div>
            
            <div className="salary-card highlight" style={{ marginBottom: 20 }}>
              <div className="salary-label">Tổng lương thực nhận</div>
              <div className="salary-value" style={{ fontSize: 26 }}>
                {(selectedRow.data.total || selectedRow.data.total_salary || (selectedRow.data.p1 || 0) + (selectedRow.data.p2 || 0)).toLocaleString('vi-VN')} đ
              </div>
            </div>

            {selectedRow.data.breakdown && (
              <div>
                <div className="section-title">Chi tiết cách tính</div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                  {Object.entries(selectedRow.data.breakdown).map(([k, v]) => (
                    <div key={k} style={{ display: 'flex', justifyContent: 'space-between', padding: '8px 12px', background: 'var(--bg-input)', borderRadius: 'var(--radius-sm)' }}>
                      <span className="text-secondary">{k}</span>
                      <span className="font-bold">{typeof v === 'number' ? v.toLocaleString('vi-VN') + ' đ' : v}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        ) : (
          <p className="text-muted">Không có thông tin chi tiết</p>
        )}
      </Modal>
    </div>
  );
}

export default PayrollPage;
