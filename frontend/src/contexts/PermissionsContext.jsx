import { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { authAdminAPI } from '../api/client.js';

const PermissionsContext = createContext({ slugs: [], loaded: false, can: () => true, refetch: () => {} });

export function PermissionsProvider({ children }) {
  const [slugs, setSlugs] = useState([]);
  const [loaded, setLoaded] = useState(false);

  const fetchPermissions = useCallback(async () => {
    const token = localStorage.getItem('auth_token');
    if (!token) {
      setSlugs([]);
      setLoaded(true);
      return;
    }
    try {
      const res = await authAdminAPI.getMePermissions();
      const perms = res.data?.data || [];
      setSlugs(perms.map((p) => p.slug));
    } catch {
      setSlugs([]);
    } finally {
      setLoaded(true);
    }
  }, []);

  useEffect(() => {
    fetchPermissions();

    // Re-fetch permissions if token changes (e.g., after login in another tab)
    const handleStorageChange = (e) => {
      if (e.key === 'auth_token') {
        fetchPermissions();
      }
    };
    window.addEventListener('storage', handleStorageChange);
    return () => window.removeEventListener('storage', handleStorageChange);
  }, [fetchPermissions]);

  const can = (slug) => {
    if (!slug) return true;
    return slugs.includes(slug);
  };

  return (
    <PermissionsContext.Provider value={{ slugs, loaded, can, refetch: fetchPermissions }}>
      {children}
    </PermissionsContext.Provider>
  );
}

export function usePermissions() {
  return useContext(PermissionsContext);
}

export default PermissionsContext;
