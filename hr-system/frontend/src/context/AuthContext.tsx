'use client';

import React, { createContext, useContext, useEffect, useState } from 'react';
import { User, Employee, LoginResponse } from '@/types';
import { authAPI } from '@/lib/api';

interface AuthContextType {
  user: User | null;
  employee: Employee | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  refreshProfile: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

interface AuthProviderProps {
  children: React.ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [employee, setEmployee] = useState<Employee | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const isAuthenticated = !!user;

  const login = async (email: string, password: string) => {
    try {
      const response = await authAPI.login({ email, password });
      const data: LoginResponse = response.data;
      
      // Store tokens
      localStorage.setItem('auth_token', data.token);
      localStorage.setItem('refresh_token', data.refreshToken);
      
      // Set user and employee data
      setUser(data.user);
      setEmployee(data.employee || null);
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Login failed');
    }
  };

  const logout = () => {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('refresh_token');
    setUser(null);
    setEmployee(null);
  };

  const refreshProfile = async () => {
    try {
      const response = await authAPI.getProfile();
      setUser(response.data.user);
      setEmployee(response.data.employee || null);
    } catch (error) {
      logout();
    }
  };

  useEffect(() => {
    const initAuth = async () => {
      const token = localStorage.getItem('auth_token');
      if (token) {
        try {
          await refreshProfile();
        } catch (error) {
          logout();
        }
      }
      setIsLoading(false);
    };

    initAuth();
  }, []);

  const value: AuthContextType = {
    user,
    employee,
    isAuthenticated,
    isLoading,
    login,
    logout,
    refreshProfile,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}