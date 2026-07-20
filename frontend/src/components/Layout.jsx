import { Sidebar } from './Sidebar.jsx';
import { Topbar } from './Topbar.jsx';
import { Outlet } from 'react-router-dom';

export function Layout() {
  return (
    <div className="app-shell">
      <Sidebar />
      <div className="main-content">
        <Topbar />
        <main className="page-content">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

export default Layout;
