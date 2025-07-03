export interface User {
  id: number;
  email: string;
  name: string;
  role: 'employee' | 'admin';
  position?: string;
  department?: string;
  created_at: string;
  updated_at: string;
}

export interface Attendance {
  id: number;
  user_id: number;
  user: User;
  date: string;
  check_in?: string;
  check_out?: string;
  notes: string;
  status: 'present' | 'absent' | 'late' | 'half_day' | 'not_checked_in';
  created_at: string;
  updated_at: string;
}

export interface AttendanceStats {
  total_days: number;
  present_days: number;
  absent_days: number;
  late_days: number;
  attendance_rate: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  name: string;
  password: string;
  position?: string;
  department?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface CheckInRequest {
  notes?: string;
}

export interface CheckOutRequest {
  notes?: string;
}

export interface UpdateProfileRequest {
  name?: string;
  position?: string;
  department?: string;
}

export interface UpdateUserRequest extends UpdateProfileRequest {
  role?: 'employee' | 'admin';
}

export interface PaginationResponse<T> {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export interface ApiError {
  error: string;
}