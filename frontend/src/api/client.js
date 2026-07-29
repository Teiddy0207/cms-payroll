import axios from 'axios';

const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

const client = axios.create({
  baseURL: BASE_URL,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - attach Bearer token
client.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor - handle 401 (token expired/invalid) and 403 (permission denied)
client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // 401 = token invalid or expired → logout
      const isAlreadyOnLogin = window.location.hash.includes('/login');
      if (!isAlreadyOnLogin) {
        localStorage.removeItem('auth_token');
        window.location.href = '/#/login';
      }
    }
    // 403 = permission denied → do NOT logout, just reject the promise
    return Promise.reject(error);
  }
);


// ===================== AUTH =====================
export const authAPI = {
  login: (identifier, password) =>
    client.post('/public/auth/login', { identifiers: identifier, password }),
  logout: () =>
    client.post('/public/auth/logout'),
};

// ===================== DEPARTMENTS =====================
export const departmentsAPI = {
  list: (params = {}) =>
    client.get('/private/payroll/departments', { params: { page_number: 1, page_size: 20, ...params } }),
  create: (data) =>
    client.post('/private/payroll/departments', data),
  update: (id, data) =>
    client.put(`/private/payroll/departments/${id}`, data),
  delete: (id) =>
    client.delete(`/private/payroll/departments/${id}`),
};

// ===================== JOB POSITIONS =====================
export const jobPositionsAPI = {
  list: (params = {}) =>
    client.get('/private/payroll/job-positions', { params: { page_number: 1, page_size: 100, ...params } }),
  create: (data) =>
    client.post('/private/payroll/job-positions', data),
  update: (id, data) =>
    client.put(`/private/payroll/job-positions/${id}`, data),
  delete: (id) =>
    client.delete(`/private/payroll/job-positions/${id}`),
  getStandards: (id) =>
    client.get(`/private/payroll/job-positions/${id}/standards`),
  assignStandard: (id, data) =>
    client.post(`/private/payroll/job-positions/${id}/standards`, data),
  removeStandard: (id, standardId) =>
    client.delete(`/private/payroll/job-positions/${id}/standards/${standardId}`),
  calculateK: () =>
    client.post('/private/payroll/job-positions/calculate-k'),
  scrapeMarketSalary: (id, data) =>
    client.post(`/private/payroll/job-positions/${id || 'new'}/fetch-market-salary`, data),
};

// ===================== COMPETENCIES =====================
export const competenciesAPI = {
  list: (params = {}) =>
    client.get('/private/payroll/competencies', { params: { page_number: 1, page_size: 100, ...params } }),
  create: (data) =>
    client.post('/private/payroll/competencies', data),
  update: (id, data) =>
    client.put(`/private/payroll/competencies/${id}`, data),
  delete: (id) =>
    client.delete(`/private/payroll/competencies/${id}`),
};

// ===================== USER PROFILES (EMPLOYEES) =====================
export const employeesAPI = {
  list: (params = {}) =>
    client.get('/private/payroll/user-profiles', { params: { page_number: 1, page_size: 20, ...params } }),
  get: (id) =>
    client.get(`/private/payroll/user-profiles/${id}`),
  create: (data) =>
    client.post('/private/payroll/user-profiles', data),
  update: (id, data) =>
    client.put(`/private/payroll/user-profiles/${id}`, data),
  delete: (id) =>
    client.delete(`/private/payroll/user-profiles/${id}`),
  getCompetencies: (id) =>
    client.get(`/private/payroll/user-profiles/${id}/competencies`),
  assignCompetency: (id, data) =>
    client.post(`/private/payroll/user-profiles/${id}/competencies`, data),
  removeCompetency: (id, competencyId) =>
    client.delete(`/private/payroll/user-profiles/${id}/competencies/${competencyId}`),
};

// ===================== CONTRACTS =====================
export const contractsAPI = {
  list: (params = {}) =>
    client.get('/private/payroll/contracts', { params: { page_number: 1, page_size: 20, ...params } }),
  create: (data) =>
    client.post('/private/payroll/contracts', data),
  update: (id, data) =>
    client.put(`/private/payroll/contracts/${id}`, data),
  delete: (id) =>
    client.delete(`/private/payroll/contracts/${id}`),
};

// ===================== SETTINGS =====================
export const settingsAPI = {
  list: () =>
    client.get('/private/payroll/settings'),
  update: (key, data) =>
    client.put(`/private/payroll/settings/${key}`, data),
};

// ===================== CALCULATOR =====================
export const calculatorAPI = {
  preview: (employeeId, period) =>
    client.get(`/private/payroll/calculator/preview/${employeeId}`, { params: { period } }),
  runCalculation: (data) =>
    client.post('/private/payroll/calculator/run', data),
  getJobStatus: (jobId) =>
    client.get(`/private/payroll/calculator/job/${jobId}`),
  getSavedRecords: (period, params = {}) =>
    client.get('/private/payroll/calculator/records', { params: { period, ...params } }),
  getFormulas: () =>
    client.get('/private/payroll/calculator/formulas'),
  createFormula: (data) =>
    client.post('/private/payroll/calculator/formulas', data),
  updateFormula: (id, data) =>
    client.put(`/private/payroll/calculator/formulas/${id}`, data),
  deleteFormula: (id) =>
    client.delete(`/private/payroll/calculator/formulas/${id}`),
};

