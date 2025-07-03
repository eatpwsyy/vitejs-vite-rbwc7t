export interface User {
  id: number;
  email: string;
  role: 'admin' | 'hr' | 'manager' | 'employee';
  isActive: boolean;
  lastLogin?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Department {
  id: number;
  name: string;
  description: string;
  isActive: boolean;
  budget: number;
  goals: string;
  head?: Employee;
  employeeCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface Employee {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  dateOfBirth: string;
  gender: 'Male' | 'Female' | 'Other';
  address: string;
  employeeId: string;
  position: string;
  department: Department;
  hireDate: string;
  salary: number;
  status: 'Active' | 'Inactive' | 'Terminated';
  manager?: Employee;
  emergencyContact: string;
  emergencyContactPhone: string;
  createdAt: string;
  updatedAt: string;
}

export interface Leave {
  id: number;
  employee: Employee;
  type: string;
  startDate: string;
  endDate: string;
  days: number;
  reason: string;
  status: 'Pending' | 'Approved' | 'Rejected' | 'Cancelled';
  approver?: Employee;
  approvedAt?: string;
  comments: string;
  attachmentUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export interface LoginResponse {
  token: string;
  refreshToken: string;
  user: User;
  employee?: Employee;
}

export interface PaginationInfo {
  page: number;
  limit: number;
  total: number;
}

export interface ApiResponse<T> {
  data?: T;
  message?: string;
  error?: string;
}

export interface EmployeesResponse {
  employees: Employee[];
  pagination: PaginationInfo;
}

export interface DepartmentsResponse {
  departments: Department[];
}

export interface LeavesResponse {
  leaves: Leave[];
  pagination: PaginationInfo;
}

export interface EmployeeCreateRequest {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  dateOfBirth: string;
  gender: string;
  address: string;
  position: string;
  departmentId: number;
  hireDate: string;
  salary: number;
  managerId?: number;
  emergencyContact: string;
  emergencyContactPhone: string;
}

export interface EmployeeUpdateRequest {
  firstName?: string;
  lastName?: string;
  email?: string;
  phone?: string;
  dateOfBirth?: string;
  gender?: string;
  address?: string;
  position?: string;
  departmentId?: number;
  salary?: number;
  status?: string;
  managerId?: number;
  emergencyContact?: string;
  emergencyContactPhone?: string;
}

export interface DepartmentCreateRequest {
  name: string;
  description: string;
  headId?: number;
  budget: number;
  goals: string;
}

export interface DepartmentUpdateRequest {
  name?: string;
  description?: string;
  isActive?: boolean;
  headId?: number;
  budget?: number;
  goals?: string;
}

export interface LeaveCreateRequest {
  employeeId: number;
  type: string;
  startDate: string;
  endDate: string;
  reason: string;
  attachmentUrl?: string;
}

export interface LeaveUpdateRequest {
  type?: string;
  startDate?: string;
  endDate?: string;
  reason?: string;
  attachmentUrl?: string;
}

export interface LeaveApprovalRequest {
  status: 'Approved' | 'Rejected';
  comments?: string;
}

export const LEAVE_TYPES = [
  'Annual',
  'Sick',
  'Maternity',
  'Paternity',
  'Personal',
  'Emergency',
  'Bereavement',
  'Study',
  'Unpaid',
] as const;

export const LEAVE_STATUSES = [
  'Pending',
  'Approved',
  'Rejected',
  'Cancelled',
] as const;

export const USER_ROLES = [
  'admin',
  'hr',
  'manager',
  'employee',
] as const;

export const EMPLOYEE_STATUSES = [
  'Active',
  'Inactive',
  'Terminated',
] as const;