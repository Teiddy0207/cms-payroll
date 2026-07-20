import { useState, useEffect } from 'react';
import { calculatorAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';

export function FormulasPage() {
	const toast = useToast();
	const [formulas, setFormulas] = useState([]);
	const [loading, setLoading] = useState(false);
	const [modalOpen, setModalOpen] = useState(false);
	const [editingItem, setEditingItem] = useState(null);

	const [formData, setFormData] = useState({
		variable_name: '',
		expression: '',
		start_date: '',
		end_date: '',
		description: ''
	});

	const loadFormulas = async () => {
		setLoading(true);
		try {
			const res = await calculatorAPI.getFormulas();
			setFormulas(res.data?.data || []);
		} catch (err) {
			console.error("loadFormulas error:", err);
			toast.error('Lỗi', 'Không thể tải danh sách công thức');
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		loadFormulas();
	}, []);

	const handleOpenCreate = () => {
		setEditingItem(null);
		setFormData({
			variable_name: 'GROSS_SALARY',
			expression: 'P1 + P2',
			start_date: new Date().toISOString().split('T')[0],
			end_date: '',
			description: ''
		});
		setModalOpen(true);
	};

	const handleOpenEdit = (item) => {
		setEditingItem(item);
		setFormData({
			variable_name: item.variable_name,
			expression: item.expression,
			start_date: item.start_date ? item.start_date.split('T')[0] : '',
			end_date: item.end_date ? item.end_date.split('T')[0] : '',
			description: item.description || ''
		});
		setModalOpen(true);
	};

	const handleInputChange = (field, value) => {
		setFormData(prev => ({ ...prev, [field]: value }));
	};

	const handleSubmit = async (e) => {
		e.preventDefault();
		if (!formData.variable_name || !formData.expression || !formData.start_date) {
			toast.warning('Chú ý', 'Vui lòng điền đầy đủ các thông tin bắt buộc');
			return;
		}

		const dataToSend = {
			variable_name: formData.variable_name,
			expression: formData.expression,
			start_date: new Date(formData.start_date).toISOString(),
			end_date: formData.end_date ? new Date(formData.end_date).toISOString() : null,
			description: formData.description
		};

		setLoading(true);
		try {
			if (editingItem) {
				await calculatorAPI.updateFormula(editingItem.id, dataToSend);
				toast.success('Thành công', 'Đã cập nhật công thức lương');
			} else {
				await calculatorAPI.createFormula(dataToSend);
				toast.success('Thành công', 'Đã tạo công thức lương mới');
			}
			setModalOpen(false);
			loadFormulas();
		} catch (err) {
			console.error("submit error:", err);
			const errMsg = err.response?.data?.message || 'Không thể lưu công thức lương';
			toast.error('Lỗi cú pháp', errMsg);
		} finally {
			setLoading(false);
		}
	};

	const handleDelete = async (id) => {
		if (!window.confirm('Bạn có chắc chắn muốn xoá công thức này?')) return;
		setLoading(true);
		try {
			await calculatorAPI.deleteFormula(id);
			toast.success('Thành công', 'Đã xoá công thức lương');
			loadFormulas();
		} catch (err) {
			console.error("delete error:", err);
			toast.error('Lỗi', 'Không thể xoá công thức');
		} finally {
			setLoading(false);
		}
	};

	const formatDate = (dateStr) => {
		if (!dateStr) return '—';
		const d = new Date(dateStr);
		return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`;
	};

	const columns = [
		{
			key: 'variable_name',
			title: 'Tên biến',
			render: (val) => <span style={{ fontWeight: 700, color: 'var(--accent)' }}>{val}</span>
		},
		{
			key: 'expression',
			title: 'Công thức toán học',
			render: (val) => <code style={{ padding: '4px 8px', background: 'var(--bg-primary)', borderRadius: 4, fontFamily: 'monospace', fontSize: 13 }}>{val}</code>
		},
		{
			key: 'start_date',
			title: 'Ngày bắt đầu',
			render: (val) => formatDate(val)
		},
		{
			key: 'end_date',
			title: 'Ngày kết thúc',
			render: (val) => val ? formatDate(val) : <Badge variant="success">Vô hạn</Badge>
		},
		{
			key: 'description',
			title: 'Mô tả',
			render: (val) => val || <span style={{ color: 'var(--text-muted)' }}>Không có</span>
		},
		{
			key: 'actions',
			title: 'Thao tác',
			render: (_, row) => (
				<div style={{ display: 'flex', gap: 8 }}>
					<button className="btn btn-sm btn-secondary" onClick={() => handleOpenEdit(row)}>Sửa</button>
					<button className="btn btn-sm btn-danger" onClick={() => handleDelete(row.id)}>Xoá</button>
				</div>
			)
		}
	];

	return (
		<div className="page-container">
			<div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
				<div>
					<h1 className="page-title">Công thức lương động</h1>
					<p className="page-subtitle" style={{ color: 'var(--text-muted)', fontSize: 14 }}>
						Cấu hình các tham số và công thức tính toán Gross, Thuế TNCN và Thực nhận theo thời gian hiệu lực
					</p>
				</div>
				<button className="btn btn-primary" onClick={handleOpenCreate}>
					➕ Thêm công thức
				</button>
			</div>

			<div className="card">
				<Table
					columns={columns}
					data={formulas}
					loading={loading}
					emptyMessage="Chưa có công thức lương nào được cấu hình"
				/>
			</div>

			<Modal
				isOpen={modalOpen}
				onClose={() => setModalOpen(false)}
				title={editingItem ? "Cập nhật công thức lương" : "Thêm công thức lương mới"}
				size="md"
			>
				<form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
					<div className="form-group">
						<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Tên biến hiển thị <span style={{ color: 'var(--error)' }}>*</span></label>
						<input
							type="text"
							className="form-input"
							placeholder="Ví dụ: GROSS_SALARY, TAX, NET_SALARY, WORKDAY_RATIO"
							value={formData.variable_name}
							onChange={e => handleInputChange('variable_name', e.target.value.toUpperCase().trim())}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
						<p style={{ color: 'var(--text-muted)', fontSize: 11, marginTop: 4 }}>
							Biến hệ thống: <code>GROSS_SALARY</code>, <code>TAX</code>, <code>NET_SALARY</code>. Hoặc nhập tên biến bất kỳ để tự định nghĩa hệ số/phụ cấp mới.
						</p>
					</div>

					<div className="form-group">
						<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Công thức tính toán <span style={{ color: 'var(--error)' }}>*</span></label>
						<input
							type="text"
							className="form-input"
							placeholder="Ví dụ: P1 + P2 + 200000"
							value={formData.expression}
							onChange={e => handleInputChange('expression', e.target.value)}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)', fontFamily: 'monospace' }}
						/>
						<p style={{ color: 'var(--text-muted)', fontSize: 11, marginTop: 4 }}>
							Biến khả dụng: <code>P1</code>, <code>P2</code>, <code>P3</code>, <code>GROSS_SALARY</code>, <code>TAX</code>, <code>NET_SALARY</code>.
						</p>
					</div>

					<div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
						<div className="form-group">
							<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Ngày bắt đầu <span style={{ color: 'var(--error)' }}>*</span></label>
							<input
								type="date"
								className="form-input"
								value={formData.start_date}
								onChange={e => handleInputChange('start_date', e.target.value)}
								style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
							/>
						</div>
						<div className="form-group">
							<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Ngày kết thúc</label>
							<input
								type="date"
								className="form-input"
								value={formData.end_date}
								onChange={e => handleInputChange('end_date', e.target.value)}
								style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
							/>
						</div>
					</div>

					<div className="form-group">
						<label className="form-label" style={{ fontWeight: 600, display: 'block', marginBottom: 6 }}>Mô tả</label>
						<textarea
							className="form-input"
							rows="3"
							placeholder="Ghi chú về đợt điều chỉnh lương..."
							value={formData.description}
							onChange={e => handleInputChange('description', e.target.value)}
							style={{ width: '100%', padding: '8px 12px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)', background: 'var(--bg-card)', color: 'var(--text-primary)' }}
						/>
					</div>

					<div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 12 }}>
						<button type="button" className="btn btn-secondary" onClick={() => setModalOpen(false)}>Huỷ</button>
						<button type="submit" className="btn btn-primary">Lưu công thức</button>
					</div>
				</form>
			</Modal>
		</div>
	);
}

export default FormulasPage;
