import axios from 'axios';

const BASE_URL = 'http://localhost:7070/api/v1';

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

// Response interceptor - handle 401
client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('auth_token');
      window.location.href = '/login';
    }
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
};

// ===================== JOB STANDARDS =====================
export const jobStandardsAPI = {
  list: (params = {}) =>
    client.get('/private/payroll/job-standards', { params: { page_number: 1, page_size: 20, ...params } }),
  create: (data) =>
    client.post('/private/payroll/job-standards', data),
  update: (id, data) =>
    client.put(`/private/payroll/job-standards/${id}`, data),
  delete: (id) =>
    client.delete(`/private/payroll/job-standards/${id}`),
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
};

export default client;
