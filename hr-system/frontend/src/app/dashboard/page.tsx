'use client';

import { useQuery } from '@tanstack/react-query';
import { employeeAPI, departmentAPI, leaveAPI } from '@/lib/api';
import Layout from '@/components/Layout';
import ProtectedRoute from '@/components/ProtectedRoute';
import { Users, Building2, Calendar, TrendingUp } from 'lucide-react';

interface StatsCardProps {
  title: string;
  value: string | number;
  icon: React.ComponentType<any>;
  color: string;
  loading?: boolean;
}

function StatsCard({ title, value, icon: Icon, color, loading }: StatsCardProps) {
  return (
    <div className="bg-white overflow-hidden shadow rounded-lg">
      <div className="p-5">
        <div className="flex items-center">
          <div className="flex-shrink-0">
            <Icon className={`h-6 w-6 ${color}`} />
          </div>
          <div className="ml-5 w-0 flex-1">
            <dl>
              <dt className="text-sm font-medium text-gray-500 truncate">{title}</dt>
              <dd className="text-lg font-medium text-gray-900">
                {loading ? (
                  <div className="animate-pulse h-6 bg-gray-200 rounded w-16"></div>
                ) : (
                  value
                )}
              </dd>
            </dl>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function DashboardPage() {
  const { data: employeesData, isLoading: employeesLoading } = useQuery({
    queryKey: ['employees', { limit: 1 }],
    queryFn: () => employeeAPI.getEmployees({ limit: 1, page: 1 }),
  });

  const { data: departmentsData, isLoading: departmentsLoading } = useQuery({
    queryKey: ['departments'],
    queryFn: () => departmentAPI.getDepartments(),
  });

  const { data: leavesData, isLoading: leavesLoading } = useQuery({
    queryKey: ['leaves', { limit: 1, status: 'Pending' }],
    queryFn: () => leaveAPI.getLeaves({ limit: 1, page: 1, status: 'Pending' }),
  });

  const { data: activeEmployeesData, isLoading: activeEmployeesLoading } = useQuery({
    queryKey: ['employees', { limit: 1, status: 'Active' }],
    queryFn: () => employeeAPI.getEmployees({ limit: 1, page: 1, status: 'Active' }),
  });

  const totalEmployees = employeesData?.data?.pagination?.total || 0;
  const totalDepartments = departmentsData?.data?.departments?.length || 0;
  const pendingLeaves = leavesData?.data?.pagination?.total || 0;
  const activeEmployees = activeEmployeesData?.data?.pagination?.total || 0;

  return (
    <ProtectedRoute>
      <Layout>
        <div className="space-y-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
            <p className="mt-1 text-sm text-gray-500">
              Welcome to the HR Management System
            </p>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
            <StatsCard
              title="Total Employees"
              value={totalEmployees}
              icon={Users}
              color="text-blue-600"
              loading={employeesLoading}
            />
            <StatsCard
              title="Active Employees"
              value={activeEmployees}
              icon={TrendingUp}
              color="text-green-600"
              loading={activeEmployeesLoading}
            />
            <StatsCard
              title="Departments"
              value={totalDepartments}
              icon={Building2}
              color="text-purple-600"
              loading={departmentsLoading}
            />
            <StatsCard
              title="Pending Leaves"
              value={pendingLeaves}
              icon={Calendar}
              color="text-orange-600"
              loading={leavesLoading}
            />
          </div>

          {/* Quick Actions */}
          <div className="bg-white shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <h3 className="text-lg leading-6 font-medium text-gray-900">
                Quick Actions
              </h3>
              <div className="mt-5 grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
                <a
                  href="/employees/new"
                  className="inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
                >
                  Add Employee
                </a>
                <a
                  href="/departments/new"
                  className="inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-purple-600 hover:bg-purple-700"
                >
                  Add Department
                </a>
                <a
                  href="/leaves/new"
                  className="inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700"
                >
                  Request Leave
                </a>
                <a
                  href="/leaves?status=Pending"
                  className="inline-flex items-center justify-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
                >
                  Review Leaves
                </a>
              </div>
            </div>
          </div>

          {/* Recent Activity */}
          <div className="bg-white shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <h3 className="text-lg leading-6 font-medium text-gray-900">
                Recent Activity
              </h3>
              <div className="mt-5">
                <div className="text-sm text-gray-500">
                  No recent activity to display.
                </div>
              </div>
            </div>
          </div>
        </div>
      </Layout>
    </ProtectedRoute>
  );
}