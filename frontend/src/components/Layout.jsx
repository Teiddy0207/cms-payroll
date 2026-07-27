import { useState } from 'react';
import { Sidebar } from './Sidebar.jsx';
import { Topbar } from './Topbar.jsx';
import { Outlet } from 'react-router-dom';

export function Layout() {
  const [mobileOpen, setMobileOpen] = useState(false);

  const toggleMobileMenu = () => {
    setMobileOpen((prev) => !prev);
  };

  const closeMobileMenu = () => {
    setMobileOpen(false);
  };

  return (
    <div className="app-shell">
      <Sidebar mobileOpen={mobileOpen} onClose={closeMobileMenu} />
      <div className="main-content">
        <Topbar onToggleMobileMenu={toggleMobileMenu} />
        <main className="page-content">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

export default Layout;

