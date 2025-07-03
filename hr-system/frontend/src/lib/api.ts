import axios, { AxiosResponse } from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('refresh_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: (credentials: { email: string; password: string }) =>
    api.post('/auth/login', credentials),
  
  register: (userData: {
    email: string;
    password: string;
    role?: string;
    employeeId?: number;
  }) => api.post('/auth/register', userData),
  
  changePassword: (data: {
    currentPassword: string;
    newPassword: string;
  }) => api.post('/auth/change-password', data),
  
  getProfile: () => api.get('/auth/profile'),
};

// Employee API
export const employeeAPI = {
  getEmployees: (params?: {
    page?: number;
    limit?: number;
    department?: string;
    status?: string;
    search?: string;
  }) => api.get('/employees', { params }),
  
  getEmployee: (id: string) => api.get(`/employees/${id}`),
  
  createEmployee: (data: any) => api.post('/employees', data),
  
  updateEmployee: (id: string, data: any) => api.put(`/employees/${id}`, data),
  
  deleteEmployee: (id: string) => api.delete(`/employees/${id}`),
};

// Department API
export const departmentAPI = {
  getDepartments: (params?: {
    search?: string;
    active?: boolean;
  }) => api.get('/departments', { params }),
  
  getDepartment: (id: string) => api.get(`/departments/${id}`),
  
  createDepartment: (data: any) => api.post('/departments', data),
  
  updateDepartment: (id: string, data: any) => api.put(`/departments/${id}`, data),
  
  deleteDepartment: (id: string) => api.delete(`/departments/${id}`),
};

// Leave API
export const leaveAPI = {
  getLeaves: (params?: {
    page?: number;
    limit?: number;
    employee?: string;
    status?: string;
    type?: string;
    startDate?: string;
    endDate?: string;
  }) => api.get('/leaves', { params }),
  
  getLeave: (id: string) => api.get(`/leaves/${id}`),
  
  createLeave: (data: any) => api.post('/leaves', data),
  
  updateLeave: (id: string, data: any) => api.put(`/leaves/${id}`, data),
  
  approveLeave: (id: string, data: {
    status: 'Approved' | 'Rejected';
    comments?: string;
  }) => api.post(`/leaves/${id}/approve`, data),
  
  deleteLeave: (id: string) => api.delete(`/leaves/${id}`),
  
  getLeaveTypes: () => api.get('/leaves/types'),
  
  getLeaveStatuses: () => api.get('/leaves/statuses'),
};

export default api;