# HR Management System

A comprehensive Human Resources Management System built with Go (Golang) backend and Next.js 15 frontend.

## Features

### Backend (Go)
- **Authentication & Authorization**: JWT-based authentication with role-based access control
- **Employee Management**: Full CRUD operations for employee records
- **Department Management**: Department creation and management
- **Leave Management**: Leave request, approval, and tracking system
- **RESTful API**: Clean REST API with proper error handling
- **Database**: SQLite with GORM ORM (easily configurable for PostgreSQL/MySQL)

### Frontend (Next.js 15)
- **Modern UI**: Clean and responsive design with Tailwind CSS
- **Authentication**: Secure login system with session management
- **Dashboard**: Overview with statistics and quick actions
- **Employee Management**: Add, edit, view employee information
- **Department Management**: Manage departments and assignments
- **Leave Management**: Request, approve, and track leaves
- **Role-based Access**: Different views and permissions based on user roles

## User Roles

1. **Admin**: Full access to all features
2. **HR**: Manage employees, departments, and approve leaves
3. **Manager**: View employees, approve team leaves
4. **Employee**: View own information, request leaves

## Tech Stack

### Backend
- Go 1.21+
- Gin Web Framework
- GORM (ORM)
- SQLite (database)
- JWT for authentication
- bcrypt for password hashing

### Frontend
- Next.js 15 (React 19)
- TypeScript
- Tailwind CSS
- TanStack Query (React Query)
- Axios for HTTP requests
- Lucide React (icons)

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Node.js 18 or higher
- npm or yarn

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd hr-system/backend
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the backend server:
   ```bash
   go run main.go
   ```

The backend will start on `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd hr-system/frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

The frontend will start on `http://localhost:3000`

## Default Login Credentials

- **Email**: admin@hr.com
- **Password**: password

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/change-password` - Change password
- `GET /api/v1/auth/profile` - Get user profile

### Employees
- `GET /api/v1/employees` - List employees (with pagination and filters)
- `GET /api/v1/employees/:id` - Get employee details
- `POST /api/v1/employees` - Create employee (HR/Admin only)
- `PUT /api/v1/employees/:id` - Update employee (HR/Admin only)
- `DELETE /api/v1/employees/:id` - Delete employee (HR/Admin only)

### Departments
- `GET /api/v1/departments` - List departments
- `GET /api/v1/departments/:id` - Get department details
- `POST /api/v1/departments` - Create department (HR/Admin only)
- `PUT /api/v1/departments/:id` - Update department (HR/Admin only)
- `DELETE /api/v1/departments/:id` - Delete department (HR/Admin only)

### Leaves
- `GET /api/v1/leaves` - List leaves (with filters)
- `GET /api/v1/leaves/:id` - Get leave details
- `POST /api/v1/leaves` - Create leave request
- `PUT /api/v1/leaves/:id` - Update leave request
- `POST /api/v1/leaves/:id/approve` - Approve/reject leave (Manager/HR/Admin)
- `DELETE /api/v1/leaves/:id` - Delete leave request
- `GET /api/v1/leaves/types` - Get leave types
- `GET /api/v1/leaves/statuses` - Get leave statuses

## Environment Variables

### Backend
Create a `.env` file in the backend directory:
```env
PORT=8080
DB_PATH=hr_system.db
JWT_SECRET=your-secret-key-change-in-production
```

### Frontend
Create a `.env.local` file in the frontend directory:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## Database Schema

The system automatically creates the following tables:
- `users` - User authentication and roles
- `employees` - Employee information
- `departments` - Department details
- `leaves` - Leave requests and approvals

## Development

### Running Backend Tests
```bash
cd backend
go test ./...
```

### Running Frontend Tests
```bash
cd frontend
npm test
```

### Building for Production

#### Backend
```bash
cd backend
go build -o hr-backend main.go
```

#### Frontend
```bash
cd frontend
npm run build
npm start
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support or questions, please open an issue on the repository.