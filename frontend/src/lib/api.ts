import axios from 'axios';
import {
  User,
  Attendance,
  AttendanceStats,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  CheckInRequest,
  CheckOutRequest,
  UpdateProfileRequest,
  UpdateUserRequest,
  PaginationResponse,
} from '@/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/auth/login';
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: (data: LoginRequest): Promise<AuthResponse> =>
    api.post('/auth/login', data).then((res) => res.data),
  
  register: (data: RegisterRequest): Promise<AuthResponse> =>
    api.post('/auth/register', data).then((res) => res.data),
  
  getProfile: (): Promise<User> =>
    api.get('/profile').then((res) => res.data),
  
  updateProfile: (data: UpdateProfileRequest): Promise<User> =>
    api.put('/profile', data).then((res) => res.data),
};

// Attendance API
export const attendanceAPI = {
  checkIn: (data: CheckInRequest): Promise<Attendance> =>
    api.post('/attendance/checkin', data).then((res) => res.data),
  
  checkOut: (data: CheckOutRequest): Promise<Attendance> =>
    api.post('/attendance/checkout', data).then((res) => res.data),
  
  getTodayAttendance: (): Promise<Attendance | { status: string; date: string; check_in: null; check_out: null }> =>
    api.get('/attendance/today').then((res) => res.data),
  
  getAttendanceHistory: (page = 1, limit = 10): Promise<PaginationResponse<Attendance>> =>
    api.get(`/attendance?page=${page}&limit=${limit}`).then((res) => res.data),
  
  getAttendanceStats: (startDate?: string, endDate?: string): Promise<AttendanceStats> => {
    const params = new URLSearchParams();
    if (startDate) params.append('start_date', startDate);
    if (endDate) params.append('end_date', endDate);
    return api.get(`/attendance/stats?${params.toString()}`).then((res) => res.data);
  },
};

// Admin API
export const adminAPI = {
  getAllUsers: (page = 1, limit = 10): Promise<PaginationResponse<User>> =>
    api.get(`/admin/users?page=${page}&limit=${limit}`).then((res) => res.data),
  
  getAllAttendance: (page = 1, limit = 10, userId?: number, startDate?: string, endDate?: string): Promise<PaginationResponse<Attendance>> => {
    const params = new URLSearchParams();
    params.append('page', page.toString());
    params.append('limit', limit.toString());
    if (userId) params.append('user_id', userId.toString());
    if (startDate) params.append('start_date', startDate);
    if (endDate) params.append('end_date', endDate);
    return api.get(`/admin/attendance?${params.toString()}`).then((res) => res.data);
  },
  
  updateUser: (userId: number, data: UpdateUserRequest): Promise<User> =>
    api.put(`/admin/users/${userId}`, data).then((res) => res.data),
  
  deleteUser: (userId: number): Promise<{ message: string }> =>
    api.delete(`/admin/users/${userId}`).then((res) => res.data),
};

export default api;