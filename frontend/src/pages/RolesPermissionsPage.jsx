import { useState, useEffect, useCallback } from 'react';
import { authAdminAPI, employeesAPI } from '../api/client.js';
import { Table } from '../components/Table.jsx';
import { Modal } from '../components/Modal.jsx';
import { ConfirmDialog } from '../components/ConfirmDialog.jsx';
import { Pagination } from '../components/Pagination.jsx';
import { Badge } from '../components/Badge.jsx';
import { useToast } from '../hooks/useToast.js';

const PAGE_SIZE = 15;

const emptyRoleForm = {
  name: '',
  slug: '',
  description: '',
};

const emptyPermissionForm = {
  name: '',
  resource: '',
  action: '',
  description: '',
};

const emptyUserForm = {
  username: '',
  email: '',
  password: '',
  is_active: true,
};

export default function RolesPermissionsPage() {
  const toast = useToast();
  const [activeTab, setActiveTab] = useState('roles');
  
  // Data States
  const [roles, setRoles] = useState([]);
  const [permissions, setPermissions] = useState([]);
  const [users, setUsers] = useState([]);
  const [employees, setEmployees] = useState([]);
  
  // Totals & Pagination
  const [totalRoles, setTotalRoles] = useState(0);
  const [totalPermissions, setTotalPermissions] = useState(0);
  const [totalUsers, setTotalUsers] = useState(0);
  
  const [pageRoles, setPageRoles] = useState(1);
  const [pagePermissions, setPagePermissions] = useState(1);
  const [pageUsers, setPageUsers] = useState(1);
  
  // Loadings
  const [loadingRoles, setLoadingRoles] = useState(false);
  const [loadingPermissions, setLoadingPermissions] = useState(false);
  const [loadingUsers, setLoadingUsers] = useState(false);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);

  // Modals Open
  const [roleModalOpen, setRoleModalOpen] = useState(false);
  const [permModalOpen, setPermModalOpen] = useState(false);
  const [userModalOpen, setUserModalOpen] = useState(false);
  const [assignRoleOpen, setAssignRoleOpen] = useState(false);
  const [assignPermsOpen, setAssignPermsOpen] = useState(false);
  
  // Confirm Delete Modals
  const [deleteRoleOpen, setDeleteRoleOpen] = useState(false);
  const [deletePermOpen, setDeletePermOpen] = useState(false);
  const [deleteUserOpen, setDeleteUserOpen] = useState(false);

  // Selected Data & Forms
  const [selectedRole, setSelectedRole] = useState(null);
  const [selectedPermission, setSelectedPermission] = useState(null);
  const [selectedUser, setSelectedUser] = useState(null);
  
  const [roleForm, setRoleForm] = useState(emptyRoleForm);
  const [permForm, setPermForm] = useState(emptyPermissionForm);
  const [userForm, setUserForm] = useState(emptyUserForm);
  const [userRoleSelect, setUserRoleSelect] = useState('');
  const [checkedPermissions, setCheckedPermissions] = useState([]);
  const [permSearchTerm, setPermSearchTerm] = useState('');

  // Load Metadata
  useEffect(() => {
    loadEmployees();
    loadPermissionsList();
  }, []);

  // Reload lists on page/tab change
  useEffect(() => {
    if (activeTab === 'roles') loadRoles(pageRoles);
    if (activeTab === 'users') loadUsers(pageUsers);
    if (activeTab === 'permissions') loadPermissions(pagePermissions);
  }, [activeTab, pageRoles, pageUsers, pagePermissions]);

  const loadEmployees = async () => {
    try {
      const res = await employeesAPI.list({ page_size: 500 });
      setEmployees(res.data?.data?.items || []);
    } catch { /* ignore */ }
  };

  const loadPermissionsList = async () => {
    try {
      const res = await authAdminAPI.listPermissions({ page_size: 500 });
      setPermissions(res.data?.data?.items || []);
    } catch { /* ignore */ }
  };

  const loadRoles = useCallback(async (p) => {
    setLoadingRoles(true);
    try {
      const res = await authAdminAPI.listRoles({ page_number: p, page_size: PAGE_SIZE });
      setRoles(res.data?.data?.items || []);
      setTotalRoles(res.data?.data?.total_items || 0);
    } catch {
      toast.error('Lỗi', 'Không thể tải danh sách vai trò');
    } finally {
      setLoadingRoles(false);
    }
  }, [toast]);

  const loadPermissions = useCallback(async (p) => {
    setLoadingPermissions(true);
    try {
      const res = await authAdminAPI.listPermissions({ page_number: p, page_size: PAGE_SIZE });
      setPermissions(res.data?.data?.items || []);
      setTotalPermissions(res.data?.data?.total_items || 0);
    } catch {
      toast.error('Lỗi', 'Không thể tải danh sách quyền');
    } finally {
      setLoadingPermissions(false);
    }
  }, [toast]);

  const loadUsers = useCallback(async (p) => {
    setLoadingUsers(true);
    try {
      const res = await authAdminAPI.listUsers({ page_number: p, page_size: PAGE_SIZE });
      setUsers(res.data?.data?.items || []);
      setTotalUsers(res.data?.data?.total_items || 0);
    } catch {
      toast.error('Lỗi', 'Không thể tải danh sách tài khoản');
    } finally {
      setLoadingUsers(false);
    }
  }, [toast]);

  // --- Handlers Role ---
  const handleRoleSubmit = async () => {
    if (!roleForm.name || !roleForm.slug) {
      toast.error('Lỗi', 'Vui lòng nhập tên và mã định danh vai trò');
      return;
    }
    setSaving(true);
    try {
      if (selectedRole) {
        await authAdminAPI.updateRole(selectedRole.id, roleForm);
        toast.success('Thành công', 'Đã cập nhật thông tin vai trò');
      } else {
        await authAdminAPI.createRole(roleForm);
        toast.success('Thành công', 'Đã tạo vai trò mới');
      }
      setRoleModalOpen(false);
      loadRoles(pageRoles);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Có lỗi xảy ra');
    } finally {
      setSaving(false);
    }
  };

  const handleRoleDelete = async () => {
    setDeleting(true);
    try {
      await authAdminAPI.deleteRole(selectedRole.id);
      toast.success('Thành công', 'Đã xóa vai trò thành công');
      setDeleteRoleOpen(false);
      loadRoles(pageRoles);
    } catch {
      toast.error('Lỗi', 'Không thể xóa vai trò mặc định hoặc đang có người dùng');
    } finally {
      setDeleting(false);
    }
  };

  // --- Handlers Permission ---
  const handlePermSubmit = async () => {
    if (!permForm.name || !permForm.resource || !permForm.action) {
      toast.error('Lỗi', 'Vui lòng điền đầy đủ thông tin quyền');
      return;
    }
    setSaving(true);
    try {
      if (selectedPermission) {
        await authAdminAPI.updatePermission(selectedPermission.id, permForm);
        toast.success('Thành công', 'Đã cập nhật quyền');
      } else {
        await authAdminAPI.createPermission(permForm);
        toast.success('Thành công', 'Đã thêm quyền mới');
      }
      setPermModalOpen(false);
      loadPermissions(pagePermissions);
      loadPermissionsList();
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Có lỗi xảy ra');
    } finally {
      setSaving(false);
    }
  };

  const handlePermDelete = async () => {
    setDeleting(true);
    try {
      await authAdminAPI.deletePermission(selectedPermission.id);
      toast.success('Thành công', 'Đã xóa quyền');
      setDeletePermOpen(false);
      loadPermissions(pagePermissions);
      loadPermissionsList();
    } catch {
      toast.error('Lỗi', 'Không thể xóa quyền này');
    } finally {
      setDeleting(false);
    }
  };

  // --- Handlers User ---
  const handleUserSubmit = async () => {
    if (!userForm.username || !userForm.email || (!selectedUser && !userForm.password)) {
      toast.error('Lỗi', 'Vui lòng nhập đầy đủ thông tin tài khoản');
      return;
    }
    setSaving(true);
    try {
      if (selectedUser) {
        await authAdminAPI.updateUser(selectedUser.id, userForm);
        toast.success('Thành công', 'Đã cập nhật tài khoản');
      } else {
        await authAdminAPI.createUser({
          ...userForm,
          confirmed_password: userForm.password
        });
        toast.success('Thành công', 'Đã tạo tài khoản thành công');
      }
      setUserModalOpen(false);
      loadUsers(pageUsers);
    } catch (err) {
      toast.error('Lỗi', err.response?.data?.message || 'Có lỗi xảy ra');
    } finally {
      setSaving(false);
    }
  };

  const handleUserDelete = async () => {
    setDeleting(true);
    try {
      await authAdminAPI.deleteUser(selectedUser.id);
      toast.success('Thành công', 'Đã xóa tài khoản');
      setDeleteUserOpen(false);
      loadUsers(pageUsers);
    } catch {
      toast.error('Lỗi', 'Không thể xóa tài khoản này');
    } finally {
      setDeleting(false);
    }
  };

  // --- Handlers Assignments ---
  const handleSaveRolePermissions = async () => {
    setSaving(false);
    try {
      await authAdminAPI.assignPermissionToRole({
        role_id: selectedRole.id,
        permission_id: checkedPermissions
      });
      toast.success('Thành công', 'Đã lưu phân quyền vai trò');
      setAssignPermsOpen(false);
      loadRoles(pageRoles);
    } catch {
      toast.error('Lỗi', 'Không thể gán quyền');
    }
  };

  const handleSaveUserRole = async () => {
    if (!userRoleSelect) {
      toast.error('Lỗi', 'Vui lòng chọn vai trò');
      return;
    }
    setSaving(true);
    try {
      const selectedRoleInfo = roles.find(r => r.id === userRoleSelect);
      await authAdminAPI.assignRoleToUser({
        id: userRoleSelect,
        user_id: selectedUser.id,
        code: selectedRoleInfo?.slug || '',
        name: selectedRoleInfo?.name || ''
      });
      toast.success('Thành công', 'Đã gán vai trò người dùng');
      setAssignRoleOpen(false);
      loadUsers(pageUsers);
    } catch {
      toast.error('Lỗi', 'Không thể gán vai trò');
    } finally {
      setSaving(false);
    }
  };

  // --- Group Permissions by Resource ---
  const getGroupedPermissions = () => {
    const term = permSearchTerm.trim().toLowerCase();
    const filtered = term
      ? permissions.filter((p) => (
          p.name?.toLowerCase().includes(term) ||
          p.slug?.toLowerCase().includes(term) ||
          p.resource?.toLowerCase().includes(term) ||
          p.action?.toLowerCase().includes(term)
        ))
      : permissions;

    const groups = {};
    filtered.forEach((p) => {
      const resName = p.resource || 'khác';
      if (!groups[resName]) groups[resName] = [];
      groups[resName].push(p);
    });
    return groups;
  };

  // --- Columns Configuration ---
  const columnsRoles = [
    {
      title: 'Tên vai trò',
      key: 'name',
      render: (_, r) => (
        <div>
          <div style={{ fontWeight: 600, color: 'var(--text-primary)' }}>{r.name}</div>
          <div style={{ fontSize: 11, color: 'var(--text-muted)' }}>Mã: <code>{r.slug}</code></div>
        </div>
      )
    },
    {
      title: 'Mô tả',
      key: 'description',
      render: (_, r) => <span className="text-muted">{r.description || '—'}</span>
    },
    {
      title: 'Loại',
      key: 'is_system',
      render: (_, r) => (
        <Badge color={r.is_system ? 'amber' : 'teal'}>
          {r.is_system ? 'Hệ thống' : 'Người dùng'}
        </Badge>
      )
    },
    {
      title: 'Thao tác',
      key: 'actions',
      style: { width: 220, textAlign: 'right' },
      render: (_, r) => (
        <div className="table-actions" style={{ display: 'flex', gap: 6, justifyContent: 'flex-end' }}>
          <button
            className="btn btn-secondary btn-sm"
            onClick={() => {
              setSelectedRole(r);
              setCheckedPermissions((r.permissions || []).map(p => p.id));
              setPermSearchTerm('');
              setAssignPermsOpen(true);
            }}
            title="Phân quyền chi tiết"
          >
            <i className="fa-solid fa-shield-halved" style={{ marginRight: 4 }}></i> Phân quyền
          </button>
          <button
            className="btn btn-secondary btn-sm"
            onClick={() => {
              setSelectedRole(r);
              setRoleForm({ name: r.name, slug: r.slug, description: r.description || '' });
              setRoleModalOpen(true);
            }}
          >
            Sửa
          </button>
          {!r.is_system && (
            <button
              className="btn btn-danger btn-sm"
              onClick={() => {
                setSelectedRole(r);
                setDeleteRoleOpen(true);
              }}
            >
              Xóa
            </button>
          )}
        </div>
      )
    }
  ];

  const columnsUsers = [
    {
      title: 'Tài khoản',
      key: 'username',
      render: (_, r) => (
        <div>
          <div style={{ fontWeight: 600, color: 'var(--text-primary)' }}>{r.username}</div>
          <div style={{ fontSize: 11, color: 'var(--text-muted)' }}>{r.email}</div>
        </div>
      )
    },
    {
      title: 'Vai trò hiện tại',
      key: 'roles',
      render: (_, r) => {
        const role = r.roles;
        if (!role) return <span className="text-muted" style={{ fontSize: 12 }}>Chưa có vai trò</span>;
        if (Array.isArray(role)) {
          return (
            <div style={{ display: 'flex', gap: 4, flexWrap: 'wrap' }}>
              {role.map(x => <Badge key={x.id} color="teal">{x.name}</Badge>)}
            </div>
          );
        }
        return <Badge color="teal">{role.name}</Badge>;
      }
    },
    {
      title: 'Họ và tên nhân viên',
      key: 'employee',
      render: (_, r) => {
        return <span>{r.user_profile?.full_name || <span className="text-muted">Tài khoản Admin</span>}</span>;
      }
    },
    {
      title: 'Trạng thái',
      key: 'is_active',
      render: (_, r) => (
        <Badge color={r.is_active ? 'teal' : 'rose'}>
          {r.is_active ? 'Hoạt động' : 'Đang khóa'}
        </Badge>
      )
    },
    {
      title: 'Thao tác',
      key: 'actions',
      style: { width: 220, textAlign: 'right' },
      render: (_, r) => (
        <div className="table-actions" style={{ display: 'flex', gap: 6, justifyContent: 'flex-end' }}>
          <button
            className="btn btn-secondary btn-sm"
            onClick={() => openAssignRole(r)}
          >
            <i className="fa-solid fa-user-tag" style={{ marginRight: 4 }}></i> Gán vai trò
          </button>
          <button
            className="btn btn-secondary btn-sm"
            onClick={() => {
              setSelectedUser(r);
              setUserForm({ username: r.username, email: r.email, password: '', is_active: r.is_active });
              setUserModalOpen(true);
            }}
          >
            Sửa
          </button>
          {r.username !== 'admin' && (
            <button
              className="btn btn-danger btn-sm"
              onClick={() => {
                setSelectedUser(r);
                setDeleteUserOpen(true);
              }}
            >
              Xóa
            </button>
          )}
        </div>
      )
    }
  ];

  const columnsPermissions = [
    {
      title: 'Tên quyền',
      key: 'name',
      render: (_, r) => <span style={{ fontWeight: 600, color: 'var(--text-primary)' }}>{r.name}</span>
    },
    {
      title: 'Đối tượng (Resource)',
      key: 'resource',
      render: (_, r) => <Badge color="indigo">{r.resource}</Badge>
    },
    {
      title: 'Hành động (Action)',
      key: 'action',
      render: (_, r) => <Badge color="sky">{r.action}</Badge>
    },
    {
      title: 'Mô tả',
      key: 'description',
      render: (_, r) => <span className="text-muted">{r.description || '—'}</span>
    },
    {
      title: 'Thao tác',
      key: 'actions',
      style: { width: 140, textAlign: 'right' },
      render: (_, r) => (
        <div className="table-actions" style={{ display: 'flex', gap: 6, justifyContent: 'flex-end' }}>
          <button
            className="btn btn-secondary btn-sm"
            onClick={() => {
              setSelectedPermission(r);
              setPermForm({ name: r.name, resource: r.resource, action: r.action, description: r.description || '' });
              setPermModalOpen(true);
            }}
          >
            Sửa
          </button>
          <button
            className="btn btn-danger btn-sm"
            onClick={() => {
              setSelectedPermission(r);
              setDeletePermOpen(true);
            }}
          >
            Xóa
          </button>
        </div>
      )
    }
  ];

  const openAssignRole = (user) => {
    setSelectedUser(user);
    const firstRole = user.roles && user.roles.length > 0 ? user.roles[0].id : '';
    setUserRoleSelect(firstRole);
    setAssignRoleOpen(true);
  };

  return (
    <div className="roles-permissions-page">
      <div className="page-header">
        <div>
          <h2 className="page-title">
            <i className="fa-solid fa-shield-halved" style={{ color: 'var(--accent)', marginRight: 10 }}></i> 
            Phân quyền hệ thống
          </h2>
          <p className="page-subtitle">Quản trị các vai trò, gán danh sách quyền hạn cho vai trò, và thiết lập tài khoản đăng nhập.</p>
        </div>
        <div style={{ display: 'flex', gap: 8 }}>
          {activeTab === 'roles' && (
            <button className="btn btn-primary" onClick={() => { setSelectedRole(null); setRoleForm(emptyRoleForm); setRoleModalOpen(true); }}>
              <i className="fa-solid fa-plus" style={{ marginRight: 6 }}></i> Thêm vai trò
            </button>
          )}
          {activeTab === 'users' && (
            <button className="btn btn-primary" onClick={() => { setSelectedUser(null); setUserForm(emptyUserForm); setUserModalOpen(true); }}>
              <i className="fa-solid fa-plus" style={{ marginRight: 6 }}></i> Tạo tài khoản
            </button>
          )}
          {activeTab === 'permissions' && (
            <button className="btn btn-primary" onClick={() => { setSelectedPermission(null); setPermForm(emptyPermissionForm); setPermModalOpen(true); }}>
              <i className="fa-solid fa-plus" style={{ marginRight: 6 }}></i> Thêm quyền hạn
            </button>
          )}
        </div>
      </div>

      {/* Tabs navigation */}
      <div className="tab-navigation">
        <button className={`tab-btn ${activeTab === 'roles' ? 'active' : ''}`} onClick={() => setActiveTab('roles')}>
          Vai trò (Roles)
        </button>
        <button className={`tab-btn ${activeTab === 'users' ? 'active' : ''}`} onClick={() => setActiveTab('users')}>
          Tài khoản (Users)
        </button>
        <button className={`tab-btn ${activeTab === 'permissions' ? 'active' : ''}`} onClick={() => setActiveTab('permissions')}>
          Quyền hạn (Permissions)
        </button>
      </div>

      {/* Main Content Area */}
      <div className="table-container">
        {activeTab === 'roles' && (
          <>
            <Table columns={columnsRoles} data={roles} loading={loadingRoles} emptyMessage="Chưa có vai trò nào" />
            <Pagination page={pageRoles} totalItems={totalRoles} pageSize={PAGE_SIZE} onPageChange={setPageRoles} />
          </>
        )}

        {activeTab === 'users' && (
          <>
            <Table columns={columnsUsers} data={users} loading={loadingUsers} emptyMessage="Chưa có tài khoản nào" />
            <Pagination page={pageUsers} totalItems={totalUsers} pageSize={PAGE_SIZE} onPageChange={setPageUsers} />
          </>
        )}

        {activeTab === 'permissions' && (
          <>
            <Table columns={columnsPermissions} data={permissions} loading={loadingPermissions} emptyMessage="Chưa có quyền hạn nào" />
            <Pagination page={pagePermissions} totalItems={totalPermissions} pageSize={PAGE_SIZE} onPageChange={setPagePermissions} />
          </>
        )}
      </div>

      {/* Create/Edit Role Modal */}
      <Modal
        isOpen={roleModalOpen}
        onClose={() => setRoleModalOpen(false)}
        title={selectedRole ? 'Cập nhật vai trò' : 'Thêm vai trò mới'}
        size="md"
        footer={
          <>
            <button className="btn btn-secondary" onClick={() => setRoleModalOpen(false)}>Hủy</button>
            <button className="btn btn-primary" onClick={handleRoleSubmit} disabled={saving}>
              {saving ? 'Đang lưu...' : 'Lưu vai trò'}
            </button>
          </>
        }
      >
        <div className="form-group">
          <label className="form-label">Tên vai trò <span className="required">*</span></label>
          <input
            className="form-control"
            value={roleForm.name}
            onChange={(e) => setRoleForm({ ...roleForm, name: e.target.value })}
            placeholder="VD: Quản lý nhân sự"
          />
        </div>
        <div className="form-group">
          <label className="form-label">Mã vai trò (Slug) <span className="required">*</span></label>
          <input
            className="form-control"
            value={roleForm.slug}
            onChange={(e) => setRoleForm({ ...roleForm, slug: e.target.value })}
            placeholder="VD: hr_manager"
            disabled={!!selectedRole}
          />
        </div>
        <div className="form-group">
          <label className="form-label">Mô tả vai trò</label>
          <textarea
            className="form-control"
            value={roleForm.description}
            onChange={(e) => setRoleForm({ ...roleForm, description: e.target.value })}
            placeholder="Mô tả các đặc quyền của vai trò này..."
            rows={3}
          />
        </div>
      </Modal>

      {/* Create/Edit Permission Modal */}
      <Modal
        isOpen={permModalOpen}
        onClose={() => setPermModalOpen(false)}
        title={selectedPermission ? 'Cập nhật quyền hạn' : 'Thêm quyền hạn mới'}
        size="md"
        footer={
          <>
            <button className="btn btn-secondary" onClick={() => setPermModalOpen(false)}>Hủy</button>
            <button className="btn btn-primary" onClick={handlePermSubmit} disabled={saving}>
              {saving ? 'Đang lưu...' : 'Lưu quyền'}
            </button>
          </>
        }
      >
        <div className="form-group">
          <label className="form-label">Tên quyền hạn <span className="required">*</span></label>
          <input
            className="form-control"
            value={permForm.name}
            onChange={(e) => setPermForm({ ...permForm, name: e.target.value })}
            placeholder="VD: Xem thông tin bảng công"
          />
        </div>
        <div className="form-group">
          <label className="form-label">Đối tượng (Resource) <span className="required">*</span></label>
          <input
            className="form-control"
            value={permForm.resource}
            onChange={(e) => setPermForm({ ...permForm, resource: e.target.value })}
            placeholder="VD: payroll, timekeeping, employees"
          />
        </div>
        <div className="form-group">
          <label className="form-label">Hành động (Action) <span className="required">*</span></label>
          <input
            className="form-control"
            value={permForm.action}
            onChange={(e) => setPermForm({ ...permForm, action: e.target.value })}
            placeholder="VD: read, write, create, delete"
          />
        </div>
        <div className="form-group">
          <label className="form-label">Mô tả chi tiết</label>
          <textarea
            className="form-control"
            value={permForm.description}
            onChange={(e) => setPermForm({ ...permForm, description: e.target.value })}
            placeholder="Giải thích chi tiết về quyền hạn này..."
            rows={3}
          />
        </div>
      </Modal>

      {/* Create/Edit User Modal */}
      <Modal
        isOpen={userModalOpen}
        onClose={() => setUserModalOpen(false)}
        title={selectedUser ? 'Sửa thông tin tài khoản' : 'Tạo tài khoản mới'}
        size="md"
        footer={
          <>
            <button className="btn btn-secondary" onClick={() => setUserModalOpen(false)}>Hủy</button>
            <button className="btn btn-primary" onClick={handleUserSubmit} disabled={saving}>
              {saving ? 'Đang lưu...' : 'Lưu tài khoản'}
            </button>
          </>
        }
      >
        <div className="form-group">
          <label className="form-label">Tên đăng nhập <span className="required">*</span></label>
          <input
            className="form-control"
            value={userForm.username}
            onChange={(e) => setUserForm({ ...userForm, username: e.target.value })}
            placeholder="VD: nguyenvanan"
            disabled={!!selectedUser}
          />
        </div>
        <div className="form-group">
          <label className="form-label">Email liên hệ <span className="required">*</span></label>
          <input
            className="form-control"
            type="email"
            value={userForm.email}
            onChange={(e) => setUserForm({ ...userForm, email: e.target.value })}
            placeholder="VD: an.nguyen@example.com"
          />
        </div>
        {!selectedUser && (
          <div className="form-group">
            <label className="form-label">Mật khẩu ban đầu <span className="required">*</span></label>
            <input
              className="form-control"
              type="password"
              value={userForm.password}
              onChange={(e) => setUserForm({ ...userForm, password: e.target.value })}
              placeholder="Nhập mật khẩu truy cập"
            />
          </div>
        )}
        <div className="form-group" style={{ display: 'flex', alignItems: 'center', gap: 8, marginTop: 12 }}>
          <input
            type="checkbox"
            id="user_active_chk"
            checked={userForm.is_active}
            onChange={(e) => setUserForm({ ...userForm, is_active: e.target.checked })}
            style={{ width: 16, height: 16, accentColor: 'var(--accent)' }}
          />
          <label htmlFor="user_active_chk" className="form-label" style={{ margin: 0, cursor: 'pointer' }}>Kích hoạt tài khoản</label>
        </div>
      </Modal>

      {/* Assign Role to User Modal */}
      <Modal
        isOpen={assignRoleOpen}
        onClose={() => setAssignRoleOpen(false)}
        title={selectedUser ? `Gán vai trò cho tài khoản "${selectedUser.username}"` : 'Gán vai trò'}
        size="md"
        footer={
          <>
            <button className="btn btn-secondary" onClick={() => setAssignRoleOpen(false)}>Hủy</button>
            <button className="btn btn-primary" onClick={handleSaveUserRole} disabled={saving}>
              {saving ? 'Đang lưu...' : 'Lưu thay đổi'}
            </button>
          </>
        }
      >
        <div className="form-group">
          <label className="form-label">Chọn vai trò chính</label>
          <select
            className="form-control"
            value={userRoleSelect}
            onChange={(e) => setUserRoleSelect(e.target.value)}
          >
            <option value="">-- Chọn vai trò hệ thống --</option>
            {roles.map((r) => (
              <option value={r.id} key={r.id}>
                {r.name} ({r.slug})
              </option>
            ))}
          </select>
        </div>
      </Modal>

      {/* Assign Permissions to Role Modal */}
      <Modal
        isOpen={assignPermsOpen}
        onClose={() => setAssignPermsOpen(false)}
        title={selectedRole ? `Thiết lập quyền cho vai trò "${selectedRole.name}"` : 'Thiết lập quyền'}
        size="lg"
        footer={
          <>
            <button className="btn btn-secondary" onClick={() => setAssignPermsOpen(false)}>Hủy</button>
            <button className="btn btn-primary" onClick={handleSaveRolePermissions}>Lưu phân quyền</button>
          </>
        }
      >
        <div className="form-group" style={{ marginBottom: 16 }}>
          <input
            className="form-control"
            value={permSearchTerm}
            onChange={(e) => setPermSearchTerm(e.target.value)}
            placeholder="Tìm theo tên quyền, mã (slug), resource hoặc action..."
          />
        </div>

        <div className="permissions-grid">
          {Object.keys(getGroupedPermissions()).length === 0 && (
            <div className="text-muted" style={{ padding: 12 }}>Không tìm thấy quyền hạn nào phù hợp</div>
          )}
          {Object.entries(getGroupedPermissions()).map(([resource, perms]) => (
            <div key={resource} className="resource-card">
              <div className="resource-card-header">{resource}</div>
              <div className="resource-card-body">
                {perms.map((p) => {
                  const isChecked = checkedPermissions.includes(p.id);
                  return (
                    <label key={p.id} className="perm-checkbox-label">
                      <input
                        type="checkbox"
                        checked={isChecked}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setCheckedPermissions([...checkedPermissions, p.id]);
                          } else {
                            setCheckedPermissions(checkedPermissions.filter(id => id !== p.id));
                          }
                        }}
                      />
                      <div>
                        <span className="perm-name">{p.name}</span>
                        <div style={{ fontSize: 10, color: 'var(--text-muted)' }}>Mã: {p.action}</div>
                      </div>
                    </label>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      </Modal>

      {/* Confirm deletes */}
      <ConfirmDialog
        isOpen={deleteRoleOpen}
        onClose={() => setDeleteRoleOpen(false)}
        onConfirm={handleRoleDelete}
        loading={deleting}
        title="Xóa vai trò"
        message={`Bạn có chắc chắn muốn xóa vai trò "${selectedRole?.name}"? Các tài khoản liên kết sẽ mất vai trò này.`}
        confirmText="Xóa vai trò"
      />

      <ConfirmDialog
        isOpen={deletePermOpen}
        onClose={() => setDeletePermOpen(false)}
        onConfirm={handlePermDelete}
        loading={deleting}
        title="Xóa quyền hạn"
        message={`Bạn có chắc chắn muốn xóa quyền hạn "${selectedPermission?.name}" khỏi hệ thống?`}
        confirmText="Xóa quyền hạn"
      />

      <ConfirmDialog
        isOpen={deleteUserOpen}
        onClose={() => setDeleteUserOpen(false)}
        onConfirm={handleUserDelete}
        loading={deleting}
        title="Xóa tài khoản"
        message={`Bạn có chắc chắn muốn xóa tài khoản "${selectedUser?.username}"? Thao tác này không thể hoàn tác.`}
        confirmText="Xóa tài khoản"
      />

      {/* Custom Styles */}
      <style>{`
        .tab-navigation {
          display: flex;
          gap: 6px;
          border-bottom: 2px solid var(--border-color);
          margin-bottom: 24px;
          padding-bottom: 2px;
        }
        .tab-btn {
          background: transparent;
          border: none;
          padding: 12px 24px;
          font-size: 14px;
          font-weight: 600;
          color: var(--text-muted);
          border-bottom: 3px solid transparent;
          transition: var(--transition);
          cursor: pointer;
        }
        .tab-btn:hover {
          color: var(--accent-hover);
          background: rgba(16, 185, 129, 0.03);
        }
        .tab-btn.active {
          color: var(--text-accent);
          border-bottom-color: var(--accent);
          font-weight: 700;
        }
        
        /* Permissions Modal Grid Styles */
        .permissions-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
          gap: 16px;
          padding: 8px 4px;
        }
        
        .resource-card {
          background: var(--bg-card);
          border: 1px solid var(--border-color);
          border-radius: var(--radius-md);
          box-shadow: var(--shadow-sm);
          overflow: hidden;
          transition: var(--transition);
        }
        
        .resource-card:hover {
          box-shadow: var(--shadow-md);
          border-color: var(--border-focus);
        }
        
        .resource-card-header {
          background: rgba(16, 185, 129, 0.04);
          padding: 10px 14px;
          font-weight: 700;
          color: var(--text-accent);
          text-transform: uppercase;
          font-size: 11px;
          letter-spacing: 0.5px;
          border-bottom: 1px solid var(--border-color);
        }
        
        .resource-card-body {
          padding: 12px 14px;
          display: flex;
          flex-direction: column;
          gap: 10px;
        }
        
        .perm-checkbox-label {
          display: flex;
          align-items: flex-start;
          gap: 8px;
          font-size: 12.5px;
          color: var(--text-primary);
          cursor: pointer;
          user-select: none;
        }
        
        .perm-checkbox-label input[type="checkbox"] {
          margin-top: 3px;
          width: 15px;
          height: 15px;
          accent-color: var(--accent);
          cursor: pointer;
        }
        
        .perm-name {
          font-weight: 500;
          line-height: 1.3;
        }
      `}</style>
    </div>
  );
}
