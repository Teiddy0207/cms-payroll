import { Navigate, Route, Routes, HashRouter } from 'react-router-dom';
import { ToastProvider } from './components/Toast.jsx';
import Layout from './components/Layout.jsx';
import LoginPage from './pages/LoginPage.jsx';
import DashboardPage from './pages/DashboardPage.jsx';
import EmployeesPage from './pages/EmployeesPage.jsx';
import PayrollPage from './pages/PayrollPage.jsx';
import PayrollRunPage from './pages/PayrollRunPage.jsx';
import JobPositionsPage from './pages/JobPositionsPage.jsx';
import CompetenciesPage from './pages/CompetenciesPage.jsx';
import DepartmentsPage from './pages/DepartmentsPage.jsx';
import ContractsPage from './pages/ContractsPage.jsx';
import SettingsPage from './pages/SettingsPage.jsx';
import FormulasPage from './pages/FormulasPage.jsx';
import FaceScanKiosk from './pages/FaceScanKiosk.jsx';
import TimesheetsPage from './pages/TimesheetsPage.jsx';
import AttendanceLogsPage from './pages/AttendanceLogsPage.jsx';
import ExplanationRequestsPage from './pages/ExplanationRequestsPage.jsx';
import OTRequestsPage from './pages/OTRequestsPage.jsx';

function ProtectedRoute({ children }) {
  const token = localStorage.getItem('auth_token');
  if (!token) {
    return <Navigate to="/login" replace />;
  }
  return children;
}

function App() {
  return (
    <ToastProvider>
      <HashRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <Layout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Navigate to="/dashboard" replace />} />
            <Route path="dashboard" element={<DashboardPage />} />
            <Route path="employees" element={<EmployeesPage />} />
            <Route path="payroll" element={<PayrollPage />} />
            <Route path="payroll-run" element={<PayrollRunPage />} />
            <Route path="payroll-formulas" element={<FormulasPage />} />
            <Route path="face-scan" element={<FaceScanKiosk />} />
            <Route path="timesheets" element={<TimesheetsPage />} />
            <Route path="attendance-logs" element={<AttendanceLogsPage />} />
            <Route path="explanation-requests" element={<ExplanationRequestsPage />} />
            <Route path="ot-requests" element={<OTRequestsPage />} />
            <Route path="contracts" element={<ContractsPage />} />
            <Route path="departments" element={<DepartmentsPage />} />
            <Route path="job-positions" element={<JobPositionsPage />} />
            <Route path="competencies" element={<CompetenciesPage />} />
            <Route path="settings" element={<SettingsPage />} />
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Route>
        </Routes>
      </HashRouter>
    </ToastProvider>
  );
}

export default App;