// ===================== TIMEKEEPING =====================
export const timekeepingAPI = {
  checkin: (data) =>
    client.post('/private/timekeeping/checkin', data),
  getSheets: (params) =>
    client.get('/private/timekeeping/sheets', { params }),
  getLogs: (params) =>
    client.get('/private/timekeeping/logs', { params }),
  calculateSheets: (data) =>
    client.post('/private/timekeeping/sheets/calculate', data),
  getExplanations: () =>
    client.get('/private/timekeeping/explanations'),
  createExplanation: (data) =>
    client.post('/private/timekeeping/explanations', data),
  updateExplanationStatus: (id, status) =>
    client.put(`/private/timekeeping/explanations/${id}/status`, { status }),
  getOTRequests: () =>
    client.get('/private/timekeeping/ot-requests'),
  createOTRequest: (data) =>
    client.post('/private/timekeeping/ot-requests', data),
  updateOTRequestStatus: (id, status) =>
    client.put(`/private/timekeeping/ot-requests/${id}/status`, { status }),
  registerFace: (data) =>
    client.post('/private/timekeeping/faces', data),
  getFaces: () =>
    client.get('/private/timekeeping/faces'),
  deleteFace: (code) =>
    client.delete(`/private/timekeeping/faces/${code}`),
};

// ===================== PDF TOOLS =====================
export const pdfAPI = {
  // Upload nhiều file, gộp thành 1 PDF
  merge: (files, outputName = 'merged_output') => {
    const form = new FormData();
    files.forEach(f => form.append('files', f));
    form.append('output_name', outputName);
    return client.post('/private/pdf/merge', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
      timeout: 120000,
    });
  },
  // Upload 1 file, nén PDF
  compress: (file) => {
    const form = new FormData();
    form.append('file', file);
    return client.post('/private/pdf/compress', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
      timeout: 120000,
    });
  },
  // Upload 1 file, thêm watermark text
  watermark: (file, text = 'BẢO MẬT') => {
    const form = new FormData();
    form.append('file', file);
    form.append('text', text);
    return client.post('/private/pdf/watermark', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
      timeout: 120000,
    });
  },
  // Upload 1 file, xoay trang
  rotate: (file, angle = 90, pages = '') => {
    const form = new FormData();
    form.append('file', file);
    form.append('angle', String(angle));
    form.append('pages', pages);
    return client.post('/private/pdf/rotate', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
      timeout: 120000,
    });
  },
  // Download file kết quả từ server với token query param
  downloadUrl: (path) => {
    const token = localStorage.getItem('auth_token') || '';
    return `${BASE_URL}/private/pdf/download?path=${encodeURIComponent(path)}&token=${encodeURIComponent(token)}`;
  },
  // Download file qua client Axios với Header Authorization
  downloadFile: (path) =>
    client.get('/private/pdf/download', {
      params: { path },
      responseType: 'blob',
    }),
};

// ===================== MEETINGS =====================
export const meetingsAPI = {
  list: () =>
    client.get('/meetings'),
  get: (id) =>
    client.get(`/meetings/${id}`),
  create: (data) =>
    client.post('/meetings', data),
  update: (id, data) =>
    client.put(`/meetings/${id}`, data),
  rsvp: (id, status, note = '') =>
    client.post(`/meetings/${id}/rsvp`, { status: status, note: note }),
  signal: (data) =>
    client.post('/meetings/signal', data),
  getSummary: (id) =>
    client.get(`/meetings/${id}/summary`),
  saveSummary: (id, data) =>
    client.post(`/meetings/${id}/summary`, data),
};

// ===================== ROLES & PERMISSIONS =====================
export const authAdminAPI = {
  // Roles CRUD
  listRoles: (params = {}) =>
    client.get('/private/auth/roles', { params: { page_number: 1, page_size: 100, ...params } }),
  createRole: (data) =>
    client.post('/private/auth/roles', data),
  updateRole: (id, data) =>
    client.put(`/private/auth/roles/${id}`, data),
  deleteRole: (id) =>
    client.delete(`/private/auth/roles/${id}`),

  // Permissions CRUD
  listPermissions: (params = {}) =>
    client.get('/private/auth/permissions', { params: { page_number: 1, page_size: 200, ...params } }),
  createPermission: (data) =>
    client.post('/private/auth/permissions', data),
  updatePermission: (id, data) =>
    client.put(`/private/auth/permissions/${id}`, data),
  deletePermission: (id) =>
    client.delete(`/private/auth/permissions/${id}`),

  // Role - Permission Assignment
  assignPermissionToRole: (data) =>
    client.post('/private/auth/roles/assign-permission', data),

  // User - Role - Permission Assignment
  assignRoleToUser: (data) =>
    client.post('/private/auth/users/assign-role', data),
  assignPermissionToUser: (data) =>
    client.post('/private/auth/users/assign-permission', data),
  getUserPermissions: (userId) =>
    client.get(`/private/auth/users/${userId}/permissions`),

  // User CRUD
  listUsers: (params = {}) =>
    client.get('/private/auth/users', { params: { page_number: 1, page_size: 100, ...params } }),
  createUser: (data) =>
    client.post('/private/auth/users', data),
  updateUser: (id, data) =>
    client.put(`/private/auth/users/${id}`, data),
  deleteUser: (id) =>
    client.delete(`/private/auth/users/${id}`),
  getMePermissions: () =>
    client.get('/private/auth/me/permissions'),
};

export default client;
